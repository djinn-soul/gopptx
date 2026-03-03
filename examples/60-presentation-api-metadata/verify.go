package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func main() {
	outDir := filepath.Join("examples", "output")
	examplePath := filepath.Join(outDir, "60_metadata_persistence.pptx")
	copyPath := filepath.Join(outDir, "60_metadata_copy.pptx")

	// First create a sample presentation
	if err := createSamplePresentation(examplePath); err != nil {
		log.Fatalf("failed to create sample presentation: %v", err)
	}

	fmt.Println("=== Opening a presentation ===")
	prs, err := pptx.Open(examplePath)
	if err != nil {
		log.Fatalf("failed to open presentation: %v", err)
	}

	prs.SetTitle("Updated Presentation Title")
	prs.SetAuthor("Jane Doe")
	prs.SetKeywords("presentation, go, pptx, metadata")

	if err := prs.Save(); err != nil {
		log.Fatalf("failed to save presentation: %v", err)
	}
	prs.Close()
	fmt.Println("Changes saved!")

	// Create a copy
	prs2, err := pptx.Open(examplePath)
	if err != nil {
		log.Fatalf("failed to reopen: %v", err)
	}
	if err := prs2.SaveAs(copyPath); err != nil {
		log.Fatalf("failed to save copy: %v", err)
	}
	prs2.Close()
	fmt.Printf("Saved copy to %s\n", copyPath)
}

func createSamplePresentation(path string) error {
	data, err := pptx.Create("Sample Presentation", 3)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
