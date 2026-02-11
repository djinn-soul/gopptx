package shapes_test

import (
	"archive/zip"
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

var shapeRunSizePattern = regexp.MustCompile(`<a:rPr[^>]*sz="([0-9]+)"`)

func TestCreateWithSlidesAutoFitsShapeTextSize(t *testing.T) {
	short := pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(1), pptx.Inches(1), pptx.Inches(3), pptx.Inches(1.1)).
		WithFill(pptx.NewShapeFill("4472C4")).
		WithText("Short title")
	long := pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(1), pptx.Inches(3), pptx.Inches(3), pptx.Inches(1.1)).
		WithFill(pptx.NewShapeFill("4472C4")).
		WithText("This is a much longer sentence that should trigger smaller auto-fit text sizing")

	data, err := pptx.CreateWithSlides("Deck", []pptx.SlideContent{
		pptx.NewSlide("").WithBlankLayout().AddShape(short).AddShape(long),
	})
	if err != nil {
		t.Fatalf("CreateWithSlides returned error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	slideXML := testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")

	if !strings.Contains(slideXML, `<a:spAutoFit/>`) {
		t.Fatalf("expected shape text body to include auto-fit tag")
	}

	matches := shapeRunSizePattern.FindAllStringSubmatch(slideXML, -1)
	if len(matches) < 2 {
		t.Fatalf("expected at least two shape text size runs, got %d", len(matches))
	}

	shortSize, err := strconv.Atoi(matches[0][1])
	if err != nil {
		t.Fatalf("parse short text size: %v", err)
	}
	longSize, err := strconv.Atoi(matches[1][1])
	if err != nil {
		t.Fatalf("parse long text size: %v", err)
	}
	if longSize >= shortSize {
		t.Fatalf("expected long text size < short text size, got long=%d short=%d", longSize, shortSize)
	}
}
