package pptx

import (
	"fmt"
)

// SlideContent describes the user-visible content of a slide.
type SlideContent struct {
	Title   string
	Bullets []string
}

// NewSlide creates a new slide with a title.
func NewSlide(title string) SlideContent {
	return SlideContent{Title: title}
}

// AddBullet appends one bullet item and returns the updated slide.
func (s SlideContent) AddBullet(text string) SlideContent {
	s.Bullets = append(s.Bullets, text)
	return s
}

func validateSlide(s SlideContent, index int) error {
	if s.Title == "" {
		return fmt.Errorf("slide %d title cannot be empty", index)
	}
	for bulletIndex, bullet := range s.Bullets {
		if bullet == "" {
			return fmt.Errorf("slide %d bullet %d cannot be empty", index, bulletIndex+1)
		}
	}
	return nil
}
