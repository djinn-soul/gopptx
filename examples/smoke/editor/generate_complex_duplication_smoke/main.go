package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

const (
	inputDir  = "examples/assets/37"
	outputDir = "examples/output"
	inputFile = "160070-labyrinth-template-16x9.pptx"
	outFile   = "37_complex_duplication.pptx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	srcPath := filepath.Join(inputDir, inputFile)
	destPath := filepath.Join(outputDir, outFile)

	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("source template missing: %s", srcPath)
	}

	// 1. Open with Editor
	editor, err := pptx.OpenPresentationEditor(srcPath)
	if err != nil {
		return fmt.Errorf("open editor: %w", err)
	}

	count := editor.SlideCount()
	fmt.Printf("Original slide count: %d\n", count)

	// 2. Complex Manipulation
	// Duplicate the first slide (usually a title) to the end
	if _, err := editor.DuplicateSlide(0, count); err != nil {
		return fmt.Errorf("duplicate title slide to end: %w", err)
	}

	// Move the now-last slide to index 1
	if err := editor.MoveSlide(editor.SlideCount()-1, 1); err != nil {
		return fmt.Errorf("move cloned title to index 1: %w", err)
	}

	// Duplicate a slide from the middle somewhere
	if count > 2 {
		if _, err := editor.DuplicateSlideAfter(2); err != nil {
			return fmt.Errorf("duplicate middle slide: %w", err)
		}
	}

	// 3. Save final result
	if err := editor.Save(destPath); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	fmt.Printf("Generated complex duplication smoke sample: %s\n", destPath)

	return nil
}
