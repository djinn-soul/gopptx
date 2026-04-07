//nolint:mnd // Scatter chart helpers use fixed geometry constants to match PPT output.
package export

import (
	"math"

	"github.com/signintech/gopdf"
)

func renderScatterAxes(
	pdf *gopdf.GoPdf,
	px, py, pw, ph,
	minX, maxX, minY, maxY, rangeX, rangeY float64,
	opts chartSeriesOpts,
) {
	renderScatterGridlines(pdf, px, py, pw, ph, minX, maxX, minY, maxY, rangeX, rangeY, opts)
	renderScatterTickLabels(pdf, px, py, pw, ph, minX, maxX, minY, maxY, rangeX, rangeY, opts)
}

func renderScatterGridlines(
	pdf *gopdf.GoPdf,
	px, py, pw, ph,
	minX, maxX, minY, maxY, rangeX, rangeY float64,
	opts chartSeriesOpts,
) {
	stepX := niceStep(rangeX)
	stepY := niceStep(rangeY)

	pdf.SetStrokeColor(90, 90, 90)
	if opts.showCatGridlines {
		for tick := minX; tick <= maxX+stepX*1e-9; tick = math.Round((tick+stepX)*1e9) / 1e9 {
			xg := px + (tick-minX)/rangeX*pw
			if xg >= px-1 && xg <= px+pw+1 {
				pdf.Line(xg, py, xg, py+ph)
			}
		}
	}
	if opts.showMajorGridlines {
		for tick := minY; tick <= maxY+stepY*1e-9; tick = math.Round((tick+stepY)*1e9) / 1e9 {
			yg := py + ph - (tick-minY)/rangeY*ph
			if yg >= py-1 && yg <= py+ph+1 {
				pdf.Line(px, yg, px+pw, yg)
			}
		}
	}
}

func renderScatterTickLabels(
	pdf *gopdf.GoPdf,
	px, py, pw, ph,
	minX, maxX, minY, maxY, rangeX, rangeY float64,
	opts chartSeriesOpts,
) {
	stepX := niceStep(rangeX)
	stepY := niceStep(rangeY)
	pdf.SetStrokeColor(30, 30, 30)
	for tick := minX; tick <= maxX+stepX*1e-9; tick = math.Round((tick+stepX)*1e9) / 1e9 {
		xTick := px + (tick-minX)/rangeX*pw
		if xTick < px-1 || xTick > px+pw+1 {
			continue
		}
		pdf.SetX(xTick - 6)
		pdf.SetY(py + ph + 8)
		_ = pdf.Cell(nil, formatTickValue(tick, opts.valueFormat))
	}
	for tick := minY; tick <= maxY+stepY*1e-9; tick = math.Round((tick+stepY)*1e9) / 1e9 {
		yTick := py + ph - (tick-minY)/rangeY*ph
		if yTick < py-1 || yTick > py+ph+1 {
			continue
		}
		pdf.SetX(px - 28)
		pdf.SetY(yTick - 3)
		_ = pdf.Cell(nil, formatTickValue(tick, opts.valueFormat))
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
