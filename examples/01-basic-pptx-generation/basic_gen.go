package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/gopptx"
)

func main() {
	pres := &gopptx.Presentation{}
	pres.AddSlide()

	outDir := "examples/output"
	if err := os.MkdirAll(outDir, 0o750); err != nil {
		fmt.Printf("Failed to create output directory: %v\n", err)
		os.Exit(1)
	}

	filename := filepath.Join(outDir, "01_hello_world.pptx")
	err := pres.Save(filename)
	if err != nil {
		fmt.Printf("Failed to save presentation: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated %s\n", filename)
}
