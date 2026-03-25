package pptxxml

import (
	"fmt"
	"strings"
)

// txStylesXML renders the <p:txStyles> block for a slide master.
func txStylesXML(spec *SlideMasterSpec) string {
	if spec == nil || spec.TxStyles == nil {
		return ""
	}
	tx := spec.TxStyles
	var b strings.Builder
	b.WriteString(`
<p:txStyles>`)
	if len(tx.TitleStyle) > 0 {
		b.WriteString(`
<p:titleStyle>`)
		b.WriteString(textLevelStylesXML(tx.TitleStyle))
		b.WriteString(`
</p:titleStyle>`)
	}
	if len(tx.BodyStyle) > 0 {
		b.WriteString(`
<p:bodyStyle>`)
		b.WriteString(textLevelStylesXML(tx.BodyStyle))
		b.WriteString(`
</p:bodyStyle>`)
	}
	if len(tx.OtherStyle) > 0 {
		b.WriteString(`
<p:otherStyle>`)
		b.WriteString(textLevelStylesXML(tx.OtherStyle))
		b.WriteString(`
</p:otherStyle>`)
	}
	b.WriteString(`
</p:txStyles>`)
	return b.String()
}

// textLevelStylesXML renders <a:lvlNpPr> elements for each text level.
//
//nolint:mnd // Contains specific level limits and point conversion factors.
func textLevelStylesXML(levels []TextLevelStyle) string {
	var b strings.Builder
	for _, lvl := range levels {
		lvlNum := max(lvl.Level+1, 1) // 0-based -> 1-based
		lvlNum = min(lvlNum, 9)

		attrs := ""
		if lvl.IndentEMU > 0 {
			attrs += fmt.Sprintf(` indent="%d"`, lvl.IndentEMU)
		}

		b.WriteString(fmt.Sprintf(`
<a:lvl%dpPr%s>`, lvlNum, attrs))
		if lvl.BulletChar != "" {
			b.WriteString(fmt.Sprintf(`<a:buChar char="%s"/>`, Escape(lvl.BulletChar)))
		}

		rprAttrs := ""
		if lvl.SizePt > 0 {
			rprAttrs += fmt.Sprintf(` sz="%d"`, lvl.SizePt*100)
		}
		if lvl.Bold {
			rprAttrs += ` b="1"`
		}
		if lvl.Italic {
			rprAttrs += ` i="1"`
		}
		b.WriteString(fmt.Sprintf(`<a:defRPr%s>`, rprAttrs))

		if lvl.Color != "" {
			color := strings.TrimPrefix(lvl.Color, "#")
			b.WriteString(fmt.Sprintf(`<a:solidFill><a:srgbClr val="%s"/></a:solidFill>`, color))
		}
		if lvl.Font != "" {
			b.WriteString(fmt.Sprintf(`<a:latin typeface="%s"/>`, Escape(lvl.Font)))
		}
		b.WriteString(`</a:defRPr>`)
		b.WriteString(fmt.Sprintf(`
</a:lvl%dpPr>`, lvlNum))
	}
	return b.String()
}

const imageRidStart = 8

// SlideMasterRelationships renders ppt/slideMasters/_rels/slideMasterN.xml.rels.
// imageTargets are optional media paths for master images (e.g. "../media/image1.png").
func SlideMasterRelationships(imageTargets []string, masterIndex int, themeIndex int) string {
	if masterIndex < 1 {
		masterIndex = 1
	}
	if themeIndex < 1 {
		themeIndex = 1
	}
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`)
	for i := 1; i <= 6; i++ {
		b.WriteString(fmt.Sprintf(`
<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/%s"/>`,
			i, slideLayoutPartName(i, masterIndex)))
	}
	b.WriteString(fmt.Sprintf(`
<Relationship Id="rId7" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="../theme/theme%d.xml"/>`,
		themeIndex))
	for i, target := range imageTargets {
		b.WriteString(fmt.Sprintf(`
<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="%s"/>`,
			imageRidStart+i, Escape(target)))
	}
	b.WriteString(`
</Relationships>`)
	return b.String()
}
