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
			WithFill(shapes.NewShapeFill("455A64")).
			WithText("Open README.md (Relative)").
			WithClickAction(action.NewHyperlink(action.HyperlinkFile("README.md")).
				WithTooltip("Open project README"))).
		AddShape(shapes.NewShape("rect", styling.Inches(5), styling.Inches(1.5), styling.Inches(4), styling.Inches(1)).
			WithFill(shapes.NewShapeFill("37474F")).
			WithText("Open Notepad (Program)").
			WithClickAction(action.NewHyperlink(action.HyperlinkProgram("C:\\Windows\\System32\\notepad.exe")).
				WithTooltip("Launch Notepad application")))
}
