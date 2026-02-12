package common

import "github.com/djinn-soul/gopptx/pkg/pptx/styling"

// Point represents a 2D coordinate in EMU.
type Point struct {
	X styling.Length
	Y styling.Length
}

// Size represents dimensions in EMU.
type Size struct {
	CX styling.Length
	CY styling.Length
}

// Box represents a rectangular region with position and size in EMU.
type Box struct {
	X  styling.Length
	Y  styling.Length
	CX styling.Length
	CY styling.Length
}
