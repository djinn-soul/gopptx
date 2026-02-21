package shapes

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// NewPieSlice creates a pie slice shape using the "pie" preset with adjustment points.
// startAngle and endAngle are in degrees (0 is at 3 o'clock, clockwise).
func NewPieSlice(x, y, cx, cy styling.Length, startAngle, endAngle float64) Shape {
	// OOXML angles are in 60,000ths of a degree
	s := NewShape(ShapeTypePie, x, y, cx, cy)
	s = s.WithAdjustmentValue("adj1", int(startAngle*60000))
	s = s.WithAdjustmentValue("adj2", int(endAngle*60000))
	return s
}
