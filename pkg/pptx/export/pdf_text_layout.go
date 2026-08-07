//nolint:mnd // Text layout heuristics intentionally use small fixed seed capacities and scaling constants.
package export

import (
	"math"
	"strings"
	"unicode/utf8"

	"github.com/signintech/gopdf"
)

const minTextAutoFitSize = 10

func fitPDFTextToBoxWithMetrics(
	pdf *gopdf.GoPdf,
	text string,
	initialSize int,
	minSize int,
	bold bool,
	italic bool,
	maxWidth float64,
	maxHeight float64,
	fontHint string,
) int {
	size := initialSize
	if size <= 0 {
		size = defaultFontSize
	}
	if minSize <= 0 {
		minSize = minTextAutoFitSize
	}
	if size <= minSize {
		return size
	}

	// Text only ever grows taller as the size goes up, so the largest size that
	// fits can be found by bisection. Stepping down one point at a time re-wrapped
	// the whole string for every size in between, which on a long paragraph cost
	// tens of full layouts to answer one question.
	fits := func(candidate int) bool {
		setPDFTextFontWithHint(pdf, candidate, bold, italic, fontHint)
		lines := wrapPDFTextWithMetrics(pdf, text, maxWidth)
		return float64(len(lines))*pdfLineHeight(candidate) <= maxHeight
	}
	if fits(size) {
		return size
	}

	best := minSize
	for low, high := minSize, size-1; low <= high; {
		mid := low + (high-low)/2
		if fits(mid) {
			best = mid
			low = mid + 1
			continue
		}
		high = mid - 1
	}
	// Leave the document on the size that was chosen, as the linear scan did.
	setPDFTextFontWithHint(pdf, best, bold, italic, fontHint)
	return best
}

func wrapPDFTextWithMetrics(pdf *gopdf.GoPdf, text string, maxWidth float64) []string {
	raw := strings.TrimSpace(text)
	if raw == "" {
		return []string{""}
	}
	lines := make([]string, 0, 4)
	for paragraph := range strings.SplitSeq(raw, "\n") {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			lines = append(lines, "")
			continue
		}
		lines = append(lines, wrapParagraph(pdf, paragraph, maxWidth)...)
	}
	if len(lines) == 0 {
		return []string{""}
	}
	return lines
}

func wrapParagraph(pdf *gopdf.GoPdf, paragraph string, maxWidth float64) []string {
	words := strings.Fields(paragraph)
	if len(words) == 0 {
		return []string{""}
	}
	lines := make([]string, 0, max(4, len(words)/6))
	current := words[0]
	for _, word := range words[1:] {
		candidate := current + " " + word
		if measuredWidth(pdf, candidate) <= maxWidth {
			current = candidate
			continue
		}
		if measuredWidth(pdf, current) > maxWidth {
			lines = append(lines, breakLongToken(pdf, current, maxWidth)...)
		} else {
			lines = append(lines, current)
		}
		current = word
	}
	if measuredWidth(pdf, current) > maxWidth {
		lines = append(lines, breakLongToken(pdf, current, maxWidth)...)
	} else {
		lines = append(lines, current)
	}
	return lines
}

// breakLongToken splits a token with no break opportunity in it — a long URL, or
// a run of CJK, which carries no spaces at all — across as many lines as it
// needs.
//
// The width of the line being built is accumulated one rune at a time rather
// than re-measured from the start on every rune. Re-measuring made the cost
// quadratic in the token's length, which a CJK paragraph pays in full because
// the whole paragraph arrives here as a single token.
//
// Summing per-rune advances ignores kerning between the runes, which can only
// make the sum slightly wider than the text really is, so a line breaks a hair
// early rather than overflowing its box.
func breakLongToken(pdf *gopdf.GoPdf, token string, maxWidth float64) []string {
	if token == "" {
		return []string{""}
	}
	parts := make([]string, 0, max(2, utf8.RuneCountInString(token)/12))
	widths := make(map[rune]float64, 32)
	var b strings.Builder
	lineWidth := 0.0
	for _, r := range token {
		width, seen := widths[r]
		if !seen {
			width = measuredWidth(pdf, string(r))
			widths[r] = width
		}
		if b.Len() == 0 || lineWidth+width <= maxWidth {
			b.WriteRune(r)
			lineWidth += width
			continue
		}
		parts = append(parts, b.String())
		b.Reset()
		b.WriteRune(r)
		lineWidth = width
	}
	if b.Len() > 0 {
		parts = append(parts, b.String())
	}
	return parts
}

// measuredWidth returns the advance width of text under the font currently set
// on pdf. gopdf sums the real hmtx advances (and, since fonts are registered
// with UseKerning, the real kern pairs), so this is exact for the embedded face
// and must not be scaled by correction factors.
func measuredWidth(pdf *gopdf.GoPdf, text string) float64 {
	w, err := pdf.MeasureTextWidth(text)
	if err != nil {
		return math.MaxFloat64
	}
	return w
}

// pdfLineHeight is the height of one line box at 100% line spacing. PowerPoint
// uses a flat 1.2 x the point size here regardless of the font, so this takes no
// font hint; only the baseline inside the box is font-dependent.
func pdfLineHeight(fontSize int) float64 {
	if fontSize <= 0 {
		fontSize = defaultFontSize
	}
	return float64(fontSize) * powerPointLineBoxFactor
}

// fontBaselineShift is the correction added to a line's top Y so that gopdf's
// Cell() lands the baseline where PowerPoint puts it. gopdf already offsets by
// OS/2 typoAscender, so only the remainder is applied here.
func fontBaselineShift(pdf *gopdf.GoPdf, fontHint string, fontSize int) float64 {
	return fontBaselineShiftInLineBox(pdf, fontHint, fontSize, 0)
}

// fontBaselineShiftInLineBox is fontBaselineShift for a line drawn into a box of
// lineBoxPt points rather than the 1.2 em default, which is what a paragraph
// with line spacing other than 100% gets. A lineBoxPt of zero means the default.
func fontBaselineShiftInLineBox(pdf *gopdf.GoPdf, fontHint string, fontSize int, lineBoxPt float64) float64 {
	if fontSize <= 0 {
		fontSize = defaultFontSize
	}
	factor := powerPointLineBoxFactor
	if lineBoxPt > 0 {
		factor = lineBoxPt / float64(fontSize)
	}
	return float64(fontSize) * metricsForFontHint(pdf, fontHint).baselineShiftFactorInLineBox(factor)
}
