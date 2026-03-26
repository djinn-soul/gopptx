package table

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

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
