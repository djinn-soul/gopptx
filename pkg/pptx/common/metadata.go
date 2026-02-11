package common

// SlideSize describes the dimensions of slides in a presentation in EMUs.
type SlideSize struct {
	Width  int64
	Height int64
}

var (
	// SlideSize4x3 is the standard 4:3 slide size (10x7.5 inches).
	SlideSize4x3 = SlideSize{Width: 9144000, Height: 6858000}
	// SlideSize16x9 is the standard 16:9 widescreen slide size (13.33x7.5 inches).
	SlideSize16x9 = SlideSize{Width: 12192000, Height: 6858000}
)

// PresentationMetadata describes summary information for a PPTX package.
type PresentationMetadata struct {
	Title       string
	Subject     string
	Creator     string
	Description string
	SlideSize   SlideSize
	SlideCount  int
}
