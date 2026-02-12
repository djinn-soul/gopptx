package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func main() {
	builder := pptx.NewPresentationBuilder("Text Enhancements Showcase")

	// Slide 1: Case Controls
	slide1 := pptx.NewSlide("Text Case Controls").
		AddBullet("This corresponds to normal text.").
		AddBulletRuns([]pptx.TextRun{
			pptx.NewTextRun("This is ").WithBold(true),
			pptx.NewTextRun("ALL CAPS").WithAllCaps(true).WithColor("FF0000"),
			pptx.NewTextRun(" text."),
		}).
		AddBulletRuns([]pptx.TextRun{
			pptx.NewTextRun("This is ").WithItalic(true),
			pptx.NewTextRun("Small Caps").WithSmallCaps(true).WithColor("0000FF"),
			pptx.NewTextRun(" text."),
		})
	builder.AddSlide(slide1)

	// Slide 2: Paragraph Indents
	// Note: Inches helper is available in root pptx package via styling_compat.go
	slide2 := pptx.NewSlide("Paragraph Indents").
		AddBullet("Default indentation for comparison (Level 0).").
		AddBulletWithStyle("Custom Left Indent (2 inches)",
			pptx.DefaultTextParagraphStyle().WithLeftIndent(pptx.Inches(2))).
		AddBulletWithStyle("Hanging Indent (0.5 inch offset)",
			pptx.DefaultTextParagraphStyle().
				WithLeftIndent(pptx.Inches(1)).
				WithHangingIndent(pptx.Inches(-0.5)))

	builder.AddSlide(slide2)

	// Save
	outDir := "smoke_samples"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}
	outPath := filepath.Join(outDir, "gopptx_text_enhancements_smoke.pptx")
	if err := builder.WriteToFile(outPath); err != nil {
		log.Fatalf("Failed to save presentation: %v", err)
	}

	log.Printf("Successfully generated smoke sample: %s", outPath)
}
