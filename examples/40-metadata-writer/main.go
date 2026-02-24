package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	const outputDir = "examples/output"
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	tmpDir, tempErr := os.MkdirTemp("", "gopptx-metadata-writer-*")
	if tempErr != nil {
		return fmt.Errorf("failed to create temp dir: %w", tempErr)
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
		return fmt.Errorf("failed to create minimal PPTX: %w", err)
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
		return fmt.Errorf("failed to open presentation: %w", err)
	}
	defer func() { _ = ppt.Close() }()

	// 3. Check initial metadata
	props := ppt.GetCoreProperties()
	log.Printf("Initial Title: %s\n", props.Title)
	if props.Title != "Initial Title" {
		return fmt.Errorf("expected 'Initial Title', got %q", props.Title)
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
		return fmt.Errorf("failed to save: %w", err)
	}

	// 6. Verify output
	if err := verifyOutput(outputFile); err != nil {
		return err
	}

	log.Println("Done! Smoke test passed.")
	return nil
}

func verifyOutput(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer func() { _ = f.Close() }()

	fi, statErr := f.Stat()
	if statErr != nil {
		return fmt.Errorf("failed to stat output file: %w", statErr)
	}
	z, err := zip.NewReader(f, fi.Size())
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}

	foundCore := false
	for _, zf := range z.File {
		if zf.Name != "docProps/core.xml" {
			continue
		}
		foundCore = true
		rc, openErr := zf.Open()
		if openErr != nil {
			return fmt.Errorf("failed to open docProps/core.xml: %w", openErr)
		}
		content, readErr := io.ReadAll(rc)
		closeErr := rc.Close()
		if readErr != nil {
			return fmt.Errorf("failed to read docProps/core.xml: %w", readErr)
		}
		if closeErr != nil {
			return fmt.Errorf("failed to close docProps/core.xml stream: %w", closeErr)
		}
		s := string(content)
		if !contains(s, "Updated Title") {
			return errors.New("output missing updated title")
		}
		if !contains(s, "Updated Description") {
			return errors.New("output missing updated description")
		}
	}
	if !foundCore {
		return errors.New("output missing docProps/core.xml")
	}
	return nil
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func createMinimalPPTX(filename string) error {
	meta := pptx.Metadata{
		Metadata: pptx.MetadataFields{
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
	return os.WriteFile(filename, data, 0o600)
}
