package shapes

import (
	"errors"
	"fmt"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// LineDashStyle represents predefined dash styles for shape lines.
type LineDashStyle string

const (
	// LineDashStyleSolid is a solid line.
	LineDashStyleSolid LineDashStyle = "solid"
	// LineDashStyleDash is a dashed line.
	LineDashStyleDash LineDashStyle = "dash"
	// LineDashStyleDot is a dotted line.
	LineDashStyleDot LineDashStyle = "dot"
	// LineDashStyleDashDot is a dash-dot line.
	LineDashStyleDashDot LineDashStyle = "dashDot"
	// LineDashStyleDashDotDot is a dash-dot-dot line.
	LineDashStyleDashDotDot LineDashStyle = "dashDotDot"
	// LineDashStyleLongDash is a long dashed line.
	LineDashStyleLongDash LineDashStyle = "lgDash"
	// LineDashStyleLongDashDot is a long dash-dot line.
	LineDashStyleLongDashDot LineDashStyle = "lgDashDot"
	// LineDashStyleSystemDash is a system dash line.
	LineDashStyleSystemDash LineDashStyle = "sysDash"
	// LineDashStyleSystemDot is a system dot line.
	LineDashStyleSystemDot LineDashStyle = "sysDot"
	// LineDashStyleSystemDashDot is a system dash-dot line.
	LineDashStyleSystemDashDot LineDashStyle = "sysDashDot"
)

// IsValidLineDashStyle returns true if the dash style is valid.
func IsValidLineDashStyle(s LineDashStyle) bool {
	switch s {
	case LineDashStyleSolid, LineDashStyleDash, LineDashStyleDot, LineDashStyleDashDot,
		LineDashStyleDashDotDot, LineDashStyleLongDash, LineDashStyleLongDashDot,
		LineDashStyleSystemDash, LineDashStyleSystemDot, LineDashStyleSystemDashDot:
		return true
	}
	return false
}

// NormalizeLineDashStyle normalizes a line dash style string.
func NormalizeLineDashStyle(s string) LineDashStyle {
	ls := LineDashStyle(s)
	if IsValidLineDashStyle(ls) {
		return ls
	}
	return LineDashStyleSolid
}

// LineCapStyle represents the cap style for line ends.
type LineCapStyle string

const (
	// LineCapStyleFlat is a flat line cap.
	LineCapStyleFlat LineCapStyle = "flat"
	// LineCapStyleRound is a round line cap.
	LineCapStyleRound LineCapStyle = "rnd"
	// LineCapStyleSquare is a square line cap.
	LineCapStyleSquare LineCapStyle = "sq"
)

// IsValidLineCapStyle returns true if the cap style is valid.
func IsValidLineCapStyle(s LineCapStyle) bool {
	switch s {
	case LineCapStyleFlat, LineCapStyleRound, LineCapStyleSquare:
		return true
	}
	return false
}

// NormalizeLineCapStyle normalizes a line cap style string.
func NormalizeLineCapStyle(s string) LineCapStyle {
	cs := LineCapStyle(s)
	if IsValidLineCapStyle(cs) {
		return cs
	}
	return LineCapStyleFlat
}

// LineJoinStyle represents the join style for line corners.
type LineJoinStyle string

const (
	// LineJoinStyleRound is a round line join.
	LineJoinStyleRound LineJoinStyle = "round"
	// LineJoinStyleBevel is a beveled line join.
	LineJoinStyleBevel LineJoinStyle = "bevel"
	// LineJoinStyleMiter is a mitered line join.
	LineJoinStyleMiter LineJoinStyle = "miter"
)

// IsValidLineJoinStyle returns true if the join style is valid.
func IsValidLineJoinStyle(s LineJoinStyle) bool {
	switch s {
	case LineJoinStyleRound, LineJoinStyleBevel, LineJoinStyleMiter:
		return true
	}
	return false
}

// NormalizeLineJoinStyle normalizes a line join style string.
func NormalizeLineJoinStyle(s string) LineJoinStyle {
	js := LineJoinStyle(s)
	if IsValidLineJoinStyle(js) {
		return js
	}
	return LineJoinStyleRound
}

// RichShapeLine provides detailed control over shape line properties.
type RichShapeLine struct {
	Color        string
	Width        styling.Length
	DashStyle    LineDashStyle
	CapStyle     LineCapStyle
	JoinStyle    LineJoinStyle
	Transparency float64 // 0.0 = opaque, 1.0 = fully transparent
}

// NewRichShapeLine creates a new line style with the specified color and width.
func NewRichShapeLine(color string, width styling.Length) *RichShapeLine {
	return &RichShapeLine{
		Color:        common.NormalizeHexColor(color),
		Width:        width,
		DashStyle:    LineDashStyleSolid,
		CapStyle:     LineCapStyleFlat,
		JoinStyle:    LineJoinStyleRound,
		Transparency: 0.0,
	}
}

// WithColor sets the line color.
func (l *RichShapeLine) WithColor(color string) *RichShapeLine {
	l.Color = common.NormalizeHexColor(color)
	return l
}

// WithWidth sets the line width.
func (l *RichShapeLine) WithWidth(width styling.Length) *RichShapeLine {
	l.Width = width
	return l
}

// WithDashStyle sets the line dash style.
func (l *RichShapeLine) WithDashStyle(style LineDashStyle) *RichShapeLine {
	l.DashStyle = style
	return l
}

// WithCapStyle sets the line cap style.
func (l *RichShapeLine) WithCapStyle(style LineCapStyle) *RichShapeLine {
	l.CapStyle = style
	return l
}

// WithJoinStyle sets the line join style.
func (l *RichShapeLine) WithJoinStyle(style LineJoinStyle) *RichShapeLine {
	l.JoinStyle = style
	return l
}

// WithTransparency sets the line transparency (0.0 to 1.0).
func (l *RichShapeLine) WithTransparency(transparency float64) *RichShapeLine {
	l.Transparency = transparency
	return l
}

// Validate checks the line style for validity.
func (l *RichShapeLine) Validate() error {
	if l == nil {
		return nil
	}

	if !common.IsHexColor(l.Color) {
		return fmt.Errorf("invalid line color: %s", l.Color)
	}

	if l.Width < 0 {
		return errors.New("line width cannot be negative")
	}

	if !IsValidLineDashStyle(l.DashStyle) {
		return fmt.Errorf("invalid dash style: %s", l.DashStyle)
	}

	if !IsValidLineCapStyle(l.CapStyle) {
		return fmt.Errorf("invalid cap style: %s", l.CapStyle)
	}

	if !IsValidLineJoinStyle(l.JoinStyle) {
		return fmt.Errorf("invalid join style: %s", l.JoinStyle)
	}

	if l.Transparency < 0 || l.Transparency > 1 {
		return errors.New("transparency must be between 0.0 and 1.0")
	}

	return nil
}
