package table

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// fontSzScale converts font size in points to OOXML hundredths-of-points.
const fontSzScale = 100

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
			"first_row":     TruthyAttr(parsed.TblPr.FirstRow),
			"first_col":     TruthyAttr(parsed.TblPr.FirstCol),
			"last_row":      TruthyAttr(parsed.TblPr.LastRow),
			"last_col":      TruthyAttr(parsed.TblPr.LastCol),
			"band_row":      TruthyAttr(parsed.TblPr.BandRow),
			"band_col":      TruthyAttr(parsed.TblPr.BandCol),
			"row_heights":   rowHeights,
			"column_widths": columnWidths,
			"cells":         cells,
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

func UpdateTableRowHeightInFrame(frame []byte, rowIdx int, height int64) ([]byte, error) {
	if height <= 0 {
		return nil, errors.New("row height must be > 0")
	}
	parsed, err := ParseTable(frame)
	if err != nil {
		return nil, err
	}
	rowCount, _ := Dimensions(parsed)
	if rowIdx < 0 || rowIdx >= rowCount {
		return nil, fmt.Errorf("table row %d out of range", rowIdx)
	}
	return MutateTableRows(frame, rowIdx, rowIdx, func(_ int, rowContent []byte) ([]byte, error) {
		tagEnd := bytes.Index(rowContent, []byte(">"))
		if tagEnd == -1 {
			return nil, errors.New("invalid table row xml")
		}
		openTag := rowContent[:tagEnd+1]
		updatedTag := SetOrInsertAttr(openTag, "h", strconv.FormatInt(height, 10))
		updated := make([]byte, 0, len(rowContent)-len(openTag)+len(updatedTag))
		updated = append(updated, updatedTag...)
		updated = append(updated, rowContent[tagEnd+1:]...)
		return updated, nil
	})
}

func UpdateTableColumnWidthInFrame(frame []byte, colIdx int, width int64) ([]byte, error) {
	if width <= 0 {
		return nil, errors.New("column width must be > 0")
	}
	parsed, err := ParseTable(frame)
	if err != nil {
		return nil, err
	}
	_, colCount := Dimensions(parsed)
	if colIdx < 0 || colIdx >= colCount {
		return nil, fmt.Errorf("table column %d out of range", colIdx)
	}
	return MutateTableElements(
		frame,
		[]byte("<a:gridCol"),
		[]byte("/>"),
		colIdx,
		colIdx,
		"column",
		func(_ int, colContent []byte) ([]byte, error) {
			updated := SetOrInsertAttr(colContent, "w", strconv.FormatInt(width, 10))
			return updated, nil
		},
	)
}

// CellContentUpdate holds the fields for a combined text+style cell update.
// A nil Text means preserve the existing cell text.
type CellContentUpdate struct {
	Text     *string
	SizePt   float64
	FontName string
}

// UpdateTableCellContentInFrame updates a cell's text and/or run-level style (font size, font name).
// When Text is nil, the existing cell text is preserved. When SizePt or FontName are set,
// the rPr element is emitted with those attributes.
func UpdateTableCellContentInFrame(frame []byte, rowIdx, colIdx int, update CellContentUpdate) ([]byte, error) {
	parsed, err := ParseTable(frame)
	if err != nil {
		return nil, err
	}
	rows, cols := Dimensions(parsed)
	if rowIdx < 0 || rowIdx >= rows || colIdx < 0 || colIdx >= cols {
		return nil, fmt.Errorf("table cell [%d,%d] out of range", rowIdx, colIdx)
	}

	textToUse := ""
	if update.Text != nil {
		textToUse = *update.Text
	} else if rowIdx < len(parsed.Rows) && colIdx < len(parsed.Rows[rowIdx].Cells) {
		cell := parsed.Rows[rowIdx].Cells[colIdx]
		var sb strings.Builder
		for _, para := range cell.TxBody.Paragraphs {
			for _, run := range para.Runs {
				sb.WriteString(run.Text)
			}
		}
		textToUse = sb.String()
	}

	rPr := buildCellRPr(update.SizePt, update.FontName)
	escapedText := common.XMLEscape(textToUse)

	return MutateTableRows(frame, rowIdx, rowIdx, func(_ int, rowContent []byte) ([]byte, error) {
		return MutateTableCells(rowContent, colIdx, colIdx, func(_ int, cellContent []byte) ([]byte, error) {
			return replaceCellTxBody(cellContent, rPr, escapedText)
		})
	})
}

