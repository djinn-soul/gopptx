package export

import (
	"math"
	"strconv"

	"github.com/signintech/gopdf"
)

func pieColor(i int) (r, g, b uint8) {
	palette := [][3]uint8{{79, 129, 189}, {192, 80, 77}, {155, 187, 89}, {128, 100, 162}, {75, 172, 198}, {247, 150, 70}}
	c := palette[i%len(palette)]
	return c[0], c[1], c[2]
}

func drawWedge(pdf *gopdf.GoPdf, cx, cy, radius, start, end float64, r, g, b uint8) {
	pdf.SetFillColor(r, g, b)
	pts := []gopdf.Point{{X: cx, Y: cy}}
	steps := max(8, int(math.Ceil((end-start)/(math.Pi/18))))
	for i := 0; i <= steps; i++ {
		t := start + (end-start)*float64(i)/float64(steps)
		pts = append(pts, gopdf.Point{X: cx + radius*math.Cos(t), Y: cy + radius*math.Sin(t)})
	}
	pdf.Polygon(pts, "F")
}

func maxFloat(values []float64) float64 {
	v := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] > v {
			v = values[i]
		}
	}
	return v
}

func sumFloat(values []float64) float64 {
	total := 0.0
	for _, v := range values {
		total += v
	}
	return total
}

func minMax(values []float64) (float64, float64) {
	minV, maxV := values[0], values[0]
	for i := 1; i < len(values); i++ {
		if values[i] < minV {
			minV = values[i]
		}
		if values[i] > maxV {
			maxV = values[i]
		}
	}
	return minV, maxV
}

func chartPlotRect(r chartRect) (x, y, w, h float64) {
	leftPad := math.Max(36, r.w*0.08)
	rightPad := math.Max(16, r.w*0.07)
	topPad := math.Max(24, r.h*0.12)
	bottomPad := math.Max(26, r.h*0.12)
	return r.x + leftPad, r.y + topPad, r.w - leftPad - rightPad, r.h - topPad - bottomPad
}

func chartPlotRectHorizontal(r chartRect) (x, y, w, h float64) {
	leftPad := math.Max(18, r.w*0.03)
	rightPad := math.Max(26, r.w*0.10)
	topPad := math.Max(24, r.h*0.12)
	bottomPad := math.Max(26, r.h*0.12)
	return r.x + leftPad, r.y + topPad, r.w - leftPad - rightPad, r.h - topPad - bottomPad
}

func drawChartFrame(pdf *gopdf.GoPdf, x, y, w, h, maxV float64) {
	pdf.SetStrokeColor(30, 30, 30)
	pdf.Line(x, y, x, y+h)
	pdf.Line(x, y+h, x+w, y+h)
	pdf.SetStrokeColor(90, 90, 90)
	for i := 1; i <= 5; i++ {
		yLine := y + h - (float64(i)/5.0)*h
		pdf.Line(x, yLine, x+w, yLine)
		pdf.SetX(x - 18)
		pdf.SetY(yLine - 3)
		_ = pdf.Cell(nil, strconv.Itoa(int(math.Round(maxV*float64(i)/5.0))))
	}
	pdf.SetX(x - 8)
	pdf.SetY(y + h - 3)
	_ = pdf.Cell(nil, "0")
}

func drawCategoryLabels(pdf *gopdf.GoPdf, x, y, w, h float64, count int) {
	if count <= 0 {
		return
	}
	slot := w / float64(count)
	for i := 0; i < count; i++ {
		cx := x + slot*float64(i) + slot/2
		label := "Q" + strconv.Itoa(i+1)
		pdf.SetX(cx - 6)
		pdf.SetY(y + h + 8)
		_ = pdf.Cell(nil, label)
	}
}

func drawHorizontalChartFrame(pdf *gopdf.GoPdf, x, y, w, h, maxV float64, count int) {
	pdf.SetStrokeColor(30, 30, 30)
	pdf.Line(x, y, x, y+h)
	pdf.Line(x, y+h, x+w, y+h)
	pdf.SetStrokeColor(90, 90, 90)
	steps := 6
	if math.Abs(maxV-100.0) < 0.001 {
		steps = 10
	}
	for i := 1; i <= steps; i++ {
		xLine := x + (float64(i)/float64(steps))*w
		pdf.Line(xLine, y, xLine, y+h)
	}
	for i := 0; i <= steps; i++ {
		xTick := x + (float64(i)/float64(steps))*w
		tickValue := int(math.Round(maxV * float64(i) / float64(steps)))
		label := strconv.Itoa(tickValue)
		if steps == 10 && math.Abs(maxV-100.0) < 0.001 {
			label += "%"
		}
		pdf.SetX(xTick - 4)
		pdf.SetY(y + h + 8)
		_ = pdf.Cell(nil, label)
	}
	if count <= 0 {
		return
	}
	slot := h / float64(count)
	for i := 0; i < count; i++ {
		// Match PowerPoint bar order for horizontal variants (last category at top).
		label := "Q" + strconv.Itoa(count-i)
		cy := y + slot*float64(i) + slot/2
		pdf.SetX(x - 18)
		pdf.SetY(cy - 3)
		_ = pdf.Cell(nil, label)
	}
}

func normalizePercentSeries(values []float64) []float64 {
	if len(values) == 0 {
		return []float64{100, 100, 100, 100}
	}
	out := make([]float64, len(values))
	for i := range values {
		out[i] = 100
	}
	return out
}

func niceAxisMax(value float64) float64 {
	if value <= 0 {
		return 1
	}
	if value <= 10 {
		return math.Ceil(value)
	}
	if value <= 100 {
		return math.Ceil(value/5.0) * 5.0
	}
	return math.Ceil(value/50.0) * 50.0
}

func drawComboLegend(pdf *gopdf.GoPdf, r chartRect) {
	x := r.x + r.w - 86
	y := r.y + r.h*0.52
	pdf.SetFillColor(79, 129, 189)
	pdf.RectFromUpperLeftWithStyle(x, y, 12, 6, "F")
	pdf.SetX(x + 15)
	pdf.SetY(y - 1)
	_ = pdf.Cell(nil, "Revenue")

	pdf.SetStrokeColor(192, 80, 77)
	pdf.SetFillColor(192, 80, 77)
	pdf.Line(x, y+16, x+12, y+16)
	pdf.RectFromUpperLeftWithStyle(x+6, y+14, 3, 3, "F")
	pdf.SetX(x + 15)
	pdf.SetY(y + 12)
	_ = pdf.Cell(nil, "Growth %")
}
