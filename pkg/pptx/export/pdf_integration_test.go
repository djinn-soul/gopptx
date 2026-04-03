package export_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/export"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

func TestPDFIntegration(t *testing.T) {
	// Setup complex presentation
	slides := []elements.SlideContent{
		elements.NewSlide("PDF Integration Test").
			WithBackgroundColor("#EFEFEF").
			AddBullet("First bullet").
			AddBullet("Second bullet"),

		elements.NewSlide("Shapes and Images").
			AddShape(shapes.NewShape(shapes.ShapeTypeRectangle, 100, 100, 200, 100).
				WithFill(shapes.NewShapeFill("FF0000")).
				WithText("Rectangle")).
			AddShape(shapes.NewShape(shapes.ShapeTypeEllipse, 400, 100, 100, 100).
				WithGradientFill(shapes.NewShapeGradientFill(shapes.ShapeGradientTypeLinear, []shapes.ShapeGradientStop{
					shapes.NewShapeGradientStop(0, "00FF00"),
					shapes.NewShapeGradientStop(100, "0000FF"),
				}))).
			AddShape(shapes.NewShape(shapes.ShapeTypeRectangle, 100, 500, 100, 100).
				WithFill(shapes.NewShapeFill("00FF00")).
				WithSoftEdges(true).
				WithShadow(true).
				WithGlow(true).
				WithReflection(true)).
			AddConnector(shapes.NewStraightConnector(styling.Inches(1), styling.Inches(1), styling.Inches(2), styling.Inches(2))).
			AddConnector(shapes.NewElbowConnector(styling.Inches(3), styling.Inches(3), styling.Inches(4), styling.Inches(4)).
				WithArrows(shapes.ArrowTypeStealth, shapes.ArrowTypeTriangle)).
			AddImage(shapes.Image{
				Data:   testutil.TinyPNG(),
				Format: "png",
				X:      400,
				Y:      300,
				CX:     200,
				CY:     200,
			}),

		elements.NewSlide("Chart Slide").
			WithBarChart(charts.BarChart{
				Categories: []string{"A", "B"},
				Values:     []float64{10, 20},
			}),

		elements.NewSlide("Gradient BG Slide").
			WithGradientBackground(shapes.ShapeGradientFill{
				Type: shapes.ShapeGradientTypeLinear,
				Stops: []shapes.ShapeGradientStop{
					shapes.NewShapeGradientStop(0, "FF0000"),
					shapes.NewShapeGradientStop(100, "0000FF"),
				},
			}),

		elements.NewSlide("Table Slide").
			WithTable(tables.NewTable([]styling.Length{styling.Inches(1), styling.Inches(1)}).
				AddRow([]string{"Header 1", "Header 2"}).
				AddRow([]string{"Cell 1-1", "Cell 1-2"})),

		elements.NewSlide("SmartArt Slide").
			AddSmartArt(
				smartart.NewSmartArt(smartart.Hierarchy).
					Position(styling.Inches(1), styling.Inches(1.2)).
					Size(styling.Inches(6.5), styling.Inches(3.5)).
					AddNode(
						smartart.NewNode("CEO").
							WithChild(smartart.NewNode("Finance")).
							WithChild(smartart.NewNode("Engineering").WithChild(smartart.NewNode("Platform"))),
					),
			),
	}

	tmpDir := t.TempDir()
	pdfPath := filepath.Join(tmpDir, "test.pdf")

	// Test native export
	opts := export.PDFOptions{Driver: export.PDFDriverNative}
	err := export.PDFWithOptions("PDF Test", slides, pdfPath, opts)
	if err != nil {
		t.Fatalf("PDFWithOptions (Native) failed: %v", err)
	}

	data, err := os.ReadFile(pdfPath)
	if err != nil {
		t.Fatalf("Failed to read generated PDF: %v", err)
	}

	if len(data) < 100 {
		t.Errorf("Generated PDF is suspiciously small: %d bytes", len(data))
	}

	// Verify PDF header
	if !bytes.HasPrefix(data, []byte("%PDF-")) {
		t.Error("Output does not have PDF header")
	}
}
