//nolint:mnd // Curved and stroked presets use fixed proportions from the DrawingML geometry.
package export

import (
	"math"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// The presets that are stroked rather than filled — brackets, braces, arcs and
// the connector lines — and the ones built from an arc band: the curved arrows,
// the ellipse ribbons, the waves, the scrolls and the gears.

// drawPDFCurveGeometry draws one of the presets below and reports whether it
// recognised the type.
func drawPDFCurveGeometry(
	pdf *gopdf.GoPdf,
	fl flipState,
	shapeType string,
	x, y, w, h float64,
	style string,
) bool {
	switch shapeType {
	case shapes.ShapeTypeLeftBracket:
		fl.polyline(pdf, bracketPoints(x, y, w, h, true))
	case shapes.ShapeTypeRightBracket:
		fl.polyline(pdf, bracketPoints(x, y, w, h, false))
	case shapes.ShapeTypeBracketPair, shapes.ShapeTypeDoubleBracket:
		fl.polyline(pdf, bracketPoints(x, y, w, h, true))
		fl.polyline(pdf, bracketPoints(x, y, w, h, false))
	case shapes.ShapeTypeLeftBrace:
		fl.polyline(pdf, bracePoints(x, y, w, h, true))
	case shapes.ShapeTypeRightBrace:
		fl.polyline(pdf, bracePoints(x, y, w, h, false))
	case shapes.ShapeTypeBracePair, shapes.ShapeTypeDoubleBrace:
		fl.polyline(pdf, bracePoints(x, y, w, h, true))
		fl.polyline(pdf, bracePoints(x, y, w, h, false))
	case shapes.ShapeTypeArc:
		fl.polyline(pdf, ellipseArcPoints(x+w/2, y+h/2, w/2, h/2, 180, 180))
	case shapes.ShapeTypeCurvedRightArrow:
		fl.polygon(pdf, curvedArrowPoints(x, y, w, h, curveRight), style)
	case shapes.ShapeTypeCurvedLeftArrow:
		fl.polygon(pdf, curvedArrowPoints(x, y, w, h, curveLeft), style)
	case shapes.ShapeTypeCurvedUpArrow:
		fl.polygon(pdf, curvedArrowPoints(x, y, w, h, curveUp), style)
	case shapes.ShapeTypeCurvedDownArrow:
		fl.polygon(pdf, curvedArrowPoints(x, y, w, h, curveDown), style)
	case shapes.ShapeTypeCurvedLeftRightArrow, shapes.ShapeTypeCurvedUpDownArrow,
		shapes.ShapeTypeCircularArrow:
		fl.polygon(pdf, circularArrowPoints(x, y, w, h), style)
	case shapes.ShapeTypeEllipseRibbon:
		fl.polygon(pdf, ellipseRibbonPoints(x, y, w, h, false), style)
	case shapes.ShapeTypeEllipseRibbon2:
		fl.polygon(pdf, ellipseRibbonPoints(x, y, w, h, true), style)
	case shapes.ShapeTypeWave:
		fl.polygon(pdf, wavePoints(x, y, w, h, 1), style)
	case shapes.ShapeTypeDoubleWave:
		fl.polygon(pdf, wavePoints(x, y, w, h, 2), style)
	case shapes.ShapeTypeHorizontalScroll:
		fl.polygon(pdf, scrollPoints(x, y, w, h, true), style)
	case shapes.ShapeTypeVerticalScroll:
		fl.polygon(pdf, scrollPoints(x, y, w, h, false), style)
	case shapes.ShapeTypeGear6:
		fl.polygon(pdf, gearPoints(x, y, w, h, 6), style)
	case shapes.ShapeTypeGear9:
		fl.polygon(pdf, gearPoints(x, y, w, h, 9), style)
	case shapes.ShapeTypeSmileyFace:
		drawSmileyFace(pdf, fl, x, y, w, h, style)
	default:
		return false
	}
	return true
}

const (
	// bracketRatio is how far a bracket's arms reach in, as a fraction of width.
	bracketRatio = 0.5
	// curvedArrowBandRatio is the thickness of a curved arrow's band, as a
	// fraction of the radius.
	curvedArrowBandRatio = 0.45
	// curvedArrowHeadRatio is how much wider the head is than the band.
	curvedArrowHeadRatio = 1.9
	// curvedArrowHeadDeg is how much of the sweep the head takes.
	curvedArrowHeadDeg = 32.0
	// waveDepthRatio is the height of a wave's crest.
	waveDepthRatio = 0.12
	// scrollCurlRatio is the size of a scroll's rolled edge.
	scrollCurlRatio = 0.22
	// gearToothRatio is how far a gear's teeth stand out from its body.
	gearToothRatio = 0.24
	// smileyFeatureRatio sizes a smiley's eyes against the face.
	smileyFeatureRatio = 0.1
)

// curveDirection is which way a curved arrow sweeps.
type curveDirection uint8

const (
	curveRight curveDirection = iota
	curveLeft
	curveUp
	curveDown
)

// ellipseArcPoints walks an arc of an ellipse, which the circular presets need
// and the circle-only arcPoints cannot express.
func ellipseArcPoints(cx, cy, rx, ry, startDeg, sweepDeg float64) []gopdf.Point {
	const segments = arcSegments
	points := make([]gopdf.Point, 0, segments+1)
	for i := range segments + 1 {
		angle := (startDeg + sweepDeg*float64(i)/float64(segments)) * math.Pi / 180
		points = append(points, gopdf.Point{
			X: cx + rx*math.Cos(angle),
			Y: cy + ry*math.Sin(angle),
		})
	}
	return points
}

// bracketPoints is one square bracket, opening right when left is true.
func bracketPoints(x, y, w, h float64, left bool) []gopdf.Point {
	arm := w * bracketRatio / 2
	spine := x
	if !left {
		spine = x + w
		arm = -arm
	}
	return []gopdf.Point{
		{X: spine + arm, Y: y}, {X: spine, Y: y},
		{X: spine, Y: y + h}, {X: spine + arm, Y: y + h},
	}
}

// bracePoints is one curly brace, with the pinch at its middle.
func bracePoints(x, y, w, h float64, left bool) []gopdf.Point {
	arm := w * bracketRatio / 2
	spine := x + arm
	tip := x
	if !left {
		spine = x + w - arm
		tip = x + w
	}
	mid := y + h/2
	return []gopdf.Point{
		{X: spine + (spine - tip), Y: y},
		{X: spine, Y: y + h*0.08},
		{X: spine, Y: mid - h*0.08},
		{X: tip, Y: mid},
		{X: spine, Y: mid + h*0.08},
		{X: spine, Y: y + h*0.92},
		{X: spine + (spine - tip), Y: y + h},
	}
}

// curvedArrowPoints is a quarter-turn band with a head on its end. The band is
// traced along its outer edge and back along its inner one, so it closes into a
// single outline.
func curvedArrowPoints(x, y, w, h float64, dir curveDirection) []gopdf.Point {
	cx, cy, startDeg := x, y+h, 270.0
	switch dir {
	case curveRight:
		cx, cy, startDeg = x, y+h, 270
	case curveLeft:
		cx, cy, startDeg = x+w, y+h, 180
	case curveUp:
		cx, cy, startDeg = x, y, 0
	case curveDown:
		cx, cy, startDeg = x+w, y, 90
	}
	outerRX, outerRY := w, h
	band := math.Min(w, h) * curvedArrowBandRatio
	innerRX, innerRY := math.Max(1, w-band), math.Max(1, h-band)

	sweep := 90.0 - curvedArrowHeadDeg
	points := ellipseArcPoints(cx, cy, outerRX, outerRY, startDeg, sweep)

	// The head: a triangle spanning the band, wider than it, at the sweep's end.
	headAngle := (startDeg + sweep) * math.Pi / 180
	tipAngle := (startDeg + 90) * math.Pi / 180
	midRX, midRY := (outerRX+innerRX)/2, (outerRY+innerRY)/2
	spread := band * (curvedArrowHeadRatio - 1) / 2
	points = append(points,
		gopdf.Point{
			X: cx + (outerRX+spread)*math.Cos(headAngle),
			Y: cy + (outerRY+spread)*math.Sin(headAngle),
		},
		gopdf.Point{X: cx + midRX*math.Cos(tipAngle), Y: cy + midRY*math.Sin(tipAngle)},
		gopdf.Point{
			X: cx + (innerRX-spread)*math.Cos(headAngle),
			Y: cy + (innerRY-spread)*math.Sin(headAngle),
		},
	)
	return append(points, ellipseArcPoints(cx, cy, innerRX, innerRY, startDeg+sweep, -sweep)...)
}

// circularArrowPoints is a nearly closed ring with a head, used for the
// circular and double curved arrows.
func circularArrowPoints(x, y, w, h float64) []gopdf.Point {
	cx, cy := x+w/2, y+h/2
	outerRX, outerRY := w/2, h/2
	band := math.Min(w, h) * curvedArrowBandRatio / 2
	innerRX, innerRY := math.Max(1, outerRX-band), math.Max(1, outerRY-band)

	const start, sweep = -60.0, 280.0
	points := ellipseArcPoints(cx, cy, outerRX, outerRY, start, sweep)

	headAngle := (start + sweep) * math.Pi / 180
	tipAngle := (start + sweep + curvedArrowHeadDeg) * math.Pi / 180
	spread := band * (curvedArrowHeadRatio - 1) / 2
	points = append(points,
		gopdf.Point{
			X: cx + (outerRX+spread)*math.Cos(headAngle),
			Y: cy + (outerRY+spread)*math.Sin(headAngle),
		},
		gopdf.Point{
			X: cx + (outerRX+innerRX)/2*math.Cos(tipAngle),
			Y: cy + (outerRY+innerRY)/2*math.Sin(tipAngle),
		},
		gopdf.Point{
			X: cx + (innerRX-spread)*math.Cos(headAngle),
			Y: cy + (innerRY-spread)*math.Sin(headAngle),
		},
	)
	return append(points, ellipseArcPoints(cx, cy, innerRX, innerRY, start+sweep, -sweep)...)
}

// ellipseRibbonPoints is the curved banner: a band that bows down for the first
// preset and up for the second.
func ellipseRibbonPoints(x, y, w, h float64, up bool) []gopdf.Point {
	bow := h * 0.3
	band := h * 0.45
	inset := w * ribbonInsetRatio
	direction := 1.0
	if up {
		direction = -1
	}
	top := make([]gopdf.Point, 0, arcSegments+1)
	bottom := make([]gopdf.Point, 0, arcSegments+1)
	for i := range arcSegments + 1 {
		t := float64(i) / float64(arcSegments)
		dip := direction * bow * math.Sin(t*math.Pi)
		top = append(top, gopdf.Point{X: x + w*t, Y: y + bow + dip})
		bottom = append(bottom, gopdf.Point{X: x + w*(1-t), Y: y + bow + band + direction*bow*math.Sin((1-t)*math.Pi)})
	}
	points := make([]gopdf.Point, 0, len(top)+len(bottom)+4)
	points = append(points, top...)
	points = append(points,
		gopdf.Point{X: x + w, Y: y + h},
		gopdf.Point{X: x + w - inset, Y: y + h - bow},
	)
	points = append(points, bottom...)
	return append(points,
		gopdf.Point{X: x + inset, Y: y + h - bow},
		gopdf.Point{X: x, Y: y + h},
	)
}

// wavePoints is a rectangle whose top and bottom edges ripple.
func wavePoints(x, y, w, h float64, periods float64) []gopdf.Point {
	depth := h * waveDepthRatio
	points := make([]gopdf.Point, 0, 2*arcSegments+2)
	for i := range arcSegments + 1 {
		t := float64(i) / float64(arcSegments)
		points = append(points, gopdf.Point{
			X: x + w*t,
			Y: y + depth + depth*math.Sin(2*math.Pi*periods*t),
		})
	}
	for i := range arcSegments + 1 {
		t := 1 - float64(i)/float64(arcSegments)
		points = append(points, gopdf.Point{
			X: x + w*t,
			Y: y + h - depth + depth*math.Sin(2*math.Pi*periods*t),
		})
	}
	return points
}

// scrollPoints is a sheet with one edge rolled up.
func scrollPoints(x, y, w, h float64, horizontal bool) []gopdf.Point {
	curl := math.Min(w, h) * scrollCurlRatio
	if horizontal {
		return []gopdf.Point{
			{X: x + curl, Y: y}, {X: x + w, Y: y}, {X: x + w, Y: y + h - curl},
			{X: x + w - curl, Y: y + h}, {X: x, Y: y + h}, {X: x, Y: y + curl},
			{X: x + curl, Y: y + curl}, {X: x + curl, Y: y},
		}
	}
	return []gopdf.Point{
		{X: x, Y: y + curl}, {X: x + curl, Y: y}, {X: x + w, Y: y},
		{X: x + w, Y: y + h - curl}, {X: x + w - curl, Y: y + h}, {X: x, Y: y + h},
		{X: x, Y: y + h - curl}, {X: x + curl, Y: y + h - curl},
	}
}

// gearPoints is a toothed wheel with the requested number of teeth.
func gearPoints(x, y, w, h float64, teeth int) []gopdf.Point {
	cx, cy := x+w/2, y+h/2
	outerRX, outerRY := w/2, h/2
	bodyRX, bodyRY := outerRX*(1-gearToothRatio), outerRY*(1-gearToothRatio)
	points := make([]gopdf.Point, 0, teeth*4)
	step := 2 * math.Pi / float64(teeth)
	for i := range teeth {
		base := float64(i) * step
		// Each tooth is a flat-topped step: up, across, down, then the gap.
		for _, spec := range [4]struct {
			offset float64
			rx, ry float64
		}{
			{0.10 * step, bodyRX, bodyRY},
			{0.18 * step, outerRX, outerRY},
			{0.42 * step, outerRX, outerRY},
			{0.50 * step, bodyRX, bodyRY},
		} {
			angle := base + spec.offset
			points = append(points, gopdf.Point{
				X: cx + spec.rx*math.Cos(angle),
				Y: cy + spec.ry*math.Sin(angle),
			})
		}
	}
	return points
}

// drawSmileyFace draws the face, then its eyes and mouth over it.
func drawSmileyFace(pdf *gopdf.GoPdf, fl flipState, x, y, w, h float64, style string) {
	cx, cy := x+w/2, y+h/2
	fl.polygon(pdf, ellipsePoints(cx, cy, w/2, h/2), style)

	eyeRX, eyeRY := w*smileyFeatureRatio/2, h*smileyFeatureRatio
	for _, dx := range []float64{-w * 0.18, w * 0.18} {
		fl.polygon(pdf, ellipsePoints(cx+dx, cy-h*0.15, eyeRX, eyeRY), "D")
	}
	// The mouth is an arc across the lower half, stroked rather than filled.
	fl.polyline(pdf, ellipseArcPoints(cx, cy+h*0.05, w*0.28, h*0.22, 20, 140))
}
