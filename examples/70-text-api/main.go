// examples/70-text-api demonstrates text formatting with paragraphs and runs.
//
// Shows Run and ParagraphStyle, rich text runs with bold/italic/underline/color/size,
// bullet styles (bullet, number, letter, roman), sub-bullets, and AddBulletRuns
// for per-run styling in a single bullet.
//
// Run with: go run ./examples/70-text-api/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

const (
	outputDir  = "examples/output"
	outputFile = "70_text_api.pptx"
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
		buildTitleStylingSlide(),
		buildBulletStylesSlide(),
		buildPerBulletStylesSlide(),
		buildRichTextRunsSlide(),
		buildRichNotesSlide(),
		buildTextSizePresetsSlide(),
	}

	outputPath := filepath.Join(outputDir, outputFile)
	data, err := pptx.CreateWithSlides("Text API Demo", slides)
	if err != nil {
		return fmt.Errorf("create presentation: %w", err)
	}
	if err := os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}

func buildTitleStylingSlide() pptx.SlideContent {
	slide1 := pptx.NewSlide("Slide Title Styling")
	slide1.TitleSize = 40
	slide1.TitleColor = "1565C0"
	slide1.TitleBold = true
	slide1.TitleItalic = false
	slide1.ContentSize = 16
	slide1.ContentColor = "333333"
	return slide1.
		AddBullet("Regular bullet – uses ContentSize and ContentColor").
		AddBullet("Second bullet – same styling")
}

func buildBulletStylesSlide() pptx.SlideContent {
	return pptx.NewSlide("Bullet Style Variants").
		AddBullet("Standard bullet point (bullet style)").
		AddNumbered("First numbered item").
		AddNumbered("Second numbered item").
		AddLettered("First lettered item (lowercase a, b, c)").
		AddLettered("Second lettered item").
		AddSubBullet(1, "Sub-bullet at indent level 1").
		AddSubBullet(2, "Sub-bullet at indent level 2").
		AddSubBullet(3, "Sub-bullet at indent level 3")
}

func buildPerBulletStylesSlide() pptx.SlideContent {
	boldStyle := elements.DefaultParagraphStyle()
	boldStyle.BulletColor = "C0504D"

	centeredStyle := elements.DefaultParagraphStyle()
	centeredStyle.Align = text.TextAlignCenter

	numberedStyle := elements.DefaultParagraphStyle().WithNumbered()

	return pptx.NewSlide("Per-Bullet Paragraph Styles").
		AddBulletWithStyle("Red bullet color (bullet color set to C0504D)", boldStyle).
		AddBulletWithStyle("Centered text bullet at 20pt", centeredStyle).
		AddBulletWithStyle("Numbered bullet via style", numberedStyle)
}

func buildRichTextRunsSlide() pptx.SlideContent {
	run1 := elements.NewRun("Bold: ")
	run1.Bold = true

	run2 := elements.NewRun("normal text, ")

	run3 := elements.NewRun("italic + blue")
	run3.Italic = true
	run3.Color = "1565C0"

	run4 := elements.NewRun(" and ")

	run5 := elements.NewRun("underlined red 18pt")
	run5.Underline = text.UnderlineStyleSingle
	run5.Color = "C0504D"
	run5.SizePt = 18

	strikeRun := elements.NewRun("Strikethrough text")
	strikeRun.Strikethrough = text.StrikethroughStyleSingle

	courierRun := elements.NewRun("Courier 14pt")
	courierRun.Font = "Courier New"
	courierRun.SizePt = 14

	return pptx.NewSlide("Rich Text Runs").
		AddBulletRuns([]elements.Run{run1, run2, run3, run4, run5}).
		AddBulletRuns([]elements.Run{
			strikeRun,
			elements.NewRun(" – "),
			courierRun,
		})
}

func buildRichNotesSlide() pptx.SlideContent {
	bulletRun := elements.NewRun("Notes paragraph with bullet")
	p1 := elements.NewParagraph()
	p1.Style.BulletStyle = text.BulletStyleBullet
	p1.Runs = []elements.Run{bulletRun}

	numberedRun := elements.NewRun("Numbered notes paragraph")
	p2 := elements.NewParagraph()
	p2.Style.BulletStyle = text.BulletStyleNumber
	p2.Runs = []elements.Run{numberedRun}

	return pptx.NewSlide("Rich Speaker Notes").
		AddBullet("This slide has rich speaker notes – check the notes panel.").
		WithRichNotes([]elements.Paragraph{p1, p2})
}

func buildTextSizePresetsSlide() pptx.SlideContent {
	slide6 := pptx.NewSlide("Text Size Presets")
	slide6.ContentSize = text.TextSizeBody
	return slide6.
		AddBullet(fmt.Sprintf("TextSizeTitle  = %d pt", text.TextSizeTitle)).
		AddBullet(fmt.Sprintf("TextSizeSubtitle = %d pt", text.TextSizeSubtitle)).
		AddBullet(fmt.Sprintf("TextSizeHeading  = %d pt", text.TextSizeHeading)).
		AddBullet(fmt.Sprintf("TextSizeBody     = %d pt", text.TextSizeBody)).
		AddBullet(fmt.Sprintf("TextSizeSmall    = %d pt", text.TextSizeSmall)).
		AddBullet(fmt.Sprintf("TextSizeCaption  = %d pt", text.TextSizeCaption)).
		AddBullet(fmt.Sprintf("TextSizeLarge    = %d pt", text.TextSizeLarge)).
		AddBullet(fmt.Sprintf("TextSizeXLarge   = %d pt", text.TextSizeXLarge))
}
