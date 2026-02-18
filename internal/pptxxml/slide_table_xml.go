package pptxxml

import (
	"strconv"
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

	// Accessibility
	AltText      string
	IsDecorative bool
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
	return `
<p:graphicFrame>
<p:nvGraphicFramePr>
<p:cNvPr id="` + strconv.Itoa(shapeID) + `" name="Table 1"` + makeCNvPrAttrs(table.AltText, table.IsDecorative) + `/>
<p:cNvGraphicFramePr><a:graphicFrameLocks noGrp="1"/></p:cNvGraphicFramePr>
<p:nvPr/>
</p:nvGraphicFramePr>
<p:xfrm>
<a:off x="` + strconv.FormatInt(table.X, 10) + `" y="` + strconv.FormatInt(table.Y, 10) + `"/>
<a:ext cx="` + strconv.FormatInt(table.CX, 10) + `" cy="` + strconv.FormatInt(table.CY, 10) + `"/>
</p:xfrm>
` + tableGraphicXML(table) + `
</p:graphicFrame>`
}

func makeCNvPrAttrs(altText string, isDecorative bool) string {
	if isDecorative || altText == "" {
		return ` descr=""`
	}
	escaped := Escape(altText)
	return ` descr="` + escaped + `" title="` + escaped + `"`
}

func tableGraphicXML(table *TableSpec) string {
	var b strings.Builder
	columnWidths := tableColumnWidthsForRender(table)
	b.WriteString(`
<a:graphic>
<a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/table">
<a:tbl>
<a:tblPr firstRow="1" bandRow="1">
<a:tableStyleId>{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}</a:tableStyleId>
</a:tblPr>
<a:tblGrid>`)

	for _, width := range columnWidths {
		b.WriteString(`
<a:gridCol w="`)
		b.WriteString(strconv.FormatInt(width, 10))
		b.WriteString(`"/>`)
	}

	b.WriteString(`
</a:tblGrid>`)

	rows := tableStyledRows(table)
	rowHeights := tableRowHeightsForRender(table, len(rows))
	for rowIndex, row := range rows {
		b.WriteString(`
<a:tr h="`)
		b.WriteString(strconv.FormatInt(rowHeights[rowIndex], 10))
		b.WriteString(`">`)
		for _, cell := range row {
			b.WriteString(`
`)
			b.WriteString(tableCellOpenTag(cell))
			b.WriteString(`
<a:txBody>
`)
			b.WriteString(tableCellBodyPrXML(cell))
			b.WriteString(`
<a:lstStyle/>
<a:p>
`)
			b.WriteString(tableCellParagraphPropsXML(cell.Align))
			b.WriteString(`
<a:r>
<a:rPr lang="en-US" dirty="0"`)
			b.WriteString(tableCellBoldAttr(cell.Bold))
			b.WriteString(`/>
<a:t>`)
			b.WriteString(Escape(cell.Text))
			b.WriteString(`</a:t>
</a:r>
</a:p>
</a:txBody>
`)
			b.WriteString(tableCellPropsXML(cell))
			b.WriteString(`
</a:tc>`)
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

func tableColumnWidthsForRender(table *TableSpec) []int64 {
	if len(table.ColumnWidths) > 0 {
		widths := make([]int64, len(table.ColumnWidths))
		copy(widths, table.ColumnWidths)
		return widths
	}

	columnCount := 0
	for _, row := range table.Rows {
		if len(row) > columnCount {
			columnCount = len(row)
		}
	}
	for _, row := range table.StyledRows {
		if len(row) > columnCount {
			columnCount = len(row)
		}
	}
	if columnCount == 0 {
		columnCount = 1
	}

	defaultWidth := int64(1828800)
	if table.CX > 0 {
		defaultWidth = table.CX / int64(columnCount)
		if defaultWidth <= 0 {
			defaultWidth = 1828800
		}
	}

	widths := make([]int64, 0, columnCount)
	for i := 0; i < columnCount; i++ {
		widths = append(widths, defaultWidth)
	}
	return widths
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
		` w="` + strconv.FormatInt(border.Width, 10) +
		`"><a:solidFill><a:srgbClr val="` + Escape(border.Color) +
		`"/></a:solidFill><a:prstDash val="` + Escape(dash) + `"/></a:` + side + `>`
}

func tableCellBorderDash(dash string) string {
	switch strings.ToLower(strings.TrimSpace(dash)) {
	case "", strokeDashSolid:
		return strokeDashSolid
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
