// Package ink models digital ink annotations: the pen strokes PowerPoint
// records when someone draws on a slide with a stylus or the Draw tab.
//
// One Annotation becomes an InkML part (`ppt/ink/inkN.xml`) plus a
// `<p:contentPart>` reference in the slide's shape tree.
package ink

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	// emuPerHiMetric converts EMU (914400 per inch) to the HiMetric units
	// (1/100 mm, 360000 per inch) InkML traces are written in.
	emuPerHiMetric = 360

	defaultPenWidthEmu = 28575 // 2.25pt, PowerPoint's default pen
	fullyOpaque        = 0
	highlighterAlpha   = 128
	maxTransparency    = 255
)

// PenTip is the shape of the pen nib.
type PenTip int

const (
	// TipEllipse is a round nib, PowerPoint's default.
	TipEllipse PenTip = iota
	// TipRectangle is a chisel nib, used by the highlighter.
	TipRectangle
)

// XMLValue returns the InkML brush property value for the tip.
func (t PenTip) XMLValue() string {
	switch t {
	case TipRectangle:
		return "rectangle"
	case TipEllipse:
		return "ellipse"
	default:
		return "ellipse"
	}
}

// Pen describes the brush one stroke is drawn with.
//
//nolint:govet // Preserve public field order for source compatibility with positional literals.
type Pen struct {
	// Color is an RRGGBB hex string without the leading '#'.
	Color string
	// WidthEmu and HeightEmu are the nib size. HeightEmu falls back to
	// WidthEmu when it is zero.
	WidthEmu  int64
	HeightEmu int64
	Tip       PenTip
	// Transparency is 0 (opaque) to 255 (invisible).
	Transparency uint8
	// IgnorePressure keeps the stroke width constant along the trace.
	IgnorePressure bool
}

// NewPen returns an opaque round pen of the given color and width.
func NewPen(color string, width styling.Length) Pen {
	w := width.Emu()
	if w <= 0 {
		w = defaultPenWidthEmu
	}
	return Pen{
		Color:          normalizeColor(color),
		WidthEmu:       w,
		HeightEmu:      w,
		Tip:            TipEllipse,
		Transparency:   fullyOpaque,
		IgnorePressure: true,
	}
}

// RedPen returns PowerPoint's default red pen.
func RedPen() Pen { return NewPen("FF0000", styling.Emu(defaultPenWidthEmu)) }

// BluePen returns PowerPoint's default blue pen.
func BluePen() Pen { return NewPen("0000FF", styling.Emu(defaultPenWidthEmu)) }

// BlackPen returns PowerPoint's default black pen.
func BlackPen() Pen { return NewPen("000000", styling.Emu(defaultPenWidthEmu)) }

// Highlighter returns a wide, semi-transparent chisel-tip yellow pen.
func Highlighter() Pen {
	const highlighterWidthEmu = 152400 // 12pt
	pen := NewPen("FFFF00", styling.Emu(highlighterWidthEmu))
	pen.Tip = TipRectangle
	pen.Transparency = highlighterAlpha
	return pen
}

// WithTip sets the nib shape.
func (p Pen) WithTip(tip PenTip) Pen {
	p.Tip = tip
	return p
}

// WithTransparency sets the stroke transparency, 0 opaque to 255 invisible.
func (p Pen) WithTransparency(transparency uint8) Pen {
	p.Transparency = transparency
	return p
}

// WithOpacity sets transparency from an opacity fraction in [0,1].
func (p Pen) WithOpacity(opacity float64) Pen {
	switch {
	case opacity <= 0:
		p.Transparency = maxTransparency
	case opacity >= 1:
		p.Transparency = fullyOpaque
	default:
		// Both guards above are exclusive, so the scaled value lands strictly
		// inside 0..255 and the rounding cannot push it out.
		p.Transparency = uint8(math.Round((1 - opacity) * maxTransparency))
	}
	return p
}

// WithSize sets the nib width and height.
func (p Pen) WithSize(width, height styling.Length) Pen {
	p.WidthEmu = width.Emu()
	p.HeightEmu = height.Emu()
	return p
}

// Point is one sampled position along a stroke, in slide coordinates.
type Point struct {
	X styling.Length
	Y styling.Length
}

// Stroke is a single continuous pen trace.
type Stroke struct {
	Pen    Pen
	Points []Point
}

