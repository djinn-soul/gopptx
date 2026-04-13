package pptx

import "github.com/djinn-soul/gopptx/pkg/pptx/shapes"

func NewShape(shapeType string, x, y, cx, cy Length) Shape {
	return shapes.NewShape(shapeType, x, y, cx, cy)
}

func NewShapeFill(color string) ShapeFill {
	return shapes.NewShapeFill(color)
}

func NewShapeLine(color string, width Length) ShapeLine {
	return shapes.NewShapeLine(color, width)
}

func NewShapeGradientStop(positionPct int, color string) ShapeGradientStop {
	return shapes.NewShapeGradientStop(positionPct, color)
}

func NewShapeGradientFill(gradientType string, stops []ShapeGradientStop) ShapeGradientFill {
	return shapes.NewShapeGradientFill(gradientType, stops)
}

func NewTextFrame() TextFrame {
	return shapes.NewTextFrame()
}

// NewSolidFill creates a new solid color fill.
func NewSolidFill(color string) *RichShapeFill {
	return shapes.NewSolidFill(color)
}

// NewNoFill creates a fill that represents "no fill" (transparent).
func NewNoFill() *RichShapeFill {
	return shapes.NewNoFill()
}

// NewPatternFill creates a new pattern fill with the specified pattern type.
func NewPatternFill(pattern PatternType) *RichShapeFill {
	return shapes.NewPatternFill(pattern)
}

// NewRichShapeLine creates a new rich line style with the specified color and width.
func NewRichShapeLine(color string, width Length) *RichShapeLine {
	return shapes.NewRichShapeLine(color, width)
}

// NewOuterShadow creates a new outer shadow with the specified color.
func NewOuterShadow(color string) *RichShapeShadow {
	return shapes.NewOuterShadow(color)
}

// NewInnerShadow creates a new inner shadow with the specified color.
func NewInnerShadow(color string) *RichShapeShadow {
	return shapes.NewInnerShadow(color)
}

// NewPerspectiveShadow creates a new perspective shadow with the specified color.
func NewPerspectiveShadow(color string) *RichShapeShadow {
	return shapes.NewPerspectiveShadow(color)
}

// NewGroupShape creates a new empty group shape at the specified position and size.
func NewGroupShape(x, y, w, h Length) GroupShape {
	return shapes.NewGroupShape(x, y, w, h)
}

// NewGroupShapeBounds creates a new group shape that auto-calculates bounds from its children.
func NewGroupShapeBounds(shapesList []Shape) GroupShape {
	return shapes.NewGroupShapeBounds(shapesList)
}

// NewFreeform creates a new freeform shape from the specified points.
func NewFreeform(points []FreeformPoint) Freeform {
	return shapes.NewFreeform(points)
}

// NewFreeformCoords creates a new freeform shape from coordinate values (in EMUs).
func NewFreeformCoords(xCoords, yCoords []int64) (Freeform, error) {
	return shapes.NewFreeformCoords(xCoords, yCoords)
}

// NewFreeformInches creates a new freeform shape from points specified in inches.
func NewFreeformInches(points [][2]float64) (Freeform, error) {
	return shapes.NewFreeformInches(points)
}

// NewFreeformClosed creates a new closed freeform shape.
func NewFreeformClosed(points []FreeformPoint) Freeform {
	return shapes.NewFreeformClosed(points)
}

// NewFreeformOpen creates a new open freeform shape (line).
func NewFreeformOpen(points []FreeformPoint) Freeform {
	return shapes.NewFreeformOpen(points)
}
