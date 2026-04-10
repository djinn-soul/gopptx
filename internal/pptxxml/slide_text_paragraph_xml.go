package pptxxml

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	indentStep = 457200
	pctFactor  = 1000
	ptFactor   = 100
)

// defaultBulletParagraphProps is the precomputed output of BulletParagraphPropsXML
// for a zero-value BulletParagraphSpec (level=0, default indent, bullet char).
// Returning a constant avoids ~6 allocations per bullet in the common case.
const defaultBulletParagraphProps = `<a:pPr lvl="0" marL="457200" indent="-457200"><a:buChar char="•"/></a:pPr>`

// BulletParagraphSpec describes paragraph formatting for one bullet line.
type BulletParagraphSpec struct {
	Align             string
	SpaceBeforePt     int
	SpaceAfterPt      int
	LineSpacingPct    int
	LineSpacingPts    int
	SpaceBeforeRaw    int
	SpaceAfterRaw     int
	LineSpacingPctRaw int
	LineSpacingPtsRaw int
	BulletStyle       string
	BulletChar        string
	BulletColor       string
	BulletSize        int
	TabStops          []int64
	Level             int
	LeftIndent        int64
	RightIndent       int64
	HangingIndent     int64
	RTL               bool
}

// IsZero reports whether this style has no explicit values set.
func (s BulletParagraphSpec) IsZero() bool {
	return s.Align == "" &&
		s.SpaceBeforePt == 0 &&
		s.SpaceAfterPt == 0 &&
		s.LineSpacingPct == 0 &&
		s.LineSpacingPts == 0 &&
		s.SpaceBeforeRaw == 0 &&
		s.SpaceAfterRaw == 0 &&
		s.LineSpacingPctRaw == 0 &&
		s.LineSpacingPtsRaw == 0 &&
		s.BulletStyle == "" &&
		s.BulletChar == "" &&
		s.BulletColor == "" &&
		s.BulletSize == 0 &&
		len(s.TabStops) == 0 &&
		s.Level == 0 &&
		s.LeftIndent == 0 &&
		s.RightIndent == 0 &&
		s.HangingIndent == 0 &&
		!s.RTL
}

func bulletStyleAt(all []BulletParagraphSpec, index int) BulletParagraphSpec {
	if len(all) == 0 || index < 0 || index >= len(all) {
		return BulletParagraphSpec{}
	}
	return all[index]
}

func BulletParagraphPropsXML(style BulletParagraphSpec) string {
	// Fast path: zero-value spec → precomputed constant (zero allocs).
	if style.IsZero() {
		return defaultBulletParagraphProps
	}
	marL, indent := bulletIndent(style.Level)
	if style.LeftIndent != 0 {
		marL = int(style.LeftIndent)
	}
	if style.HangingIndent != 0 {
		indent = int(style.HangingIndent)
	}
	// Note: Right indent (marR) is also supported in a:pPr
	marRXML := ""
	if style.RightIndent != 0 {
		marRXML = ` marR="` + strconv.FormatInt(style.RightIndent, 10) + `"`
	}

	base := `<a:pPr lvl="` + strconv.Itoa(
		style.Level,
	) + `" marL="` + strconv.Itoa(
		marL,
	) + `" indent="` + strconv.Itoa(
		indent,
	) + `"` + marRXML
	if style.RTL {
		base += ` rtl="1"`
	}
	if style.Align != "" {
		base += ` algn="` + Escape(style.Align) + `"`
	}
	base += ">"

	if value, ok := paragraphPctValue(style); ok {
		base += `<a:lnSpc><a:spcPct val="` + strconv.Itoa(value) + `"/></a:lnSpc>`
	} else if value, ok := paragraphPtValue(style.LineSpacingPtsRaw, style.LineSpacingPts); ok {
		base += `<a:lnSpc><a:spcPts val="` + strconv.Itoa(value) + `"/></a:lnSpc>`
	}
	if value, ok := paragraphPtValue(style.SpaceBeforeRaw, style.SpaceBeforePt); ok {
		base += `<a:spcBef><a:spcPts val="` + strconv.Itoa(value) + `"/></a:spcBef>`
	}
	if value, ok := paragraphPtValue(style.SpaceAfterRaw, style.SpaceAfterPt); ok {
		base += `<a:spcAft><a:spcPts val="` + strconv.Itoa(value) + `"/></a:spcAft>`
	}

	base += bulletNodeXML(style)
	if len(style.TabStops) > 0 {
		var tabs strings.Builder
		tabs.WriteString(`<a:tabLst>`)
		for _, pos := range style.TabStops {
			tabs.WriteString(`<a:tab pos="`)
			tabs.WriteString(strconv.FormatInt(pos, 10))
			tabs.WriteString(`"/>`)
		}
		tabs.WriteString(`</a:tabLst>`)
		base += tabs.String()
	}

	return base + `</a:pPr>`
}

