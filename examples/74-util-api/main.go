// examples/74-util-api demonstrates unit conversion helpers.
//
// Shows the styling.Length type and the Inches, Centimeters, Points, and Emu
// constructor functions, plus the reverse conversion methods (.Inches(), .Cm(),
// .Pt(), .Emu()) and FontSize().
//
// Run with: go run ./examples/74-util-api/main.go
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
	outputFile = "74_util_api.pptx"
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

	// --- Demonstrate unit conversions ---

	// Inches <-> EMU
	oneInch := styling.Inches(1)
	log.Printf("1 inch = %d EMU\n", oneInch.Emu())
	log.Printf("914400 EMU = %.4f inches\n", styling.Emu(914400).Inches())

	// Centimeters <-> EMU
	oneCm := styling.Centimeters(1)
	log.Printf("1 cm = %d EMU\n", oneCm.Emu())
	log.Printf("360000 EMU = %.4f cm\n", styling.Emu(360000).Cm())

	// Points <-> EMU
	onePt := styling.Points(1)
	log.Printf("1 pt = %d EMU\n", onePt.Emu())
	log.Printf("12700 EMU = %.4f pt\n", styling.Emu(12700).Pt())

	// Direct EMU
	directEmu := styling.Emu(914400)
	log.Printf("Emu(914400) inches = %.4f\n", directEmu.Inches())

	// FontSize (pt -> hundredths-of-a-point for OOXML)
	log.Printf("FontSize(18pt) = %d (hundredths of a point)\n", styling.FontSize(18))
	log.Printf("FontSize(44pt) = %d\n", styling.FontSize(44))

	// MaxEMU boundary
	log.Printf("MaxEMU = %d\n", styling.MaxEMU.Emu())

	// --- Build a slide using every unit type ---
	rectInches := pptx.NewRectangle(
		styling.Inches(0.5).Inches(),
		styling.Inches(1.5).Inches(),
		styling.Inches(2.5).Inches(),
		styling.Inches(1).Inches(),
	).WithFill(pptx.NewShapeFill("4472C4")).WithText("Inches")

	rectCm := pptx.NewShape(
		pptx.ShapeTypeRectangle,
		styling.Centimeters(8),
		styling.Centimeters(4),
		styling.Centimeters(5),
		styling.Centimeters(2),
	).WithFill(pptx.NewShapeFill("C0504D")).WithText("Centimeters")

	rectPt := pptx.NewShape(
		pptx.ShapeTypeRectangle,
		styling.Points(50),
		styling.Points(400),
		styling.Points(200),
		styling.Points(80),
	).WithFill(pptx.NewShapeFill("9BBB59")).WithText("Points")

	rectEmu := pptx.NewShape(
		pptx.ShapeTypeRectangle,
		styling.Emu(5000000),
		styling.Emu(3000000),
		styling.Emu(2500000),
		styling.Emu(800000),
	).WithFill(pptx.NewShapeFill("F79646")).WithText("EMU")

	infoSlide := pptx.NewSlide("Unit Conversion Results").
		AddBullet(fmt.Sprintf("1 inch = %d EMU", oneInch.Emu())).
		AddBullet(fmt.Sprintf("1 cm   = %d EMU", oneCm.Emu())).
		AddBullet(fmt.Sprintf("1 pt   = %d EMU", onePt.Emu())).
		AddBullet(fmt.Sprintf("FontSize(18pt) = %d", styling.FontSize(18))).
		AddBullet(fmt.Sprintf("MaxEMU = %d", styling.MaxEMU.Emu()))

	shapesSlide := pptx.NewSlide("Shapes Placed with Different Units").
		AddShape(rectInches).
		AddShape(rectCm).
		AddShape(rectPt).
		AddShape(rectEmu)

	slides := []pptx.SlideContent{infoSlide, shapesSlide}

	outputPath := filepath.Join(outputDir, outputFile)
	data, err := pptx.CreateWithSlides("Util API Demo", slides)
	if err != nil {
		return fmt.Errorf("create presentation: %w", err)
	}
	if err := os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}
