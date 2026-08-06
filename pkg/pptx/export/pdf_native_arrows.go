//nolint:mnd // Arrow geometry uses fixed visual proportion constants for correct arrow shape rendering.
package export

import (
	"math"

	"github.com/signintech/gopdf"
)

func rightArrowPoints(x, y, w, h float64, geom arrowGeometry) []gopdf.Point {
	bw := w - geom.head // where the head starts
	shaftTop := y + (h-geom.shaft)/2
	shaftBottom := shaftTop + geom.shaft
	return []gopdf.Point{
		{X: x, Y: shaftTop},
		{X: x + bw, Y: shaftTop},
		{X: x + bw, Y: y},
		{X: x + w, Y: y + h/2},
		{X: x + bw, Y: y + h},
		{X: x + bw, Y: shaftBottom},
		{X: x, Y: shaftBottom},
	}
}

func leftArrowPoints(x, y, w, h float64, geom arrowGeometry) []gopdf.Point {
	aw := geom.head // where the head ends
	shaftTop := y + (h-geom.shaft)/2
	shaftBottom := shaftTop + geom.shaft
	return []gopdf.Point{
		{X: x + aw, Y: shaftTop},
		{X: x + w, Y: shaftTop},
		{X: x + w, Y: shaftBottom},
		{X: x + aw, Y: shaftBottom},
		{X: x + aw, Y: y + h},
		{X: x, Y: y + h/2},
		{X: x + aw, Y: y},
	}
}

func upArrowPoints(x, y, w, h float64, geom arrowGeometry) []gopdf.Point {
	sh := geom.head // head height, measured from the top
	lx := x + (w-geom.shaft)/2
	rx := lx + geom.shaft
	return []gopdf.Point{
		{X: lx, Y: y + h},
		{X: rx, Y: y + h},
		{X: rx, Y: y + sh},
		{X: x + w, Y: y + sh},
		{X: x + w/2, Y: y},
		{X: x, Y: y + sh},
		{X: lx, Y: y + sh},
	}
}

func downArrowPoints(x, y, w, h float64, geom arrowGeometry) []gopdf.Point {
	sh := h - geom.head // where the head starts
	lx := x + (w-geom.shaft)/2
	rx := lx + geom.shaft
	return []gopdf.Point{
		{X: lx, Y: y},
		{X: rx, Y: y},
		{X: rx, Y: sh + y},
		{X: x + w, Y: sh + y},
		{X: x + w/2, Y: y + h},
		{X: x, Y: sh + y},
		{X: lx, Y: sh + y},
	}
}

func leftRightArrowPoints(x, y, w, h float64, geom arrowGeometry) []gopdf.Point {
	// Both heads have to fit inside the width, so each takes at most half of it.
	hw := math.Min(geom.head, w/2)
	shy := geom.shaft / 2 // shaft half-height from centre
	cy := y + h/2
	return []gopdf.Point{
		{X: x, Y: cy},
		{X: x + hw, Y: y},
		{X: x + hw, Y: cy - shy},
		{X: x + w - hw, Y: cy - shy},
		{X: x + w - hw, Y: y},
		{X: x + w, Y: cy},
		{X: x + w - hw, Y: y + h},
		{X: x + w - hw, Y: cy + shy},
		{X: x + hw, Y: cy + shy},
		{X: x + hw, Y: y + h},
	}
}

func upDownArrowPoints(x, y, w, h float64, geom arrowGeometry) []gopdf.Point {
	hh := math.Min(geom.head, h/2) // head height, one at each end
	shx := geom.shaft / 2          // shaft half-width from centre
	cx := x + w/2
	return []gopdf.Point{
		{X: cx, Y: y},
		{X: x + w, Y: y + hh},
		{X: cx + shx, Y: y + hh},
		{X: cx + shx, Y: y + h - hh},
		{X: x + w, Y: y + h - hh},
		{X: cx, Y: y + h},
		{X: x, Y: y + h - hh},
		{X: cx - shx, Y: y + h - hh},
		{X: cx - shx, Y: y + hh},
		{X: x, Y: y + hh},
	}
}

func chevronPoints(x, y, w, h float64) []gopdf.Point {
	notch := w * 0.25
	tip := w * 0.75
	return []gopdf.Point{
		{X: x, Y: y},
		{X: x + tip, Y: y},
		{X: x + w, Y: y + h/2},
		{X: x + tip, Y: y + h},
		{X: x, Y: y + h},
		{X: x + notch, Y: y + h/2},
	}
}
