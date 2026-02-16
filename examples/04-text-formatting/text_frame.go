//go:build ignore

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func main() { //nolint:unused
	//nolint:unused
	outDir := "examples/output"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	slide := pptx.NewSlide("Text Frame Properties").
		AddBullet("Custom Margins").
		AddBullet("Vertical Alignment (Anchor)").
		AddBullet("Word Wrap toggling").
		AddBullet("Auto-fit behaviors")

	// 1. Large Margin
	slide.AddShape(
		pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(0.5), pptx.Inches(2), pptx.Inches(2), pptx.Inches(1.5)).
			WithText("Large internal margins (0.5 in)").
			WithFill(pptx.NewShapeFill("FFC000")).
			WithTextMargins(pptx.Inches(0.5), pptx.Inches(0.5), pptx.Inches(0.5), pptx.Inches(0.5)),
	)

	// 2. Top Anchor
	slide.AddShape(
		pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(3), pptx.Inches(2), pptx.Inches(2), pptx.Inches(1.5)).
			WithText("Top Anchored Text").
			WithFill(pptx.NewShapeFill("5B9BD5")).
			WithVerticalAnchor(pptx.TextAnchorTop),
	)

	// 3. Bottom Anchor
	slide.AddShape(
		pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(5.5), pptx.Inches(2), pptx.Inches(2), pptx.Inches(1.5)).
			WithText("Bottom Anchored Text").
			WithFill(pptx.NewShapeFill("70AD47")).
			WithVerticalAnchor(pptx.TextAnchorBottom),
	)

	// 4. No Wrap
	slide.AddShape(
		pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(0.5), pptx.Inches(4), pptx.Inches(2), pptx.Inches(0.5)).
			WithText("This text should NOT wrap and spill out").
			WithFill(pptx.NewShapeFill("ED7D31")).
			WithTextWrap(pptx.TextWrapNone),
	)

	// 5. Shrink Text (normAutoFit)
	slide.AddShape(
		pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(3), pptx.Inches(4), pptx.Inches(2), pptx.Inches(1)).
			WithText("This is a lot of text that should shrink to fit inside the box without expanding the box itself.").
			WithFill(pptx.NewShapeFill("A5A5A5")).
			WithAutoFit(pptx.TextAutoFitNormal),
	)

	data, buildErr := pptx.CreateWithSlides("Text Frame Smoke", []pptx.SlideContent{slide})
	if buildErr != nil {
		log.Fatalf("Failed to generate PPTX: %v", buildErr)
	}

	path := filepath.Join(outDir, "04_text_frame_smoke.pptx")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	log.Printf("Successfully generated smoke sample: %s\n", path)
}
