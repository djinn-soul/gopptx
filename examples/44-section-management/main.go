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
	log.Println("Section Management Smoke Test: SUCCESS")
	log.Println(
		"To verify visually: Open examples/output/44_section_management.pptx in PowerPoint and switch to 'Slide Sorter' view.",
	)
}

func run() error {
	const outputDir = "examples/output"
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	tmpDir, tempErr := os.MkdirTemp("", "gopptx-section-management-*")
	if tempErr != nil {
		return fmt.Errorf("create temp dir: %w", tempErr)
	}
	defer func() {
		if cleanupErr := os.RemoveAll(tmpDir); cleanupErr != nil {
			log.Printf("warning: failed to remove temp dir %s: %v", tmpDir, cleanupErr)
		}
	}()

	inputFile := filepath.Join(tmpDir, "section_test_input.pptx")
	outFile := filepath.Join(outputDir, "44_section_management.pptx")

	if err := createInputDeck(inputFile); err != nil {
		return err
	}

	edit, err := editor.OpenPresentationEditor(inputFile)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := edit.Close(); closeErr != nil {
			log.Printf("warning: failed to close editor: %v", closeErr)
		}
	}()

	if err := addAndRenameSections(edit); err != nil {
		return err
	}

	if err := edit.Save(outFile); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	if err := verifySections(outFile); err != nil {
		return err
	}
	return nil
}

func createInputDeck(inputFile string) error {
	deck := pptx.NewPresentationBuilder("Section Demo")
	deck.AddTitleSlide("Intro Slide")
	deck.AddTitleSlide("Detail 1")
	deck.AddTitleSlide("Detail 2")
	deck.AddTitleSlide("Appendix")
	if err := deck.WriteToFile(inputFile); err != nil {
		return err
	}
	return nil
}

func addAndRenameSections(edit *editor.PresentationEditor) error {
	log.Println("Adding 'Introduction' section...")
	if err := edit.AddSection("Introduction", []int{0}); err != nil {
		return fmt.Errorf("failed to add intro section: %w", err)
	}

	log.Println("Adding 'Core Content' section...")
	if err := edit.AddSection("Core Content", []int{1, 2}); err != nil {
		return fmt.Errorf("failed to add core section: %w", err)
	}

	log.Println("Adding 'Appendix' section...")
	if err := edit.AddSection("Appendix", []int{3}); err != nil {
		return fmt.Errorf("failed to add appendix section: %w", err)
	}

	log.Println("Renaming 'Appendix' to 'Back Matter'...")
	if err := edit.RenameSection("Appendix", "Back Matter"); err != nil {
		return fmt.Errorf("failed to rename section: %w", err)
	}
	return nil
}

func verifySections(outFile string) error {
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
	log.Printf("Found %d sections:\n", len(sections))
	for _, s := range sections {
		log.Printf("- %s (ID: %s, Slides: %d)\n", s.Name, s.GUID, len(s.SlideIDs))
	}

	if len(sections) != 3 {
		return fmt.Errorf("expected 3 sections, got %d", len(sections))
	}
	if sections[2].Name != "Back Matter" {
		return fmt.Errorf("expected last section to be 'Back Matter', got %q", sections[2].Name)
	}
	return nil
}
