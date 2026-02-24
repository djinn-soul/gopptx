package main

import (
	"log"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func main() {
	log.Println("Generating Protected Presentation...")

	builder := pptx.NewPresentationBuilder("Protected Presentation")

	// Add some content
	builder.AddTitleSlide("Confidential Content").
		AddBulletSlide("Features", []string{
			"Password to Modify: 'test'",
			"Marked as Final",
		})

	// Apply Protection
	builder.WithModifyPassword("test").
		WithMarkAsFinal(true).
		WithSignaturesEnabled(true)

	// Save
	outputPath := "examples/output/46_presentation_protection.pptx"
	if err := builder.WriteToFile(outputPath); err != nil {
		log.Fatalf("Failed to write presentation: %v", err)
	}

	log.Printf("Successfully generated protected presentation at: %s", outputPath)
	log.Println("Please verify manually in PowerPoint:")
	log.Println("1. It should prompt for password 'test' to modify.")
	log.Println("2. It should show a 'Marked as Final' banner.")
}
