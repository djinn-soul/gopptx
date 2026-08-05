package export

import (
	"math"

	"github.com/signintech/gopdf"
)

const (
	// roundRectCornerPoints is how many intermediate points gopdf interpolates
	// per corner. gopdf falls back to a plain polygon when this is <= 0, which
	// is why every rounded box used to come out square.
	roundRectCornerPoints = 8

	// minRoundRectRadius is the smallest corner worth drawing. Below it the
	// rounding is invisible and a plain rectangle is cheaper.
	minRoundRectRadius = 0.35
)

// drawRoundedRect draws PowerPoint's roundRect preset.
//
// gopdf's Rectangle rejects a radius larger than either side and silently
// degrades to a square polygon when given no corner points, so both are
// checked here rather than left to it. A failure falls back to a plain
// rectangle so the shape is still drawn.
func drawRoundedRect(pdf *gopdf.GoPdf, x, y, w, h float64, style string) {
	radius := roundRectRadius(w, h)
	if radius < minRoundRectRadius {
		pdf.RectFromUpperLeftWithStyle(x, y, w, h, style)
		return
	}
	if err := pdf.Rectangle(x, y, x+w, y+h, style, radius, roundRectCornerPoints); err != nil {
		pdf.RectFromUpperLeftWithStyle(x, y, w, h, style)
	}
}

// roundRectRadius is the corner radius for a box, clamped so it can never
// exceed either side — gopdf errors out rather than clamping for us.
func roundRectRadius(w, h float64) float64 {
	shorter := math.Min(w, h)
	if shorter <= 0 {
		return 0
	}
	radius := shorter * defaultRadiusFactor
	// A radius equal to half the side is a full semicircle; beyond that the
	// corners would overlap.
	return math.Min(radius, shorter/2)
}
