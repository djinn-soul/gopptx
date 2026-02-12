package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

type (
	// SlideContent describes the user-visible content of a slide.
	SlideContent = elements.SlideContent
	// PlaceholderContent describes overridden content for a slide layout placeholder.
	PlaceholderContent = shapes.PlaceholderContent

	// SlideBackground defines how a slide's background is rendered.
	SlideBackground = elements.SlideBackground
	// SlideBackgroundType defines the filling method for a slide background.
	SlideBackgroundType = elements.SlideBackgroundType
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

	SlideBackgroundSolid    = elements.SlideBackgroundSolid
	SlideBackgroundGradient = elements.SlideBackgroundGradient
	SlideBackgroundPicture  = elements.SlideBackgroundPicture
)

func NewSolidBackground(color string) SlideBackground {
	return elements.NewSolidBackground(color)
}

func NewGradientBackground(gradient shapes.ShapeGradientFill) SlideBackground {
	return elements.NewGradientBackground(gradient)
}

func NewPictureBackground(img shapes.Image) SlideBackground {
	return elements.NewPictureBackground(img)
}
