//nolint:mnd // Native chart rendering uses tuned numeric drawing constants.
package export

import (
	"math"

	"github.com/signintech/gopdf"
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
		drawBarDataLabel(pdf, labelX, by+bh/2-3, values[valueIndex])
	}
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
	bw := math.Max(8, slot*0.40)
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
		labelY := barTop - 5
		if v < 0 {
			labelY = barTop + barH + 3
		}
		// Clamp: if the label would fall above the plot area, draw it inside the bar top.
		if labelY < py {
			labelY = barTop + 2
		}
		drawBarDataLabel(pdf, bx+bw/2, labelY, v)
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
		plotR = chartRectWithLegendMargin(r, opts.legendPosition)
	}
	px, py, pw, ph := chartPlotRect(plotR, opts.titleOverlay)
	if horizontal {
		px, py, pw, ph = chartPlotRectHorizontal(plotR, opts.titleOverlay)
	}

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

	barR, barG, barB := uint8(79), uint8(129), uint8(189)
	if opts.color != "" {
		barR, barG, barB = hexToRGB(opts.color)
	}

	if horizontal {
		drawHorizontalChartFrame(pdf, px, py, pw, ph, minV, maxV, categories, opts.showCatGridlines, opts.valueFormat)
	} else {
		drawChartFrame(pdf, px, py, pw, ph, minV, maxV, opts.showMajorGridlines, opts.valueFormat)
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
		plotR = chartRectWithLegendMargin(r, opts.legendPosition)
	}
	px, py, pw, ph := chartPlotRect(plotR, opts.titleOverlay)

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

	lineR, lineG, lineB := uint8(79), uint8(129), uint8(189)
	if opts.color != "" {
		lineR, lineG, lineB = hexToRGB(opts.color)
	}

	drawChartFrame(pdf, px, py, pw, ph, minV, maxV, opts.showMajorGridlines, opts.valueFormat)
	pdf.SetStrokeColor(lineR, lineG, lineB)
	pdf.SetFillColor(lineR, lineG, lineB)

	// Build the raw data points.
	rawPts := make([]gopdf.Point, len(values))
	for i, v := range values {
		rawPts[i] = gopdf.Point{
			X: px + (float64(i)*pw)/float64(len(values)-1),
			Y: py + ph - ((v-minV)/rangeV)*(ph-4),
		}
	}

	// Draw connecting lines: straight or Catmull-Rom smooth.
	drawPts := rawPts
	if opts.smooth && len(rawPts) >= 2 {
		drawPts = catmullRomPoints(rawPts, 8)
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
			drawBarDataLabel(pdf, pt.X, pt.Y-8, values[i])
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
		plotR = chartRectWithLegendMargin(r, opts.legendPosition)
	}
	px, py, pw, ph := chartPlotRect(plotR, opts.titleOverlay)

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

	areaR, areaG, areaB := uint8(79), uint8(129), uint8(189)
	if opts.color != "" {
		areaR, areaG, areaB = hexToRGB(opts.color)
	}

	drawChartFrame(pdf, px, py, pw, ph, minV, maxV, opts.showMajorGridlines, opts.valueFormat)

	zeroY := py + ph*maxV/rangeV
	pts := make([]gopdf.Point, 0, len(values)+2)
	pts = append(pts, gopdf.Point{X: px, Y: zeroY})
	for i, v := range values {
		x := px + (float64(i)*pw)/float64(len(values)-1)
		y := py + ph - ((v-minV)/rangeV)*(ph-4)
		pts = append(pts, gopdf.Point{X: x, Y: y})
	}
	pts = append(pts, gopdf.Point{X: px + pw, Y: zeroY})

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
