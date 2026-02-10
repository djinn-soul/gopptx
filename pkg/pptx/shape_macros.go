package pptx

// NewRectangle creates a rectangle shape with given inch dimensions.
func NewRectangle(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeRectangle, Inches(x), Inches(y), Inches(w), Inches(h))
}

// NewEllipse creates an ellipse shape with given inch dimensions.
func NewEllipse(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeEllipse, Inches(x), Inches(y), Inches(w), Inches(h))
}

// NewTextBox creates a text box shape with given text and inch dimensions.
func NewTextBox(text string, x, y, w, h float64) Shape {
	return NewRectangle(x, y, w, h).
		WithText(text)
}

// NewRoundedRectangle creates a rounded rectangle shape with given inch dimensions.
func NewRoundedRectangle(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeRoundedRectangle, Inches(x), Inches(y), Inches(w), Inches(h))
}

// NewTriangle creates a triangle shape with given inch dimensions.
func NewTriangle(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeTriangle, Inches(x), Inches(y), Inches(w), Inches(h))
}

// NewRightTriangle creates a right triangle shape with given inch dimensions.
func NewRightTriangle(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeRightTriangle, Inches(x), Inches(y), Inches(w), Inches(h))
}

// NewDiamond creates a diamond shape with given inch dimensions.
func NewDiamond(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeDiamond, Inches(x), Inches(y), Inches(w), Inches(h))
}

// NewPentagon creates a pentagon shape with given inch dimensions.
func NewPentagon(x, y, w, h float64) Shape {
	return NewShape(ShapeTypePentagon, Inches(x), Inches(y), Inches(w), Inches(h))
}

// NewHexagon creates a hexagon shape with given inch dimensions.
func NewHexagon(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeHexagon, Inches(x), Inches(y), Inches(w), Inches(h))
}

// NewParallelogram creates a parallelogram shape with given inch dimensions.
func NewParallelogram(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeParallelogram, Inches(x), Inches(y), Inches(w), Inches(h))
}

// NewFlowChartProcess creates a flowchart process shape with given inch dimensions.
func NewFlowChartProcess(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeFlowChartProcess, Inches(x), Inches(y), Inches(w), Inches(h))
}

// NewFlowChartDecision creates a flowchart decision shape with given inch dimensions.
func NewFlowChartDecision(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeFlowChartDecision, Inches(x), Inches(y), Inches(w), Inches(h))
}

// NewFlowChartTerminator creates a flowchart terminator shape with given inch dimensions.
func NewFlowChartTerminator(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeFlowChartTerminator, Inches(x), Inches(y), Inches(w), Inches(h))
}

// NewRightArrow creates a right arrow shape with given inch dimensions.
func NewRightArrow(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeRightArrow, Inches(x), Inches(y), Inches(w), Inches(h))
}

// NewLeftArrow creates a left arrow shape with given inch dimensions.
func NewLeftArrow(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeLeftArrow, Inches(x), Inches(y), Inches(w), Inches(h))
}

// NewUpArrow creates an up arrow shape with given inch dimensions.
func NewUpArrow(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeUpArrow, Inches(x), Inches(y), Inches(w), Inches(h))
}

// NewDownArrow creates a down arrow shape with given inch dimensions.
func NewDownArrow(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeDownArrow, Inches(x), Inches(y), Inches(w), Inches(h))
}

// NewCloud creates a cloud shape with given inch dimensions.
func NewCloud(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeCloud, Inches(x), Inches(y), Inches(w), Inches(h))
}

// NewCircle creates a circle (ellipse with equal width and height) with given diameter in inches.
func NewCircle(x, y, diameter float64) Shape {
	return NewShape(ShapeTypeEllipse, Inches(x), Inches(y), Inches(diameter), Inches(diameter))
}

// NewStar creates a 5-pointed star shape with given size in inches.
func NewStar(x, y, size float64) Shape {
	return NewShape(ShapeTypeStar5, Inches(x), Inches(y), Inches(size), Inches(size))
}

// NewHeart creates a heart shape with given size in inches.
func NewHeart(x, y, size float64) Shape {
	return NewShape(ShapeTypeHeart, Inches(x), Inches(y), Inches(size), Inches(size))
}

// NewFlowChartDocument creates a flowchart document shape with given inch dimensions.
func NewFlowChartDocument(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeFlowChartDocument, Inches(x), Inches(y), Inches(w), Inches(h))
}

// NewFlowChartData creates a flowchart data shape (parallelogram) with given inch dimensions.
func NewFlowChartData(x, y, w, h float64) Shape {
	return NewShape(ShapeTypeFlowChartData, Inches(x), Inches(y), Inches(w), Inches(h))
}

// NewBadge creates a badge (rounded rectangle with text) at a default size (1.5x0.4 inches).
func NewBadge(text string, x, y float64, color string) Shape {
	return NewShape(ShapeTypeRoundedRectangle, Inches(x), Inches(y), Inches(1.5), Inches(0.4)).
		WithFill(NewShapeFill(color)).
		WithText(text)
}
