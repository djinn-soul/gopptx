package editor

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/presentation/protection"
)

func TestEditorSave_WithEncryptionPassword_WritesCFB(t *testing.T) {
	base := writeDeckFixture(t, "encryption-base.pptx", []elements.SlideContent{
		elements.NewSlide("Slide 1"),
	})

	ed, err := OpenPresentationEditor(base)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	ed.Metadata().Protection.EncryptPassword = "Secret123!"
	out := filepath.Join(t.TempDir(), "encrypted-output.pptx")

	if !protection.CanEncryptAgile() {
		// Encryption runtime unavailable: verify the error is returned correctly.
		saveErr := ed.Save(out)
		if saveErr == nil {
			t.Fatal("expected error when agile encryption unavailable, got nil")
		}
		if !strings.Contains(saveErr.Error(), "encryption") &&
			!strings.Contains(saveErr.Error(), "PowerPoint") &&
			!strings.Contains(saveErr.Error(), "unavailable") {
			t.Fatalf("unexpected error message: %v", saveErr)
		}
		return
	}

	if err := ed.Save(out); err != nil {
		// In some Windows environments PowerPoint COM is callable but rejects
		// synthetic test fixtures with runtime-specific HRESULTs.
		if strings.Contains(err.Error(), "0x80070570") ||
			strings.Contains(err.Error(), "BadImageFormatException") {
			t.Logf("PowerPoint COM rejected test fixture in this runtime: %v", err)
			return
		}
		t.Fatalf("save encrypted deck: %v", err)
	}

	data := readFile(t, out)
	if len(data) < 8 {
		t.Fatalf("encrypted output too short: %d", len(data))
	}
	if bytes.Equal(data[:4], []byte("PK\x03\x04")) {
		t.Fatal("expected encrypted CFB output, got zip header")
	}
	if !bytes.Equal(data[:8], []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}) {
		t.Fatalf("expected CFB signature, got %x", data[:8])
	}
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
