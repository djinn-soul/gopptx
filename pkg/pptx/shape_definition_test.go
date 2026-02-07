package pptx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

type testBadgeShape struct {
	label string
}

func (b testBadgeShape) ToShape() Shape {
	return NewShape(ShapeTypeRoundedRectangle, Inches(1), Inches(1.4), Inches(2.6), Inches(1)).
		WithFill(NewShapeFill("2F5597")).
		WithText(b.label)
}

func TestCreateWithSlidesAcceptsShapeDefinition(t *testing.T) {
	slide := NewSlide("").WithBlankLayout().AddShape(testBadgeShape{label: "Interface"})

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
		`<a:srgbClr val="2F5597">`,
		`<a:t>Interface</a:t>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide XML", needle)
		}
	}
}
