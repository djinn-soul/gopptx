//nolint:mnd // Gradient renderer uses tuned numeric banding/orientation constants for visual parity.
package export

import (
	"math"
	"sort"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

const gradientBands = 96

type gradientStop struct {
	pos   float64
	color rgbColor
	alpha float64
}

func renderPDFGradientBackground(pdf *gopdf.GoPdf, grad *shapes.ShapeGradientFill) bool {
	return renderPDFLinearGradientRect(pdf, grad, 0, 0, slideWidthPt, slideHeightPt)
}

func renderPDFShapeGradient(
	pdf *gopdf.GoPdf,
	s shapes.Shape,
	x, y, w, h float64,
) bool {
	if s.GradientFill == nil || len(s.GradientFill.Stops) < 2 {
		return false
	}
	switch s.Type {
	case shapes.ShapeTypeRectangle, shapes.ShapeTypeRoundedRectangle:
		return renderPDFLinearGradientRect(pdf, s.GradientFill, x, y, w, h)
	case shapes.ShapeTypeEllipse:
		return renderPDFLinearGradientEllipse(pdf, s.GradientFill, x, y, w, h)
	default:
		return false
	}
}

// renderPDFLinearGradientEllipse renders a linear gradient clipped to an ellipse
// by drawing banded rectangles trimmed to the ellipse silhouette at each position.
func renderPDFLinearGradientEllipse(
	pdf *gopdf.GoPdf,
	grad *shapes.ShapeGradientFill,
	x, y, w, h float64,
) bool {
	stops := gradientStopsFromFill(grad)
	if len(stops) == 0 {
		return false
	}
	angleDeg := 0.0
	if grad != nil && grad.AngleDeg != nil {
		angleDeg = float64(*grad.AngleDeg)
	}
	vertical := isMostlyVerticalGradient(angleDeg)
	bands := gradientBands
	cx, cy := x+w/2, y+h/2
	rx, ry := w/2, h/2

	for i := range bands {
		t0 := float64(i) / float64(bands)
		t1 := float64(i+1) / float64(bands)
		c := interpolateGradient(stops, (t0+t1)/2)
		rgb := blendOverWhite(c, c.alpha)
		pdf.SetFillColor(rgb.r, rgb.g, rgb.b)
		if vertical {
			yy0 := y + h*t0
			yy1 := y + h*t1
			ymid := (yy0 + yy1) / 2
			xw := rx * math.Sqrt(math.Max(0, 1-math.Pow((ymid-cy)/ry, 2)))
			if xw > 0 {
				pdf.RectFromUpperLeftWithStyle(cx-xw, yy0, 2*xw, yy1-yy0, "F")
			}
		} else {
			xx0 := x + w*t0
			xx1 := x + w*t1
			xmid := (xx0 + xx1) / 2
			yw := ry * math.Sqrt(math.Max(0, 1-math.Pow((xmid-cx)/rx, 2)))
			if yw > 0 {
				pdf.RectFromUpperLeftWithStyle(xx0, cy-yw, xx1-xx0, 2*yw, "F")
			}
		}
	}
	return true
}

func renderPDFLinearGradientRect(
	pdf *gopdf.GoPdf,
	grad *shapes.ShapeGradientFill,
	x, y, w, h float64,
) bool {
	stops := gradientStopsFromFill(grad)
	if len(stops) == 0 {
		return false
	}

	angleDeg := 0.0
	if grad != nil && grad.AngleDeg != nil {
		angleDeg = float64(*grad.AngleDeg)
	}
	vertical := isMostlyVerticalGradient(angleDeg)
	bands := gradientBands
	bands = max(bands, 8)

	for i := range bands {
		t0 := float64(i) / float64(bands)
		t1 := float64(i+1) / float64(bands)
		c := interpolateGradient(stops, (t0+t1)/2)
		// Blend alpha over white page background.
		rgb := blendOverWhite(c, c.alpha)
		pdf.SetFillColor(rgb.r, rgb.g, rgb.b)
		if vertical {
			yy := y + h*t0
			hh := h * (t1 - t0)
			pdf.RectFromUpperLeftWithStyle(x, yy, w, hh, "F")
		} else {
			xx := x + w*t0
			ww := w * (t1 - t0)
			pdf.RectFromUpperLeftWithStyle(xx, y, ww, h, "F")
		}
	}
	return true
}

func gradientStopsFromFill(grad *shapes.ShapeGradientFill) []gradientStop {
	if grad == nil || len(grad.Stops) == 0 {
		return nil
	}
	out := make([]gradientStop, 0, len(grad.Stops))
	for _, stop := range grad.Stops {
		r, g, b, ok := resolveOOXMLColorToken(stop.Color)
		if !ok {
			r, g, b = 0, 0, 0
		}
		alpha := 1.0
		if stop.Transparency != nil {
			alpha = 1.0 - *stop.Transparency
		}
		out = append(out, gradientStop{
			pos:   clamp01(float64(stop.PositionPct) / 100.0),
			color: rgbColor{r: r, g: g, b: b},
			alpha: clamp01(alpha),
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].pos < out[j].pos })
	return out
}

func interpolateGradient(stops []gradientStop, t float64) gradientStop {
	if len(stops) == 0 {
		return gradientStop{color: rgbColor{r: 0, g: 0, b: 0}, alpha: 1}
	}
	tt := clamp01(t)
	if tt <= stops[0].pos {
		return stops[0]
	}
	last := stops[len(stops)-1]
	if tt >= last.pos {
		return last
	}
	for i := range len(stops) - 1 {
		a := stops[i]
		b := stops[i+1]
		if tt < a.pos || tt > b.pos {
			continue
		}
		span := b.pos - a.pos
		if span <= 0 {
			return b
		}
		u := (tt - a.pos) / span
		return gradientStop{
			pos: tt,
			color: rgbColor{
				r: lerpByte(a.color.r, b.color.r, u),
				g: lerpByte(a.color.g, b.color.g, u),
				b: lerpByte(a.color.b, b.color.b, u),
			},
			alpha: a.alpha + (b.alpha-a.alpha)*u,
		}
	}
	return last
}

func blendOverWhite(c gradientStop, alpha float64) rgbColor {
	a := clamp01(alpha)
	return rgbColor{
		r: clampFloatToUint8(float64(c.color.r)*a + 255*(1-a)),
		g: clampFloatToUint8(float64(c.color.g)*a + 255*(1-a)),
		b: clampFloatToUint8(float64(c.color.b)*a + 255*(1-a)),
	}
}

func lerpByte(a, b uint8, t float64) uint8 {
	v := float64(a) + (float64(b)-float64(a))*clamp01(t)
	return clampFloatToUint8(v)
}

func clampFloatToUint8(v float64) uint8 {
	rounded := math.Round(v)
	switch {
	case rounded <= 0:
		return 0
	case rounded >= 255:
		return 255
	default:
		return uint8(rounded)
	}
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func isMostlyVerticalGradient(angleDeg float64) bool {
	a := math.Mod(angleDeg, 360)
	if a < 0 {
		a += 360
	}
	// 0/180 -> horizontal, 90/270 -> vertical.
	return (a > 45 && a < 135) || (a > 225 && a < 315)
}
