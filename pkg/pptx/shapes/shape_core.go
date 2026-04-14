package shapes

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

// ShapeAdjustment represents one geometry adjustment point (<a:gd>) entry.
type ShapeAdjustment struct {
	Name    string
	Formula string
}

// Shape is one auto shape.
type Shape struct {
	Type           string
	X              styling.Length
	Y              styling.Length
	CX             styling.Length
	CY             styling.Length
	Fill           *ShapeFill
	Line           *ShapeLine
	GradientFill   *ShapeGradientFill
	Text           string
	RotationDeg    *int
	Hyperlink      *action.Hyperlink // Legacy: mapped to ClickAction
	ClickAction    *action.Hyperlink
	HoverAction    *action.Hyperlink
	AltText        string
	IsDecorative   bool
	TextFrame      *TextFrame
	TextParagraphs []text.Paragraph
	Name           string
	Adjustments    []ShapeAdjustment
	Effects        *ShapeEffects
	RichFill       *RichShapeFill
	RichLine       *RichShapeLine
	RichShadow     *RichShapeShadow
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
	s.Hyperlink = &link
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

// cloneTextFrame returns a copy of the shape's TextFrame (or a fresh one),
// then points s.TextFrame at the copy so mutations never alias a shared instance.
func (s *Shape) cloneTextFrame() TextFrame {
	if s.TextFrame != nil {
		return *s.TextFrame
	}
	return NewTextFrame()
}

// WithTextMargins sets EMU margins for the shape text frame.
func (s Shape) WithTextMargins(left, top, right, bottom styling.Length) Shape {
	tf := s.cloneTextFrame()
	tf.MarginLeft = left
	tf.MarginTop = top
	tf.MarginRight = right
	tf.MarginBottom = bottom
	s.TextFrame = &tf
	return s
}

// WithVerticalAnchor sets the vertical alignment of text in the shape.
func (s Shape) WithVerticalAnchor(anchor TextFrameAnchor) Shape {
	tf := s.cloneTextFrame()
	tf.Anchor = anchor
	s.TextFrame = &tf
	return s
}

// WithTextWrap sets the text wrapping mode.
func (s Shape) WithTextWrap(wrap TextFrameWrap) Shape {
	tf := s.cloneTextFrame()
	tf.Wrap = wrap
	s.TextFrame = &tf
	return s
}

// WithAutoFit sets the auto-fit behavior for text within the shape.
func (s Shape) WithAutoFit(autoFit TextFrameAutoFit) Shape {
	tf := s.cloneTextFrame()
	tf.AutoFit = autoFit
	s.TextFrame = &tf
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
	s.Fill = nil
	s.GradientFill = nil
	return s
}

// WithRichLine applies a rich line style to a shape.
func (s Shape) WithRichLine(line *RichShapeLine) Shape {
	s.RichLine = line
	s.Line = nil
	return s
}

// WithRichShadow applies a rich shadow effect to a shape.
func (s Shape) WithRichShadow(shadow *RichShapeShadow) Shape {
	s.RichShadow = shadow
	if s.Effects == nil {
		s.Effects = &ShapeEffects{}
	}
	s.Effects.Shadow = shadow != nil
	return s
}
