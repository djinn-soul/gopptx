package main

import (
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/gopptx"
)

func main() {
	pres := &gopptx.Presentation{}
	pres.AddSlide()

	outDir := "examples/output"
	if err := os.MkdirAll(outDir, 0o750); err != nil {
		log.Printf("Failed to create output directory: %v", err)
		os.Exit(1)
	}

	filename := filepath.Join(outDir, "01_hello_world.pptx")
	err := pres.Save(filename)
	if err != nil {
		log.Printf("Failed to save presentation: %v", err)
		os.Exit(1)
	}

	log.Printf("Successfully generated %s", filename)
}
