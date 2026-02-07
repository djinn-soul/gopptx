package pptxxml

import (
	"fmt"
	"strings"
)

// TableSpec describes one table in a slide.
type TableSpec struct {
	X            int64
	Y            int64
	CX           int64
	CY           int64
	ColumnWidths []int64
	Rows         [][]string
	StyledRows   [][]TableCellSpec
}

// TableCellSpec describes one table cell with optional style.
type TableCellSpec struct {
	Text            string
	Bold            bool
	BackgroundColor string
	Align           string
	VAlign          string
	BorderColor     string
	BorderWidth     int64
}

func tableShape(table *TableSpec, shapeID int) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`
<p:graphicFrame>
<p:nvGraphicFramePr>
<p:cNvPr id="%d" name="Table 1"/>
<p:cNvGraphicFramePr><a:graphicFrameLocks noGrp="1"/></p:cNvGraphicFramePr>
<p:nvPr/>
</p:nvGraphicFramePr>
<p:xfrm>
<a:off x="%d" y="%d"/>
<a:ext cx="%d" cy="%d"/>
</p:xfrm>
<a:graphic>
<a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/table">
<a:tbl>
<a:tblPr firstRow="1" bandRow="1">
<a:tableStyleId>{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}</a:tableStyleId>
</a:tblPr>
<a:tblGrid>`,
		shapeID,
		table.X,
		table.Y,
		table.CX,
		table.CY,
	))

	for _, width := range table.ColumnWidths {
		b.WriteString(fmt.Sprintf(`
<a:gridCol w="%d"/>`, width))
	}

	b.WriteString(`
</a:tblGrid>`)

	rows := tableStyledRows(table)
	rowHeight := table.CY / int64(len(rows))
	if rowHeight <= 0 {
		rowHeight = 1
	}
	for _, row := range rows {
		b.WriteString(fmt.Sprintf(`
<a:tr h="%d">`, rowHeight))
		for _, cell := range row {
			b.WriteString(fmt.Sprintf(`
<a:tc>
<a:txBody>
<a:bodyPr/>
<a:lstStyle/>
<a:p>
%s
<a:r>
<a:rPr lang="en-US" dirty="0"%s/>
<a:t>%s</a:t>
</a:r>
</a:p>
</a:txBody>
%s
</a:tc>`, tableCellParagraphPropsXML(cell.Align), tableCellBoldAttr(cell.Bold), Escape(cell.Text), tableCellPropsXML(cell)))
		}
		b.WriteString(`
</a:tr>`)
	}

	b.WriteString(`
</a:tbl>
</a:graphicData>
</a:graphic>
</p:graphicFrame>`)
	return b.String()
}

func tableStyledRows(table *TableSpec) [][]TableCellSpec {
	if len(table.StyledRows) == len(table.Rows) && len(table.StyledRows) > 0 {
		rows := make([][]TableCellSpec, len(table.StyledRows))
		for i := range table.StyledRows {
			row := make([]TableCellSpec, len(table.StyledRows[i]))
			copy(row, table.StyledRows[i])
			rows[i] = row
		}
		return rows
	}

	rows := make([][]TableCellSpec, len(table.Rows))
	for i := range table.Rows {
		cells := make([]TableCellSpec, len(table.Rows[i]))
		for j, text := range table.Rows[i] {
			cells[j] = TableCellSpec{Text: text}
		}
		rows[i] = cells
	}
	return rows
}

func tableCellBoldAttr(bold bool) string {
	if bold {
		return ` b="1"`
	}
	return ""
}

func tableCellParagraphPropsXML(align string) string {
	if strings.TrimSpace(align) == "" {
		return ""
	}
	return `<a:pPr algn="` + Escape(align) + `"/>`
}

func tableCellPropsXML(cell TableCellSpec) string {
	hasFill := strings.TrimSpace(cell.BackgroundColor) != ""
	hasVAlign := strings.TrimSpace(cell.VAlign) != ""
	hasBorder := cell.BorderWidth > 0 && strings.TrimSpace(cell.BorderColor) != ""
	if !hasFill && !hasVAlign && !hasBorder {
		return "<a:tcPr/>"
	}

	var b strings.Builder
	b.WriteString("<a:tcPr")
	if hasVAlign {
		b.WriteString(` anchor="`)
		b.WriteString(Escape(cell.VAlign))
		b.WriteString(`"`)
	}
	b.WriteString(">")

	if hasFill {
		b.WriteString(`<a:solidFill><a:srgbClr val="`)
		b.WriteString(Escape(cell.BackgroundColor))
		b.WriteString(`"/></a:solidFill>`)
	}
	if hasBorder {
		b.WriteString(tableCellBorderXML("lnL", cell.BorderWidth, cell.BorderColor))
		b.WriteString(tableCellBorderXML("lnR", cell.BorderWidth, cell.BorderColor))
		b.WriteString(tableCellBorderXML("lnT", cell.BorderWidth, cell.BorderColor))
		b.WriteString(tableCellBorderXML("lnB", cell.BorderWidth, cell.BorderColor))
	}

	b.WriteString("</a:tcPr>")
	return b.String()
}

func tableCellBorderXML(side string, width int64, color string) string {
	return `<a:` + side + ` w="` + fmt.Sprintf("%d", width) + `"><a:solidFill><a:srgbClr val="` + Escape(color) + `"/></a:solidFill><a:prstDash val="solid"/></a:` + side + `>`
}
