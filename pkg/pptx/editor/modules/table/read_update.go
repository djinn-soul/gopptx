package table

import (
	"bytes"
	"errors"
)

func FlagAttributeName(flag string) (string, bool) {
	switch flag {
	case keyFirstRow, "firstRow":
		return "firstRow", true
	case "band_row", "bandRow":
		return "bandRow", true
	case "first_col", "firstCol":
		return "firstCol", true
	case "last_row", "lastRow":
		return "lastRow", true
	case "last_col", "lastCol":
		return "lastCol", true
	case keyBandCol, "bandCol":
		return "bandCol", true
	default:
		return "", false
	}
}

func BuildTableInfo(frame []byte) (map[string]any, error) {
	parsed, err := ParseTable(frame)
	if err != nil {
		return nil, err
	}

	rowCount, colCount := Dimensions(parsed)
	cells := make([]map[string]any, 0, rowCount*max(colCount, 1))
	rowHeights := make([]int64, rowCount)
	for i := range rowCount {
		if i < len(parsed.Rows) {
			rowHeights[i] = parsed.Rows[i].Height
		}
	}
	columnWidths := make([]int64, colCount)
	for i := range colCount {
		if i < len(parsed.Grid.Cols) {
			columnWidths[i] = parsed.Grid.Cols[i].Width
		}
	}

	for rIdx, row := range parsed.Rows {
		for cIdx, cell := range row.Cells {
			cells = append(cells, buildTableCellInfo(rIdx, cIdx, cell))
		}
	}
	rowsView, colsView := buildTableTraversalViews(
		rowCount,
		colCount,
		cells,
		rowHeights,
		columnWidths,
	)

	return map[string]any{
		"table": map[string]any{
			"row_count":     rowCount,
			"col_count":     colCount,
			keyFirstRow:     TruthyAttr(parsed.TblPr.FirstRow),
			"first_col":     TruthyAttr(parsed.TblPr.FirstCol),
			"last_row":      TruthyAttr(parsed.TblPr.LastRow),
			"last_col":      TruthyAttr(parsed.TblPr.LastCol),
			"band_row":      TruthyAttr(parsed.TblPr.BandRow),
			keyBandCol:      TruthyAttr(parsed.TblPr.BandCol),
			"row_heights":   rowHeights,
			"column_widths": columnWidths,
			keyCells:        cells,
			"rows":          rowsView,
			"columns":       colsView,
		},
	}, nil
}

func UpdateTableFlagsInFrame(frame []byte, flags map[string]any) ([]byte, error) {
	tblPrStart := bytes.Index(frame, []byte("<a:tblPr"))
	if tblPrStart == -1 {
		return nil, errors.New("table properties not found")
	}
	tblPrRelEnd := bytes.Index(frame[tblPrStart:], []byte(">"))
	if tblPrRelEnd == -1 {
		return nil, errors.New("invalid tblPr")
	}
	tblPrEnd := tblPrStart + tblPrRelEnd + 1
	tblPrXML := append([]byte(nil), frame[tblPrStart:tblPrEnd]...)

	for k, v := range flags {
		xmlKey, ok := FlagAttributeName(k)
		if !ok {
			continue
		}
		boolVal, ok := v.(bool)
		if !ok {
			continue
		}
		val := "0"
		if boolVal {
			val = "1"
		}
		tblPrXML = SetOrInsertAttr(tblPrXML, xmlKey, val)
	}

	updatedFrame := make([]byte, 0, len(frame)-((tblPrEnd-tblPrStart)-len(tblPrXML)))
	updatedFrame = append(updatedFrame, frame[:tblPrStart]...)
	updatedFrame = append(updatedFrame, tblPrXML...)
	updatedFrame = append(updatedFrame, frame[tblPrEnd:]...)
	return updatedFrame, nil
}
