package shapes

import (
	"fmt"

	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
)

// ShapeFill configures solid fill properties for one shape.
type ShapeFill struct {
	Color           string
	TransparencyPct *int
}

// NewShapeFill creates a solid fill using a 6-digit RGB color.
func NewShapeFill(color string) ShapeFill {
	return ShapeFill{Color: common.NormalizeHexColor(color)}
}

// WithTransparency sets fill transparency percentage in the range [0,100].
func (f ShapeFill) WithTransparency(percent int) ShapeFill {
	value := percent
	f.TransparencyPct = &value
	return f
}

// Validate checks for validity of fill parameters.
func (f ShapeFill) Validate() error {
	if !common.IsHexColor(f.Color) {
		return fmt.Errorf("invalid color %q", f.Color)
	}
	if f.TransparencyPct != nil && (*f.TransparencyPct < 0 || *f.TransparencyPct > 100) {
		return fmt.Errorf("transparency must be between 0 and 100")
	}
	return nil
}

// ShapeLine configures line style for one shape or connector.
type ShapeLine struct {
	Color string
	Width int64
	Dash  string
}

// NewShapeLine creates a line style with RGB color and EMU width.
func NewShapeLine(color string, width int64) ShapeLine {
	return ShapeLine{
		Color: common.NormalizeHexColor(color),
		Width: width,
		Dash:  LineDashSolid,
	}
}

// WithDash sets line dash style.
func (l ShapeLine) WithDash(dash string) ShapeLine {
	l.Dash = NormalizeDrawingLineDash(dash)
	return l
}

// Validate checks for validity of line parameters.
func (l ShapeLine) Validate() error {
	if !common.IsHexColor(l.Color) {
		return fmt.Errorf("line color must be 6-digit RGB hex")
	}
	if l.Width <= 0 {
		return fmt.Errorf("line width must be > 0")
	}
	if !IsDrawingLineDash(l.Dash) {
		return fmt.Errorf("invalid line dash %q", l.Dash)
	}
	return nil
}

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
		Color:       common.NormalizeHexColor(color),
	}
}

// WithTransparency sets stop transparency percentage in the range [0,100].
func (s ShapeGradientStop) WithTransparency(percent int) ShapeGradientStop {
	value := percent
	s.TransparencyPct = &value
	return s
}

// Validate checks for validity of gradient stop parameters.
func (s ShapeGradientStop) Validate() error {
	if s.PositionPct < 0 || s.PositionPct > 100 {
		return fmt.Errorf("position must be between 0 and 100")
	}
	if !common.IsHexColor(s.Color) {
		return fmt.Errorf("invalid color %q", s.Color)
	}
	if s.TransparencyPct != nil && (*s.TransparencyPct < 0 || *s.TransparencyPct > 100) {
		return fmt.Errorf("transparency must be between 0 and 100")
	}
	return nil
}

// ShapeGradientFill configures gradient fill properties for one shape.
type ShapeGradientFill struct {
	Type     string
	Stops    []ShapeGradientStop
	AngleDeg *int
}

// NewShapeGradientFill creates one gradient fill.
func NewShapeGradientFill(gradientType string, stops []ShapeGradientStop) ShapeGradientFill {
	copiedStops := make([]ShapeGradientStop, len(stops))
	copy(copiedStops, stops)
	return ShapeGradientFill{
		Type:  NormalizeShapeGradientType(gradientType),
		Stops: copiedStops,
	}
}

// WithLinearAngle sets the linear gradient angle in degrees.
func (f ShapeGradientFill) WithLinearAngle(degrees int) ShapeGradientFill {
	value := degrees
	f.AngleDeg = &value
	return f
}

// Validate checks for validity of gradient fill parameters.
func (f ShapeGradientFill) Validate() error {
	if !IsShapeGradientType(f.Type) {
		return fmt.Errorf("invalid gradient type %q", f.Type)
	}
	if len(f.Stops) < 2 {
		return fmt.Errorf("gradient must have at least 2 stops")
	}
	for i, stop := range f.Stops {
		if err := stop.Validate(); err != nil {
			return fmt.Errorf("stop %d invalid: %w", i, err)
		}
	}
	// Validate stop constraints
	for i := 1; i < len(f.Stops); i++ {
		if f.Stops[i].PositionPct <= f.Stops[i-1].PositionPct {
			return fmt.Errorf("gradient stop positions must be strictly increasing")
		}
	}

	if f.AngleDeg != nil && f.Type != ShapeGradientTypeLinear {
		return fmt.Errorf("gradient angle is only supported for linear gradients")
	}
	return nil
}

