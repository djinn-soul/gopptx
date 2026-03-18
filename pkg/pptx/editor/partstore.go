package editor

import (
	"archive/zip"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
)

const (
	maxUnknownZipEntryBytes = 256 * 1024 * 1024
	zipReadChunkBytes       = 32 * 1024
)

//nolint:gochecknoglobals // Shared reusable buffers reduce allocations in zip reads.
var zipReadChunkPool = sync.Pool{
	New: func() any {
		b := make([]byte, zipReadChunkBytes)
		return &b
	},
}

type inflightRead struct {
	ch   chan struct{}
	data []byte
	err  error
}

// PartStore provides lazy-loading access to PPTX package parts.
// Parts from the original zip are read on demand; mutations are tracked
// separately so untouched parts never leave the zip archive.
type PartStore struct {
	mu   sync.RWMutex
	file *os.File // kept open for lazy reads; nil for in-memory stores

	// index maps canonical part names to zip entries (data not yet read).
	index map[string]*zip.File

	// cache holds data loaded from zip (lazy-populated by Get).
	cache map[string][]byte

	// modified holds parts written or replaced after opening.
	modified map[string][]byte

	// deleted tracks parts that have been removed.
	deleted map[string]bool

	// keysCache stores sorted keys and is invalidated by mutations.
	keysCache []string
	keysDirty bool

	// inflight deduplicates concurrent lazy reads for the same part.
	inflight map[string]*inflightRead

	// allKeys maintains the set of all active part names to avoid O(N) map scans in Keys().
	allKeys map[string]struct{}
}

// Get returns the data for the named part. Modified data takes priority,
// then cached data, then lazy-reads from the zip archive.
func (ps *PartStore) Get(name string) ([]byte, bool) {
	ps.mu.RLock()
	if data, ok, pending := ps.getPriorityDataLocked(name); ok {
		ps.mu.RUnlock()
		return data, true
	} else if pending != nil {
		ps.mu.RUnlock()
		return waitInflightRead(pending)
	}
	ps.mu.RUnlock()

	ps.mu.Lock()
	if data, ok, pending := ps.getPriorityDataLocked(name); ok {
		ps.mu.Unlock()
		return data, true
	} else if pending != nil {
		ps.mu.Unlock()
		return waitInflightRead(pending)
	}
	entry, ok := ps.index[name]
	if !ok {
		ps.mu.Unlock()
		return nil, false
	}

	pending := &inflightRead{ch: make(chan struct{})}
	ps.inflight[name] = pending
	ps.mu.Unlock()

	data, err := readZipEntry(entry)

	ps.mu.Lock()
	delete(ps.inflight, name)
	if err == nil {
		if ps.deleted[name] {
			err = fmt.Errorf("part %q was deleted during read", name)
		} else if modifiedData, modifiedOK := ps.modified[name]; modifiedOK {
			data = modifiedData
		} else if cachedData, cachedOK := ps.cache[name]; cachedOK {
			data = cachedData
		} else {
			ps.cache[name] = data
		}
	}
	pending.data = data
	pending.err = err
	close(pending.ch)
	ps.mu.Unlock()

	if err != nil {
		return nil, false
	}
	return data, true
}

// Set writes or replaces a part's data.
func (ps *PartStore) Set(name string, data []byte) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	delete(ps.deleted, name)
	ps.modified[name] = data
	ps.allKeys[name] = struct{}{}
	ps.invalidateKeysLocked()
}

// Delete removes a part from the store.
func (ps *PartStore) Delete(name string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	delete(ps.modified, name)
	delete(ps.cache, name)
	delete(ps.allKeys, name)
	if _, ok := ps.index[name]; ok {
		ps.deleted[name] = true
	} else {
		delete(ps.deleted, name)
	}
	ps.invalidateKeysLocked()
}

