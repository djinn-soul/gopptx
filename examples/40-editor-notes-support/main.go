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
	outputFile = "40_editor_notes_support.pptx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	// Step 1: build a base presentation with speaker notes already embedded.
	tmpDir, err := os.MkdirTemp("", "gopptx-notes-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	tmpPath := filepath.Join(tmpDir, "base.pptx")

	builder := pptx.NewPresentationBuilder("Notes Demo")
	builder.AddSlide(
		pptx.NewSlide("Slide 1").
			AddBullet("Point A").
			AddBullet("Point B").
			WithNotes("Speaker notes for slide 1.\nThese will be overwritten by the editor."),
	)
	builder.AddSlide(
		pptx.NewSlide("Slide 2").
			AddBullet("Point C").
			AddBullet("Point D").
			WithNotes("Speaker notes for slide 2.\nThese remain unchanged."),
	)

	if err := builder.WriteToFile(tmpPath); err != nil {
		return fmt.Errorf("write base: %w", err)
	}

	// Step 2: open the saved file with the presentation editor.
	e, err := pptx.OpenPresentationEditor(tmpPath)
	if err != nil {
		return fmt.Errorf("open editor: %w", err)
	}
	defer func() {
		if closeErr := e.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "warning: close editor: %v\n", closeErr)
		}
	}()

	// Step 3: overwrite speaker notes on slide 0 via SetNotes.
	if err := e.SetNotes(0, "Updated notes for slide 1 – written via editor.SetNotes"); err != nil {
		return fmt.Errorf("set notes slide 0: %w", err)
	}

	// Step 4: add a brand-new slide that includes notes from the start.
	_, err = e.AddSlide(
		pptx.NewSlide("Slide 3 (Added via Editor)").
			AddBullet("Added programmatically").
			AddBullet("Includes notes written at add time").
			WithNotes("Notes for the new slide – set during AddSlide"),
	)
	if err != nil {
		return fmt.Errorf("add slide: %w", err)
	}

	// Step 5: save the edited presentation to the output directory.
	outputPath := filepath.Join(outputDir, outputFile)
	if err := e.Save(outputPath); err != nil {
		return fmt.Errorf("save: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}
