package main

import (
	"fmt"
	"log"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func main() {
	fmt.Println("Generating Protected Presentation...")

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

	fmt.Printf("Successfully generated protected presentation at: %s\n", outputPath)
	fmt.Println("Please verify manually in PowerPoint:")
	fmt.Println("1. It should prompt for password 'test' to modify.")
	fmt.Println("2. It should show a 'Marked as Final' banner.")
}