// Has checks whether a part exists without loading its data.
func (ps *PartStore) Has(name string) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	if ps.deleted[name] {
		return false
	}
	if _, ok := ps.modified[name]; ok {
		return true
	}
	if _, ok := ps.cache[name]; ok {
		return true
	}
	_, ok := ps.index[name]
	return ok
}

// Keys returns all part names in sorted order.
func (ps *PartStore) Keys() []string {
	ps.mu.RLock()
	if !ps.keysDirty {
		keys := append([]string(nil), ps.keysCache...)
		ps.mu.RUnlock()
		return keys
	}
	ps.mu.RUnlock()

	ps.mu.Lock()
	defer ps.mu.Unlock()
	if !ps.keysDirty {
		return append([]string(nil), ps.keysCache...)
	}

	out := make([]string, 0, len(ps.allKeys))
	for name := range ps.allKeys {
		out = append(out, name)
	}
	sort.Strings(out)
	ps.keysCache = out
	ps.keysDirty = false
	return append([]string(nil), out...)
}

// Close releases the underlying file handle, if any.
func (ps *PartStore) Close() error {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if ps.file != nil {
		err := ps.file.Close()
		ps.file = nil
		return err
	}
	return nil
}

func waitInflightRead(pending *inflightRead) ([]byte, bool) {
	<-pending.ch
	if pending.err != nil {
		return nil, false
	}
	return pending.data, true
}

// getPriorityDataLocked checks deleted/modified/cache/inflight in priority order.
// Callers must hold at least a read lock on ps.mu.
func (ps *PartStore) getPriorityDataLocked(name string) ([]byte, bool, *inflightRead) {
	if ps.deleted[name] {
		return nil, false, nil
	}
	if data, ok := ps.modified[name]; ok {
		return data, true, nil
	}
	if data, ok := ps.cache[name]; ok {
		return data, true, nil
	}
	if pending, ok := ps.inflight[name]; ok {
		return nil, false, pending
	}
	return nil, false, nil
}

// KeysWithPrefix returns part names that start with the given prefix.
// It operates directly on the sorted keysCache — no full-copy allocation.
// Binary search locates the first candidate; linear scan stops at the first non-match.
// Fast path (clean cache) uses RLock so concurrent readers are not blocked.
func (ps *PartStore) KeysWithPrefix(prefix string) []string {
	// Fast path: cache is clean, use shared read lock.
	ps.mu.RLock()
	if !ps.keysDirty {
		result := prefixSearchSorted(ps.keysCache, prefix)
		ps.mu.RUnlock()
		return result
	}
	ps.mu.RUnlock()

	// Slow path: rebuild cache under exclusive write lock, then search.
	ps.mu.Lock()
	if ps.keysDirty {
		out := make([]string, 0, len(ps.allKeys))
		for name := range ps.allKeys {
			out = append(out, name)
		}
		sort.Strings(out)
		ps.keysCache = out
		ps.keysDirty = false
	}
	result := prefixSearchSorted(ps.keysCache, prefix)
	ps.mu.Unlock()
	return result
}

func prefixSearchSorted(cache []string, prefix string) []string {
	start := sort.SearchStrings(cache, prefix)
	out := make([]string, 0, 8)
	for i := start; i < len(cache); i++ {
		if !strings.HasPrefix(cache[i], prefix) {
			break
		}
		out = append(out, cache[i])
	}
	return out
}

// Snapshot returns a deep-copy map of all parts (forces all lazy reads).
// Used for merge operations where the source store may be closed.
func (ps *PartStore) Snapshot() map[string][]byte {
	keys := ps.Keys()
	out := make(map[string][]byte, len(keys))
	for _, k := range keys {
		data, ok := ps.Get(k)
		if ok {
			clone := make([]byte, len(data))
			copy(clone, data)
			out[k] = clone
		}
	}
	return out
}

func (ps *PartStore) invalidateKeysLocked() {
	ps.keysDirty = true
	ps.keysCache = nil
}
