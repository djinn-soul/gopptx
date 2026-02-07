package pptx

// ShapeDefinition allows external shape builders to plug into slide composition.
type ShapeDefinition interface {
	ToShape() Shape
}

// ToShape returns the shape itself and satisfies ShapeDefinition.
func (s Shape) ToShape() Shape {
	return s
}
