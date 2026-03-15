package editor

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

func TestMediaSaveIncludesWmvWmaOggContentTypes(t *testing.T) {
	base := writeDeckFixture(
		t,
		"media-format-base.pptx",
		[]elements.SlideContent{elements.NewSlide("Slide 1")},
	)
	editor, err := OpenPresentationEditor(base)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	if _, err := editor.AddVideo(0, []byte("video-wmv"), testutil.TinyPNG(), "video/x-ms-wmv", 10, 10, 100, 80); err != nil {
		t.Fatalf("AddVideo wmv failed: %v", err)
	}
	if _, err := editor.AddAudio(0, []byte("audio-wma"), "audio/x-ms-wma", 120, 10, 80, 40); err != nil {
		t.Fatalf("AddAudio wma failed: %v", err)
	}
	if _, err := editor.AddAudio(0, []byte("audio-ogg"), "audio/ogg", 210, 10, 80, 40); err != nil {
		t.Fatalf("AddAudio ogg failed: %v", err)
	}

	out := filepath.Join(t.TempDir(), "media-format-out.pptx")
	if err := editor.Save(out); err != nil {
		t.Fatalf("save output: %v", err)
	}

	raw, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	zr, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		t.Fatalf("open output zip: %v", err)
	}

	contentTypes := testutil.ReadZipFile(t, zr, "[Content_Types].xml")
	mustContain(t, contentTypes, `Extension="wmv" ContentType="video/x-ms-wmv"`)
	mustContain(t, contentTypes, `Extension="wma" ContentType="audio/x-ms-wma"`)
	mustContain(t, contentTypes, `Extension="ogg" ContentType="audio/ogg"`)
}

func mustContain(t *testing.T, body string, needle string) {
	t.Helper()
	if !strings.Contains(body, needle) {
		t.Fatalf("expected content to include %q, got: %s", needle, body)
	}
}
