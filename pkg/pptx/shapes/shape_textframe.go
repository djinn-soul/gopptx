package shapes

import (
	"strings"

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
	Orientation  string
	Columns      int
	RotationDeg  *float64
}

const defaultTextMarginInches = 0.05

// NewTextFrame creates a text frame with default margins (0.05 inches).
func NewTextFrame() TextFrame {
	return TextFrame{
		MarginLeft:   styling.Inches(defaultTextMarginInches),
		MarginRight:  styling.Inches(defaultTextMarginInches),
		MarginTop:    styling.Inches(defaultTextMarginInches),
		MarginBottom: styling.Inches(defaultTextMarginInches),
		Anchor:       TextAnchorMiddle,
		Wrap:         TextWrapSquare,
		AutoFit:      TextAutoFitShape,
	}
}

// WithRotation sets text-frame rotation in degrees.
func (t TextFrame) WithRotation(degrees float64) TextFrame {
	value := degrees
	t.RotationDeg = &value
	return t
}

// WithOrientation sets the OOXML text orientation token (for example, "vert270").
func (t TextFrame) WithOrientation(orientation string) TextFrame {
	t.Orientation = strings.TrimSpace(orientation)
	return t
}

// WithColumns sets the number of text columns in the frame.
func (t TextFrame) WithColumns(columns int) TextFrame {
	t.Columns = columns
	return t
}