// replaceCellTxBody rewrites the <a:txBody> of a single cell, preserving
// existing <a:bodyPr> and <a:lstStyle> children.
func replaceCellTxBody(cellContent []byte, rPr, escapedText string) ([]byte, error) {
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
		[]byte("<a:p><a:r>"),
		[]byte(rPr),
		[]byte("<a:t>"),
		[]byte(escapedText),
		[]byte("</a:t></a:r></a:p></a:txBody>"),
	}, nil)

	updated := make([]byte, 0, len(cellContent)-((txEnd-txStart)-len(newTxBody)))
	updated = append(updated, cellContent[:txStart]...)
	updated = append(updated, newTxBody...)
	updated = append(updated, cellContent[txEnd:]...)
	return updated, nil
}

// AddTableRowInFrame appends a new empty row to the table XML.
// height is in EMU; pass 0 to omit the h attribute (PowerPoint will auto-size).
func AddTableRowInFrame(frame []byte, height int64) ([]byte, error) {
	parsed, err := ParseTable(frame)
	if err != nil {
		return nil, err
	}
	_, colCount := Dimensions(parsed)
	if colCount == 0 {
		return nil, errors.New("table has no columns")
	}

	// Build empty cell XML
	emptyCell := `<a:tc><a:txBody><a:bodyPr/><a:lstStyle/><a:p><a:r><a:rPr/><a:t></a:t></a:r></a:p></a:txBody><a:tcPr/></a:tc>`

	var newRow strings.Builder
	if height > 0 {
		newRow.WriteString(`<a:tr h="`)
		newRow.WriteString(strconv.FormatInt(height, 10))
		newRow.WriteString(`">`)
	} else {
		newRow.WriteString(`<a:tr>`)
	}
	for range colCount {
		newRow.WriteString(emptyCell)
	}
	newRow.WriteString(`</a:tr>`)

	// Insert before </a:tbl>
	tblEnd := bytes.Index(frame, []byte("</a:tbl>"))
	if tblEnd == -1 {
		return nil, errors.New("invalid table xml: missing </a:tbl>")
	}
	newRowBytes := []byte(newRow.String())
	updated := make([]byte, 0, len(frame)+len(newRowBytes))
	updated = append(updated, frame[:tblEnd]...)
	updated = append(updated, newRowBytes...)
	updated = append(updated, frame[tblEnd:]...)
	return updated, nil
}

// AddTableColumnInFrame appends a new empty column to the table XML.
// width is in EMU.
func AddTableColumnInFrame(frame []byte, width int64) ([]byte, error) {
	if width <= 0 {
		return nil, errors.New("column width must be > 0")
	}
	parsed, err := ParseTable(frame)
	if err != nil {
		return nil, err
	}
	rowCount, _ := Dimensions(parsed)

	// 1. Insert new <a:gridCol> before </a:tblGrid>
	gridColXML := []byte(`<a:gridCol w="` + strconv.FormatInt(width, 10) + `"/>`)
	tblGridEnd := bytes.Index(frame, []byte("</a:tblGrid>"))
	if tblGridEnd == -1 {
		return nil, errors.New("invalid table xml: missing </a:tblGrid>")
	}
	step1 := make([]byte, 0, len(frame)+len(gridColXML))
	step1 = append(step1, frame[:tblGridEnd]...)
	step1 = append(step1, gridColXML...)
	step1 = append(step1, frame[tblGridEnd:]...)

	if rowCount == 0 {
		return step1, nil
	}

	// 2. Append an empty <a:tc> to every row
	const emptyCellXML = `<a:tc><a:txBody><a:bodyPr/><a:lstStyle/>` +
		`<a:p><a:r><a:rPr/><a:t></a:t></a:r></a:p></a:txBody><a:tcPr/></a:tc>`
	emptyCell := []byte(emptyCellXML)
	result, err := MutateTableRows(step1, 0, rowCount-1, func(_ int, rowContent []byte) ([]byte, error) {
		closeRow := []byte("</a:tr>")
		pos := bytes.Index(rowContent, closeRow)
		if pos == -1 {
			return nil, errors.New("invalid row xml: missing </a:tr>")
		}
		updated := make([]byte, 0, len(rowContent)+len(emptyCell))
		updated = append(updated, rowContent[:pos]...)
		updated = append(updated, emptyCell...)
		updated = append(updated, rowContent[pos:]...)
		return updated, nil
	})
	return result, err
}

