package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func mainEditorChartSmoke() {
	templatePath := "ppt-rs/chart_example.pptx"
	outputDirChart := "examples/output"
	if err := os.MkdirAll(outputDirChart, 0o755); err != nil {
		log.Fatalf("failed to create output dir: %v", err)
	}
	outPath := filepath.Join(outputDirChart, "39_editor_chart_smoke.pptx")

	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		log.Fatalf("template file %s not found. Please ensure it exists.", templatePath)
	}

	// 1. Open with Editor
	editor, openErr := pptx.OpenPresentationEditor(templatePath)
	if openErr != nil {
		log.Fatalf("failed to open editor: %v", openErr)
	}

	// 2. Duplicate the slide with the chart (assuming it's the first slide)
	// Duplicate slide 0 to index 1
	newIndex, err := editor.DuplicateSlide(0, 1)
	if err != nil {
		log.Fatalf("failed to duplicate slide: %v", err)
	}

	log.Printf("Duplicated chart slide to index %d\n", newIndex)

	// 3. Save
	if err := editor.Save(outPath); err != nil {
		log.Fatalf("failed to save edited pptx: %v", err)
	}

	log.Printf("Successfully generated %s\n", outPath)
}
