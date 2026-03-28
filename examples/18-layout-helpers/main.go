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
	outputFile = "18_layout_helpers.pptx"
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

	builder := pptx.NewPresentationBuilder("Layout Helpers Demo").
		WithSlideSize(pptx.SlideSize16x9())

	// Slide 1: demonstrate unit conversion values
	convSlide := buildUnitConversionSlide()
	builder.AddSlide(convSlide)

	// Slide 2: shapes positioned using styling.Inches
	posSlide := buildPositionedShapesSlide()
	builder.AddSlide(posSlide)

	// Slide 3: demonstrate Emu and retrieval helpers
	emuSlide := buildEmuHelpersSlide()
	builder.AddSlide(emuSlide)

	outputPath := filepath.Join(outputDir, outputFile)
	if err := builder.WriteToFile(outputPath); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}

// buildUnitConversionSlide shows the unit-conversion helpers and their values.
func buildUnitConversionSlide() pptx.SlideContent {
	// styling.Inches / styling.InchesToEMU / styling.Emu
	oneInch := styling.Inches(1)
	twoInch := styling.InchesToEMU(2) // alias for styling.Inches(2)
	rawEmu := styling.Emu(914400)     // 1 inch in EMU
	_ = twoInch

	slide := pptx.NewSlide("Unit Conversion Helpers").
		AddBullet(fmt.Sprintf("styling.Inches(1)     = %d EMU (%.4f in)", oneInch.Emu(), oneInch.Inches())).
		AddBullet(fmt.Sprintf("styling.InchesToEMU(2) = %d EMU (%.4f in)", twoInch.Emu(), twoInch.Inches())).
		AddBullet(fmt.Sprintf("styling.Emu(914400)   = %d EMU (%.4f in)", rawEmu.Emu(), rawEmu.Inches())).
		AddBullet(fmt.Sprintf("styling.Centimeters(2.54) = %d EMU", styling.Centimeters(2.54).Emu())).
		AddBullet(fmt.Sprintf("styling.Points(72)    = %d EMU", styling.Points(72).Emu()))

	return slide
}

// buildPositionedShapesSlide places two shapes using styling.Inches for precise coordinates.
func buildPositionedShapesSlide() pptx.SlideContent {
	// Shape at (1in, 1in) with size 2in × 1in
	x1 := styling.Inches(1)
	y1 := styling.Inches(1)
	w := styling.Inches(2)
	h := styling.Inches(1)

	// Shape at (4in, 2in) with size 2in × 1in
	x2 := styling.Inches(4)
	y2 := styling.Inches(2)

	slide := pptx.NewSlide("Precise Positioning with styling.Inches").
		AddShape(
			pptx.NewShape(pptx.ShapeTypeRectangle, x1, y1, w, h).
				WithFill(pptx.NewShapeFill("4472C4")).
				WithText(fmt.Sprintf("(%.0fin, %.0fin) 2×1in", x1.Inches(), y1.Inches())),
		).
		AddShape(
			pptx.NewShape(pptx.ShapeTypeRoundedRectangle, x2, y2, w, h).
				WithFill(pptx.NewShapeFill("ED7D31")).
				WithText(fmt.Sprintf("(%.0fin, %.0fin) 2×1in", x2.Inches(), y2.Inches())),
		)

	return slide
}

// buildEmuHelpersSlide demonstrates retrieving values back from a Length.
func buildEmuHelpersSlide() pptx.SlideContent {
	oneInch := styling.Inches(1)
	rawEmu := styling.Emu(914400)

	slide := pptx.NewSlide("Retrieving Values from Length").
		AddBullet(fmt.Sprintf("oneInch.Inches() = %.6f", oneInch.Inches())).
		AddBullet(fmt.Sprintf("oneInch.Emu()    = %d", oneInch.Emu())).
		AddBullet(fmt.Sprintf("oneInch.Cm()     = %.4f", oneInch.Cm())).
		AddBullet(fmt.Sprintf("oneInch.Pt()     = %.4f", oneInch.Pt())).
		AddBullet(fmt.Sprintf("rawEmu.Inches()  = %.6f (should equal 1.0)", rawEmu.Inches())).
		AddBullet(fmt.Sprintf("rawEmu.Emu()     = %d", rawEmu.Emu()))

	return slide
}
