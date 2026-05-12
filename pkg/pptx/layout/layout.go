package layout

import (
	"errors"
	"fmt"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	// SlideWidth is the standard width of a 4:3 slide in EMU.
	SlideWidth styling.Length = 9144000
	// SlideHeight is the standard height of a 4:3 slide in EMU.
	SlideHeight styling.Length = 6858000
)

const (
	// OrientationHorizontal represents left-to-right or horizontal distribution.
	OrientationHorizontal = "horizontal"
	// OrientationVertical represents top-to-bottom or vertical distribution.
	OrientationVertical = "vertical"

	centerDivisor styling.Length = 2
)

// Center calculates the (X, Y) coordinates to center an element of size (cx, cy)
// within the standard 4:3 slide bounds.
func Center(cx, cy styling.Length) (styling.Length, styling.Length) {
	return CenterInSize(cx, cy, SlideWidth, SlideHeight)
}

// CenterInSize calculates the (X, Y) coordinates to center an element of size (cx, cy)
// within total dimensions (totalW, totalH).
func CenterInSize(cx, cy, totalW, totalH styling.Length) (styling.Length, styling.Length) {
	return (totalW - cx) / centerDivisor, (totalH - cy) / centerDivisor
}

// CenterInBox calculates the (X, Y) coordinates to center an element of size (cx, cy)
// within a specific bounding box.
func CenterInBox(cx, cy styling.Length, bounds common.Box) (styling.Length, styling.Length) {
	x := bounds.X + (bounds.CX-cx)/centerDivisor
	y := bounds.Y + (bounds.CY-cy)/centerDivisor
	return x, y
}

// Grid calculates a sequence of bounding boxes for a grid layout.
// rows and cols must be greater than zero. margin is the spacing between elements.
// The grid fills the standard slide area.
func Grid(rows, cols int, margin styling.Length) ([]common.Box, error) {
	return GridInBox(rows, cols, margin, common.Box{X: 0, Y: 0, CX: SlideWidth, CY: SlideHeight})
}

// GridInBox calculates a sequence of bounding boxes for a grid layout within specific bounds.
func GridInBox(rows, cols int, margin styling.Length, bounds common.Box) ([]common.Box, error) {
	if rows <= 0 || cols <= 0 {
		return nil, errors.New("rows and cols must be greater than zero")
	}

	totalMarginX := margin * styling.Length(cols-1)
	totalMarginY := margin * styling.Length(rows-1)

	if totalMarginX >= bounds.CX || totalMarginY >= bounds.CY {
		return nil, errors.New("margins exceed bounding box dimensions")
	}

	elementCX := (bounds.CX - totalMarginX) / styling.Length(cols)
	elementCY := (bounds.CY - totalMarginY) / styling.Length(rows)

	boxes := make([]common.Box, 0, rows*cols)
	for r := range rows {
		for c := range cols {
			x := bounds.X + styling.Length(c)*(elementCX+margin)
			y := bounds.Y + styling.Length(r)*(elementCY+margin)

			boxes = append(boxes, common.Box{
				X:  x,
				Y:  y,
				CX: elementCX,
				CY: elementCY,
			})
		}
	}

	return boxes, nil
}

// Stack calculates the starting points for elements placed sequentially with a fixed gap.
// orientation can be "horizontal" or "vertical".
func Stack(
	orientation string,
	start common.Point,
	gap styling.Length,
	elements ...common.Size,
) ([]common.Point, error) {
	points := make([]common.Point, 0, len(elements))
	currentX := start.X
	currentY := start.Y

	for _, el := range elements {
		points = append(points, common.Point{X: currentX, Y: currentY})
		switch orientation {
		case OrientationHorizontal:
			currentX += el.CX + gap
		case OrientationVertical:
			currentY += el.CY + gap
		default:
			return nil, fmt.Errorf("invalid orientation: %s", orientation)
		}
	}

	return points, nil
}

// DistributeNonUniform calculates the top or left coordinates to evenly space elements
// of variable sizes within a bound. Gaps between elements are uniform; element sizes are not.
// orientation can be "horizontal" or "vertical".
func DistributeNonUniform(
	orientation string,
	bounds common.Box,
	sizes []styling.Length,
) ([]styling.Length, error) {
	if len(sizes) == 0 {
		return nil, errors.New("sizes must contain at least one element")
	}
	if orientation != OrientationHorizontal && orientation != OrientationVertical {
		return nil, fmt.Errorf("invalid orientation: %s", orientation)
	}

	var totalAvailable, startCoord styling.Length
	switch orientation {
	case OrientationHorizontal:
		totalAvailable = bounds.CX
		startCoord = bounds.X
	case OrientationVertical:
		totalAvailable = bounds.CY
		startCoord = bounds.Y
	}

	var totalSize styling.Length
	for _, s := range sizes {
		if s < 0 {
			return nil, errors.New("element sizes must be non-negative")
		}
		totalSize += s
	}
	if totalSize > totalAvailable {
		return nil, errors.New("elements exceed available space")
	}

	if len(sizes) == 1 {
		coord := startCoord + (totalAvailable-sizes[0])/centerDivisor
		return []styling.Length{coord}, nil
	}

	gap := (totalAvailable - totalSize) / styling.Length(len(sizes)-1)
	coords := make([]styling.Length, len(sizes))
	cursor := startCoord
	for i, s := range sizes {
		coords[i] = cursor
		cursor += s + gap
	}
	return coords, nil
}

// DistributeUniform calculates the top or left coordinates to evenly space elements of identical size within a bound.
// orientation can be "horizontal" or "vertical".
// count is the number of elements to distribute.
func DistributeUniform(
	orientation string,
	bounds common.Box,
	count int,
	elementSize styling.Length,
) ([]styling.Length, error) {
	if count <= 0 {
		return nil, errors.New("count must be greater than zero")
	}
	if orientation != OrientationHorizontal && orientation != OrientationVertical {
		return nil, fmt.Errorf("invalid orientation: %s", orientation)
	}
	if count == 1 {
		switch orientation {
		case OrientationHorizontal:
			x, _ := CenterInBox(elementSize, 0, bounds)
			return []styling.Length{x}, nil
		case OrientationVertical:
			_, y := CenterInBox(0, elementSize, bounds)
			return []styling.Length{y}, nil
		}
	}

	var totalAvailable styling.Length
	var startCoord styling.Length
	switch orientation {
	case OrientationHorizontal:
		totalAvailable = bounds.CX
		startCoord = bounds.X
	case OrientationVertical:
		totalAvailable = bounds.CY
		startCoord = bounds.Y
	}

	totalElementSize := elementSize * styling.Length(count)
	if totalElementSize > totalAvailable {
		return nil, errors.New("elements exceed available space")
	}

	gap := (totalAvailable - totalElementSize) / styling.Length(count-1)
	coords := make([]styling.Length, count)
	for i := range count {
		coords[i] = startCoord + styling.Length(i)*(elementSize+gap)
	}

	return coords, nil
}
