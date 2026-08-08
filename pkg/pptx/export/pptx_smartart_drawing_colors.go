//nolint:mnd // OOXML colour transforms are stated in spec-defined units (1/100000, 1/60000 degree).
package export

import (
	"math"

	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

// A cached drawing states its colours as theme slots carrying the hue,
// saturation and luminance offsets the diagram's colour style spreads across its
// nodes — that is how one accent becomes six distinct node fills in a "colorful"
// style. The theme lives in the deck, not in the cache, so the offsets are
// resolved here, at read time, and the renderer only ever sees literal RGB.

// resolveSmartArtDrawingColors returns shapes whose colour references are
// literal RGB. A reference that names a slot this theme does not define is left
// unset, so the renderer treats it as "no colour stated" rather than painting
// black.
func resolveSmartArtDrawingColors(shapes []smartart.DrawingShape, theme deckTheme) []smartart.DrawingShape {
	if len(shapes) == 0 {
		return nil
	}
	out := make([]smartart.DrawingShape, 0, len(shapes))
	for _, shape := range shapes {
		shape.Fill = resolvedColorRef(shape.Fill, theme)
		shape.Line = resolvedColorRef(shape.Line, theme)
		shape.Paragraphs = resolvedParagraphColors(shape.Paragraphs, theme)
		out = append(out, shape)
	}
	return out
}

func resolvedParagraphColors(paragraphs []smartart.DrawingParagraph, theme deckTheme) []smartart.DrawingParagraph {
	if len(paragraphs) == 0 {
		return nil
	}
	out := make([]smartart.DrawingParagraph, 0, len(paragraphs))
	for _, paragraph := range paragraphs {
		runs := make([]smartart.DrawingRun, 0, len(paragraph.Runs))
		for _, run := range paragraph.Runs {
			run.Color = resolvedColorRef(run.Color, theme)
			runs = append(runs, run)
		}
		paragraph.Runs = runs
		out = append(out, paragraph)
	}
	return out
}

// resolvedColorRef collapses one reference to a literal RGB, dropping the slot
// and the transforms now folded into it.
func resolvedColorRef(ref smartart.ColorRef, theme deckTheme) smartart.ColorRef {
	if !ref.IsSet() {
		return smartart.ColorRef{}
	}
	base, ok := smartArtColorRefBase(ref, theme)
	if !ok {
		return smartart.ColorRef{}
	}
	return smartart.ColorRef{
		SRGB:     rgbToHex(applySmartArtColorOffsets(base, ref)),
		AlphaPct: ref.AlphaPct,
	}
}

func smartArtColorRefBase(ref smartart.ColorRef, theme deckTheme) (rgbColor, bool) {
	if ref.SRGB != "" {
		return parseHexRGB(ref.SRGB)
	}
	hexValue, ok := theme.themeColor(ref.Scheme)
	if !ok {
		return rgbColor{}, false
	}
	return parseHexRGB(hexValue)
}

// applySmartArtColorOffsets applies the shade and tint the shared OOXML helper
// already knows, then the HSL modulations and offsets on top. Modulation comes
// before offset, as ECMA-376 specifies.
func applySmartArtColorOffsets(c rgbColor, ref smartart.ColorRef) rgbColor {
	c = applyColorTransforms(c, colorTransforms{shade: ref.Shade, tint: ref.Tint})
	if ref.HueOff == 0 && ref.SatOff == 0 && ref.LumOff == 0 && ref.LumMod == 0 && ref.SatMod == 0 {
		return c
	}
	hue, saturation, luminance := smartArtRGBToHSL(c)
	if ref.SatMod != 0 {
		saturation *= float64(ref.SatMod) / 100000
	}
	if ref.LumMod != 0 {
		luminance *= float64(ref.LumMod) / 100000
	}
	hue = math.Mod(hue+float64(ref.HueOff)/60000, 360)
	if hue < 0 {
		hue += 360
	}
	saturation = clampFloat(saturation+float64(ref.SatOff)/100000, 0, 1)
	luminance = clampFloat(luminance+float64(ref.LumOff)/100000, 0, 1)
	return smartArtHSLToRGB(hue, saturation, luminance)
}

// smartArtRGBToHSL converts to hue in degrees and saturation/luminance in [0,1].
func smartArtRGBToHSL(c rgbColor) (float64, float64, float64) {
	red := float64(c.r) / 255
	green := float64(c.g) / 255
	blue := float64(c.b) / 255
	maxC := math.Max(red, math.Max(green, blue))
	minC := math.Min(red, math.Min(green, blue))
	luminance := (maxC + minC) / 2
	if maxC == minC {
		return 0, 0, luminance
	}
	delta := maxC - minC
	saturation := delta / (2 - maxC - minC)
	if luminance < 0.5 {
		saturation = delta / (maxC + minC)
	}
	var hue float64
	switch maxC {
	case red:
		hue = math.Mod((green-blue)/delta, 6)
	case green:
		hue = (blue-red)/delta + 2
	default:
		hue = (red-green)/delta + 4
	}
	hue *= 60
	if hue < 0 {
		hue += 360
	}
	return hue, saturation, luminance
}

func smartArtHSLToRGB(hue, saturation, luminance float64) rgbColor {
	if saturation <= 0 {
		return rgbColor{r: colorByte(luminance), g: colorByte(luminance), b: colorByte(luminance)}
	}
	q := luminance + saturation - luminance*saturation
	if luminance < 0.5 {
		q = luminance * (1 + saturation)
	}
	p := 2*luminance - q
	third := 1.0 / 3.0
	return rgbColor{
		r: colorByte(hueToChannel(p, q, hue/360+third)),
		g: colorByte(hueToChannel(p, q, hue/360)),
		b: colorByte(hueToChannel(p, q, hue/360-third)),
	}
}

func hueToChannel(p, q, t float64) float64 {
	t = math.Mod(t, 1)
	if t < 0 {
		t++
	}
	switch {
	case t < 1.0/6.0:
		return p + (q-p)*6*t
	case t < 0.5:
		return q
	case t < 2.0/3.0:
		return p + (q-p)*(2.0/3.0-t)*6
	default:
		return p
	}
}

func colorByte(v float64) uint8 {
	return uint8(math.Round(clampFloat(v, 0, 1) * 255))
}
