package tables

func tableRowsForRender(table Table) [][]TableCell {
	if len(table.renderRows) > 0 {
		return copyTableRows(table.renderRows)
	}
	if len(table.StyledRows) > 0 {
		if len(table.Rows) == 0 {
			return copyTableRows(table.StyledRows)
		}
		if len(table.StyledRows) == len(table.Rows) {
			return mergeStyledAndPlainRows(table.Rows, table.StyledRows)
		}
	}

	rows := make([][]TableCell, len(table.Rows))
	for i := range table.Rows {
		rows[i] = plainRowToCells(table.Rows[i])
	}
	return rows
}

func mergeStyledAndPlainRows(plainRows [][]string, styledRows [][]TableCell) [][]TableCell {
	rows := make([][]TableCell, len(plainRows))
	for i := range plainRows {
		styled := copyTableCells(styledRows[i])
		if len(styled) == len(plainRows[i]) {
			rows[i] = styled
		} else {
			rows[i] = plainRowToCells(plainRows[i])
		}
	}
	return rows
}

func plainRowToCells(cells []string) []TableCell {
	row := make([]TableCell, len(cells))
	for i, text := range cells {
		row[i] = NewTableCell(text)
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
	for i := range cells {
		row[i] = cells[i]
		row[i].BorderLeft = cloneTableCellBorder(cells[i].BorderLeft)
		row[i].BorderRight = cloneTableCellBorder(cells[i].BorderRight)
		row[i].BorderTop = cloneTableCellBorder(cells[i].BorderTop)
		row[i].BorderBottom = cloneTableCellBorder(cells[i].BorderBottom)
		row[i].MarginLeftPt = CloneFloat64Pointer(cells[i].MarginLeftPt)
		row[i].MarginRightPt = CloneFloat64Pointer(cells[i].MarginRightPt)
		row[i].MarginTopPt = CloneFloat64Pointer(cells[i].MarginTopPt)
		row[i].MarginBottomPt = CloneFloat64Pointer(cells[i].MarginBottomPt)
		row[i].WrapText = CloneBoolPointer(cells[i].WrapText)
	}
	return row
}
