package media_test

import (
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func TestSampleImageIntegration(t *testing.T) {
	imagePath := filepath.Join("..", "..", "..", "examples", "assets", "55", "repository-open-graph-template.png")

	// Use the image in a normal way
	slide1 := pptx.NewSlide("Sample Image Showcase").
		AddImage(pptx.NewImage(imagePath, 1*914400, 1*914400, 4*914400, 2*914400).
			WithRotation(10).
			WithFlip(true, false))

	// Use it as a placeholder override
	slide2 := pptx.NewSlide("Sample Image Placeholder").
		WithPlaceholderImageAs(1, "picture", pptx.NewImage(imagePath, 0, 0, 0, 0))

	slides := []pptx.SlideContent{slide1, slide2}
	_, err := pptx.CreateWithSlides("Sample Integration", slides)
	if err != nil {
		t.Fatalf("failed to create presentation with sample image: %v", err)
	}
}
