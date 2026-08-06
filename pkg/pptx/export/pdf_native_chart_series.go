//nolint:mnd // Chart series rendering uses tuned visual constants for native PDF fidelity.
package export

import (
	"math"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

// renderStockLike renders a High/Low/Close or Open/High/Low/Close stock chart.
// Pass nil for openVals to get an HLC chart.
func renderStockLike(
	pdf *gopdf.GoPdf,
	title string,
	r chartRect,
	openVals, highVals, lowVals, closeVals []float64,
	categories []string,
	opts chartSeriesOpts,
) {
	renderChartTitle(pdf, title, r)
	n := min(len(highVals), min(len(lowVals), len(closeVals)))
	if n == 0 {
		return
	}
	// Determine axis range from all OHLC values.
	// Use niceAxisRange so the axis starts at 0 for non-negative data (matches PowerPoint default).
	all := append(append([]float64{}, highVals[:n]...), lowVals[:n]...)
	minV, maxV := niceAxisRange(all)
	if opts.minValue != nil {
		minV = *opts.minValue
	}
	if opts.maxValue != nil {
		maxV = *opts.maxValue
	}
	if maxV <= minV {
		maxV = minV + 1
	}
	rangeV := maxV - minV

	// The tick labels the axis range produces are what the plot area is sized
	// around, so the range has to be resolved first.
	layout := solveChartLayout(
		pdf, r, opts.titleOverlay,
		title,
		chartAxisSpec{MinV: minV, MaxV: maxV, ValueFormat: opts.valueFormat},
		chartAxisSpec{Categories: categories},
	)
	px, py, pw, ph := layout.X, layout.Y, layout.W, layout.H

	drawChartFrame(pdf, px, py, pw, ph, minV, maxV, opts.showMajorGridlines, opts.valueFormat, layout.Vertical)
	// Complete the box with top and right borders.
	pdf.SetStrokeColor(30, 30, 30)
	pdf.Line(px, py, px+pw, py)
	pdf.Line(px+pw, py, px+pw, py+ph)

	for i := range n {
		x := px + (float64(i)+0.5)*pw/float64(n)
		yHigh := py + ph - ((highVals[i]-minV)/rangeV)*(ph-2)
		yLow := py + ph - ((lowVals[i]-minV)/rangeV)*(ph-2)
		yClose := py + ph - ((closeVals[i]-minV)/rangeV)*(ph-2)
		pdf.SetStrokeColor(30, 30, 30)
		pdf.Line(x, yHigh, x, yLow)
		pdf.Line(x, yClose, x+5, yClose)
		if i < len(openVals) {
			yOpen := py + ph - ((openVals[i]-minV)/rangeV)*(ph-2)
			pdf.Line(x-5, yOpen, x, yOpen)
		}
	}

	drawCategoryLabels(pdf, px, py, pw, ph, categories)
	if opts.catAxisTitle != "" {
		drawCategoryAxisTitle(pdf, px, py, pw, ph, opts.catAxisTitle)
	}
	if opts.valAxisTitle != "" {
		drawValueAxisTitle(pdf, px, py, pw, ph, opts.valAxisTitle)
	}
	if opts.showLegend {
		drawChartLegend(pdf, r, opts.legendPosition, []legendEntry{
			{Name: "High", R: 30, G: 30, B: 30},
			{Name: "Low", R: 90, G: 90, B: 90},
			{Name: "Close", R: 60, G: 60, B: 60},
		})
	}
}

// drawComboBarSeries draws all bar series for a combo chart as grouped bars.
func drawComboBarSeries(
	pdf *gopdf.GoPdf, barSeries []charts.Series,
	px, py, ph, catSlot float64, nCats, nBarSeries int,
	maxV, rangeV, zeroY float64,
) {
	for si, bs := range barSeries {
		colR, colG, colB := pieColor(si)
		pdf.SetFillColor(colR, colG, colB)
		for ci, v := range bs.Values {
			if ci >= nCats {
				break
			}
			// All the bars of one category share the band PowerPoint's default
			// gap width leaves them, centred in the category slot. Giving each
			// series 85% of a full sub-slot instead drew a bar chart with almost
			// no gaps between categories.
			band := catSlot * chartBarBandFraction
			bw := math.Max(4, band/float64(max(1, nBarSeries)))
			bx := px + catSlot*float64(ci) + (catSlot-band)/2 + bw*float64(si)
			valueY := py + ph*(maxV-v)/rangeV
			barTop := math.Min(zeroY, valueY)
			barH := math.Abs(zeroY - valueY)
			if barH < 0.5 {
				barH = 0.5
			}
			pdf.RectFromUpperLeftWithStyle(bx, barTop, bw, barH, "F")
		}
	}
}

// drawComboLineSeries draws all line series overlaid on a combo chart.
func drawComboLineSeries(
	pdf *gopdf.GoPdf, lineSeries []charts.Series,
	px, py, ph, catSlot float64, nCats, nBarSeries int, minV, rangeV float64,
) {
	for si, ls := range lineSeries {
		colR, colG, colB := pieColor(si + nBarSeries)
		pdf.SetStrokeColor(colR, colG, colB)
		n := min(len(ls.Values), nCats)
		if n < 1 {
			continue
		}
		points := make([]gopdf.Point, 0, n)
		for i := range n {
			points = append(points, gopdf.Point{
				X: px + catSlot*(float64(i)+0.5),
				Y: chartValueY(ls.Values[i], minV, rangeV, py, ph),
			})
		}
		// OOXML's <c:smooth> defaults to on, and a combo chart's line series
		// almost never states it, so PowerPoint curves the line where this used
		// to join the points with straight segments.
		drawPoints := points
		if len(points) >= 2 {
			drawPoints = catmullRomPoints(points, chartSmoothSegments)
		}
		pdf.SetLineWidth(chartSeriesLineWidthPt)
		for i := 1; i < len(drawPoints); i++ {
			pdf.Line(drawPoints[i-1].X, drawPoints[i-1].Y, drawPoints[i].X, drawPoints[i].Y)
		}
		pdf.SetLineWidth(1)
		for _, point := range points {
			drawFilledCircle(pdf, point.X, point.Y, 2.5, colR, colG, colB)
		}
	}
}

// comboSeriesNames is every legend entry a combo chart will draw, which is what
// its legend has to be wide enough for.
func comboSeriesNames(barSeries, lineSeries []charts.Series) []string {
	names := make([]string, 0, len(barSeries)+len(lineSeries))
	for _, series := range barSeries {
		names = append(names, series.Name)
	}
	for _, series := range lineSeries {
		names = append(names, series.Name)
	}
	return names
}

// renderComboLike renders a combination bar+line chart with multiple series per type.
func renderComboLike(
	pdf *gopdf.GoPdf,
	title string,
	r chartRect,
	barSeries, lineSeries []charts.Series,
	categories []string,
	opts chartSeriesOpts,
) {
	renderChartTitle(pdf, title, r)

	plotR := r
	if opts.showLegend {
		plotR = chartRectWithLegendMargin(pdf, r, opts.legendPosition, comboSeriesNames(barSeries, lineSeries))
	}
	// Unified axis range across all series.
	var allVals []float64
	for _, s := range barSeries {
		allVals = append(allVals, s.Values...)
	}
	for _, s := range lineSeries {
		allVals = append(allVals, s.Values...)
	}
	minV, maxV := niceAxisRange(allVals)
	if opts.minValue != nil {
		minV = *opts.minValue
	}
	if opts.maxValue != nil {
		maxV = *opts.maxValue
	}
	if maxV <= minV {
		maxV = minV + 1
	}
	rangeV := maxV - minV

	// The tick labels the axis range produces are what the plot area is sized
	// around, so the range has to be resolved first.
	layout := solveChartLayout(
		pdf, plotR, opts.titleOverlay,
		title,
		chartAxisSpec{MinV: minV, MaxV: maxV, ValueFormat: opts.valueFormat},
		chartAxisSpec{Categories: categories},
	)
	px, py, pw, ph := layout.X, layout.Y, layout.W, layout.H

	drawChartFrame(pdf, px, py, pw, ph, minV, maxV, opts.showMajorGridlines, opts.valueFormat, layout.Vertical)

	nCats := len(categories)
	if nCats == 0 {
		if len(barSeries) > 0 {
			nCats = len(barSeries[0].Values)
		} else if len(lineSeries) > 0 {
			nCats = len(lineSeries[0].Values)
		}
	}
	if nCats == 0 {
		return
	}

	catSlot := pw / float64(nCats)
	nBarSeries := len(barSeries)
	zeroY := py + ph*maxV/rangeV

	drawComboBarSeries(pdf, barSeries, px, py, ph, catSlot, nCats, nBarSeries, maxV, rangeV, zeroY)
	drawComboLineSeries(pdf, lineSeries, px, py, ph, catSlot, nCats, nBarSeries, minV, rangeV)

	drawCategoryLabels(pdf, px, py, pw, ph, categories)
	if opts.catAxisTitle != "" {
		drawCategoryAxisTitle(pdf, px, py, pw, ph, opts.catAxisTitle)
	}
	if opts.valAxisTitle != "" {
		drawValueAxisTitle(pdf, px, py, pw, ph, opts.valAxisTitle)
	}

	if opts.showLegend {
		entries := make([]legendEntry, 0, len(barSeries)+len(lineSeries))
		for si, bs := range barSeries {
			eR, eG, eB := pieColor(si)
			entries = append(entries, legendEntry{Name: bs.Name, R: eR, G: eG, B: eB})
		}
		for si, ls := range lineSeries {
			eR, eG, eB := pieColor(si + nBarSeries)
			entries = append(entries, legendEntry{Name: ls.Name, R: eR, G: eG, B: eB})
		}
		drawChartLegend(pdf, r, opts.legendPosition, entries)
	}
}
