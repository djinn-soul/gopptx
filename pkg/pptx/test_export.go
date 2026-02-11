package pptx

import (
	"archive/zip"
	"bytes"
	"path/filepath"
	"testing"
)

// RootTestdataPath returns the absolute path to testdata from the package root.
func RootTestdataPath(parts ...string) string {
	base := []string{"..", "..", "testdata"}
	base = append(base, parts...)
	return filepath.Join(base...)
}

// ReadZipFile is a test helper that extracts one file from a zip.Reader.
func ReadZipFile(t *testing.T, zr *zip.Reader, name string) string {
	t.Helper()
	for _, f := range zr.File {
		if f.Name != name {
			continue
		}
		r, err := f.Open()
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}
		defer func() { _ = r.Close() }()
		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(r); err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		return buf.String()
	}
	t.Fatalf("file %s not found in zip", name)
	return ""
}
