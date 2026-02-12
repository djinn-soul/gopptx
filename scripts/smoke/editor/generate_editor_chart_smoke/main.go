package main

import (
	"fmt"
	"log"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func main() {
	templatePath := "ppt-rs/chart_example.pptx"
	outPath := "editor_chart_smoke.pptx"

	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		log.Fatalf("template file %s not found. Please ensure it exists.", templatePath)
	}

	// 1. Open with Editor
	editor, err := pptx.OpenEditor(templatePath)
	if err != nil {
		log.Fatalf("failed to open editor: %v", err)
	}

	// 2. Duplicate the slide with the chart (assuming it's the first slide)
	// Duplicate slide 0 to index 1
	newIndex, err := editor.DuplicateSlide(0, 1)
	if err != nil {
		log.Fatalf("failed to duplicate slide: %v", err)
	}

	fmt.Printf("Duplicated chart slide to index %d\n", newIndex)

	// 3. Save
	if err := editor.Save(outPath); err != nil {
		log.Fatalf("failed to save edited pptx: %v", err)
	}

	fmt.Printf("Successfully generated %s\n", outPath)
}
