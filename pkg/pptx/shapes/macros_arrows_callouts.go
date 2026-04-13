package shapes

func NewRightArrow(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeRightArrow, x, y, w, h)
}

func NewLeftArrow(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeLeftArrow, x, y, w, h)
}

func NewUpArrow(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeUpArrow, x, y, w, h)
}

func NewDownArrow(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeDownArrow, x, y, w, h)
}

func NewLeftRightArrow(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeLeftRightArrow, x, y, w, h)
}

func NewUpDownArrow(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeUpDownArrow, x, y, w, h)
}

func NewQuadArrow(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeQuadArrow, x, y, w, h)
}

func NewBentArrow(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeBentArrow, x, y, w, h)
}

func NewUturnArrow(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeUturnArrow, x, y, w, h)
}

func NewCircularArrow(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeCircularArrow, x, y, w, h)
}

func NewChevron(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeChevronArrow, x, y, w, h)
}

func NewWedgeRectCallout(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeWedgeRectCallout, x, y, w, h)
}

func NewWedgeRRectCallout(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeWedgeRRectCallout, x, y, w, h)
}

func NewWedgeEllipseCallout(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeWedgeEllipseCallout, x, y, w, h)
}

func NewCloudCallout(x, y, w, h float64) Shape {
	return newShapeInches(ShapeTypeCloudCallout, x, y, w, h)
}
