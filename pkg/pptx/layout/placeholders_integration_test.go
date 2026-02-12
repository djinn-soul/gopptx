package layout_test

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func TestPlaceholderOverrides(t *testing.T) {
	slides := []pptx.SlideContent{
		pptx.NewSlide("Slide 1").
			WithPlaceholderTextAs(0, "title", "Title Override").
			WithPlaceholderTextAs(1, "body", "Body Override"),
		pptx.NewSlide("Slide 2").
			WithPlaceholderImageAs(1, "picture", pptx.NewImageFromBytes([]byte("fake png"), "png", 0, 0, 0, 0)),
		pptx.NewSlide("Slide 3").
			WithPlaceholderTableAs(1, "body", pptx.Table{
				Rows: [][]string{
					{"A", "B"},
					{"1", "2"},
				},
			}),
	}

	_, err := pptx.CreateWithSlides("Test Pres", slides)
	if err != nil {
		t.Fatalf("failed to create presentation: %v", err)
	}
}
