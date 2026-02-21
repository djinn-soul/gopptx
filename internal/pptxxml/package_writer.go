package pptxxml

import (
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"
	"strings"
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

// WriteTo writes all collected parts to the provided [zip.Writer].
func (pw *PackageWriter) WriteTo(zw *zip.Writer) error {
	// Note: In an OPC package, order doesn't strictly matter for most tools,
	// but [Content_Types].xml and _rels/.rels are usually first.
	// For now, we just iterate the map.
	for path, content := range pw.textParts {
		method := packageZipMethod(path)
		var (
			w         io.Writer
			createErr error
		)
		if method == zip.Deflate {
			w, createErr = zw.Create(path)
		} else {
			header := &zip.FileHeader{Name: path, Method: method}
			w, createErr = zw.CreateHeader(header)
		}
		if createErr != nil {
			return createErr
		}
		if _, err := io.WriteString(w, content); err != nil {
			return fmt.Errorf("write package part %q: %w", path, err)
		}
	}
	for path, content := range pw.binaryParts {
		method := packageZipMethod(path)
		var (
			w         io.Writer
			createErr error
		)
		if method == zip.Deflate {
			w, createErr = zw.Create(path)
		} else {
			header := &zip.FileHeader{Name: path, Method: method}
			w, createErr = zw.CreateHeader(header)
		}
		if createErr != nil {
			return createErr
		}
		if _, err := w.Write(content); err != nil {
			return fmt.Errorf("write package part %q: %w", path, err)
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

func binaryZipMethod(path string) uint16 {
	return packageZipMethod(path)
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
