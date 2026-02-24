package main

import (
	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func mainSectionSmoke() {
	templatePath := "examples/assets/37/162301-moneybox-template-16x9.pptx"
	outputPath := "examples/output/39_modular_sections.pptx"

	log.Printf("Opening %s...\n", templatePath)
	e, openErr := pptx.OpenPresentationEditor(templatePath)
	if openErr != nil {
		log.Fatalf("Failed to open template: %v", openErr)
	}

	log.Println("Defining sections...")
	// Section 1: Intro (Slide 1)
	if err := e.AddSection("Intro", []int{0}); err != nil {
		log.Fatalf("AddSection Intro failed: %v", err)
	}

	// Section 2: Core Analysis (Slides 2, 3, 4)
	if err := e.AddSection("Core Analysis", []int{1, 2, 3}); err != nil {
		log.Fatalf("AddSection Analysis failed: %v", err)
	}

	// Section 3: Summary (Remaining slides)
	remaining := make([]int, 0)
	for i := 4; i < e.SlideCount(); i++ {
		remaining = append(remaining, i)
	}
	if err := e.AddSection("Summary", remaining); err != nil {
		log.Fatalf("AddSection Summary failed: %v", err)
	}

	log.Printf("Saving to %s...\n", outputPath)
	if err := pptx.Save(e, outputPath); err != nil {
		log.Fatalf("Failed to save: %v", err)
	}

	log.Println("Done!")
}

var _ = mainSectionSmoke
