package shapes_test

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

type testBadgeShape struct {
	label string
}

func (b testBadgeShape) ToShape() pptx.Shape {
	return pptx.NewShape(pptx.ShapeTypeRoundedRectangle, pptx.Inches(1), pptx.Inches(1.4), pptx.Inches(2.6), pptx.Inches(1)).
		WithFill(pptx.NewShapeFill("2F5597")).
		WithText(b.label)
}

func TestCreateWithSlidesAcceptsShapeDefinition(t *testing.T) {
	slide := pptx.NewSlide("").WithBlankLayout().AddShape(testBadgeShape{label: "Interface"})

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
		`<a:srgbClr val="2F5597">`,
		`<a:t>Interface</a:t>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide XML", needle)
		}
	}
}
