package editor

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"sort"
	"sync"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
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

// newPartStoreFromZip creates a lazy store backed by the given zip reader.
// The [os.File] must remain open for the lifetime of the store.
func newPartStoreFromZip(file *os.File, zr *zip.Reader) *PartStore {
	allKeys := make(map[string]struct{}, len(zr.File))
	index := make(map[string]*zip.File, len(zr.File))
	for _, entry := range zr.File {
		if entry.FileInfo().IsDir() {
			continue
		}
		name := common.CanonicalPartPath(entry.Name)
		index[name] = entry
		allKeys[name] = struct{}{}
	}

	return &PartStore{
		file:      file,
		index:     index,
		cache:     make(map[string][]byte),
		modified:  make(map[string][]byte),
		deleted:   make(map[string]bool),
		keysDirty: true,
		inflight:  make(map[string]*inflightRead),
		allKeys:   allKeys,
	}
}

// NewPartStore creates a new, empty in-memory part store.
// NewPartStore creates a new, empty in-memory part store.
func NewPartStore() *PartStore {
	return newPartStoreFromMap(nil)
}

// newPartStoreFromMap creates an in-memory store from a pre-loaded map.
// Used by tests and by newPresentationEditorFromParts when parts are
// already loaded (e.g. from MergeFromEditor).
func newPartStoreFromMap(parts map[string][]byte) *PartStore {
	cached := make(map[string][]byte, len(parts))
	maps.Copy(cached, parts)
	allKeys := make(map[string]struct{}, len(parts))
	for k := range parts {
		allKeys[k] = struct{}{}
	}
	return &PartStore{
		index:     make(map[string]*zip.File),
		cache:     cached,
		modified:  make(map[string][]byte),
		deleted:   make(map[string]bool),
		keysDirty: true,
		inflight:  make(map[string]*inflightRead),
		allKeys:   allKeys,
	}
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
func (ps *PartStore) KeysWithPrefix(prefix string) []string {
	all := ps.Keys()
	out := make([]string, 0)
	for _, name := range all {
		if len(name) >= len(prefix) && name[:len(prefix)] == prefix {
			out = append(out, name)
		}
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

func readZipEntry(entry *zip.File) ([]byte, error) {
	rc, err := entry.Open()
	if err != nil {
		return nil, fmt.Errorf("open zip entry %q: %w", entry.Name, err)
	}
	defer func() { _ = rc.Close() }()

	const maxInt = int(^uint(0) >> 1)
	size64 := entry.UncompressedSize64
	if size64 > 0 && size64 <= uint64(maxInt) {
		// Pre-allocate the exact size when available.
		data := make([]byte, int(size64))
		if _, err := io.ReadFull(rc, data); err != nil {
			return nil, fmt.Errorf("read zip entry %q: %w", entry.Name, err)
		}
		return data, nil
	}

	// Fallback for unknown sizes: use a pooled read buffer to reduce temporary allocations.
	chunkPtr, ok := zipReadChunkPool.Get().(*[]byte)
	if !ok || chunkPtr == nil {
		return nil, errors.New("zip read pool returned invalid buffer")
	}
	defer zipReadChunkPool.Put(chunkPtr)
	chunk := *chunkPtr
	var out bytes.Buffer
	limitedReader := io.LimitReader(rc, maxUnknownZipEntryBytes+1)
	if _, err := io.CopyBuffer(&out, limitedReader, chunk); err != nil {
		return nil, fmt.Errorf("read zip entry %q: %w", entry.Name, err)
	}
	if out.Len() > maxUnknownZipEntryBytes {
		return nil, fmt.Errorf("zip entry %q exceeds max size %d bytes", entry.Name, maxUnknownZipEntryBytes)
	}
	return out.Bytes(), nil
}

// Materialize eagerly reads all remaining lazy zip entries into the cache
// and then closes the underlying file handle. After this call the store
// is fully in-memory and the source file is no longer locked.
func (ps *PartStore) Materialize() error {
	ps.mu.RLock()
	toLoad := make([]string, 0, len(ps.index))
	for name := range ps.index {
		if ps.deleted[name] {
			continue
		}
		if _, ok := ps.cache[name]; ok {
			continue
		}
		if _, ok := ps.modified[name]; ok {
			continue
		}
		toLoad = append(toLoad, name)
	}
	ps.mu.RUnlock()

	for _, name := range toLoad {
		if _, ok := ps.Get(name); !ok {
			return fmt.Errorf("materialize part %q: read failed", name)
		}
	}

	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.index = make(map[string]*zip.File)
	ps.invalidateKeysLocked()
	if ps.file != nil {
		err := ps.file.Close()
		ps.file = nil
		return err
	}
	return nil
}
