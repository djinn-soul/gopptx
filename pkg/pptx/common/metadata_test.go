package common

import (
	"testing"
	"time"
)

func TestSlideSize(t *testing.T) {
	s43 := GetSlideSize4x3()
	if s43.Width != width4x3 || s43.Height != height4x3 {
		t.Errorf("expected 4:3 size (%d, %d), got (%d, %d)", width4x3, height4x3, s43.Width, s43.Height)
	}

	s169 := GetSlideSize16x9()
	if s169.Width != width16x9 || s169.Height != height4x3 {
		t.Errorf("expected 16:9 size (%d, %d), got (%d, %d)", width16x9, height4x3, s169.Width, s169.Height)
	}
}

func TestMetadata(t *testing.T) {
	now := time.Now()
	meta := Metadata{
		Title:         "Test Title",
		Creator:       "Test Creator",
		GeneratedDate: now,
	}

	if meta.Title != "Test Title" {
		t.Errorf("expected title 'Test Title', got '%s'", meta.Title)
	}
	if meta.Creator != "Test Creator" {
		t.Errorf("expected creator 'Test Creator', got '%s'", meta.Creator)
	}
	if !meta.GeneratedDate.Equal(now) {
		t.Error("expected generated date to match")
	}
}

func TestNormalizeHexColor(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"#FFFFFF", "FFFFFF"},
		{"FFFFFF", "FFFFFF"},
		{"#abc", "AABBCC"},
		{"abc", "AABBCC"},
		{"  #123456  ", "123456"},
	}

	for _, tt := range tests {
		got := NormalizeHexColor(tt.input)
		if got != tt.expected {
			t.Errorf("NormalizeHexColor(%q) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}

func TestIsHexColor(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"FFFFFF", true},
		{"#FFFFFF", true},
		{"#abc", true},
		{"GHIJKL", false},
		{"12345", false},
		{"1234567", false},
	}

	for _, tt := range tests {
		got := IsHexColor(tt.input)
		if got != tt.expected {
			t.Errorf("IsHexColor(%q) = %v; want %v", tt.input, got, tt.expected)
		}
	}
}

func TestNewGUID(t *testing.T) {
	guid1, err := NewGUID()
	if err != nil {
		t.Fatalf("NewGUID failed: %v", err)
	}
	guid2, err := NewGUID()
	if err != nil {
		t.Fatalf("NewGUID failed: %v", err)
	}

	if guid1 == guid2 {
		t.Error("expected unique GUIDs")
	}

	if len(guid1) != 38 { // {XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX}
		t.Errorf("expected GUID length 38, got %d", len(guid1))
	}

	if guid1[0] != '{' || guid1[37] != '}' {
		t.Errorf("expected GUID to be wrapped in {}, got %s", guid1)
	}
}
