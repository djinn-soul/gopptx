package pptxxml

import (
	"bytes"
	"testing"
)

func TestPackageWriter_StandaloneWriteFile(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteFile(&buf, "hello"); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	if buf.String() != "hello" {
		t.Error("content mismatch")
	}
}

func TestPackageZipMethod(t *testing.T) {
	if packageZipMethod("test.xml") != 8 { // Deflate
		t.Error("expected 8 for xml")
	}
	if packageZipMethod("test.png") != 0 { // Store
		t.Error("expected 0 for png")
	}
	if packageZipMethod("ppt/notes/note1.xml") != 0 {
		t.Error("expected 0 for notes")
	}
}
