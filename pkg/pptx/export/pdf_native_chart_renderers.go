//nolint:mnd // Native chart rendering uses tuned numeric drawing constants.
package export

import (
	"math"

	"github.com/signintech/gopdf"
)

func renderBarLike(
	pdf *gopdf.GoPdf,
	title string,
	r chartRect,
	values []float64,
	categories []string,
	horizontal bool,
) {
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
		drawHorizontalChartFrame(pdf, px, py, pw, ph, maxV, categories)
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
		drawCategoryLabels(pdf, px, py, pw, ph, categories)
	}
}

func renderLineLike(
	pdf *gopdf.GoPdf,
	title string,
	r chartRect,
	values []float64,
	categories []string,
	markers bool,
) {
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
	drawCategoryLabels(pdf, px, py, pw, ph, categories)
}

func renderAreaLike(
	pdf *gopdf.GoPdf,
	title string,
	r chartRect,
	values []float64,
	categories []string,
) {
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
	drawCategoryLabels(pdf, px, py, pw, ph, categories)
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
	for i := range n {
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
