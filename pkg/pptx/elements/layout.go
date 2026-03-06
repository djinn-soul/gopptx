package elements

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// SlideLayout defines a slide layout with shapes and placeholders.
type SlideLayout struct {
	Name         string
	Shapes       []shapes.ShapeDefinition
	Images       []shapes.Image
	Placeholders []shapes.Placeholder
}

// NewSlideLayout creates a new slide layout with the given name.
func NewSlideLayout(name string) *SlideLayout {
	return &SlideLayout{Name: name}
}

// AddShape adds a shape to the slide layout.
func (l *SlideLayout) AddShape(sd shapes.ShapeDefinition) *SlideLayout {
	l.Shapes = append(l.Shapes, sd)
	return l
}

// AddImage adds an image to the slide layout.
func (l *SlideLayout) AddImage(img shapes.Image) *SlideLayout {
	l.Images = append(l.Images, img)
	return l
}

// AddPlaceholder adds a placeholder to the slide layout.
func (l *SlideLayout) AddPlaceholder(ph shapes.Placeholder) *SlideLayout {
	l.Placeholders = append(l.Placeholders, ph)
	return l
}
