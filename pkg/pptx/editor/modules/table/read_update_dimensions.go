package table

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

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

// InsertTableRowInFrame inserts a new empty row before the row at atIndex.
// If atIndex equals the current row count the row is appended at the end.
// height is in EMU; pass 0 to omit the h attribute (PowerPoint will auto-size).
func InsertTableRowInFrame(frame []byte, atIndex int, height int64) ([]byte, error) {
	parsed, err := ParseTable(frame)
	if err != nil {
		return nil, err
	}
	rowCount, colCount := Dimensions(parsed)
	if atIndex < 0 || atIndex > rowCount {
		return nil, fmt.Errorf("insert index %d out of range [0, %d]", atIndex, rowCount)
	}
	if colCount == 0 {
		return nil, errors.New("table has no columns")
	}
	if atIndex == rowCount {
		return AddTableRowInFrame(frame, height)
	}

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
	newRowBytes := []byte(newRow.String())

	return insertBeforeNthElement(frame, []byte("<a:tr"), []byte("</a:tr>"), newRowBytes, atIndex)
}

// RemoveTableRowInFrame removes the row at atIndex from the table XML.
func RemoveTableRowInFrame(frame []byte, atIndex int) ([]byte, error) {
	parsed, err := ParseTable(frame)
	if err != nil {
		return nil, err
	}
	rowCount, _ := Dimensions(parsed)
	if atIndex < 0 || atIndex >= rowCount {
		return nil, fmt.Errorf("row index %d out of range [0, %d]", atIndex, rowCount-1)
	}
	return MutateTableElements(frame, []byte("<a:tr"), []byte("</a:tr>"), atIndex, atIndex, "row",
		func(_ int, _ []byte) ([]byte, error) { return nil, nil },
	)
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

	const emptyCellXML = `<a:tc><a:txBody><a:bodyPr/><a:lstStyle/>` +
		`<a:p><a:r><a:rPr/><a:t></a:t></a:r></a:p></a:txBody><a:tcPr/></a:tc>`
	emptyCell := []byte(emptyCellXML)
	return MutateTableRows(step1, 0, rowCount-1, func(_ int, rowContent []byte) ([]byte, error) {
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
}

// InsertTableColumnInFrame inserts a new empty column before the column at atIndex.
// If atIndex equals the current column count the column is appended at the end.
// width is in EMU.
func InsertTableColumnInFrame(frame []byte, atIndex int, width int64) ([]byte, error) {
	if width <= 0 {
		return nil, errors.New("column width must be > 0")
	}
	parsed, err := ParseTable(frame)
	if err != nil {
		return nil, err
	}
	rowCount, colCount := Dimensions(parsed)
	if atIndex < 0 || atIndex > colCount {
		return nil, fmt.Errorf("insert index %d out of range [0, %d]", atIndex, colCount)
	}
	if atIndex == colCount {
		return AddTableColumnInFrame(frame, width)
	}

	gridColXML := []byte(`<a:gridCol w="` + strconv.FormatInt(width, 10) + `"/>`)
	step1, err := insertBeforeNthElement(frame, []byte("<a:gridCol"), []byte("/>"), gridColXML, atIndex)
	if err != nil {
		return nil, err
	}
	if rowCount == 0 {
		return step1, nil
	}

	emptyCell := []byte(`<a:tc><a:txBody><a:bodyPr/><a:lstStyle/>` +
		`<a:p><a:r><a:rPr/><a:t></a:t></a:r></a:p></a:txBody><a:tcPr/></a:tc>`)
	return MutateTableRows(step1, 0, rowCount-1, func(_ int, rowContent []byte) ([]byte, error) {
		return insertBeforeNthElement(rowContent, []byte("<a:tc"), []byte("</a:tc>"), emptyCell, atIndex)
	})
}

// RemoveTableColumnInFrame removes the column at atIndex from the table XML.
func RemoveTableColumnInFrame(frame []byte, atIndex int) ([]byte, error) {
	parsed, err := ParseTable(frame)
	if err != nil {
		return nil, err
	}
	rowCount, colCount := Dimensions(parsed)
	if atIndex < 0 || atIndex >= colCount {
		return nil, fmt.Errorf("column index %d out of range [0, %d]", atIndex, colCount-1)
	}
	if colCount == 1 {
		return nil, errors.New("cannot remove the last column")
	}

	// Remove the gridCol entry from tblGrid.
	step1, err := MutateTableElements(frame, []byte("<a:gridCol"), []byte("/>"), atIndex, atIndex, "column",
		func(_ int, _ []byte) ([]byte, error) { return nil, nil },
	)
	if err != nil {
		return nil, err
	}
	if rowCount == 0 {
		return step1, nil
	}

	// Remove the cell at atIndex from each row.
	return MutateTableRows(step1, 0, rowCount-1, func(_ int, rowContent []byte) ([]byte, error) {
		return MutateTableElements(rowContent, []byte("<a:tc"), []byte("</a:tc>"), atIndex, atIndex, "cell",
			func(_ int, _ []byte) ([]byte, error) { return nil, nil },
		)
	})
}
