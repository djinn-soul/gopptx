package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Smoke test failed: %v", err)
	}
	log.Println("Smoke test completed successfully!")
}

func run() error {
	const outputDir = "examples/output"
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	// 1. Create a base presentation with a known shape
	deck := pptx.NewPresentationBuilder("Smoke Test")

	// Add a TextBox (Inches)
	// 1 Inch = 914400 EMUs
	// Shape at 1,1 inches. Size 5x2 inches.
	tb := pptx.NewTextBox("Original Text", 1.0, 1.0, 5.0, 2.0)
	deck.AddShapesSlide("Title Slide", tb)

	inputFile := filepath.Join(outputDir, "41_shape_original.pptx")
	if err := deck.WriteToFile(inputFile); err != nil {
		return fmt.Errorf("failed to save input file: %w", err)
	}
	log.Printf("Created base file: %s\n", inputFile)

	// 2. Open with Editor
	edit, err := editor.OpenPresentationEditor(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open editor: %w", err)
	}
	defer func() {
		if closeErr := edit.Close(); closeErr != nil {
			log.Printf("warning: failed to close editor: %v", closeErr)
		}
	}()
	shapes, err := edit.GetShapes(0)
	if err != nil {
		return fmt.Errorf("GetShapes failed: %w", err)
	}

	targetIndex := -1
	for i, s := range shapes {
		log.Printf("Shape %d: ID=%d Name=%q Text=%q\n", i, s.ID, s.Name, s.Text)
		if s.Text == "Original Text" {
			targetIndex = i
		}
	}

	if targetIndex == -1 {
		return errors.New("could not find shape with text 'Original Text'")
	}

	// 4. Update the shape
	// Note: editor works with raw EMUs.
	// We want to verify update persistence.
	// Let's set it to exactly 500,000 EMUs (approx 0.5 inches) just to have a clean integer number to verify.
	log.Println("Updating shape...")
	updatedX := 500000
	updatedY := 500000

	newText := "Edited Text"
	// Create vars for pointers
	newXVal := updatedX
	newYVal := updatedY

	err = edit.UpdateShape(0, shapes[targetIndex].ID, common.ShapeUpdate{
		Text: &newText,
		X:    &newXVal,
		Y:    &newYVal,
	})
	if err != nil {
		return fmt.Errorf("UpdateShape failed: %w", err)
	}

	// 5. Save
	outputFile := filepath.Join(outputDir, "41_shape_edited.pptx")
	if err := edit.Save(outputFile); err != nil {
		return fmt.Errorf("failed to save output file: %w", err)
	}
	log.Printf("Saved edited file: %s\n", outputFile)

	// 6. Verify by re-opening
	verifyEdit, err := editor.OpenPresentationEditor(outputFile)
	if err != nil {
		return fmt.Errorf("failed to open verification file: %w", err)
	}
	defer func() {
		if closeErr := verifyEdit.Close(); closeErr != nil {
			log.Printf("warning: failed to close verification editor: %v", closeErr)
		}
	}()

	vShapes, err := verifyEdit.GetShapes(0)
	if err != nil {
		return fmt.Errorf("verification GetShapes failed: %w", err)
	}

	foundEdited := false
	for _, s := range vShapes {
		log.Printf("Verification Shape: ID=%d Text=%q X=%d Y=%d\n", s.ID, s.Text, s.X, s.Y)
		if s.Text == "Edited Text" && s.X == updatedX && s.Y == updatedY {
			foundEdited = true
			break
		}
	}

	if !foundEdited {
		return fmt.Errorf("verification failed: could not find shape with 'Edited Text' at %d,%d", updatedX, updatedY)
	}

	return nil
}
