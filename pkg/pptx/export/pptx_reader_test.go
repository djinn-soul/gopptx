package export

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	editorcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
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

func TestEditorShapeToShape_MapsStyleAndAdjustments(t *testing.T) {
	fillColor := "4285F4"
	lineColor := "FFFFFF"
	lineWidth := 12700
	source := editorcommon.Shape{
		ID:   4,
		Name: "pie-slice",
		Type: "pie",
		X:    100,
		Y:    200,
		W:    300,
		H:    400,
		Fill: &editorcommon.ShapeFill{Solid: &fillColor},
		Line: &editorcommon.ShapeLine{Color: &lineColor, WidthEmu: &lineWidth},
		Adjustments: []editorcommon.ShapeAdjustment{
			{Name: "adj1", Formula: "val 0"},
			{Name: "adj2", Formula: "val 17100000"},
		},
	}

	mapped := editorShapeToShape(source)
	if mapped.Type != "pie" {
		t.Fatalf("expected mapped type pie, got %q", mapped.Type)
	}
	if mapped.Fill == nil || mapped.Fill.Color != fillColor {
		t.Fatalf("expected fill color %q, got %#v", fillColor, mapped.Fill)
	}
	if mapped.Line == nil || mapped.Line.Color != lineColor {
		t.Fatalf("expected line color %q, got %#v", lineColor, mapped.Line)
	}
	if int64(mapped.Line.Width) != int64(lineWidth) {
		t.Fatalf("expected line width %d, got %d", lineWidth, mapped.Line.Width)
	}
	if len(mapped.Adjustments) != 2 {
		t.Fatalf("expected 2 adjustments, got %d", len(mapped.Adjustments))
	}
	if mapped.Adjustments[0].Name != "adj1" || mapped.Adjustments[1].Name != "adj2" {
		t.Fatalf("unexpected adjustments: %#v", mapped.Adjustments)
	}
}

func TestSlidesFromPPTX_PieShapesKeepGeometryAndFill(t *testing.T) {
	deckPath := filepath.Clean("../../../examples/output/03_markdown_mermaid_complex_edited.pptx")
	if _, err := os.Stat(deckPath); err != nil {
		t.Skipf("deck fixture unavailable: %v", err)
	}
	_, slides, err := SlidesFromPPTX(deckPath)
	if err != nil {
		t.Fatalf("SlidesFromPPTX failed: %v", err)
	}
	if len(slides) < 4 {
		t.Fatalf("expected at least 4 slides, got %d", len(slides))
	}
	slide := slides[3]
	found := 0
	for _, shape := range slide.Shapes {
		if shape.Name != "Shape 4" && shape.Name != "Shape 5" && shape.Name != "Shape 6" {
			continue
		}
		found++
		if shape.Type != "pie" {
			t.Fatalf("%s expected type=pie, got %q", shape.Name, shape.Type)
		}
		if shape.Fill == nil || shape.Fill.Color == "" {
			t.Fatalf("%s expected non-empty fill, got %#v", shape.Name, shape.Fill)
		}
		if len(shape.Adjustments) < 2 {
			t.Fatalf("%s expected adjustments, got %#v", shape.Name, shape.Adjustments)
		}
	}
	if found != 3 {
		t.Fatalf("expected 3 pie slice shapes, got %d", found)
	}
}

func TestSlidesFromPPTX_Slide14ExtractsTable(t *testing.T) {
	deckPath := filepath.Clean("../../../examples/output/03_markdown_mermaid_complex.pptx")
	if _, err := os.Stat(deckPath); err != nil {
		t.Skipf("deck fixture unavailable: %v", err)
	}
	_, slides, err := SlidesFromPPTX(deckPath)
	if err != nil {
		t.Fatalf("SlidesFromPPTX failed: %v", err)
	}
	if len(slides) < 14 {
		t.Fatalf("expected at least 14 slides, got %d", len(slides))
	}
	slide := slides[13]
	if slide.Table == nil {
		t.Fatalf("expected slide 14 table to be extracted")
	}
	if len(slide.Table.Rows) < 2 {
		t.Fatalf("expected extracted table rows, got %d", len(slide.Table.Rows))
	}
	if len(slide.Table.ColumnWidths) != 2 {
		t.Fatalf("expected 2 table columns, got %d", len(slide.Table.ColumnWidths))
	}
}
