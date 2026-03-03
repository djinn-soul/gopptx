package common

import (
	"archive/zip"
	"bytes"
	"testing"
)

func TestColor_Normalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"#FF0000", "FF0000"},
		{"abc", "AABBCC"},
		{"  123456  ", "123456"},
		{"red", "RREEDD"}, // expanded because it has 3 chars
	}
	for _, tt := range tests {
		if got := NormalizeHexColor(tt.input); got != tt.expected {
			t.Errorf("NormalizeHexColor(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestColor_IsValid(t *testing.T) {
	if !IsHexColor("#FF0000") { t.Error("Valid hex failed") }
	if !IsHexColor("abc") { t.Error("Valid 3-digit failed") }
	if IsHexColor("red") { t.Error("Invalid color passed") }
}

func TestGeometry_SlideSize(t *testing.T) {
	s43 := GetSlideSize4x3()
	if s43.Width != 9144000 { t.Error("4x3 width mismatch") }

	s169 := GetSlideSize16x9()
	if s169.Width != 12192000 { t.Error("16x9 width mismatch") }
}

func TestGUID(t *testing.T) {
	g, err := NewGUID()
	if err != nil {
		t.Fatalf("NewGUID failed: %v", err)
	}
	if len(g) != 38 {
		t.Errorf("expected GUID length 38 ({...}), got %d", len(g))
	}
}

func TestZipHelpers(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	if err := WriteFile(zw, "test.txt", "hello"); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	if err := WriteBinaryFile(zw, "test.bin", []byte{1, 2, 3}); err != nil {
		t.Fatalf("WriteBinaryFile failed: %v", err)
	}

	zw.Close()

	zr, _ := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if len(zr.File) != 2 {
		t.Errorf("expected 2 files in zip, got %d", len(zr.File))
	}
}
