package pptx

import (
	"path/filepath"
	"testing"
)

func TestSampleImageIntegration(t *testing.T) {
	imagePath := filepath.Join("..", "..", "smoke_samples", "sampleimage", "repository-open-graph-template.png")

	// Use the image in a normal way
	slide1 := NewSlide("Sample Image Showcase").
		AddImage(NewImage(imagePath, 1*914400, 1*914400, 4*914400, 2*914400).
			WithRotation(10).
			WithFlip(true, false))

	// Use it as a placeholder override
	slide2 := NewSlide("Sample Image Placeholder").
		WithPlaceholderImage(1, "picture", NewImage(imagePath, 0, 0, 0, 0))

	slides := []SlideContent{slide1, slide2}
	_, err := CreateWithSlides("Sample Integration", slides)
	if err != nil {
		t.Fatalf("failed to create presentation with sample image: %v", err)
	}
}