// CellBorderSideUpdate holds the new border properties for a single cell border side.
// Pass a nil pointer to UpdateTableCellBordersInFrame to clear (remove) the border.
type CellBorderSideUpdate struct {
	Width int64
	Color string
	Dash  string
}

// UpdateTableCellBordersInFrame updates a single border side of a table cell.
// side must be "left", "right", "top", or "bottom".
// update=nil removes the border element for that side.
func UpdateTableCellBordersInFrame(
	frame []byte, rowIdx, colIdx int, side string, update *CellBorderSideUpdate,
) ([]byte, error) {
	tag, ok := borderSideTag(side)
	if !ok {
		return nil, fmt.Errorf("invalid border side %q", side)
	}
	parsed, err := ParseTable(frame)
	if err != nil {
		return nil, err
	}
	rows, cols := Dimensions(parsed)
	if rowIdx < 0 || rowIdx >= rows || colIdx < 0 || colIdx >= cols {
		return nil, fmt.Errorf("table cell [%d,%d] out of range", rowIdx, colIdx)
	}
	return MutateTableRows(frame, rowIdx, rowIdx, func(_ int, rowContent []byte) ([]byte, error) {
		return MutateTableCells(rowContent, colIdx, colIdx, func(_ int, cellContent []byte) ([]byte, error) {
			return applyCellBorder(cellContent, tag, update)
		})
	})
}

func borderSideTag(side string) (string, bool) {
	switch side {
	case "left":
		return "lnL", true
	case "right":
		return "lnR", true
	case "top":
		return "lnT", true
	case "bottom":
		return "lnB", true
	}
	return "", false
}

func buildBorderLineXML(tag string, update *CellBorderSideUpdate) string {
	var b strings.Builder
	b.WriteString(`<a:`)
	b.WriteString(tag)
	if update.Width > 0 {
		b.WriteString(` w="`)
		b.WriteString(strconv.FormatInt(update.Width, 10))
		b.WriteString(`"`)
	}
	if update.Color == "" && update.Dash == "" {
		b.WriteString(`/>`)
		return b.String()
	}
	b.WriteString(`>`)
	if update.Color != "" {
		b.WriteString(`<a:solidFill><a:srgbClr val="`)
		b.WriteString(update.Color)
		b.WriteString(`"/></a:solidFill>`)
	}
	if update.Dash != "" {
		b.WriteString(`<a:prstDash val="`)
		b.WriteString(update.Dash)
		b.WriteString(`"/>`)
	}
	b.WriteString(`</a:`)
	b.WriteString(tag)
	b.WriteString(`>`)
	return b.String()
}

