package export

import (
	"math"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

const (
	tableMinUsableWidthPt = 24.0
	tableCenterDivisor    = 2.0

	// defaultTableFontSize is what PowerPoint uses for table text when the cell
	// carries no explicit size. It is 18pt, not the renderer's general 14pt
	// default, which is why unsized cells used to come out visibly small.
	defaultTableFontSize = 18
)

type tableCellBorderSpec struct {
	widthPt float64
	color   string
}

// tableCellPaddingLR and tableCellPaddingTB default to the OOXML a:tcPr insets
// (marL/marR 91440 EMU = 7.2pt, marT/marB 45720 EMU = 3.6pt) rather than one
// padding for all four sides.
func tableCellPaddingLR(marginPt *float64) float64 {
	return tableCellPaddingOr(marginPt, defaultTextInsetLRPt)
}

func tableCellPaddingTB(marginPt *float64) float64 {
	return tableCellPaddingOr(marginPt, defaultTextInsetTBPt)
}

func tableCellPaddingOr(marginPt *float64, fallback float64) float64 {
	if marginPt == nil || *marginPt <= 0 {
		return fallback
	}
	return *marginPt
}

// tableCellFontSize honours the cell's own size and otherwise falls back to
// PowerPoint's table default.
func tableCellFontSize(cell tables.TableCell) int {
	if cell.SizePt > 0 {
		return int(cell.SizePt)
	}
	return defaultTableFontSize
}

func tableCellTextX(pdf *gopdf.GoPdf, cell tables.TableCell, line string, cellX, cellW float64) float64 {
	leftPad := tableCellPaddingLR(cell.MarginLeftPt)
	rightPad := tableCellPaddingLR(cell.MarginRightPt)
	contentW := max(cellW-leftPad-rightPad, tableMinUsableWidthPt)
	switch cell.Align {
	case tables.TableAlignCenter:
		lineW := measuredWidth(pdf, line)
		return cellX + leftPad + math.Max((contentW-lineW)/tableCenterDivisor, 0)
	case tables.TableAlignRight:
		lineW := measuredWidth(pdf, line)
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
