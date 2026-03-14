package editor

import (
	"archive/zip"
	"maps"
	"os"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

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
