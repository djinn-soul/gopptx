package pptx

import "fmt"

const (
	// MinFontSizePt is the minimum allowed font size in points.
	MinFontSizePt = 1
	// MaxFontSizePt is the maximum allowed font size in points.
	MaxFontSizePt = 400
)

func validateSlideStyle(s SlideContent, slideIndex int) error {
	if s.TitleSize != 0 && (s.TitleSize < MinFontSizePt || s.TitleSize > MaxFontSizePt) {
		return fmt.Errorf("slide %d title size must be between %d and %d pt (got %d)", slideIndex, MinFontSizePt, MaxFontSizePt, s.TitleSize)
	}
	if s.TitleColor != "" && !isHexColor(s.TitleColor) {
		return fmt.Errorf("slide %d title color must be 6-digit RGB hex", slideIndex)
	}
	if s.ContentSize != 0 && (s.ContentSize < MinFontSizePt || s.ContentSize > MaxFontSizePt) {
		return fmt.Errorf("slide %d content size must be between %d and %d pt (got %d)", slideIndex, MinFontSizePt, MaxFontSizePt, s.ContentSize)
	}
	if s.ContentColor != "" && !isHexColor(s.ContentColor) {
		return fmt.Errorf("slide %d content color must be 6-digit RGB hex", slideIndex)
	}
	return nil
}
