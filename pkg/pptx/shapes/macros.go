package shapes

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// NewRectangle creates a rectangle shape with given inch dimensions.
func NewRectangle(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeRectangle, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewEllipse creates an ellipse shape with given inch dimensions.
func NewEllipse(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeEllipse, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewTextBox creates a text box shape with given text and inch dimensions.
func NewTextBox(text string, x, y, w, h float64) Shape {
	return NewRectangle(x, y, w, h).
		WithText(text)
}

// NewRoundedRectangle creates a rounded rectangle shape with given inch dimensions.
func NewRoundedRectangle(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeRoundedRectangle, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewTriangle creates a triangle shape with given inch dimensions.
func NewTriangle(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeTriangle, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewRightTriangle creates a right triangle shape with given inch dimensions.
func NewRightTriangle(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeRightTriangle, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewDiamond creates a diamond shape with given inch dimensions.
func NewDiamond(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeDiamond, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewPentagon creates a pentagon shape with given inch dimensions.
func NewPentagon(x, y, w, h float64) Shape {
	return NewShape(ShapeTypePentagon, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewHexagon creates a hexagon shape with given inch dimensions.
func NewHexagon(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeHexagon, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewParallelogram creates a parallelogram shape with given inch dimensions.
func NewParallelogram(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeParallelogram, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewFlowChartProcess creates a flowchart process shape with given inch dimensions.
func NewFlowChartProcess(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeFlowChartProcess, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewFlowChartDecision creates a flowchart decision shape with given inch dimensions.
func NewFlowChartDecision(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeFlowChartDecision, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewFlowChartTerminator creates a flowchart terminator shape with given inch dimensions.
func NewFlowChartTerminator(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeFlowChartTerminator, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewRightArrow creates a right arrow shape with given inch dimensions.
func NewRightArrow(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeRightArrow, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewLeftArrow creates a left arrow shape with given inch dimensions.
func NewLeftArrow(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeLeftArrow, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewUpArrow creates an up arrow shape with given inch dimensions.
func NewUpArrow(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeUpArrow, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewDownArrow creates a down arrow shape with given inch dimensions.
func NewDownArrow(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeDownArrow, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewCloud creates a cloud shape with given inch dimensions.
func NewCloud(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeCloud, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewCircle creates a circle (ellipse with equal width and height) with given diameter in inches.
func NewCircle(x, y, diameter float64) Shape {
	return NewShape(ShapeTypeEllipse, styling.Inches(x), styling.Inches(y), styling.Inches(diameter), styling.Inches(diameter))
}

// NewStar creates a 5-pointed star shape with given size in inches.
func NewStar(x, y, size float64) Shape {
	return NewShape(ShapeTypeStar5, styling.Inches(x), styling.Inches(y), styling.Inches(size), styling.Inches(size))
}

// NewHeart creates a heart shape with given size in inches.
func NewHeart(x, y, size float64) Shape {
	return NewShape(ShapeTypeHeart, styling.Inches(x), styling.Inches(y), styling.Inches(size), styling.Inches(size))
}

// NewFlowChartDocument creates a flowchart document shape with given inch dimensions.
func NewFlowChartDocument(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeFlowChartDocument, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewFlowChartData creates a flowchart data shape (parallelogram) with given inch dimensions.
func NewFlowChartData(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeFlowChartData, styling.Inches(x), styling.Inches(y), styling.Inches(w), styling.Inches(h))
}

// NewBadge creates a badge (rounded rectangle with text) at a default size (1.5x0.4 inches).
func NewBadge(text string, x, y float64, color string) Shape {
	if color == "" {
		color = styling.ColorMaterialGreen
	}
	return NewShape(ShapeTypeRoundedRectangle, styling.Inches(x), styling.Inches(y), styling.Inches(1.5), styling.Inches(0.4)).
		WithFill(NewShapeFill(color)).
		WithText(text)
}
