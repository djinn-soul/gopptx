package export

import (
	"math"
	"strings"
	"unicode"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

// Character-level decorations PowerPoint draws itself rather than leaving to the
// font: underline, strike-through, and the size/offset of script runs.
//
// The geometry comes from the font's own post and OS/2 tables (see
// pdf_font_metrics.go) so an underline sits where the type designer put it
// instead of at a guessed fraction of the point size.

const (
	// smallCapsSizeFactor is the size lowercase letters take in a small-caps
	// run. PowerPoint draws them at 80% of the run's size, uppercased.
	smallCapsSizeFactor = 0.8

	// Gap between the two lines of a double underline or double strike, as a
	// fraction of the point size.
	doubleLineGapFactor = 0.06

	// Dotted underlines are drawn as a dash pattern of this many multiples of
	// the line thickness.
	dottedDashFactor = 2.0

	minDecorationThicknessPt = 0.35

	// gopdf's line types; anything else draws solid.
	pdfLineTypeDotted = "dotted"
	pdfLineTypeDashed = "dashed"

	// heavyDecorationWeight is how much thicker OOXML's "Heavy" underline and
	// strike variants are drawn than their plain counterparts.
	heavyDecorationWeight = 2.0

	// Wavy underline geometry, as multiples of the stroke thickness, with floors
	// so a hairline rule still visibly waves.
	wavyAmplitudeFactor  = 1.2
	wavyWavelengthFactor = 4.0
	wavySegmentsPerWave  = 8
	minWavyAmplitudePt   = 0.6
	minWavyWavelengthPt  = 2.4
)

// effectiveRunSizePt is the point size a run is actually drawn at. It differs
// from the run's stated size for script and small-caps runs, both of which
// PowerPoint renders smaller than the surrounding text.
func effectiveRunSizePt(pdf *gopdf.GoPdf, run pdfStyledRun) int {
	size := runSizePt(run)
	metrics := metricsForFontHint(pdf, run.FontHint)
	factor := 1.0
	switch {
	case run.Superscript:
		factor = metrics.superscriptSizeFactor()
	case run.Subscript:
		factor = metrics.subscriptSizeFactor()
	case run.SmallCaps:
		factor = smallCapsSizeFactor
	}
	if factor >= 1 {
		return size
	}
	return max(int(math.Round(float64(size)*factor)), 1)
}

// scriptBaselineOffsetPt is how far a script run is moved off the line's
// baseline: negative raises a superscript, positive drops a subscript.
func scriptBaselineOffsetPt(pdf *gopdf.GoPdf, run pdfStyledRun) float64 {
	size := float64(runSizePt(run))
	metrics := metricsForFontHint(pdf, run.FontHint)
	switch {
	case run.Superscript:
		return -size * metrics.superscriptOffsetFactor()
	case run.Subscript:
		return size * metrics.subscriptOffsetFactor()
	default:
		return 0
	}
}

// gopdfBaselineDropPt is how far below the Y it is given gopdf's Cell() puts the
// baseline: the font's own typoAscender.
func gopdfBaselineDropPt(pdf *gopdf.GoPdf, run pdfStyledRun) float64 {
	size := effectiveRunSizePt(pdf, run)
	return float64(size) * metricsForFontHint(pdf, run.FontHint).typoAscenderFactor()
}

// runBaselineOffset is where one run's font wants its baseline, measured down
// from the top of the line box.
func runBaselineOffset(pdf *gopdf.GoPdf, run pdfStyledRun) float64 {
	size := effectiveRunSizePt(pdf, run)
	return fontBaselineShift(pdf, run.FontHint, size) + gopdfBaselineDropPt(pdf, run)
}

// styledLineBaselineOffset is the single baseline every run of a line is drawn
// on, measured down from the top of the line box. PowerPoint puts a line's runs
// on one baseline no matter how their fonts differ, so the deepest run's
// baseline is the one the whole line adopts. Script runs are excluded: they are
// offset from the line's baseline rather than defining it.
func styledLineBaselineOffset(pdf *gopdf.GoPdf, line []pdfStyledRun) float64 {
	baseline, scriptBaseline := 0.0, 0.0
	for _, run := range line {
		offset := runBaselineOffset(pdf, run)
		if run.Subscript || run.Superscript {
			scriptBaseline = math.Max(scriptBaseline, offset)
			continue
		}
		baseline = math.Max(baseline, offset)
	}
	if baseline == 0 {
		return scriptBaseline
	}
	return baseline
}

// renderPDFStyledRunDecorations draws the run's underline and strike-through.
// baseline is the y of the run's own baseline, already shifted for a script run.
func renderPDFStyledRunDecorations(pdf *gopdf.GoPdf, run pdfStyledRun, x, baseline float64) {
	underline := normalizeDecoration(run.Underline)
	strike := normalizeDecoration(run.Strikethrough)
	if underline == "" && strike == "" {
		return
	}
	width := measureStyledRunWidth(pdf, run)
	if width <= 0 {
		return
	}

	size := float64(effectiveRunSizePt(pdf, run))
	metrics := metricsForFontHint(pdf, run.FontHint)
	pdf.SetStrokeColor(run.Color[0], run.Color[1], run.Color[2])

	if underline != "" {
		drawTextDecorationLine(
			pdf, textDecoration{
				x:         x,
				width:     width,
				y:         baseline + size*metrics.underlineOffsetFactor(),
				thickness: size * metrics.underlineThicknessFactor(),
				gap:       size * doubleLineGapFactor,
				style:     underline,
				// A second underline is drawn further from the baseline.
				downward: true,
			},
		)
	}
	if strike != "" {
		drawTextDecorationLine(
			pdf, textDecoration{
				x:         x,
				width:     width,
				y:         baseline - size*metrics.strikeoutOffsetFactor(),
				thickness: size * metrics.strikeoutThicknessFactor(),
				gap:       size * doubleLineGapFactor,
				style:     strike,
			},
		)
	}
	pdf.SetStrokeColor(0, 0, 0)
	pdf.SetLineType("")
	pdf.SetLineWidth(1)
}

type textDecoration struct {
	x         float64
	width     float64
	y         float64
	thickness float64
	gap       float64
	style     string
	downward  bool
}

func drawTextDecorationLine(pdf *gopdf.GoPdf, d textDecoration) {
	thickness := math.Max(d.thickness*decorationWeight(d.style), minDecorationThicknessPt)
	pdf.SetLineWidth(thickness)
	pdf.SetLineType(decorationLineType(d.style))

	drawOneDecorationStroke(pdf, d, d.y, thickness)
	if isDoubleDecoration(d.style) {
		offset := math.Max(d.gap, thickness*dottedDashFactor)
		if !d.downward {
			offset = -offset
		}
		drawOneDecorationStroke(pdf, d, d.y+offset, thickness)
	}
	pdf.SetLineType("")
}

// drawOneDecorationStroke draws a single rule at y: a straight line, or a sine
// wave for the wavy styles, which have no line-type equivalent and would
// otherwise be indistinguishable from a plain underline.
func drawOneDecorationStroke(pdf *gopdf.GoPdf, d textDecoration, y, thickness float64) {
	if !isWavyDecoration(d.style) {
		pdf.Line(d.x, y, d.x+d.width, y)
		return
	}
	amplitude := math.Max(thickness*wavyAmplitudeFactor, minWavyAmplitudePt)
	wavelength := math.Max(amplitude*wavyWavelengthFactor, minWavyWavelengthPt)
	steps := max(int(math.Ceil(d.width/wavelength*wavySegmentsPerWave)), wavySegmentsPerWave)

	prevX, prevY := d.x, y
	for i := 1; i <= steps; i++ {
		t := float64(i) / float64(steps)
		nextX := d.x + d.width*t
		nextY := y + amplitude*math.Sin(2*math.Pi*d.width*t/wavelength)
		pdf.Line(prevX, prevY, nextX, nextY)
		prevX, prevY = nextX, nextY
	}
}

// decorationLineType maps an OOXML underline style to a gopdf line type. The
// heavy and wavy variants are not line types: weight and wave are applied by
// decorationWeight and drawOneDecorationStroke respectively.
func decorationLineType(style string) string {
	switch style {
	case text.UnderlineStyleDotted, "dottedHeavy", "dotDash", "dotDashHeavy",
		"dotDotDash", "dotDotDashHeavy":
		return pdfLineTypeDotted
	case shapes.LineDashDash, "dashHeavy", "dashLong", "dashLongHeavy":
		return pdfLineTypeDashed
	default:
		return ""
	}
}

// decorationWeight is the stroke-width multiplier of a style. OOXML's "Heavy"
// variants are the same rule drawn thicker, which is how PowerPoint tells
// sngHeavy from sng.
func decorationWeight(style string) float64 {
	if strings.Contains(strings.ToLower(style), "heavy") {
		return heavyDecorationWeight
	}
	return 1
}

func isWavyDecoration(style string) bool {
	return strings.HasPrefix(style, "wavy")
}

func isDoubleDecoration(style string) bool {
	return style == text.UnderlineStyleDouble || style == text.StrikethroughStyleDouble
}

// normalizeDecoration returns "" for an absent or explicitly disabled
// decoration, so callers can test the result directly.
func normalizeDecoration(style string) string {
	trimmed := strings.TrimSpace(style)
	if trimmed == "" || trimmed == text.UnderlineStyleNone {
		return ""
	}
	return trimmed
}

// applyRunCapsTransform returns the text as PowerPoint draws it under a:caps.
// Small caps are not handled here: they change size as well as case and are
// split into their own runs by splitSmallCapsRun.
func applyRunCapsTransform(value string, allCaps bool) string {
	if !allCaps {
		return value
	}
	return strings.ToUpper(value)
}

// splitSmallCapsRun turns one small-caps run into the runs PowerPoint actually
// draws: letters that were already capitals keep the run's size, and everything
// else is uppercased and drawn at smallCapsSizeFactor. Rendering the whole run
// small would shrink the capitals too, which is visibly wrong for a name like
// "McDonald".
func splitSmallCapsRun(run pdfStyledRun) []pdfStyledRun {
	if !run.SmallCaps || run.Text == "" {
		return []pdfStyledRun{run}
	}
	out := make([]pdfStyledRun, 0, 2)
	var segment strings.Builder
	segmentIsSmall := false

	flush := func() {
		if segment.Len() == 0 {
			return
		}
		part := run
		part.Text = strings.ToUpper(segment.String())
		part.SmallCaps = segmentIsSmall
		out = append(out, part)
		segment.Reset()
	}
	for _, r := range run.Text {
		// Anything that is not already a capital is drawn as a small capital,
		// which is what PowerPoint does with digits and punctuation too: they
		// are unaffected by case, so they stay at the run's own size.
		small := unicode.IsLower(r)
		if segment.Len() > 0 && small != segmentIsSmall {
			flush()
		}
		segmentIsSmall = small
		segment.WriteRune(r)
	}
	flush()
	if len(out) == 0 {
		return []pdfStyledRun{run}
	}
	return out
}
