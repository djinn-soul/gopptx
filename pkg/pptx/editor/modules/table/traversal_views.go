package table

func buildTableTraversalViews(
	rowCount int,
	colCount int,
	cells []map[string]any,
	rowHeights []int64,
	columnWidths []int64,
) ([]map[string]any, []map[string]any) {
	rowsView := initRowsView(rowCount, colCount, rowHeights)
	colsView := initColsView(rowCount, colCount, columnWidths)
	attachTraversalCells(rowsView, colsView, cells, rowCount, colCount)
	return rowsView, colsView
}

func initRowsView(rowCount, colCount int, rowHeights []int64) []map[string]any {
	rowsView := make([]map[string]any, rowCount)
	for r := range rowCount {
		rowsView[r] = map[string]any{
			"index":  r,
			"height": traversalMeasure(rowHeights, r),
			"cells":  make([]map[string]any, 0, colCount),
		}
	}
	return rowsView
}

func initColsView(rowCount, colCount int, columnWidths []int64) []map[string]any {
	colsView := make([]map[string]any, colCount)
	for c := range colCount {
		colsView[c] = map[string]any{
			"index": c,
			"width": traversalMeasure(columnWidths, c),
			"cells": make([]map[string]any, 0, rowCount),
		}
	}
	return colsView
}

func traversalMeasure(values []int64, index int) int64 {
	if index < 0 || index >= len(values) {
		return 0
	}
	return values[index]
}

func attachTraversalCells(
	rowsView []map[string]any,
	colsView []map[string]any,
	cells []map[string]any,
	rowCount, colCount int,
) {
	for _, cell := range cells {
		rowIndex, colIndex, ok := cellCoordinates(cell)
		if !ok {
			continue
		}
		appendTraversalRowCell(rowsView, rowIndex, rowCount, cell)
		appendTraversalColCell(colsView, colIndex, colCount, cell)
	}
}

func cellCoordinates(cell map[string]any) (int, int, bool) {
	rowIndex, rowOK := cell["row"].(int)
	colIndex, colOK := cell["col"].(int)
	if !rowOK || !colOK {
		return 0, 0, false
	}
	return rowIndex, colIndex, true
}

func appendTraversalRowCell(
	rowsView []map[string]any,
	rowIndex, rowCount int,
	cell map[string]any,
) {
	if rowIndex < 0 || rowIndex >= rowCount {
		return
	}
	rowCells, ok := rowsView[rowIndex]["cells"].([]map[string]any)
	if !ok {
		return
	}
	rowsView[rowIndex]["cells"] = append(rowCells, cell)
}

func appendTraversalColCell(
	colsView []map[string]any,
	colIndex, colCount int,
	cell map[string]any,
) {
	if colIndex < 0 || colIndex >= colCount {
		return
	}
	colCells, ok := colsView[colIndex]["cells"].([]map[string]any)
	if !ok {
		return
	}
	colsView[colIndex]["cells"] = append(colCells, cell)
}
