package pptx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestCreateWithSlidesUsesLightTextOnDarkShapeFill(t *testing.T) {
	shape := NewShape(ShapeTypeRectangle, Inches(1), Inches(1), Inches(2.5), Inches(1.2)).
		WithFill(NewShapeFill("1F4E78")).
		WithText("Dark")

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
		`<a:t>Dark</a:t>`,
		`<a:srgbClr val="FFFFFF"/>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide XML", needle)
		}
	}
}

func TestCreateWithSlidesUsesDarkTextOnLightShapeFill(t *testing.T) {
	shape := NewShape(ShapeTypeRectangle, Inches(1), Inches(1), Inches(2.5), Inches(1.2)).
		WithFill(NewShapeFill("EAF2FB")).
		WithText("Light")

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
		`<a:t>Light</a:t>`,
		`<a:srgbClr val="000000"/>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide XML", needle)
		}
	}
}

func TestCreateWithSlidesUsesLightTextOnDarkGradientShapeFill(t *testing.T) {
	gradient := NewShapeGradientFill(
		ShapeGradientTypeLinear,
		[]ShapeGradientStop{
			NewShapeGradientStop(0, "173A5E"),
			NewShapeGradientStop(100, "2A5D8F"),
		},
	)
	shape := NewShape(ShapeTypeRectangle, Inches(1), Inches(1), Inches(2.5), Inches(1.2)).
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
		`<a:t>Gradient</a:t>`,
		`<a:srgbClr val="FFFFFF"/>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide XML", needle)
		}
	}
}
