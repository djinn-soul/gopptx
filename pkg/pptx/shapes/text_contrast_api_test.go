package shapes_test

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

func TestCreateWithSlidesUsesLightTextOnDarkShapeFill(t *testing.T) {
	shape := pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(1), pptx.Inches(1), pptx.Inches(2.5), pptx.Inches(1.2)).
		WithFill(pptx.NewShapeFill("1F4E78")).
		WithText("Dark")

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
	shape := pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(1), pptx.Inches(1), pptx.Inches(2.5), pptx.Inches(1.2)).
		WithFill(pptx.NewShapeFill("EAF2FB")).
		WithText("Light")

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
	gradient := pptx.NewShapeGradientFill(
		pptx.ShapeGradientTypeLinear,
		[]pptx.ShapeGradientStop{
			pptx.NewShapeGradientStop(0, "173A5E"),
			pptx.NewShapeGradientStop(100, "2A5D8F"),
		},
	)
	shape := pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(1), pptx.Inches(1), pptx.Inches(2.5), pptx.Inches(1.2)).
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
		`<a:t>Gradient</a:t>`,
		`<a:srgbClr val="FFFFFF"/>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide XML", needle)
		}
	}
}