// Shape is one auto shape.
type Shape struct {
	Type         string
	X            int64
	Y            int64
	CX           int64
	CY           int64
	Fill         *ShapeFill
	Line         *ShapeLine
	GradientFill *ShapeGradientFill
	Text         string
	RotationDeg  *int
	Hyperlink    *action.Hyperlink
	AltText      string
	IsDecorative bool
}

// NewShape creates one shape.
func NewShape(shapeType string, x, y, cx, cy int64) Shape {
	return Shape{
		Type: NormalizeShapeType(shapeType),
		X:    x,
		Y:    y,
		CX:   cx,
		CY:   cy,
	}
}

// WithFill applies solid fill to a shape.
func (s Shape) WithFill(fill ShapeFill) Shape {
	s.Fill = &fill
	s.GradientFill = nil
	return s
}

// WithLine applies a line style to a shape.
func (s Shape) WithLine(line ShapeLine) Shape {
	s.Line = &line
	return s
}

// WithText sets text rendered inside the shape.
func (s Shape) WithText(text string) Shape {
	s.Text = text
	return s
}

// WithRotation rotates shape geometry in degrees.
func (s Shape) WithRotation(degrees int) Shape {
	value := degrees
	s.RotationDeg = &value
	return s
}

// ShapeDefinition allows external shape builders to plug into slide composition.
type ShapeDefinition interface {
	ToShape() Shape
}

// ToShape returns the shape itself and satisfies ShapeDefinition.
func (s Shape) ToShape() Shape {
	return s
}

// WithHyperlink attaches a clickable hyperlink to the shape.
func (s Shape) WithHyperlink(hyperlink action.Hyperlink) Shape {
	s.Hyperlink = &hyperlink
	return s
}

// WithAltText sets the alternative text for accessibility.
func (s Shape) WithAltText(text string) Shape {
	s.AltText = text
	return s
}

// WithDecorative marks the shape as decorative (ignored by screen readers).
func (s Shape) WithDecorative(enabled bool) Shape {
	s.IsDecorative = enabled
	return s
}

// WithGradientFill applies gradient fill to a shape.
func (s Shape) WithGradientFill(fill ShapeGradientFill) Shape {
	s.GradientFill = &fill
	s.Fill = nil
	return s
}

// Validate checks for validity of shape parameters.
func (s Shape) Validate(slideIndex, shapeIndex int) error {
	if s.X < 0 || s.Y < 0 {
		return fmt.Errorf("shape %d on slide %d position cannot be negative", shapeIndex, slideIndex)
	}
	if s.CX <= 0 || s.CY <= 0 {
		return fmt.Errorf("shape %d on slide %d size must be > 0", shapeIndex, slideIndex)
	}

	if !IsShapeType(s.Type) {
		return fmt.Errorf("shape %d type %q is invalid on slide %d", shapeIndex, s.Type, slideIndex)
	}

	if s.Fill != nil {
		if s.GradientFill != nil {
			return fmt.Errorf("shape %d (type %q) on slide %d cannot set both solid and gradient fill", shapeIndex, s.Type, slideIndex)
		}
		if err := s.Fill.Validate(); err != nil {
			return fmt.Errorf("shape %d on slide %d has invalid fill: %w", shapeIndex, slideIndex, err)
		}
	}

	if s.Line != nil {
		if err := s.Line.Validate(); err != nil {
			return fmt.Errorf("shape %d on slide %d has invalid line: %w", shapeIndex, slideIndex, err)
		}
	}

	if s.GradientFill != nil {
		if err := s.GradientFill.Validate(); err != nil {
			return fmt.Errorf("shape %d on slide %d has invalid gradient fill: %w", shapeIndex, slideIndex, err)
		}
	}

	if s.RotationDeg != nil {
		if *s.RotationDeg < -360 || *s.RotationDeg > 360 {
			return fmt.Errorf("shape %d on slide %d rotation must be in [-360,360]", shapeIndex, slideIndex)
		}
	}

	if s.Hyperlink != nil {
		if err := s.Hyperlink.Validate(); err != nil {
			return err
		}
	}
	return nil
}
