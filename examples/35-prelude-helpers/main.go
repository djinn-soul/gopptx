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
	outputFile = "35_prelude_helpers.pptx"
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

	builder := pptx.NewPresentationBuilder("Prelude Helpers Demo")

	// Set widescreen 16:9 slide size via the custom SlideSize helper.
	// WithSlideSize accepts a SlideSize value; SlideSize16x9() returns the
	// standard widescreen preset (approximately 16in × 9in in EMU).
	builder.WithSlideSize(pptx.SlideSize16x9())

	// Apply a theme using the fluent builder.
	builder.WithTheme(styling.ThemeModern)

	// Slide 1: welcome / overview
	builder.AddSlide(
		pptx.NewSlide("Welcome").
			AddBullet("Built with PresentationBuilder").
			AddBullet("Fluent API for creating presentations").
			AddBullet("Chain methods for concise setup"),
	)

	// Slide 2: numbered list of features
	builder.AddSlide(
		pptx.NewSlide("Features").
			AddNumbered("Builder pattern").
			AddNumbered("Theme support").
			AddNumbered("Slide size control").
			AddNumbered("Metadata helpers"),
	)

	// Slide 3: slide size awareness – show the EMU values of a 16:9 slide
	w16 := styling.Inches(16)
	h9 := styling.Inches(9)
	builder.AddSlide(
		pptx.NewSlide("Slide Size Awareness").
			AddBullet(fmt.Sprintf("16 inches wide = %d EMU", w16.Emu())).
			AddBullet(fmt.Sprintf(" 9 inches tall = %d EMU", h9.Emu())).
			AddBullet("styling.Inches() converts from human-readable to EMU").
			AddBullet("SlideSize16x9() returns the standard preset"),
	)

	// Slide 4: metadata and slide numbers
	builder.WithSlideNumbers(true)
	builder.AddSlide(
		pptx.NewSlide("Builder Options").
			AddBullet("WithSlideNumbers(true) – page numbers on every slide").
			AddBullet("WithFooter(text) – footer across all slides").
			AddBullet("WithTheme(theme) – apply a color+font palette").
			AddBullet("WriteToFile(path) – build and persist in one call"),
	)

	outputPath := filepath.Join(outputDir, outputFile)
	if err := builder.WriteToFile(outputPath); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}
