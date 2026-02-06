package pptx

import "fmt"

// Table is a simple slide table with fixed columns and text cells.
type Table struct {
	X            int64
	Y            int64
	CX           int64
	CY           int64
	ColumnWidths []int64
	Rows         [][]string
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
	}
}

// AddRow appends one row and returns the updated table.
func (t Table) AddRow(cells []string) Table {
	row := make([]string, len(cells))
	copy(row, cells)
	t.Rows = append(t.Rows, row)
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
	if len(table.Rows) == 0 {
		return fmt.Errorf("slide %d table must define at least one row", slideIndex)
	}
	for rowIndex, row := range table.Rows {
		if len(row) != len(table.ColumnWidths) {
			return fmt.Errorf(
				"slide %d table row %d has %d cells; expected %d",
				slideIndex,
				rowIndex+1,
				len(row),
				len(table.ColumnWidths),
			)
		}
	}
	return nil
}
