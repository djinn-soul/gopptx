package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Smoke test failed: %v", err)
	}
	fmt.Println("Section Management Smoke Test: SUCCESS")
	fmt.Println("To verify visually: Open examples/output/44_section_management.pptx in PowerPoint and switch to 'Slide Sorter' view.")
}

func run() error {
	const outputDir = "examples/output"
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "gopptx-section-management-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer func() {
		if cleanupErr := os.RemoveAll(tmpDir); cleanupErr != nil {
			log.Printf("warning: failed to remove temp dir %s: %v", tmpDir, cleanupErr)
		}
	}()

	inputFile := filepath.Join(tmpDir, "section_test_input.pptx")
	outFile := filepath.Join(outputDir, "44_section_management.pptx")

	// 1. Create a presentation with multiple slides
	deck := pptx.NewPresentationBuilder("Section Demo")
	deck.AddTitleSlide("Intro Slide") // Index 0
	deck.AddTitleSlide("Detail 1")    // Index 1
	deck.AddTitleSlide("Detail 2")    // Index 2
	deck.AddTitleSlide("Appendix")    // Index 3

	if err := deck.WriteToFile(inputFile); err != nil {
		return err
	}

	// 2. Open editor and add sections
	// Sections must be contiguous logic usually, but the API allows arbitrary indices.
	// PPT usually requires sections to cover all slides or start from a point.
	// Let's try to group them:
	// Section 1: Intro (Index 0)
	// Section 2: Details (Index 1, 2)
	// Section 3: Appendix (Index 3)

	edit, err := editor.OpenPresentationEditor(inputFile)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := edit.Close(); closeErr != nil {
			log.Printf("warning: failed to close editor: %v", closeErr)
		}
	}()

	fmt.Println("Adding 'Introduction' section...")
	if err := edit.AddSection("Introduction", []int{0}); err != nil {
		return fmt.Errorf("failed to add intro section: %w", err)
	}

	fmt.Println("Adding 'Core Content' section...")
	if err := edit.AddSection("Core Content", []int{1, 2}); err != nil {
		return fmt.Errorf("failed to add core section: %w", err)
	}

	fmt.Println("Adding 'Appendix' section...")
	if err := edit.AddSection("Appendix", []int{3}); err != nil {
		return fmt.Errorf("failed to add appendix section: %w", err)
	}
	fmt.Println("Renaming 'Appendix' to 'Back Matter'...")
	if err := edit.RenameSection("Appendix", "Back Matter"); err != nil {
		return fmt.Errorf("failed to rename section: %w", err)
	}

	// 3. Save
	if err := edit.Save(outFile); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	// 4. Verification Check
	// Re-open and check internal state using the Sections() API.
	checkEdit, err := editor.OpenPresentationEditor(outFile)
	if err != nil {
		return fmt.Errorf("failed to open result: %w", err)
	}
	defer func() {
		if closeErr := checkEdit.Close(); closeErr != nil {
			log.Printf("warning: failed to close verification editor: %v", closeErr)
		}
	}()

	sections := checkEdit.Sections()
	fmt.Printf("Found %d sections:\n", len(sections))
	for _, s := range sections {
		fmt.Printf("- %s (ID: %s, Slides: %d)\n", s.Name, s.GUID, len(s.SlideIDs))
	}

	if len(sections) != 3 {
		return fmt.Errorf("expected 3 sections, got %d", len(sections))
	}
	if sections[2].Name != "Back Matter" {
		return fmt.Errorf("expected last section to be 'Back Matter', got %q", sections[2].Name)
	}

	return nil
}
