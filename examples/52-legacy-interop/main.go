package main

import (
	"fmt"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	log "github.com/djinn-soul/gopptx/pkg/stdlog"
)

const (
	outputDir  = "examples/output"
	outputFile = "52_legacy_interop.pptx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

//nolint:gosec // Example writes temp files with non-sensitive fixture content.
func run() error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	// 1. Create a presentation and save it.
	original := []pptx.SlideContent{
		pptx.NewSlide("Original Slide").
			AddBullet("This slide will survive a round-trip").
			AddBullet("Through the editor and back"),
		pptx.NewSlide("Second Slide").
			AddBullet("Unknown XML parts are preserved").
			AddBullet("Only touched parts are rewritten"),
	}

	data, err := pptx.CreateWithSlides("Legacy Interop Demo", original)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}

	tmpPath := os.TempDir() + "/gopptx_52_interop_src.pptx"
	if err = os.WriteFile(tmpPath, data, 0o644); err != nil {
		return fmt.Errorf("write tmp: %w", err)
	}
	defer func() { _ = os.Remove(tmpPath) }()

	// 2. Open with the editor (simulates working with an existing / legacy file).
	ed, err := pptx.OpenPresentationEditor(tmpPath)
	if err != nil {
		return fmt.Errorf("open editor: %w", err)
	}
	defer func() { _ = ed.Close() }()

	// 3. Add a new slide without touching the originals.
	newSlide := pptx.NewSlide("Added via Editor").AddBullet("Appended after round-trip")
	if _, err = ed.AddSlide(newSlide); err != nil {
		return fmt.Errorf("add slide: %w", err)
	}

	count := ed.SlideCount()
	log.Printf("Slide count after round-trip: %d\n", count)

	// 4. Save the modified file.
	outputPath := outputDir + "/" + outputFile
	if err = ed.Save(outputPath); err != nil {
		return fmt.Errorf("save: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}
