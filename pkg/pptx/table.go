package pptx

import (
	"fmt"
	"math"
	"strings"
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

const tableBorderPtToEMU = 12700.0

// TableCell stores text and optional style for one table cell.
type TableCell struct {
	Text            string
	Bold            bool
	BackgroundColor string
	Align           string
	VAlign          string
	BorderColor     string
	BorderWidthPt   float64
}

// NewTableCell creates a styled table cell with text.
func NewTableCell(text string) TableCell {
	return TableCell{Text: text}
}

// WithBold sets bold text for this cell.
func (c TableCell) WithBold(enabled bool) TableCell {
	c.Bold = enabled
	return c
}

// WithBackgroundColor sets cell background fill using RGB hex.
func (c TableCell) WithBackgroundColor(color string) TableCell {
	c.BackgroundColor = normalizeHexColor(color)
	return c
}

// WithAlign sets horizontal text alignment.
func (c TableCell) WithAlign(align string) TableCell {
	c.Align = normalizeTableAlign(align)
	return c
}

// WithAlignLeft sets horizontal text alignment to left.
func (c TableCell) WithAlignLeft() TableCell {
	return c.WithAlign(TableAlignLeft)
}

// WithAlignCenter sets horizontal text alignment to center.
func (c TableCell) WithAlignCenter() TableCell {
	return c.WithAlign(TableAlignCenter)
}

// WithAlignRight sets horizontal text alignment to right.
func (c TableCell) WithAlignRight() TableCell {
	return c.WithAlign(TableAlignRight)
}

// WithAlignJustify sets horizontal text alignment to justify.
func (c TableCell) WithAlignJustify() TableCell {
	return c.WithAlign(TableAlignJustify)
}

// WithVAlign sets vertical text alignment.
func (c TableCell) WithVAlign(vAlign string) TableCell {
	c.VAlign = normalizeTableVAlign(vAlign)
	return c
}

// WithVAlignTop sets vertical text alignment to top.
func (c TableCell) WithVAlignTop() TableCell {
	return c.WithVAlign(TableVAlignTop)
}

// WithVAlignMiddle sets vertical text alignment to middle.
func (c TableCell) WithVAlignMiddle() TableCell {
	return c.WithVAlign(TableVAlignMiddle)
}

// WithVAlignBottom sets vertical text alignment to bottom.
func (c TableCell) WithVAlignBottom() TableCell {
	return c.WithVAlign(TableVAlignBottom)
}

// WithBorder sets a uniform border (all 4 sides) in points and RGB hex color.
func (c TableCell) WithBorder(widthPt float64, color string) TableCell {
	c.BorderWidthPt = widthPt
	c.BorderColor = normalizeHexColor(color)
	return c
}

func (c TableCell) borderWidthEMU() int64 {
	if c.BorderWidthPt <= 0 {
		return 0
	}
	width := int64(math.Round(c.BorderWidthPt * tableBorderPtToEMU))
	if width < 1 {
		return 1
	}
	return width
}

// Table is a simple slide table with fixed columns and text cells.
type Table struct {
	X            int64
	Y            int64
	CX           int64
	CY           int64
	ColumnWidths []int64
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

func validateTable(table Table, slideIndex int) error {
	if table.X < 0 || table.Y < 0 {
		return fmt.Errorf("slide %d table position cannot be negative", slideIndex)
	}
	if table.CX <= 0 || table.CY <= 0 {
		return fmt.Errorf("slide %d table size must be > 0", slideIndex)
	}
	if len(table.ColumnWidths) == 0 {
		return fmt.Errorf("slide %d table must define at least one column", slideIndex)
	}
	for columnIndex, width := range table.ColumnWidths {
		if width <= 0 {
			return fmt.Errorf("slide %d table column %d width must be > 0", slideIndex, columnIndex+1)
		}
	}

	rows := tableRowsForRender(table)
	if len(rows) == 0 {
		return fmt.Errorf("slide %d table must define at least one row", slideIndex)
	}
	for rowIndex, row := range rows {
		if len(row) != len(table.ColumnWidths) {
			return fmt.Errorf(
				"slide %d table row %d has %d cells; expected %d",
				slideIndex,
				rowIndex+1,
				len(row),
				len(table.ColumnWidths),
			)
		}
		for cellIndex, cell := range row {
			if err := validateTableCell(cell, slideIndex, rowIndex+1, cellIndex+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateTableCell(cell TableCell, slideIndex int, rowIndex int, cellIndex int) error {
	if color := strings.TrimSpace(cell.BackgroundColor); color != "" && !isHexColor(color) {
		return fmt.Errorf("slide %d table row %d cell %d background color must be 6-digit RGB hex", slideIndex, rowIndex, cellIndex)
	}
	if align := strings.TrimSpace(cell.Align); align != "" && !isTableAlign(align) {
		return fmt.Errorf("slide %d table row %d cell %d align must be one of l|ctr|r|just", slideIndex, rowIndex, cellIndex)
	}
	if vAlign := strings.TrimSpace(cell.VAlign); vAlign != "" && !isTableVAlign(vAlign) {
		return fmt.Errorf("slide %d table row %d cell %d valign must be one of t|ctr|b", slideIndex, rowIndex, cellIndex)
	}
	if math.IsNaN(cell.BorderWidthPt) || math.IsInf(cell.BorderWidthPt, 0) {
		return fmt.Errorf("slide %d table row %d cell %d border width must be finite", slideIndex, rowIndex, cellIndex)
	}
	if cell.BorderWidthPt < 0 {
		return fmt.Errorf("slide %d table row %d cell %d border width must be >= 0", slideIndex, rowIndex, cellIndex)
	}
	if cell.BorderWidthPt > 0 && !isHexColor(cell.BorderColor) {
		return fmt.Errorf("slide %d table row %d cell %d border color must be 6-digit RGB hex", slideIndex, rowIndex, cellIndex)
	}
	if strings.TrimSpace(cell.BorderColor) != "" && cell.BorderWidthPt <= 0 {
		return fmt.Errorf("slide %d table row %d cell %d border width must be > 0 when border color is set", slideIndex, rowIndex, cellIndex)
	}
	return nil
}

func tableRowsForRender(table Table) [][]TableCell {
	if len(table.renderRows) > 0 {
		return copyTableRows(table.renderRows)
	}
	if len(table.StyledRows) > 0 && len(table.StyledRows) == len(table.Rows) {
		rows := make([][]TableCell, len(table.Rows))
		for i := range table.Rows {
			styled := copyTableCells(table.StyledRows[i])
			if len(styled) == len(table.Rows[i]) {
				rows[i] = styled
				continue
			}
			rows[i] = plainRowToCells(table.Rows[i])
		}
		return rows
	}
	if len(table.StyledRows) > 0 && len(table.Rows) == 0 {
		return copyTableRows(table.StyledRows)
	}

	rows := make([][]TableCell, len(table.Rows))
	for i := range table.Rows {
		rows[i] = plainRowToCells(table.Rows[i])
	}
	return rows
}

func plainRowToCells(cells []string) []TableCell {
	row := make([]TableCell, len(cells))
	for i, text := range cells {
		row[i] = TableCell{Text: text}
	}
	return row
}

func copyTableRows(rows [][]TableCell) [][]TableCell {
	out := make([][]TableCell, len(rows))
	for i := range rows {
		out[i] = copyTableCells(rows[i])
	}
	return out
}

func copyTableCells(cells []TableCell) []TableCell {
	row := make([]TableCell, len(cells))
	copy(row, cells)
	return row
}

func normalizeTableAlign(align string) string {
	return strings.ToLower(strings.TrimSpace(align))
}

func normalizeTableVAlign(vAlign string) string {
	return strings.ToLower(strings.TrimSpace(vAlign))
}

func isTableAlign(align string) bool {
	switch normalizeTableAlign(align) {
	case TableAlignLeft, TableAlignCenter, TableAlignRight, TableAlignJustify:
		return true
	default:
		return false
	}
}

func isTableVAlign(vAlign string) bool {
	switch normalizeTableVAlign(vAlign) {
	case TableVAlignTop, TableVAlignMiddle, TableVAlignBottom:
		return true
	default:
		return false
	}
}
