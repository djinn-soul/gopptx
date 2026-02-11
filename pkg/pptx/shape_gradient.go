package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

type (
	// ShapeGradientStop configures one gradient stop for a shape fill.
	ShapeGradientStop = elements.ShapeGradientStop
	// ShapeGradientFill configures gradient fill properties for one shape.
	ShapeGradientFill = elements.ShapeGradientFill
)

const (
	ShapeGradientTypeLinear      = elements.ShapeGradientTypeLinear
	ShapeGradientTypeRadial      = elements.ShapeGradientTypeRadial
	ShapeGradientTypeRectangular = elements.ShapeGradientTypeRectangular
	ShapeGradientTypePath        = elements.ShapeGradientTypePath
)

func NewShapeGradientStop(positionPct int, color string) ShapeGradientStop {
	return elements.NewShapeGradientStop(positionPct, color)
}

func NewShapeGradientFill(gradientType string, stops []ShapeGradientStop) ShapeGradientFill {
	return elements.NewShapeGradientFill(gradientType, stops)
}

func normalizeShapeGradientType(gradientType string) string {
	return elements.NormalizeShapeGradientType(gradientType)
}
