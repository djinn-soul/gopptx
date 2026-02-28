package export

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestRenderShapesSVG_Rect(t *testing.T) {
	s := shapes.NewShape(shapes.ShapeTypeRectangle, styling.Inches(1), styling.Inches(1), styling.Inches(2), styling.Inches(2)).
		WithFill(shapes.NewShapeFill("FF0000"))

	out := renderShapesSVG([]shapes.Shape{s})

	if !strings.Contains(out, "<svg") {
		t.Error("missing <svg tag")
	}
	if !strings.Contains(out, "<rect") {
		t.Error("missing <rect tag")
	}
	if !strings.Contains(out, `fill="#FF0000"`) {
		t.Error("missing fill in rect")
	}
	if strings.Contains(out, "NaN") {
		t.Error("rendered NaN coordinates")
	}
}

func TestRenderShapesSVG_Ellipse(t *testing.T) {
	s := shapes.NewShape(shapes.ShapeTypeEllipse, styling.Inches(1), styling.Inches(1), styling.Inches(2), styling.Inches(2)).
		WithFill(shapes.NewShapeFill("00FF00")).
		WithText("Hello World")

	out := renderShapesSVG([]shapes.Shape{s})

	if !strings.Contains(out, "<ellipse") {
		t.Error("missing <ellipse tag")
	}
	if !strings.Contains(out, "Hello World</text>") {
		t.Error("missing text rendering")
	}
}

func TestRenderShapesSVG_Gradient(t *testing.T) {
	grad := shapes.NewShapeGradientFill(shapes.ShapeGradientTypeLinear, []shapes.ShapeGradientStop{
		shapes.NewShapeGradientStop(0, "FF0000"),
		shapes.NewShapeGradientStop(100, "0000FF"),
	})
	s := shapes.NewShape(shapes.ShapeTypeRectangle, styling.Inches(1), styling.Inches(1), styling.Inches(2), styling.Inches(2)).
		WithGradientFill(grad)

	out := renderShapesSVG([]shapes.Shape{s})

	if !strings.Contains(out, "<linearGradient") {
		t.Error("missing <linearGradient tag for gradient fill")
	}
	if !strings.Contains(out, `stop-color="#FF0000"`) {
		t.Error("missing gradient stop color")
	}
	if !strings.Contains(out, "url(#grad-") {
		t.Error("missing url(#grad-...) fill reference in parsed shape")
	}
}

func TestRenderShapesSVG_RotationTransform(t *testing.T) {
	s := shapes.NewShape(shapes.ShapeTypeRectangle, styling.Inches(1), styling.Inches(1), styling.Inches(2), styling.Inches(2)).
		WithRotation(45)

	out := renderShapesSVG([]shapes.Shape{s})

	if !strings.Contains(out, `transform="rotate(45,`) {
		t.Error("missing transform rotation")
	}
}
