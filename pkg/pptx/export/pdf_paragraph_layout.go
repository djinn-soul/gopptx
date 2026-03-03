package export

import (
	"fmt"
	"math"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

const (
	percentScale     = 100.0
	alphabetRuneSpan = 26
	minLineSpacing   = 0.6
)

func paragraphLineSpacingFactor(style text.ParagraphStyle) float64 {
	if style.LineSpacingPct <= 0 {
		return 1.0
	}
	return math.Max(float64(style.LineSpacingPct)/percentScale, minLineSpacing)
}

func paragraphStartGap(index int, prevAfter float64, style text.ParagraphStyle) float64 {
	before := float64(max(style.SpaceBeforePt, 0))
	if index == 0 {
		return before
	}
	return math.Max(prevAfter, before)
}

func paragraphAfterGap(style text.ParagraphStyle) float64 {
	return float64(max(style.SpaceAfterPt, 0))
}

func bulletPrefix(style text.ParagraphStyle, idx int) string {
	switch text.NormalizeBulletStyle(style.BulletStyle) {
	case text.BulletStyleNone:
		return ""
	case text.BulletStyleNumber:
		return fmt.Sprintf("%d.", idx+1)
	case text.BulletStyleLetterLower:
		return fmt.Sprintf("%c.", 'a'+(idx%alphabetRuneSpan))
	case text.BulletStyleLetterUpper:
		return fmt.Sprintf("%c.", 'A'+(idx%alphabetRuneSpan))
	case text.BulletStyleRomanLower:
		return strings.ToLower(romanNumeral(idx + 1))
	case text.BulletStyleRomanUpper:
		return romanNumeral(idx + 1)
	case text.BulletStyleCustom:
		if style.BulletChar != "" {
			return style.BulletChar
		}
		return "•"
	default:
		return "•"
	}
}

func romanNumeral(n int) string {
	if n <= 0 {
		return ""
	}
	table := []struct {
		value int
		sym   string
	}{
		{1000, "M"}, {900, "CM"}, {500, "D"}, {400, "CD"},
		{100, "C"}, {90, "XC"}, {50, "L"}, {40, "XL"},
		{10, "X"}, {9, "IX"}, {5, "V"}, {4, "IV"}, {1, "I"},
	}
	var out strings.Builder
	for _, entry := range table {
		for n >= entry.value {
			out.WriteString(entry.sym)
			n -= entry.value
		}
	}
	return out.String()
}
