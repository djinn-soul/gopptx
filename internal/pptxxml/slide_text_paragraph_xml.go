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
	base += `>` + bulletNodeXML(style)

	if style.LineSpacingPct > 0 {
		base += `<a:lnSpc><a:spcPct val="` + strconv.Itoa(style.LineSpacingPct*1000) + `"/></a:lnSpc>`
	}
	if style.SpaceBeforePt > 0 {
		base += `<a:spcBef><a:spcPts val="` + strconv.Itoa(style.SpaceBeforePt*100) + `"/></a:spcBef>`
	}
	if style.SpaceAfterPt > 0 {
		base += `<a:spcAft><a:spcPts val="` + strconv.Itoa(style.SpaceAfterPt*100) + `"/></a:spcAft>`
	}

	return base + `</a:pPr>`
}

func bulletIndent(level int) (int, int) {
	indent := 457200 + (level * 457200)
	marginLeft := level*457200 + indent
	return marginLeft, -indent
}

func bulletNodeXML(style BulletParagraphSpec) string {
	switch normalizeBulletStyle(style.BulletStyle) {
	case "", "bullet":
		return `<a:buChar char="•"/>`
	case "number":
		return `<a:buAutoNum type="arabicPeriod"/>`
	case "letter_lower":
		return `<a:buAutoNum type="alphaLcPeriod"/>`
	case "letter_upper":
		return `<a:buAutoNum type="alphaUcPeriod"/>`
	case "roman_lower":
		return `<a:buAutoNum type="romanLcPeriod"/>`
	case "roman_upper":
		return `<a:buAutoNum type="romanUcPeriod"/>`
	case "custom":
		return `<a:buChar char="` + Escape(strings.TrimSpace(style.BulletChar)) + `"/>`
	case "none":
		return `<a:buNone/>`
	default:
		panic(fmt.Sprintf("unsupported bullet style: %q", style.BulletStyle))
	}
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
