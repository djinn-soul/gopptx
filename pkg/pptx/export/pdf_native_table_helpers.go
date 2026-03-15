package export

import (
	"math"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

const (
	tableDefaultPaddingPt = 4.0
	tableMinUsableWidthPt = 24.0
	tableCenterDivisor    = 2.0
)

type tableCellBorderSpec struct {
	widthPt float64
	color   string
}

func tableCellPadding(marginPt *float64) float64 {
	if marginPt == nil || *marginPt <= 0 {
		return tableDefaultPaddingPt
	}
	return *marginPt
}

func tableCellTextX(pdf *gopdf.GoPdf, cell tables.TableCell, line string, cellX, cellW float64) float64 {
	leftPad := tableCellPadding(cell.MarginLeftPt)
	rightPad := tableCellPadding(cell.MarginRightPt)
	contentW := max(cellW-leftPad-rightPad, tableMinUsableWidthPt)
	switch cell.Align {
	case tables.TableAlignCenter:
		lineW := measuredWidthWithMetrics(pdf, line, "")
		return cellX + leftPad + math.Max((contentW-lineW)/tableCenterDivisor, 0)
	case tables.TableAlignRight:
		lineW := measuredWidthWithMetrics(pdf, line, "")
		return cellX + leftPad + math.Max(contentW-lineW, 0)
	default:
		return cellX + leftPad
	}
}

func drawPDFTableCellBorders(
	pdf *gopdf.GoPdf,
	cell tables.TableCell,
	x, y, w, h float64,
	defaultColor string,
	defaultWidth float64,
) {
	left := resolveTableCellBorder(cell, cell.BorderLeft, defaultColor, defaultWidth)
	right := resolveTableCellBorder(cell, cell.BorderRight, defaultColor, defaultWidth)
	top := resolveTableCellBorder(cell, cell.BorderTop, defaultColor, defaultWidth)
	bottom := resolveTableCellBorder(cell, cell.BorderBottom, defaultColor, defaultWidth)
	drawTableBorderLine(pdf, x, y, x, y+h, left)
	drawTableBorderLine(pdf, x+w, y, x+w, y+h, right)
	drawTableBorderLine(pdf, x, y, x+w, y, top)
	drawTableBorderLine(pdf, x, y+h, x+w, y+h, bottom)
}

func resolveTableCellBorder(
	cell tables.TableCell,
	side *tables.TableCellBorder,
	defaultColor string,
	defaultWidth float64,
) tableCellBorderSpec {
	spec := tableCellBorderSpec{color: defaultColor, widthPt: defaultWidth}
	if cell.BorderWidthPt > 0 {
		spec.widthPt = cell.BorderWidthPt
	}
	if cell.BorderColor != "" {
		spec.color = cell.BorderColor
	}
	if side != nil {
		if side.WidthPt <= 0 {
			spec.widthPt = 0
		} else {
			spec.widthPt = side.WidthPt
		}
		if side.Color != "" {
			spec.color = side.Color
		}
	}
	return spec
}

func drawTableBorderLine(pdf *gopdf.GoPdf, x1, y1, x2, y2 float64, border tableCellBorderSpec) {
	if border.widthPt <= 0 {
		return
	}
	r, g, b := hexToRGB(border.color)
	pdf.SetStrokeColor(r, g, b)
	pdf.SetLineWidth(border.widthPt)
	pdf.Line(x1, y1, x2, y2)
}
