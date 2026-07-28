package table

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Border side names accepted by UpdateTableCellBordersInFrame.
const (
	borderSideLeft   = "left"
	borderSideRight  = "right"
	borderSideTop    = "top"
	borderSideBottom = "bottom"
)

// Line join names shared by the writer and the reader.
const (
	borderJoinRound = "round"
	borderJoinBevel = "bevel"
	borderJoinMiter = "miter"
)

// borderCapRound is the OOXML token for a round cap, also accepted as input.
const borderCapRound = "rnd"

// CellBorderSideUpdate holds the new border properties for a single cell border side.
// Pass a nil pointer to UpdateTableCellBordersInFrame to clear (remove) the border.
type CellBorderSideUpdate struct {
	Width int64
	Color string
	Dash  string
	// Cap is the line cap: "rnd", "sq" or "flat".
	Cap string
	// Join is the line join: "round", "bevel" or "miter".
	Join string
	// MiterLimitPct scales a miter join; ignored for other joins.
	MiterLimitPct float64
	// Compound is the compound line type: "sng", "dbl", "thickThin",
	// "thinThick" or "tri".
	Compound string
	// Inset draws the pen inside the cell boundary instead of centred on it.
	Inset bool
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
	case borderSideLeft:
		return "lnL", true
	case borderSideRight:
		return "lnR", true
	case borderSideTop:
		return "lnT", true
	case borderSideBottom:
		return "lnB", true
	}
	return "", false
}

// buildBorderLineXML renders one a:lnL/R/T/B. CT_LineProperties fixes the child
// order: fill, then prstDash, then the join.
func buildBorderLineXML(tag string, update *CellBorderSideUpdate) string {
	var b strings.Builder
	b.WriteString(`<a:`)
	b.WriteString(tag)
	if update.Width > 0 {
		b.WriteString(` w="`)
		b.WriteString(strconv.FormatInt(update.Width, 10))
		b.WriteString(`"`)
	}
	if lineCap := normalizeBorderCap(update.Cap); lineCap != "" {
		b.WriteString(` cap="` + lineCap + `"`)
	}
	if compound := normalizeBorderCompound(update.Compound); compound != "" {
		b.WriteString(` cmpd="` + compound + `"`)
	}
	if update.Inset {
		b.WriteString(` algn="in"`)
	}

	joinXML := borderJoinXML(update)
	if update.Color == "" && update.Dash == "" && joinXML == "" {
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
	b.WriteString(joinXML)
	b.WriteString(`</a:`)
	b.WriteString(tag)
	b.WriteString(`>`)
	return b.String()
}

func normalizeBorderCap(lineCap string) string {
	switch strings.ToLower(strings.TrimSpace(lineCap)) {
	case borderCapRound, borderJoinRound:
		return borderCapRound
	case "sq", "square":
		return "sq"
	default:
		return ""
	}
}

func normalizeBorderCompound(compound string) string {
	switch strings.ToLower(strings.TrimSpace(compound)) {
	case "dbl", "double":
		return "dbl"
	case "thickthin", "thick-thin", "thick_thin":
		return "thickThin"
	case "thinthick", "thin-thick", "thin_thick":
		return "thinThick"
	case "tri", "triple":
		return "tri"
	default:
		return ""
	}
}

const defaultBorderMiterLimitPct = 800000.0

func borderJoinXML(update *CellBorderSideUpdate) string {
	switch strings.ToLower(strings.TrimSpace(update.Join)) {
	case borderJoinRound, borderCapRound:
		return `<a:round/>`
	case borderJoinBevel:
		return `<a:bevel/>`
	case borderJoinMiter:
		limit := update.MiterLimitPct
		if limit <= 0 {
			limit = defaultBorderMiterLimitPct
		}
		return `<a:miter lim="` + strconv.FormatInt(int64(limit), 10) + `"/>`
	default:
		return ""
	}
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
		inner = insertBorderInSchemaOrder(inner, tag, buildBorderLineXML(tag, update))
	}

	openFull := cellContent[tcPrStart : tcPrTagEnd+1]
	newTcPr := string(openFull) + string(inner) + "</a:tcPr>"
	result := make([]byte, 0, len(cellContent)-(tcPrEnd-tcPrStart)+len(newTcPr))
	result = append(result, cellContent[:tcPrStart]...)
	result = append(result, []byte(newTcPr)...)
	result = append(result, cellContent[tcPrEnd:]...)
	return result, nil
}
