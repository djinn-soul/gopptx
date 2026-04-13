package shapes

import "github.com/djinn-soul/gopptx/pkg/pptx/styling"

const (
	defaultBadgeWidthInches  = 1.5
	defaultBadgeHeightInches = 0.4
)

func newShapeInches(shapeType string, x, y, w, h float64) Shape {
	return NewShape(
		shapeType,
		styling.Inches(x),
		styling.Inches(y),
		styling.Inches(w),
		styling.Inches(h),
	)
}

func newSquareShape(shapeType string, x, y, size float64) Shape {
	return newShapeInches(shapeType, x, y, size, size)
}
