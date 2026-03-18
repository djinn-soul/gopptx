package editor

import (
	"archive/zip"
	"path/filepath"
	"runtime"
	"strings"
)

// SourceZipEntry returns the original zip entry for a part when it is still
// unmodified/deleted in the store. Callers can raw-copy it into a new archive.
func (ps *PartStore) SourceZipEntry(name string) (*zip.File, bool) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.deleted[name] {
		return nil, false
	}
	if _, modified := ps.modified[name]; modified {
		return nil, false
	}
	entry, ok := ps.index[name]
	if !ok {
		return nil, false
	}
	return entry, true
}

// IsBackedByPath reports whether this store is currently backed by the same
// on-disk path (best-effort normalized compare).
func (ps *PartStore) IsBackedByPath(targetPath string) bool {
	ps.mu.RLock()
	file := ps.file
	ps.mu.RUnlock()
	if file == nil {
		return false
	}

	sourcePath := file.Name()
	if filepath.IsAbs(sourcePath) && filepath.IsAbs(targetPath) {
		sourcePath = filepath.Clean(sourcePath)
		targetPath = filepath.Clean(targetPath)
		if runtime.GOOS == "windows" {
			return strings.EqualFold(sourcePath, targetPath)
		}
		return sourcePath == targetPath
	}

	srcAbs, srcErr := filepath.Abs(sourcePath)
	dstAbs, dstErr := filepath.Abs(targetPath)
	if srcErr != nil || dstErr != nil {
		return false
	}
	srcNorm := filepath.Clean(srcAbs)
	dstNorm := filepath.Clean(dstAbs)
	if runtime.GOOS == "windows" {
		return strings.EqualFold(srcNorm, dstNorm)
	}
	return srcNorm == dstNorm
}
