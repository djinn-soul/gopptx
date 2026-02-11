package pptxxml

import (
	"fmt"
	"strconv"
	"strings"
)

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
}

func bulletStyleAt(all []BulletParagraphSpec, index int) BulletParagraphSpec {
	if len(all) == 0 || index < 0 || index >= len(all) {
		return BulletParagraphSpec{}
	}
	return all[index]
}

func bulletParagraphPropsXML(style BulletParagraphSpec) string {
	marL, indent := bulletIndent(style.Level)
	base := `<a:pPr lvl="` + strconv.Itoa(style.Level) + `" marL="` + strconv.Itoa(marL) + `" indent="` + strconv.Itoa(indent) + `"`
	if style.Align != "" {
		base += ` algn="` + Escape(style.Align) + `"`
	}
	base += ">"

	if style.LineSpacingPct > 0 {
		base += `<a:lnSpc><a:spcPct val="` + strconv.Itoa(style.LineSpacingPct*1000) + `"/></a:lnSpc>`
	}
	if style.SpaceBeforePt > 0 {
		base += `<a:spcBef><a:spcPts val="` + strconv.Itoa(style.SpaceBeforePt*100) + `"/></a:spcBef>`
	}
	if style.SpaceAfterPt > 0 {
		base += `<a:spcAft><a:spcPts val="` + strconv.Itoa(style.SpaceAfterPt*100) + `"/></a:spcAft>`
	}

	base += bulletNodeXML(style)

	return base + `</a:pPr>`
}

func bulletIndent(level int) (int, int) {
	indent := 457200 + (level * 457200)
	marginLeft := level*457200 + indent
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
		b.WriteString(strconv.Itoa(style.BulletSize * 1000))
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
	case "none":
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
