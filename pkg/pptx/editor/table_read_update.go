package editor

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
)

// GetTable reads a table's structure entirely from XML.
func (e *PresentationEditor) GetTable(slideIndex, shapeID int) (map[string]any, error) {
	_, _, _, _, frame, err := getSlideTableFrame(e, slideIndex, shapeID)
	if err != nil {
		return nil, err
	}
	parsed, err := parseTable(frame)
	if err != nil {
		return nil, err
	}

	rowCount, colCount := tableDimensions(parsed)
	cells := make([]map[string]any, 0, rowCount*max(colCount, 1))

	for rIdx, row := range parsed.Rows {
		for cIdx, cell := range row.Cells {
			var textBuf bytes.Buffer
			for _, p := range cell.TxBody.Paragraphs {
				for _, r := range p.Runs {
					textBuf.WriteString(r.Text)
				}
			}

			rowSpan := cell.RowSpan
			if rowSpan <= 0 {
				rowSpan = 1
			}
			colSpan := cell.GridSpan
			if colSpan <= 0 {
				colSpan = 1
			}

			isOrigin := rowSpan > 1 || colSpan > 1
			isSpanned := truthyAttr(cell.VMerge) || truthyAttr(cell.HMerge)

			cells = append(cells, map[string]any{
				"row":             rIdx,
				"col":             cIdx,
				"row_span":        rowSpan,
				"col_span":        colSpan,
				"v_merge":         truthyAttr(cell.VMerge),
				"h_merge":         truthyAttr(cell.HMerge),
				"is_merge_origin": isOrigin,
				"is_spanned":      isSpanned,
				"text":            textBuf.String(),
			})
		}
	}

	return map[string]any{
		"table": map[string]any{
			"row_count": rowCount,
			"col_count": colCount,
			"first_row": truthyAttr(parsed.TblPr.FirstRow),
			"first_col": truthyAttr(parsed.TblPr.FirstCol),
			"last_row":  truthyAttr(parsed.TblPr.LastRow),
			"last_col":  truthyAttr(parsed.TblPr.LastCol),
			"band_row":  truthyAttr(parsed.TblPr.BandRow),
			"band_col":  truthyAttr(parsed.TblPr.BandCol),
			"cells":     cells,
		},
	}, nil
}

// UpdateTableFlags modifies properties of the table like firstRow, bandRow, etc.
func (e *PresentationEditor) UpdateTableFlags(slideIndex, shapeID int, flags map[string]any) error {
	partPath, slideContent, frameStart, frameEnd, frame, err := getSlideTableFrame(e, slideIndex, shapeID)
	if err != nil {
		return err
	}

	tblPrStart := bytes.Index(frame, []byte("<a:tblPr"))
	if tblPrStart == -1 {
		return errors.New("table properties not found")
	}
	tblPrRelEnd := bytes.Index(frame[tblPrStart:], []byte(">"))
	if tblPrRelEnd == -1 {
		return errors.New("invalid tblPr")
	}
	tblPrEnd := tblPrStart + tblPrRelEnd + 1
	tblPrXML := append([]byte(nil), frame[tblPrStart:tblPrEnd]...)

	camelMap := map[string]string{
		"first_row": "firstRow",
		"band_row":  "bandRow",
		"first_col": "firstCol",
		"last_row":  "lastRow",
		"last_col":  "lastCol",
		"band_col":  "bandCol",
	}

	for k, v := range flags {
		xmlKey, ok := camelMap[k]
		if !ok {
			xmlKey = k
		}
		boolVal, ok := v.(bool)
		if !ok {
			continue
		}
		val := "0"
		if boolVal {
			val = "1"
		}
		tblPrXML = setOrInsertAttr(tblPrXML, xmlKey, val)
	}

	updatedFrame := make([]byte, 0, len(frame)-((tblPrEnd-tblPrStart)-len(tblPrXML)))
	updatedFrame = append(updatedFrame, frame[:tblPrStart]...)
	updatedFrame = append(updatedFrame, tblPrXML...)
	updatedFrame = append(updatedFrame, frame[tblPrEnd:]...)

	e.parts.Set(partPath, replaceTableFrame(slideContent, frameStart, frameEnd, updatedFrame))
	return nil
}

func (e *PresentationEditor) UpdateTableCellText(slideIndex, shapeID, rowIdx, colIdx int, text string) error {
	partPath, slideContent, frameStart, frameEnd, frame, err := getSlideTableFrame(e, slideIndex, shapeID)
	if err != nil {
		return err
	}
	parsed, err := parseTable(frame)
	if err != nil {
		return err
	}
	rows, cols := tableDimensions(parsed)
	if rowIdx < 0 || rowIdx >= rows || colIdx < 0 || colIdx >= cols {
		return fmt.Errorf("table cell [%d,%d] out of range", rowIdx, colIdx)
	}

	var escaped bytes.Buffer
	if err := xml.EscapeText(&escaped, []byte(text)); err != nil {
		return fmt.Errorf("escape cell text: %w", err)
	}
	newTxBody := []byte(`<a:txBody><a:bodyPr/><a:lstStyle/><a:p><a:r><a:rPr/><a:t>` + escaped.String() + `</a:t></a:r></a:p></a:txBody>`)

	updatedFrame, err := mutateTableRows(frame, rowIdx, rowIdx, func(_ int, rowContent []byte) ([]byte, error) {
		return mutateTableCells(rowContent, colIdx, colIdx, func(_ int, cellContent []byte) ([]byte, error) {
			txStart := bytes.Index(cellContent, []byte("<a:txBody>"))
			if txStart == -1 {
				return nil, errors.New("txBody not found in cell")
			}
			txEndRel := bytes.Index(cellContent[txStart:], []byte("</a:txBody>"))
			if txEndRel == -1 {
				return nil, errors.New("invalid txBody xml")
			}
			txEnd := txStart + txEndRel + len("</a:txBody>")

			updated := make([]byte, 0, len(cellContent)-((txEnd-txStart)-len(newTxBody)))
			updated = append(updated, cellContent[:txStart]...)
			updated = append(updated, newTxBody...)
			updated = append(updated, cellContent[txEnd:]...)
			return updated, nil
		})
	})
	if err != nil {
		return err
	}

	e.parts.Set(partPath, replaceTableFrame(slideContent, frameStart, frameEnd, updatedFrame))
	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
