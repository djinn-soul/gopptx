package pptx_test

import (
	"archive/zip"
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

func TestAccessibility(t *testing.T) {
	const (
		shapeAltText = "A red rectangle"
		imageAltText = "A sample image"
	)

	imgData := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	img := pptx.NewImageFromBytes(imgData, "png", 0, 0, 100, 100).
		WithAltText(imageAltText).
		WithDecorative(true)

	shape := pptx.NewRectangle(1, 1, 2, 2).
		WithAltText(shapeAltText).
		WithDecorative(false)

	pb := pptx.NewPresentationBuilder("Access Test").
		AddSlide(pptx.NewSlide("Slide 1").
			AddShape(shape).
			AddImage(img))

	pptxBytes, err := pb.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(pptxBytes), int64(len(pptxBytes)))
	if err != nil {
		t.Fatalf("Failed to create zip reader: %v", err)
	}

	slideXML := testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")

	if !strings.Contains(slideXML, fmt.Sprintf(`descr="%s"`, shapeAltText)) {
		t.Errorf("Shape AltText not found in XML. Expected descr=%q", shapeAltText)
	}

	if !strings.Contains(slideXML, `descr=""`) {
		t.Errorf("Decorative image did not have descr=\"\" in XML")
	}
	if strings.Contains(slideXML, fmt.Sprintf(`descr="%s"`, imageAltText)) {
		t.Errorf("Decorative image should not have AltText in XML")
	}
}
