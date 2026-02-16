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
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

func main() {
	slides := []pptx.SlideContent{
		buildClickActionSlide(),
		buildHoverActionSlide(),
		buildTextRunHyperlinkSlide(),
	}

	data, buildErr := pptx.CreateWithSlides("Action & Hyperlink Smoke Test", slides)
	if buildErr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", buildErr)
		os.Exit(1)
	}

	const outputDir = "examples/output"
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir error: %v\n", err)
		os.Exit(1)
	}
	outPath := filepath.Join(outputDir, "56_action_smoke.pptx")
	if err := os.WriteFile(outPath, data, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
		os.Exit(1)
	}
	log.Printf("Created %s\n", outPath)
}

func buildClickActionSlide() pptx.SlideContent {
	return pptx.NewSlide("Click Actions").
		AddShape(shapes.NewShape("roundRect", styling.Inches(0.5), styling.Inches(1.5), styling.Inches(3), styling.Inches(0.8)).
			WithFill(shapes.NewShapeFill("1565C0")).
			WithText("Open Example.com").
			WithClickAction(action.NewHyperlink(action.HyperlinkURL("https://example.com")).
				WithTooltip("Visit example.com"))).
		AddShape(shapes.NewShape("rect", styling.Inches(4), styling.Inches(1.5), styling.Inches(3), styling.Inches(0.8)).
			WithFill(shapes.NewShapeFill("2E7D32")).
			WithText("Go to Slide 3").
			WithClickAction(action.NewHyperlink(action.HyperlinkSlide(3)))).
		AddShape(shapes.NewShape("rect", styling.Inches(0.5), styling.Inches(3), styling.Inches(3), styling.Inches(0.8)).
			WithFill(shapes.NewShapeFill("E65100")).
			WithText("Send Email").
			WithClickAction(action.NewHyperlink(action.HyperlinkEmailWithSubject("test@example.com", "Hello from gopptx"))),
		).
		AddShape(shapes.NewShape("rect", styling.Inches(4), styling.Inches(3), styling.Inches(3), styling.Inches(0.8)).
			WithFill(shapes.NewShapeFill("6A1B9A")).
			WithText("Next Slide →").
			WithClickAction(action.NewHyperlink(action.HyperlinkNextSlide())))
}

func buildHoverActionSlide() pptx.SlideContent {
	return pptx.NewSlide("Hover Actions").
		AddShape(shapes.NewShape("roundRect", styling.Inches(1), styling.Inches(2), styling.Inches(4), styling.Inches(1)).
			WithFill(shapes.NewShapeFill("0288D1")).
			WithText("Hover me for tooltip").
			WithHoverAction(action.NewHyperlink(action.HyperlinkURL("https://example.com")).
				WithTooltip("You hovered!"))).
		AddShape(shapes.NewShape("ellipse", styling.Inches(5.5), styling.Inches(2), styling.Inches(3), styling.Inches(1)).
			WithFill(shapes.NewShapeFill("C62828")).
			WithText("Click AND Hover").
			WithClickAction(action.NewHyperlink(action.HyperlinkURL("https://example.com"))).
			WithHoverAction(action.NewHyperlink(action.HyperlinkLastSlide()).
				WithTooltip("Hover goes to last slide")))
}

func buildTextRunHyperlinkSlide() pptx.SlideContent {
	return pptx.NewSlide("Text Run Hyperlinks").
		AddBulletRuns([]text.TextRun{
			text.NewTextRun("Click "),
			text.NewTextRun("this link").
				WithHyperlink(action.NewHyperlink(action.HyperlinkURL("https://example.com")).
					WithTooltip("External link")).
				WithColor("1565C0").
				WithUnderline(true),
			text.NewTextRun(" to visit example.com"),
		}).
		AddBulletRuns([]text.TextRun{
			text.NewTextRun("Jump to "),
			text.NewTextRun("slide 1").
				WithHyperlink(action.NewHyperlink(action.HyperlinkSlide(1))).
				WithColor("2E7D32").
				WithUnderline(true),
			text.NewTextRun(" (internal navigation)"),
		}).
		AddBulletRuns([]text.TextRun{
			text.NewTextRun("Hover over "),
			text.NewTextRun("this text").
				WithHoverAction(action.NewHyperlink(action.HyperlinkNextSlide()).
					WithTooltip("Hover tooltip on text")).
				WithColor("E65100").
				WithBold(true),
			text.NewTextRun(" to see a tooltip"),
		})
}
