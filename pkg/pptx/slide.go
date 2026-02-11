package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

type (
	// SlideContent describes the user-visible content of a slide.
	SlideContent = elements.SlideContent
	// PlaceholderContent describes overridden content for a slide layout placeholder.
	PlaceholderContent = elements.PlaceholderContent
)

func NewSlide(title string) SlideContent {
	return elements.NewSlide(title)
}

const (
	SlideLayoutTitleAndContent    = elements.SlideLayoutTitleAndContent
	SlideLayoutTitleOnly          = elements.SlideLayoutTitleOnly
	SlideLayoutBlank              = elements.SlideLayoutBlank
	SlideLayoutCenteredTitle      = elements.SlideLayoutCenteredTitle
	SlideLayoutTitleAndBigContent = elements.SlideLayoutTitleAndBigContent
	SlideLayoutTwoColumn          = elements.SlideLayoutTwoColumn
)
