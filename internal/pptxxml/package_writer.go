package pptxxml

import (
	"archive/zip"
	"io"
)

// PackageWriter accumulates parts and handles writing them to a zip archive.
// This provides a more robust OPC-compliant infrastructure than direct zip writing.
type PackageWriter struct {
	parts map[string][]byte
}

// NewPackageWriter creates a new PackageWriter.
func NewPackageWriter() *PackageWriter {
	return &PackageWriter{
		parts: make(map[string][]byte),
	}
}

// AddPart adds a text part to the package.
func (pw *PackageWriter) AddPart(path string, content string) {
	pw.parts[path] = []byte(content)
}

// AddBinaryPart adds a binary part to the package.
func (pw *PackageWriter) AddBinaryPart(path string, content []byte) {
	pw.parts[path] = content
}

// WriteTo writes all collected parts to the provided [zip.Writer].
func (pw *PackageWriter) WriteTo(zw *zip.Writer) error {
	// Note: In an OPC package, order doesn't strictly matter for most tools,
	// but [Content_Types].xml and _rels/.rels are usually first.
	// For now, we just iterate the map.
	for path, content := range pw.parts {
		w, createErr := zw.Create(path)
		if createErr != nil {
			return createErr
		}
		if _, err := w.Write(content); err != nil {
			// TODO: Verify resource cleanup procedures on write failure.
			return err
		}
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
