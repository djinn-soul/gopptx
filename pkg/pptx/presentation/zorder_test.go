package presentation

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func zShape(name string, zOrder int) shapes.Shape {
	s := shapes.NewShape("rect", styling.Inches(1), styling.Inches(1), styling.Inches(2), styling.Inches(1)).
		WithText(name)
	s.ZOrder = zOrder
	return s
}

// The reader records each element's position in the shape tree and the PDF
// exporter paints from it; the PPTX generator emitted the slice order instead,
// so a deck read and written back could paint back-to-front.
func TestGeneratorEmitsShapesInZOrder(t *testing.T) {
	slide := elements.NewSlide("Paint order")
	slide.Shapes = []shapes.Shape{zShape("ZA", 2), zShape("ZB", 1)}

	parts := buildPackageParts(t, Metadata{}, []elements.SlideContent{slide})
	slideXML := parts["ppt/slides/slide1.xml"]

	first, second := strings.Index(slideXML, "ZB"), strings.Index(slideXML, "ZA")
	if first < 0 || second < 0 {
		t.Fatalf("both shapes should be present: %s", slideXML)
	}
	if first > second {
		t.Fatalf("ZB (ZOrder 1) must be written before ZA (ZOrder 2)")
	}
}

// Shapes that state no order keep the order the caller gave them.
func TestGeneratorKeepsSliceOrderWithoutZOrder(t *testing.T) {
	slide := elements.NewSlide("Paint order")
	slide.Shapes = []shapes.Shape{zShape("ZA", 0), zShape("ZB", 0)}

	parts := buildPackageParts(t, Metadata{}, []elements.SlideContent{slide})
	slideXML := parts["ppt/slides/slide1.xml"]

	if strings.Index(slideXML, "ZA") > strings.Index(slideXML, "ZB") {
		t.Fatal("slice order should be preserved when no ZOrder is set")
	}
}
