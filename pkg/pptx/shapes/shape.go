package shapes

import (
	"errors"
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// TextFrameAnchor specifies the vertical alignment of text within its shape.
type TextFrameAnchor string

const (
	TextAnchorTop    TextFrameAnchor = "t"
	TextAnchorMiddle TextFrameAnchor = "ctr"
	TextAnchorBottom TextFrameAnchor = "b"
)

// TextFrameWrap specifies how text wraps within the shape's text frame.
type TextFrameWrap string

const (
	TextWrapNone   TextFrameWrap = "none"
	TextWrapSquare TextFrameWrap = "square"
)

// TextFrameAutoFit specifies how text is automatically resized or how the shape is resized.
type TextFrameAutoFit string

const (
	TextAutoFitNone   TextFrameAutoFit = "none"
	TextAutoFitShape  TextFrameAutoFit = "spAutoFit"
	TextAutoFitNormal TextFrameAutoFit = "normAutoFit"
)

// TextFrame configures the text layout within a shape.
type TextFrame struct {
	MarginLeft   styling.Length // EMU
	MarginRight  styling.Length
	MarginTop    styling.Length
	MarginBottom styling.Length
	Anchor       TextFrameAnchor
	Wrap         TextFrameWrap
	AutoFit      TextFrameAutoFit
}

const (
	defaultTextMarginInches = 0.05
	minGradientStops        = 2
)

// NewTextFrame creates a text frame with default margins (0.05 inches).
func NewTextFrame() TextFrame {
	return TextFrame{
		MarginLeft:   styling.Inches(defaultTextMarginInches),
		MarginRight:  styling.Inches(defaultTextMarginInches),
		MarginTop:    styling.Inches(defaultTextMarginInches),
		MarginBottom: styling.Inches(defaultTextMarginInches),

		Anchor:  TextAnchorMiddle,
		Wrap:    TextWrapSquare,
		AutoFit: TextAutoFitShape,
	}
}

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
	// Validate stop constraints
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

// ShapeAdjustment represents one geometry adjustment point (<a:gd>) entry.
type ShapeAdjustment struct {
	Name    string
	Formula string
}

// Shape is one auto shape.
type Shape struct {
	Type         string
	X            styling.Length
	Y            styling.Length
	CX           styling.Length
	CY           styling.Length
	Fill         *ShapeFill
	Line         *ShapeLine
	GradientFill *ShapeGradientFill
	Text         string
	RotationDeg  *int
	Hyperlink    *action.Hyperlink // Legacy: mapped to ClickAction
	ClickAction  *action.Hyperlink
	HoverAction  *action.Hyperlink
	AltText      string
	IsDecorative bool
	TextFrame    *TextFrame
	Name         string
	Adjustments  []ShapeAdjustment
	Effects      *ShapeEffects
	// Rich formatting properties (new)
	RichFill   *RichShapeFill
	RichLine   *RichShapeLine
	RichShadow *RichShapeShadow
}

// NewShape creates one shape.
func NewShape(shapeType string, x, y, cx, cy styling.Length) Shape {
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

// WithAdjustment appends one geometry adjustment point.
func (s Shape) WithAdjustment(name, formula string) Shape {
	s.Adjustments = append(s.Adjustments, ShapeAdjustment{
		Name:    strings.TrimSpace(name),
		Formula: strings.TrimSpace(formula),
	})
	return s
}

// WithAdjustmentValue appends one "val" adjustment helper entry.
func (s Shape) WithAdjustmentValue(name string, value int) Shape {
	return s.WithAdjustment(name, fmt.Sprintf("val %d", value))
}

// WithClickAction adds a click behavior to the shape.
func (s Shape) WithClickAction(link action.Hyperlink) Shape {
	s.ClickAction = &link
	s.Hyperlink = &link // Keep legacy field in sync
	return s
}

// WithHoverAction adds a hover behavior to the shape.
func (s Shape) WithHoverAction(link action.Hyperlink) Shape {
	s.HoverAction = &link
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
	return s.WithClickAction(hyperlink)
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

// WithTextFrame applies custom text frame properties to a shape.
func (s Shape) WithTextFrame(frame TextFrame) Shape {
	s.TextFrame = &frame
	return s
}

// WithTextMargins sets EMU margins for the shape text frame.
func (s Shape) WithTextMargins(left, top, right, bottom styling.Length) Shape {
	if s.TextFrame == nil {
		tf := NewTextFrame()
		s.TextFrame = &tf
	}
	s.TextFrame.MarginLeft = left
	s.TextFrame.MarginTop = top
	s.TextFrame.MarginRight = right
	s.TextFrame.MarginBottom = bottom
	return s
}

// WithVerticalAnchor sets the vertical alignment of text in the shape.
func (s Shape) WithVerticalAnchor(anchor TextFrameAnchor) Shape {
	if s.TextFrame == nil {
		tf := NewTextFrame()
		s.TextFrame = &tf
	}
	s.TextFrame.Anchor = anchor
	return s
}

// WithTextWrap sets the text wrapping mode.
func (s Shape) WithTextWrap(wrap TextFrameWrap) Shape {
	if s.TextFrame == nil {
		tf := NewTextFrame()
		s.TextFrame = &tf
	}
	s.TextFrame.Wrap = wrap
	return s
}

// WithAutoFit sets the auto-fit behavior for text within the shape.
func (s Shape) WithAutoFit(autoFit TextFrameAutoFit) Shape {
	if s.TextFrame == nil {
		tf := NewTextFrame()
		s.TextFrame = &tf
	}
	s.TextFrame.AutoFit = autoFit
	return s
}

// WithName sets the name of the shape for Morph transitions and selection pane.
func (s Shape) WithName(name string) Shape {
	s.Name = name
	return s
}

// WithRichFill applies a rich fill (solid, gradient, pattern, or no-fill) to a shape.
func (s Shape) WithRichFill(fill *RichShapeFill) Shape {
	s.RichFill = fill
	// Clear legacy fill when using rich fill
	s.Fill = nil
	s.GradientFill = nil
	return s
}

// WithRichLine applies a rich line style to a shape.
func (s Shape) WithRichLine(line *RichShapeLine) Shape {
	s.RichLine = line
	// Clear legacy line when using rich line
	s.Line = nil
	return s
}

// WithRichShadow applies a rich shadow effect to a shape.
func (s Shape) WithRichShadow(shadow *RichShapeShadow) Shape {
	s.RichShadow = shadow
	// Update legacy effects if present
	if s.Effects == nil {
		s.Effects = &ShapeEffects{}
	}
	s.Effects.Shadow = shadow != nil
	return s
}

// Validate checks for validity of shape parameters.
func (s Shape) Validate(slideIndex, shapeIndex int) error {
	if !s.IsDecorative && len(s.AltText) > common.MaxAltTextLength {
		return fmt.Errorf(
			"shape %d on slide %d alt text exceeds %d characters",
			shapeIndex,
			slideIndex,
			common.MaxAltTextLength,
		)
	}

	if err := s.validateShapeBounds(slideIndex, shapeIndex); err != nil {
		return err
	}
	if !IsShapeType(s.Type) {
		return fmt.Errorf("shape %d type %q is invalid on slide %d", shapeIndex, s.Type, slideIndex)
	}

	if err := s.validateFills(slideIndex, shapeIndex); err != nil {
		return err
	}
	if err := s.validateLinesAndRotation(slideIndex, shapeIndex); err != nil {
		return err
	}
	return s.validateActions(slideIndex, shapeIndex)
}

func (s Shape) validateShapeBounds(slideIndex, shapeIndex int) error {
	if s.X < 0 || s.Y < 0 {
		return fmt.Errorf("shape %d on slide %d position cannot be negative", shapeIndex, slideIndex)
	}
	if s.CX <= 0 || s.CY <= 0 {
		return fmt.Errorf("shape %d on slide %d size must be > 0", shapeIndex, slideIndex)
	}
	return nil
}

func (s Shape) validateFills(slideIndex, shapeIndex int) error {
	// Check for conflicts between legacy and rich fill
	if s.RichFill != nil && (s.Fill != nil || s.GradientFill != nil) {
		return fmt.Errorf("shape %d (type %q) on slide %d cannot set both rich fill and legacy fill",
			shapeIndex, s.Type, slideIndex)
	}
	if s.Fill != nil && s.GradientFill != nil {
		return fmt.Errorf("shape %d (type %q) on slide %d cannot set both solid and gradient fill",
			shapeIndex, s.Type, slideIndex)
	}
	if s.Fill != nil {
		if err := s.Fill.Validate(); err != nil {
			return fmt.Errorf("shape %d on slide %d has invalid fill: %w", shapeIndex, slideIndex, err)
		}
	}
	if s.GradientFill != nil {
		if err := s.GradientFill.Validate(); err != nil {
			return fmt.Errorf("shape %d on slide %d has invalid gradient fill: %w", shapeIndex, slideIndex, err)
		}
	}
	if s.RichFill != nil {
		if err := s.RichFill.Validate(); err != nil {
			return fmt.Errorf("shape %d on slide %d has invalid rich fill: %w", shapeIndex, slideIndex, err)
		}
	}
	return nil
}

func (s Shape) validateLinesAndRotation(slideIndex, shapeIndex int) error {
	// Check for conflicts between legacy and rich line
	if s.RichLine != nil && s.Line != nil {
		return fmt.Errorf("shape %d (type %q) on slide %d cannot set both rich line and legacy line",
			shapeIndex, s.Type, slideIndex)
	}
	if s.Line != nil {
		if err := s.Line.Validate(); err != nil {
			return fmt.Errorf("shape %d on slide %d has invalid line: %w", shapeIndex, slideIndex, err)
		}
	}
	if s.RichLine != nil {
		if err := s.RichLine.Validate(); err != nil {
			return fmt.Errorf("shape %d on slide %d has invalid rich line: %w", shapeIndex, slideIndex, err)
		}
	}
	if s.RichShadow != nil {
		if err := s.RichShadow.Validate(); err != nil {
			return fmt.Errorf("shape %d on slide %d has invalid rich shadow: %w", shapeIndex, slideIndex, err)
		}
	}
	if s.RotationDeg != nil {
		if *s.RotationDeg < -360 || *s.RotationDeg > 360 {
			return fmt.Errorf("shape %d on slide %d rotation must be in [-360,360]", shapeIndex, slideIndex)
		}
	}
	return nil
}

func (s Shape) validateActions(slideIndex, shapeIndex int) error {
	if s.ClickAction != nil {
		if err := s.ClickAction.Validate(); err != nil {
			return fmt.Errorf("shape %d on slide %d has invalid click action: %w", shapeIndex, slideIndex, err)
		}
	} else if s.Hyperlink != nil {
		if err := s.Hyperlink.Validate(); err != nil {
			return fmt.Errorf("shape %d on slide %d has invalid hyperlink: %w", shapeIndex, slideIndex, err)
		}
	}
	if s.HoverAction != nil {
		if err := s.HoverAction.Validate(); err != nil {
			return fmt.Errorf("shape %d on slide %d has invalid hover action: %w", shapeIndex, slideIndex, err)
		}
	}
	return nil
}
