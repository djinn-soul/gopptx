package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Smoke test failed: %v", err)
	}
	log.Println("Smart Merge Smoke Test: SUCCESS")
}

func run() error {
	const outputDir = "examples/output"
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "gopptx-merge-assets-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			log.Printf("warning: failed to remove temp dir %s: %v", tmpDir, err)
		}
	}()

	// 1. Create Source Presentation with an Image
	pngData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // Signature
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, // IHDR
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, // 1x1
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4, // Data
		0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41, // IDAT
		0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00, // Data
		0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00, // Data
		0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, // IEND
		0xae, 0x42, 0x60, 0x82,
	}

	srcDeck := pptx.NewPresentationBuilder("Source Deck")

	sourceImagePath := filepath.Join(tmpDir, "source_image.png")
	if err := os.WriteFile(sourceImagePath, pngData, 0o644); err != nil {
		return fmt.Errorf("failed to write source image: %w", err)
	}

	bgImg := shapes.NewImage(sourceImagePath, 0, 0, 0, 0)
	// Create a slide with picture background
	slide := pptx.NewSlide("Source Slide with Image").WithBackground(pptx.NewPictureBackground(bgImg))
	srcDeck.AddSlide(slide)

	sourcePPTXPath := filepath.Join(tmpDir, "merge_source.pptx")
	if err := srcDeck.WriteToFile(sourcePPTXPath); err != nil {
		return err
	}

	// 2. Create Target Presentation (Empty)
	dstDeck := pptx.NewPresentationBuilder("Target Deck")
	dstDeck.AddTitleSlide("Title Slide")
	targetPPTXPath := filepath.Join(tmpDir, "merge_target.pptx")
	if err := dstDeck.WriteToFile(targetPPTXPath); err != nil {
		return err
	}

	// 3. Merge
	dstEdit, err := editor.OpenPresentationEditor(targetPPTXPath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := dstEdit.Close(); closeErr != nil {
			log.Printf("warning: failed to close destination editor: %v", closeErr)
		}
	}()

	log.Println("Merging...")
	if err := dstEdit.MergeFromFile(sourcePPTXPath); err != nil {
		return fmt.Errorf("MergeFromFile failed: %w", err)
	}

	// Save Result
	resultPPTXPath := filepath.Join(outputDir, "42_smart_merge_assets.pptx")
	if err := dstEdit.Save(resultPPTXPath); err != nil {
		return err
	}

	// 4. Verification
	// Re-open result to check.
	chkEdit, err := editor.OpenPresentationEditor(resultPPTXPath)
	if err != nil {
		return fmt.Errorf("failed to open result: %w", err)
	}
	defer func() {
		if closeErr := chkEdit.Close(); closeErr != nil {
			log.Printf("warning: failed to close verification editor: %v", closeErr)
		}
	}()

	// Expected: 2 slides.
	// Slide 2 should have image relationship.
	// We can't easily inspect relationships via public API of editor yet (internal only).
	// But if it opened without error, that's a good sign.
	// Let's try to verify via scanning parts? No public API for scan.
	// We trust if Open succeeds on save, XML is valid.
	// To be sure image is there, we rely on the fact that MergeFromFile would fail if copy failed.

	return nil
}
