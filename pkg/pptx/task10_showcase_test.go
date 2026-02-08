package pptx

import (
	"testing"
)

func TestTask10Showcase(t *testing.T) {
	// Replicate the slide from the user's screenshot
	slide := NewSlide("Images").
		AddBullet("Image embedding supported").
		AddBullet("PNG, JPEG, GIF formats").
		AddBullet("Position and size control")

	// Add an image to demonstrate embedding and control
	// X: 4 inches, Y: 1 inch, CX: 2 inches, CY: 2 inches (in EMUs)
	// 1 inch = 914400 EMUs
	slide = slide.AddImage(NewImageFromBytes([]byte("fake png content"), "png", 4*914400, 1*914400, 2*914400, 2*914400))

	slides := []SlideContent{slide}
	_, err := CreateWithSlides("Task 10 Showcase", slides)
	if err != nil {
		t.Fatalf("failed to create showcase presentation: %v", err)
	}
}
