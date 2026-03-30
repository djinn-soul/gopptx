// examples/64-slides-api demonstrates the slides API:
// creating slides with different layouts, setting titles and content,
// slide numbering, footers, hidden slides, and adding slides via the editor.
//
// Run with: go run ./examples/64-slides-api/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	log "github.com/djinn-soul/gopptx/pkg/stdlog"
)

const (
	outputDir  = "examples/output"
	outputFile = "64_slides_api.pptx"
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

	// --- Slide 1: TitleAndContent (default layout) ---
	slide1 := pptx.NewSlide("Slide Layout: TitleAndContent").
		AddBullet("This is the default layout: title_and_content").
		AddBullet("All standard bullet slides use this layout")
	slide1.Layout = pptx.SlideLayoutTitleAndContent

	// --- Slide 2: TitleOnly ---
	slide2 := pptx.NewSlide("Slide Layout: TitleOnly")
	slide2.Layout = pptx.SlideLayoutTitleOnly

	// --- Slide 3: Blank ---
	slide3 := pptx.NewSlide("")
	slide3.Layout = pptx.SlideLayoutBlank

	// --- Slide 4: CenteredTitle ---
	slide4 := pptx.NewSlide("Centered Title Slide")
	slide4.Layout = pptx.SlideLayoutCenteredTitle

	// --- Slide 5: TwoColumn ---
	slide5 := pptx.NewSlide("Two Column Layout").
		AddBullet("Left column content goes here").
		AddBullet("Second bullet in left column")
	slide5.Layout = pptx.SlideLayoutTwoColumn

	// --- Slide 6: TitleAndBigContent ---
	slide6 := pptx.NewSlide("Big Content Layout").
		AddBullet("Uses a larger content area").
		AddBullet("Ideal for a single key point")
	slide6.Layout = pptx.SlideLayoutTitleAndBigContent

	// --- Slide 7: Show slide number and footer ---
	slide7 := pptx.NewSlide("Slide Number & Footer Demo").
		AddBullet("Slide numbers are enabled for this slide").
		AddBullet("Footer text is shown below").
		WithSlideNumber(true)
	slide7.FooterText = "© 2025 gopptx Examples"
	slide7.ShowSlideNumber = true

	// --- Slide 8: Title size and color customization ---
	slide8 := pptx.NewSlide("Styled Title Slide")
	slide8.TitleSize = 36
	slide8.TitleColor = "1565C0"
	slide8.TitleBold = true
	slide8.AddBullet("Title is styled with a custom size and color").
		AddBullet("TitleSize=36pt, TitleColor=1565C0, TitleBold=true")
	slide8 = slide8.AddBullet("These properties control the slide title appearance")

	// --- Slide 9: Hidden slide ---
	slide9 := pptx.NewSlide("Hidden Slide (will be skipped in show)")
	slide9.Hidden = true
	slide9 = slide9.AddBullet("This slide is marked hidden")

	// --- Slide 10: Numbered and lettered bullets ---
	slide10 := pptx.NewSlide("Bullet Variants").
		AddNumbered("First numbered item").
		AddNumbered("Second numbered item").
		AddLettered("First lettered item").
		AddLettered("Second lettered item").
		AddSubBullet(1, "Sub-bullet at level 1").
		AddSubBullet(2, "Sub-bullet at level 2")

	slides := []pptx.SlideContent{
		slide1, slide2, slide3, slide4, slide5,
		slide6, slide7, slide8, slide9, slide10,
	}

	outputPath := filepath.Join(outputDir, outputFile)
	data, err := pptx.CreateWithSlides("Slides API Demo", slides)
	if err != nil {
		return fmt.Errorf("create presentation: %w", err)
	}
	if err := os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	// --- Demonstrate adding a slide via PresentationEditor ---
	ed, err := pptx.OpenPresentationEditor(outputPath)
	if err != nil {
		return fmt.Errorf("open editor: %w", err)
	}
	defer ed.Close()

	log.Printf("Slide count before AddSlide: %d\n", ed.SlideCount())

	_, err = ed.AddSlide(
		pptx.NewSlide("Added via Editor").
			AddBullet("This slide was appended via PresentationEditor.AddSlide"),
	)
	if err != nil {
		return fmt.Errorf("add slide via editor: %w", err)
	}

	log.Printf("Slide count after AddSlide: %d\n", ed.SlideCount())

	if err := ed.Save(outputPath); err != nil {
		return fmt.Errorf("save after editor: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}
