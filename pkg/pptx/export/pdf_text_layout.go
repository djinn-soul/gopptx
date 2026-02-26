//nolint:mnd // Text layout heuristics intentionally use small fixed seed capacities and scaling constants.
package export

import (
	"strings"

	"github.com/signintech/gopdf"
)

const (
	minTextAutoFitSize = 10
)

func fitPDFTextSize(
	pdf *gopdf.GoPdf,
	text string,
	initialSize int,
	minSize int,
	bold bool,
	italic bool,
	maxWidth float64,
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
		w, err := pdf.MeasureTextWidth(text)
		if err != nil {
			break
		}
		if w <= maxWidth {
			return size
		}
		size--
	}
	return size
}

func wrapPDFText(pdf *gopdf.GoPdf, text string, maxWidth float64) []string {
	raw := strings.TrimSpace(text)
	if raw == "" {
		return []string{""}
	}
	words := strings.Fields(raw)
	if len(words) == 0 {
		return []string{""}
	}
	lines := make([]string, 0, 4)
	current := words[0]
	for _, word := range words[1:] {
		candidate := current + " " + word
		w, err := pdf.MeasureTextWidth(candidate)
		if err != nil {
			current = candidate
			continue
		}
		if w <= maxWidth {
			current = candidate
			continue
		}
		lines = append(lines, current)
		current = word
	}
	lines = append(lines, current)
	return lines
}
