package shapes

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// GroupShape represents a group of shapes that move and transform together.
type GroupShape struct {
	X       styling.Length
	Y       styling.Length
	CX      styling.Length
	CY      styling.Length
	Name    string
	AltText string
	Shapes  []Shape
}

// NewGroupShape creates a new empty group shape at the specified position and size.
func NewGroupShape(x, y, cx, cy styling.Length) GroupShape {
	return GroupShape{
		X:      x,
		Y:      y,
		CX:     cx,
		CY:     cy,
		Name:   "",
		Shapes: []Shape{},
	}
}

// NewGroupShapeBounds creates a new group shape that auto-calculates bounds from its children.
func NewGroupShapeBounds(shapes []Shape) GroupShape {
	if len(shapes) == 0 {
		return NewGroupShape(0, 0, 0, 0)
	}

	// Calculate bounds from children
	minX := shapes[0].X
	minY := shapes[0].Y
	maxX := shapes[0].X + shapes[0].CX
	maxY := shapes[0].Y + shapes[0].CY

	for i := 1; i < len(shapes); i++ {
		s := shapes[i]
		if s.X < minX {
			minX = s.X
		}
		if s.Y < minY {
			minY = s.Y
		}
		if s.X+s.CX > maxX {
			maxX = s.X + s.CX
		}
		if s.Y+s.CY > maxY {
			maxY = s.Y + s.CY
		}
	}

	return GroupShape{
		X:      minX,
		Y:      minY,
		CX:     maxX - minX,
		CY:     maxY - minY,
		Name:   "",
		Shapes: shapes,
	}
}

// WithName sets the name of the group shape.
func (g GroupShape) WithName(name string) GroupShape {
	g.Name = name
	return g
}

// WithAltText sets the alternative text for accessibility.
func (g GroupShape) WithAltText(text string) GroupShape {
	g.AltText = text
	return g
}

// AddShape adds a shape to the group.
func (g GroupShape) AddShape(shape Shape) GroupShape {
	g.Shapes = append(g.Shapes, shape)
	return g
}

// AddShapes adds multiple shapes to the group.
func (g GroupShape) AddShapes(shapes ...Shape) GroupShape {
	g.Shapes = append(g.Shapes, shapes...)
	return g
}

// Validate checks the group shape for validity.
func (g GroupShape) Validate() error {
	if g.CX <= 0 || g.CY <= 0 {
		return nil // Allow zero-size groups, they'll be recalculated
	}
	return nil
}

// ToShape converts the group shape to a Shape for compatibility.
// Note: This is a limited conversion - for full group support, use the GroupShape directly.
func (g GroupShape) ToShape() Shape {
	s := NewShape(ShapeTypeGroup, g.X, g.Y, g.CX, g.CY)
	s.Name = g.Name
	s.AltText = g.AltText
	return s
}

// IsGroupShape returns true for group shapes.
func (s Shape) IsGroupShape() bool {
	return s.Type == ShapeTypeGroup
}

// ShapeTypeGroup is the shape type constant for groups.
const ShapeTypeGroup = "grpSp"
