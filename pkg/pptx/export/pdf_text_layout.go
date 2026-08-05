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
	for size > minSize {
		setPDFTextFontWithHint(pdf, size, bold, italic, fontHint)
		lines := wrapPDFTextWithMetrics(pdf, text, maxWidth)
		textH := float64(len(lines)) * pdfLineHeight(size)
		if textH <= maxHeight {
			return size
		}
		size--
	}
	return size
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

func breakLongToken(pdf *gopdf.GoPdf, token string, maxWidth float64) []string {
	if token == "" {
		return []string{""}
	}
	parts := make([]string, 0, max(2, utf8.RuneCountInString(token)/12))
	var b strings.Builder
	for _, r := range token {
		next := b.String() + string(r)
		if measuredWidth(pdf, next) <= maxWidth || b.Len() == 0 {
			b.WriteRune(r)
			continue
		}
		parts = append(parts, b.String())
		b.Reset()
		b.WriteRune(r)
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
	if fontSize <= 0 {
		fontSize = defaultFontSize
	}
	return float64(fontSize) * metricsForFontHint(pdf, fontHint).baselineShiftFactor()
}
