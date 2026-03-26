package pptxxml

import (
	"strconv"
	"strings"
)

func makeCNvPrAttrs(altText string, isDecorative bool) string {
	if isDecorative || altText == "" {
		return ` descr=""`
	}
	escaped := Escape(altText)
	return ` descr="` + escaped + `" title="` + escaped + `"`
}

func tableCellParagraphPropsXML(align string) string {
	if strings.TrimSpace(align) == "" {
		return ""
	}
	return `<a:pPr algn="` + Escape(align) + `"/>`
}

func tableCellOpenTag(cell TableCellSpec) string {
	var b strings.Builder
	b.WriteString("<a:tc")
	if cell.RowSpan > 1 {
		b.WriteString(` rowSpan="`)
		b.WriteString(strconv.Itoa(cell.RowSpan))
		b.WriteString(`"`)
	}
	if cell.ColSpan > 1 {
		b.WriteString(` gridSpan="`)
		b.WriteString(strconv.Itoa(cell.ColSpan))
		b.WriteString(`"`)
	}
	if cell.VMerge {
		b.WriteString(` vMerge="1"`)
	}
	if cell.HMerge {
		b.WriteString(` hMerge="1"`)
	}
	b.WriteString(">")
	return b.String()
}
