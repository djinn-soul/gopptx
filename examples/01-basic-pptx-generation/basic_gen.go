package main

import (
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func main() {
	outDir := "examples/output"
	if err := os.MkdirAll(outDir, 0o750); err != nil {
		log.Printf("Failed to create output directory: %v", err)
		os.Exit(1)
	}

	filename := filepath.Join(outDir, "01_hello_world.pptx")
	err := pptx.NewPresentationBuilder("Presentation").
		AddSlide(pptx.NewSlide("Slide 1")).
		WriteToFile(filename)
	if err != nil {
		log.Printf("Failed to save presentation: %v", err)
		os.Exit(1)
	}

	log.Printf("Successfully generated %s", filename)
}
