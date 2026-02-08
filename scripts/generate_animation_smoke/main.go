package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/djinn09/gopptx/pkg/pptx"
)

const (
	outputDir  = "smoke_samples"
	outputFile = "animations.pptx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	transPushRight := pptx.TransitionOptions{
		Type:       pptx.TransitionPush,
		Direction:  pptx.TransitionDirRight,
		DurationMS: 1000,
	}

	slides := []pptx.SlideContent{
		// Slide 1: Title
		pptx.NewSlide("Animations & Transitions Demo").
			AddBullet("This presentation demonstrates new features in gopptx.").
			AddBullet("Slide Transitions").
			AddBullet("Shape Animations"),

		// Slide 2: Transition Demo (Push Right)
		pptx.NewSlide("Transition: Push Right (1s)").
			WithTransitionOptions(transPushRight).
			AddBullet("This slide entered with a Push Right transition.").
			AddBullet("Duration: 1000ms"),

		// Slide 3: Animations Demo (No Transition to ensure animations are seen clearly)
		pptx.NewSlide("Animations Demo").
			AddShape(pptx.NewShape("rect", 1000000, 1500000, 2000000, 2000000).
				WithFill(pptx.NewShapeFill("FF0000")).
				WithText("Fade In")).
			AddShape(pptx.NewShape("ellipse", 4000000, 1500000, 2000000, 2000000).
				WithFill(pptx.NewShapeFill("00FF00")).
				WithText("Fly In")).
			AddShape(pptx.NewShape("triangle", 7000000, 1500000, 2000000, 2000000).
				WithFill(pptx.NewShapeFill("0000FF")).
				WithText("Pulse")).

			// Animations matching the shapes (1-based index based on AddShape calls)
			// Shape 1 (Rect): Fade In
			AddAnimation(pptx.NewAnimation(1, pptx.AnimationEntranceFade).
				WithTrigger(pptx.AnimationOnClick)).

			// Shape 2 (Ellipse): Fly In (After Previous)
			AddAnimation(pptx.NewAnimation(2, pptx.AnimationEntranceFlyIn).
				WithTrigger(pptx.AnimationAfterPrevious).
				WithDuration(1000)).

			// Shape 3 (Triangle): Pulse (After Previous)
			AddAnimation(pptx.NewAnimation(3, pptx.AnimationEmphasisPulse).
				WithTrigger(pptx.AnimationAfterPrevious)),

		// Slide 4: Conclusion
		pptx.NewSlide("End of Demo").
			WithTransition(pptx.TransitionSplit).
			AddBullet("Thank you!"),
	}

	path := filepath.Join(outputDir, outputFile)
	if err := pptx.WriteFile(path, "Animations Demo", slides); err != nil {
		return fmt.Errorf("failed to save presentation: %w", err)
	}

	fmt.Printf("Generated %s\n", path)
	return nil
}



