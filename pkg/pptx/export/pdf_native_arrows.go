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
