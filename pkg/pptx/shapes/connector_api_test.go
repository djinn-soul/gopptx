package shapes_test

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

func TestCreateWithSlidesRendersCustomShape(t *testing.T) {
	shape := pptx.NewShape(pptx.ShapeTypeRoundedRectangle, pptx.Inches(1.2), pptx.Inches(1.4), pptx.Inches(2.8), pptx.Inches(1.3)).
		WithFill(pptx.NewShapeFill("4472C4").WithTransparency(0.25)).
		WithLine(pptx.NewShapeLine("1F4E78", pptx.Points(1.5)).WithDash(pptx.LineDashDash)).
		WithText("Service")

	slide := pptx.NewSlide("").WithBlankLayout().AddShape(shape)
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
	left := pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(1), pptx.Inches(2), pptx.Inches(2), pptx.Inches(1)).WithText("Start")
	right := pptx.NewShape(pptx.ShapeTypeDiamond, pptx.Inches(5), pptx.Inches(2), pptx.Inches(2), pptx.Inches(1)).WithText("Decision")
	connector := pptx.NewStraightConnector(pptx.Inches(3), pptx.Inches(2.5), pptx.Inches(5), pptx.Inches(2.5)).
		WithLine(
			pptx.NewShapeLine("5B9BD5", pptx.Points(1.25)).
				WithDash(pptx.LineDashDash).
				WithCap(pptx.LineCapSquare).
				WithJoin(pptx.LineJoinBevel),
		).
		WithArrows(pptx.ArrowTypeNone, pptx.ArrowTypeTriangle).
		WithArrowSize(pptx.ArrowSizeLarge).
		ConnectStart(1, pptx.ConnectionSiteRight).
		ConnectEnd(2, pptx.ConnectionSiteLeft).
		WithLabel("next")

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
		`<p:cxnSp>`,
		`<a:stCxn id="2" idx="1"/>`,
		`<a:endCxn id="3" idx="3"/>`,
		`<a:prstGeom prst="straightConnector1"><a:avLst/></a:prstGeom>`,
		`<a:ln w="15875" cap="sq">`,
		`<a:prstDash val="dash"/>`,
		`<a:bevel/>`,
		`<a:tailEnd type="triangle" w="lg" len="lg"/>`,
		`<a:t>next</a:t>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide XML", needle)
		}
	}
}

func TestCreateWithSlidesRendersConnectorAdjustmentPoints(t *testing.T) {
	connector := pptx.NewElbowConnector(pptx.Inches(1), pptx.Inches(1), pptx.Inches(4), pptx.Inches(3)).
		WithAdjustmentValue("adj1", 25000).
		WithAdjustment("adj2", "val 50000")

	slide := pptx.NewSlide("").WithBlankLayout().AddConnector(connector)
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
		`<a:prstGeom prst="bentConnector3"><a:avLst><a:gd name="adj1" fmla="val 25000"/><a:gd name="adj2" fmla="val 50000"/></a:avLst></a:prstGeom>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide XML", needle)
		}
	}
}

func TestCreateWithSlidesRejectsUnsupportedShapeType(t *testing.T) {
	slide := pptx.NewSlide("").WithBlankLayout().AddShape(pptx.NewShape("doesNotExist", 1, 1, 1, 1))
	_, err := pptx.CreateWithSlides("Deck", []pptx.SlideContent{slide})
	if err == nil {
		t.Fatalf("expected shape validation error")
	}
	if !strings.Contains(err.Error(), "shape 1 type") {
		t.Fatalf("expected shape type validation error, got %v", err)
	}
}

func TestCreateWithSlidesRejectsConnectorAnchorOutOfRange(t *testing.T) {
	connector := pptx.NewStraightConnector(1, 1, 2, 2).ConnectEnd(1, pptx.ConnectionSiteLeft)
	slide := pptx.NewSlide("").WithBlankLayout().AddConnector(connector)
	_, err := pptx.CreateWithSlides("Deck", []pptx.SlideContent{slide})
	if err == nil {
		t.Fatalf("expected connector validation error")
	}
	if !strings.Contains(err.Error(), "end shape index 1 is out of range") {
		t.Fatalf("expected connector anchor validation error, got %v", err)
	}
}

func TestUnitsHelpers(t *testing.T) {
	if pptx.Inches(1) != 914400 {
		t.Fatalf("expected 1 inch to equal 914400 EMU, got %d", pptx.Inches(1))
	}
	if pptx.Centimeters(2.54) != 914400 {
		t.Fatalf("expected 2.54 cm to equal 914400 EMU, got %d", pptx.Centimeters(2.54))
	}
	if pptx.Points(1) != 12700 {
		t.Fatalf("expected 1 pt to equal 12700 EMU, got %d", pptx.Points(1))
	}
}
