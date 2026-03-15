package interop

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConvertFromPpt_NotFound(t *testing.T) {
	_, err := ConvertFromPpt("non_existent_file.ppt", "")
	if err == nil {
		t.Errorf("Expected error when file does not exist")
	}
	if !strings.Contains(err.Error(), "input file not found") {
		t.Errorf("Expected 'input file not found' error, got %v", err)
	}
}

func TestConvertFromPpt_EmptyInput(t *testing.T) {
	_, err := ConvertFromPpt("   ", "")
	if err == nil {
		t.Errorf("Expected error on empty input")
	}
	if !strings.Contains(err.Error(), "inputPath is empty") {
		t.Errorf("Expected 'inputPath is empty' error, got %v", err)
	}
}

func TestConvertFromPpt_WithValidFakeFileSkipped(t *testing.T) {
	// We can't guarantee LibreOffice is installed in CI or the user's dev env for this test.
	// But we can test the path up to `findSoffice` or execution by creating a dummy input file.
	tmpDir := t.TempDir()
	dummyPPT := filepath.Join(tmpDir, "dummy.ppt")
	err := os.WriteFile(dummyPPT, []byte("fake ole2 data"), 0o644)
	if err != nil {
		t.Fatalf("failed to create dummy file: %v", err)
	}

	outDir := filepath.Join(tmpDir, "out")

	// Call it
	relPath, err := ConvertFromPpt(dummyPPT, outDir)

	if err == nil {
		// If it somehow worked (unlikely with fake data, but maybe LibreOffice just generates an empty one)
		if filepath.Base(relPath) != "dummy.pptx" {
			t.Errorf("Expected output file dummy.pptx, got %s", relPath)
		}
		if _, statErr := os.Stat(relPath); statErr != nil {
			t.Errorf("Expected output file to exist on success: %v", statErr)
		}
		return
	}

	switch {
	case strings.Contains(err.Error(), "libreoffice required"),
		strings.Contains(err.Error(), "soffice binary not found"):
		t.Skipf("Skipping integration test: LibreOffice not installed on host. Error: %v", err)
	case strings.Contains(err.Error(), "conversion failed"):
		t.Skipf("Skipping integration test: LibreOffice failed to digest the fake PPT. Error: %v", err)
	default:
		t.Errorf("Unexpected error during conversion attempt: %v", err)
	}
}

func TestConvertFromPpt_OutDirError(t *testing.T) {
	tmpDir := t.TempDir()
	dummyPPT := filepath.Join(tmpDir, "dummy.ppt")
	if err := os.WriteFile(dummyPPT, []byte("fake ole2 data"), 0o644); err != nil {
		t.Fatalf("failed: %v", err)
	}

	// Make outDir an existing file instead of a directory to force MkdirAll error
	badOutDir := filepath.Join(tmpDir, "bad_out")
	if err := os.WriteFile(badOutDir, []byte("file"), 0o644); err != nil {
		t.Fatalf("failed: %v", err)
	}

	_, err := ConvertFromPpt(dummyPPT, badOutDir)
	if err == nil || !strings.Contains(err.Error(), "failed to create outDir") {
		t.Errorf("expected outDir creation error, got %v", err)
	}
}

func TestConvertFromPpt_EmptyOutDir(t *testing.T) {
	tmpDir := t.TempDir()
	dummyPPT := filepath.Join(tmpDir, "dummy.ppt")
	if err := os.WriteFile(dummyPPT, []byte("fake ole2 data"), 0o644); err != nil {
		t.Fatalf("failed: %v", err)
	}

	// Test with empty outDir to hit the `outDir == ""` branch
	_, err := ConvertFromPpt(dummyPPT, "")
	if err != nil &&
		!strings.Contains(err.Error(), "libreoffice required") &&
		!strings.Contains(err.Error(), "conversion failed") &&
		!strings.Contains(err.Error(), "soffice binary not found") {
		t.Errorf("unexpected error: %v", err)
	}
}
