package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func main() {
	slides := []pptx.SlideContent{
		buildAdvancedHyperlinkSlide(),
	}

	data, buildErr := pptx.CreateWithSlides("Advanced Hyperlink Smoke Test", slides)
	if buildErr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", buildErr)
		os.Exit(1)
	}

	const outputDir = "examples/output"
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir error: %v\n", err)
		os.Exit(1)
	}
	outPath := filepath.Join(outputDir, "31_advanced_hyperlink_smoke.pptx")
	if err := os.WriteFile(outPath, data, 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
		os.Exit(1)
	}
	log.Printf("Created %s\n", outPath)
}

func buildAdvancedHyperlinkSlide() pptx.SlideContent {
	return pptx.NewSlide("Advanced Hyperlinks").
		AddShape(shapes.NewShape("rect", styling.Inches(0.5), styling.Inches(1.5), styling.Inches(4), styling.Inches(1)).
			WithFill(shapes.NewShapeFill("DDE7F0")).
			WithClickAction(action.NewHyperlink(action.HyperlinkFile("README.md")).
				WithTooltip("Open project README"))).
		AddShape(shapes.NewTextBox("Open README.md (Relative)", 0.75, 1.8, 3.5, 0.35)).
		AddShape(shapes.NewShape("rect", styling.Inches(5), styling.Inches(1.5), styling.Inches(4), styling.Inches(1)).
			WithFill(shapes.NewShapeFill("E8F1E6")).
			WithClickAction(action.NewHyperlink(action.HyperlinkProgram("notepad.exe")).
				WithTooltip("Launch Notepad application"))).
		AddShape(shapes.NewTextBox("Open Notepad (Program)", 5.25, 1.8, 3.5, 0.35))
}
