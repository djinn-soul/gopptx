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

// paragraphSpacingDefaults is the spacing PowerPoint falls back to when neither
// the paragraph nor the slide master states its own. The values differ by the
// kind of text box, so callers pass the set that matches what they are drawing.
type paragraphSpacingDefaults struct {
	LineSpacingFactor float64
	SpaceBeforePt     float64
}

// bodyPlaceholderSpacing is what PowerPoint applies inside a body placeholder.
// Both values come from the Office theme's bodyStyle lvl1pPr defaults
// (lnSpc spcPct 90000, spcBef spcPts 1000) and were confirmed against
// PowerPoint's own render: 18pt bullets come out on a 29.2pt paragraph pitch,
// which is a 19.4pt line (21.6 x 0.9) plus 10pt of space-before.
const (
	bodyDefaultLineSpacingFactor = 0.9
	bodyDefaultSpaceBeforePt     = 10.0
)

// titleDefaultLineSpacingFactor is the line spacing a title placeholder
// inherits from the Office theme's titleStyle (lnSpc spcPct 90000). Confirmed
// against PowerPoint: a 44pt title in a centre-anchored placeholder sits 2.5pt
// higher than the same text drawn at 100% spacing.
const titleDefaultLineSpacingFactor = 0.9

func bodyPlaceholderSpacing() paragraphSpacingDefaults {
	return paragraphSpacingDefaults{
		LineSpacingFactor: bodyDefaultLineSpacingFactor,
		SpaceBeforePt:     bodyDefaultSpaceBeforePt,
	}
}

// shapeTextSpacing is what PowerPoint applies to text in an ordinary shape or
// text box, which inherits the theme's otherStyle rather than bodyStyle: single
// spacing and no space between paragraphs.
func shapeTextSpacing() paragraphSpacingDefaults {
	return paragraphSpacingDefaults{LineSpacingFactor: 1.0, SpaceBeforePt: 0}
}

func paragraphLineSpacingFactor(style text.ParagraphStyle, defaults paragraphSpacingDefaults) float64 {
	if style.LineSpacingPct <= 0 {
		return defaults.LineSpacingFactor
	}
	return math.Max(float64(style.LineSpacingPct)/percentScale, minLineSpacing)
}

func paragraphStartGap(
	index int,
	prevAfter float64,
	style text.ParagraphStyle,
	defaults paragraphSpacingDefaults,
) float64 {
	before := float64(max(style.SpaceBeforePt, 0))
	if before == 0 {
		before = defaults.SpaceBeforePt
	}
	// PowerPoint does not lead the first paragraph in a box with space-before,
	// so the text starts flush against the top inset.
	if index == 0 {
		return 0
	}
	return math.Max(prevAfter, before)
}

func paragraphAfterGap(style text.ParagraphStyle) float64 {
	return float64(max(style.SpaceAfterPt, 0))
}

func bulletPrefix(style text.ParagraphStyle, idx int) string {
	switch text.NormalizeBulletStyle(style.BulletStyle) {
	case "":
		// Unset means the paragraph never asked for a bullet. This used to fall
		// through to the default and mark every plain shape paragraph with a
		// bullet, which PowerPoint draws on top of the first letter. The bullet
		// list renderer supplies its own default for the body placeholder, where
		// a bullet is the expected look.
		return ""
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
