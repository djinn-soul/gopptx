package layout

import (
	"fmt"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
)

const (
	// SlideWidth is the standard width of a 4:3 slide in EMU.
	SlideWidth int64 = 9144000
	// SlideHeight is the standard height of a 4:3 slide in EMU.
	SlideHeight int64 = 6858000
)

const (
	// OrientationHorizontal represents left-to-right or horizontal distribution.
	OrientationHorizontal = "horizontal"
	// OrientationVertical represents top-to-bottom or vertical distribution.
	OrientationVertical = "vertical"
)

// Center calculates the (X, Y) coordinates to center an element of size (cx, cy)
// within the standard slide bounds.
func Center(cx, cy int64) (x, y int64) {
	x = (SlideWidth - cx) / 2
	y = (SlideHeight - cy) / 2
	return x, y
}

// CenterInBox calculates the (X, Y) coordinates to center an element of size (cx, cy)
// within a specific bounding box.
func CenterInBox(cx, cy int64, bounds common.Box) (x, y int64) {
	x = bounds.X + (bounds.CX-cx)/2
	y = bounds.Y + (bounds.CY-cy)/2
	return x, y
}

// Grid calculates a sequence of bounding boxes for a grid layout.
// rows and cols must be greater than zero. margin is the spacing between elements.
// The grid fills the standard slide area.
func Grid(rows, cols int, margin int64) ([]common.Box, error) {
	return GridInBox(rows, cols, margin, common.Box{X: 0, Y: 0, CX: SlideWidth, CY: SlideHeight})
}

// GridInBox calculates a sequence of bounding boxes for a grid layout within specific bounds.
func GridInBox(rows, cols int, margin int64, bounds common.Box) ([]common.Box, error) {
	if rows <= 0 || cols <= 0 {
		return nil, fmt.Errorf("rows and cols must be greater than zero")
	}

	totalMarginX := margin * int64(cols-1)
	totalMarginY := margin * int64(rows-1)

	if totalMarginX >= bounds.CX || totalMarginY >= bounds.CY {
		return nil, fmt.Errorf("margins exceed bounding box dimensions")
	}

	elementCX := (bounds.CX - totalMarginX) / int64(cols)
	elementCY := (bounds.CY - totalMarginY) / int64(rows)

	boxes := make([]common.Box, 0, rows*cols)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			x := bounds.X + int64(c)*(elementCX+margin)
			y := bounds.Y + int64(r)*(elementCY+margin)
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
func Stack(orientation string, start common.Point, gap int64, elements ...common.Size) ([]common.Point, error) {
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

// Distribute calculates the top or left coordinates to evenly space elements within a bound.
// orientation can be "horizontal" or "vertical".
// count is the number of elements to distribute.
func Distribute(orientation string, bounds common.Box, count int, elementSize int64) ([]int64, error) {
	if count <= 0 {
		return nil, fmt.Errorf("count must be greater than zero")
	}
	if orientation != OrientationHorizontal && orientation != OrientationVertical {
		return nil, fmt.Errorf("invalid orientation: %s", orientation)
	}
	if count == 1 {
		switch orientation {
		case OrientationHorizontal:
			x, _ := CenterInBox(elementSize, 0, bounds)
			return []int64{x}, nil
		case OrientationVertical:
			_, y := CenterInBox(0, elementSize, bounds)
			return []int64{y}, nil
		}
	}

	var totalAvailable int64
	var startCoord int64
	switch orientation {
	case OrientationHorizontal:
		totalAvailable = bounds.CX
		startCoord = bounds.X
	case OrientationVertical:
		totalAvailable = bounds.CY
		startCoord = bounds.Y
	}

	totalElementSize := elementSize * int64(count)
	if totalElementSize > totalAvailable {
		return nil, fmt.Errorf("elements exceed available space")
	}

	gap := (totalAvailable - totalElementSize) / int64(count-1)
	coords := make([]int64, count)
	for i := 0; i < count; i++ {
		coords[i] = startCoord + int64(i)*(elementSize+gap)
	}

	return coords, nil
}
