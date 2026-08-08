//nolint:mnd // Shape effects are approximation constants tuned for native PDF parity.
package export

import (
	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

const (
	shapeShadowOffsetPt   = 3.0
	shapeShadowAlpha      = 0.22
	shapeGlowAlpha        = 0.20
	shapeReflectionAlpha  = 0.22
	shapeReflectionGapPt  = 2.0
	shapeReflectionHeight = 0.28
	shapeSoftEdgesAlpha   = 0.85
	// gopdf names its blend modes with the PDF dictionary spelling, slash and
	// all. This read "Normal", which no blend mode matches, so every call that
	// used it failed and returned the zero transparency — silently, because the
	// error is swallowed. Soft edges, shadows, glows, reflections and the
	// translucent Venn circles were all being dropped on that.
	shapeEffectsBlendMode       = "/Normal"
	shapeDefaultAccentR   uint8 = 68
	shapeDefaultAccentG   uint8 = 114
	shapeDefaultAccentB   uint8 = 196
)

func applyPDFShapeSoftEdges(pdf *gopdf.GoPdf, s shapes.Shape) bool {
	if s.Effects == nil || !s.Effects.SoftEdges {
		return false
	}
	alpha, err := gopdf.NewTransparency(shapeSoftEdgesAlpha, shapeEffectsBlendMode)
	if err != nil {
		return false
	}
	_ = pdf.SetTransparency(alpha)
	return true
}

// applyPDFShapeFillAlpha makes a part-transparent fill actually paint that way,
// and reports whether it set anything the caller has to clear.
//
// ShapeFill has carried a Transparency since long before this, and the native
// renderer ignored it: a fill the deck asked to show through painted solid.
// Soft edges already own the transparency state when they are on, so the fill's
// own alpha defers to them rather than overwriting it.
func applyPDFShapeFillAlpha(pdf *gopdf.GoPdf, s shapes.Shape, softEdgesApplied bool) bool {
	if softEdgesApplied || s.Fill == nil || s.Fill.Transparency == nil {
		return false
	}
	opacity := 1 - *s.Fill.Transparency
	if opacity >= 1 || opacity < 0 {
		return false
	}
	alpha, err := gopdf.NewTransparency(opacity, shapeEffectsBlendMode)
	if err != nil {
		return false
	}
	_ = pdf.SetTransparency(alpha)
	return true
}

func renderPDFShapeEffects(
	pdf *gopdf.GoPdf,
	s shapes.Shape,
	x, y, w, h float64,
	hasFill bool,
) {
	if s.Effects == nil {
		return
	}
	if s.Effects.Shadow {
		renderPDFShapeShadow(pdf, s, x, y, w, h)
	}
	if s.Effects.Glow {
		renderPDFShapeGlow(pdf, s, x, y, w, h)
	}
	if s.Effects.Reflection {
		renderPDFShapeReflection(pdf, s, x, y, w, h, hasFill)
	}
}

func renderPDFShapeShadow(pdf *gopdf.GoPdf, s shapes.Shape, x, y, w, h float64) {
	alpha, err := gopdf.NewTransparency(shapeShadowAlpha, shapeEffectsBlendMode)
	if err == nil {
		_ = pdf.SetTransparency(alpha)
	}
	pdf.SetFillColor(0, 0, 0)
	drawPDFGeometry(pdf, s, x+shapeShadowOffsetPt, y+shapeShadowOffsetPt, w, h, "F")
	if err == nil {
		pdf.ClearTransparency()
	}
}

func renderPDFShapeGlow(pdf *gopdf.GoPdf, s shapes.Shape, x, y, w, h float64) {
	r, g, b := shapeGlowColor(s)
	alpha, err := gopdf.NewTransparency(shapeGlowAlpha, shapeEffectsBlendMode)
	if err == nil {
		_ = pdf.SetTransparency(alpha)
	}
	pdf.SetStrokeColor(r, g, b)
	pdf.SetLineWidth(4)
	drawPDFGeometry(pdf, s, x, y, w, h, "D")
	if err == nil {
		pdf.ClearTransparency()
	}
}

func shapeGlowColor(s shapes.Shape) (uint8, uint8, uint8) {
	if s.Line != nil && s.Line.Color != "" {
		return hexToRGB(s.Line.Color)
	}
	if s.Fill != nil && s.Fill.Color != "" {
		return hexToRGB(s.Fill.Color)
	}
	if s.GradientFill != nil && len(s.GradientFill.Stops) > 0 {
		return hexToRGB(s.GradientFill.Stops[0].Color)
	}
	return shapeDefaultAccentR, shapeDefaultAccentG, shapeDefaultAccentB
}

func renderPDFShapeReflection(
	pdf *gopdf.GoPdf,
	s shapes.Shape,
	x, y, w, h float64,
	hasFill bool,
) {
	refH := h * shapeReflectionHeight
	if refH <= 1 || !hasFill {
		return
	}
	alpha, err := gopdf.NewTransparency(shapeReflectionAlpha, shapeEffectsBlendMode)
	if err != nil {
		return
	}
	_ = pdf.SetTransparency(alpha)
	style := "F"
	if s.Line != nil && s.Line.Width > 0 {
		style = "DF"
	}
	drawPDFGeometry(pdf, s, x, y+h+shapeReflectionGapPt, w, refH, style)
	pdf.ClearTransparency()
}
