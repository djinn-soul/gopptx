package pptx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestCreateWithSlidesRendersLinearShapeGradientFill(t *testing.T) {
	gradient := NewShapeGradientFill(
		ShapeGradientTypeLinear,
		[]ShapeGradientStop{
			NewShapeGradientStop(0, "1F4E78"),
			NewShapeGradientStop(100, "8FB9E0").WithTransparency(30),
		},
	).WithLinearAngle(45)

	shape := NewShape(ShapeTypeRectangle, Inches(1.1), Inches(1.2), Inches(3), Inches(1.4)).
		WithGradientFill(gradient).
		WithText("Gradient")

	data, err := CreateWithSlides("Deck", []SlideContent{
		NewSlide("").WithBlankLayout().AddShape(shape),
	})
	if err != nil {
		t.Fatalf("CreateWithSlides returned error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")

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
	gradient := NewShapeGradientFill(
		ShapeGradientTypeRadial,
		[]ShapeGradientStop{
			NewShapeGradientStop(0, "FFFFFF"),
			NewShapeGradientStop(100, "4472C4"),
		},
	)

	shape := NewShape(ShapeTypeEllipse, Inches(2), Inches(2), Inches(2), Inches(2)).
		WithGradientFill(gradient)

	data, err := CreateWithSlides("Deck", []SlideContent{
		NewSlide("").WithBlankLayout().AddShape(shape),
	})
	if err != nil {
		t.Fatalf("CreateWithSlides returned error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")

	if !strings.Contains(slideXML, `<a:path path="circle"/>`) {
		t.Fatalf("expected radial gradient path xml in slide output")
	}
}

func TestCreateWithSlidesRejectsInvalidShapeGradientStops(t *testing.T) {
	gradient := NewShapeGradientFill(
		ShapeGradientTypeLinear,
		[]ShapeGradientStop{
			NewShapeGradientStop(50, "4472C4"),
			NewShapeGradientStop(30, "8FB9E0"),
		},
	)

	slide := NewSlide("").WithBlankLayout().
		AddShape(NewShape(ShapeTypeRectangle, 1, 1, 1, 1).WithGradientFill(gradient))
	_, err := CreateWithSlides("Deck", []SlideContent{slide})
	if err == nil {
		t.Fatalf("expected gradient validation error")
	}
	if !strings.Contains(err.Error(), "gradient stop positions must be strictly increasing") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateWithSlidesRejectsGradientAngleOnNonLinearFill(t *testing.T) {
	gradient := NewShapeGradientFill(
		ShapeGradientTypeRadial,
		[]ShapeGradientStop{
			NewShapeGradientStop(0, "4472C4"),
			NewShapeGradientStop(100, "8FB9E0"),
		},
	).WithLinearAngle(90)

	slide := NewSlide("").WithBlankLayout().
		AddShape(NewShape(ShapeTypeRectangle, 1, 1, 1, 1).WithGradientFill(gradient))
	_, err := CreateWithSlides("Deck", []SlideContent{slide})
	if err == nil {
		t.Fatalf("expected gradient angle validation error")
	}
	if !strings.Contains(err.Error(), "gradient angle is only supported for linear gradients") {
		t.Fatalf("unexpected error: %v", err)
	}
}
