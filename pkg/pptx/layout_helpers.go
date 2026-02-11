package pptx

import (
	"fmt"
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

// Point represents a 2D coordinate in EMU.
type Point struct {
	X int64
	Y int64
}

// Size represents dimensions in EMU.
type Size struct {
	CX int64
	CY int64
}

// Box represents a rectangular region with position and size in EMU.
type Box struct {
	X  int64
	Y  int64
	CX int64
	CY int64
}

// Center calculates the (X, Y) coordinates to center an element of size (cx, cy)
// within the standard slide bounds.
//
// Example:
//
//	x, y := Center(Inches(4), Inches(3))
func Center(cx, cy int64) (x, y int64) {
	x = (SlideWidth - cx) / 2
	y = (SlideHeight - cy) / 2
	return x, y
}

// CenterInBox calculates the (X, Y) coordinates to center an element of size (cx, cy)
// within a specific bounding box.
func CenterInBox(cx, cy int64, bounds Box) (x, y int64) {
	x = bounds.X + (bounds.CX-cx)/2
	y = bounds.Y + (bounds.CY-cy)/2
	return x, y
}

// Grid calculates a sequence of bounding boxes for a grid layout.
// rows and cols must be greater than zero. margin is the spacing between elements.
// The grid fills the standard slide area.
func Grid(rows, cols int, margin int64) ([]Box, error) {
	return GridInBox(rows, cols, margin, Box{0, 0, SlideWidth, SlideHeight})
}

// GridInBox calculates a sequence of bounding boxes for a grid layout within specific bounds.
func GridInBox(rows, cols int, margin int64, bounds Box) ([]Box, error) {
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

	boxes := make([]Box, 0, rows*cols)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			x := bounds.X + int64(c)*(elementCX+margin)
			y := bounds.Y + int64(r)*(elementCY+margin)
			boxes = append(boxes, Box{
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
func Stack(orientation string, start Point, gap int64, elements ...Size) ([]Point, error) {
	points := make([]Point, 0, len(elements))
	currentX := start.X
	currentY := start.Y

	for _, el := range elements {
		points = append(points, Point{currentX, currentY})
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
func Distribute(orientation string, bounds Box, count int, elementSize int64) ([]int64, error) {
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
