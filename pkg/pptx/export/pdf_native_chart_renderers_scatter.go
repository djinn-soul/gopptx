//nolint:mnd // Chart scatter/pie/radar rendering uses tuned visual constants.
package export

import (
	"math"
	"strconv"

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

		if opts.showDataLabels {
			midAngle := (start + end) / 2
			if opts.showCatName && i < len(categories) && categories[i] != "" {
				drawPieSliceLabel(pdf, cx, cy, radius, midAngle, categories[i])
			} else {
				pct := frac * 100
				label := strconv.FormatFloat(pct, 'f', 1, 64) + "%"
				drawPieSliceLabel(pdf, cx, cy, radius, midAngle, label)
			}
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
	rangeX := maxX - minX
	rangeY := maxY - minY

	// Draw plot frame.
	pdf.SetStrokeColor(30, 30, 30)
	pdf.RectFromUpperLeftWithStyle(px, py, pw, ph, "D")
	renderScatterAxes(pdf, px, py, pw, ph, minX, maxX, minY, maxY, rangeX, rangeY, opts)

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
