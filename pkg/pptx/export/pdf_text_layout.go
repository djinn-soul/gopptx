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
		lines := wrapPDFTextWithMetrics(pdf, text, maxWidth, fontHint)
		textH := float64(len(lines)) * pdfLineHeight(size)
		if textH <= maxHeight {
			return size
		}
		size--
	}
	return size
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
	profile := fontMetricProfile(fontHint)
	kerning := kerningAdjustment(text, profile, base)
	return (base * profile.WidthFactor) + kerning
}

type fontMetricsProfile struct {
	WidthFactor      float64
	SpaceFactor      float64
	NarrowRuneFactor float64
	WideRuneFactor   float64
	KernPairFactor   float64
}

func fontMetricProfile(fontHint string) fontMetricsProfile {
	name := strings.ToLower(strings.TrimSpace(fontHint))
	switch {
	case strings.Contains(name, "calibri"):
		return fontMetricsProfile{
			WidthFactor:      0.972,
			SpaceFactor:      0.88,
			NarrowRuneFactor: 0.78,
			WideRuneFactor:   1.12,
			KernPairFactor:   0.16,
		}
	case strings.Contains(name, "times"):
		return fontMetricsProfile{
			WidthFactor:      0.955,
			SpaceFactor:      0.84,
			NarrowRuneFactor: 0.74,
			WideRuneFactor:   1.15,
			KernPairFactor:   0.20,
		}
	case strings.Contains(name, "courier"), strings.Contains(name, "mono"), strings.Contains(name, "consolas"):
		return fontMetricsProfile{
			WidthFactor:      1.00,
			SpaceFactor:      1.00,
			NarrowRuneFactor: 1.00,
			WideRuneFactor:   1.00,
			KernPairFactor:   0,
		}
	default:
		return fontMetricsProfile{
			WidthFactor:      1.00,
			SpaceFactor:      0.92,
			NarrowRuneFactor: 0.82,
			WideRuneFactor:   1.10,
			KernPairFactor:   0.14,
		}
	}
}

func kerningAdjustment(text string, profile fontMetricsProfile, measured float64) float64 {
	runes := []rune(text)
	if len(runes) < 2 {
		return 0
	}
	perRune := measured / math.Max(float64(len(runes)), 1)
	adj := 0.0
	for i, cur := range runes {
		switch {
		case cur == ' ':
			adj += (profile.SpaceFactor - 1.0) * perRune
		case isNarrowRune(cur):
			adj += (profile.NarrowRuneFactor - 1.0) * perRune * 0.5
		case isWideRune(cur):
			adj += (profile.WideRuneFactor - 1.0) * perRune * 0.5
		}
		if i == 0 || profile.KernPairFactor == 0 {
			continue
		}
		if isTightPair(runes[i-1], cur) {
			adj -= perRune * profile.KernPairFactor
		}
	}
	return adj
}

func isNarrowRune(r rune) bool {
	return strings.ContainsRune("iljftI1|!:;.,'`", r)
}

func isWideRune(r rune) bool {
	return strings.ContainsRune("MWQ@#%&8", r)
}

func isTightPair(prev, cur rune) bool {
	return strings.ContainsRune("TAVWLYF", prev) && strings.ContainsRune("aoeu.,", cur)
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