func applyCellBorder(cellContent []byte, tag string, update *CellBorderSideUpdate) ([]byte, error) {
	openTag := []byte("<a:tcPr")
	tcPrStart := bytes.Index(cellContent, openTag)
	if tcPrStart == -1 {
		if update == nil {
			return cellContent, nil
		}
		closeTC := []byte("</a:tc>")
		pos := bytes.Index(cellContent, closeTC)
		if pos == -1 {
			return nil, errors.New("invalid cell xml: missing </a:tc>")
		}
		newTcPr := "<a:tcPr>" + buildBorderLineXML(tag, update) + "</a:tcPr>"
		result := make([]byte, 0, len(cellContent)+len(newTcPr))
		result = append(result, cellContent[:pos]...)
		result = append(result, []byte(newTcPr)...)
		result = append(result, cellContent[pos:]...)
		return result, nil
	}

	tagEndRel := bytes.Index(cellContent[tcPrStart:], []byte(">"))
	if tagEndRel == -1 {
		return nil, errors.New("invalid tcPr xml: missing >")
	}
	tcPrTagEnd := tcPrStart + tagEndRel

	if cellContent[tcPrTagEnd-1] == '/' {
		// Self-closing <a:tcPr/>
		if update == nil {
			return cellContent, nil
		}
		newTcPr := "<a:tcPr>" + buildBorderLineXML(tag, update) + "</a:tcPr>"
		result := make([]byte, 0, len(cellContent)-(tcPrTagEnd+1-tcPrStart)+len(newTcPr))
		result = append(result, cellContent[:tcPrStart]...)
		result = append(result, []byte(newTcPr)...)
		result = append(result, cellContent[tcPrTagEnd+1:]...)
		return result, nil
	}

	// <a:tcPr>...inner...</a:tcPr>
	closeTcPr := []byte("</a:tcPr>")
	tcPrCloseRel := bytes.Index(cellContent[tcPrStart:], closeTcPr)
	if tcPrCloseRel == -1 {
		return nil, errors.New("invalid tcPr xml: missing </a:tcPr>")
	}
	tcPrCloseStart := tcPrStart + tcPrCloseRel
	tcPrEnd := tcPrCloseStart + len(closeTcPr)

	inner := append([]byte(nil), cellContent[tcPrTagEnd+1:tcPrCloseStart]...)
	inner = removeSingleXMLElement(inner, []byte("<a:"+tag), []byte("</a:"+tag+">"))
	if update != nil {
		inner = append(inner, []byte(buildBorderLineXML(tag, update))...)
	}

	openFull := cellContent[tcPrStart : tcPrTagEnd+1]
	newTcPr := string(openFull) + string(inner) + "</a:tcPr>"
	result := make([]byte, 0, len(cellContent)-(tcPrEnd-tcPrStart)+len(newTcPr))
	result = append(result, cellContent[:tcPrStart]...)
	result = append(result, []byte(newTcPr)...)
	result = append(result, cellContent[tcPrEnd:]...)
	return result, nil
}

func removeSingleXMLElement(content []byte, openTag []byte, closeTag []byte) []byte {
	start := bytes.Index(content, openTag)
	if start == -1 {
		return content
	}
	tagEndRel := bytes.Index(content[start:], []byte(">"))
	if tagEndRel == -1 {
		return content
	}
	tagEnd := start + tagEndRel
	var removeEnd int
	if content[tagEnd-1] == '/' {
		removeEnd = tagEnd + 1
	} else {
		closeRel := bytes.Index(content[start:], closeTag)
		if closeRel == -1 {
			return content
		}
		removeEnd = start + closeRel + len(closeTag)
	}
	result := make([]byte, 0, len(content)-(removeEnd-start))
	result = append(result, content[:start]...)
	result = append(result, content[removeEnd:]...)
	return result
}

func buildCellRPr(sizePt float64, fontName string) string {
	if sizePt <= 0 && strings.TrimSpace(fontName) == "" {
		return "<a:rPr/>"
	}
	var b strings.Builder
	b.WriteString(`<a:rPr lang="en-US" dirty="0"`)
	if sizePt > 0 {
		b.WriteString(` sz="`)
		b.WriteString(strconv.Itoa(int(sizePt * fontSzScale)))
		b.WriteString(`"`)
	}
	b.WriteString(`>`)
	if strings.TrimSpace(fontName) != "" {
		b.WriteString(`<a:latin typeface="`)
		b.WriteString(common.XMLEscape(fontName))
		b.WriteString(`"/>`)
	}
	b.WriteString(`</a:rPr>`)
	return b.String()
}
