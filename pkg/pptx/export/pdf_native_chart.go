//nolint:mnd // Chart rendering uses fixed visual defaults for native PDF output.
package export

import (
	"math"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

type chartRect struct{ x, y, w, h float64 }

func renderNativePDFSlideCharts(pdf *gopdf.GoPdf, slide elements.SlideContent) {
	if slide.Chart != nil {
		renderBarLike(pdf, slide.Chart.Title, chartRectFromLength(slide.Chart.X.Emu(), slide.Chart.Y.Emu(), slide.Chart.CX.Emu(), slide.Chart.CY.Emu()), slide.Chart.Values, false)
	}
	if slide.BarHorizontal != nil {
		renderBarLike(pdf, slide.BarHorizontal.Title, chartRectFromLength(slide.BarHorizontal.X.Emu(), slide.BarHorizontal.Y.Emu(), slide.BarHorizontal.CX.Emu(), slide.BarHorizontal.CY.Emu()), slide.BarHorizontal.Values, true)
	}
	if slide.BarStacked != nil {
		renderBarLike(pdf, slide.BarStacked.Title, chartRectFromLength(slide.BarStacked.X.Emu(), slide.BarStacked.Y.Emu(), slide.BarStacked.CX.Emu(), slide.BarStacked.CY.Emu()), slide.BarStacked.Values, true)
	}
	if slide.BarStacked100 != nil {
		renderBarLike(
			pdf,
			slide.BarStacked100.Title,
			chartRectFromLength(slide.BarStacked100.X.Emu(), slide.BarStacked100.Y.Emu(), slide.BarStacked100.CX.Emu(), slide.BarStacked100.CY.Emu()),
			normalizePercentSeries(slide.BarStacked100.Values),
			true,
		)
	}
	if slide.Line != nil {
		renderLineLike(pdf, slide.Line.Title, chartRectFromLength(slide.Line.X.Emu(), slide.Line.Y.Emu(), slide.Line.CX.Emu(), slide.Line.CY.Emu()), slide.Line.Values, false)
	}
	if slide.LineMarkers != nil {
		renderLineLike(pdf, slide.LineMarkers.Title, chartRectFromLength(slide.LineMarkers.X.Emu(), slide.LineMarkers.Y.Emu(), slide.LineMarkers.CX.Emu(), slide.LineMarkers.CY.Emu()), slide.LineMarkers.Values, true)
	}
	if slide.LineStacked != nil {
		renderLineLike(pdf, slide.LineStacked.Title, chartRectFromLength(slide.LineStacked.X.Emu(), slide.LineStacked.Y.Emu(), slide.LineStacked.CX.Emu(), slide.LineStacked.CY.Emu()), slide.LineStacked.Values, true)
	}
	if slide.Area != nil {
		renderAreaLike(pdf, slide.Area.Title, chartRectFromLength(slide.Area.X.Emu(), slide.Area.Y.Emu(), slide.Area.CX.Emu(), slide.Area.CY.Emu()), slide.Area.Values)
	}
	if slide.AreaStacked != nil {
		renderAreaLike(pdf, slide.AreaStacked.Title, chartRectFromLength(slide.AreaStacked.X.Emu(), slide.AreaStacked.Y.Emu(), slide.AreaStacked.CX.Emu(), slide.AreaStacked.CY.Emu()), slide.AreaStacked.Values)
	}
	if slide.AreaStacked100 != nil {
		renderAreaLike(
			pdf,
			slide.AreaStacked100.Title,
			chartRectFromLength(slide.AreaStacked100.X.Emu(), slide.AreaStacked100.Y.Emu(), slide.AreaStacked100.CX.Emu(), slide.AreaStacked100.CY.Emu()),
			normalizePercentSeries(slide.AreaStacked100.Values),
		)
	}
	if slide.Pie != nil {
		renderPieLike(pdf, slide.Pie.Title, chartRectFromLength(slide.Pie.X.Emu(), slide.Pie.Y.Emu(), slide.Pie.CX.Emu(), slide.Pie.CY.Emu()), slide.Pie.Values, false)
	}
	if slide.Doughnut != nil {
		renderPieLike(pdf, slide.Doughnut.Title, chartRectFromLength(slide.Doughnut.X.Emu(), slide.Doughnut.Y.Emu(), slide.Doughnut.CX.Emu(), slide.Doughnut.CY.Emu()), slide.Doughnut.Values, true)
	}
	if slide.Scatter != nil {
		renderScatterLike(pdf, slide.Scatter.Title, chartRectFromLength(slide.Scatter.X.Emu(), slide.Scatter.Y.Emu(), slide.Scatter.CX.Emu(), slide.Scatter.CY.Emu()), slide.Scatter.XValues, slide.Scatter.YValues, nil)
	}
	if slide.Bubble != nil {
		renderScatterLike(pdf, slide.Bubble.Title, chartRectFromLength(slide.Bubble.X, slide.Bubble.Y, slide.Bubble.CX, slide.Bubble.CY), slide.Bubble.XValues, slide.Bubble.YValues, slide.Bubble.BubbleSizes)
	}
	if slide.Radar != nil {
		renderRadarLike(pdf, slide.Radar.Title, chartRectFromLength(slide.Radar.X.Emu(), slide.Radar.Y.Emu(), slide.Radar.CX.Emu(), slide.Radar.CY.Emu()), slide.Radar.Values, false)
	}
	if slide.RadarFilled != nil {
		renderRadarLike(pdf, slide.RadarFilled.Title, chartRectFromLength(slide.RadarFilled.X.Emu(), slide.RadarFilled.Y.Emu(), slide.RadarFilled.CX.Emu(), slide.RadarFilled.CY.Emu()), slide.RadarFilled.Values, true)
	}
	if slide.StockHLC != nil {
		renderStockLike(pdf, slide.StockHLC.Title, chartRectFromLength(slide.StockHLC.X.Emu(), slide.StockHLC.Y.Emu(), slide.StockHLC.CX.Emu(), slide.StockHLC.CY.Emu()), nil, slide.StockHLC.HighValues, slide.StockHLC.LowValues, slide.StockHLC.CloseValues)
	}
	if slide.StockOHLC != nil {
		renderStockLike(pdf, slide.StockOHLC.Title, chartRectFromLength(slide.StockOHLC.X.Emu(), slide.StockOHLC.Y.Emu(), slide.StockOHLC.CX.Emu(), slide.StockOHLC.CY.Emu()), slide.StockOHLC.OpenValues, slide.StockOHLC.HighValues, slide.StockOHLC.LowValues, slide.StockOHLC.CloseValues)
	}
	if slide.Combo != nil {
		renderComboLike(pdf, slide.Combo.Title, chartRectFromLength(slide.Combo.X.Emu(), slide.Combo.Y.Emu(), slide.Combo.CX.Emu(), slide.Combo.CY.Emu()), slide.Combo.BarSeries, slide.Combo.LineSeries)
	}
}

func chartRectFromLength(x, y, w, h int64) chartRect {
	return chartRect{emuToPt(x), emuToPt(y), emuToPt(w), emuToPt(h)}
}

func renderChartTitle(pdf *gopdf.GoPdf, title string, r chartRect) {
	if title == "" {
		return
	}
	pdf.SetTextColor(40, 40, 40)
	pdf.SetX(r.x + 6)
	pdf.SetY(r.y + 4)
	_ = pdf.Cell(nil, title)
}

func renderBarLike(pdf *gopdf.GoPdf, title string, r chartRect, values []float64, horizontal bool) {
	renderChartTitle(pdf, title, r)
	if len(values) == 0 {
		return
	}
	px, py, pw, ph := chartPlotRect(r)
	if horizontal {
		px, py, pw, ph = chartPlotRectHorizontal(r)
	}
	maxV := niceAxisMax(maxFloat(values))
	if maxV <= 0 {
		maxV = 1
	}
	if horizontal {
		drawHorizontalChartFrame(pdf, px, py, pw, ph, maxV, len(values))
	} else {
		drawChartFrame(pdf, px, py, pw, ph, maxV)
	}
	pdf.SetFillColor(79, 129, 189)
	for i, v := range values {
		if horizontal {
			slot := ph / float64(len(values))
			bh := math.Max(8, slot*0.42)
			by := py + slot*float64(i) + (slot-bh)/2
			valueIndex := len(values) - 1 - i
			if valueIndex < 0 || valueIndex >= len(values) {
				valueIndex = i
			}
			bw := (values[valueIndex] / maxV) * (pw - 2)
			pdf.RectFromUpperLeftWithStyle(px+1, by, bw, bh, "F")
		} else {
			slot := pw / float64(len(values))
			bw := math.Max(8, slot*0.40)
			bx := px + slot*float64(i) + (slot-bw)/2
			bh := (v / maxV) * (ph - 4)
			pdf.RectFromUpperLeftWithStyle(bx, py+ph-bh, bw, bh, "F")
		}
	}
	if !horizontal {
		drawCategoryLabels(pdf, px, py, pw, ph, len(values))
	}
}

func renderLineLike(pdf *gopdf.GoPdf, title string, r chartRect, values []float64, markers bool) {
	renderChartTitle(pdf, title, r)
	if len(values) < 2 {
		return
	}
	px, py, pw, ph := chartPlotRect(r)
	maxV := niceAxisMax(maxFloat(values))
	if maxV <= 0 {
		maxV = 1
	}
	drawChartFrame(pdf, px, py, pw, ph, maxV)
	pdf.SetStrokeColor(192, 80, 77)
	for i := 1; i < len(values); i++ {
		x1 := px + (float64(i-1)*pw)/float64(len(values)-1)
		y1 := py + ph - (values[i-1]/maxV)*(ph-4)
		x2 := px + (float64(i)*pw)/float64(len(values)-1)
		y2 := py + ph - (values[i]/maxV)*(ph-4)
		pdf.Line(x1, y1, x2, y2)
		if markers {
			pdf.Oval(x2-2, y2-2, x2+2, y2+2)
		}
	}
	drawCategoryLabels(pdf, px, py, pw, ph, len(values))
}

func renderAreaLike(pdf *gopdf.GoPdf, title string, r chartRect, values []float64) {
	renderChartTitle(pdf, title, r)
	if len(values) < 2 {
		return
	}
	px, py, pw, ph := chartPlotRect(r)
	maxV := niceAxisMax(maxFloat(values))
	if maxV <= 0 {
		maxV = 1
	}
	drawChartFrame(pdf, px, py, pw, ph, maxV)
	pts := make([]gopdf.Point, 0, len(values)+2)
	pts = append(pts, gopdf.Point{X: px, Y: py + ph})
	for i, v := range values {
		x := px + (float64(i)*pw)/float64(len(values)-1)
		y := py + ph - (v/maxV)*(ph-4)
		pts = append(pts, gopdf.Point{X: x, Y: y})
	}
	pts = append(pts, gopdf.Point{X: px + pw, Y: py + ph})
	pdf.SetFillColor(155, 187, 89)
	pdf.Polygon(pts, "F")
	drawCategoryLabels(pdf, px, py, pw, ph, len(values))
}

func renderPieLike(pdf *gopdf.GoPdf, title string, r chartRect, values []float64, doughnut bool) {
	renderChartTitle(pdf, title, r)
	if len(values) == 0 {
		return
	}
	cx := r.x + r.w/2
	cy := r.y + r.h/2 + 8
	radius := math.Min(r.w, r.h) * 0.35
	total := sumFloat(values)
	if total <= 0 {
		return
	}
	start := -math.Pi / 2
	for i, v := range values {
		frac := v / total
		end := start + frac*2*math.Pi
		rC, gC, bC := pieColor(i)
		drawWedge(pdf, cx, cy, radius, start, end, rC, gC, bC)
		start = end
	}
	if doughnut {
		pdf.SetFillColor(255, 255, 255)
		inner := radius * 0.45
		pdf.Oval(cx-inner, cy-inner, cx+inner, cy+inner)
	}
}

func renderScatterLike(pdf *gopdf.GoPdf, title string, r chartRect, xs, ys, sizes []float64) {
	renderChartTitle(pdf, title, r)
	if len(xs) == 0 || len(ys) == 0 {
		return
	}
	n := min(len(xs), len(ys))
	px, py, pw, ph := r.x+8, r.y+24, r.w-16, r.h-32
	minX, maxX := minMax(xs[:n])
	minY, maxY := minMax(ys[:n])
	if maxX <= minX {
		maxX = minX + 1
	}
	if maxY <= minY {
		maxY = minY + 1
	}
	pdf.RectFromUpperLeftWithStyle(px, py, pw, ph, "D")
	for i := 0; i < n; i++ {
		xf := (xs[i] - minX) / (maxX - minX)
		yf := (ys[i] - minY) / (maxY - minY)
		x := px + xf*(pw-6) + 3
		y := py + ph - yf*(ph-6) - 3
		rad := 3.0
		if i < len(sizes) && sizes[i] > 0 {
			rad = math.Min(8, 2+math.Sqrt(sizes[i])*0.8)
		}
		pdf.Oval(x-rad, y-rad, x+rad, y+rad)
	}
}

func renderRadarLike(pdf *gopdf.GoPdf, title string, r chartRect, values []float64, filled bool) {
	renderChartTitle(pdf, title, r)
	if len(values) < 3 {
		return
	}
	cx := r.x + r.w/2
	cy := r.y + r.h/2 + 6
	radius := math.Min(r.w, r.h) * 0.35
	maxV := maxFloat(values)
	if maxV <= 0 {
		maxV = 1
	}
	pts := make([]gopdf.Point, 0, len(values))
	for i, v := range values {
		angle := -math.Pi/2 + (2*math.Pi*float64(i))/float64(len(values))
		scale := v / maxV
		pts = append(pts, gopdf.Point{X: cx + math.Cos(angle)*radius*scale, Y: cy + math.Sin(angle)*radius*scale})
	}
	if filled {
		pdf.SetFillColor(155, 187, 89)
		pdf.Polygon(pts, "FD")
	} else {
		pdf.Polygon(pts, "D")
	}
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
	for i := 0; i < n; i++ {
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

func renderComboLike(pdf *gopdf.GoPdf, title string, r chartRect, barSeries, lineSeries []charts.Series) {
	comboRect := r
	if comboRect.w > 160 {
		comboRect.w -= 120
	}
	values := []float64{1, 2, 3}
	if len(barSeries) > 0 && len(barSeries[0].Values) > 0 {
		values = barSeries[0].Values
	}
	renderBarLike(pdf, title, comboRect, values, false)
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
	drawComboLegend(pdf, r)
}
