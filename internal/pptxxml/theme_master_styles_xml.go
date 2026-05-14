package pptxxml

import (
	"strconv"
	"strings"
)

const slideMasterRelsGrowCap = 1024

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
	b.Grow(256 * len(levels))
	for _, lvl := range levels {
		lvlNum := max(lvl.Level+1, 1) // 0-based -> 1-based
		lvlNum = min(lvlNum, 9)
		lvlNumStr := strconv.Itoa(lvlNum)

		b.WriteString("\n<a:lvl")
		b.WriteString(lvlNumStr)
		b.WriteString("pPr")
		if lvl.IndentEMU > 0 {
			b.WriteString(` indent="`)
			b.WriteString(strconv.FormatInt(lvl.IndentEMU, 10))
			b.WriteString(`"`)
		}
		b.WriteString(`>`)
		if lvl.BulletChar != "" {
			b.WriteString(`<a:buChar char="`)
			b.WriteString(Escape(lvl.BulletChar))
			b.WriteString(`"/>`)
		}

		b.WriteString(`<a:defRPr`)
		if lvl.SizePt > 0 {
			b.WriteString(` sz="`)
			b.WriteString(strconv.Itoa(lvl.SizePt * 100))
			b.WriteString(`"`)
		}
		if lvl.Bold {
			b.WriteString(` b="1"`)
		}
		if lvl.Italic {
			b.WriteString(` i="1"`)
		}
		b.WriteString(`>`)

		if lvl.Color != "" {
			color := strings.TrimPrefix(lvl.Color, "#")
			b.WriteString(`<a:solidFill><a:srgbClr val="`)
			b.WriteString(color)
			b.WriteString(`"/></a:solidFill>`)
		}
		if lvl.Font != "" {
			b.WriteString(`<a:latin typeface="`)
			b.WriteString(Escape(lvl.Font))
			b.WriteString(`"/>`)
		}
		b.WriteString(`</a:defRPr>`)
		b.WriteString("\n</a:lvl")
		b.WriteString(lvlNumStr)
		b.WriteString("pPr>")
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
	b.Grow(slideMasterRelsGrowCap)
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`)
	for i := 1; i <= 6; i++ {
		b.WriteString("\n<Relationship Id=\"rId")
		b.WriteString(strconv.Itoa(i))
		b.WriteString(
			`" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout"` +
				` Target="../slideLayouts/`,
		)
		b.WriteString(slideLayoutPartName(i, masterIndex))
		b.WriteString(`"/>`)
	}
	b.WriteString(
		"\n<Relationship Id=\"rId7\"" +
			" Type=\"http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme\"" +
			" Target=\"../theme/theme",
	)
	b.WriteString(strconv.Itoa(themeIndex))
	b.WriteString(`.xml"/>`)
	for i, target := range imageTargets {
		b.WriteString("\n<Relationship Id=\"rId")
		b.WriteString(strconv.Itoa(imageRidStart + i))
		b.WriteString(`" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="`)
		b.WriteString(Escape(target))
		b.WriteString(`"/>`)
	}
	b.WriteString("\n</Relationships>")
	return b.String()
}
