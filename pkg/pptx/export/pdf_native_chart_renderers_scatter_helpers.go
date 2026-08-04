//nolint:mnd // Scatter chart helpers use fixed geometry constants to match PPT output.
package export

import (
	"math"

	"github.com/signintech/gopdf"
)

// scatterAxisGeometry bundles the plot rect and both axis ranges, which every
// scatter axis helper needs together.
type scatterAxisGeometry struct {
	px, py, pw, ph         float64
	minX, maxX, minY, maxY float64
	rangeX, rangeY         float64
	densityX, densityY     chartAxisTickDensity
}

func renderScatterAxes(pdf *gopdf.GoPdf, g scatterAxisGeometry, opts chartSeriesOpts) {
	renderScatterGridlines(pdf, g, opts)
	renderScatterTickLabels(pdf, g, opts)
}

func renderScatterGridlines(pdf *gopdf.GoPdf, g scatterAxisGeometry, opts chartSeriesOpts) {
	pdf.SetStrokeColor(chartGridlineGrey, chartGridlineGrey, chartGridlineGrey)
	if opts.showCatGridlines {
		for _, tick := range chartAxisTicks(g.minX, g.maxX, g.densityX) {
			xg := g.px + (tick-g.minX)/g.rangeX*g.pw
			if xg >= g.px-1 && xg <= g.px+g.pw+1 {
				pdf.Line(xg, g.py, xg, g.py+g.ph)
			}
		}
	}
	if opts.showMajorGridlines {
		for _, tick := range chartAxisTicks(g.minY, g.maxY, g.densityY) {
			yg := g.py + g.ph - (tick-g.minY)/g.rangeY*g.ph
			if yg >= g.py-1 && yg <= g.py+g.ph+1 {
				pdf.Line(g.px, yg, g.px+g.pw, yg)
			}
		}
	}
}

func renderScatterTickLabels(pdf *gopdf.GoPdf, g scatterAxisGeometry, opts chartSeriesOpts) {
	pdf.SetStrokeColor(30, 30, 30)
	for _, tick := range chartAxisTicks(g.minX, g.maxX, g.densityX) {
		xTick := g.px + (tick-g.minX)/g.rangeX*g.pw
		if xTick < g.px-1 || xTick > g.px+g.pw+1 {
			continue
		}
		drawChartAxisTick(pdf, xTick, g.py+g.ph, chartAxisHorizontal)
		drawChartLabel(
			pdf, formatTickValue(tick, opts.valueFormat),
			xTick, g.py+g.ph+chartTickMarkPt+chartAxisLabelGapPt+pdfLineHeight(chartLabelFontSize)/2,
			chartLabelFontSize, chartTextCenter,
		)
	}
	for _, tick := range chartAxisTicks(g.minY, g.maxY, g.densityY) {
		yTick := g.py + g.ph - (tick-g.minY)/g.rangeY*g.ph
		if yTick < g.py-1 || yTick > g.py+g.ph+1 {
			continue
		}
		drawChartAxisTick(pdf, g.px, yTick, chartAxisVertical)
		drawChartLabel(
			pdf, formatTickValue(tick, opts.valueFormat),
			g.px-chartTickMarkPt-chartAxisLabelGapPt, yTick,
			chartLabelFontSize, chartTextRight,
		)
	}
}

func drawScatterConnectingLines(pdf *gopdf.GoPdf, plotPts []gopdf.Point, style string, n int, ptR, ptG, ptB uint8) {
	if style != "lineMarker" && style != "smoothMarker" {
		return
	}
	linePts := plotPts
	if style == "smoothMarker" && n >= 2 {
		linePts = catmullRomPoints(plotPts, 8)
	}
	pdf.SetStrokeColor(ptR, ptG, ptB)
	for i := 1; i < len(linePts); i++ {
		pdf.Line(linePts[i-1].X, linePts[i-1].Y, linePts[i].X, linePts[i].Y)
	}
}

func drawScatterPoints(
	pdf *gopdf.GoPdf,
	plotPts []gopdf.Point,
	sizes []float64,
	bubbleScale float64,
	ptR, ptG, ptB uint8,
	pw, ph float64,
) {
	isBubble := len(sizes) > 0
	var maxSize float64
	if isBubble {
		for _, s := range sizes {
			if s > maxSize {
				maxSize = s
			}
		}
	}
	refRadius := math.Min(pw, ph) * 0.15
	if bubbleScale > 0 {
		refRadius *= bubbleScale
	}
	for i, pt := range plotPts {
		rad := 3.0
		if isBubble && i < len(sizes) && sizes[i] > 0 && maxSize > 0 {
			rad = refRadius * math.Sqrt(sizes[i]/maxSize)
		}
		drawFilledCircle(pdf, pt.X, pt.Y, rad, ptR, ptG, ptB)
	}
}
