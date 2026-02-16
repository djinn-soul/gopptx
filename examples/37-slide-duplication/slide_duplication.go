package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

const (
	outputDir = "examples/output"
	fileName  = "37_slide_duplication.pptx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	finalPath := filepath.Join(outputDir, fileName)

	// 1. Create a base presentation in memory using CreateWithSlides
	baseSlides := []pptx.SlideContent{
		pptx.NewSlide("Original Slide A").AddBullet("This slide will be duplicated."),
		pptx.NewSlide("Original Slide B").AddBullet("This slide will stay as is."),
		pptx.NewSlide("Original Slide C").AddBullet("This slide will be moved to the beginning."),
	}

	// We need a file on disk for PresentationEditor (currently)
	tempBase := filepath.Join(os.TempDir(), "duplication_base.pptx")
	if err := pptx.WriteFile(tempBase, "Duplication Base", baseSlides); err != nil {
		return fmt.Errorf("create base file: %w", err)
	}
	defer func() {
		if err := os.Remove(tempBase); err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "warning: cleanup temp file %s: %v\n", tempBase, err)
		}
	}()

	// 2. Open with Editor
	editor, openErr := pptx.OpenPresentationEditor(tempBase)
	if openErr != nil {
		return fmt.Errorf("open editor: %w", openErr)
	}

	// 3. Duplicate Slide A (index 0) and place it after Slide B (becomes index 2)
	if _, err := editor.DuplicateSlide(0, 2); err != nil {
		return fmt.Errorf("duplicate slide 0 to 2: %w", err)
	}
	// Current: [A, B, A (Copy), C]

	// 4. Move Slide C (index 3) to the beginning (index 0)
	if err := editor.MoveSlide(3, 0); err != nil {
		return fmt.Errorf("move slide 3 to 0: %w", err)
	}
	// Current: [C, A, B, A (Copy)]

	// 5. Use the ergonomic helper to duplicate Slide B (now index 2)
	if _, err := editor.DuplicateSlideAfter(2); err != nil {
		return fmt.Errorf("duplicate slide after 2: %w", err)
	}
	// Final: [C, A, B, B (Copy), A (Copy)]

	// 6. Save final result
	if err := editor.Save(finalPath); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	log.Printf("Generated slide duplication smoke sample: %s\n", finalPath)

	return nil
}
