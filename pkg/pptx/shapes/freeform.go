package shapes

import (
	"errors"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const minFreeformPoints = 2

// FreeformPoint represents a point in a freeform shape path.
type FreeformPoint struct {
	X styling.Length
	Y styling.Length
}

// Freeform represents a custom-geometry (freeform) shape.
type Freeform struct {
	Points       []FreeformPoint
	ClosePath    bool
	Fill         *ShapeFill
	Line         *ShapeLine
	GradientFill *ShapeGradientFill
	RichFill     *RichShapeFill
	RichLine     *RichShapeLine
	RichShadow   *RichShapeShadow
	Text         string
	RotationDeg  *int
	Name         string
	AltText      string
	IsDecorative bool
	TextFrame    *TextFrame
	Effects      *ShapeEffects
}

// NewFreeform creates a new freeform shape from the specified points.
func NewFreeform(points []FreeformPoint) Freeform {
	return Freeform{
		Points:    points,
		ClosePath: true,
	}
}

// NewFreeformCoords creates a new freeform shape from coordinate values (in EMUs).
func NewFreeformCoords(xCoords, yCoords []int64) (Freeform, error) {
	if len(xCoords) != len(yCoords) {
		return Freeform{}, errors.New("x and y coordinate slices must have the same length")
	}
	if len(xCoords) < minFreeformPoints {
		return Freeform{}, errors.New("freeform requires at least 2 points")
	}

	points := make([]FreeformPoint, len(xCoords))
	for i := range xCoords {
		points[i] = FreeformPoint{
			X: styling.Emu(xCoords[i]),
			Y: styling.Emu(yCoords[i]),
		}
	}

	return NewFreeform(points), nil
}

// NewFreeformInches creates a new freeform shape from points specified in inches.
func NewFreeformInches(points [][2]float64) (Freeform, error) {
	if len(points) < minFreeformPoints {
		return Freeform{}, errors.New("freeform requires at least 2 points")
	}

	fp := make([]FreeformPoint, len(points))
	for i, p := range points {
		fp[i] = FreeformPoint{
			X: styling.Inches(p[0]),
			Y: styling.Inches(p[1]),
		}
	}

	return Freeform{
		Points:    fp,
		ClosePath: true,
	}, nil
}

// NewFreeformClosed creates a new closed freeform shape.
func NewFreeformClosed(points []FreeformPoint) Freeform {
	return Freeform{
		Points:    points,
		ClosePath: true,
	}
}

// NewFreeformOpen creates a new open freeform shape (line).
func NewFreeformOpen(points []FreeformPoint) Freeform {
	return Freeform{
		Points:    points,
		ClosePath: false,
	}
}

// WithClosePath sets whether the freeform path should be closed.
func (f Freeform) WithClosePath(closed bool) Freeform {
	f.ClosePath = closed
	return f
}

// WithFill applies solid fill to the freeform.
func (f Freeform) WithFill(fill ShapeFill) Freeform {
	f.Fill = &fill
	f.RichFill = nil
	f.GradientFill = nil
	return f
}

// WithGradientFill applies gradient fill to the freeform.
func (f Freeform) WithGradientFill(fill ShapeGradientFill) Freeform {
	f.GradientFill = &fill
	f.Fill = nil
	f.RichFill = nil
	return f
}

// WithRichFill applies a rich fill to the freeform.
func (f Freeform) WithRichFill(fill *RichShapeFill) Freeform {
	f.RichFill = fill
	f.Fill = nil
	f.GradientFill = nil
	return f
}

// WithLine applies a line style to the freeform.
func (f Freeform) WithLine(line ShapeLine) Freeform {
	f.Line = &line
	f.RichLine = nil
	return f
}

// WithRichLine applies a rich line style to the freeform.
func (f Freeform) WithRichLine(line *RichShapeLine) Freeform {
	f.RichLine = line
	f.Line = nil
	return f
}

// WithRichShadow applies a rich shadow effect to the freeform.
func (f Freeform) WithRichShadow(shadow *RichShapeShadow) Freeform {
	f.RichShadow = shadow
	if f.Effects == nil {
		f.Effects = &ShapeEffects{}
	}
	f.Effects.Shadow = shadow != nil
	return f
}

// WithText sets text rendered inside the freeform.
func (f Freeform) WithText(text string) Freeform {
	f.Text = text
	return f
}

// WithRotation rotates the freeform geometry in degrees.
func (f Freeform) WithRotation(degrees int) Freeform {
	f.RotationDeg = &degrees
	return f
}

// WithName sets the name of the freeform.
func (f Freeform) WithName(name string) Freeform {
	f.Name = name
	return f
}

// WithAltText sets the alternative text for accessibility.
func (f Freeform) WithAltText(text string) Freeform {
	f.AltText = text
	return f
}

// WithDecorative marks the freeform as decorative.
func (f Freeform) WithDecorative(enabled bool) Freeform {
	f.IsDecorative = enabled
	return f
}

// WithTextFrame applies custom text frame properties.
func (f Freeform) WithTextFrame(frame TextFrame) Freeform {
	f.TextFrame = &frame
	return f
}

// WithEffects enables visual effects on the freeform.
func (f Freeform) WithEffects(effects ShapeEffects) Freeform {
	f.Effects = &effects
	return f
}

// ToShape converts the freeform to a Shape for slide addition.
// Note: This creates a basic shape - for full freeform support, use the dedicated renderer.
func (f Freeform) ToShape() Shape {
	// Calculate bounds from points
	if len(f.Points) == 0 {
		return Shape{}
	}

	minX := f.Points[0].X
	minY := f.Points[0].Y
	maxX := f.Points[0].X
	maxY := f.Points[0].Y

	for i := 1; i < len(f.Points); i++ {
		if f.Points[i].X < minX {
			minX = f.Points[i].X
		}
		if f.Points[i].Y < minY {
			minY = f.Points[i].Y
		}
		if f.Points[i].X > maxX {
			maxX = f.Points[i].X
		}
		if f.Points[i].Y > maxY {
			maxY = f.Points[i].Y
		}
	}

	s := NewShape(ShapeTypeRectangle, minX, minY, maxX-minX, maxY-minY)
	s.Name = f.Name
	s.AltText = f.AltText
	s.IsDecorative = f.IsDecorative
	s.Text = f.Text
	s.RotationDeg = f.RotationDeg
	s.Fill = f.Fill
	s.Line = f.Line
	s.GradientFill = f.GradientFill
	s.RichFill = f.RichFill
	s.RichLine = f.RichLine
	s.RichShadow = f.RichShadow
	s.TextFrame = f.TextFrame
	s.Effects = f.Effects

	return s
}
