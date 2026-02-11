package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/layout"
)

const (
	// SlideWidth is the standard width of a 4:3 slide in EMU.
	SlideWidth = layout.SlideWidth
	// SlideHeight is the standard height of a 4:3 slide in EMU.
	SlideHeight = layout.SlideHeight

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

var (
	// Center calculates the (X, Y) coordinates to center an element of size (cx, cy)
	// within the standard slide bounds.
	Center = layout.Center
	// CenterInBox calculates the (X, Y) coordinates to center an element of size (cx, cy)
	// within a specific bounding box.
	CenterInBox = layout.CenterInBox
	// Grid calculates a sequence of bounding boxes for a grid layout.
	Grid = layout.Grid
	// GridInBox calculates a sequence of bounding boxes for a grid layout within specific bounds.
	GridInBox = layout.GridInBox
	// Stack calculates the starting points for elements placed sequentially with a fixed gap.
	Stack = layout.Stack
	// Distribute calculates the top or left coordinates to evenly space elements within a bound.
	Distribute = layout.Distribute
)
