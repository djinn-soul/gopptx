package pptx

import (
	"archive/zip"
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestAccessibility(t *testing.T) {
	// 1. Create a presentation with accessible elements
	const (
		shapeAltText = "A red rectangle"
		imageAltText = "A sample image"
	)

	// Create a dummy image - marked as decorative
	imgData := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A} // PNG header
	img := NewImageFromBytes(imgData, "png", 0, 0, 100, 100).
		WithAltText(imageAltText).
		WithDecorative(true)

	// Create a shape - with alt text
	shape := NewRectangle(1, 1, 2, 2).
		WithAltText(shapeAltText).
		WithDecorative(false)

	pb := NewPresentationBuilder("Access Test").
		AddSlide(NewSlide("Slide 1").
			AddShape(shape).
			AddImage(img))

	// 2. Build the presentation
	pptxBytes, err := pb.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// 3. Inspect the XML
	zr, err := zip.NewReader(bytes.NewReader(pptxBytes), int64(len(pptxBytes)))
	if err != nil {
		t.Fatalf("Failed to create zip reader: %v", err)
	}

	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")

	// 4. Verify Shape AltText (Not decorative, has AltText)
	if !strings.Contains(slideXML, fmt.Sprintf(`descr="%s"`, shapeAltText)) {
		t.Errorf("Shape AltText not found in XML. Expected descr=%q", shapeAltText)
	}

	// 5. Verify Image is Decorative (Has descr="")
	if !strings.Contains(slideXML, `descr=""`) {
		t.Errorf("Decorative image did not have descr=\"\" in XML")
	}
	if strings.Contains(slideXML, fmt.Sprintf(`descr="%s"`, imageAltText)) {
		t.Errorf("Decorative image should not have AltText in XML")
	}

	// 6. Verify common structure properties (just to be sure we are looking at cNvPr)
	// We expect something like <p:cNvPr id="..." name="..." descr="..."/>
	// Check for the presence of descr attribute in general context if strict check fails?
	// The strings.Contains check above is sufficient for now.
}
