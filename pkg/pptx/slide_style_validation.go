package pptx

import "fmt"

func validateSlideStyle(s SlideContent, slideIndex int) error {
	if s.TitleSize < 0 || s.TitleSize > 400 {
		return fmt.Errorf("slide %d title size must be between 1 and 400 pt (got %d)", slideIndex, s.TitleSize)
	}
	if s.TitleColor != "" && !isHexColor(s.TitleColor) {
		return fmt.Errorf("slide %d title color must be 6-digit RGB hex", slideIndex)
	}
	if s.ContentSize < 0 || s.ContentSize > 400 {
		return fmt.Errorf("slide %d content size must be between 1 and 400 pt (got %d)", slideIndex, s.ContentSize)
	}
	if s.ContentColor != "" && !isHexColor(s.ContentColor) {
		return fmt.Errorf("slide %d content color must be 6-digit RGB hex", slideIndex)
	}
	return nil
}
