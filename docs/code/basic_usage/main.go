// basic_usage demonstrates creating a simple PowerPoint presentation
// with a title slide using gopptx.
//
// Run: go run ./docs/code/basic_usage/
// Output: docs/assets/pptx/basic_usage.pptx
package main

import (
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/gopptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

const (
	outFile = "docs/assets/pptx/basic_usage.pptx"
	boxX    = 0.5
	boxY    = 4.5
	boxW    = 9.0
	boxH    = 0.9
)

func main() {
	if err := os.MkdirAll(filepath.Dir(outFile), 0o750); err != nil {
		log.Fatalf("create output dir: %v", err)
	}

	pres := &gopptx.Presentation{Title: "gopptx – Basic Usage"}

	slide := pres.AddSlide()
	slide.Title = "Hello from gopptx"
	slide.AddBullet("Create slides programmatically in Go.")
	slide.AddBullet("Add shapes, text, charts, images, and more.")
	slide.AddBullet("Export to .pptx in milliseconds.")

	// Add a coloured call-out box
	box := shapes.NewRectangle(boxX, boxY, boxW, boxH).
		WithText("Open-source • High-performance • Go + Python").
		WithFill(shapes.NewShapeFill("2E4057"))
	slide.AddShape(box)

	if err := pres.Save(outFile); err != nil {
		log.Fatalf("save: %v", err)
	}
	log.Printf("Saved %s", outFile)
}
