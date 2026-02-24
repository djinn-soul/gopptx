package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	const outputDir = "examples/output"
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}
	outPath := filepath.Join(outputDir, "22_editor_notes_smoke.pptx")

	// 1. Create a simple presentation first
	builder := pptx.NewPresentationBuilder("Notes Template")
	builder.AddSlide(pptx.NewSlide("Slide with Notes").
		AddBullet("Bullet 1").
		AddBullet("Bullet 2").
		WithNotes("This is a speaker note for slide 1.\nIt has multiple lines."))

	tmpPath := filepath.Join(outputDir, "40_editor_notes_template.pptx")
	if err := builder.WriteToFile(tmpPath); err != nil {
		return fmt.Errorf("failed to save template: %w", err)
	}
	defer func() { _ = os.Remove(tmpPath) }()

	// 2. Open with Editor
	editor, openErr := pptx.OpenEditor(tmpPath)
	if openErr != nil {
		return fmt.Errorf("failed to open editor: %w", openErr)
	}

	// 3. Add a new slide with notes
	slide2 := pptx.NewSlide("New Slide with Notes").
		AddBullet("New Bullet 1").
		WithNotes("Secret speaker notes for slide 2.")

	if _, addErr := editor.AddSlide(slide2); addErr != nil {
		return fmt.Errorf("failed to add slide: %w", addErr)
	}

	// 4. Update existing slide notes
	slide1Updated := pptx.NewSlide("Slide with Updated Notes").
		AddBullet("Updated Bullet 1").
		WithNotes("Updated notes content.")

	if updateErr := editor.UpdateSlide(0, slide1Updated); updateErr != nil {
		return fmt.Errorf("failed to update slide: %w", updateErr)
	}

	// 5. Save
	if err := editor.Save(outPath); err != nil {
		return fmt.Errorf("failed to save edited pptx: %w", err)
	}

	log.Printf("Successfully generated %s\n", outPath)
	return nil
}
