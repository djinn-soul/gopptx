package fonts

import (
	"bytes"
	"testing"
)

func TestFontStyle_XMLElement(t *testing.T) {
	tests := []struct {
		style    FontStyle
		expected string
	}{
		{StyleRegular, "regular"},
		{StyleBold, "bold"},
		{StyleItalic, "italic"},
		{StyleBoldItalic, "boldItalic"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.style.XMLElement(); got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestEmbeddedFont_Builder(t *testing.T) {
	data := []byte{1, 2, 3}
	f := New("Arial", StyleBold, data).
		WithCharset(CharsetRussian).
		WithPanose("020F0502020204030204").
		WithPitchFamily(0x34)

	if f.Typeface != "Arial" {
		t.Errorf("expected Arial, got %v", f.Typeface)
	}
	if f.Style != StyleBold {
		t.Errorf("expected StyleBold, got %v", f.Style)
	}
	if f.Charset != CharsetRussian {
		t.Errorf("expected CharsetRussian, got %v", f.Charset)
	}
	if f.Panose != "020F0502020204030204" {
		t.Errorf("expected panose, got %v", f.Panose)
	}
	if f.PitchFamily != 0x34 {
		t.Errorf("expected 0x34, got %v", f.PitchFamily)
	}
	if !bytes.Equal(f.Data, data) {
		t.Errorf("data mismatch")
	}
}

func TestObfuscateFont(t *testing.T) {
	// A mock valid GUID "12345678-9ABC-DEF0-1234-56789ABCDEF0"
	guid := "12345678-9ABC-DEF0-1234-56789ABCDEF0"

	// GUID bytes would be: 12 34 56 78 9A BC DE F0 12 34 56 78 9A BC DE F0
	// Reversed key: F0 DE BC 9A 78 56 34 12 F0 DE BC 9A 78 56 34 12

	fontData := make([]byte, 64)
	for i := range 64 {
		fontData[i] = byte(i) // Predictable dummy data
	}

	obfuscated := ObfuscateFont(fontData, guid)

	// Since obfuscation only touches the first 32 bytes, the remaining should be identical.
	if !bytes.Equal(fontData[32:], obfuscated[32:]) {
		t.Errorf("obfuscation modified data beyond the first 32 bytes")
	}

	// Verify the obfuscation is symmetric (XORing again with same GUID yields original data)
	restored := ObfuscateFont(obfuscated, guid)

	if !bytes.Equal(fontData, restored) {
		t.Errorf("obfuscation is not symmetric")
	}

	// Test short data (less than 32 bytes)
	shortData := []byte{1, 2, 3, 4, 5}
	shortObfuscated := ObfuscateFont(shortData, guid)
	if !bytes.Equal(shortData, shortObfuscated) {
		t.Errorf("short data should be returned as-is")
	}

	// Test invalid GUID
	invalidGUID := "not-a-guid"
	invalidObfuscated := ObfuscateFont(fontData, invalidGUID)
	if !bytes.Equal(fontData, invalidObfuscated) {
		t.Errorf("invalid GUID should return data as-is")
	}
}
