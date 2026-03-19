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

// defaultBulletParagraphProps is the precomputed output of bulletParagraphPropsXML
// for a zero-value BulletParagraphSpec (level=0, default indent, bullet char).
// Returning a constant avoids ~6 allocations per bullet in the common case.
const defaultBulletParagraphProps = `<a:pPr lvl="0" marL="457200" indent="-457200"><a:buChar char="•"/></a:pPr>`

// BulletParagraphSpec describes paragraph formatting for one bullet line.
type BulletParagraphSpec struct {
	Align          string
	SpaceBeforePt  int
	SpaceAfterPt   int
	LineSpacingPct int
	BulletStyle    string
	BulletChar     string
	BulletColor    string
	BulletSize     int
	Level          int
	LeftIndent     int64
	RightIndent    int64
	HangingIndent  int64
	RTL            bool
}

func bulletStyleAt(all []BulletParagraphSpec, index int) BulletParagraphSpec {
	if len(all) == 0 || index < 0 || index >= len(all) {
		return BulletParagraphSpec{}
	}
	return all[index]
}

func bulletParagraphPropsXML(style BulletParagraphSpec) string {
	// Fast path: zero-value spec → precomputed constant (zero allocs).
	if style == (BulletParagraphSpec{}) {
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

	if style.LineSpacingPct > 0 {
		base += `<a:lnSpc><a:spcPct val="` + strconv.Itoa(style.LineSpacingPct*pctFactor) + `"/></a:lnSpc>`
	}
	if style.SpaceBeforePt > 0 {
		base += `<a:spcBef><a:spcPts val="` + strconv.Itoa(style.SpaceBeforePt*ptFactor) + `"/></a:spcBef>`
	}
	if style.SpaceAfterPt > 0 {
		base += `<a:spcAft><a:spcPts val="` + strconv.Itoa(style.SpaceAfterPt*ptFactor) + `"/></a:spcAft>`
	}

	base += bulletNodeXML(style)

	return base + `</a:pPr>`
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
