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
	return MutateTableRows(frame, rowIdx, rowIdx, func(_ int, rowContent []byte) ([]byte, error) {
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
