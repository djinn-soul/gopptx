package pptx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestCreateWithSlidesRendersCustomShape(t *testing.T) {
	shape := NewShape(ShapeTypeRoundedRectangle, Inches(1.2), Inches(1.4), Inches(2.8), Inches(1.3)).
		WithFill(NewShapeFill("4472C4").WithTransparency(25)).
		WithLine(NewShapeLine("1F4E78", Points(1.5)).WithDash(LineDashDash)).
		WithText("Service")

	slide := NewSlide("").WithBlankLayout().AddShape(shape)
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
		`<a:prstGeom prst="roundRect"><a:avLst/></a:prstGeom>`,
		`<a:srgbClr val="4472C4"><a:alpha val="75000"/></a:srgbClr>`,
		`<a:ln w="19050">`,
		`<a:prstDash val="dash"/>`,
		`<a:t>Service</a:t>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide XML", needle)
		}
	}
}

func TestCreateWithSlidesRendersConnectorWithShapeAnchors(t *testing.T) {
	left := NewShape(ShapeTypeRectangle, Inches(1), Inches(2), Inches(2), Inches(1)).WithText("Start")
	right := NewShape(ShapeTypeDiamond, Inches(5), Inches(2), Inches(2), Inches(1)).WithText("Decision")
	connector := NewStraightConnector(Inches(3), Inches(2.5), Inches(5), Inches(2.5)).
		WithLine(NewShapeLine("5B9BD5", Points(1.25)).WithDash(LineDashDash)).
		WithArrows(ArrowTypeNone, ArrowTypeTriangle).
		WithArrowSize(ArrowSizeLarge).
		ConnectStart(1, ConnectionSiteRight).
		ConnectEnd(2, ConnectionSiteLeft).
		WithLabel("next")

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
		`<p:cxnSp>`,
		`<a:stCxn id="2" idx="1"/>`,
		`<a:endCxn id="3" idx="3"/>`,
		`<a:prstGeom prst="straightConnector1"><a:avLst/></a:prstGeom>`,
		`<a:prstDash val="dash"/>`,
		`<a:tailEnd type="triangle" w="lg" len="lg"/>`,
		`<a:t>next</a:t>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide XML", needle)
		}
	}
}

func TestCreateWithSlidesRejectsUnsupportedShapeType(t *testing.T) {
	slide := NewSlide("").WithBlankLayout().AddShape(NewShape("doesNotExist", 1, 1, 1, 1))
	_, err := CreateWithSlides("Deck", []SlideContent{slide})
	if err == nil {
		t.Fatalf("expected shape validation error")
	}
	if !strings.Contains(err.Error(), "shape 1 type") {
		t.Fatalf("expected shape type validation error, got %v", err)
	}
}

func TestCreateWithSlidesRejectsConnectorAnchorOutOfRange(t *testing.T) {
	connector := NewStraightConnector(1, 1, 2, 2).ConnectEnd(1, ConnectionSiteLeft)
	slide := NewSlide("").WithBlankLayout().AddConnector(connector)
	_, err := CreateWithSlides("Deck", []SlideContent{slide})
	if err == nil {
		t.Fatalf("expected connector validation error")
	}
	if !strings.Contains(err.Error(), "end shape index 1 is out of range") {
		t.Fatalf("expected connector anchor validation error, got %v", err)
	}
}

func TestUnitsHelpers(t *testing.T) {
	if Inches(1) != 914400 {
		t.Fatalf("expected 1 inch to equal 914400 EMU, got %d", Inches(1))
	}
	if Centimeters(2.54) != 914400 {
		t.Fatalf("expected 2.54 cm to equal 914400 EMU, got %d", Centimeters(2.54))
	}
	if Points(1) != 12700 {
		t.Fatalf("expected 1 pt to equal 12700 EMU, got %d", Points(1))
	}
}
