package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/layout"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	// SlideWidth is the standard width of a 4:3 slide in EMU.
	SlideWidth Length = layout.SlideWidth
	// SlideHeight is the standard height of a 4:3 slide in EMU.
	SlideHeight Length = layout.SlideHeight

	// OrientationHorizontal represents left-to-right or horizontal distribution.
	OrientationHorizontal = layout.OrientationHorizontal
	// OrientationVertical represents top-to-bottom or vertical distribution.
	OrientationVertical = layout.OrientationVertical
)

type (
	// Point represents a 2D coordinate in EMU.
	Point = common.Point
	// Size represents dimensions in EMU.
	Size = common.Size
	// Box represents a rectangular region with position and size in EMU.
	Box = common.Box
)

// Center calculates the (X, Y) coordinates to center an element of size (cx, cy)
// within the standard slide bounds.
func Center(cx, cy styling.Length) (styling.Length, styling.Length) {
	return layout.Center(cx, cy)
}

// CenterInSize calculates the (X, Y) coordinates to center an element of size (cx, cy)
// within arbitrary dimensions.
func CenterInSize(cx, cy, totalW, totalH styling.Length) (styling.Length, styling.Length) {
	return layout.CenterInSize(cx, cy, totalW, totalH)
}

// CenterInBox calculates the (X, Y) coordinates to center an element of size (cx, cy)
// within a specific bounding box.
func CenterInBox(cx, cy styling.Length, bounds common.Box) (styling.Length, styling.Length) {
	return layout.CenterInBox(cx, cy, bounds)
}

// Grid calculates a sequence of bounding boxes for a grid layout.
func Grid(rows, cols int, margin styling.Length) ([]common.Box, error) {
	return layout.Grid(rows, cols, margin)
}

// GridInBox calculates a sequence of bounding boxes for a grid layout within specific bounds.
func GridInBox(rows, cols int, margin styling.Length, bounds common.Box) ([]common.Box, error) {
	return layout.GridInBox(rows, cols, margin, bounds)
}

// Stack calculates the starting points for elements placed sequentially with a fixed gap.
func Stack(
	orientation string,
	start common.Point,
	gap styling.Length,
	elements ...common.Size,
) ([]common.Point, error) {
	return layout.Stack(orientation, start, gap, elements...)
}

// DistributeUniform calculates the top or left coordinates to evenly space elements of identical size within a bound.
func DistributeUniform(
	orientation string,
	bounds common.Box,
	count int,
	elementSize styling.Length,
) ([]styling.Length, error) {
	return layout.DistributeUniform(orientation, bounds, count, elementSize)
}

// DistributeNonUniform calculates the top or left coordinates to space variable-sized elements with uniform gaps within a bound.
func DistributeNonUniform(
	orientation string,
	bounds common.Box,
	sizes []styling.Length,
) ([]styling.Length, error) {
	return layout.DistributeNonUniform(orientation, bounds, sizes)
}
