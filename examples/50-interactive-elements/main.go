package main

import (
	"fmt"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	log "github.com/djinn-soul/gopptx/pkg/stdlog"
)

const (
	outputDir  = "examples/output"
	outputFile = "50_interactive_elements.pptx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	slides := []pptx.SlideContent{
		// Slide 1: navigation buttons.
		pptx.NewSlide("Navigation Buttons").
			AddShape(pptx.NewShape(pptx.ShapeTypeRoundedRectangle,
				pptx.Inches(1), pptx.Inches(2), pptx.Inches(3), pptx.Emu(800000)).
				WithFill(pptx.NewShapeFill("0078D4")).
				WithText("← Previous Slide").
				WithClickAction(pptx.NewHyperlink(pptx.HyperlinkPreviousSlide()))).
			AddShape(pptx.NewShape(pptx.ShapeTypeRoundedRectangle,
				pptx.Inches(5.5), pptx.Inches(2), pptx.Inches(3), pptx.Emu(800000)).
				WithFill(pptx.NewShapeFill("107C10")).
				WithText("Next Slide →").
				WithClickAction(pptx.NewHyperlink(pptx.HyperlinkNextSlide()))).
			AddBullet("Click shapes to navigate between slides"),

		// Slide 2: URL and email hyperlinks.
		pptx.NewSlide("URL & Email Links").
			AddShape(pptx.NewShape(pptx.ShapeTypeRoundedRectangle,
				pptx.Inches(1), pptx.Inches(1.5), pptx.Inches(4), pptx.Emu(800000)).
				WithFill(pptx.NewShapeFill("D83B01")).
				WithText("Open gopptx on GitHub").
				WithClickAction(pptx.NewHyperlink(pptx.HyperlinkURL("https://github.com/djinn-soul/gopptx")).
					WithTooltip("View the gopptx repository"))).
			AddShape(pptx.NewShape(pptx.ShapeTypeRoundedRectangle,
				pptx.Inches(1), pptx.Inches(3), pptx.Inches(4), pptx.Emu(800000)).
				WithFill(pptx.NewShapeFill("5B9BD5")).
				WithText("Send Email").
				WithClickAction(pptx.NewHyperlink(pptx.HyperlinkURL("mailto:hello@example.com")))),

		// Slide 3: text-run hyperlinks.
		pptx.NewSlide("Inline Text Hyperlinks").
			AddBulletRuns([]pptx.Run{
				pptx.NewRun("Visit "),
				pptx.NewRun("the gopptx repo").
					WithHyperlink(pptx.NewHyperlink(pptx.HyperlinkURL("https://github.com/djinn-soul/gopptx"))).
					WithColor("0563C1").
					WithUnderline(true),
				pptx.NewRun(" for documentation."),
			}).
			AddBullet("Inline hyperlinks on text runs").
			AddBullet("URL and email targets supported"),
	}

	data, err := pptx.CreateWithSlides("Task 50: Interactive Elements", slides)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}

	outputPath := outputDir + "/" + outputFile
	if err = os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}
