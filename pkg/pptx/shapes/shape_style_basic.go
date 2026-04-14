package shapes

import (
	"errors"
	"fmt"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const minGradientStops = 2

// ShapeFill configures solid fill properties for one shape.
type ShapeFill struct {
	Color        string
	Transparency *float64 // 0.0=opaque, 1.0=transparent
}

// NewShapeFill creates a solid fill using a 6-digit RGB color.
func NewShapeFill(color string) ShapeFill {
	return ShapeFill{Color: common.NormalizeHexColor(color)}
}

// WithTransparency sets fill transparency in the range [0.0, 1.0] (0.0=opaque, 1.0=transparent).
func (f ShapeFill) WithTransparency(percent float64) ShapeFill {
	value := percent
	f.Transparency = &value
	return f
}

// Validate checks for validity of fill parameters.
func (f ShapeFill) Validate() error {
	if !common.IsHexColor(f.Color) {
		return fmt.Errorf("invalid color %q", f.Color)
	}
	if f.Transparency != nil && (*f.Transparency < 0 || *f.Transparency > 1) {
		return errors.New("transparency must be between 0.0 and 1.0")
	}
	return nil
}

// ShapeLine configures line style for one shape or connector.
type ShapeLine struct {
	Color string
	Width styling.Length
	Dash  string
	Cap   string
	Join  string
}

// NewShapeLine creates a line style with RGB color and EMU width.
func NewShapeLine(color string, width styling.Length) ShapeLine {
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

// WithCap sets line cap style.
func (l ShapeLine) WithCap(lineCap string) ShapeLine {
	l.Cap = NormalizeLineCap(lineCap)
	return l
}

// WithJoin sets line join style.
func (l ShapeLine) WithJoin(join string) ShapeLine {
	l.Join = NormalizeLineJoin(join)
	return l
}

// Validate checks for validity of line parameters.
func (l ShapeLine) Validate() error {
	if !common.IsHexColor(l.Color) {
		return errors.New("line color must be 6-digit RGB hex")
	}
	if l.Width <= 0 {
		return errors.New("line width must be > 0")
	}
	if !IsDrawingLineDash(l.Dash) {
		return fmt.Errorf("invalid line dash %q", l.Dash)
	}
	if !IsLineCap(l.Cap) {
		return fmt.Errorf("invalid line cap %q", l.Cap)
	}
	if !IsLineJoin(l.Join) {
		return fmt.Errorf("invalid line join %q", l.Join)
	}
	return nil
}

// ShapeGradientStop configures one gradient stop for a shape fill.
type ShapeGradientStop struct {
	PositionPct  int
	Color        string
	Transparency *float64 // 0.0=opaque, 1.0=transparent
}

// NewShapeGradientStop creates a gradient stop at one position in [0,100].
func NewShapeGradientStop(positionPct int, color string) ShapeGradientStop {
	return ShapeGradientStop{
		PositionPct: positionPct,
		Color:       common.NormalizeHexColor(color),
	}
}

// WithTransparency sets stop transparency in the range [0.0, 1.0] (0.0=opaque, 1.0=transparent).
func (s ShapeGradientStop) WithTransparency(percent float64) ShapeGradientStop {
	value := percent
	s.Transparency = &value
	return s
}

// Validate checks for validity of gradient stop parameters.
func (s ShapeGradientStop) Validate() error {
	if s.PositionPct < 0 || s.PositionPct > 100 {
		return errors.New("position must be between 0 and 100")
	}
	if !common.IsHexColor(s.Color) {
		return fmt.Errorf("invalid color %q", s.Color)
	}
	if s.Transparency != nil && (*s.Transparency < 0 || *s.Transparency > 1) {
		return errors.New("transparency must be between 0.0 and 1.0")
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
	if len(f.Stops) < minGradientStops {
		return errors.New("gradient must have at least 2 stops")
	}
	for i, stop := range f.Stops {
		if err := stop.Validate(); err != nil {
			return fmt.Errorf("stop %d invalid: %w", i, err)
		}
	}
	for i := 1; i < len(f.Stops); i++ {
		if f.Stops[i].PositionPct <= f.Stops[i-1].PositionPct {
			return errors.New("gradient stop positions must be strictly increasing")
		}
	}
	if f.AngleDeg != nil && f.Type != ShapeGradientTypeLinear {
		return errors.New("gradient angle is only supported for linear gradients")
	}
	return nil
}
