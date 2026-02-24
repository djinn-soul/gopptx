package shapes

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const ooxmlAngleUnitsPerDegree = 60000

// NewArc creates an arc shape using the "arc" preset with adjustment points.
// startAngle and endAngle are in degrees (0 is at 3 o'clock, clockwise).
func NewArc(x, y, cx, cy styling.Length, startAngle, endAngle float64) Shape {
	// OOXML angles are in 60,000ths of a degree
	s := NewShape(ShapeTypeArc, x, y, cx, cy)
	s = s.WithAdjustmentValue("adj1", int(startAngle*ooxmlAngleUnitsPerDegree))
	s = s.WithAdjustmentValue("adj2", int(endAngle*ooxmlAngleUnitsPerDegree))
	return s
}
