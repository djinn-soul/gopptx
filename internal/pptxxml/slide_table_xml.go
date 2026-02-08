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
	RowHeights   []int64
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
	MarginLeft      *int64
	MarginRight     *int64
	MarginTop       *int64
	MarginBottom    *int64
	WrapText        *bool
	RowSpan         int
	ColSpan         int
	VMerge          bool
	HMerge          bool
	BorderColor     string
	BorderWidth     int64
	BorderLeft      *TableCellBorderSpec
	BorderRight     *TableCellBorderSpec
	BorderTop       *TableCellBorderSpec
	BorderBottom    *TableCellBorderSpec
}

// TableCellBorderSpec describes one table border side style.
type TableCellBorderSpec struct {
	Width int64
	Color string
	Dash  string
}

type tableCellBorderSet struct {
	Left   *TableCellBorderSpec
	Right  *TableCellBorderSpec
	Top    *TableCellBorderSpec
	Bottom *TableCellBorderSpec
}

func tableShape(table *TableSpec, shapeID int) string {
	return fmt.Sprintf(`
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
%s
</p:graphicFrame>`,
		shapeID,
		table.X,
		table.Y,
		table.CX,
		table.CY,
		tableGraphicXML(table),
	)
}

func tableGraphicXML(table *TableSpec) string {
	var b strings.Builder
	b.WriteString(`
<a:graphic>
<a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/table">
<a:tbl>
<a:tblPr firstRow="1" bandRow="1">
<a:tableStyleId>{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}</a:tableStyleId>
</a:tblPr>
<a:tblGrid>`)

	for _, width := range table.ColumnWidths {
		b.WriteString(fmt.Sprintf(`
<a:gridCol w="%d"/>`, width))
	}

	b.WriteString(`
</a:tblGrid>`)

	rows := tableStyledRows(table)
	rowHeights := tableRowHeightsForRender(table, len(rows))
	for rowIndex, row := range rows {
		b.WriteString(fmt.Sprintf(`
<a:tr h="%d">`, rowHeights[rowIndex]))
		for _, cell := range row {
			b.WriteString(fmt.Sprintf(`
%s
<a:txBody>
%s
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
</a:tc>`, tableCellOpenTag(cell), tableCellBodyPrXML(cell), tableCellParagraphPropsXML(cell.Align), tableCellBoldAttr(cell.Bold), Escape(cell.Text), tableCellPropsXML(cell)))
		}
		b.WriteString(`
</a:tr>`)
	}

	b.WriteString(`
</a:tbl>
</a:graphicData>
</a:graphic>`)
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

func tableCellOpenTag(cell TableCellSpec) string {
	var b strings.Builder
	b.WriteString("<a:tc")
	if cell.RowSpan > 1 {
		b.WriteString(` rowSpan="`)
		b.WriteString(fmt.Sprintf("%d", cell.RowSpan))
		b.WriteString(`"`)
	}
	if cell.ColSpan > 1 {
		b.WriteString(` gridSpan="`)
		b.WriteString(fmt.Sprintf("%d", cell.ColSpan))
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

func tableCellPropsXML(cell TableCellSpec) string {
	borders := tableCellBordersForRender(cell)
	hasFill := strings.TrimSpace(cell.BackgroundColor) != ""
	hasVAlign := strings.TrimSpace(cell.VAlign) != ""
	hasMargins := hasTableCellMargins(cell)
	hasBorder := borders.Left != nil || borders.Right != nil || borders.Top != nil || borders.Bottom != nil
	if !hasFill && !hasVAlign && !hasMargins && !hasBorder {
		return "<a:tcPr/>"
	}

	var b strings.Builder
	b.WriteString("<a:tcPr")
	if hasVAlign {
		b.WriteString(` anchor="`)
		b.WriteString(Escape(cell.VAlign))
		b.WriteString(`"`)
	}
	appendTableCellMarginAttrs(&b, cell)
	b.WriteString(">")

	if hasFill {
		b.WriteString(`<a:solidFill><a:srgbClr val="`)
		b.WriteString(Escape(cell.BackgroundColor))
		b.WriteString(`"/></a:solidFill>`)
	}
	if hasBorder {
		if borders.Left != nil {
			b.WriteString(tableCellBorderXML("lnL", *borders.Left))
		}
		if borders.Right != nil {
			b.WriteString(tableCellBorderXML("lnR", *borders.Right))
		}
		if borders.Top != nil {
			b.WriteString(tableCellBorderXML("lnT", *borders.Top))
		}
		if borders.Bottom != nil {
			b.WriteString(tableCellBorderXML("lnB", *borders.Bottom))
		}
	}

	b.WriteString("</a:tcPr>")
	return b.String()
}

func tableCellBordersForRender(cell TableCellSpec) tableCellBorderSet {
	borders := tableCellBorderSet{
		Left:   cloneTableCellBorderSpec(cell.BorderLeft),
		Right:  cloneTableCellBorderSpec(cell.BorderRight),
		Top:    cloneTableCellBorderSpec(cell.BorderTop),
		Bottom: cloneTableCellBorderSpec(cell.BorderBottom),
	}
	if borders.Left == nil && borders.Right == nil && borders.Top == nil && borders.Bottom == nil {
		if cell.BorderWidth > 0 && strings.TrimSpace(cell.BorderColor) != "" {
			legacy := &TableCellBorderSpec{Width: cell.BorderWidth, Color: cell.BorderColor, Dash: "solid"}
			borders.Left = cloneTableCellBorderSpec(legacy)
			borders.Right = cloneTableCellBorderSpec(legacy)
			borders.Top = cloneTableCellBorderSpec(legacy)
			borders.Bottom = cloneTableCellBorderSpec(legacy)
		}
	}
	return borders
}

func cloneTableCellBorderSpec(border *TableCellBorderSpec) *TableCellBorderSpec {
	if border == nil {
		return nil
	}
	clone := *border
	return &clone
}

func tableCellBorderXML(side string, border TableCellBorderSpec) string {
	dash := tableCellBorderDash(border.Dash)
	return `<a:` + side +
		` w="` + fmt.Sprintf("%d", border.Width) +
		`"><a:solidFill><a:srgbClr val="` + Escape(border.Color) +
		`"/></a:solidFill><a:prstDash val="` + Escape(dash) + `"/></a:` + side + `>`
}

func tableCellBorderDash(dash string) string {
	switch strings.ToLower(strings.TrimSpace(dash)) {
	case "", "solid":
		return "solid"
	case "dash":
		return "dash"
	case "dot":
		return "dot"
	case "dashdot", "dash-dot", "dash_dot":
		return "dashDot"
	case "lgdash", "lg-dash", "longdash", "long-dash", "long_dash":
		return "lgDash"
	default:
		return strings.TrimSpace(dash)
	}
}
