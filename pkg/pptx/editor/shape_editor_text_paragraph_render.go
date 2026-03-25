package editor

import (
	"errors"
	"fmt"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

const (
	alignCenter       = "center"
	alignCtr          = "ctr"
	alignJust         = "just"
	alignDist         = "dist"
	bulletDot         = "•"
	bulletStyleNone   = "none"
	bulletStyleCustom = "custom"
	bulletSizeScale   = 1000
)

func renderParagraphPropsXML(paragraph *common.Paragraph) (string, error) {
	if paragraph == nil {
		return "", nil
	}
	attrs, err := buildParagraphAttributes(paragraph)
	if err != nil {
		return "", err
	}
	children, err := renderParagraphChildren(paragraph)
	if err != nil {
		return "", err
	}
	if attrs == "" && children == "" {
		return "", nil
	}
	if children == "" {
		return `<a:pPr` + attrs + `/>`, nil
	}
	return `<a:pPr` + attrs + `>` + children + `</a:pPr>`, nil
}

func buildParagraphAttributes(paragraph *common.Paragraph) (string, error) {
	var attrs strings.Builder
	if paragraph.Indent != nil {
		attrs.WriteString(fmt.Sprintf(` marL="%d"`, *paragraph.Indent))
	}
	if paragraph.Hanging != nil {
		if *paragraph.Hanging < 0 {
			return "", errors.New("paragraph.hanging must be >= 0")
		}
		attrs.WriteString(fmt.Sprintf(` indent="%d"`, -*paragraph.Hanging))
	}
	if paragraph.Alignment != nil {
		normalized, err := normalizeParagraphAlignment(*paragraph.Alignment)
		if err != nil {
			return "", err
		}
		attrs.WriteString(fmt.Sprintf(` algn="%s"`, normalized))
	}
	if paragraph.Level != nil {
		if *paragraph.Level < 0 || *paragraph.Level > 8 {
			return "", errors.New("paragraph.level must be between 0 and 8")
		}
		attrs.WriteString(fmt.Sprintf(` lvl="%d"`, *paragraph.Level))
	}
	return attrs.String(), nil
}

func renderParagraphChildren(paragraph *common.Paragraph) (string, error) {
	var out strings.Builder
	if err := writeParagraphBulletNodes(&out, paragraph); err != nil {
		return "", err
	}
	if len(paragraph.TabStops) > 0 {
		if err := writeTabStops(&out, paragraph.TabStops); err != nil {
			return "", err
		}
	}
	if err := writeParagraphSpacingNodes(&out, paragraph); err != nil {
		return "", err
	}
	return out.String(), nil
}

func writeTabStops(out *strings.Builder, tabStops []int) error {
	out.WriteString(`<a:tabLst>`)
	for _, pos := range tabStops {
		if pos < 0 {
			return errors.New("paragraph.tab_stops values must be >= 0")
		}
		_, _ = fmt.Fprintf(out, `<a:tab pos="%d"/>`, pos)
	}
	out.WriteString(`</a:tabLst>`)
	return nil
}

func writeParagraphSpacingNodes(out *strings.Builder, paragraph *common.Paragraph) error {
	lnSp, err := renderParagraphSpacingNode("lnSp", paragraph.LineSpacingPct, paragraph.LineSpacingPts)
	if err != nil {
		return err
	}
	out.WriteString(lnSp)

	spcBef, err := renderParagraphPtsNode("spcBef", paragraph.SpaceBeforePts)
	if err != nil {
		return err
	}
	out.WriteString(spcBef)

	spcAft, err := renderParagraphPtsNode("spcAft", paragraph.SpaceAfterPts)
	if err != nil {
		return err
	}
	out.WriteString(spcAft)
	return nil
}

func writeParagraphBulletNodes(out *strings.Builder, paragraph *common.Paragraph) error {
	if paragraph.BulletStyle == nil {
		return nil
	}
	bulletXML, err := renderParagraphBulletStyleXML(*paragraph.BulletStyle, paragraph.BulletChar)
	if err != nil {
		return err
	}
	out.WriteString(bulletXML)
	if paragraph.BulletColor != nil {
		color, err := editorshape.NormalizeHexColor(*paragraph.BulletColor)
		if err != nil {
			return fmt.Errorf("paragraph.bullet_color: %w", err)
		}
		_, _ = fmt.Fprintf(out, `<a:buClr><a:srgbClr val="%s"/></a:buClr>`, color)
	}
	if paragraph.BulletSizePct != nil {
		if *paragraph.BulletSizePct < 0 {
			return errors.New("paragraph.bullet_size_pct must be >= 0")
		}
		_, _ = fmt.Fprintf(out, `<a:buSzPct val="%d"/>`, *paragraph.BulletSizePct*bulletSizeScale)
	}
	return nil
}

func renderParagraphBulletStyleXML(rawStyle string, bulletChar *string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(rawStyle)) {
	case bulletStyleNone:
		return `<a:buNone/>`, nil
	case "bullet":
		return fmt.Sprintf(`<a:buChar char="%s"/>`, common.XMLEscape(bulletDot)), nil
	case bulletStyleCustom:
		if bulletChar == nil || strings.TrimSpace(*bulletChar) == "" {
			return "", errors.New("paragraph.bullet_char is required when paragraph.bullet_style=custom")
		}
		return fmt.Sprintf(`<a:buChar char="%s"/>`, common.XMLEscape(*bulletChar)), nil
	case "number":
		return `<a:buAutoNum type="arabicPeriod"/>`, nil
	case "letter_lower":
		return `<a:buAutoNum type="alphaLcPeriod"/>`, nil
	case "letter_upper":
		return `<a:buAutoNum type="alphaUcPeriod"/>`, nil
	case "roman_lower":
		return `<a:buAutoNum type="romanLcPeriod"/>`, nil
	case "roman_upper":
		return `<a:buAutoNum type="romanUcPeriod"/>`, nil
	default:
		return "", fmt.Errorf("unsupported paragraph.bullet_style %q", rawStyle)
	}
}

func renderParagraphSpacingNode(tag string, pct *int, pts *int) (string, error) {
	if pct == nil && pts == nil {
		return "", nil
	}
	if pct != nil && pts != nil {
		return "", errors.New("paragraph line spacing cannot set both percent and points")
	}
	if pct != nil {
		if *pct < 0 {
			return "", errors.New("paragraph.line_spacing_pct must be >= 0")
		}
		return fmt.Sprintf(`<a:%s><a:spcPct val="%d"/></a:%s>`, tag, *pct, tag), nil
	}
	if *pts < 0 {
		return "", errors.New("paragraph.line_spacing_pts must be >= 0")
	}
	return fmt.Sprintf(`<a:%s><a:spcPts val="%d"/></a:%s>`, tag, *pts, tag), nil
}

func renderParagraphPtsNode(tag string, pts *int) (string, error) {
	if pts == nil {
		return "", nil
	}
	if *pts < 0 {
		return "", fmt.Errorf("paragraph.%s must be >= 0", tag)
	}
	return fmt.Sprintf(`<a:%s><a:spcPts val="%d"/></a:%s>`, tag, *pts, tag), nil
}

func normalizeParagraphAlignment(raw string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "l", "left":
		return "l", nil
	case alignCtr, alignCenter, "middle":
		return alignCtr, nil
	case "r", "right":
		return "r", nil
	case alignJust, "justify":
		return alignJust, nil
	case alignDist, "distribute":
		return alignDist, nil
	case "thaidist":
		return "thaiDist", nil
	case "justlow":
		return "justLow", nil
	default:
		return "", fmt.Errorf("unsupported paragraph.alignment %q", raw)
	}
}
