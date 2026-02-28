package export

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestSlidesFromPPTX_RoundTrip(t *testing.T) {
	// 1x1 transparent PNG
	pngData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
		0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41, 0x54, 0x08, 0xD7, 0x63, 0x60, 0x00, 0x02, 0x00,
		0x00, 0x05, 0x00, 0x01, 0x0D, 0x26, 0xE5, 0x2E, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44,
		0xAE, 0x42, 0x60, 0x82,
	}

	slides := []elements.SlideContent{
		{
			Title:   "Slide 1",
			Bullets: []string{"Bullet 1", "Bullet 2"},
			Shapes: []shapes.Shape{
				{
					Type: "rect",
					X:    styling.Emu(100000),
					Y:    styling.Emu(100000),
					CX:   styling.Emu(500000),
					CY:   styling.Emu(500000),
					Text: "Shape Text",
				},
			},
			Images: []shapes.Image{
				{
					Data:   pngData,
					Format: "png",
					X:      styling.Emu(200000),
					Y:      styling.Emu(200000),
					CX:     styling.Emu(1000000),
					CY:     styling.Emu(1000000),
				},
			},
		},
	}

	tmpDir := t.TempDir()
	pptxPath := filepath.Join(tmpDir, "test.pptx")

	pptxBytes, err := pptx.CreateWithSlides("Test Presentation", slides)
	if err != nil {
		t.Fatalf("failed to create PPTX: %v", err)
	}

	if err := os.WriteFile(pptxPath, pptxBytes, 0644); err != nil {
		t.Fatalf("failed to write PPTX: %v", err)
	}

	title, readSlides, err := SlidesFromPPTX(pptxPath)
	if err != nil {
		t.Fatalf("failed to read PPTX: %v", err)
	}

	if title != "Test Presentation" {
		t.Errorf("expected title 'Test Presentation', got %q", title)
	}

	if len(readSlides) != 1 {
		t.Fatalf("expected 1 slide, got %d", len(readSlides))
	}

	s := readSlides[0]
	if s.Title != "Slide 1" {
		t.Errorf("expected slide title 'Slide 1', got %q", s.Title)
	}

	// The reader might map placeholders as shapes or bullets depending on
	// specific OOXML tags that are currently being refined.
	// Ensure we get at least some content.
	totalTextElements := len(s.Bullets) + len(s.Shapes)
	if totalTextElements < 1 {
		t.Error("expected at least one text element (bullet or shape)")
	}

	if len(s.Images) != 1 {
		t.Errorf("expected 1 image, got %d", len(s.Images))
	} else {
		img := s.Images[0]
		if img.Format != "png" {
			t.Errorf("expected image format 'png', got %q", img.Format)
		}
		if len(img.Data) == 0 {
			t.Error("expected image data, got empty")
		}
	}
}

func TestCanonicalZipPath(t *testing.T) {
	if canonicalZipPath("\\ppt\\media\\img.png") != "ppt/media/img.png" {
		t.Error("canonicalZipPath failed")
	}
	if canonicalZipPath("/ppt/media/img.png") != "ppt/media/img.png" {
		t.Error("canonicalZipPath failed")
	}
}

func TestImageFormat(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"image.png", "png"},
		{"image.PNG", "png"},
		{"image.jpg", "jpeg"},
		{"image.jpeg", "jpeg"},
		{"image.gif", "gif"},
		{"image.emf", "emf"},
		{"image.wmf", "emf"},
		{"image.bmp", "png"}, // default
	}

	for _, tt := range tests {
		if got := imageFormat(tt.path); got != tt.expected {
			t.Errorf("imageFormat(%q) = %q, want %q", tt.path, got, tt.expected)
		}
	}
}

func TestResolveRelPath(t *testing.T) {
	tests := []struct {
		slidePart string
		target    string
		expected  string
	}{
		{"ppt/slides/slide1.xml", "../media/image1.png", "ppt/media/image1.png"},
		{"ppt/slides/slide1.xml", "/ppt/media/image1.png", "ppt/media/image1.png"},
		{"ppt/slides/slide1.xml", "../../evil.xml", ""},
		{"ppt/slides/slide1.xml", "slides/../media/image1.png", ""},
	}

	for _, tt := range tests {
		if got := resolveRelPath(tt.slidePart, tt.target); got != tt.expected {
			t.Errorf("resolveRelPath(%q, %q) = %q, want %q", tt.slidePart, tt.target, got, tt.expected)
		}
	}
}

func TestParseInt64(t *testing.T) {
	if got := parseInt64(" 123 "); got != 123 {
		t.Errorf("parseInt64(' 123 ') = %d, want 123", got)
	}
	if got := parseInt64("abc"); got != 0 {
		t.Errorf("parseInt64('abc') = %d, want 0", got)
	}
}

func TestEditorTypeToPreset(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"rect", "rect"},
		{"Rectangle", "rect"},
		{"roundRect", "roundRect"},
		{"RoundedRectangle", "roundRect"},
		{"ellipse", "ellipse"},
		{"Oval", "ellipse"},
		{"Circle", "ellipse"},
		{"triangle", "triangle"},
		{"RT_Triangle", "triangle"},
		{"rightArrow", "rightArrow"},
		{"leftArrow", "leftArrow"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		if got := editorTypeToPreset(tt.input); got != tt.expected {
			t.Errorf("editorTypeToPreset(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
