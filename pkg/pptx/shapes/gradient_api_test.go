package shapes_test

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

func TestCreateWithSlidesRendersLinearShapeGradientFill(t *testing.T) {
	gradient := pptx.NewShapeGradientFill(
		pptx.ShapeGradientTypeLinear,
		[]pptx.ShapeGradientStop{
			pptx.NewShapeGradientStop(0, "1F4E78"),
			pptx.NewShapeGradientStop(100, "8FB9E0").WithTransparency(30),
		},
	).WithLinearAngle(45)

	shape := pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(1.1), pptx.Inches(1.2), pptx.Inches(3), pptx.Inches(1.4)).
		WithGradientFill(gradient).
		WithText("Gradient")

	data, err := pptx.CreateWithSlides("Deck", []pptx.SlideContent{
		pptx.NewSlide("").WithBlankLayout().AddShape(shape),
	})
	if err != nil {
		t.Fatalf("CreateWithSlides returned error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	slideXML := testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")

	checks := []string{
		`<a:gradFill rotWithShape="1">`,
		`<a:gs pos="0">`,
		`<a:srgbClr val="1F4E78">`,
		`<a:gs pos="100000">`,
		`<a:srgbClr val="8FB9E0"><a:alpha val="70000"/></a:srgbClr>`,
		`<a:lin ang="2700000" scaled="1"/>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide XML", needle)
		}
	}
}

func TestCreateWithSlidesRendersRadialShapeGradientFill(t *testing.T) {
	gradient := pptx.NewShapeGradientFill(
		pptx.ShapeGradientTypeRadial,
		[]pptx.ShapeGradientStop{
			pptx.NewShapeGradientStop(0, "FFFFFF"),
			pptx.NewShapeGradientStop(100, "4472C4"),
		},
	)

	shape := pptx.NewShape(pptx.ShapeTypeEllipse, pptx.Inches(2), pptx.Inches(2), pptx.Inches(2), pptx.Inches(2)).
		WithGradientFill(gradient)

	data, err := pptx.CreateWithSlides("Deck", []pptx.SlideContent{
		pptx.NewSlide("").WithBlankLayout().AddShape(shape),
	})
	if err != nil {
		t.Fatalf("CreateWithSlides returned error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	slideXML := testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")

	if !strings.Contains(slideXML, `<a:path path="circle"/>`) {
		t.Fatalf("expected radial gradient path xml in slide output")
	}
}

func TestCreateWithSlidesRejectsInvalidShapeGradientStops(t *testing.T) {
	gradient := pptx.NewShapeGradientFill(
		pptx.ShapeGradientTypeLinear,
		[]pptx.ShapeGradientStop{
			pptx.NewShapeGradientStop(50, "4472C4"),
			pptx.NewShapeGradientStop(30, "8FB9E0"),
		},
	)

	slide := pptx.NewSlide("").WithBlankLayout().
		AddShape(pptx.NewShape(pptx.ShapeTypeRectangle, 1, 1, 1, 1).WithGradientFill(gradient))
	_, err := pptx.CreateWithSlides("Deck", []pptx.SlideContent{slide})
	if err == nil {
		t.Fatalf("expected gradient validation error")
	}
	if !strings.Contains(err.Error(), "gradient stop positions must be strictly increasing") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateWithSlidesRejectsGradientAngleOnNonLinearFill(t *testing.T) {
	gradient := pptx.NewShapeGradientFill(
		pptx.ShapeGradientTypeRadial,
		[]pptx.ShapeGradientStop{
			pptx.NewShapeGradientStop(0, "4472C4"),
			pptx.NewShapeGradientStop(100, "8FB9E0"),
		},
	).WithLinearAngle(90)

	slide := pptx.NewSlide("").WithBlankLayout().
		AddShape(pptx.NewShape(pptx.ShapeTypeRectangle, 1, 1, 1, 1).WithGradientFill(gradient))
	_, err := pptx.CreateWithSlides("Deck", []pptx.SlideContent{slide})
	if err == nil {
		t.Fatalf("expected gradient angle validation error")
	}
	if !strings.Contains(err.Error(), "gradient angle is only supported for linear gradients") {
		t.Fatalf("unexpected error: %v", err)
	}
}
