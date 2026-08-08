//nolint:mnd // Preset silhouettes use fixed proportions from the DrawingML geometry.
package export

import (
	"math"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// The remaining presets that draw as a single closed outline: the solids, the
// banner and ribbon family, the frames and corners, and the rings. Each used to
// fall through to a plain rectangle.

// drawPDFSolidGeometry draws one of the presets below and reports whether it
// recognised the type.
func drawPDFSolidGeometry(
	pdf *gopdf.GoPdf,
	fl flipState,
	shapeType string,
	x, y, w, h float64,
	style string,
) bool {
	switch shapeType {
	case shapes.ShapeTypeCube:
		fl.polygon(pdf, cubePoints(x, y, w, h), style)
	case shapes.ShapeTypeFoldedCorner:
		fl.polygon(pdf, foldedCornerPoints(x, y, w, h), style)
	case shapes.ShapeTypeCorner:
		fl.polygon(pdf, cornerPoints(x, y, w, h), style)
	case shapes.ShapeTypeDiagStripe:
		fl.polygon(pdf, diagStripePoints(x, y, w, h), style)
	case shapes.ShapeTypeHalfFrame:
		fl.polygon(pdf, halfFramePoints(x, y, w, h), style)
	case shapes.ShapeTypePlaque:
		fl.polygon(pdf, plaquePoints(x, y, w, h), style)
	case shapes.ShapeTypeChevron:
		fl.polygon(pdf, chevronPoints(x, y, w, h), style)
	case shapes.ShapeTypeNotchedRightArrow:
		fl.polygon(pdf, notchedRightArrowPoints(x, y, w, h), style)
	case shapes.ShapeTypeTeardrop:
		fl.polygon(pdf, teardropPoints(x, y, w, h), style)
	case shapes.ShapeTypeMoon:
		fl.polygon(pdf, moonPoints(x, y, w, h), style)
	case shapes.ShapeTypeRibbon:
		fl.polygon(pdf, ribbonPoints(x, y, w, h, false), style)
	case shapes.ShapeTypeRibbon2:
		fl.polygon(pdf, ribbonPoints(x, y, w, h, true), style)
	default:
		return false
	}
	return true
}

const (
	// cubeDepthRatio is how much of the shorter side the cube's top and side
	// faces take.
	cubeDepthRatio = 0.25
	// foldRatio is the size of a folded corner, as a fraction of the shorter side.
	foldRatio = 0.2
	// frameThicknessRatio is how thick a frame's border is.
	frameThicknessRatio = 0.25
	// stripeRatio is the width of a diagonal stripe across the box.
	stripeRatio = 0.5
	// plaqueCutRatio is how deep a plaque's concave corners cut in.
	plaqueCutRatio = 0.16667
	// notchRatio is how deep the tail notch of a notched arrow cuts in.
	notchRatio = 0.5
	// ribbonTailRatio is how far the ribbon's tails hang below its band.
	ribbonTailRatio = 0.25
	// arcSegments is the resolution of a curved silhouette.
	arcSegments = 24
)

// cubePoints is the silhouette of a cube seen from the front-left: the front
// face with the top and right faces behind it.
func cubePoints(x, y, w, h float64) []gopdf.Point {
	d := math.Min(w, h) * cubeDepthRatio
	return []gopdf.Point{
		{X: x, Y: y + d}, {X: x + d, Y: y}, {X: x + w, Y: y},
		{X: x + w, Y: y + h - d}, {X: x + w - d, Y: y + h}, {X: x, Y: y + h},
	}
}

// foldedCornerPoints is a rectangle with its bottom-right corner turned up.
func foldedCornerPoints(x, y, w, h float64) []gopdf.Point {
	fold := math.Min(w, h) * foldRatio
	return []gopdf.Point{
		{X: x, Y: y}, {X: x + w, Y: y},
		{X: x + w, Y: y + h - fold}, {X: x + w - fold, Y: y + h},
		{X: x, Y: y + h},
	}
}

// cornerPoints is the L-shape of the corner preset.
func cornerPoints(x, y, w, h float64) []gopdf.Point {
	tw := w * frameThicknessRatio
	th := h * frameThicknessRatio
	return []gopdf.Point{
		{X: x, Y: y}, {X: x + tw, Y: y},
		{X: x + tw, Y: y + h - th}, {X: x + w, Y: y + h - th},
		{X: x + w, Y: y + h}, {X: x, Y: y + h},
	}
}

// diagStripePoints is the band running corner to corner across the box.
func diagStripePoints(x, y, w, h float64) []gopdf.Point {
	return []gopdf.Point{
		{X: x, Y: y + h*stripeRatio}, {X: x + w*stripeRatio, Y: y},
		{X: x + w, Y: y}, {X: x, Y: y + h},
	}
}

// halfFramePoints is the two-sided frame: the top and left borders only.
func halfFramePoints(x, y, w, h float64) []gopdf.Point {
	tw := w * frameThicknessRatio
	th := h * frameThicknessRatio
	return []gopdf.Point{
		{X: x, Y: y}, {X: x + w, Y: y}, {X: x + w - tw, Y: y + th},
		{X: x + tw, Y: y + th}, {X: x + tw, Y: y + h - th}, {X: x, Y: y + h},
	}
}

// plaquePoints is a rectangle whose four corners curve inwards.
func plaquePoints(x, y, w, h float64) []gopdf.Point {
	cut := math.Min(w, h) * plaqueCutRatio
	points := make([]gopdf.Point, 0, 4*(arcSegments/4+1))
	// Each corner is a quarter arc turning the other way, so the outline bites
	// into the box rather than rounding off it.
	corners := [4]struct {
		centre   gopdf.Point
		startDeg float64
	}{
		{gopdf.Point{X: x, Y: y}, 90},
		{gopdf.Point{X: x + w, Y: y}, 180},
		{gopdf.Point{X: x + w, Y: y + h}, 270},
		{gopdf.Point{X: x, Y: y + h}, 0},
	}
	for _, corner := range corners {
		points = append(points, arcPoints(corner.centre, cut, corner.startDeg, -90, arcSegments/4)...)
	}
	return points
}

// notchedRightArrowPoints is a right arrow with a V cut into its tail.
func notchedRightArrowPoints(x, y, w, h float64) []gopdf.Point {
	geom := arrowGeometryFor(nil, w, h)
	bodyTop := y + (h-geom.shaft)/2
	bodyBottom := bodyTop + geom.shaft
	shaftEnd := x + w - geom.head
	notch := geom.head * notchRatio
	return []gopdf.Point{
		{X: x, Y: bodyTop}, {X: shaftEnd, Y: bodyTop},
		{X: shaftEnd, Y: y}, {X: x + w, Y: y + h/2},
		{X: shaftEnd, Y: y + h}, {X: shaftEnd, Y: bodyBottom},
		{X: x, Y: bodyBottom}, {X: x + notch, Y: y + h/2},
	}
}

// teardropPoints is a circle with its top-right quadrant drawn out to a point.
func teardropPoints(x, y, w, h float64) []gopdf.Point {
	cx, cy := x+w/2, y+h/2
	rx, ry := w/2, h/2
	// Three quarters of the ellipse, then out to the top-right corner.
	points := make([]gopdf.Point, 0, arcSegments+2)
	for i := range arcSegments + 1 {
		angle := -math.Pi/2 + (3*math.Pi/2)*float64(i)/float64(arcSegments)
		points = append(points, gopdf.Point{X: cx + rx*math.Cos(angle), Y: cy + ry*math.Sin(angle)})
	}
	return append(points, gopdf.Point{X: x + w, Y: y})
}

// moonPoints is a crescent: the outer half-circle closed by a shallower inner curve.
func moonPoints(x, y, w, h float64) []gopdf.Point {
	cx, cy := x+w, y+h/2
	points := make([]gopdf.Point, 0, 2*arcSegments+2)
	for i := range arcSegments + 1 {
		angle := math.Pi/2 + math.Pi*float64(i)/float64(arcSegments)
		points = append(points, gopdf.Point{X: cx + w*math.Cos(angle), Y: cy + (h/2)*math.Sin(angle)})
	}
	for i := range arcSegments + 1 {
		angle := -math.Pi/2 - math.Pi*float64(i)/float64(arcSegments)
		points = append(points, gopdf.Point{
			X: cx + (w * moonInnerRatio * math.Cos(angle)),
			Y: cy + (h / 2 * math.Sin(angle)),
		})
	}
	return points
}

// moonInnerRatio is how far the crescent's inner curve reaches back, as a
// fraction of the width.
const moonInnerRatio = 0.5

// ribbonPoints is the banner family: a band across the box with two tails,
// pointing down for ribbon and up for ribbon2.
func ribbonPoints(x, y, w, h float64, up bool) []gopdf.Point {
	tail := h * ribbonTailRatio
	inset := w * ribbonInsetRatio
	// ribbon hangs its tails below the band; ribbon2 folds them above it, which
	// is the whole difference between the two presets.
	bandTop, bandBottom, tailEdge := y+tail, y+h-tail, y+h
	if up {
		bandTop, bandBottom, tailEdge = y+tail, y+h-tail, y
	}
	return []gopdf.Point{
		{X: x, Y: bandTop}, {X: x + inset, Y: y + h/2}, {X: x, Y: bandBottom},
		{X: x + inset, Y: bandBottom}, {X: x + inset, Y: tailEdge},
		{X: x + w - inset, Y: tailEdge}, {X: x + w - inset, Y: bandBottom},
		{X: x + w, Y: bandBottom}, {X: x + w - inset, Y: y + h/2}, {X: x + w, Y: bandTop},
		{X: x + w - inset, Y: bandTop},
	}
}

// ribbonInsetRatio is how far in from each side the ribbon's tails fold.
const ribbonInsetRatio = 0.15
