// examples/71-action-api demonstrates hyperlinks and click actions on shapes.
//
// Shows URL links, slide navigation (first/last/next/prev/specific slide),
// email links, hover actions, and how to attach hyperlinks to text runs.
//
// Run with: go run ./examples/71-action-api/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

const (
	outputDir  = "examples/output"
	outputFile = "71_action_api.pptx"
)

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
		buildURLHyperlinkSlide(),
		buildNavigationSlide(),
		buildJumpToSlideSlide(),
		buildEmailHyperlinkSlide(),
		buildHoverActionSlide(),
		buildTextRunHyperlinkSlide(),
	}

	outputPath := filepath.Join(outputDir, outputFile)
	data, err := pptx.CreateWithSlides("Action API Demo", slides)
	if err != nil {
		return fmt.Errorf("create presentation: %w", err)
	}
	if err := os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}

func buildURLHyperlinkSlide() pptx.SlideContent {
	urlLink := action.NewHyperlink(
		action.HyperlinkURL("https://github.com/djinn-soul/gopptx"),
	).
		WithTooltip("Visit gopptx on GitHub").
		WithHighlightClick(true)

	urlShape := pptx.NewRectangle(1, 1.5, 5, 1).
		WithFill(pptx.NewShapeFill("4472C4")).
		WithText("Click to open GitHub (URL hyperlink)").
		WithClickAction(urlLink)

	return pptx.NewSlide("URL Hyperlink on Shape").
		AddShape(urlShape).
		AddBullet("Shapes can carry click actions via WithClickAction.")
}

func buildNavigationSlide() pptx.SlideContent {
	nextBtn := pptx.NewRectangle(0.5, 4.5, 2, 0.8).
		WithFill(pptx.NewShapeFill("9BBB59")).
		WithText("Next Slide").
		WithClickAction(action.NewHyperlink(action.HyperlinkNextSlide()))

	prevBtn := pptx.NewRectangle(3, 4.5, 2, 0.8).
		WithFill(pptx.NewShapeFill("C0504D")).
		WithText("Prev Slide").
		WithClickAction(action.NewHyperlink(action.HyperlinkPreviousSlide()))

	firstBtn := pptx.NewRectangle(5.5, 4.5, 2, 0.8).
		WithFill(pptx.NewShapeFill("F79646")).
		WithText("First Slide").
		WithClickAction(action.NewHyperlink(action.HyperlinkFirstSlide()))

	lastBtn := pptx.NewRectangle(0.5, 5.5, 2, 0.8).
		WithFill(pptx.NewShapeFill("8064A2")).
		WithText("Last Slide").
		WithClickAction(action.NewHyperlink(action.HyperlinkLastSlide()))

	endShowBtn := pptx.NewRectangle(3, 5.5, 2, 0.8).
		WithFill(pptx.NewShapeFill("333333")).
		WithText("End Show").
		WithClickAction(action.NewHyperlink(action.HyperlinkEndShow()))

	return pptx.NewSlide("Slide Navigation Actions").
		AddShape(nextBtn).
		AddShape(prevBtn).
		AddShape(firstBtn).
		AddShape(lastBtn).
		AddShape(endShowBtn).
		AddBullet("Use HyperlinkNextSlide, PreviousSlide, FirstSlide, LastSlide, EndShow.")
}

func buildJumpToSlideSlide() pptx.SlideContent {
	jumpBtn := pptx.NewRectangle(1, 2, 4, 1).
		WithFill(pptx.NewShapeFill("4BACC6")).
		WithText("Jump to Slide 5").
		WithClickAction(action.NewHyperlink(action.HyperlinkSlide(5)))

	return pptx.NewSlide("Jump to Specific Slide").
		AddShape(jumpBtn).
		AddBullet("HyperlinkSlide(n) jumps to slide number n in the show.")
}

func buildEmailHyperlinkSlide() pptx.SlideContent {
	emailLink := action.NewHyperlink(
		action.HyperlinkEmailWithSubject("hello@example.com", "Question about gopptx"),
	).WithTooltip("Send us an email")

	emailShape := pptx.NewRectangle(1, 2, 6, 1).
		WithFill(pptx.NewShapeFill("EBF1DE")).
		WithText("Click to email hello@example.com").
		WithClickAction(emailLink)

	return pptx.NewSlide("Email Hyperlink").
		AddShape(emailShape).
		AddBullet("HyperlinkEmailWithSubject sets the mailto: href with a subject.")
}

func buildHoverActionSlide() pptx.SlideContent {
	hoverShape := pptx.NewRectangle(1, 2, 6, 1).
		WithFill(pptx.NewShapeFill("FDE9D9")).
		WithText("Hover over me!").
		WithHoverAction(action.NewHyperlink(action.HyperlinkURL("https://example.com")))

	return pptx.NewSlide("Hover Action").
		AddShape(hoverShape).
		AddBullet("Shapes can also have hover actions via WithHoverAction.")
}

func buildTextRunHyperlinkSlide() pptx.SlideContent {
	linkRun := elements.NewRun("click here").WithHyperlink(
		action.NewHyperlink(action.HyperlinkURL("https://pkg.go.dev")),
	)
	plainRun1 := elements.NewRun("Visit Go docs – ")
	plainRun2 := elements.NewRun(" – for package documentation.")

	return pptx.NewSlide("Hyperlink on Text Run").
		AddBulletRuns([]elements.Run{plainRun1, linkRun, plainRun2}).
		AddBullet("Text runs support WithHyperlink for inline links.")
}
