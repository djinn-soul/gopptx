package table

import (
	"bytes"
	"errors"
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func FlagAttributeName(flag string) (string, bool) {
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

	rowCount, colCount := Dimensions(parsed)
	cells := make([]map[string]any, 0, rowCount*max(colCount, 1))

	for rIdx, row := range parsed.Rows {
		for cIdx, cell := range row.Cells {
			cells = append(cells, buildTableCellInfo(rIdx, cIdx, cell))
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

func buildTableCellInfo(rowIndex, colIndex int, cell CellXML) map[string]any {
	rowSpan := normalizeSpan(cell.RowSpan)
	colSpan := normalizeSpan(cell.GridSpan)
	vMerge := TruthyAttr(cell.VMerge)
	hMerge := TruthyAttr(cell.HMerge)
	return map[string]any{
		"row":             rowIndex,
		"col":             colIndex,
		"row_span":        rowSpan,
		"col_span":        colSpan,
		"v_merge":         vMerge,
		"h_merge":         hMerge,
		"is_merge_origin": rowSpan > 1 || colSpan > 1,
		"is_spanned":      vMerge || hMerge,
		"text":            tableCellText(cell),
	}
}

func normalizeSpan(span int) int {
	if span <= 0 {
		return 1
	}
	return span
}

func tableCellText(cell CellXML) string {
	var textBuf bytes.Buffer
	for i, p := range cell.TxBody.Paragraphs {
		if i > 0 {
			textBuf.WriteString("\n")
		}
		for _, r := range p.Runs {
			textBuf.WriteString(r.Text)
		}
	}
	return textBuf.String()
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

// UpdateTableCellTextInFrame modifies the text of a single table cell.
// NOTE: This implementation replaces all existing paragraphs and runs within the cell
// with a single paragraph containing the new text, while preserving cell-level
// formatting properties like vertical alignment.
func UpdateTableCellTextInFrame(frame []byte, rowIdx, colIdx int, text string) ([]byte, error) {
	parsed, err := ParseTable(frame)
	if err != nil {
		return nil, err
	}
	rows, cols := Dimensions(parsed)
	if rowIdx < 0 || rowIdx >= rows || colIdx < 0 || colIdx >= cols {
		return nil, fmt.Errorf("table cell [%d,%d] out of range", rowIdx, colIdx)
	}

	escapedText := common.XMLEscape(text)

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
				[]byte(escapedText),
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
	// In a production environment, use a proper XML stack-based scanner if nesting is possible.
	end := bytes.Index(content[start:], closeTag)
	if end == -1 {
		// Fallback to just the opening tag if no closing tag found (invalid XML but avoids panic)
		return content[start : start+tagEnd+1]
	}
	return content[start : start+end+len(closeTag)]
}
