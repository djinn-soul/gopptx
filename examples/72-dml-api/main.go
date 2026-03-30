// examples/72-dml-api demonstrates DrawingML fill, line, and color formatting.
//
// Shows solid fill, gradient fill (linear, radial), no-fill, pattern fill,
// rich line styling (color, width, dash, cap, join), outer/inner/perspective
// shadows, and shape effects.
//
// Run with: go run ./examples/72-dml-api/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	outputDir  = "examples/output"
	outputFile = "72_dml_api.pptx"
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
		buildSolidFillSlide(),
		buildGradientFillSlide(),
		buildPatternFillSlide(),
		buildLineStylesSlide(),
		buildShadowEffectsSlide(),
		buildLineDashConstantsSlide(),
	}

	outputPath := filepath.Join(outputDir, outputFile)
	data, err := pptx.CreateWithSlides("DML API Demo", slides)
	if err != nil {
		return fmt.Errorf("create presentation: %w", err)
	}
	if err := os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}

func buildSolidFillSlide() pptx.SlideContent {
	solidOpaque := pptx.NewRectangle(0.5, 1.5, 2.5, 1.5).
		WithFill(pptx.NewShapeFill("4472C4")).
		WithText("Solid Fill")

	solidTransp := pptx.NewRectangle(3.5, 1.5, 2.5, 1.5).
		WithFill(pptx.NewShapeFill("C0504D").WithTransparency(0.5)).
		WithText("50% Transparent")

	noFill := pptx.NewRectangle(6.5, 1.5, 2.5, 1.5).
		WithRichFill(pptx.NewNoFill()).
		WithLine(pptx.NewShapeLine("4472C4", styling.Points(2))).
		WithText("No Fill")

	return pptx.NewSlide("Solid Fill & Transparency").
		AddShape(solidOpaque).
		AddShape(solidTransp).
		AddShape(noFill)
}

func buildGradientFillSlide() pptx.SlideContent {
	linearStops := []pptx.ShapeGradientStop{
		pptx.NewShapeGradientStop(0, "4472C4"),
		pptx.NewShapeGradientStop(50, "A7C7E7"),
		pptx.NewShapeGradientStop(100, "FFFFFF"),
	}
	linearGrad := pptx.NewShapeGradientFill("linear", linearStops).WithLinearAngle(45)

	radialStops := []pptx.ShapeGradientStop{
		pptx.NewShapeGradientStop(0, "FF6F00"),
		pptx.NewShapeGradientStop(100, "FFF9C4"),
	}
	radialGrad := pptx.NewShapeGradientFill("radial", radialStops)

	gradLinear := pptx.NewRectangle(0.5, 1.5, 3.5, 2).
		WithGradientFill(linearGrad).
		WithText("Linear Gradient 45°")

	gradRadial := pptx.NewRectangle(5, 1.5, 3.5, 2).
		WithGradientFill(radialGrad).
		WithText("Radial Gradient")

	return pptx.NewSlide("Gradient Fills").
		AddShape(gradLinear).
		AddShape(gradRadial)
}

func buildPatternFillSlide() pptx.SlideContent {
	patternShape := pptx.NewRectangle(0.5, 1.5, 3, 2).
		WithRichFill(pptx.NewPatternFill(shapes.PatternTypeDiagStripe)).
		WithText("Diagonal Stripe Pattern")

	patternShape2 := pptx.NewRectangle(4.5, 1.5, 3, 2).
		WithRichFill(pptx.NewPatternFill(shapes.PatternTypeHorz)).
		WithText("Horizontal Lines Pattern")

	return pptx.NewSlide("Pattern Fills").
		AddShape(patternShape).
		AddShape(patternShape2)
}

func buildLineStylesSlide() pptx.SlideContent {
	solidLine := pptx.NewRectangle(0.5, 1.5, 2.5, 1.2).
		WithFill(pptx.NewShapeFill("EBF1DE")).
		WithLine(pptx.NewShapeLine("4472C4", styling.Points(3))).
		WithText("Solid 3pt")

	dashLine := pptx.NewRectangle(3.5, 1.5, 2.5, 1.2).
		WithFill(pptx.NewShapeFill("EBF1DE")).
		WithLine(pptx.NewShapeLine("C0504D", styling.Points(2)).WithDash(styling.LineDashDash)).
		WithText("Dashed 2pt")

	dotLine := pptx.NewRectangle(6.5, 1.5, 2.5, 1.2).
		WithFill(pptx.NewShapeFill("EBF1DE")).
		WithLine(pptx.NewShapeLine("9BBB59", styling.Points(2)).WithDash(styling.LineDashDot)).
		WithText("Dotted 2pt")

	richLine := pptx.NewRichShapeLine("8064A2", styling.Points(4))
	richLineShape := pptx.NewRectangle(0.5, 3.5, 2.5, 1.2).
		WithFill(pptx.NewShapeFill("F0F0F0")).
		WithRichLine(richLine).
		WithText("Rich Line 4pt")

	return pptx.NewSlide("Line Styles").
		AddShape(solidLine).
		AddShape(dashLine).
		AddShape(dotLine).
		AddShape(richLineShape)
}

func buildShadowEffectsSlide() pptx.SlideContent {
	outerShadow := pptx.NewOuterShadow("333333")
	outerShadowShape := pptx.NewRectangle(1, 1.5, 3, 1.5).
		WithFill(pptx.NewShapeFill("4472C4")).
		WithText("Outer Shadow").
		WithRichShadow(outerShadow)

	innerShadow := pptx.NewInnerShadow("000000")
	innerShadowShape := pptx.NewRectangle(5, 1.5, 3, 1.5).
		WithFill(pptx.NewShapeFill("C0504D")).
		WithText("Inner Shadow").
		WithRichShadow(innerShadow)

	return pptx.NewSlide("Shadow Effects").
		AddShape(outerShadowShape).
		AddShape(innerShadowShape)
}

func buildLineDashConstantsSlide() pptx.SlideContent {
	return pptx.NewSlide("Line Dash Style Constants").
		AddBullet(fmt.Sprintf("LineDashSolid      = %q", styling.LineDashSolid)).
		AddBullet(fmt.Sprintf("LineDashDash       = %q", styling.LineDashDash)).
		AddBullet(fmt.Sprintf("LineDashDot        = %q", styling.LineDashDot)).
		AddBullet(fmt.Sprintf("LineDashDashDot    = %q", styling.LineDashDashDot)).
		AddBullet(fmt.Sprintf("LineDashDashDotDot = %q", styling.LineDashDashDotDot)).
		AddBullet(fmt.Sprintf("LineDashLongDash   = %q", styling.LineDashLongDash))
}
