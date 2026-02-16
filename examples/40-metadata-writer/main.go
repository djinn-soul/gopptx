package main

import (
	"archive/zip"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
)

func main() {
	const outputDir = "examples/output"
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatalf("Failed to create output dir: %v", err)
	}

	tmpDir, err := os.MkdirTemp("", "gopptx-metadata-writer-*")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			log.Printf("warning: failed to remove temp dir %s: %v", tmpDir, err)
		}
	}()

	inputFile := filepath.Join(tmpDir, "metadata_input.pptx")
	outputFile := filepath.Join(outputDir, "40_metadata_output.pptx")

	// 1. Create a minimal valid PPTX file
	log.Printf("Generating minimal PPTX: %s...\n", inputFile)
	if err := createMinimalPPTX(inputFile); err != nil {
		log.Fatalf("Failed to create minimal PPTX: %v", err)
	}
	defer func() {
		// optional cleanup
		if err := os.Remove(inputFile); err != nil && !os.IsNotExist(err) {
			log.Printf("warning: failed to remove input file %s: %v", inputFile, err)
		}
	}()

	// 2. Open it
	log.Printf("Opening %s...\n", inputFile)
	ppt, err := editor.OpenPresentationEditor(inputFile)
	if err != nil {
		log.Fatalf("Failed to open presentation: %v", err)
	}
	defer func() { _ = ppt.Close() }()

	// 3. Check initial metadata
	props := ppt.GetCoreProperties()
	log.Printf("Initial Title: %s\n", props.Title)
	if props.Title != "Initial Title" {
		log.Fatalf("Expected 'Initial Title', got '%s'", props.Title)
	}

	// 4. Update metadata
	log.Println("Updating metadata...")
	newProps := common.CoreProperties{
		Title:       "Updated Title",
		Subject:     "Updated Subject",
		Creator:     "Updated Creator",
		Description: "Updated Description",
		Keywords:    "test, metadata",
	}
	ppt.SetCoreProperties(newProps)

	// 5. Save
	log.Printf("Saving to %s...\n", outputFile)
	if err := ppt.Save(outputFile); err != nil {
		log.Fatalf("Failed to save: %v", err)
	}

	// 6. Verify output
	verifyOutput(outputFile)

	log.Println("Done! Smoke test passed.")
}

func verifyOutput(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open output file: %v", err)
	}
	defer func() { _ = f.Close() }()

	fi, _ := f.Stat()
	z, err := zip.NewReader(f, fi.Size())
	if err != nil {
		log.Fatalf("Failed to open zip: %v", err)
	}

	foundCore := false
	for _, f := range z.File {
		if f.Name == "docProps/core.xml" {
			foundCore = true
			rc, _ := f.Open()
			content := make([]byte, f.UncompressedSize64)
			_, _ = rc.Read(content)
			_ = rc.Close()
			s := string(content)
			if !contains(s, "Updated Title") {
				log.Fatalf("Output missing updated title")
			}
			if !contains(s, "Updated Description") {
				log.Fatalf("Output missing updated description")
			}
		}
	}
	if !foundCore {
		log.Fatalf("Output missing docProps/core.xml")
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func createMinimalPPTX(filename string) error {
	meta := pptx.PresentationMetadata{
		PresentationMetadata: pptx.PresentationMetadataFields{
			Title:   "Initial Title",
			Creator: "Initial Creator",
		},
	}
	data, err := pptx.CreateWithMetadata(meta, []pptx.SlideContent{
		pptx.NewSlide("Metadata Base"),
	})
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0o644)
}
