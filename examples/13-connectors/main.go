// examples/13-connectors/main.go demonstrates the three connector types in gopptx.
//
// Shows straight, elbow, and curved connectors placed on separate slides
// alongside the shapes they connect.
//
// Run with: go run ./examples/13-connectors/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	log "github.com/djinn-soul/gopptx/pkg/stdlog"
)

const (
	outputDir  = "examples/output"
	outputFile = "13_connectors.pptx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	// --- Slide 1: overview ---
	slide1 := pptx.NewSlide("Connectors Demo").
		AddBullet("gopptx supports three connector types:").
		AddBullet("  • Straight — direct line between two points").
		AddBullet("  • Elbow   — right-angle (bent) connector").
		AddBullet("  • Curved  — smooth bezier-curve connector")

	// --- Slide 2: straight connector ---
	straight := pptx.NewStraightConnector(
		styling.Inches(1), styling.Inches(2), // startX, startY
		styling.Inches(5), styling.Inches(2), // endX, endY
	)

	leftBox := pptx.NewShape(pptx.ShapeTypeRectangle,
		styling.Inches(0.25), styling.Inches(1.75),
		styling.Inches(0.75), styling.Inches(0.5),
	).WithFill(pptx.NewShapeFill("4472C4")).WithText("A")

	rightBox := pptx.NewShape(pptx.ShapeTypeRectangle,
		styling.Inches(5), styling.Inches(1.75),
		styling.Inches(0.75), styling.Inches(0.5),
	).WithFill(pptx.NewShapeFill("4472C4")).WithText("B")

	slide2 := pptx.NewSlide("Straight Connector").
		AddShape(leftBox).
		AddShape(rightBox).
		AddConnector(straight).
		AddBullet("pptx.NewStraightConnector(startX, startY, endX, endY)")

	// --- Slide 3: elbow connector ---
	elbow := pptx.NewElbowConnector(
		styling.Inches(1), styling.Inches(2),
		styling.Inches(5), styling.Inches(4),
	)

	topBox := pptx.NewShape(pptx.ShapeTypeRectangle,
		styling.Inches(0.25), styling.Inches(1.75),
		styling.Inches(0.75), styling.Inches(0.5),
	).WithFill(pptx.NewShapeFill("ED7D31")).WithText("C")

	bottomBox := pptx.NewShape(pptx.ShapeTypeRectangle,
		styling.Inches(5), styling.Inches(3.75),
		styling.Inches(0.75), styling.Inches(0.5),
	).WithFill(pptx.NewShapeFill("ED7D31")).WithText("D")

	slide3 := pptx.NewSlide("Elbow Connector").
		AddShape(topBox).
		AddShape(bottomBox).
		AddConnector(elbow).
		AddBullet("pptx.NewElbowConnector(startX, startY, endX, endY)")

	// --- Slide 4: curved connector ---
	curved := pptx.NewCurvedConnector(
		styling.Inches(1), styling.Inches(3),
		styling.Inches(5), styling.Inches(5),
	)

	topLeft := pptx.NewShape(pptx.ShapeTypeEllipse,
		styling.Inches(0.25), styling.Inches(2.75),
		styling.Inches(0.75), styling.Inches(0.5),
	).WithFill(pptx.NewShapeFill("A9D18E")).WithText("E")

	bottomRight := pptx.NewShape(pptx.ShapeTypeEllipse,
		styling.Inches(5), styling.Inches(4.75),
		styling.Inches(0.75), styling.Inches(0.5),
	).WithFill(pptx.NewShapeFill("A9D18E")).WithText("F")

	slide4 := pptx.NewSlide("Curved Connector").
		AddShape(topLeft).
		AddShape(bottomRight).
		AddConnector(curved).
		AddBullet("pptx.NewCurvedConnector(startX, startY, endX, endY)")

	slides := []pptx.SlideContent{slide1, slide2, slide3, slide4}

	outputPath := filepath.Join(outputDir, outputFile)
	if err := pptx.WriteFile(outputPath, "Connectors Demo", slides); err != nil {
		return fmt.Errorf("write presentation: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}
