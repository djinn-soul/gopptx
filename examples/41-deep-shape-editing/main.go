package main

import (
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
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	inputFile := filepath.Join(outputDir, "41_shape_original.pptx")
	if err := createBaseDeck(inputFile); err != nil {
		return err
	}
	log.Printf("Created base file: %s\n", inputFile)

	edit, openErr := editor.OpenPresentationEditor(inputFile)
	if openErr != nil {
		return fmt.Errorf("failed to open editor: %w", openErr)
	}
	defer func() {
		if closeErr := edit.Close(); closeErr != nil {
			log.Printf("warning: failed to close editor: %v", closeErr)
		}
	}()
	shapes, targetIndex, err := getShapesAndTarget(edit, "Original Text")
	if err != nil {
		return err
	}

	updatedX := 500000
	updatedY := 500000
	if err := updateTargetShape(edit, shapes[targetIndex].ID, "Edited Text", updatedX, updatedY); err != nil {
		return err
	}

	outputFile := filepath.Join(outputDir, "41_shape_edited.pptx")
	if err := edit.Save(outputFile); err != nil {
		return fmt.Errorf("failed to save output file: %w", err)
	}
	log.Printf("Saved edited file: %s\n", outputFile)

	return verifyEditedShape(outputFile, "Edited Text", updatedX, updatedY)
}

func createBaseDeck(inputFile string) error {
	deck := pptx.NewPresentationBuilder("Smoke Test")
	tb := pptx.NewTextBox("Original Text", 1.0, 1.0, 5.0, 2.0)
	deck.AddShapesSlide("Title Slide", tb)
	if err := deck.WriteToFile(inputFile); err != nil {
		return fmt.Errorf("failed to save input file: %w", err)
	}
	return nil
}

func getShapesAndTarget(edit *editor.PresentationEditor, text string) ([]common.Shape, int, error) {
	shapes, err := edit.GetShapes(0)
	if err != nil {
		return nil, -1, fmt.Errorf("GetShapes failed: %w", err)
	}
	targetIndex := findShapeByText(shapes, text)
	if targetIndex == -1 {
		return nil, -1, fmt.Errorf("could not find shape with text %q", text)
	}
	return shapes, targetIndex, nil
}

func findShapeByText(shapes []common.Shape, text string) int {
	for i, s := range shapes {
		log.Printf("Shape %d: ID=%d Name=%q Text=%q\n", i, s.ID, s.Name, s.Text)
		if s.Text == text {
			return i
		}
	}
	return -1
}

func updateTargetShape(edit *editor.PresentationEditor, shapeID int, newText string, updatedX int, updatedY int) error {
	log.Println("Updating shape...")
	newXVal := updatedX
	newYVal := updatedY
	if err := edit.UpdateShape(0, shapeID, common.ShapeUpdate{
		Text: &newText,
		X:    &newXVal,
		Y:    &newYVal,
	}); err != nil {
		return fmt.Errorf("UpdateShape failed: %w", err)
	}
	return nil
}

func verifyEditedShape(outputFile string, expectedText string, expectedX int, expectedY int) error {
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
	if !hasEditedShape(vShapes, expectedText, expectedX, expectedY) {
		return fmt.Errorf(
			"verification failed: could not find shape with %q at %d,%d",
			expectedText,
			expectedX,
			expectedY,
		)
	}
	return nil
}

func hasEditedShape(shapes []common.Shape, expectedText string, expectedX int, expectedY int) bool {
	for _, s := range shapes {
		log.Printf("Verification Shape: ID=%d Text=%q X=%d Y=%d\n", s.ID, s.Text, s.X, s.Y)
		if s.Text == expectedText && s.X == expectedX && s.Y == expectedY {
			return true
		}
	}
	return false
}
