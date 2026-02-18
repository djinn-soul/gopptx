//go:build ignore

package main

import (
	"fmt"
	"log"
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
	log.Printf("Original slide count: %d\n", count)

	// 2. Complex Manipulation
	// Duplicate the first slide (usually a title) to the end
	if _, dupErr := editor.DuplicateSlide(0, count); dupErr != nil {
		return fmt.Errorf("duplicate title slide to end: %w", dupErr)
	}

	// Move the now-last slide to index 1
	if moveErr := editor.MoveSlide(editor.SlideCount()-1, 1); moveErr != nil {
		return fmt.Errorf("move cloned title to index 1: %w", moveErr)
	}

	// Duplicate a slide from the middle somewhere
	if count > 2 {
		if _, dupErr := editor.DuplicateSlideAfter(2); dupErr != nil {
			return fmt.Errorf("duplicate middle slide: %w", dupErr)
		}
	}

	// 3. Save final result
	if saveErr := editor.Save(destPath); saveErr != nil {
		return fmt.Errorf("save: %w", saveErr)
	}
	log.Printf("Generated complex duplication smoke sample: %s\n", destPath)

	return nil
}
