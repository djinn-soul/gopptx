package editor

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"sort"
	"sync"

	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
)

//nolint:gochecknoglobals // Shared reusable buffers reduce allocations on save hot paths.
var rawZipCopyBufferPool = sync.Pool{
	New: func() any {
		b := make([]byte, rawZipCopyBufferSize)
		return &b
	},
}

// mergedPartNames returns the sorted union of existing part keys and any extra
// keys present in updatedParts that are not already in keys.
func mergedPartNames(keys []string, updatedParts map[string][]byte) []string {
	var extraKeys []string
	for k := range updatedParts {
		if i := sort.SearchStrings(keys, k); i >= len(keys) || keys[i] != k {
			extraKeys = append(extraKeys, k)
		}
	}
	if len(extraKeys) == 0 {
		return keys
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

// writeZipData writes data into a new zip entry, choosing Store vs Deflate
// based on the file name.
func writeZipData(zw *zip.Writer, name string, data []byte) error {
	var (
		w         io.Writer
		createErr error
	)
	if editorslide.SaveZipMethod(name) == zip.Store {
		header := &zip.FileHeader{Name: name, Method: zip.Store}
		w, createErr = zw.CreateHeader(header)
	} else {
		w, createErr = zw.Create(name)
	}
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
	zw := zip.NewWriter(w)
	poolBuf, ok := rawZipCopyBufferPool.Get().(*[]byte)
	if !ok || poolBuf == nil || cap(*poolBuf) < rawZipCopyBufferSize {
		fresh := make([]byte, rawZipCopyBufferSize)
		poolBuf = &fresh
	}
	rawZipCopyBuffer := (*poolBuf)[:rawZipCopyBufferSize]
	defer rawZipCopyBufferPool.Put(poolBuf)

	for _, name := range allNames {
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
