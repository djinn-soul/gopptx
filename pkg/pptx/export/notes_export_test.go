package export

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func slidesWithNotes() []elements.SlideContent {
	return []elements.SlideContent{
		elements.NewSlide("First").
			AddBullet("Visible bullet").
			WithNotes("REMEMBER_THE_DEMO"),
		elements.NewSlide("Second").AddBullet("No notes here"),
	}
}

// No export path could emit speaker notes, so a deck's notes were lost on the
// way to HTML even though the model carried them.
func TestHTMLExportCanIncludeNotes(t *testing.T) {
	opts := DefaultHTMLOptions()
	opts.IncludeNotes = true

	withNotes := HTMLWithOptions("Deck", slidesWithNotes(), opts)
	if !strings.Contains(withNotes, "REMEMBER_THE_DEMO") {
		t.Fatal("notes missing from the HTML export")
	}
	if !strings.Contains(withNotes, "slide-notes") {
		t.Fatal("notes are not wrapped in their own block")
	}

	// Off by default, so existing output is unchanged.
	if strings.Contains(HTML("Deck", slidesWithNotes()), "REMEMBER_THE_DEMO") {
		t.Fatal("notes should not appear unless asked for")
	}
}

func TestNativePDFCanIncludeNotesAndFrontmatter(t *testing.T) {
	out := filepath.Join(t.TempDir(), "deck.pdf")
	opts := PDFOptions{Driver: PDFDriverNative, IncludeNotes: true, IncludeFrontmatter: true}

	if err := PDFWithOptions("Deck Title", slidesWithNotes(), out, opts); err != nil {
		t.Fatalf("PDFWithOptions: %v", err)
	}
	info, err := os.Stat(out)
	if err != nil {
		t.Fatalf("stat output: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("empty PDF")
	}

	// The frontmatter and the one notes page add two pages to the two slides.
	plain := filepath.Join(t.TempDir(), "plain.pdf")
	if err := PDFWithOptions("Deck Title", slidesWithNotes(), plain, PDFOptions{Driver: PDFDriverNative}); err != nil {
		t.Fatalf("PDFWithOptions plain: %v", err)
	}
	plainInfo, err := os.Stat(plain)
	if err != nil {
		t.Fatalf("stat plain output: %v", err)
	}
	if info.Size() <= plainInfo.Size() {
		t.Fatalf("notes/frontmatter export (%d bytes) is not larger than plain (%d bytes)",
			info.Size(), plainInfo.Size())
	}
}
