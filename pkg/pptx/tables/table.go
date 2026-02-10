package tables

import (
	"github.com/djinn-soul/gopptx/internal/pptxxml"
)

const (
	// TableAlignLeft sets horizontal text alignment to left.
	TableAlignLeft = "l"
	// TableAlignCenter sets horizontal text alignment to center.
	TableAlignCenter = "ctr"
	// TableAlignRight sets horizontal text alignment to right.
	TableAlignRight = "r"
	// TableAlignJustify sets horizontal text alignment to justify.
	TableAlignJustify = "just"

	// TableVAlignTop sets vertical text alignment to top.
	TableVAlignTop = "t"
	// TableVAlignMiddle sets vertical text alignment to middle.
	TableVAlignMiddle = "ctr"
	// TableVAlignBottom sets vertical text alignment to bottom.
	TableVAlignBottom = "b"
)

const (
	// TableBorderDashSolid emits a solid line.
	TableBorderDashSolid = "solid"
	// TableBorderDashDash emits a dashed line.
	TableBorderDashDash = "dash"
	// TableBorderDashDot emits a dotted line.
	TableBorderDashDot = "dot"
	// TableBorderDashDashDot emits dash-dot line.
	TableBorderDashDashDot = "dashDot"
	// TableBorderDashLongDash emits long-dash line.
	TableBorderDashLongDash = "lgDash"
)

const (
	borderSideLeft   = "left"
	borderSideRight  = "right"
	borderSideTop    = "top"
	borderSideBottom = "bottom"
)

// Table is a simple slide table with fixed columns and text cells.
type Table struct {
	X            int64
	Y            int64
	CX           int64
	CY           int64
	ColumnWidths []int64
	RowHeights   []int64
	Rows         [][]string
	StyledRows   [][]TableCell
	renderRows   [][]TableCell
}

// NewTable creates a table with default placement and size.
func NewTable(columnWidths []int64) Table {
	widths := make([]int64, len(columnWidths))
	copy(widths, columnWidths)
	return Table{
		X:            457200,
		Y:            1600200,
		CX:           8230200,
		CY:           3200400,
		ColumnWidths: widths,
		RowHeights:   nil,
		Rows:         make([][]string, 0),
		StyledRows:   make([][]TableCell, 0),
		renderRows:   make([][]TableCell, 0),
	}
}

// AddRow appends one plain-text row and returns the updated table.
func (t Table) AddRow(cells []string) Table {
	row := make([]string, len(cells))
	copy(row, cells)
	t.Rows = append(t.Rows, row)
	t.renderRows = append(t.renderRows, plainRowToCells(row))
	return t
}

// AddStyledRow appends one styled row and returns the updated table.
func (t Table) AddStyledRow(cells []TableCell) Table {
	row := copyTableCells(cells)
	t.StyledRows = append(t.StyledRows, row)
	t.renderRows = append(t.renderRows, row)

	textRow := make([]string, len(row))
	for i, cell := range row {
		textRow[i] = cell.Text
	}
	t.Rows = append(t.Rows, textRow)
	return t
}

// Position sets table position in EMU.
func (t Table) Position(x int64, y int64) Table {
	t.X = x
	t.Y = y
	return t
}

// Size sets table size in EMU.
func (t Table) Size(cx int64, cy int64) Table {
	t.CX = cx
	t.CY = cy
	return t
}

// WithRowHeights sets explicit row heights in EMU. Length must match row count.
func (t Table) WithRowHeights(heights []int64) Table {
	if len(heights) == 0 {
		t.RowHeights = nil
		return t
	}
	out := make([]int64, len(heights))
	copy(out, heights)
	t.RowHeights = out
	return t
}

// Validate checks the table content for common constraints.
func (t Table) Validate(slideIndex int) error {
	return validateTable(t, slideIndex)
}

// ToTableSpec converts Table to internal XML spec.
func (t Table) ToTableSpec(slideNumber int) (*pptxxml.TableSpec, error) {
	styledRows, err := TableRowsWithMerges(t, slideNumber)
	if err != nil {
		return nil, err
	}
	rows := make([][]string, 0, len(styledRows))
	styledSpecRows := make([][]pptxxml.TableCellSpec, 0, len(styledRows))
	for _, srcRow := range styledRows {
		row := make([]string, len(srcRow))
		specRow := make([]pptxxml.TableCellSpec, len(srcRow))
		for i, cell := range srcRow {
			borders := cell.bordersForRender()
			row[i] = cell.Text
			specRow[i] = pptxxml.TableCellSpec{
				Text:            cell.Text,
				Bold:            cell.Bold,
				BackgroundColor: cell.BackgroundColor,
				Align:           cell.Align,
				VAlign:          cell.VAlign,
				MarginLeft:      TableMarginEMU(cell.MarginLeftPt),
				MarginRight:     TableMarginEMU(cell.MarginRightPt),
				MarginTop:       TableMarginEMU(cell.MarginTopPt),
				MarginBottom:    TableMarginEMU(cell.MarginBottomPt),
				WrapText:        CloneBoolPointer(cell.WrapText),
				RowSpan:         cell.RowSpan,
				ColSpan:         cell.ColSpan,
				VMerge:          cell.VMerge,
				HMerge:          cell.HMerge,
				BorderColor:     cell.BorderColor,
				BorderWidth:     TableBorderWidthEMU(cell.BorderWidthPt),
				BorderLeft:      toXMLTableBorderSpec(borders.Left),
				BorderRight:     toXMLTableBorderSpec(borders.Right),
				BorderTop:       toXMLTableBorderSpec(borders.Top),
				BorderBottom:    toXMLTableBorderSpec(borders.Bottom),
			}
		}
		rows = append(rows, row)
		styledSpecRows = append(styledSpecRows, specRow)
	}
	columnWidths := make([]int64, len(t.ColumnWidths))
	copy(columnWidths, t.ColumnWidths)
	rowHeights := make([]int64, len(t.RowHeights))
	copy(rowHeights, t.RowHeights)

	return &pptxxml.TableSpec{
		X:            t.X,
		Y:            t.Y,
		CX:           t.CX,
		CY:           t.CY,
		ColumnWidths: columnWidths,
		RowHeights:   rowHeights,
		Rows:         rows,
		StyledRows:   styledSpecRows,
	}, nil
}

func toXMLTableBorderSpec(border *TableCellBorder) *pptxxml.TableCellBorderSpec {
	if border == nil || border.WidthPt <= 0 {
		return nil
	}
	return &pptxxml.TableCellBorderSpec{
		Width: TableBorderWidthEMU(border.WidthPt),
		Color: border.Color,
		Dash:  border.Dash,
	}
}
