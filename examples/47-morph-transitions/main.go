package main

import (
	"fmt"
	"log"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func main() {
	fmt.Println("Generating Morph Transition Presentation...")

	builder := pptx.NewPresentationBuilder("Morph Transition Demo")

	// Slide 1: Start position
	rect1 := shapes.NewRectangle(1.0, 1.0, 2.0, 2.0).
		WithFill(shapes.NewShapeFill("#FF0000")).
		WithName("Morphed-Rect")

	slide1 := pptx.NewSlide("Slide 1: Start").
		AddShape(rect1)

	// Slide 2: End position (Morphed)
	rect2 := shapes.NewRectangle(5.0, 3.0, 4.0, 1.0).
		WithFill(shapes.NewShapeFill("#0000FF")).
		WithName("Morphed-Rect")

	slide2 := pptx.NewSlide("Slide 2: End").
		AddShape(rect2).
		WithMorphTransition()

	builder.AddSlide(slide1).AddSlide(slide2)

	// Save
	outputPath := "examples/output/47_morph_transition.pptx"
	if err := builder.WriteToFile(outputPath); err != nil {
		log.Fatalf("Failed to write presentation: %v", err)
	}

	fmt.Printf("Successfully generated morph transition presentation at: %s\n", outputPath)
	fmt.Println("Please verify manually in PowerPoint:")
	fmt.Println("1. Go to slide 2.")
	fmt.Println("2. You should see a Morph transition animation from the red square to the blue rectangle.")
}
