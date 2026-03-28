// examples/48-accessibility-alt-text/main.go demonstrates how to set alt-text
// and decorative flags on shapes and images for screen-reader accessibility.
//
// Run with: go run ./examples/48-accessibility-alt-text/main.go
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
	outputFile = "48_accessibility_alt_text.pptx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Minimal 1x1 white PNG used as a placeholder image.
	whitePNG := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41,
		0x54, 0x08, 0xD7, 0x63, 0xF8, 0xFF, 0xFF, 0x3F,
		0x00, 0x05, 0xFE, 0x02, 0xFE, 0xDC, 0x44, 0x74,
		0x06, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E,
		0x44, 0xAE, 0x42, 0x60, 0x82,
	}

	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	// Shape with descriptive alt text and decorative=false (announced by screen readers).
	shape := pptx.NewShape("rect",
		styling.Inches(1), styling.Inches(1),
		styling.Inches(3), styling.Inches(2),
	).WithFill(pptx.NewShapeFill("4472C4")).
		WithText("Accessible Shape").
		WithAltText("A blue rectangle labeled Accessible Shape").
		WithDecorative(false)

	// Image with alt text (content image — not decorative).
	img := pptx.NewImageFromBytes(whitePNG, "png",
		styling.Inches(5), styling.Inches(1),
		styling.Inches(2), styling.Inches(2),
	).WithAltText("A white square placeholder image")

	// Decorative shape: purely visual, screen readers skip it entirely.
	decorative := pptx.NewShape("ellipse",
		styling.Inches(1), styling.Inches(4),
		styling.Inches(1), styling.Inches(1),
	).WithFill(pptx.NewShapeFill("FF0000")).
		WithDecorative(true)

	// Slide 1: mix of accessible and decorative elements.
	slide1 := pptx.NewSlide("Accessibility Demo").
		AddShape(shape).
		AddImage(img).
		AddShape(decorative)

	// Slide 2: explain the accessibility model in bullets.
	slide2 := pptx.NewSlide("Accessibility Guidelines").
		AddBullet("Set WithAltText() on every meaningful shape or image").
		AddBullet("Use WithDecorative(true) for purely visual/ornamental elements").
		AddBullet("Screen readers announce alt text; decorative elements are skipped").
		AddBullet("Alt text should be concise and describe purpose, not appearance")

	data, err := pptx.CreateWithSlides("Task 48: Accessibility Alt Text", []pptx.SlideContent{slide1, slide2})
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}

	outputPath := filepath.Join(outputDir, outputFile)
	if err = os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}
