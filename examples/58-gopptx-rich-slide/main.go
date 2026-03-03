package main

import (
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/gopptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/animations"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	outputDir  = "examples/output"
	outputFile = "58_gopptx_rich_slide.pptx"
)

func main() {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	pres := &gopptx.Presentation{Title: "gopptx Rich Slide Demo"}

	slide := pres.AddSlide()
	slide.Title = "gopptx Rich Slide"
	slide.AddBullet("gopptx now exposes helpers for notes, placeholders, and animations.")
	slide.AddBullet("Reuse the same slide object to keep your API surface stable.")
	slide.SetRichNotes([]elements.Paragraph{
		elements.NewParagraph().AddRun(elements.NewRun("Use these speaker notes to describe how the slide should feel.")),
	})

	slide.AddPlaceholderText(0, "Title placeholder override")
	slide.AddPlaceholderOverride(shapes.PlaceholderContent{
		Index: 1,
		Type:  "body",
		Text:  "Body placeholder text with layout overrides and custom styling.",
		Override: &shapes.PlaceholderOverrideOptions{
			X:  ptrLength(styling.Inches(0.7)),
			Y:  ptrLength(styling.Inches(3.2)),
			CX: ptrLength(styling.Inches(6.8)),
			CY: ptrLength(styling.Inches(2.4)),
			TextStyle: &shapes.PlaceholderTextStyle{
				Font:   strPtr("Segoe UI"),
				SizePt: intPtr(28),
				Bold:   boolPtr(true),
				Color:  strPtr("1F4E79"),
			},
		},
	})

	shape := shapes.NewRectangle(1.2, 2.9, 3.2, 1.1).
		WithText("Animated callout").
		WithFill(shapes.NewShapeFill("F4C542"))
	slide.AddShape(shape)

	slide.AddAnimation(
		animations.NewAnimation(1, animations.AnimationEntranceFade).
			WithTrigger(animations.AnimationOnClick).
			WithDelay(250),
	)
	slide.AddAnimationDefinition(
		animations.NewAnimation(1, animations.AnimationEntranceZoom).
			WithTrigger(animations.AnimationAfterPrevious).
			WithDuration(900),
	)

	outputPath := filepath.Join(outputDir, outputFile)
	if err := pres.Save(outputPath); err != nil {
		log.Fatalf("failed to save rich slide: %v", err)
	}

	log.Printf("generated %s", outputPath)
}

func ptrLength(v styling.Length) *styling.Length { return &v }
func intPtr(v int) *int                          { return &v }
func boolPtr(v bool) *bool                       { return &v }
func strPtr(v string) *string                    { return &v }
