package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/export"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

func main() {
	outDir := filepath.Join("examples", "output")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatalf("failed to create output dir: %v", err)
	}

	htmlPath := filepath.Join(outDir, "21_export_demo.html")
	pdfPath := filepath.Join(outDir, "21_export_demo.pdf")

	slides := make([]elements.SlideContent, 0, 3)

	// Slide 1: Welcome and basic shapes
	s1 := elements.NewSlide("High-Fidelity Export")
	s1.Bullets = []string{
		"This presentation demonstrates gopptx's enhanced export capabilities.",
		"Features include SVG shape rendering, pure CSS styling, and headless PDF generation.",
	}

	// Add a rounded rect with solid fill
	rect := shapes.NewShape(shapes.ShapeTypeRoundedRectangle, styling.Inches(1), styling.Inches(3), styling.Inches(2), styling.Inches(1)).
		WithFill(shapes.NewShapeFill("0078D4")).
		WithText("SVG Rect").
		WithLine(shapes.NewShapeLine("000000", styling.Points(2)))
	s1.Shapes = append(s1.Shapes, rect)

	slides = append(slides, s1)

	// Slide 2: Gradients and Table
	s2 := elements.NewSlide("Gradients and Tables")

	// Add a gradient ellipse
	grad := shapes.NewShapeGradientFill(shapes.ShapeGradientTypeLinear, []shapes.ShapeGradientStop{
		shapes.NewShapeGradientStop(0, "FF5555"),
		shapes.NewShapeGradientStop(100, "5555FF"),
	})
	ellipse := shapes.NewShape(shapes.ShapeTypeEllipse, styling.Inches(7), styling.Inches(1), styling.Inches(2), styling.Inches(2)).
		WithGradientFill(grad).
		WithText("Gradient")
	s2.Shapes = append(s2.Shapes, ellipse)

	// Add a small table
	tab := tables.NewTable([]styling.Length{styling.Inches(3), styling.Inches(3)}).
		Position(styling.Inches(1), styling.Inches(2))

	header1 := tables.TableCell{Text: "Feature", Bold: true, BackgroundColor: "#0078D4", Color: "#FFFFFF"}
	header2 := tables.TableCell{Text: "Status", Bold: true, BackgroundColor: "#0078D4", Color: "#FFFFFF"}
	tab = tab.AddStyledRow([]tables.TableCell{header1, header2})

	row1c1 := tables.TableCell{Text: "HTML SVG Export"}
	row1c2 := tables.TableCell{Text: "Done"}
	tab = tab.AddStyledRow([]tables.TableCell{row1c1, row1c2})

	row2c1 := tables.TableCell{Text: "Headless PDF Chrome"}
	row2c2 := tables.TableCell{Text: "Done"}
	tab = tab.AddStyledRow([]tables.TableCell{row2c1, row2c2})

	s2.Table = &tab
	slides = append(slides, s2)

	// Slide 3: Transforms and Ending
	s3 := elements.NewSlide("Rotations & Custom Options")

	// A rotated arrow
	arrow := shapes.NewShape(shapes.ShapeTypeRightArrow, styling.Inches(4), styling.Inches(2), styling.Inches(2), styling.Inches(1)).
		WithFill(shapes.NewShapeFill("28A745")).
		WithRotation(45).
		WithText("Rotated")
	s3.Shapes = append(s3.Shapes, arrow)

	slides = append(slides, s3)

	fmt.Println("Generating HTML export...")
	opts := export.DefaultHTMLOptions()
	opts.Theme = &export.ThemeColors{
		TitleColor:  "#222",
		AccentColor: "#0078D4",
	}
	htmlStr := export.HTMLWithOptions("Export Demo", slides, opts)
	if err := os.WriteFile(htmlPath, []byte(htmlStr), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Wrote %s\n", htmlPath)

	fmt.Println("Generating PDF export (native gopdf engine)...")
	if err := export.PDF("Export Demo", slides, pdfPath); err != nil {
		log.Fatalf("PDF export failed: %v", err)
	}
	fmt.Printf("Wrote %s\n", pdfPath)
}
