package opc

import (
	"archive/zip"
	"bytes"
	"testing"
)

func TestWriter(t *testing.T) {
	buf := new(bytes.Buffer)
	w := NewWriter(buf)

	err := w.AddFile("test.xml", []byte("<test></test>"))
	if err != nil {
		t.Fatalf("Failed to add file: %v", err)
	}

	err = w.Close()
	if err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}

	// Verify ZIP content
	r, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("Failed to create zip reader: %v", err)
	}

	if len(r.File) != 1 {
		t.Errorf("Expected 1 file, got %d", len(r.File))
	}

	if r.File[0].Name != "test.xml" {
		t.Errorf("Expected test.xml, got %s", r.File[0].Name)
	}
}
