package table

import (
	"errors"
	"fmt"
	"strconv"
)

func MergeCellsInFrame(frame []byte, row1, col1, row2, col2 int) ([]byte, error) {
	rowSpan, colSpan, err := validateMergeRange(frame, row1, col1, row2, col2)
	if err != nil {
		return nil, err
	}

	updatedFrame, err := MutateTableRows(frame, row1, row2, func(r int, rowContent []byte) ([]byte, error) {
		return MutateTableCells(rowContent, col1, col2, func(c int, cellContent []byte) ([]byte, error) {
			return mergeCellContent(cellContent, r, c, row1, col1, rowSpan, colSpan), nil
		})
	})
	if err != nil {
		return nil, err
	}
	return updatedFrame, nil
}

func validateMergeRange(frame []byte, row1, col1, row2, col2 int) (int, int, error) {
	if row1 < 0 || col1 < 0 || row2 < 0 || col2 < 0 {
		return 0, 0, errors.New("merge coordinates must be non-negative")
	}
	if row1 > row2 || col1 > col2 {
		return 0, 0, errors.New("merge coordinates must be ordered: row1<=row2 and col1<=col2")
	}

	parsed, err := ParseTable(frame)
	if err != nil {
		return 0, 0, err
	}
	rows, cols := Dimensions(parsed)
	if row2 >= rows || col2 >= cols {
		return 0, 0, fmt.Errorf(
			"merge range [%d:%d,%d:%d] out of table bounds %dx%d",
			row1,
			row2,
			col1,
			col2,
			rows,
			cols,
		)
	}
	return row2 - row1 + 1, col2 - col1 + 1, nil
}

func mergeCellContent(cellContent []byte, row, col, row1, col1, rowSpan, colSpan int) []byte {
	cellContent = clearMergeAttrs(cellContent)
	if row == row1 && col == col1 {
		return setOriginMergeAttrs(cellContent, rowSpan, colSpan)
	}
	if row > row1 {
		cellContent = SetTcAttr(cellContent, "vMerge", "1")
	}
	if col > col1 {
		cellContent = SetTcAttr(cellContent, "hMerge", "1")
	}
	return cellContent
}

func clearMergeAttrs(cellContent []byte) []byte {
	cellContent = RemoveTcAttr(cellContent, "rowSpan")
	cellContent = RemoveTcAttr(cellContent, "gridSpan")
	cellContent = RemoveTcAttr(cellContent, "vMerge")
	cellContent = RemoveTcAttr(cellContent, "hMerge")
	return cellContent
}

func setOriginMergeAttrs(cellContent []byte, rowSpan, colSpan int) []byte {
	if rowSpan > 1 {
		cellContent = SetTcAttr(cellContent, "rowSpan", strconv.Itoa(rowSpan))
	}
	if colSpan > 1 {
		cellContent = SetTcAttr(cellContent, "gridSpan", strconv.Itoa(colSpan))
	}
	return cellContent
}

func SplitCellInFrame(frame []byte, row, col int) ([]byte, error) {
	parsed, err := ParseTable(frame)
	if err != nil {
		return nil, err
	}
	rows, cols := Dimensions(parsed)
	if row < 0 || row >= rows || col < 0 || col >= cols {
		return nil, fmt.Errorf("table cell [%d,%d] out of range", row, col)
	}

	cell := parsed.Rows[row].Cells[col]
	rowSpan := cell.RowSpan
	if rowSpan <= 0 {
		rowSpan = 1
	}
	colSpan := cell.GridSpan
	if colSpan <= 0 {
		colSpan = 1
	}
	if rowSpan == 1 && colSpan == 1 {
		return nil, fmt.Errorf("cell [%d,%d] is not a merge origin", row, col)
	}

	rowEnd := row + rowSpan - 1
	colEnd := col + colSpan - 1
	if rowEnd >= rows || colEnd >= cols {
		return nil, fmt.Errorf("merged span at [%d,%d] exceeds table bounds", row, col)
	}

	updatedFrame, err := MutateTableRows(frame, row, rowEnd, func(r int, rowContent []byte) ([]byte, error) {
		return MutateTableCells(rowContent, col, colEnd, func(c int, cellContent []byte) ([]byte, error) {
			cellContent = RemoveTcAttr(cellContent, "rowSpan")
			cellContent = RemoveTcAttr(cellContent, "gridSpan")
			cellContent = RemoveTcAttr(cellContent, "vMerge")
			cellContent = RemoveTcAttr(cellContent, "hMerge")
			if r == row && c == col {
				return cellContent, nil
			}
			return cellContent, nil
		})
	})
	if err != nil {
		return nil, err
	}
	return updatedFrame, nil
}
