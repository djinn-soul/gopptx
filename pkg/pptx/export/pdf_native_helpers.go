//nolint:mnd // Table/shape helper math uses fixed visual constants for parity with PPT defaults.
package export

import (
	"math"
	"strconv"
	"strings"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

const (
	pdfTableStripeOdd  = "D0D8E8"
	pdfTableStripeEven = "E9EDF4"
	pdfTableGridColor  = "E9EDF4"
	pdfTableHeaderText = "FFFFFF"
)

//nolint:funlen,gocognit // Table rendering keeps per-cell style/layout resolution in one pass for deterministic output.
func renderPDFTable(pdf *gopdf.GoPdf, tab tables.Table) {
	startX := emuToPt(tab.X.Emu())
	startY := emuToPt(tab.Y.Emu())

	colPts := make([]float64, len(tab.ColumnWidths))
	for i, cw := range tab.ColumnWidths {
		colPts[i] = emuToPt(cw.Emu())
	}

	hasStyledHeader := len(tab.StyledRows) > 0
	defaultRowH := 20.0
	if len(tab.Rows) > 0 {
		totalH := emuToPt(tab.CY.Emu())
		if totalH > 0 {
			defaultRowH = totalH / float64(len(tab.Rows))
		}
	}
	for ri, row := range tab.Rows {
		rowH := defaultRowH
		if ri < len(tab.RowHeights) {
			h := emuToPt(tab.RowHeights[ri].Emu())
			if h > 0 {
				rowH = h
			}
		}
		cellYOffset := 0.0
		for pi := range ri {
			if pi < len(tab.RowHeights) {
				h := emuToPt(tab.RowHeights[pi].Emu())
				if h > 0 {
					cellYOffset += h
					continue
				}
			}
			cellYOffset += defaultRowH
		}
		cellY := startY + cellYOffset
		cellX := startX
		for ci, text := range row {
			cw := colPts[0]
			if ci < len(colPts) {
				cw = colPts[ci]
			}
			cell := tables.NewTableCell(text)
			if ri < len(tab.StyledRows) && ci < len(tab.StyledRows[ri]) {
				cell = tab.StyledRows[ri][ci]
				text = cell.Text
			}
			switch {
			case cell.BackgroundColor != "":
				r, g, b := hexToRGB(cell.BackgroundColor)
				pdf.SetFillColor(r, g, b)
			case hasStyledHeader && ri > 0:
				if ri%2 == 1 {
					r, g, b := hexToRGB(pdfTableStripeOdd)
					pdf.SetFillColor(r, g, b)
				} else {
					r, g, b := hexToRGB(pdfTableStripeEven)
					pdf.SetFillColor(r, g, b)
				}
			default:
				pdf.SetFillColor(255, 255, 255)
			}
			if hasStyledHeader {
				r, g, b := hexToRGB(pdfTableGridColor)
				pdf.SetStrokeColor(r, g, b)
			} else {
				pdf.SetStrokeColor(180, 180, 180)
			}
			pdf.SetLineWidth(0.5)
			pdf.RectFromUpperLeftWithStyle(cellX, cellY, cw, rowH, "DF")
			switch {
			case cell.Color != "":
				r, g, b := hexToRGB(cell.Color)
				pdf.SetTextColor(r, g, b)
			case hasStyledHeader && ri == 0 && cell.BackgroundColor != "":
				r, g, b := hexToRGB(pdfTableHeaderText)
				pdf.SetTextColor(r, g, b)
			default:
				pdf.SetTextColor(0, 0, 0)
			}
			usableW := cw - 8
			if usableW < 24 {
				usableW = 24
			}
			fontSize := fitPDFTextSize(pdf, text, defaultFontSize, minTextAutoFitSize, cell.Bold, false, usableW)
			setPDFTextFont(pdf, fontSize, cell.Bold, false)
			lines := wrapPDFText(pdf, text, usableW)
			lineHeight := math.Max(float64(fontSize)+2, 12)
			totalTextH := lineHeight * float64(len(lines))
			textY := cellY + math.Max((rowH-totalTextH)/2, 2)
			for _, line := range lines {
				pdf.SetX(cellX + 4)
				pdf.SetY(textY)
				_ = pdf.Cell(nil, line)
				textY += lineHeight
				if textY > cellY+rowH-2 {
					break
				}
			}
			setPDFTextFont(pdf, defaultFontSize, false, false)
			cellX += cw
		}
	}
}

func drawStyle(hasFill, hasStroke bool) string {
	if hasFill && hasStroke {
		return "DF"
	}
	if hasFill {
		return "F"
	}
	if hasStroke {
		return "D"
	}
	return ""
}

func pieAnglesFromAdjustments(adjustments []shapes.ShapeAdjustment) (float64, float64) {
	startDeg := 0.0
	endDeg := 360.0
	for _, adj := range adjustments {
		formula := strings.TrimSpace(adj.Formula)
		if !strings.HasPrefix(formula, "val ") {
			continue
		}
		numText := strings.TrimSpace(strings.TrimPrefix(formula, "val "))
		val, err := strconv.ParseFloat(numText, 64)
		if err != nil {
			continue
		}
		deg := val / 60000.0
		switch strings.TrimSpace(adj.Name) {
		case "adj1":
			startDeg = deg
		case "adj2":
			endDeg = deg
		}
	}
	return startDeg, endDeg
}

func drawPieShape(pdf *gopdf.GoPdf, s shapes.Shape, x, y, w, h float64, style string) {
	startDeg, endDeg := pieAnglesFromAdjustments(s.Adjustments)
	for endDeg < startDeg {
		endDeg += 360
	}
	if endDeg-startDeg > 360 {
		endDeg = startDeg + 360
	}
	cx := x + w/2
	cy := y + h/2
	rx := w / 2
	ry := h / 2

	points := make([]gopdf.Point, 0, 96)
	points = append(points, gopdf.Point{X: cx, Y: cy})
	step := 4.0
	for deg := startDeg; deg <= endDeg; deg += step {
		rad := deg * math.Pi / 180.0
		px := cx + rx*math.Cos(rad)
		py := cy + ry*math.Sin(rad)
		points = append(points, gopdf.Point{X: px, Y: py})
	}
	endRad := endDeg * math.Pi / 180.0
	points = append(points, gopdf.Point{
		X: cx + rx*math.Cos(endRad),
		Y: cy + ry*math.Sin(endRad),
	})
	pdf.Polygon(points, style)
}

func hexToRGB(color string) (uint8, uint8, uint8) {
	r, g, b, ok := resolveOOXMLColorToken(color)
	if !ok {
		return 0, 0, 0
	}
	return r, g, b
}

func stripHash(c string) string {
	if len(c) > 0 && c[0] == '#' {
		return c[1:]
	}
	return c
}

func rightArrowPoints(x, y, w, h float64) []gopdf.Point {
	aw := w * 0.5
	bw := w - aw
	hh := h * 0.5
	bh := h * 0.5
	return []gopdf.Point{
		{X: x, Y: y + (h-bh)/2},
		{X: x + bw, Y: y + (h-bh)/2},
		{X: x + bw, Y: y},
		{X: x + w, Y: y + hh},
		{X: x + bw, Y: y + h},
		{X: x + bw, Y: y + h - (h-bh)/2},
		{X: x, Y: y + h - (h-bh)/2},
	}
}

func leftArrowPoints(x, y, w, h float64) []gopdf.Point {
	aw := w * 0.5
	hh := h * 0.5
	bh := h * 0.5
	return []gopdf.Point{
		{X: x + aw, Y: y + (h-bh)/2},
		{X: x + w, Y: y + (h-bh)/2},
		{X: x + w, Y: y + h - (h-bh)/2},
		{X: x + aw, Y: y + h - (h-bh)/2},
		{X: x + aw, Y: y + h},
		{X: x, Y: y + hh},
		{X: x + aw, Y: y},
	}
}
