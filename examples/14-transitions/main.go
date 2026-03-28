// examples/14-transitions/main.go demonstrates slide transitions in gopptx.
//
// Shows Fade, Push, Wipe, Split, and Zoom transitions applied to individual slides
// using both the simple WithTransition helper and the full TransitionOptions API.
//
// Run with: go run ./examples/14-transitions/main.go
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
	outputFile = "14_transitions.pptx"
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

	fade := pptx.TransitionOptions{
		Type:       pptx.TransitionFade,
		DurationMS: 500,
	}

	push := pptx.TransitionOptions{
		Type:       pptx.TransitionPush,
		Direction:  pptx.TransitionDirLeft,
		DurationMS: 700,
	}

	wipe := pptx.TransitionOptions{
		Type:       pptx.TransitionWipe,
		Direction:  pptx.TransitionDirRight,
		DurationMS: 800,
	}

	zoom := pptx.TransitionOptions{
		Type:       pptx.TransitionZoom,
		DurationMS: 600,
	}

	slides := []pptx.SlideContent{
		// Slide 1: no transition — overview
		pptx.NewSlide("Slide Transitions Demo").
			AddBullet("Each following slide demonstrates a different transition.").
			AddBullet("Open in PowerPoint and advance slides to see the effects."),

		// Slide 2: Fade
		pptx.NewSlide("Fade Transition").
			WithTransitionOptions(fade).
			AddBullet("Type: TransitionFade").
			AddBullet("Duration: 500 ms").
			AddBullet("Fades the current slide into the next."),

		// Slide 3: Push Left
		pptx.NewSlide("Push Left Transition").
			WithTransitionOptions(push).
			AddBullet("Type: TransitionPush, Direction: TransitionDirLeft").
			AddBullet("Duration: 700 ms").
			AddBullet("The new slide pushes in from the right side."),

		// Slide 4: Wipe Right
		pptx.NewSlide("Wipe Right Transition").
			WithTransitionOptions(wipe).
			AddBullet("Type: TransitionWipe, Direction: TransitionDirRight").
			AddBullet("Duration: 800 ms").
			AddBullet("The new slide wipes over from the left."),

		// Slide 5: Zoom
		pptx.NewSlide("Zoom Transition").
			WithTransitionOptions(zoom).
			AddBullet("Type: TransitionZoom").
			AddBullet("Duration: 600 ms").
			AddBullet("The new slide zooms in from the center."),

		// Slide 6: Split (simple helper form)
		pptx.NewSlide("Split Transition").
			WithTransition(pptx.TransitionSplit).
			AddBullet("Applied with the shorthand WithTransition(pptx.TransitionSplit).").
			AddBullet("Default duration is used."),
	}

	outputPath := filepath.Join(outputDir, outputFile)
	if err := pptx.WriteFile(outputPath, "Slide Transitions Demo", slides); err != nil {
		return fmt.Errorf("write presentation: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}
