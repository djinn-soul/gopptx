package editor

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestSetSlideHeaderFooterInjectsVisibleOverlayShapes(t *testing.T) {
	basePath := writeDeckFixture(t, "header-footer-base.pptx", []elements.SlideContent{
		elements.NewSlide("Headers and Footers"),
	})

	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	if err := editor.SetSlideHeaderFooter(0, SlideHeaderFooter{
		Footer:       "Confidential",
		ShowFooter:   true,
		ShowSlideNum: true,
		ShowDateTime: true,
		DateTimeText: "2026-03-20",
	}); err != nil {
		t.Fatalf("set slide header/footer: %v", err)
	}

	outPath := filepath.Join(t.TempDir(), "header-footer-visible.pptx")
	if err := editor.Save(outPath); err != nil {
		t.Fatalf("save edited deck: %v", err)
	}

	slideXML := string(readZipFileBytes(t, outPath, "ppt/slides/slide1.xml"))
	expectedSnippets := []string{
		"Confidential",
		"2026-03-20",
		"Slide Number Visible",
		"Footer Visible",
		"Date Visible",
	}
	for _, snippet := range expectedSnippets {
		if !strings.Contains(slideXML, snippet) {
			t.Fatalf("expected slide XML to contain %q", snippet)
		}
	}
}
