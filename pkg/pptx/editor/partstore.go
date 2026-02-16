package editor

import (
	"archive/zip"
	"fmt"
	"maps"
	"os"
	"sort"
	"sync"

	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

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
}

// newPartStoreFromZip creates a lazy store backed by the given zip reader.
// The os.File must remain open for the lifetime of the store.
func newPartStoreFromZip(file *os.File, zr *zip.Reader) *PartStore {
	index := make(map[string]*zip.File, len(zr.File))
	for _, entry := range zr.File {
		if entry.FileInfo().IsDir() {
			continue
		}
		index[common.CanonicalPartPath(entry.Name)] = entry
	}
	return &PartStore{
		file:      file,
		index:     index,
		cache:     make(map[string][]byte),
		modified:  make(map[string][]byte),
		deleted:   make(map[string]bool),
		keysDirty: true,
		inflight:  make(map[string]*inflightRead),
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
	return &PartStore{
		index:     make(map[string]*zip.File),
		cache:     cached,
		modified:  make(map[string][]byte),
		deleted:   make(map[string]bool),
		keysDirty: true,
		inflight:  make(map[string]*inflightRead),
	}
}

// Get returns the data for the named part. Modified data takes priority,
// then cached data, then lazy-reads from the zip archive.
func (ps *PartStore) Get(name string) ([]byte, bool) {
	ps.mu.RLock()
	if ps.deleted[name] {
		ps.mu.RUnlock()
		return nil, false
	}
	if data, ok := ps.modified[name]; ok {
		ps.mu.RUnlock()
		return data, true
	}
	if data, ok := ps.cache[name]; ok {
		ps.mu.RUnlock()
		return data, true
	}
	if pending, ok := ps.inflight[name]; ok {
		p := pending
		ch := p.ch
		ps.mu.RUnlock()
		<-ch
		if p.err != nil {
			return nil, false
		}
		return p.data, true
	}
	ps.mu.RUnlock()

	ps.mu.Lock()
	if ps.deleted[name] {
		ps.mu.Unlock()
		return nil, false
	}
	if data, ok := ps.modified[name]; ok {
		ps.mu.Unlock()
		return data, true
	}
	if data, ok := ps.cache[name]; ok {
		ps.mu.Unlock()
		return data, true
	}
	if pending, ok := ps.inflight[name]; ok {
		p := pending
		ch := p.ch
		ps.mu.Unlock()
		<-ch
		if p.err != nil {
			return nil, false
		}
		return p.data, true
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
	ps.invalidateKeysLocked()
}

// Delete removes a part from the store.
func (ps *PartStore) Delete(name string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	delete(ps.modified, name)
	delete(ps.cache, name)
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

	seen := make(map[string]struct{})
	for name := range ps.index {
		if !ps.deleted[name] {
			seen[name] = struct{}{}
		}
	}
	for name := range ps.cache {
		if !ps.deleted[name] {
			seen[name] = struct{}{}
		}
	}
	for name := range ps.modified {
		seen[name] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for name := range seen {
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
