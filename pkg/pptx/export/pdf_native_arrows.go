//nolint:mnd // Arrow geometry uses fixed visual proportion constants for correct arrow shape rendering.
package export

import "github.com/signintech/gopdf"

func rightArrowPoints(x, y, w, h float64) []gopdf.Point {
	aw := w * 0.5
	bw := w - aw
	hh := h * 0.5
	bh := h * 0.5
	return []gopdf.Point{
		{X: x, Y: y + (h-bh)/2},
		{X: x + bw, Y: y + (h-bh)/2},
		{X: x + bw, Y: y},
		{X: x + w, Y: y + hh},
		{X: x + bw, Y: y + h},
		{X: x + bw, Y: y + h - (h-bh)/2},
		{X: x, Y: y + h - (h-bh)/2},
	}
}

func leftArrowPoints(x, y, w, h float64) []gopdf.Point {
	aw := w * 0.5
	hh := h * 0.5
	bh := h * 0.5
	return []gopdf.Point{
		{X: x + aw, Y: y + (h-bh)/2},
		{X: x + w, Y: y + (h-bh)/2},
		{X: x + w, Y: y + h - (h-bh)/2},
		{X: x + aw, Y: y + h - (h-bh)/2},
		{X: x + aw, Y: y + h},
		{X: x, Y: y + hh},
		{X: x + aw, Y: y},
	}
}

func upArrowPoints(x, y, w, h float64) []gopdf.Point {
	sh := h * 0.5 // shaft height
	lx := x + w*0.3
	rx := x + w*0.7
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

func downArrowPoints(x, y, w, h float64) []gopdf.Point {
	sh := h * 0.5 // shaft height
	lx := x + w*0.3
	rx := x + w*0.7
	return []gopdf.Point{
		{X: lx, Y: y},
		{X: rx, Y: y},
		{X: rx, Y: y + sh},
		{X: x + w, Y: y + sh},
		{X: x + w/2, Y: y + h},
		{X: x, Y: y + sh},
		{X: lx, Y: y + sh},
	}
}

func leftRightArrowPoints(x, y, w, h float64) []gopdf.Point {
	hw := w * 0.3   // head width
	shy := h * 0.25 // shaft half-height from centre
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

func upDownArrowPoints(x, y, w, h float64) []gopdf.Point {
	hh := h * 0.3   // head height
	shx := w * 0.25 // shaft half-width from centre
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
