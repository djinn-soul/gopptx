package main

import (
	"fmt"
	"log"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/handout"
)

const outputDir = "examples/output"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	slides := []pptx.SlideContent{
		pptx.NewSlide("Handout Master Demo").AddBullet("Demonstrates handout master support."),
		pptx.NewSlide("Slide 2").AddBullet("More content here."),
		pptx.NewSlide("Slide 3").AddBullet("Third slide."),
	}

	// --- 1. Default handout master (1-up, all placeholders visible) ---
	if err := pptx.NewPresentationBuilder("Handout Default").
		WithMetadata(pptx.Metadata{HandoutMaster: handout.New()}).
		AddSlide(slides[0]).AddSlide(slides[1]).AddSlide(slides[2]).
		WriteToFile(outputDir + "/29_handout_master_default.pptx"); err != nil {
		return fmt.Errorf("save default: %w", err)
	}
	log.Println("Saved: 29_handout_master_default.pptx")

	// --- 2. 6-up layout with custom header/footer, date hidden ---
	if err := pptx.NewPresentationBuilder("Handout 6-Up").
		WithMetadata(pptx.Metadata{
			HandoutMaster: handout.New().
				WithLayout(handout.Layout6Up).
				WithHeader("Acme Corp — Q1 Review").
				WithFooter("Confidential").
				HideDate(),
		}).
		AddSlide(slides[0]).AddSlide(slides[1]).AddSlide(slides[2]).
		WriteToFile(outputDir + "/29_handout_master_6up.pptx"); err != nil {
		return fmt.Errorf("save 6-up: %w", err)
	}
	log.Println("Saved: 29_handout_master_6up.pptx")

	// --- 3. Outline layout, all placeholders hidden ---
	if err := pptx.NewPresentationBuilder("Handout Outline").
		WithMetadata(pptx.Metadata{
			HandoutMaster: handout.New().
				WithLayout(handout.LayoutOutline).
				HideHeader().
				HideFooter().
				HideDate().
				HidePageNumber(),
		}).
		AddSlide(slides[0]).
		WriteToFile(outputDir + "/29_handout_master_outline.pptx"); err != nil {
		return fmt.Errorf("save outline: %w", err)
	}
	log.Println("Saved: 29_handout_master_outline.pptx")

	return nil
}
