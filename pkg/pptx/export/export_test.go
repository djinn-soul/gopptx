package export_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/export"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func TestExportHTML(t *testing.T) {
	// Create a temporary image for testing
	tmpDir := t.TempDir()
	imgPath := filepath.Join(tmpDir, "test.png")
	// Minimal valid PNG
	pngData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
		0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
		0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00,
		0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
		0x42, 0x60, 0x82,
	}
	if err := os.WriteFile(imgPath, pngData, 0o600); err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}

	slides := []elements.SlideContent{
		elements.NewSlide("Slide 1").
			AddBullet("Bullet point 1").
			AddBullet("Bullet point 2"),
		elements.NewSlide("Slide 2").
			AddShape(shapes.NewShape(shapes.ShapeTypeRectangle, 0, 0, 100, 100).WithText("Shape Text")),
		elements.NewSlide("Slide 3").
			AddImage(shapes.Image{Path: imgPath}),
	}

	html := export.HTML("Test Presentation", slides)

	// Verify structure
	if !strings.Contains(html, "<title>Test Presentation</title>") {
		t.Error("HTML missing title")
	}
	if !strings.Contains(html, "<h1>Test Presentation</h1>") {
		t.Error("HTML missing h1 title")
	}

	// Slide 1
	if !strings.Contains(html, "<h2>Slide 1</h2>") {
		t.Error("HTML missing Slide 1 title")
	}
	if !strings.Contains(html, "<li>Bullet point 1</li>") {
		t.Error("HTML missing bullet 1")
	}

	// Slide 2
	if !strings.Contains(html, "<h2>Slide 2</h2>") {
		t.Error("HTML missing Slide 2 title")
	}
	if !strings.Contains(html, "<p>Shape Text</p>") {
		t.Error("HTML missing shape text")
	}

	// Slide 3
	if !strings.Contains(html, "data:image/png;base64,") {
		t.Error("HTML missing base64 image data")
	}
}
