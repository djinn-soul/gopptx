package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

const (
	outputDir  = "examples/output"
	outputFile = "53_slide_properties.pptx"
)

func main() {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		fail("create output directory", err)
	}

	slides := []pptx.SlideContent{
		pptx.NewSlide("Slide Background Color").
			WithBackgroundColor("D9E1F2"). // Light blue
			AddBullet("This slide has a solid background color (D9E1F2).").
			WithSlideNumber(true),

		pptx.NewSlide("Right Aligned Title").
			WithTitleAlign("r").
			AddBullet("The title of this slide is right-aligned.").
			WithSlideNumber(true),

		pptx.NewSlide("Centered Title (justified)").
			WithTitleAlign("just").
			AddBullet("The title of this slide is justified (though hard to see with short text).").
			WithSlideNumber(true),

		pptx.NewSlide("Bottom Aligned Content").
			WithContentVAlign("b").
			AddBullet("The text content of this slide is vertically aligned to the bottom.").
			AddBullet("This is useful for specific layout requirements.").
			WithSlideNumber(true),

		pptx.NewSlide("Middle Aligned Content").
			WithContentVAlign("ctr").
			AddBullet("The text content of this slide is vertically aligned to the middle.").
			WithSlideNumber(true),

		pptx.NewSlide("Consolas Title Font").
			WithTitleFont("Consolas").
			WithTitleSize(32).
			AddBullet("The title uses the 'Consolas' typeface.").
			WithSlideNumber(true),

		pptx.NewSlide("Everything Combined").
			WithBackgroundColor("FCE4D6"). // Light orange
			WithTitleFont("Impact").
			WithTitleAlign("ctr").
			WithTitleColor("AA0000").
			WithContentVAlign("ctr").
			AddBullet("Middle aligned content").
			AddBullet("Centered Impact title").
			AddBullet("Custom background color").
			WithSlideNumber(true),
	}

	data, buildErr := pptx.CreateWithSlides("gopptx Slide Properties Smoke", slides)
	if buildErr != nil {
		fail("create presentation", buildErr)
	}

	path := filepath.Join(outputDir, outputFile)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		fail("write output", err)
	}

	log.Printf("Successfully generated smoke sample: %s\n", path)
}

func fail(step string, err error) {
	fmt.Fprintf(os.Stderr, "error: %s: %v\n", step, err)
	os.Exit(1)
}
