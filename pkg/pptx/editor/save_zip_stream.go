package editor

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"sort"
	"sync"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/internal/zipfast"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
)

//nolint:gochecknoglobals // Shared reusable buffers reduce allocations on save hot paths.
var rawZipCopyBufferPool = sync.Pool{
	New: func() any {
		b := make([]byte, rawZipCopyBufferSize)
		return &b
	},
}

// packageRelsPartName is the package-level relationships part.
const packageRelsPartName = "_rels/.rels"

// mergedPartNames returns the sorted union of existing part keys and any extra
// keys present in updatedParts that are not already in keys.
func mergedPartNames(keys []string, updatedParts map[string][]byte) []string {
	// First pass: count extras to pre-size slice (avoids append re-growth).
	extraCount := 0
	for k := range updatedParts {
		if i := sort.SearchStrings(keys, k); i >= len(keys) || keys[i] != k {
			extraCount++
		}
	}
	if extraCount == 0 {
		return keys
	}
	extraKeys := make([]string, 0, extraCount)
	for k := range updatedParts {
		if i := sort.SearchStrings(keys, k); i >= len(keys) || keys[i] != k {
			extraKeys = append(extraKeys, k)
		}
	}
	sort.Strings(extraKeys)
	merged := make([]string, 0, len(keys)+len(extraKeys))
	ki, ei := 0, 0
	for ki < len(keys) && ei < len(extraKeys) {
		if keys[ki] <= extraKeys[ei] {
			merged = append(merged, keys[ki])
			ki++
		} else {
			merged = append(merged, extraKeys[ei])
			ei++
		}
	}
	merged = append(merged, keys[ki:]...)
	merged = append(merged, extraKeys[ei:]...)
	return merged
}

// orderPackageNames moves the content types stream to the front, then the
// package relationships. OPC requires the content types stream to be the first
// part in the archive, so a consumer reading the package as a stream knows a
// part's content type before it reaches the part.
//
// Sorted names put "[Content_Types].xml" first already for every part path seen
// in practice, because "[" sorts below "_" and the lowercase directory names.
// A source deck carrying a part whose name starts with an uppercase letter
// would break that, so the order is made explicit rather than inherited.
func orderPackageNames(names []string) []string {
	firstIdx, relsIdx := -1, -1
	for i, name := range names {
		switch name {
		case pptxxml.ContentTypesPartName:
			firstIdx = i
		case packageRelsPartName:
			relsIdx = i
		}
	}
	if firstIdx <= 0 && (relsIdx < 0 || relsIdx <= 1) {
		return names // already in order
	}

	ordered := make([]string, 0, len(names))
	if firstIdx >= 0 {
		ordered = append(ordered, names[firstIdx])
	}
	if relsIdx >= 0 {
		ordered = append(ordered, names[relsIdx])
	}
	for i, name := range names {
		if i == firstIdx || i == relsIdx {
			continue
		}
		ordered = append(ordered, name)
	}
	return ordered
}

// writeZipData writes data into a new zip entry, choosing Store vs Deflate
// based on the file name.
func writeZipData(zw *zip.Writer, name string, data []byte) error {
	w, createErr := zipfast.CreateEntry(zw, name, editorslide.SaveZipMethod(name))
	if createErr != nil {
		return fmt.Errorf("create zip entry %q: %w", name, createErr)
	}
	if _, writeErr := w.Write(data); writeErr != nil {
		return fmt.Errorf("write zip entry %q: %w", name, writeErr)
	}
	return nil
}

// buildZipStream writes all parts into a new zip archive and returns the raw bytes.
func (e *PresentationEditor) buildZipStream(
	allNames []string,
	updatedParts map[string][]byte,
) ([]byte, error) {
	var zipBuf bytes.Buffer
	if err := e.buildZipToWriter(&zipBuf, allNames, updatedParts); err != nil {
		return nil, err
	}
	return zipBuf.Bytes(), nil
}

func (e *PresentationEditor) buildZipToWriter(
	w io.Writer,
	allNames []string,
	updatedParts map[string][]byte,
) error {
	zw := zipfast.NewWriter(w)
	poolBuf, ok := rawZipCopyBufferPool.Get().(*[]byte)
	if !ok || poolBuf == nil || cap(*poolBuf) < rawZipCopyBufferSize {
		fresh := make([]byte, rawZipCopyBufferSize)
		poolBuf = &fresh
	}
	rawZipCopyBuffer := (*poolBuf)[:rawZipCopyBufferSize]
	defer rawZipCopyBufferPool.Put(poolBuf)

	for _, name := range orderPackageNames(allNames) {
		if updated, updatedOK := updatedParts[name]; updatedOK {
			if err := writeZipData(zw, name, updated); err != nil {
				return err
			}
			continue
		}
		if sourceEntry, sourceOK := e.parts.SourceZipEntry(name); sourceOK {
			if err := copyZipEntryRaw(zw, sourceEntry, rawZipCopyBuffer); err != nil {
				return fmt.Errorf("copy source zip entry %q: %w", name, err)
			}
			continue
		}
		content, partOK := e.parts.Get(name)
		if !partOK {
			return fmt.Errorf("failed to retrieve part %q during save", name)
		}
		if err := writeZipData(zw, name, content); err != nil {
			return err
		}
	}

	if err := zw.Close(); err != nil {
		return fmt.Errorf("finalize zip stream: %w", err)
	}
	return nil
}

func copyZipEntryRaw(zw *zip.Writer, sourceEntry *zip.File, copyBuffer []byte) error {
	header := sourceEntry.FileHeader
	writer, err := zw.CreateRaw(&header)
	if err != nil {
		return err
	}
	reader, err := sourceEntry.OpenRaw()
	if err != nil {
		return err
	}
	if _, err := io.CopyBuffer(writer, reader, copyBuffer); err != nil {
		return err
	}
	return nil
}
