package main

import (
	"fmt"
	"log"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func main() {
	outPath := "editor_notes_smoke.pptx"

	// 1. Create a simple presentation first
	builder := pptx.NewPresentationBuilder("Notes Template")
	builder.AddSlide(pptx.NewSlide("Slide with Notes").
		AddBullet("Bullet 1").
		AddBullet("Bullet 2").
		WithNotes("This is a speaker note for slide 1.\nIt has multiple lines."))

	tmpPath := "notes_template.pptx"
	if err := builder.WriteToFile(tmpPath); err != nil {
		log.Fatalf("failed to save template: %v", err)
	}
	defer func() { _ = os.Remove(tmpPath) }()

	// 2. Open with Editor
	editor, err := pptx.OpenEditor(tmpPath)
	if err != nil {
		log.Fatalf("failed to open editor: %v", err)
	}

	// 3. Add a new slide with notes
	slide2 := pptx.NewSlide("New Slide with Notes").
		AddBullet("New Bullet 1").
		WithNotes("Secret speaker notes for slide 2.")

	_, err = editor.AddSlide(slide2)
	if err != nil {
		log.Fatalf("failed to add slide: %v", err)
	}

	// 4. Update existing slide notes
	slide1Updated := pptx.NewSlide("Slide with Updated Notes").
		AddBullet("Updated Bullet 1").
		WithNotes("Updated notes content.")

	err = editor.UpdateSlide(0, slide1Updated)
	if err != nil {
		log.Fatalf("failed to update slide: %v", err)
	}

	// 5. Save
	if err := editor.Save(outPath); err != nil {
		log.Fatalf("failed to save edited pptx: %v", err)
	}

	fmt.Printf("Successfully generated %s\n", outPath)
}
