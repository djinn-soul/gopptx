package editor

import (
	"archive/zip"
	"fmt"
)

// Materialize eagerly reads all remaining lazy zip entries into the cache
// and then closes the underlying file handle. After this call the store
// is fully in-memory and the source file is no longer locked.
func (ps *PartStore) Materialize() error {
	// Snapshot only parts that still need lazy loading.
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
	// Index is no longer needed — parts are all in cache.
	ps.index = make(map[string]*zip.File)
	ps.invalidateKeysLocked()
	if ps.file != nil {
		err := ps.file.Close()
		ps.file = nil
		return err
	}
	return nil
}
