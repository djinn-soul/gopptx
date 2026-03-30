// examples/65-shapes-api demonstrates adding shapes to slides.
//
// Shows rectangles, ellipses, text boxes, arrows, callouts, flowchart shapes,
// stars, and shape properties: fill, line, gradient fill, rotation, text, and
// the text frame (anchor, wrap, autofit).
//
// Run with: go run ./examples/65-shapes-api/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	outputDir  = "examples/output"
	outputFile = "65_shapes_api.pptx"
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

	slides := []pptx.SlideContent{
		buildBasicShapesSlide(),
		buildGradientRotationSlide(),
		buildArrowsSlide(),
		buildFlowchartSlide(),
		buildStarsCalloutsSlide(),
		buildTextFrameSlide(),
	}

	outputPath := filepath.Join(outputDir, outputFile)
	data, err := pptx.CreateWithSlides("Shapes API Demo", slides)
	if err != nil {
		return fmt.Errorf("create presentation: %w", err)
	}
	if err := os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}

func buildBasicShapesSlide() pptx.SlideContent {
	rect := pptx.NewRectangle(0.5, 1.5, 2, 1.2).
		WithFill(pptx.NewShapeFill("4472C4")).
		WithLine(pptx.NewShapeLine("1F497D", styling.Points(2))).
		WithText("Rectangle").
		WithAltText("A blue filled rectangle")

	ellipse := pptx.NewEllipse(3.0, 1.5, 2, 1.2).
		WithFill(pptx.NewShapeFill("C0504D")).
		WithText("Ellipse")

	circle := pptx.NewCircle(5.5, 1.5, 1.2).
		WithFill(pptx.NewShapeFill("9BBB59")).
		WithText("Circle")

	textBox := pptx.NewTextBox("This is a text box – no fill by default.", 0.5, 3.0, 5, 0.8)

	roundedRect := pptx.NewRoundedRectangle(6.0, 1.5, 2.5, 1.2).
		WithFill(pptx.NewShapeFill("F79646")).
		WithText("Rounded")

	return pptx.NewSlide("Basic Shapes").
		AddShape(rect).
		AddShape(ellipse).
		AddShape(circle).
		AddShape(textBox).
		AddShape(roundedRect)
}

func buildGradientRotationSlide() pptx.SlideContent {
	stops := []pptx.ShapeGradientStop{
		pptx.NewShapeGradientStop(0, "4472C4"),
		pptx.NewShapeGradientStop(100, "FFFFFF"),
	}
	gradFill := pptx.NewShapeGradientFill("linear", stops).WithLinearAngle(90)

	gradShape := pptx.NewRectangle(1, 1.5, 3, 2).
		WithGradientFill(gradFill).
		WithText("Gradient Fill").
		WithName("gradient-rect")

	rotatedDiamond := pptx.NewDiamond(5, 1.5, 2, 2).
		WithFill(pptx.NewShapeFill("E91E63")).
		WithText("Rotated 15°").
		WithRotation(15)

	return pptx.NewSlide("Gradient Fill & Rotation").
		AddShape(gradShape).
		AddShape(rotatedDiamond)
}

func buildArrowsSlide() pptx.SlideContent {
	rightArrow := pptx.NewRightArrow(0.5, 1.5, 2, 1).
		WithFill(pptx.NewShapeFill("4472C4")).
		WithText("Right")

	leftArrow := pptx.NewLeftArrow(3.0, 1.5, 2, 1).
		WithFill(pptx.NewShapeFill("C0504D")).
		WithText("Left")

	upArrow := pptx.NewUpArrow(5.5, 1.5, 1, 2).
		WithFill(pptx.NewShapeFill("9BBB59"))

	downArrow := pptx.NewDownArrow(7.0, 1.5, 1, 2).
		WithFill(pptx.NewShapeFill("F79646"))

	chevron := pptx.NewChevron(0.5, 3.5, 2.5, 1).
		WithFill(pptx.NewShapeFill("8064A2")).
		WithText("Chevron")

	return pptx.NewSlide("Arrow Shapes").
		AddShape(rightArrow).
		AddShape(leftArrow).
		AddShape(upArrow).
		AddShape(downArrow).
		AddShape(chevron)
}

func buildFlowchartSlide() pptx.SlideContent {
	process := pptx.NewFlowChartProcess(0.5, 1.5, 2.5, 1).
		WithFill(pptx.NewShapeFill("DCE6F1")).
		WithText("Process")

	decision := pptx.NewFlowChartDecision(3.5, 1.3, 2.5, 1.4).
		WithFill(pptx.NewShapeFill("EBF1DE")).
		WithText("Decision?")

	terminator := pptx.NewFlowChartTerminator(6.5, 1.5, 2.5, 1).
		WithFill(pptx.NewShapeFill("FDE9D9")).
		WithText("End")

	doc := pptx.NewFlowChartDocument(0.5, 3.5, 2.5, 1.2).
		WithFill(pptx.NewShapeFill("E2EFDA")).
		WithText("Document")

	return pptx.NewSlide("Flowchart Shapes").
		AddShape(process).
		AddShape(decision).
		AddShape(terminator).
		AddShape(doc)
}

func buildStarsCalloutsSlide() pptx.SlideContent {
	star := pptx.NewStar(1.5, 2, 1.5).
		WithFill(pptx.NewShapeFill("FFD700")).
		WithText("Star")

	star4 := pptx.NewStar4(3.5, 2, 1.2).
		WithFill(pptx.NewShapeFill("FF6347")).
		WithText("4pt")

	callout := pptx.NewWedgeRectCallout(5.5, 1.5, 2.5, 1.5).
		WithFill(pptx.NewShapeFill("F0F0F0")).
		WithLine(pptx.NewShapeLine("555555", styling.Points(1))).
		WithText("Callout text!")

	heart := pptx.NewHeart(0.5, 3.5, 1.5).
		WithFill(pptx.NewShapeFill("FF69B4"))

	return pptx.NewSlide("Stars, Callouts & More").
		AddShape(star).
		AddShape(star4).
		AddShape(callout).
		AddShape(heart)
}

func buildTextFrameSlide() pptx.SlideContent {
	tfShape := pptx.NewRectangle(1, 1.5, 6, 3).
		WithFill(pptx.NewShapeFill("E8F4FD")).
		WithText("Text with custom text frame:\nTop anchor, no wrap").
		WithVerticalAnchor(pptx.TextAnchorTop).
		WithTextWrap(pptx.TextWrapNone).
		WithAutoFit(pptx.TextAutoFitNormal).
		WithTextMargins(
			styling.Points(10),
			styling.Points(10),
			styling.Points(10),
			styling.Points(10),
		)

	return pptx.NewSlide("Text Frame Properties").
		AddShape(tfShape)
}
