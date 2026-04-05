//nolint:mnd // Native chart rendering uses tuned numeric drawing constants.
package export

import (
	"math"
	"strconv"

	"github.com/signintech/gopdf"
)

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

	pdf.SetFillColor(barR, barG, barB)
	for i, v := range values {
		if horizontal {
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
		} else {
			slot := pw / float64(len(values))
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
				drawBarDataLabel(pdf, bx+bw/2, labelY, v)
			}
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

	lineR, lineG, lineB := uint8(192), uint8(80), uint8(77)
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

	areaR, areaG, areaB := uint8(155), uint8(187), uint8(89)
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

// renderPieLike renders a pie or doughnut chart.
// holeSizePct controls the doughnut hole (0 = solid pie, 10–90 = doughnut).
func renderPieLike(
	pdf *gopdf.GoPdf,
	title string,
	r chartRect,
	values []float64,
	holeSizePct int,
	categories []string,
	opts chartSeriesOpts,
) {
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
	entries := make([]legendEntry, 0, len(values))
	for i, v := range values {
		frac := v / total
		end := start + frac*2*math.Pi
		rC, gC, bC := pieColor(i)
		drawWedge(pdf, cx, cy, radius, start, end, rC, gC, bC)
		entries = append(entries, legendEntry{Name: categoryLabel(categories, i), R: rC, G: gC, B: bC})

		if opts.showDataLabels || (len(categories) > 0 && !opts.showDataLabels) {
			midAngle := (start + end) / 2
			if opts.showDataLabels {
				pct := frac * 100
				label := strconv.FormatFloat(pct, 'f', 1, 64) + "%"
				drawPieSliceLabel(pdf, cx, cy, radius, midAngle, label)
			} else if i < len(categories) && categories[i] != "" {
				drawPieSliceLabel(pdf, cx, cy, radius, midAngle, categories[i])
			}
		}
		start = end
	}
	isDoughnut := holeSizePct >= 10
	if isDoughnut {
		pdf.SetFillColor(255, 255, 255)
		inner := radius * float64(clampHoleSize(holeSizePct)) / 100.0
		pdf.Oval(cx-inner, cy-inner, cx+inner, cy+inner)
	}
	if opts.showLegend {
		drawChartLegend(pdf, r, opts.legendPosition, entries)
	}
}

func clampHoleSize(pct int) int {
	if pct < 10 {
		return 10
	}
	if pct > 90 {
		return 90
	}
	return pct
}

// renderScatterLike renders an XY scatter or bubble chart.
// opts.scatterStyle controls the visual style: "marker" (default), "lineMarker", "smoothMarker".
// opts.bubbleScale scales bubble radii (1–300 percent; 0 uses the default of 100).
func renderScatterLike(pdf *gopdf.GoPdf, title string, r chartRect, xs, ys, sizes []float64, opts chartSeriesOpts) {
	renderChartTitle(pdf, title, r)
	if len(xs) == 0 || len(ys) == 0 {
		return
	}
	n := min(len(xs), len(ys))

	plotR := r
	if opts.showLegend {
		plotR = chartRectWithLegendMargin(r, opts.legendPosition)
	}
	px, py, pw, ph := chartPlotRect(plotR, opts.titleOverlay)

	minX, maxX := minMax(xs[:n])
	minY, maxY := minMax(ys[:n])
	if maxX <= minX {
		maxX = minX + 1
	}
	if maxY <= minY {
		maxY = minY + 1
	}

	// Outer box.
	pdf.SetStrokeColor(30, 30, 30)
	pdf.RectFromUpperLeftWithStyle(px, py, pw, ph, "D")

	// Optional gridlines.
	if opts.showMajorGridlines {
		pdf.SetStrokeColor(90, 90, 90)
		for i := 1; i < 5; i++ {
			yg := py + float64(i)*ph/5
			pdf.Line(px, yg, px+pw, yg)
			xg := px + float64(i)*pw/5
			pdf.Line(xg, py, xg, py+ph)
		}
	}

	ptR, ptG, ptB := uint8(79), uint8(129), uint8(189)
	if opts.color != "" {
		ptR, ptG, ptB = hexToRGB(opts.color)
	}

	// Map data values to plot coordinates.
	plotPts := make([]gopdf.Point, n)
	for i := range n {
		xf := (xs[i] - minX) / (maxX - minX)
		yf := (ys[i] - minY) / (maxY - minY)
		plotPts[i] = gopdf.Point{
			X: px + xf*(pw-6) + 3,
			Y: py + ph - yf*(ph-6) - 3,
		}
	}

	style := opts.scatterStyle
	if style == "" {
		style = "marker"
	}

	// Draw connecting lines for lineMarker / smoothMarker styles.
	if style == "lineMarker" || style == "smoothMarker" {
		linePts := plotPts
		if style == "smoothMarker" && n >= 2 {
			linePts = catmullRomPoints(plotPts, 8)
		}
		pdf.SetStrokeColor(ptR, ptG, ptB)
		for i := 1; i < len(linePts); i++ {
			pdf.Line(linePts[i-1].X, linePts[i-1].Y, linePts[i].X, linePts[i].Y)
		}
	}

	// Bubble scale factor (default 100%).
	bubbleScale := 1.0
	if opts.bubbleScale > 0 {
		bubbleScale = float64(opts.bubbleScale) / 100.0
	}

	// Draw markers at original data coordinates.
	for i, pt := range plotPts {
		rad := 3.0
		if i < len(sizes) && sizes[i] > 0 {
			rad = math.Min(16, (2+math.Sqrt(sizes[i])*0.8)*bubbleScale)
		}
		drawFilledCircle(pdf, pt.X, pt.Y, rad, ptR, ptG, ptB)
	}

	if opts.catAxisTitle != "" {
		drawCategoryAxisTitle(pdf, px, py, pw, ph, opts.catAxisTitle)
	}
	if opts.valAxisTitle != "" {
		drawValueAxisTitle(pdf, px, py, pw, ph, opts.valAxisTitle)
	}
	if opts.showLegend {
		drawChartLegend(pdf, r, opts.legendPosition, []legendEntry{
			{Name: opts.seriesName, R: ptR, G: ptG, B: ptB},
		})
	}
}

// renderRadarLike renders a radar (spider) chart with concentric grid rings,
// spoke lines, perimeter category labels, and an optional legend.
func renderRadarLike(pdf *gopdf.GoPdf, title string, r chartRect, values []float64, categories []string, filled bool, opts chartSeriesOpts) {
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
	n := len(values)

	// Draw concentric grid rings.
	const numRings = 4
	for ring := 1; ring <= numRings; ring++ {
		ringRadius := radius * float64(ring) / numRings
		ringPts := make([]gopdf.Point, 0, n)
		for i := range n {
			angle := -math.Pi/2 + (2*math.Pi*float64(i))/float64(n)
			ringPts = append(ringPts, gopdf.Point{
				X: cx + math.Cos(angle)*ringRadius,
				Y: cy + math.Sin(angle)*ringRadius,
			})
		}
		pdf.SetStrokeColor(180, 180, 180)
		pdf.Polygon(ringPts, "D")
	}
	// Draw spokes from centre to perimeter.
	for i := range n {
		angle := -math.Pi/2 + (2*math.Pi*float64(i))/float64(n)
		pdf.SetStrokeColor(180, 180, 180)
		pdf.Line(cx, cy, cx+math.Cos(angle)*radius, cy+math.Sin(angle)*radius)
	}

	// Data polygon.
	lineR, lineG, lineB := uint8(70), uint8(120), uint8(180)
	if opts.color != "" {
		lineR, lineG, lineB = hexToRGB(opts.color)
	}
	pts := make([]gopdf.Point, 0, n)
	for i, v := range values {
		angle := -math.Pi/2 + (2*math.Pi*float64(i))/float64(n)
		scale := v / maxV
		pts = append(pts, gopdf.Point{X: cx + math.Cos(angle)*radius*scale, Y: cy + math.Sin(angle)*radius*scale})
	}
	pdf.SetStrokeColor(lineR, lineG, lineB)
	if filled {
		pdf.SetFillColor(lineR, lineG, lineB)
		pdf.Polygon(pts, "FD")
	} else {
		pdf.Polygon(pts, "D")
	}

	// Perimeter category labels.
	pdf.SetTextColor(40, 40, 40)
	for i := range n {
		angle := -math.Pi/2 + (2*math.Pi*float64(i))/float64(n)
		labelR := radius + 14
		lx := cx + math.Cos(angle)*labelR
		ly := cy + math.Sin(angle)*labelR
		label := categoryLabel(categories, i)
		pdf.SetX(lx - float64(len(label))*3)
		pdf.SetY(ly - 4)
		_ = pdf.Cell(nil, label)
	}
	pdf.SetTextColor(0, 0, 0)

	if opts.showLegend {
		drawChartLegend(pdf, r, opts.legendPosition, []legendEntry{
			{Name: opts.seriesName, R: lineR, G: lineG, B: lineB},
		})
	}
}
