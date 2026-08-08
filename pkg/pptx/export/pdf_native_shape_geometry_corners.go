//nolint:mnd // Corner treatments are fixed proportions from the DrawingML preset geometry.
package export

import (
	"math"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// The rounded- and snipped-corner rectangles are one family: each names which
// corners are treated and whether the treatment is an arc or a cut. They all
// fell through to a plain rectangle, so a snipped tab and a rounded card came
// out as the same box.

// cornerTreatment is what happens at one corner of a rectangle.
type cornerTreatment uint8

const (
	cornerSquare cornerTreatment = iota
	cornerRounded
	cornerSnipped
)

// cornerRatio is how much of the shorter side a treated corner takes, matching
// the presets' default 16667 adjustment.
const cornerRatio = 0.16667

// cornerResolution is how many segments a rounded corner is drawn with.
const cornerResolution = 8

// drawPDFCornerGeometry draws the rounded/snipped rectangle family and reports
// whether it recognised the type. Corners are listed clockwise from top-left.
func drawPDFCornerGeometry(
	pdf *gopdf.GoPdf,
	fl flipState,
	shapeType string,
	x, y, w, h float64,
	style string,
) bool {
	corners, ok := cornersForPreset(shapeType)
	if !ok {
		return false
	}
	fl.polygon(pdf, roundedRectPoints(x, y, w, h, corners), style)
	return true
}

// cornersForPreset maps a preset to its four corner treatments, clockwise from
// the top-left.
func cornersForPreset(shapeType string) ([4]cornerTreatment, bool) {
	switch shapeType {
	case shapes.ShapeTypeRound1Rect:
		return [4]cornerTreatment{cornerSquare, cornerRounded, cornerSquare, cornerSquare}, true
	case shapes.ShapeTypeRound2SameRect:
		return [4]cornerTreatment{cornerRounded, cornerRounded, cornerSquare, cornerSquare}, true
	case shapes.ShapeTypeRound2DiagRect:
		return [4]cornerTreatment{cornerRounded, cornerSquare, cornerRounded, cornerSquare}, true
	case shapes.ShapeTypeSnip1Rect:
		return [4]cornerTreatment{cornerSquare, cornerSnipped, cornerSquare, cornerSquare}, true
	case shapes.ShapeTypeSnip2SameRect:
		return [4]cornerTreatment{cornerSnipped, cornerSnipped, cornerSquare, cornerSquare}, true
	case shapes.ShapeTypeSnip2DiagRect:
		return [4]cornerTreatment{cornerSnipped, cornerSquare, cornerSnipped, cornerSquare}, true
	case shapes.ShapeTypeSnipRoundRect:
		return [4]cornerTreatment{cornerRounded, cornerSnipped, cornerSquare, cornerSquare}, true
	default:
		return [4]cornerTreatment{}, false
	}
}

// roundedRectPoints traces a rectangle whose corners are squared, rounded or
// cut, clockwise from the top-left.
func roundedRectPoints(x, y, w, h float64, corners [4]cornerTreatment) []gopdf.Point {
	cut := math.Min(w, h) * cornerRatio
	points := make([]gopdf.Point, 0, 4*(cornerResolution+2))

	// Each corner is described by its own vertex and the centre an arc would
	// turn about, so one loop covers all four.
	type cornerSpec struct {
		vertex   gopdf.Point
		centre   gopdf.Point
		startDeg float64
	}
	specs := [4]cornerSpec{
		{gopdf.Point{X: x, Y: y}, gopdf.Point{X: x + cut, Y: y + cut}, 180},
		{gopdf.Point{X: x + w, Y: y}, gopdf.Point{X: x + w - cut, Y: y + cut}, 270},
		{gopdf.Point{X: x + w, Y: y + h}, gopdf.Point{X: x + w - cut, Y: y + h - cut}, 0},
		{gopdf.Point{X: x, Y: y + h}, gopdf.Point{X: x + cut, Y: y + h - cut}, 90},
	}

	for i, spec := range specs {
		switch corners[i] {
		case cornerSquare:
			points = append(points, spec.vertex)
		case cornerSnipped:
			points = append(points, snipCornerPoints(spec.vertex, x, y, w, h, cut)...)
		case cornerRounded:
			points = append(points, arcPoints(spec.centre, cut, spec.startDeg, 90, cornerResolution)...)
		}
	}
	return points
}

// snipCornerPoints replaces a corner with the straight cut across it, in
// clockwise order.
func snipCornerPoints(vertex gopdf.Point, x, y, w, h, cut float64) []gopdf.Point {
	left := math.Abs(vertex.X-x) < nearZeroEpsilon
	top := math.Abs(vertex.Y-y) < nearZeroEpsilon
	switch {
	case left && top:
		return []gopdf.Point{{X: x, Y: y + cut}, {X: x + cut, Y: y}}
	case !left && top:
		return []gopdf.Point{{X: x + w - cut, Y: y}, {X: x + w, Y: y + cut}}
	case !left && !top:
		return []gopdf.Point{{X: x + w, Y: y + h - cut}, {X: x + w - cut, Y: y + h}}
	default:
		return []gopdf.Point{{X: x + cut, Y: y + h}, {X: x, Y: y + h - cut}}
	}
}

// arcPoints walks an arc of the given radius about centre, clockwise in the
// renderer's y-down space, starting at startDeg and turning sweepDeg.
func arcPoints(centre gopdf.Point, radius, startDeg, sweepDeg float64, segments int) []gopdf.Point {
	points := make([]gopdf.Point, 0, segments+1)
	for i := range segments + 1 {
		angle := (startDeg + sweepDeg*float64(i)/float64(segments)) * math.Pi / 180
		points = append(points, gopdf.Point{
			X: centre.X + radius*math.Cos(angle),
			Y: centre.Y + radius*math.Sin(angle),
		})
	}
	return points
}
