//nolint:mnd // Native chart rendering uses tuned numeric drawing constants.
package export

import (
	"math"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

// drawHorizontalBarItem draws a single bar in a horizontal bar chart at row i.
func drawHorizontalBarItem(
	pdf *gopdf.GoPdf, i int, values []float64,
	px, py, pw, ph, minV, maxV float64, opts chartSeriesOpts,
) {
	slot := ph / float64(len(values))
	bh := math.Max(8, slot*0.42)
	by := py + slot*float64(i) + (slot-bh)/2
	valueIndex := len(values) - 1 - i
	bx, bw := horizontalBarGeometry(values[valueIndex], minV, maxV, px+1, pw-2)
	if bw < 0.5 && values[valueIndex] != 0 {
		bw = 0.5
	}
	pdf.RectFromUpperLeftWithStyle(bx, by, bw, bh, "F")
	if opts.showDataLabels && bw > 0 {
		labelX := bx + bw + 3
		if values[valueIndex] < 0 {
			labelX = bx - 3
		}
		if dataLabelPosition(opts) == charts.DataLabelPositionInsideEnd {
			labelX = bx + bw - dataLabelOutsideGapPt
		}
		drawChartDataLabel(pdf, opts, chartDataLabelAt(opts, valueIndex, values[valueIndex]), labelX, by+bh/2-3)
	}
}

// chartSmoothSegments is how many straight segments each span of a smoothed
// line is drawn with. Eight is enough that the curve reads as one at slide
// sizes.
const chartSmoothSegments = 8

// chartSeriesLineWidthPt is how thick PowerPoint draws a line series: its
// default series stroke is 28575 EMU. The renderer used to leave the line at
// gopdf's 1pt default, which came out as a hairline beside PowerPoint's.
const chartSeriesLineWidthPt = 2.25

// chartValueY maps a data value onto the plot area.
//
// The full plot height is used: subtracting a few points from it, as this used
// to, squashed the series so its topmost point fell short of the axis value it
// was supposed to touch.
func chartValueY(value, minV, rangeV, plotY, plotH float64) float64 {
	if rangeV == 0 {
		return plotY + plotH
	}
	return plotY + plotH - ((value-minV)/rangeV)*plotH
}

// chartBarBandFraction is how much of a category slot the bars fill.
//
// PowerPoint's default gap width is 150% of the bar band, so the band is
// slot/(1+1.5) — measured bar for bar against its own export of a four-category
// column chart.
const chartBarBandFraction = 0.40

// categoryPointX is where a line or area chart draws the point for category i.
//
// PowerPoint gives each category a slot of the plot width and puts its point in
// the middle of it, the same slot a bar chart would centre its bar in — the
// category labels line up under the points because of it. Spreading the points
// from edge to edge instead stretched the series past its first and last
// categories.
func categoryPointX(px, pw float64, i, count int) float64 {
	if count <= 0 {
		return px
	}
	slot := pw / float64(count)
	return px + slot*(float64(i)+0.5)
}

// drawVerticalBarItem draws a single bar in a vertical bar chart at column i.
func drawVerticalBarItem(
	pdf *gopdf.GoPdf,
	i int,
	v, px, py, pw, ph, maxV, rangeV float64,
	nValues int,
	opts chartSeriesOpts,
) {
	slot := pw / float64(nValues)
	bw := math.Max(8, slot*chartBarBandFraction)
	bx := px + slot*float64(i) + (slot-bw)/2
	zeroY := py + ph*maxV/rangeV
	valueY := py + ph*(maxV-v)/rangeV
	barTop := math.Min(zeroY, valueY)
	barH := math.Abs(zeroY - valueY)
	if barH < 0.5 {
		barH = 0.5
	}
	pdf.RectFromUpperLeftWithStyle(bx, barTop, bw, barH, "F")
	if opts.showDataLabels {
		labelY := barLabelY(opts, barTop, barTop+barH, v < 0)
		// Clamp: if the label would fall above the plot area, draw it inside the bar top.
		if labelY < py {
			labelY = barTop + dataLabelInsideGapPt
		}
		drawChartDataLabel(pdf, opts, chartDataLabelAt(opts, i, v), bx+bw/2, labelY)
	}
}

// renderBarLike renders a vertical or horizontal bar chart.
// opts.color controls the bar fill (empty = default blue).
// opts.minValue / opts.maxValue override the computed axis range.
func renderBarLike(
	pdf *gopdf.GoPdf,
	title string,
	r chartRect,
	values []float64,
	categories []string,
	horizontal bool,
	opts chartSeriesOpts,
) {
	renderChartTitle(pdf, title, r)
	if len(values) == 0 {
		return
	}

	plotR := r
	if opts.showLegend {
		plotR = chartRectWithLegendMargin(pdf, r, opts.legendPosition, []string{opts.seriesName})
	}
	opts = withChartLabelData(opts, categories, values)
	minV, maxV := niceAxisRange(values)
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

	// The axis range has to be known before the plot rect, because the tick
	// labels it produces are what the plot area is sized around.
	valueAxis := chartAxisSpec{MinV: minV, MaxV: maxV, ValueFormat: opts.valueFormat}
	categoryAxis := chartAxisSpec{Categories: categories}
	verticalAxis, horizontalAxis := valueAxis, categoryAxis
	if horizontal {
		verticalAxis, horizontalAxis = categoryAxis, valueAxis
	}
	layout := solveChartLayout(pdf, plotR, opts.titleOverlay, title, verticalAxis, horizontalAxis)
	px, py, pw, ph := layout.X, layout.Y, layout.W, layout.H
	opts = withChartPlotArea(opts, px, py, pw, ph)

	barR, barG, barB := uint8(79), uint8(129), uint8(189)
	if opts.color != "" {
		barR, barG, barB = hexToRGB(opts.color)
	}

	if horizontal {
		drawHorizontalChartFrame(
			pdf, px, py, pw, ph, minV, maxV, categories,
			opts.showCatGridlines, opts.valueFormat, layout.Horizontal,
		)
	} else {
		drawChartFrame(
			pdf, px, py, pw, ph, minV, maxV,
			opts.showMajorGridlines, opts.valueFormat, layout.Vertical,
		)
	}

	for i, v := range values {
		// Re-set fill colour every iteration: gopdf shares the fill/text colour
		// register, so any pdf.SetTextColor call inside data-label helpers would
		// clobber the bar colour for the next iteration.
		pdf.SetFillColor(barR, barG, barB)
		if horizontal {
			drawHorizontalBarItem(pdf, i, values, px, py, pw, ph, minV, maxV, opts)
		} else {
			drawVerticalBarItem(pdf, i, v, px, py, pw, ph, maxV, rangeV, len(values), opts)
		}
	}

	if !horizontal {
		drawCategoryLabels(pdf, px, py, pw, ph, categories)
		if opts.catAxisTitle != "" {
			drawCategoryAxisTitle(pdf, px, py, pw, ph, opts.catAxisTitle)
		}
		if opts.valAxisTitle != "" {
			drawValueAxisTitle(pdf, px, py, pw, ph, opts.valAxisTitle)
		}
	}
	if opts.showLegend {
		drawChartLegend(pdf, r, opts.legendPosition, []legendEntry{
			{Name: opts.seriesName, R: barR, G: barG, B: barB},
		})
	}
}

// renderLineLike renders a line chart with optional markers and smooth curves.
func renderLineLike(
	pdf *gopdf.GoPdf,
	title string,
	r chartRect,
	values []float64,
	categories []string,
	markers bool,
	opts chartSeriesOpts,
) {
	renderChartTitle(pdf, title, r)
	if len(values) < 2 {
		return
	}

	plotR := r
	if opts.showLegend {
		plotR = chartRectWithLegendMargin(pdf, r, opts.legendPosition, []string{opts.seriesName})
	}
	opts = withChartLabelData(opts, categories, values)
	minV, maxV := niceAxisRange(values)
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
	opts = withChartPlotArea(opts, px, py, pw, ph)

	lineR, lineG, lineB := uint8(79), uint8(129), uint8(189)
	if opts.color != "" {
		lineR, lineG, lineB = hexToRGB(opts.color)
	}

	drawChartFrame(pdf, px, py, pw, ph, minV, maxV, opts.showMajorGridlines, opts.valueFormat, layout.Vertical)
	pdf.SetStrokeColor(lineR, lineG, lineB)
	pdf.SetFillColor(lineR, lineG, lineB)
	pdf.SetLineWidth(chartSeriesLineWidthPt)
	defer pdf.SetLineWidth(1)

	// Build the raw data points.
	rawPts := make([]gopdf.Point, len(values))
	for i, v := range values {
		rawPts[i] = gopdf.Point{
			X: categoryPointX(px, pw, i, len(values)),
			Y: chartValueY(v, minV, rangeV, py, ph),
		}
	}

	// Draw connecting lines: straight or Catmull-Rom smooth.
	drawPts := rawPts
	if opts.smooth && len(rawPts) >= 2 {
		drawPts = catmullRomPoints(rawPts, chartSmoothSegments)
	}
	for i := 1; i < len(drawPts); i++ {
		pdf.Line(drawPts[i-1].X, drawPts[i-1].Y, drawPts[i].X, drawPts[i].Y)
	}

	// Draw markers and data labels at original data points.
	for i, pt := range rawPts {
		if markers {
			drawFilledCircle(pdf, pt.X, pt.Y, 2.5, lineR, lineG, lineB)
		}
		if opts.showDataLabels {
			labelX, labelY := dataLabelPointAnchor(opts, pt.X, pt.Y)
			drawChartDataLabel(pdf, opts, chartDataLabelAt(opts, i, values[i]), labelX, labelY)
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
			{Name: opts.seriesName, R: lineR, G: lineG, B: lineB},
		})
	}
}

