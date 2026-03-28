// examples/02-slide-layouts/main.go demonstrates all available slide layout types.
//
// Shows how to select title-only, blank, centered-title, two-column, and
// the default title-and-content layouts using the fluent SlideContent API.
//
// Run with: go run ./examples/02-slide-layouts/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

const (
	outputDir  = "examples/output"
	outputFile = "02_slide_layouts.pptx"
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

	// Slide 1: Default title-and-content layout
	slideDefault := pptx.NewSlide("Title and Content (Default)").
		AddBullet("This is the default layout.").
		AddBullet("Used when no layout is specified.")

	// Slide 2: Title-only layout (no content placeholder)
	slideTitleOnly := pptx.NewSlide("Title Only Layout").
		WithTitleOnlyLayout()

	// Slide 3: Blank layout (no placeholders at all)
	slideBlank := pptx.NewSlide("").
		WithBlankLayout()

	// Slide 4: Centered title layout
	slideCentered := pptx.NewSlide("Centered Title Layout").
		WithCenteredTitleLayout()

	// Slide 5: Two-column layout with bullet content
	slideTwoColumn := pptx.NewSlide("Two Column Layout").
		WithTwoColumnLayout().
		AddBullet("First item in the content area.").
		AddBullet("Second item in the content area.")

	slides := []pptx.SlideContent{
		slideDefault,
		slideTitleOnly,
		slideBlank,
		slideCentered,
		slideTwoColumn,
	}

	data, err := pptx.CreateWithSlides("Slide Layout Types", slides)
	if err != nil {
		return fmt.Errorf("create presentation: %w", err)
	}

	outputPath := filepath.Join(outputDir, outputFile)
	if err = os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}
