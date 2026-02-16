package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

const (
	outputDir   = "examples/output"
	baseFile    = "19_editor_base.pptx"
	finalFile   = "19_editor_modified.pptx"
	basicSample = "examples/assets/01/01_basic_pptx.pptx"
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

	basePath := filepath.Join(outputDir, baseFile)
	finalPath := filepath.Join(outputDir, finalFile)

	// 1. Create a base presentation
	baseSlides := []pptx.SlideContent{
		pptx.NewSlide("Base Slide 1").AddBullet("Original Content"),
		pptx.NewSlide("Base Slide 2").AddBullet("To be removed"),
		pptx.NewSlide("Base Slide 3").AddBullet("To be updated"),
	}
	if err := pptx.WriteFile(basePath, "Editor Base Demo", baseSlides); err != nil {
		return fmt.Errorf("create base file: %w", err)
	}
	log.Printf("1. Created base: %s\n", basePath)

	// 2. Open with Editor
	editor, err := pptx.OpenPresentationEditor(basePath)
	if err != nil {
		return fmt.Errorf("open editor: %w", err)
	}
	log.Println("2. Opened presentation with Editor")

	// 3. Update a slide (Slide 3 is at index 2)
	updated := pptx.NewSlide("Updated Slide 3").AddBullet("Content has been changed via UpdateSlide")
	if err := editor.UpdateSlide(2, updated); err != nil {
		return fmt.Errorf("update slide: %w", err)
	}
	log.Println("3. Updated Slide 3")

	// 4. Remove a slide (Slide 2 is at index 1)
	if err := editor.RemoveSlide(1); err != nil {
		return fmt.Errorf("remove slide: %w", err)
	}
	log.Println("4. Removed Slide 2")

	// 5. Add a new slide
	newSlide := pptx.NewSlide("Newly Added Slide").AddBullet("Added via AddSlide")
	if _, err := editor.AddSlide(newSlide); err != nil {
		return fmt.Errorf("add slide: %w", err)
	}
	log.Println("5. Added a new slide")

	// 6. Merge from another file (reusing a simple sample if it exists)
	if _, err := os.Stat(basicSample); err == nil {
		if err := editor.MergeFromFile(basicSample); err != nil {
			log.Printf("Warning: merge failed (likely due to asset constraints): %v\n", err)
		} else {
			log.Println("6. Merged slides from 01_basic_pptx.pptx")
		}
	} else {
		log.Println("6. Skipping merge (01_basic_pptx.pptx not found)")
	}

	// 7. Save final result
	if err := editor.Save(finalPath); err != nil {
		return fmt.Errorf("save modified: %w", err)
	}
	log.Printf("7. Saved final modified presentation: %s\n", finalPath)

	return nil
}
