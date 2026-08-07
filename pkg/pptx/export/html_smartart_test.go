package export

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

func smartArtSlide() elements.SlideContent {
	diagram := smartart.NewSmartArt(smartart.BasicProcess).
		Position(styling.Inches(1), styling.Inches(1)).
		Size(styling.Inches(8), styling.Inches(2)).
		AddNode(smartart.NewNode("Discover")).
		AddNode(smartart.NewNode("Design")).
		AddNode(smartart.NewNode("Deliver"))
	return elements.NewSlide("Process").AddSmartArt(diagram)
}

func TestHTMLDrawsSmartArtDiagrams(t *testing.T) {
	out := HTML("Deck", []elements.SlideContent{smartArtSlide()})

	if !strings.Contains(out, "Discover") || !strings.Contains(out, "Deliver") {
		t.Error("the diagram's captions are missing from the HTML")
	}
	if strings.Count(out, "<svg") == 0 {
		t.Error("no SVG was emitted for the diagram")
	}
	if !strings.Contains(out, "#"+smartArtNodeFill) {
		t.Errorf("the nodes were not filled with the accent %q", smartArtNodeFill)
	}
	if !strings.Contains(out, `fill="#`+smartArtNodeTextColor+`"`) {
		t.Error("the captions were not drawn in the node text colour")
	}
}

func TestHTMLDrawsSmartArtFromItsCachedLayout(t *testing.T) {
	diagram := smartart.NewSmartArt(smartart.AccentProcess).
		Position(styling.Inches(1), styling.Inches(1)).
		Size(styling.Inches(8), styling.Inches(2)).
		AddNode(smartart.NewNode("Cached")).
		WithDrawing([]smartart.DrawingShape{{
			X: 0, Y: 0, CX: 1828800, CY: 914400,
			PresetGeom: "ellipse",
			Fill:       smartart.ColorRef{SRGB: "FF0000"},
			Paragraphs: []smartart.DrawingParagraph{{
				Runs: []smartart.DrawingRun{{Text: "Cached", SizePt: 18}},
			}},
		}})

	out := HTML("Deck", []elements.SlideContent{elements.NewSlide("Cached").AddSmartArt(diagram)})
	if !strings.Contains(out, "<ellipse") {
		t.Error("the cached ellipse was not drawn, so the cache was ignored")
	}
	if !strings.Contains(out, "#FF0000") {
		t.Error("the cached fill colour was not used")
	}
}

func TestHTMLSkipsDecorativeSmartArt(t *testing.T) {
	slide := smartArtSlide()
	slide.SmartArtDiagrams[0] = slide.SmartArtDiagrams[0].WithDecorative(true)

	if got := smartArtSVGShapes(slide); got != nil {
		t.Errorf("got %d shapes for a decorative diagram, want none", len(got))
	}
}

func TestHTMLDrawsEveryTableOnTheSlide(t *testing.T) {
	first := tables.NewTable([]styling.Length{styling.Inches(1)}).AddRow([]string{"first table"})
	second := tables.NewTable([]styling.Length{styling.Inches(1)}).AddRow([]string{"second table"})

	slide := elements.NewSlide("Tables").WithTable(first)
	slide.Tables = append(slide.Tables, second)

	out := HTML("Deck", []elements.SlideContent{slide})
	if !strings.Contains(out, "first table") {
		t.Error("the primary table is missing")
	}
	if !strings.Contains(out, "second table") {
		t.Error("the second table is missing — only the first used to be drawn")
	}
}
