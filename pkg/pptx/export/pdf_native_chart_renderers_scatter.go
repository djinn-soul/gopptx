//nolint:mnd // Chart scatter/pie/radar rendering uses tuned visual constants.
package export

import (
	"math"

	"github.com/signintech/gopdf"
)

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
	// The pie fills its plot area rather than a fixed fraction of the frame:
	// PowerPoint draws a 360pt-tall chart's pie at a 154pt radius, which is half
	// the height left once the top and bottom margins are taken off. The old
	// 0.35 × min(w,h) drew it a fifth too small.
	plot := piePlotRect(pdf, r, opts, categories)
	cx := plot.x + plot.w/2
	cy := plot.y + plot.h/2
	radius := math.Min(plot.w, plot.h) / 2
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

		if opts.showDataLabels {
			midAngle := (start + end) / 2
			drawPieSliceLabel(pdf, cx, cy, radius, midAngle, pieSliceLabelText(opts, categories, i, v, total, frac))
		}
		start = end
	}
	isDoughnut := holeSizePct >= 10
	if isDoughnut {
		inner := radius * float64(clampHoleSize(holeSizePct)) / 100.0
		drawFilledCircle(pdf, cx, cy, inner, 255, 255, 255)
	}
	if opts.showLegend {
		drawChartLegend(pdf, r, opts.legendPosition, entries)
	}
}

// piePlotRect is the area a pie or doughnut is drawn into: the chart frame less
// the legend, the title and the margins PowerPoint keeps around the plot.
func piePlotRect(pdf *gopdf.GoPdf, r chartRect, opts chartSeriesOpts, names []string) chartRect {
	if opts.showLegend {
		r = chartRectWithLegendMargin(pdf, r, opts.legendPosition, names)
	}
	top := chartMinTopPadPt
	if opts.titleOverlay {
		top = chartOverlayTopPadPt
	}
	return chartRect{
		x: r.x + piePlotSideMarginPt,
		y: r.y + top,
		w: math.Max(r.w-2*piePlotSideMarginPt, 1),
		h: math.Max(r.h-top-piePlotBottomMarginPt, 1),
	}
}

// piePlotSideMarginPt and piePlotBottomMarginPt are the clearance PowerPoint
// leaves beside and below a pie, measured off its own export.
const (
	piePlotSideMarginPt   = 10.0
	piePlotBottomMarginPt = 10.0
)

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
		plotR = chartRectWithLegendMargin(pdf, r, opts.legendPosition, []string{opts.seriesName})
	}
	// Use XY-specific axis range (adds ~20% headroom) matching PowerPoint's auto-axis
	// for scatter/bubble charts — ensures data points never sit on the axis edge.
	minX, maxX := niceAxisRangeXY(xs[:n])
	minY, maxY := niceAxisRangeXY(ys[:n])
	if maxX <= minX {
		maxX = minX + 1
	}
	if maxY <= minY {
		maxY = minY + 1
	}

	// Both axes carry numeric tick labels here, so the plot area is sized around
	// them; the ranges therefore have to be resolved first.
	layout := solveChartLayout(
		pdf, plotR, opts.titleOverlay,
		title,
		chartAxisSpec{MinV: minY, MaxV: maxY, ValueFormat: opts.valueFormat},
		chartAxisSpec{MinV: minX, MaxV: maxX, ValueFormat: opts.valueFormat},
	)
	px, py, pw, ph := layout.X, layout.Y, layout.W, layout.H
	rangeX := maxX - minX
	rangeY := maxY - minY

	// Draw plot frame.
	pdf.SetStrokeColor(30, 30, 30)
	pdf.RectFromUpperLeftWithStyle(px, py, pw, ph, "D")
	renderScatterAxes(pdf, scatterAxisGeometry{
		px: px, py: py, pw: pw, ph: ph,
		minX: minX, maxX: maxX, minY: minY, maxY: maxY,
		rangeX: rangeX, rangeY: rangeY,
		densityX: layout.Horizontal, densityY: layout.Vertical,
	}, opts)

	ptR, ptG, ptB := uint8(79), uint8(129), uint8(189)
	if opts.color != "" {
		ptR, ptG, ptB = hexToRGB(opts.color)
	}

	plotPts := make([]gopdf.Point, n)
	for i := range n {
		xf := (xs[i] - minX) / rangeX
		yf := (ys[i] - minY) / rangeY
		plotPts[i] = gopdf.Point{
			X: px + xf*pw,
			Y: py + ph - yf*ph,
		}
	}

	style := opts.scatterStyle
	if style == "" {
		style = "marker"
	}
	drawScatterConnectingLines(pdf, plotPts, style, n, ptR, ptG, ptB)

	bubbleScale := 1.0
	if opts.bubbleScale > 0 {
		bubbleScale = float64(opts.bubbleScale) / 100.0
	}
	drawScatterPoints(pdf, plotPts, sizes, bubbleScale, ptR, ptG, ptB, pw, ph)

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
func renderRadarLike(
	pdf *gopdf.GoPdf, title string, r chartRect,
	values []float64, categories []string, filled bool, opts chartSeriesOpts,
) {
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
		drawChartLabel(pdf, categoryLabel(categories, i), lx, ly, chartLabelFontSize, chartTextCenter)
	}
	pdf.SetTextColor(0, 0, 0)

	if opts.showLegend {
		drawChartLegend(pdf, r, opts.legendPosition, []legendEntry{
			{Name: opts.seriesName, R: lineR, G: lineG, B: lineB},
		})
	}
}
