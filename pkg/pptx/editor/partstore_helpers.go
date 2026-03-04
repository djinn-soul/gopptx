package editor

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
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
