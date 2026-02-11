package testutil

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// ReadZipFile extracts one file from a zip.Reader and returns its content as a string.
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

// ZipHasFile returns true if the zip archive contains a file with the given name.
func ZipHasFile(zr *zip.Reader, name string) bool {
	for _, f := range zr.File {
		if f.Name == name {
			return true
		}
	}
	return false
}

// RootTestdataPath returns the path to a file inside the testdata directory,
// searching upward from the current working directory.
func RootTestdataPath(parts ...string) string {
	base := "testdata"
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(base); err == nil {
			return filepath.Join(append([]string{base}, parts...)...)
		}
		base = filepath.Join("..", base)
	}
	return filepath.Join(append([]string{base}, parts...)...)
}
