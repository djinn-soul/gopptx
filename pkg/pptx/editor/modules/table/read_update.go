package table

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
)

func TableFlagAttributeName(flag string) (string, bool) {
	switch flag {
	case "first_row", "firstRow":
		return "firstRow", true
	case "band_row", "bandRow":
		return "bandRow", true
	case "first_col", "firstCol":
		return "firstCol", true
	case "last_row", "lastRow":
		return "lastRow", true
	case "last_col", "lastCol":
		return "lastCol", true
	case "band_col", "bandCol":
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

	rowCount, colCount := TableDimensions(parsed)
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
			isSpanned := TruthyAttr(cell.VMerge) || TruthyAttr(cell.HMerge)

			cells = append(cells, map[string]any{
				"row":             rIdx,
				"col":             cIdx,
				"row_span":        rowSpan,
				"col_span":        colSpan,
				"v_merge":         TruthyAttr(cell.VMerge),
				"h_merge":         TruthyAttr(cell.HMerge),
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
			"first_row": TruthyAttr(parsed.TblPr.FirstRow),
			"first_col": TruthyAttr(parsed.TblPr.FirstCol),
			"last_row":  TruthyAttr(parsed.TblPr.LastRow),
			"last_col":  TruthyAttr(parsed.TblPr.LastCol),
			"band_row":  TruthyAttr(parsed.TblPr.BandRow),
			"band_col":  TruthyAttr(parsed.TblPr.BandCol),
			"cells":     cells,
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
		xmlKey, ok := TableFlagAttributeName(k)
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

// UpdateTableCellTextInFrame modifies the text of a single table cell.
// NOTE: This implementation replaces all existing paragraphs and runs within the cell
// with a single paragraph containing the new text, while preserving cell-level
// formatting properties like vertical alignment.
func UpdateTableCellTextInFrame(frame []byte, rowIdx, colIdx int, text string) ([]byte, error) {
	parsed, err := ParseTable(frame)
	if err != nil {
		return nil, err
	}
	rows, cols := TableDimensions(parsed)
	if rowIdx < 0 || rowIdx >= rows || colIdx < 0 || colIdx >= cols {
		return nil, fmt.Errorf("table cell [%d,%d] out of range", rowIdx, colIdx)
	}

	var escaped bytes.Buffer
	if err := xml.EscapeText(&escaped, []byte(text)); err != nil {
		return nil, fmt.Errorf("escape cell text: %w", err)
	}

	updatedFrame, err := MutateTableRows(frame, rowIdx, rowIdx, func(_ int, rowContent []byte) ([]byte, error) {
		return MutateTableCells(rowContent, colIdx, colIdx, func(_ int, cellContent []byte) ([]byte, error) {
			txStart := bytes.Index(cellContent, []byte("<a:txBody>"))
			if txStart == -1 {
				return nil, errors.New("txBody not found in cell")
			}
			txEndRel := bytes.Index(cellContent[txStart:], []byte("</a:txBody>"))
			if txEndRel == -1 {
				return nil, errors.New("invalid txBody xml")
			}
			txEnd := txStart + txEndRel + len("</a:txBody>")

			oldTxBody := cellContent[txStart:txEnd]

			// Extract bodyPr and lstStyle to preserve cell-level formatting (like vertical alignment)
			bodyPr := extractXMLElement(oldTxBody, []byte("<a:bodyPr"))
			if len(bodyPr) == 0 {
				bodyPr = []byte("<a:bodyPr/>")
			}
			lstStyle := extractXMLElement(oldTxBody, []byte("<a:lstStyle"))
			if len(lstStyle) == 0 {
				lstStyle = []byte("<a:lstStyle/>")
			}

			newTxBody := bytes.Join([][]byte{
				[]byte("<a:txBody>"),
				bodyPr,
				lstStyle,
				[]byte("<a:p><a:r><a:rPr/><a:t>"),
				escaped.Bytes(),
				[]byte("</a:t></a:r></a:p></a:txBody>"),
			}, nil)

			updated := make([]byte, 0, len(cellContent)-((txEnd-txStart)-len(newTxBody)))
			updated = append(updated, cellContent[:txStart]...)
			updated = append(updated, newTxBody...)
			updated = append(updated, cellContent[txEnd:]...)
			return updated, nil
		})
	})
	if err != nil {
		return nil, err
	}
	return updatedFrame, nil
}

func extractXMLElement(content []byte, tagOpen []byte) []byte {
	start := bytes.Index(content, tagOpen)
	if start == -1 {
		return nil
	}

	// Check if it's self-closing: <tag ... />
	tagEnd := bytes.Index(content[start:], []byte(">"))
	if tagEnd == -1 {
		return nil
	}
	if tagEnd > 0 && content[start+tagEnd-1] == '/' {
		return content[start : start+tagEnd+1]
	}

	// Not self-closing, need to find the matching </tag>
	// This is a simple scanner, assumes no nested same-name tags for bodyPr/lstStyle
	tagName := bytes.TrimPrefix(tagOpen, []byte("<"))
	closeTag := append([]byte("</"), append(tagName, []byte(">")...)...)
	end := bytes.Index(content[start:], closeTag)
	if end == -1 {
		// Fallback to just the opening tag if no closing tag found (invalid XML but avoids panic)
		return content[start : start+tagEnd+1]
	}
	return content[start : start+end+len(closeTag)]
}