// renderAreaLike renders a filled area chart with a stroke outline.
func renderAreaLike(
	pdf *gopdf.GoPdf,
	title string,
	r chartRect,
	values []float64,
	categories []string,
	opts chartSeriesOpts,
) {
	renderChartTitle(pdf, title, r)
	if len(values) < 2 {
		return
	}

	plotR := r
	if opts.showLegend {
		plotR = chartRectWithLegendMargin(pdf, r, opts.legendPosition, []string{opts.seriesName})
	}
	opts = withChartLabelData(opts, categories, values)
	minV, maxV := niceAxisRange(values)
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
	opts = withChartPlotArea(opts, px, py, pw, ph)

	areaR, areaG, areaB := uint8(79), uint8(129), uint8(189)
	if opts.color != "" {
		areaR, areaG, areaB = hexToRGB(opts.color)
	}

	drawChartFrame(pdf, px, py, pw, ph, minV, maxV, opts.showMajorGridlines, opts.valueFormat, layout.Vertical)

	zeroY := py + ph*maxV/rangeV
	firstX := categoryPointX(px, pw, 0, len(values))
	lastX := categoryPointX(px, pw, len(values)-1, len(values))
	pts := make([]gopdf.Point, 0, len(values)+2)
	pts = append(pts, gopdf.Point{X: firstX, Y: zeroY})
	for i, v := range values {
		x := categoryPointX(px, pw, i, len(values))
		pts = append(pts, gopdf.Point{X: x, Y: chartValueY(v, minV, rangeV, py, ph)})
	}
	pts = append(pts, gopdf.Point{X: lastX, Y: zeroY})

	// Darken fill colour slightly for the stroke outline.
	strokeR := uint8(math.Max(0, float64(areaR)*0.7))
	strokeG := uint8(math.Max(0, float64(areaG)*0.7))
	strokeB := uint8(math.Max(0, float64(areaB)*0.7))
	pdf.SetStrokeColor(strokeR, strokeG, strokeB)
	pdf.SetFillColor(areaR, areaG, areaB)
	pdf.Polygon(pts, "FD")

	drawCategoryLabels(pdf, px, py, pw, ph, categories)
	if opts.catAxisTitle != "" {
		drawCategoryAxisTitle(pdf, px, py, pw, ph, opts.catAxisTitle)
	}
	if opts.valAxisTitle != "" {
		drawValueAxisTitle(pdf, px, py, pw, ph, opts.valAxisTitle)
	}
	if opts.showLegend {
		drawChartLegend(pdf, r, opts.legendPosition, []legendEntry{
			{Name: opts.seriesName, R: areaR, G: areaG, B: areaB},
		})
	}
}
