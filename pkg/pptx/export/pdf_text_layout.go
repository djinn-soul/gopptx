//nolint:mnd // Text layout heuristics intentionally use small fixed seed capacities and scaling constants.
package export

import (
	"math"
	"strings"
	"unicode/utf8"

	"github.com/signintech/gopdf"
)

const minTextAutoFitSize = 10

func fitPDFTextToBox(
	pdf *gopdf.GoPdf,
	text string,
	initialSize int,
	minSize int,
	bold bool,
	italic bool,
	maxWidth float64,
	maxHeight float64,
) int {
	return fitPDFTextToBoxWithMetrics(pdf, text, initialSize, minSize, bold, italic, maxWidth, maxHeight, "")
}

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
		setPDFTextFont(pdf, size, bold, italic)
		lines := wrapPDFTextWithMetrics(pdf, text, maxWidth, fontHint)
		textH := float64(len(lines)) * pdfLineHeight(size)
		if textH <= maxHeight {
			return size
		}
		size--
	}
	return size
}

func wrapPDFText(pdf *gopdf.GoPdf, text string, maxWidth float64) []string {
	return wrapPDFTextWithMetrics(pdf, text, maxWidth, "")
}

func wrapPDFTextWithMetrics(pdf *gopdf.GoPdf, text string, maxWidth float64, fontHint string) []string {
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
		lines = append(lines, wrapParagraph(pdf, paragraph, maxWidth, fontHint)...)
	}
	if len(lines) == 0 {
		return []string{""}
	}
	return lines
}

func wrapParagraph(pdf *gopdf.GoPdf, paragraph string, maxWidth float64, fontHint string) []string {
	words := strings.Fields(paragraph)
	if len(words) == 0 {
		return []string{""}
	}
	lines := make([]string, 0, max(4, len(words)/6))
	current := words[0]
	for _, word := range words[1:] {
		candidate := current + " " + word
		if measuredWidthWithMetrics(pdf, candidate, fontHint) <= maxWidth {
			current = candidate
			continue
		}
		if measuredWidthWithMetrics(pdf, current, fontHint) > maxWidth {
			lines = append(lines, breakLongToken(pdf, current, maxWidth, fontHint)...)
		} else {
			lines = append(lines, current)
		}
		current = word
	}
	if measuredWidthWithMetrics(pdf, current, fontHint) > maxWidth {
		lines = append(lines, breakLongToken(pdf, current, maxWidth, fontHint)...)
	} else {
		lines = append(lines, current)
	}
	return lines
}

func breakLongToken(pdf *gopdf.GoPdf, token string, maxWidth float64, fontHint string) []string {
	if token == "" {
		return []string{""}
	}
	parts := make([]string, 0, max(2, utf8.RuneCountInString(token)/12))
	var b strings.Builder
	for _, r := range token {
		next := b.String() + string(r)
		if measuredWidthWithMetrics(pdf, next, fontHint) <= maxWidth || b.Len() == 0 {
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

func measuredWidth(pdf *gopdf.GoPdf, text string) float64 {
	w, err := pdf.MeasureTextWidth(text)
	if err != nil {
		return math.MaxFloat64
	}
	return w
}

func measuredWidthWithMetrics(pdf *gopdf.GoPdf, text string, fontHint string) float64 {
	base := measuredWidth(pdf, text)
	if base == math.MaxFloat64 || text == "" {
		return base
	}
	factor := fontWidthFactor(fontHint)
	kerning := kerningAdjustment(text, fontHint)
	return (base * factor) + kerning
}

func fontWidthFactor(fontHint string) float64 {
	name := strings.ToLower(strings.TrimSpace(fontHint))
	switch {
	case strings.Contains(name, "calibri"):
		return 0.97
	case strings.Contains(name, "times"):
		return 0.95
	case strings.Contains(name, "courier"), strings.Contains(name, "mono"), strings.Contains(name, "consolas"):
		return 1.02
	default:
		return 1.0
	}
}

func kerningAdjustment(text string, fontHint string) float64 {
	name := strings.ToLower(strings.TrimSpace(fontHint))
	if strings.Contains(name, "mono") || strings.Contains(name, "courier") || strings.Contains(name, "consolas") {
		return 0
	}
	adj := 0.0
	prev := rune(0)
	for _, cur := range text {
		if prev == 0 {
			prev = cur
			continue
		}
		switch string([]rune{prev, cur}) {
		case "To", "Ta", "Te", "Yo", "VA", "WA", "LT", "Ty", "AV":
			adj -= 0.22
		case "ll", "ii", "rr":
			adj -= 0.08
		}
		prev = cur
	}
	return adj
}

func pdfLineHeight(fontSize int) float64 {
	if fontSize <= 0 {
		fontSize = defaultFontSize
	}
	return math.Max(float64(fontSize)*1.18, 12)
}

func fontBaselineShift(fontHint string, fontSize int) float64 {
	if fontSize <= 0 {
		fontSize = defaultFontSize
	}
	name := strings.ToLower(strings.TrimSpace(fontHint))
	switch {
	case strings.Contains(name, "calibri"):
		return float64(fontSize) * 0.07
	case strings.Contains(name, "times"):
		return float64(fontSize) * 0.08
	case strings.Contains(name, "mono"), strings.Contains(name, "courier"), strings.Contains(name, "consolas"):
		return float64(fontSize) * 0.04
	default:
		return float64(fontSize) * 0.06
	}
}
