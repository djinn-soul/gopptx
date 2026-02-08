package pptx

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
