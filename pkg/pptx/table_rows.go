package pptx

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
	copy(row, cells)
	return row
}
