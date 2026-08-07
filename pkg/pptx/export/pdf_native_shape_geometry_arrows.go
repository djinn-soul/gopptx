//nolint:mnd // Arrow and callout outlines use fixed proportions from the DrawingML geometry.
package export

import (
	"math"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// The rest of the arrow family — the bent, quad and striped arrows, and the
// arrow callouts — plus the block callouts and the tabbed action buttons.

// drawPDFArrowGeometry draws one of the presets below and reports whether it
// recognised the type.
func drawPDFArrowGeometry(
	pdf *gopdf.GoPdf,
	fl flipState,
	shapeType string,
	x, y, w, h float64,
	style string,
) bool {
	switch shapeType {
	case shapes.ShapeTypeStripedRightArrow:
		drawStripedRightArrow(pdf, fl, x, y, w, h, style)
	case shapes.ShapeTypeQuadArrow:
		fl.polygon(pdf, quadArrowPoints(x, y, w, h), style)
	case shapes.ShapeTypeLeftRightUpArrow:
		fl.polygon(pdf, leftRightUpArrowPoints(x, y, w, h), style)
	case shapes.ShapeTypeBentArrow:
		fl.polygon(pdf, bentArrowPoints(x, y, w, h), style)
	case shapes.ShapeTypeBentUpArrow:
		fl.polygon(pdf, bentUpArrowPoints(x, y, w, h), style)
	case shapes.ShapeTypeLeftUpArrow:
		fl.polygon(pdf, leftUpArrowPoints(x, y, w, h), style)
	case shapes.ShapeTypeUturnArrow:
		fl.polygon(pdf, uturnArrowPoints(x, y, w, h), style)
	case shapes.ShapeTypeSwooshArrow:
		fl.polygon(pdf, swooshArrowPoints(x, y, w, h), style)
	case shapes.ShapeTypeRightArrowCallout:
		fl.polygon(pdf, arrowCalloutPoints(x, y, w, h, arrowRight), style)
	case shapes.ShapeTypeLeftArrowCallout:
		fl.polygon(pdf, arrowCalloutPoints(x, y, w, h, arrowLeft), style)
	case shapes.ShapeTypeUpArrowCallout:
		fl.polygon(pdf, arrowCalloutPoints(x, y, w, h, arrowUp), style)
	case shapes.ShapeTypeDownArrowCallout:
		fl.polygon(pdf, arrowCalloutPoints(x, y, w, h, arrowDown), style)
	case shapes.ShapeTypeLeftRightArrowCallout:
		fl.polygon(pdf, leftRightArrowCalloutPoints(x, y, w, h), style)
	case shapes.ShapeTypeUpDownArrowCallout:
		fl.polygon(pdf, upDownArrowCalloutPoints(x, y, w, h), style)
	case shapes.ShapeTypeQuadArrowCallout:
		fl.polygon(pdf, quadArrowCalloutPoints(x, y, w, h), style)
	default:
		return false
	}
	return true
}

const (
	// arrowCalloutBodyRatio is how much of the box the callout's text body takes,
	// leaving the rest for the arrow.
	arrowCalloutBodyRatio = 0.6
	// arrowCalloutShaftRatio is the shaft thickness of a callout's arrow.
	arrowCalloutShaftRatio = 0.25
	// arrowCalloutHeadRatio is the head width of a callout's arrow.
	arrowCalloutHeadRatio = 0.45
	// bentArrowShaftRatio is the shaft thickness of the bent arrows.
	bentArrowShaftRatio = 0.3
	// quadArrowShaftRatio is the shaft thickness of the quad arrow.
	quadArrowShaftRatio = 0.22
	// quadArrowHeadRatio is how far in from each edge a quad arrow's head starts.
	quadArrowHeadRatio = 0.28
	// stripedArrowStripes is how many stripes the striped arrow's tail carries.
	stripedArrowStripes = 3
	// stripedArrowTailRatio is how much of the width the striped tail takes.
	stripedArrowTailRatio = 0.25
	// uturnThicknessRatio is the shaft thickness of the U-turn arrow.
	uturnThicknessRatio = 0.28
)

// arrowDirection is which way one of the arrow callouts points.
type arrowDirection uint8

const (
	arrowRight arrowDirection = iota
	arrowLeft
	arrowUp
	arrowDown
)

// quadArrowPoints is the four-headed arrow: a cross with a point on each arm.
func quadArrowPoints(x, y, w, h float64) []gopdf.Point {
	shaftX := w * quadArrowShaftRatio / 2
	shaftY := h * quadArrowShaftRatio / 2
	headX := w * quadArrowHeadRatio
	headY := h * quadArrowHeadRatio
	cx, cy := x+w/2, y+h/2
	return []gopdf.Point{
		{X: cx, Y: y}, {X: cx + headX, Y: y + headY}, {X: cx + shaftX, Y: y + headY},
		{X: cx + shaftX, Y: cy - shaftY}, {X: x + w - headX, Y: cy - shaftY},
		{X: x + w - headX, Y: cy - headY}, {X: x + w, Y: cy},
		{X: x + w - headX, Y: cy + headY}, {X: x + w - headX, Y: cy + shaftY},
		{X: cx + shaftX, Y: cy + shaftY}, {X: cx + shaftX, Y: y + h - headY},
		{X: cx + headX, Y: y + h - headY}, {X: cx, Y: y + h},
		{X: cx - headX, Y: y + h - headY}, {X: cx - shaftX, Y: y + h - headY},
		{X: cx - shaftX, Y: cy + shaftY}, {X: x + headX, Y: cy + shaftY},
		{X: x + headX, Y: cy + headY}, {X: x, Y: cy},
		{X: x + headX, Y: cy - headY}, {X: x + headX, Y: cy - shaftY},
		{X: cx - shaftX, Y: cy - shaftY}, {X: cx - shaftX, Y: y + headY},
		{X: cx - headX, Y: y + headY},
	}
}

// leftRightUpArrowPoints is the three-headed arrow: left, right and up.
func leftRightUpArrowPoints(x, y, w, h float64) []gopdf.Point {
	shaftX := w * quadArrowShaftRatio / 2
	shaftY := h * quadArrowShaftRatio / 2
	headX := w * quadArrowHeadRatio
	headY := h * quadArrowHeadRatio
	cx := x + w/2
	base := y + h
	mid := base - shaftY*2
	return []gopdf.Point{
		{X: cx, Y: y}, {X: cx + headX, Y: y + headY}, {X: cx + shaftX, Y: y + headY},
		{X: cx + shaftX, Y: mid}, {X: x + w - headX, Y: mid},
		{X: x + w - headX, Y: mid - (headY - shaftY)}, {X: x + w, Y: (mid + base) / 2},
		{X: x + w - headX, Y: base}, {X: x + w - headX, Y: base - shaftY},
		{X: x + headX, Y: base - shaftY}, {X: x + headX, Y: base},
		{X: x, Y: (mid + base) / 2}, {X: x + headX, Y: mid - (headY - shaftY)},
		{X: x + headX, Y: mid}, {X: cx - shaftX, Y: mid},
		{X: cx - shaftX, Y: y + headY}, {X: cx - headX, Y: y + headY},
	}
}

// bentArrowPoints is the arrow that turns a right angle and points up.
func bentArrowPoints(x, y, w, h float64) []gopdf.Point {
	shaft := h * bentArrowShaftRatio
	head := math.Min(w, h) * 0.4
	return []gopdf.Point{
		{X: x, Y: y + h - shaft}, {X: x + w - head, Y: y + h - shaft},
		{X: x + w - head, Y: y + h - shaft - head/2}, {X: x + w, Y: y + h - shaft - head},
		{X: x + w - head, Y: y + h - shaft - head*1.5}, {X: x + w - head, Y: y + h - shaft*2 - head},
		{X: x + shaft, Y: y + h - shaft*2 - head}, {X: x + shaft, Y: y + h},
		{X: x, Y: y + h},
	}
}

// bentUpArrowPoints is the arrow that runs along the bottom then turns up.
func bentUpArrowPoints(x, y, w, h float64) []gopdf.Point {
	shaft := h * bentArrowShaftRatio
	head := math.Min(w, h) * 0.45
	stemX := x + w - head
	return []gopdf.Point{
		{X: x, Y: y + h - shaft}, {X: stemX - head/2, Y: y + h - shaft},
		{X: stemX - head/2, Y: y + head}, {X: stemX - head, Y: y + head},
		{X: stemX, Y: y}, {X: stemX + head, Y: y + head},
		{X: stemX + head/2, Y: y + head}, {X: stemX + head/2, Y: y + h},
		{X: x, Y: y + h},
	}
}

// leftUpArrowPoints is the two-headed corner arrow: left and up.
func leftUpArrowPoints(x, y, w, h float64) []gopdf.Point {
	shaft := math.Min(w, h) * bentArrowShaftRatio
	head := math.Min(w, h) * 0.45
	return []gopdf.Point{
		{X: x, Y: y + h - shaft/2 - head/2}, {X: x + head, Y: y + h - shaft/2 - head},
		{X: x + head, Y: y + h - shaft}, {X: x + w - shaft, Y: y + h - shaft},
		{X: x + w - shaft, Y: y + head}, {X: x + w - shaft - head/2, Y: y + head},
		{X: x + w - shaft/2, Y: y}, {X: x + w, Y: y + head},
		{X: x + w - shaft/2 + head/2, Y: y + head}, {X: x + w - shaft/2 + head/2, Y: y + h},
		{X: x + head, Y: y + h}, {X: x + head, Y: y + h - shaft/2 + head/2},
	}
}

// uturnArrowPoints is the arrow that doubles back on itself.
func uturnArrowPoints(x, y, w, h float64) []gopdf.Point {
	shaft := w * uturnThicknessRatio
	head := shaft * 1.6
	right := x + w
	return []gopdf.Point{
		{X: x, Y: y + h}, {X: x, Y: y + shaft},
		{X: x + shaft*1.5, Y: y}, {X: right - shaft*1.5, Y: y},
		{X: right, Y: y + shaft}, {X: right, Y: y + h - head},
		{X: right + shaft/2 - shaft, Y: y + h - head}, {X: right - shaft/2, Y: y + h},
		{X: right - shaft*1.5 - shaft/2, Y: y + h - head},
		{X: right - shaft, Y: y + h - head}, {X: right - shaft, Y: y + shaft*1.2},
		{X: x + shaft, Y: y + shaft*1.2}, {X: x + shaft, Y: y + h},
	}
}

// swooshArrowPoints is the curved sweep with a head on its end.
func swooshArrowPoints(x, y, w, h float64) []gopdf.Point {
	points := make([]gopdf.Point, 0, 2*arcSegments+4)
	// The outer sweep, then the inner one back, gives the tapering band.
	for i := range arcSegments + 1 {
		t := float64(i) / float64(arcSegments)
		points = append(points, gopdf.Point{
			X: x + w*t,
			Y: y + h - h*math.Sin(t*math.Pi/2),
		})
	}
	points = append(points,
		gopdf.Point{X: x + w, Y: y + h*0.35},
		gopdf.Point{X: x + w*0.72, Y: y + h*0.16},
	)
	for i := range arcSegments + 1 {
		t := 1 - float64(i)/float64(arcSegments)
		points = append(points, gopdf.Point{
			X: x + w*t*0.92,
			Y: y + h - h*math.Sin(t*math.Pi/2)*0.78,
		})
	}
	return points
}

// drawStripedRightArrow draws a right arrow whose tail is broken into stripes.
func drawStripedRightArrow(pdf *gopdf.GoPdf, fl flipState, x, y, w, h float64, style string) {
	geom := arrowGeometryFor(nil, w, h)
	top := y + (h-geom.shaft)/2
	tail := w * stripedArrowTailRatio
	stripe := tail / (stripedArrowStripes*2 - 1)
	for i := range stripedArrowStripes {
		left := x + float64(i)*stripe*2
		fl.polygon(pdf, []gopdf.Point{
			{X: left, Y: top}, {X: left + stripe, Y: top},
			{X: left + stripe, Y: top + geom.shaft}, {X: left, Y: top + geom.shaft},
		}, style)
	}
	fl.polygon(pdf, rightArrowPoints(x+tail, y, w-tail, h, arrowGeometryFor(nil, w-tail, h)), style)
}

// arrowCalloutPoints is a text box with an arrow growing out of one side.
func arrowCalloutPoints(x, y, w, h float64, dir arrowDirection) []gopdf.Point {
	switch dir {
	case arrowRight:
		return rightArrowCalloutOutline(x, y, w, h)
	case arrowLeft:
		return mirrorHorizontally(rightArrowCalloutOutline(x, y, w, h), x, w)
	case arrowUp:
		return transposeAbout(rightArrowCalloutOutline(x, y, h, w), x, y, w, h)
	default:
		return mirrorVertically(transposeAbout(rightArrowCalloutOutline(x, y, h, w), x, y, w, h), y, h)
	}
}

// rightArrowCalloutOutline is the box-plus-arrow with the arrow on the right.
func rightArrowCalloutOutline(x, y, w, h float64) []gopdf.Point {
	bodyW := w * arrowCalloutBodyRatio
	shaft := h * arrowCalloutShaftRatio
	head := h * arrowCalloutHeadRatio
	cy := y + h/2
	return []gopdf.Point{
		{X: x, Y: y}, {X: x + bodyW, Y: y},
		{X: x + bodyW, Y: cy - shaft}, {X: x + w - (w - bodyW), Y: cy - shaft},
		{X: x + w - (h * arrowCalloutHeadRatio), Y: cy - shaft},
		{X: x + w - (h * arrowCalloutHeadRatio), Y: cy - head},
		{X: x + w, Y: cy},
		{X: x + w - (h * arrowCalloutHeadRatio), Y: cy + head},
		{X: x + w - (h * arrowCalloutHeadRatio), Y: cy + shaft},
		{X: x + bodyW, Y: cy + shaft}, {X: x + bodyW, Y: y + h}, {X: x, Y: y + h},
	}
}

// leftRightArrowCalloutPoints is a box with an arrow on each side.
func leftRightArrowCalloutPoints(x, y, w, h float64) []gopdf.Point {
	shaft := h * arrowCalloutShaftRatio
	head := h * arrowCalloutHeadRatio
	inset := w * 0.2
	cy := y + h/2
	return []gopdf.Point{
		{X: x + inset, Y: y}, {X: x + w - inset, Y: y},
		{X: x + w - inset, Y: cy - shaft}, {X: x + w - inset, Y: cy - head},
		{X: x + w, Y: cy}, {X: x + w - inset, Y: cy + head},
		{X: x + w - inset, Y: cy + shaft}, {X: x + w - inset, Y: y + h},
		{X: x + inset, Y: y + h}, {X: x + inset, Y: cy + shaft},
		{X: x + inset, Y: cy + head}, {X: x, Y: cy},
		{X: x + inset, Y: cy - head}, {X: x + inset, Y: cy - shaft},
	}
}

// upDownArrowCalloutPoints is a box with an arrow above and below.
func upDownArrowCalloutPoints(x, y, w, h float64) []gopdf.Point {
	return transposeAbout(leftRightArrowCalloutPoints(x, y, h, w), x, y, w, h)
}

// quadArrowCalloutPoints is a box with an arrow on all four sides.
func quadArrowCalloutPoints(x, y, w, h float64) []gopdf.Point {
	shaftX := w * arrowCalloutShaftRatio / 2
	shaftY := h * arrowCalloutShaftRatio / 2
	headX := w * 0.16
	headY := h * 0.16
	insetX := w * 0.22
	insetY := h * 0.22
	cx, cy := x+w/2, y+h/2
	return []gopdf.Point{
		{X: cx - shaftX, Y: y + insetY}, {X: cx - headX, Y: y + insetY},
		{X: cx, Y: y}, {X: cx + headX, Y: y + insetY}, {X: cx + shaftX, Y: y + insetY},
		{X: x + w - insetX, Y: y + insetY}, {X: x + w - insetX, Y: cy - shaftY},
		{X: x + w - insetX, Y: cy - headY}, {X: x + w, Y: cy},
		{X: x + w - insetX, Y: cy + headY}, {X: x + w - insetX, Y: cy + shaftY},
		{X: x + w - insetX, Y: y + h - insetY}, {X: cx + shaftX, Y: y + h - insetY},
		{X: cx + headX, Y: y + h - insetY}, {X: cx, Y: y + h},
		{X: cx - headX, Y: y + h - insetY}, {X: cx - shaftX, Y: y + h - insetY},
		{X: x + insetX, Y: y + h - insetY}, {X: x + insetX, Y: cy + shaftY},
		{X: x + insetX, Y: cy + headY}, {X: x, Y: cy},
		{X: x + insetX, Y: cy - headY}, {X: x + insetX, Y: cy - shaftY},
		{X: x + insetX, Y: y + insetY},
	}
}

// mirrorHorizontally reflects an outline about the vertical centre of its box,
// which turns a right-pointing preset into its left-pointing twin.
func mirrorHorizontally(points []gopdf.Point, x, w float64) []gopdf.Point {
	out := make([]gopdf.Point, len(points))
	for i, p := range points {
		out[i] = gopdf.Point{X: 2*x + w - p.X, Y: p.Y}
	}
	return out
}

func mirrorVertically(points []gopdf.Point, y, h float64) []gopdf.Point {
	out := make([]gopdf.Point, len(points))
	for i, p := range points {
		out[i] = gopdf.Point{X: p.X, Y: 2*y + h - p.Y}
	}
	return out
}

// transposeAbout turns an outline built for a landscape box a quarter turn, so
// a right-pointing preset becomes the up-pointing one.
func transposeAbout(points []gopdf.Point, x, y, w, h float64) []gopdf.Point {
	out := make([]gopdf.Point, len(points))
	for i, p := range points {
		out[i] = gopdf.Point{X: x + (p.Y-y)*w/h, Y: y + (p.X-x)*h/w}
	}
	return out
}
