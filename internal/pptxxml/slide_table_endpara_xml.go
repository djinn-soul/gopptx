package pptxxml

import (
	"strconv"
	"strings"
)

// tableCellEndParaRPrXML emits <a:endParaRPr> for the paragraph.
// PowerPoint uses endParaRPr to determine the font size/family of the phantom
// cursor at the end of a paragraph, including empty cells without run text.
func tableCellEndParaRPrXML(cell TableCellSpec) string {
	hasFontStyle := cell.SizePt > 0 || strings.TrimSpace(cell.FontName) != "" ||
		cell.Bold || strings.TrimSpace(cell.Color) != ""
	if !hasFontStyle {
		return `<a:endParaRPr lang="en-US" dirty="0"/>`
	}
	var b strings.Builder
	b.WriteString(`<a:endParaRPr lang="en-US" dirty="0"`)
	if cell.Bold {
		b.WriteString(` b="1"`)
	}
	if cell.SizePt > 0 {
		b.WriteString(` sz="`)
		b.WriteString(strconv.Itoa(int(cell.SizePt * fontSzScale)))
		b.WriteString(`"`)
	}
	b.WriteString(`>`)
	if strings.TrimSpace(cell.Color) != "" {
		b.WriteString(`<a:solidFill><a:srgbClr val="`)
		b.WriteString(Escape(cell.Color))
		b.WriteString(`"/></a:solidFill>`)
	}
	if strings.TrimSpace(cell.FontName) != "" {
		b.WriteString(`<a:latin typeface="`)
		b.WriteString(Escape(cell.FontName))
		b.WriteString(`"/>`)
	}
	b.WriteString(`</a:endParaRPr>`)
	return b.String()
}
