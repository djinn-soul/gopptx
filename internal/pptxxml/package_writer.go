package pptxxml

import (
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"github.com/djinn-soul/gopptx/internal/zipfast"
)

// PackageWriter accumulates parts and handles writing them to a zip archive.
// This provides a more robust OPC-compliant infrastructure than direct zip writing.
type PackageWriter struct {
	textParts   map[string]string
	binaryParts map[string][]byte
}

// NewPackageWriter creates a new PackageWriter.
func NewPackageWriter() *PackageWriter {
	return &PackageWriter{
		textParts:   make(map[string]string),
		binaryParts: make(map[string][]byte),
	}
}

// AddPart adds a text part to the package.
func (pw *PackageWriter) AddPart(path string, content string) {
	pw.textParts[path] = content
	delete(pw.binaryParts, path)
}

// AddBinaryPart adds a binary part to the package.
func (pw *PackageWriter) AddBinaryPart(path string, content []byte) {
	pw.binaryParts[path] = content
	delete(pw.textParts, path)
}

// ContentTypesPartName is the OPC content types stream, which must be the first
// entry in the package.
const ContentTypesPartName = "[Content_Types].xml"

// WriteTo writes all collected parts to the provided [zip.Writer].
//
// Entries are written in a fixed order: [Content_Types].xml first, then the
// package relationships, then everything else sorted by path. OPC requires the
// content types stream to be the first part in the archive — a consumer reading
// the package as a stream has to know a part's content type before it reaches
// the part. PowerPoint tolerates finding it anywhere; stricter readers, such as
// the previewers used by webmail, do not. Sorting the rest also makes the
// package byte-reproducible, which iterating the maps directly was not.
func (pw *PackageWriter) WriteTo(zw *zip.Writer) error {
	for _, path := range pw.orderedPartNames() {
		if content, ok := pw.textParts[path]; ok {
			if err := writePart(zw, path, []byte(content)); err != nil {
				return err
			}
			continue
		}
		if err := writePart(zw, path, pw.binaryParts[path]); err != nil {
			return err
		}
	}
	return nil
}

// orderedPartNames returns every part name in package write order.
func (pw *PackageWriter) orderedPartNames() []string {
	paths := make([]string, 0, len(pw.textParts)+len(pw.binaryParts))
	for path := range pw.textParts {
		paths = append(paths, path)
	}
	for path := range pw.binaryParts {
		paths = append(paths, path)
	}
	sort.Slice(paths, func(i, j int) bool {
		ri, rj := partWriteRank(paths[i]), partWriteRank(paths[j])
		if ri != rj {
			return ri < rj
		}
		return paths[i] < paths[j]
	})
	return paths
}

func partWriteRank(path string) int {
	switch path {
	case ContentTypesPartName:
		return 0
	case "_rels/.rels":
		return 1
	default:
		return 2
	}
}

func writePart(zw *zip.Writer, path string, content []byte) error {
	w, err := zipfast.CreateEntry(zw, path, packageZipMethod(path))
	if err != nil {
		return fmt.Errorf("create package part %q: %w", path, err)
	}
	if _, err := w.Write(content); err != nil {
		return fmt.Errorf("write package part %q: %w", path, err)
	}
	return nil
}

// WriteFile is a convenience helper to write string content to a writer (non-buffered).
//
// Deprecated: used for incremental migration.
func WriteFile(w io.Writer, content string) error {
	_, err := io.WriteString(w, content)
	return err
}

func packageZipMethod(path string) uint16 {
	if strings.HasPrefix(strings.ToLower(path), "ppt/notes") {
		return zip.Store
	}
	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tif", ".tiff", ".mp3", ".m4a", ".wav", ".mp4", ".avi":
		return zip.Store
	default:
		return zip.Deflate
	}
}
