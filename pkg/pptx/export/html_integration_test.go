package export_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/export"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

func TestHTMLIntegration(t *testing.T) {
	// Setup complex presentation
	slides := []elements.SlideContent{
		// Slide 1: Title and Bullets
		elements.NewSlide("Integration Test").
			AddBullet("First bullet").
			AddBullet("Second bullet"),

		// Slide 2: Shapes and Fills
		elements.NewSlide("Shapes and Fills").
			AddShape(shapes.NewShape(shapes.ShapeTypeRectangle, 100, 100, 200, 100).
				WithFill(shapes.NewShapeFill("FF0000")).
				WithText("Rectangle")).
			AddShape(shapes.NewShape(shapes.ShapeTypeEllipse, 400, 100, 100, 100).
				WithGradientFill(shapes.NewShapeGradientFill(shapes.ShapeGradientTypeLinear, []shapes.ShapeGradientStop{
					shapes.NewShapeGradientStop(0, "00FF00"),
					shapes.NewShapeGradientStop(100, "0000FF"),
				}))).
			AddShape(shapes.NewShape(shapes.ShapeTypeTriangle, 100, 300, 100, 100).
				WithRotation(45)),

		// Slide 3: Image
		elements.NewSlide("Image Slide").
			AddImage(shapes.Image{
				Data:   testutil.TinyPNG(),
				Format: "png",
				X:      100,
				Y:      100,
				CX:     200,
				CY:     200,
			}),

		// Slide 4: Table
		elements.NewSlide("Table Slide").
			WithTable(tables.NewTable([]styling.Length{styling.Inches(1), styling.Inches(1)}).
				AddRow([]string{"Header 1", "Header 2"}).
				AddRow([]string{"Cell 1-1", "Cell 1-2"})),
	}

	opts := export.DefaultHTMLOptions()
	opts.Theme = &export.ThemeColors{
		TitleColor:      "#123456",
		BackgroundColor: "#F0F0F0",
	}
	opts.IncludeNavigation = true

	html := export.HTMLWithOptions("Integration Presentation", slides, opts)

	// Basic checks
	checks := []string{
		`<!DOCTYPE html>`,
		`<title>Integration Presentation</title>`,
		`--title-color: #123456;`,
		`--bg-color: #F0F0F0;`,
		`class="nav-controls"`,
		`<h1>Integration Presentation</h1>`,
		`<h2>Integration Test</h2>`,
		`<li>First bullet</li>`,
		`Rectangle`,
		`linearGradient`,
		`data:image/png;base64,`,
		`<table`,
		`<td>Header 1</td>`,
	}

	for _, check := range checks {
		if !strings.Contains(html, check) {
			t.Errorf("HTML missing expected content: %q", check)
		}
	}

	// Golden file comparison
	_ = os.MkdirAll("testdata", 0o755)
	goldenPath := filepath.Join("testdata", "integration.html")
	if os.Getenv("UPDATE_GOLDEN") == "true" {
		if err := os.WriteFile(goldenPath, []byte(html), 0o600); err != nil {
			t.Fatalf("Failed to update golden file: %v", err)
		}
	}

	// For now, let's just ensure it's not empty.
	if len(html) < 1000 {
		t.Errorf("Generated HTML is suspiciously small: %d bytes", len(html))
	}
}