// NewStroke starts an empty stroke drawn with pen.
func NewStroke(pen Pen) Stroke {
	return Stroke{Pen: pen}
}

// AddPoint appends one sampled position.
func (s Stroke) AddPoint(x, y styling.Length) Stroke {
	s.Points = append(s.Points, Point{X: x, Y: y})
	return s
}

// AddPoints appends several sampled positions.
func (s Stroke) AddPoints(points []Point) Stroke {
	s.Points = append(s.Points, points...)
	return s
}

// Len returns the number of sampled points.
func (s Stroke) Len() int { return len(s.Points) }

// IsEmpty reports whether the stroke has nothing to draw. A single point is
// still nothing to draw: InkML needs at least two positions for a trace.
func (s Stroke) IsEmpty() bool { return len(s.Points) < 2 }

// Annotation is the set of strokes stored in one InkML part, together with the
// frame the ink is placed in on the slide.
//
//nolint:govet // Preserve public field order for source compatibility with positional literals.
type Annotation struct {
	Strokes []Stroke
	// X, Y, CX and CY position the ink frame on the slide. When CX or CY is
	// zero, Bounds() is used instead.
	X  styling.Length
	Y  styling.Length
	CX styling.Length
	CY styling.Length
	// Name is the shape name PowerPoint shows in the selection pane.
	Name string
}

// New returns an empty annotation.
func New() *Annotation {
	return &Annotation{Name: "Ink"}
}

// AddStroke appends one stroke.
func (a *Annotation) AddStroke(stroke Stroke) *Annotation {
	a.Strokes = append(a.Strokes, stroke)
	return a
}

// WithFrame sets the ink frame explicitly instead of deriving it from Bounds.
func (a *Annotation) WithFrame(x, y, cx, cy styling.Length) *Annotation {
	a.X, a.Y, a.CX, a.CY = x, y, cx, cy
	return a
}

// WithName sets the shape name.
func (a *Annotation) WithName(name string) *Annotation {
	a.Name = name
	return a
}

// Len returns the number of strokes.
func (a *Annotation) Len() int { return len(a.Strokes) }

// IsEmpty reports whether the annotation would draw nothing.
func (a *Annotation) IsEmpty() bool {
	if a == nil {
		return true
	}
	for _, s := range a.Strokes {
		if !s.IsEmpty() {
			return false
		}
	}
	return true
}

// Clear drops every stroke.
func (a *Annotation) Clear() {
	a.Strokes = nil
}

// Bounds returns the smallest box containing every stroke point as
// x, y, cx, cy in EMU. It returns zeros for an empty annotation.
func (a *Annotation) Bounds() (int64, int64, int64, int64) {
	first := true
	var minX, minY, maxX, maxY int64
	for _, stroke := range a.Strokes {
		for _, p := range stroke.Points {
			px, py := p.X.Emu(), p.Y.Emu()
			if first {
				minX, minY, maxX, maxY = px, py, px, py
				first = false
				continue
			}
			minX = min(minX, px)
			minY = min(minY, py)
			maxX = max(maxX, px)
			maxY = max(maxY, py)
		}
	}
	if first {
		return 0, 0, 0, 0
	}
	return minX, minY, maxX - minX, maxY - minY
}

// frame returns the placement used for the content part as x, y, cx, cy,
// falling back to the stroke bounds when no explicit frame was set.
func (a *Annotation) frame() (int64, int64, int64, int64) {
	if a.CX.Emu() > 0 && a.CY.Emu() > 0 {
		return a.X.Emu(), a.Y.Emu(), a.CX.Emu(), a.CY.Emu()
	}
	return a.Bounds()
}

func normalizeColor(color string) string {
	c := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(color), "#"))
	if c == "" {
		return "000000"
	}
	return strings.ToUpper(c)
}

func emuToHiMetric(emu int64) int64 {
	return emu / emuPerHiMetric
}

func itoa(v int64) string {
	return strconv.FormatInt(v, 10)
}

// PartPath returns the package path of the nth ink part, 1-based.
func PartPath(index int) string {
	return fmt.Sprintf("ppt/ink/ink%d.xml", index)
}

// RelationshipTarget returns the target a slide uses to reach the nth ink part.
func RelationshipTarget(index int) string {
	return fmt.Sprintf("../ink/ink%d.xml", index)
}
