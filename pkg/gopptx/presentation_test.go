package gopptx

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPresentation_Save(t *testing.T) {
	pres := &Presentation{}
	pres.AddSlide()

	filename := filepath.Join(t.TempDir(), "test_save.pptx")
	defer os.Remove(filename)

	err := pres.Save(filename)
	if err != nil {
		t.Fatalf("Failed to save presentation: %v", err)
	}

	if _, statErr := os.Stat(filename); os.IsNotExist(statErr) {
		t.Errorf("Expected file %s to be created", filename)
	}
}

func TestPresentation_SaveError(t *testing.T) {
	pres := &Presentation{}
	// Invalid path (directory doesn't exist)
	err := pres.Save("nonexistent_dir/test.pptx")
	if err == nil {
		t.Error("Expected error when saving to nonexistent directory, got nil")
	}
}
