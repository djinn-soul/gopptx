package editor

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestVerifyMediaInventoryChecksumsParallel_Pass(t *testing.T) {
	e := newMediaEditorFixture()
	if _, err := e.RegisterMedia([]byte("media-bytes"), "mp3"); err != nil {
		t.Fatalf("RegisterMedia failed: %v", err)
	}
	if err := e.verifyMediaInventoryChecksumsParallel(); err != nil {
		t.Fatalf("verifyMediaInventoryChecksumsParallel failed: %v", err)
	}
}

func TestSaveFailsWhenMediaChecksumMismatches(t *testing.T) {
	base := writeDeckFixture(t, "media-checksum-base.pptx", []elements.SlideContent{
		elements.NewSlide("Media checksum test"),
	})
	editor, err := OpenPresentationEditor(base)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	if _, err := editor.AddAudio(0, []byte("audio-bytes"), "audio/mpeg", 10, 10, 300, 120); err != nil {
		t.Fatalf("AddAudio failed: %v", err)
	}

	mediaParts := editor.parts.KeysWithPrefix("ppt/media/")
	if len(mediaParts) == 0 {
		t.Fatal("expected media part after AddAudio")
	}
	editor.parts.Set(mediaParts[0], []byte("tampered-audio-payload"))

	out := filepath.Join(t.TempDir(), "tampered-media-save.pptx")
	err = editor.Save(out)
	if err == nil || !strings.Contains(err.Error(), "checksum mismatch") {
		t.Fatalf("expected checksum mismatch error, got: %v", err)
	}
}
