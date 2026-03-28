// examples/05-bullet-styles/main.go demonstrates all available bullet list styles.
//
// Shows standard bullets, numbered lists, lettered lists, and sub-bullets at
// various indent levels using the SlideContent fluent API.
//
// Run with: go run ./examples/05-bullet-styles/main.go
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
	outputFile = "05_bullet_styles.pptx"
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

	// Slide 1: Standard bullet points
	slideStandard := pptx.NewSlide("Standard Bullet Points").
		AddBullet("First standard bullet item.").
		AddBullet("Second standard bullet item.").
		AddBullet("Third standard bullet item.")

	// Slide 2: Numbered bullet list
	slideNumbered := pptx.NewSlide("Numbered List").
		AddNumbered("First numbered item.").
		AddNumbered("Second numbered item.").
		AddNumbered("Third numbered item.")

	// Slide 3: Lettered bullet list
	slideLettered := pptx.NewSlide("Lettered List").
		AddLettered("Item alpha.").
		AddLettered("Item beta.").
		AddLettered("Item gamma.")

	// Slide 4: Sub-bullets at multiple indent levels
	slideSubBullets := pptx.NewSlide("Sub-Bullet Levels").
		AddBullet("Top-level bullet (level 0).").
		AddSubBullet(1, "Sub-bullet at level 1.").
		AddSubBullet(2, "Sub-bullet at level 2.").
		AddBullet("Back to top-level.").
		AddSubBullet(1, "Another level 1 sub-bullet.")

	// Slide 5: Mixed styles on one slide
	slideMixed := pptx.NewSlide("Mixed Bullet Styles").
		AddBullet("Standard bullet item.").
		AddNumbered("Numbered item follows.").
		AddLettered("Lettered item next.").
		AddSubBullet(1, "Indented sub-point.")

	slides := []pptx.SlideContent{
		slideStandard,
		slideNumbered,
		slideLettered,
		slideSubBullets,
		slideMixed,
	}

	data, err := pptx.CreateWithSlides("Bullet Styles Demo", slides)
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
