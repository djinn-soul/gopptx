package table

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// UpdateTableCellTextInFrame modifies the text of a single table cell.
// NOTE: This implementation replaces all existing paragraphs and runs within the cell
// with a single paragraph containing the new text, while preserving cell-level
// formatting properties like vertical alignment.
func UpdateTableCellTextInFrame(frame []byte, rowIdx, colIdx int, text string) ([]byte, error) {
	return UpdateTableCellContentInFrame(frame, rowIdx, colIdx, CellContentUpdate{Text: &text})
}

func extractXMLElement(content []byte, tagOpen []byte) []byte {
	start := bytes.Index(content, tagOpen)
	if start == -1 {
		return nil
	}

	tagEnd := bytes.Index(content[start:], []byte(">"))
	if tagEnd == -1 {
		return nil
	}
	if tagEnd > 0 && content[start+tagEnd-1] == '/' {
		return content[start : start+tagEnd+1]
	}

	tagName := bytes.TrimPrefix(tagOpen, []byte("<"))
	closeTag := append([]byte("</"), append(tagName, []byte(">")...)...)
	end := bytes.Index(content[start:], closeTag)
	if end == -1 {
		return content[start : start+tagEnd+1]
	}
	return content[start : start+end+len(closeTag)]
}

// CellContentUpdate holds the fields for a combined text+style cell update.
// A nil Text means preserve the existing cell text; nil style pointers mean
// leave that aspect alone.
type CellContentUpdate struct {
	Text     *string
	SizePt   float64
	FontName string
	// Bold, Italic and Underline are tri-state: nil leaves the run as it is.
	Bold      *bool
	Italic    *bool
	Underline *bool
	// Color is the run's text colour as a hex string, e.g. "C00000".
	Color string
	// BackgroundColor fills the cell via <a:tcPr>. "none" clears the fill.
	BackgroundColor string
}

// HasRunStyle reports whether the update touches run-level formatting.
func (u CellContentUpdate) HasRunStyle() bool {
	return u.SizePt > 0 ||
		strings.TrimSpace(u.FontName) != "" ||
		u.Bold != nil || u.Italic != nil || u.Underline != nil ||
		strings.TrimSpace(u.Color) != ""
}

// UpdateTableCellContentInFrame updates a cell's text and/or run-level style.
// When Text is nil the existing cell text is preserved. Any style field the
// caller leaves unset keeps whatever the cell already had: the new run
// properties are merged onto the cell's first <a:rPr>, and its <a:pPr> is
// carried over, so refreshing a deck's numbers does not flatten hand-applied
// formatting (upstream issue #1037).
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

	escapedText := common.XMLEscape(textToUse)

	return MutateTableRows(frame, rowIdx, rowIdx, func(_ int, rowContent []byte) ([]byte, error) {
		return MutateTableCells(rowContent, colIdx, colIdx, func(_ int, cellContent []byte) ([]byte, error) {
			if update.Text == nil && !update.HasRunStyle() {
				return applyCellFill(cellContent, update.BackgroundColor)
			}
			updated, err := replaceCellTxBody(cellContent, update, escapedText)
			if err != nil {
				return nil, err
			}
			return applyCellFill(updated, update.BackgroundColor)
		})
	})
}

// replaceCellTxBody rewrites the <a:txBody> of a single cell, preserving the
// existing <a:bodyPr>, <a:lstStyle>, <a:pPr> and run properties.
func replaceCellTxBody(cellContent []byte, update CellContentUpdate, escapedText string) ([]byte, error) {
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
	// The paragraph's own properties (alignment, indent, bullet) survive a data
	// refresh, as do the run properties not named by this update.
	pPr := extractXMLElement(oldTxBody, []byte("<a:pPr"))
	rPr := mergeCellRPr(extractXMLElement(oldTxBody, []byte("<a:rPr")), update)
	// An empty paragraph carries its formatting on <a:endParaRPr>, which is the
	// only place to recover it once a run is added.
	if rPr == "<a:rPr/>" {
		if endPara := extractXMLElement(oldTxBody, []byte("<a:endParaRPr")); len(endPara) > 0 {
			rPr = mergeCellRPr(renameElement(endPara, "endParaRPr", "rPr"), update)
		}
	}

	newTxBody := bytes.Join([][]byte{
		[]byte("<a:txBody>"),
		bodyPr,
		lstStyle,
		[]byte("<a:p>"),
		pPr,
		[]byte("<a:r>"),
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

// renameElement swaps an element's tag name, keeping its attributes and body.
func renameElement(element []byte, from, to string) []byte {
	out := bytes.ReplaceAll(element, []byte("<a:"+from), []byte("<a:"+to))
	return bytes.ReplaceAll(out, []byte("</a:"+from+">"), []byte("</a:"+to+">"))
}
