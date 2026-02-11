package pptx_test

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func TestTask10Showcase(t *testing.T) {
	slide := pptx.NewSlide("Images").
		AddBullet("Image embedding supported").
		AddBullet("PNG, JPEG, GIF formats").
		AddBullet("Position and size control")

	slide = slide.AddImage(pptx.NewImageFromBytes([]byte("fake png content"), "png", 4*914400, 1*914400, 2*914400, 2*914400))

	slides := []pptx.SlideContent{slide}
	_, err := pptx.CreateWithSlides("Task 10 Showcase", slides)
	if err != nil {
		t.Fatalf("failed to create showcase presentation: %v", err)
	}
}
