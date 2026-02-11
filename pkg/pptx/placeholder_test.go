package pptx

import (
	"testing"
)

func TestPlaceholderOverrides(t *testing.T) {
	slides := []SlideContent{
		NewSlide("Slide 1").
			WithPlaceholderText(0, "title", "Title Override").
			WithPlaceholderText(1, "body", "Body Override"),
		NewSlide("Slide 2").
			WithPlaceholderImage(1, "picture", NewImageFromBytes([]byte("fake png"), "png", 0, 0, 0, 0)),
		NewSlide("Slide 3").
			WithPlaceholderTable(1, "body", Table{
				Rows: [][]string{
					{"A", "B"},
					{"1", "2"},
				},
			}),
	}

	_, err := CreateWithSlides("Test Pres", slides)
	if err != nil {
		t.Fatalf("failed to create presentation: %v", err)
	}
}
