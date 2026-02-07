package pptx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestCreateWithSlidesRendersConnectorAutoSitesHorizontal(t *testing.T) {
	left := NewShape(ShapeTypeRectangle, Inches(1), Inches(2), Inches(2), Inches(1)).WithText("Left")
	right := NewShape(ShapeTypeRectangle, Inches(5), Inches(2), Inches(2), Inches(1)).WithText("Right")
	connector := NewStraightConnector(Inches(3), Inches(2.5), Inches(5), Inches(2.5)).
		ConnectStartAuto(1).
		ConnectEndAuto(2)

	slide := NewSlide("").WithBlankLayout().
		AddShape(left).
		AddShape(right).
		AddConnector(connector)

	data, err := CreateWithSlides("Deck", []SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides returned error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")

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
	top := NewShape(ShapeTypeRectangle, Inches(3), Inches(1), Inches(2), Inches(1)).WithText("Top")
	bottom := NewShape(ShapeTypeRectangle, Inches(3), Inches(4), Inches(2), Inches(1)).WithText("Bottom")
	connector := NewStraightConnector(Inches(4), Inches(2), Inches(4), Inches(4)).
		ConnectStartAuto(1).
		ConnectEndAuto(2)

	slide := NewSlide("").WithBlankLayout().
		AddShape(top).
		AddShape(bottom).
		AddConnector(connector)

	data, err := CreateWithSlides("Deck", []SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides returned error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")

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
