package export

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

// The HTML export used to drop SmartArt entirely: a deck full of diagrams came
// out as a deck of titles and bullets with nothing where the diagrams were. The
// diagrams are drawn the same way the PDF exporter draws them — from the layout
// PowerPoint cached, and from the built-in layouts when there is no cache — and
// then emitted through the same SVG writer the slide's own shapes use.

// smartArtSVGShapes flattens a slide's diagrams into shapes the SVG writer can
// draw.
func smartArtSVGShapes(slide elements.SlideContent) []shapes.Shape {
	if len(slide.SmartArtDiagrams) == 0 {
		return nil
	}
	out := make([]shapes.Shape, 0, len(slide.SmartArtDiagrams))
	for _, diagram := range slide.SmartArtDiagrams {
		if diagram.IsDecorative {
			// A decorative diagram carries no information, and the HTML export
			// is the accessible rendering of the deck.
			continue
		}
		out = append(out, smartArtDiagramShapes(diagram)...)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func smartArtDiagramShapes(diagram smartart.SmartArt) []shapes.Shape {
	if len(diagram.Drawing) > 0 && !smartArtCacheIsStale(diagram) {
		out := make([]shapes.Shape, 0, len(diagram.Drawing))
		for _, cached := range diagram.Drawing {
			shape, ok := smartArtCachedShape(diagram, cached)
			if !ok {
				continue
			}
			out = append(out, shape)
		}
		if len(out) > 0 {
			return out
		}
	}

	boxes, _ := layoutSmartArt(diagram)
	out := make([]shapes.Shape, 0, len(boxes))
	for _, box := range boxes {
		out = append(out, smartArtNodeShape(box))
	}
	return out
}
