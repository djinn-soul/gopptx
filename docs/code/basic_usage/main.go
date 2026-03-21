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
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const outFile = "docs/assets/pptx/basic_usage.pptx"

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
	box := shapes.NewRectangle(0.5, 4.5, 9.0, 0.9).
		WithText("Open-source • High-performance • Go + Python").
		WithFill(shapes.NewShapeFill("2E4057"))
	slide.AddShape(box)

	// Tint the box text white
	_ = styling.Inches(0) // keep import used

	if err := pres.Save(outFile); err != nil {
		log.Fatalf("save: %v", err)
	}
	log.Printf("Saved %s", outFile)
}
