package presentation

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

// The reader fills Shape.TextParagraphs and the PDF exporter renders it, but the
// PPTX generator only ever looked at the flat Text field — so a shape read from
// a deck and written back came out with no text at all.
func TestGeneratorEmitsShapeTextParagraphs(t *testing.T) {
	shape := shapes.NewShape("rect", styling.Inches(1), styling.Inches(1), styling.Inches(4), styling.Inches(2))
	shape.TextParagraphs = []text.Paragraph{
		{Runs: []text.Run{{Text: "PARA_ONE", Bold: true}}},
		{Runs: []text.Run{{Text: "PARA_TWO"}}},
	}

	slide := elements.NewSlide("Shapes")
	slide.Shapes = []shapes.Shape{shape}

	parts := buildPackageParts(t, Metadata{}, []elements.SlideContent{slide})
	slideXML := parts["ppt/slides/slide1.xml"]

	for _, want := range []string{"PARA_ONE", "PARA_TWO"} {
		if !strings.Contains(slideXML, want) {
			t.Fatalf("slide1.xml does not contain %q: %s", want, slideXML)
		}
	}
	if !strings.Contains(slideXML, `b="1"`) {
		t.Fatalf("bold run formatting lost: %s", slideXML)
	}
}

// A shape with only the flat Text keeps rendering as before.
func TestGeneratorStillEmitsFlatShapeText(t *testing.T) {
	shape := shapes.NewShape("rect", styling.Inches(1), styling.Inches(1), styling.Inches(4), styling.Inches(2)).
		WithText("FLAT_TEXT")

	slide := elements.NewSlide("Shapes")
	slide.Shapes = []shapes.Shape{shape}

	parts := buildPackageParts(t, Metadata{}, []elements.SlideContent{slide})
	if !strings.Contains(parts["ppt/slides/slide1.xml"], "FLAT_TEXT") {
		t.Fatalf("flat shape text lost: %s", parts["ppt/slides/slide1.xml"])
	}
}
