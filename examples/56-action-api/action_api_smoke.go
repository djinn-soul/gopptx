//go:build ignore

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func main() {
	builder := pptx.NewPresentationBuilder("Action API Showcase")

	// Slide 1: Shape Actions
	actionURL := pptx.HyperlinkURL("https://github.com/djinn-soul/gopptx")
	actionNext := pptx.HyperlinkNextSlide()

	slide1 := pptx.NewSlide("Shape Move & Hover").
		AddShape(pptx.NewRectangle(1, 1, 3, 2).
			WithText("Click Me (URL)").
			WithFill(pptx.NewShapeFill("00FF00")).
			WithClickAction(pptx.NewHyperlink(actionURL))).
		AddShape(pptx.NewEllipse(5, 1, 3, 2).
			WithText("Hover Me (Next Slide)").
			WithFill(pptx.NewShapeFill("FFFF00")).
			WithHoverAction(pptx.NewHyperlink(actionNext)))
	builder.AddSlide(slide1)

	// Slide 2: Text Run Actions
	slide2 := pptx.NewSlide("Text Run Actions")
	run1 := pptx.NewTextRun("This is a ").WithBold(true)
	run2 := pptx.NewTextRun("clickable/hoverable").
		WithColor("FF0000").
		WithUnderline(true).
		WithHyperlink(pptx.NewHyperlink(actionURL)).
		WithHoverAction(pptx.NewHyperlink(actionNext))
	run3 := pptx.NewTextRun(" word.")

	slide2 = slide2.AddBulletRuns([]pptx.TextRun{run1, run2, run3})
	builder.AddSlide(slide2)

	// Save
	outDir := "examples/output"
	if err := os.MkdirAll(outDir, 0o750); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}
	outPath := filepath.Join(outDir, "56_action_api_smoke.pptx")
	if err := builder.WriteToFile(outPath); err != nil {
		log.Fatalf("Failed to save presentation: %v", err)
	}

	log.Printf("Successfully generated smoke sample: %s", outPath)
}
