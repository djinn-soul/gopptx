package main

import (
	"fmt"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	log "github.com/djinn-soul/gopptx/pkg/stdlog"
)

const (
	outputDir  = "examples/output"
	outputFile = "35_layout_helpers.pptx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	stackSlide, err := buildStackSlide()
	if err != nil {
		return err
	}
	distSlide, err := buildDistributeSlide()
	if err != nil {
		return err
	}
	gridSlide, err := buildGridSlide()
	if err != nil {
		return err
	}
	centerSlide := buildCenterSlide()

	slides := []pptx.SlideContent{stackSlide, distSlide, gridSlide, centerSlide}

	data, err := pptx.CreateWithSlides("Task 35: Layout Helpers", slides)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}

	outputPath := outputDir + "/" + outputFile
	if err = os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}

func buildStackSlide() (pptx.SlideContent, error) {
	slide := pptx.NewSlide("Stack Layout")
	h := pptx.Inches(0.6)
	w := pptx.Inches(8)
	pts, err := pptx.Stack(
		pptx.OrientationVertical,
		pptx.Point{X: pptx.Inches(1), Y: pptx.Inches(1.2)},
		pptx.Inches(0.12),
		pptx.Size{CX: w, CY: h},
		pptx.Size{CX: w, CY: h},
		pptx.Size{CX: w, CY: h},
		pptx.Size{CX: w, CY: h},
	)
	if err != nil {
		return slide, fmt.Errorf("stack: %w", err)
	}
	colors := []string{"4472C4", "ED7D31", "A9D18E", "FF0000"}
	labels := []string{"Init", "Build", "Validate", "Export"}
	for i, pt := range pts {
		slide = slide.AddShape(
			pptx.NewShape(pptx.ShapeTypeRectangle, pt.X, pt.Y, w, h).
				WithFill(pptx.NewShapeFill(colors[i])).
				WithText(labels[i]),
		)
	}
	return slide, nil
}

func buildDistributeSlide() (pptx.SlideContent, error) {
	slide := pptx.NewSlide("Distribute Uniform")
	bounds := pptx.Box{
		X: pptx.Inches(0.5), Y: pptx.Inches(2),
		CX: pptx.Inches(9), CY: pptx.Inches(1.5),
	}
	elemW := pptx.Inches(1.5)
	coords, err := pptx.DistributeUniform(pptx.OrientationHorizontal, bounds, 5, elemW)
	if err != nil {
		return slide, fmt.Errorf("distribute: %w", err)
	}
	colors := []string{"4472C4", "ED7D31", "A9D18E", "FF0000", "FFC000"}
	labels := []string{"Alpha", "Beta", "Gamma", "Delta", "Epsilon"}
	for i, x := range coords {
		slide = slide.AddShape(
			pptx.NewShape(pptx.ShapeTypeRoundedRectangle, x, bounds.Y, elemW, bounds.CY).
				WithFill(pptx.NewShapeFill(colors[i])).
				WithText(labels[i]).
				WithTextWrap(pptx.TextWrapNone).
				WithAutoFit(pptx.TextAutoFitNone),
		)
	}
	return slide, nil
}

func buildGridSlide() (pptx.SlideContent, error) {
	slide := pptx.NewSlide("Grid Layout (2×3)").WithTitleOnlyLayout()
	boxes, err := pptx.GridInBox(
		2,
		3,
		pptx.Inches(0.2),
		pptx.Box{
			X:  pptx.Inches(0.5),
			Y:  pptx.Inches(1.9),
			CX: pptx.Inches(9),
			CY: pptx.Inches(4.9),
		},
	)
	if err != nil {
		return slide, fmt.Errorf("grid: %w", err)
	}
	colors := []string{"5B9BD5", "ED7D31", "A9D18E", "FF0000", "FFC000", "7030A0"}
	for i, b := range boxes {
		slide = slide.AddShape(
			pptx.NewShape(pptx.ShapeTypeRectangle, b.X, b.Y, b.CX, b.CY).
				WithFill(pptx.NewShapeFill(colors[i])).
				WithText(fmt.Sprintf("Cell %d", i+1)),
		)
	}
	return slide, nil
}

func buildCenterSlide() pptx.SlideContent {
	slide := pptx.NewSlide("Center Helper")
	cx, cy := pptx.Inches(4), pptx.Inches(2)
	x, y := pptx.Center(cx, cy)
	return slide.AddShape(
		pptx.NewShape(pptx.ShapeTypeEllipse, x, y, cx, cy).
			WithFill(pptx.NewShapeFill("1B6CA8")).
			WithText("Centered Ellipse"),
	)
}
