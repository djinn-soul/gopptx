package shapes_test

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

func TestCreateWithSlidesRendersConnectorAutoSitesHorizontal(t *testing.T) {
	left := pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(1), pptx.Inches(2), pptx.Inches(2), pptx.Inches(1)).WithText("Left")
	right := pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(5), pptx.Inches(2), pptx.Inches(2), pptx.Inches(1)).WithText("Right")
	connector := pptx.NewStraightConnector(pptx.Inches(3), pptx.Inches(2.5), pptx.Inches(5), pptx.Inches(2.5)).
		ConnectStartAuto(1).
		ConnectEndAuto(2)

	slide := pptx.NewSlide("").WithBlankLayout().
		AddShape(left).
		AddShape(right).
		AddConnector(connector)

	data, err := pptx.CreateWithSlides("Deck", []pptx.SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides returned error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	slideXML := testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")

	checks := []string{
		`<a:stCxn id="2" idx="1"/>`,
		`<a:endCxn id="3" idx="3"/>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide XML", needle)
		}
	}
}

func TestCreateWithSlidesRendersConnectorAutoSitesVertical(t *testing.T) {
	top := pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(3), pptx.Inches(1), pptx.Inches(2), pptx.Inches(1)).WithText("Top")
	bottom := pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(3), pptx.Inches(4), pptx.Inches(2), pptx.Inches(1)).WithText("Bottom")
	connector := pptx.NewStraightConnector(pptx.Inches(4), pptx.Inches(2), pptx.Inches(4), pptx.Inches(4)).
		ConnectStartAuto(1).
		ConnectEndAuto(2)

	slide := pptx.NewSlide("").WithBlankLayout().
		AddShape(top).
		AddShape(bottom).
		AddConnector(connector)

	data, err := pptx.CreateWithSlides("Deck", []pptx.SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides returned error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	slideXML := testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")

	checks := []string{
		`<a:stCxn id="2" idx="2"/>`,
		`<a:endCxn id="3" idx="0"/>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide XML", needle)
		}
	}
}

func TestCreateWithSlidesRendersConnectorAutoSitesDiagonal(t *testing.T) {
	start := pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(1), pptx.Inches(1), pptx.Inches(2), pptx.Inches(1)).WithText("Start")
	end := pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(5), pptx.Inches(4), pptx.Inches(2), pptx.Inches(1)).WithText("End")
	connector := pptx.NewStraightConnector(pptx.Inches(3), pptx.Inches(2), pptx.Inches(5), pptx.Inches(4)).
		ConnectStartAuto(1).
		ConnectEndAuto(2)

	slide := pptx.NewSlide("").WithBlankLayout().
		AddShape(start).
		AddShape(end).
		AddConnector(connector)

	data, err := pptx.CreateWithSlides("Deck", []pptx.SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides returned error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	slideXML := testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")

	checks := []string{
		`<a:stCxn id="2" idx="6"/>`,
		`<a:endCxn id="3" idx="4"/>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide XML", needle)
		}
	}
}
