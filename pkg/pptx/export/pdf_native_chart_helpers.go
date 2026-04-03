//nolint:mnd // Chart helper math uses tuned visual constants for native PDF fidelity.
package export

import (
	"math"
	"strconv"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

// categoryLabel returns the label for index i from the provided slice,
// falling back to "Q<i+1>" when the slice is nil or too short.
func categoryLabel(categories []string, i int) string {
	if i < len(categories) && categories[i] != "" {
		return categories[i]
	}
	return "Q" + strconv.Itoa(i+1)
}

func pieColor(i int) (uint8, uint8, uint8) {
	palette := [][3]uint8{
		{79, 129, 189},
		{192, 80, 77},
		{155, 187, 89},
		{128, 100, 162},
		{75, 172, 198},
		{247, 150, 70},
	}
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
	if len(values) == 0 {
		return 0
	}
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
	if len(values) == 0 {
		return 0, 0
	}
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

func chartPlotRect(r chartRect) (float64, float64, float64, float64) {
	leftPad := math.Max(36, r.w*0.08)
	rightPad := math.Max(16, r.w*0.07)
	topPad := math.Max(24, r.h*0.12)
	bottomPad := math.Max(26, r.h*0.12)
	return r.x + leftPad, r.y + topPad, r.w - leftPad - rightPad, r.h - topPad - bottomPad
}

func chartPlotRectHorizontal(r chartRect) (float64, float64, float64, float64) {
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

func drawCategoryLabels(pdf *gopdf.GoPdf, x, y, w, h float64, categories []string) {
	count := len(categories)
	if count <= 0 {
		return
	}
	slot := w / float64(count)
	for i := range count {
		cx := x + slot*float64(i) + slot/2
		pdf.SetX(cx - 6)
		pdf.SetY(y + h + 8)
		_ = pdf.Cell(nil, categoryLabel(categories, i))
	}
}

func drawHorizontalChartFrame(pdf *gopdf.GoPdf, x, y, w, h, maxV float64, categories []string) {
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
	count := len(categories)
	if count <= 0 {
		return
	}
	slot := h / float64(count)
	for i := range count {
		// Match PowerPoint bar order for horizontal variants (last category at top).
		cy := y + slot*float64(i) + slot/2
		pdf.SetX(x - 18)
		pdf.SetY(cy - 3)
		_ = pdf.Cell(nil, categoryLabel(categories, count-1-i))
	}
}

// normalizePercentSeries converts a single-series 100%-stacked chart's values
// to their percentage of the total so the renderer can display proportional bars.
// Each value becomes (v / total) * 100; zero or negative totals yield an empty slice.
func normalizePercentSeries(values []float64) []float64 {
	if len(values) == 0 {
		return nil
	}
	total := sumFloat(values)
	if total <= 0 {
		return nil
	}
	out := make([]float64, len(values))
	for i, v := range values {
		out[i] = (v / total) * 100
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

func renderStockLike(pdf *gopdf.GoPdf, title string, r chartRect, openVals, highVals, lowVals, closeVals []float64) {
	renderChartTitle(pdf, title, r)
	n := min(len(highVals), min(len(lowVals), len(closeVals)))
	if n == 0 {
		return
	}
	px, py, pw, ph := r.x+8, r.y+24, r.w-16, r.h-32
	pdf.RectFromUpperLeftWithStyle(px, py, pw, ph, "D")
	all := append(append([]float64{}, highVals[:n]...), lowVals[:n]...)
	minV, maxV := minMax(all)
	if maxV <= minV {
		maxV = minV + 1
	}
	for i := range n {
		x := px + (float64(i)+0.5)*pw/float64(n)
		yHigh := py + ph - ((highVals[i]-minV)/(maxV-minV))*(ph-6) - 3
		yLow := py + ph - ((lowVals[i]-minV)/(maxV-minV))*(ph-6) - 3
		yClose := py + ph - ((closeVals[i]-minV)/(maxV-minV))*(ph-6) - 3
		pdf.Line(x, yHigh, x, yLow)
		pdf.Line(x, yClose, x+4, yClose)
		if i < len(openVals) {
			yOpen := py + ph - ((openVals[i]-minV)/(maxV-minV))*(ph-6) - 3
			pdf.Line(x-4, yOpen, x, yOpen)
		}
	}
}

func renderComboLike(
	pdf *gopdf.GoPdf,
	title string,
	r chartRect,
	barSeries, lineSeries []charts.Series,
	categories []string,
) {
	comboRect := r
	if comboRect.w > 160 {
		comboRect.w -= 120
	}
	values := []float64{1, 2, 3}
	if len(barSeries) > 0 && len(barSeries[0].Values) > 0 {
		values = barSeries[0].Values
	}
	renderBarLike(pdf, title, comboRect, values, categories, false)
	if len(lineSeries) > 0 && len(lineSeries[0].Values) > 1 {
		px, py, pw, ph := chartPlotRect(comboRect)
		scaleMax := niceAxisMax(max(maxFloat(values), maxFloat(lineSeries[0].Values)))
		if scaleMax <= 0 {
			scaleMax = 1
		}
		pdf.SetStrokeColor(192, 80, 77)
		pdf.SetFillColor(192, 80, 77)
		for i := 1; i < len(lineSeries[0].Values); i++ {
			x1 := px + (float64(i-1)*pw)/float64(len(lineSeries[0].Values)-1)
			y1 := py + ph - (lineSeries[0].Values[i-1]/scaleMax)*(ph-4)
			x2 := px + (float64(i)*pw)/float64(len(lineSeries[0].Values)-1)
			y2 := py + ph - (lineSeries[0].Values[i]/scaleMax)*(ph-4)
			pdf.Line(x1, y1, x2, y2)
			pdf.RectFromUpperLeftWithStyle(x2-2, y2-2, 4, 4, "F")
		}
	}
	barLabel := ""
	if len(barSeries) > 0 {
		barLabel = barSeries[0].Name
	}
	lineLabel := ""
	if len(lineSeries) > 0 {
		lineLabel = lineSeries[0].Name
	}
	drawComboLegend(pdf, r, barLabel, lineLabel)
}

func drawComboLegend(pdf *gopdf.GoPdf, r chartRect, barLabel, lineLabel string) {
	if barLabel == "" {
		barLabel = "Series 1"
	}
	if lineLabel == "" {
		lineLabel = "Series 2"
	}
	x := r.x + r.w - 86
	y := r.y + r.h*0.52
	pdf.SetFillColor(79, 129, 189)
	pdf.RectFromUpperLeftWithStyle(x, y, 12, 6, "F")
	pdf.SetX(x + 15)
	pdf.SetY(y - 1)
	_ = pdf.Cell(nil, barLabel)

	pdf.SetStrokeColor(192, 80, 77)
	pdf.SetFillColor(192, 80, 77)
	pdf.Line(x, y+16, x+12, y+16)
	pdf.RectFromUpperLeftWithStyle(x+6, y+14, 3, 3, "F")
	pdf.SetX(x + 15)
	pdf.SetY(y + 12)
	_ = pdf.Cell(nil, lineLabel)
}
