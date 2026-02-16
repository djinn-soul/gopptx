package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/gopptx"
)

func main() {
	pres := &gopptx.Presentation{}
	pres.AddSlide()

	outDir := "examples/output"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	filename := filepath.Join(outDir, "01_hello_world.pptx")
	err := pres.Save(filename)
	if err != nil {
		log.Fatalf("Failed to save presentation: %v", err)
	}

	log.Printf("Successfully generated %s\n", filename)
}
