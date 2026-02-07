package pptx

import "strings"

const (
	// ShapeGradientTypeLinear renders a linear gradient.
	ShapeGradientTypeLinear = "linear"
	// ShapeGradientTypeRadial renders a radial gradient.
	ShapeGradientTypeRadial = "radial"
	// ShapeGradientTypeRectangular renders a rectangular gradient.
	ShapeGradientTypeRectangular = "rectangular"
	// ShapeGradientTypePath renders a shape-path gradient.
	ShapeGradientTypePath = "path"
)

// ShapeGradientStop configures one gradient stop for a shape fill.
type ShapeGradientStop struct {
	PositionPct     int
	Color           string
	TransparencyPct *int
}

// NewShapeGradientStop creates a gradient stop at one position in [0,100].
func NewShapeGradientStop(positionPct int, color string) ShapeGradientStop {
	return ShapeGradientStop{
		PositionPct: positionPct,
		Color:       normalizeHexColor(color),
	}
}

// WithTransparency sets stop transparency percentage in the range [0,100].
func (s ShapeGradientStop) WithTransparency(percent int) ShapeGradientStop {
	value := percent
	s.TransparencyPct = &value
	return s
}

// ShapeGradientFill configures gradient fill properties for one shape.
type ShapeGradientFill struct {
	Type     string
	Stops    []ShapeGradientStop
	AngleDeg *int
}

// NewShapeGradientFill creates one gradient fill with explicit type and stops.
func NewShapeGradientFill(gradientType string, stops []ShapeGradientStop) ShapeGradientFill {
	copiedStops := make([]ShapeGradientStop, len(stops))
	copy(copiedStops, stops)
	return ShapeGradientFill{
		Type:  normalizeShapeGradientType(gradientType),
		Stops: copiedStops,
	}
}

// WithLinearAngle sets the linear gradient angle in degrees.
func (f ShapeGradientFill) WithLinearAngle(degrees int) ShapeGradientFill {
	value := degrees
	f.AngleDeg = &value
	return f
}

// WithGradientFill applies gradient fill to a shape.
func (s Shape) WithGradientFill(fill ShapeGradientFill) Shape {
	value := ShapeGradientFill{
		Type:  normalizeShapeGradientType(fill.Type),
		Stops: append([]ShapeGradientStop(nil), fill.Stops...),
	}
	if fill.AngleDeg != nil {
		angle := *fill.AngleDeg
		value.AngleDeg = &angle
	}
	s.GradientFill = &value
	s.Fill = nil
	return s
}

func normalizeShapeGradientType(gradientType string) string {
	switch strings.ToLower(strings.TrimSpace(gradientType)) {
	case ShapeGradientTypeLinear:
		return ShapeGradientTypeLinear
	case ShapeGradientTypeRadial, "radial-gradient", "radial_gradient":
		return ShapeGradientTypeRadial
	case ShapeGradientTypeRectangular, "rectangular-gradient", "rectangular_gradient", "rect":
		return ShapeGradientTypeRectangular
	case ShapeGradientTypePath, "path-gradient", "path_gradient":
		return ShapeGradientTypePath
	default:
		return strings.TrimSpace(gradientType)
	}
}

func isShapeGradientType(gradientType string) bool {
	switch normalizeShapeGradientType(gradientType) {
	case ShapeGradientTypeLinear, ShapeGradientTypeRadial, ShapeGradientTypeRectangular, ShapeGradientTypePath:
		return true
	default:
		return false
	}
}