func paragraphPctValue(style BulletParagraphSpec) (int, bool) {
	if style.LineSpacingPctRaw > 0 {
		return style.LineSpacingPctRaw, true
	}
	if style.LineSpacingPct > 0 {
		return style.LineSpacingPct * pctFactor, true
	}
	return 0, false
}

func paragraphPtValue(raw int, whole int) (int, bool) {
	if raw > 0 {
		return raw, true
	}
	if whole > 0 {
		return whole * ptFactor, true
	}
	return 0, false
}

func bulletIndent(level int) (int, int) {
	indent := indentStep + (level * indentStep)
	marginLeft := level*indentStep + indent
	return marginLeft, -indent
}

func bulletNodeXML(style BulletParagraphSpec) string {
	var b strings.Builder

	if style.BulletColor != "" {
		b.WriteString(`<a:buClr><a:srgbClr val="`)
		b.WriteString(Escape(style.BulletColor))
		b.WriteString(`"/></a:buClr>`)
	}
	if style.BulletSize > 0 {
		b.WriteString(`<a:buSzPct val="`)
		b.WriteString(strconv.Itoa(style.BulletSize * pctFactor))
		b.WriteString(`"/>`)
	}

	switch normalizeBulletStyle(style.BulletStyle) {
	case "", "bullet":
		b.WriteString(`<a:buChar char="•"/>`)
	case "number":
		b.WriteString(`<a:buAutoNum type="arabicPeriod"/>`)
	case "letter_lower":
		b.WriteString(`<a:buAutoNum type="alphaLcPeriod"/>`)
	case "letter_upper":
		b.WriteString(`<a:buAutoNum type="alphaUcPeriod"/>`)
	case "roman_lower":
		b.WriteString(`<a:buAutoNum type="romanLcPeriod"/>`)
	case "roman_upper":
		b.WriteString(`<a:buAutoNum type="romanUcPeriod"/>`)
	case "custom":
		b.WriteString(`<a:buChar char="`)
		b.WriteString(Escape(strings.TrimSpace(style.BulletChar)))
		b.WriteString(`"/>`)
	case arrowTypeNone:
		b.WriteString(`<a:buNone/>`)
	default:
		panic(fmt.Sprintf("unsupported bullet style: %q", style.BulletStyle))
	}
	return b.String()
}

func normalizeBulletStyle(style string) string {
	normalized := strings.ToLower(strings.TrimSpace(style))
	normalized = strings.ReplaceAll(normalized, "-", "_")
	normalized = strings.ReplaceAll(normalized, " ", "_")
	switch normalized {
	case "numbered":
		return "number"
	case "lettered", "letter", "letterlower", "alphalower":
		return "letter_lower"
	case "letterupper", "alphaupper":
		return "letter_upper"
	case "roman", "romanupper":
		return "roman_upper"
	case "romanlower":
		return "roman_lower"
	default:
		return normalized
	}
}
