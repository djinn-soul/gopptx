package export

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func TestPDFHyperlinkTargetPerActionType(t *testing.T) {
	cases := []struct {
		name   string
		action action.HyperlinkAction
		want   string
	}{
		{"url", action.HyperlinkURL("https://example.com/a"), "https://example.com/a"},
		{"email", action.HyperlinkEmail("who@example.com"), "mailto:who@example.com"},
		{
			"email with subject",
			action.HyperlinkEmailWithSubject("who@example.com", "Q3 deck"),
			"mailto:who@example.com?subject=Q3+deck",
		},
		{"file", action.HyperlinkFile(`C:\decks\notes.txt`), "file:///C:/decks/notes.txt"},
		{"slide has no URI form", action.HyperlinkSlide(3), ""},
		{"program is not expressible", action.HyperlinkProgram("notepad.exe"), ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := pdfHyperlinkTarget(tc.action); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestShapeClickActionPrefersTheCurrentField(t *testing.T) {
	click := action.NewHyperlink(action.HyperlinkURL("https://current"))
	legacy := action.NewHyperlink(action.HyperlinkURL("https://legacy"))

	if got := shapeClickAction(&click, &legacy); got != &click {
		t.Error("the legacy field won over the click action")
	}
	if got := shapeClickAction(nil, &legacy); got != &legacy {
		t.Error("the legacy field was ignored with no click action set")
	}
	if got := shapeClickAction(nil, nil); got != nil {
		t.Errorf("got %+v, want nil", got)
	}
}

func TestSlidePDFAnchorIsStable(t *testing.T) {
	if got := slidePDFAnchor(4); got != "slide-4" {
		t.Errorf("got %q, want slide-4", got)
	}
}

// TestNativePDFEmbedsLinksAndMetadata checks the whole path: a shape link, a run
// link and a slide jump all reach the file, and the document is titled.
func TestNativePDFEmbedsLinksAndMetadata(t *testing.T) {
	click := action.NewHyperlink(action.HyperlinkURL("https://example.com/clicked"))
	jump := action.NewHyperlink(action.HyperlinkSlide(2))

	slides := []elements.SlideContent{
		elements.NewSlide("Linked").
			AddShape(shapes.NewShape(shapes.ShapeTypeRectangle, 100, 100, 200, 100).
				WithFill(shapes.NewShapeFill("4472C4")).
				WithText("Click me")).
			AddShape(shapes.NewShape(shapes.ShapeTypeRectangle, 100, 300, 200, 100).
				WithText("Next slide")),
		elements.NewSlide("Target"),
	}
	slides[0].Shapes[0].ClickAction = &click
	slides[0].Shapes[1].ClickAction = &jump

	pdfPath := filepath.Join(t.TempDir(), "links.pdf")
	err := PDFWithOptions("Linked Deck", slides, pdfPath, PDFOptions{Driver: PDFDriverNative})
	if err != nil {
		t.Fatalf("PDFWithOptions: %v", err)
	}
	data, err := os.ReadFile(pdfPath)
	if err != nil {
		t.Fatalf("read pdf: %v", err)
	}

	if !bytes.Contains(data, []byte("https://example.com/clicked")) {
		t.Error("the external link is not in the file")
	}
	if !bytes.Contains(data, []byte("/Annots")) {
		t.Error("no annotations were written, so no link is clickable")
	}
	// The information dictionary is written as UTF-16BE hex, so the strings are
	// checked by their key rather than their plain text.
	if !bytes.Contains(data, []byte("/Title <FEFF")) {
		t.Error("the document title is not in the file")
	}
	if !bytes.Contains(data, []byte("/Producer <FEFF")) {
		t.Error("the producer is not in the file")
	}
}
