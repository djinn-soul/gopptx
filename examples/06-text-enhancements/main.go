// examples/06-text-enhancements/main.go demonstrates text run enhancements.
//
// Shows strikethrough, highlight, all-caps, small-caps, subscript, superscript,
// bold, italic, and underline applied to individual text runs within a bullet.
//
// Run with: go run ./examples/06-text-enhancements/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

const (
	outputDir  = "examples/output"
	outputFile = "06_text_enhancements.pptx"
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

	// Slide 1: Strikethrough and highlight
	slide1 := pptx.NewSlide("Strikethrough and Highlight").
		AddBulletRuns([]pptx.Run{
			pptx.NewRun("Normal text, then "),
			pptx.NewRun("strikethrough text").WithStrikethrough(true),
			pptx.NewRun(" at the end."),
		}).
		AddBulletRuns([]pptx.Run{
			pptx.NewRun("This word is "),
			pptx.NewRun("highlighted yellow").WithHighlight("FFFF00"),
			pptx.NewRun(" for emphasis."),
		})

	// Slide 2: All-caps and small-caps
	slide2 := pptx.NewSlide("Capitalization Styles").
		AddBulletRuns([]pptx.Run{
			pptx.NewRun("Normal, then "),
			pptx.NewRun("all caps mode").WithAllCaps(true),
			pptx.NewRun(" applied here."),
		}).
		AddBulletRuns([]pptx.Run{
			pptx.NewRun("Normal, then "),
			pptx.NewRun("small caps mode").WithSmallCaps(true),
			pptx.NewRun(" applied here."),
		})

	// Slide 3: Subscript and superscript
	slide3 := pptx.NewSlide("Subscript and Superscript").
		AddBulletRuns([]pptx.Run{
			pptx.NewRun("H"),
			pptx.NewRun("2").WithSubscript(true),
			pptx.NewRun("O is the chemical formula for water."),
		}).
		AddBulletRuns([]pptx.Run{
			pptx.NewRun("E = mc"),
			pptx.NewRun("2").WithSuperscript(true),
			pptx.NewRun(" is Einstein's mass-energy equation."),
		})

	// Slide 4: Bold, italic, and underline
	slide4 := pptx.NewSlide("Bold, Italic, and Underline").
		AddBulletRuns([]pptx.Run{
			pptx.NewRun("This run is "),
			pptx.NewRun("bold").WithBold(true),
			pptx.NewRun(" for emphasis."),
		}).
		AddBulletRuns([]pptx.Run{
			pptx.NewRun("This run is "),
			pptx.NewRun("italic").WithItalic(true),
			pptx.NewRun(" for style."),
		}).
		AddBulletRuns([]pptx.Run{
			pptx.NewRun("This run is "),
			pptx.NewRun("underlined").WithUnderline(true),
			pptx.NewRun(" for attention."),
		})

	// Slide 5: Combined enhancements in a single run
	slide5 := pptx.NewSlide("Combined Enhancements").
		AddBulletRuns([]pptx.Run{
			pptx.NewRun("Plain, "),
			pptx.NewRun("bold-italic-underline").
				WithBold(true).
				WithItalic(true).
				WithUnderline(true),
			pptx.NewRun(", plain again."),
		}).
		AddBulletRuns([]pptx.Run{
			pptx.NewRun("Highlighted "),
			pptx.NewRun("and bold").
				WithHighlight("FFFF00").
				WithBold(true),
			pptx.NewRun(" together."),
		})

	slides := []pptx.SlideContent{slide1, slide2, slide3, slide4, slide5}

	data, err := pptx.CreateWithSlides("Text Enhancements Demo", slides)
	if err != nil {
		return fmt.Errorf("create presentation: %w", err)
	}

	outputPath := filepath.Join(outputDir, outputFile)
	if err = os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}
