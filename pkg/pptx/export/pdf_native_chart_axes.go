//nolint:mnd // Chart axis drawing uses tuned visual constants for native PDF fidelity.
package export

import (
	"math"
	"strconv"
	"strings"

	"github.com/signintech/gopdf"
)

// niceStep returns a human-readable tick interval for an axis of the given range,
// matching PowerPoint's own "nice number" axis algorithm.
// It picks the smallest step from {1,2,5,10,20,25,50,…} that gives 4–8 ticks.
func niceStep(rangeV float64) float64 {
	if rangeV <= 0 {
		return 1
	}
	// Scale candidates to the right magnitude.
	magnitude := math.Pow(10, math.Floor(math.Log10(rangeV)))
	for _, m := range []float64{0.1, 0.2, 0.5, 1, 2, 2.5, 5, 10, 20, 25, 50, 100} {
		step := m * magnitude
		if step <= 0 {
			continue
		}
		ticks := math.Ceil(rangeV / step)
		if ticks >= 4 && ticks <= 8 {
			return step
		}
	}
	// Fallback: divide range by 5.
	return math.Ceil(rangeV/5.0) * (magnitude / 10.0)
}

// drawChartFrame draws the Y-axis, baseline, optional horizontal gridlines, and value
// labels. showGridlines controls whether tick-lines cross the plot area. valueFormat is
// applied to all Y-axis labels.
func drawChartFrame(pdf *gopdf.GoPdf, x, y, w, h, minV, maxV float64, showGridlines bool, valueFormat string) {
	pdf.SetStrokeColor(30, 30, 30)
	pdf.Line(x, y, x, y+h)
	pdf.Line(x, y+h, x+w, y+h)
	rangeV := maxV - minV
	if rangeV <= 0 {
		rangeV = 1
	}
	step := rangeV / 5.0
	// Walk upward from minV in step increments, drawing gridline + label per tick.
	for tick := minV; tick <= maxV+step*1e-9; tick = math.Round((tick+step)*1e9) / 1e9 {
		yLine := y + h - ((tick - minV) / rangeV * h)
		if yLine < y-1 || yLine > y+h+1 {
			continue
		}
		if showGridlines {
			pdf.SetStrokeColor(90, 90, 90)
			pdf.Line(x, yLine, x+w, yLine)
		}
		pdf.SetX(x - 28)
		pdf.SetY(yLine - 3)
		_ = pdf.Cell(nil, formatTickValue(tick, valueFormat))
	}
	// Prominent zero line when the range spans positive and negative values.
	if minV < 0 && maxV > 0 {
		zeroY := y + h*maxV/rangeV
		pdf.SetStrokeColor(30, 30, 30)
		pdf.Line(x, zeroY, x+w, zeroY)
	}
}

func drawCategoryLabels(pdf *gopdf.GoPdf, x, y, w, h float64, categories []string) {
	count := len(categories)
	if count <= 0 {
		return
	}
	slot := w / float64(count)
	for i := range count {
		cx := x + slot*float64(i) + slot/2
		pdf.SetX(cx - 6)
		pdf.SetY(y + h + 8)
		_ = pdf.Cell(nil, categoryLabel(categories, i))
	}
}

// drawHorizontalChartFrame draws axes, optional vertical gridlines, X-axis tick labels,
// and category labels for a horizontal bar chart.
func drawHorizontalChartFrame(
	pdf *gopdf.GoPdf,
	x, y, w, h, minV, maxV float64,
	categories []string,
	showGridlines bool,
	valueFormat string,
) {
	pdf.SetStrokeColor(30, 30, 30)
	pdf.Line(x, y, x, y+h)
	pdf.Line(x, y+h, x+w, y+h)
	rangeV := maxV - minV
	if rangeV <= 0 {
		rangeV = 1
	}
	step := rangeV / 5.0
	// Draw vertical gridlines and X-axis labels at step positions.
	for tick := minV; tick <= maxV+step*1e-9; tick = math.Round((tick+step)*1e9) / 1e9 {
		xTick := x + ((tick - minV) / rangeV * w)
		if xTick < x-1 || xTick > x+w+1 {
			continue
		}
		if showGridlines {
			pdf.SetStrokeColor(90, 90, 90)
			pdf.Line(xTick, y, xTick, y+h)
		}
		pdf.SetX(xTick - 4)
		pdf.SetY(y + h + 8)
		_ = pdf.Cell(nil, formatTickValue(tick, valueFormat))
	}
	if minV < 0 && maxV > 0 {
		zeroX := x + ((0-minV)/rangeV)*w
		pdf.SetStrokeColor(30, 30, 30)
		pdf.Line(zeroX, y, zeroX, y+h)
	}
	count := len(categories)
	if count <= 0 {
		return
	}
	slot := h / float64(count)
	for i := range count {
		// Match PowerPoint bar order for horizontal variants (last category at top).
		cy := y + slot*float64(i) + slot/2
		pdf.SetX(x - 18)
		pdf.SetY(cy - 3)
		_ = pdf.Cell(nil, categoryLabel(categories, count-1-i))
	}
}

// fullBars returns a slice of n values all equal to 100, representing 100% bars
// for a single-series 100%-stacked chart (each category = 100% of itself).
func fullBars(n int) []float64 {
	out := make([]float64, n)
	for i := range out {
		out[i] = 100
	}
	return out
}

func niceAxisMax(value float64) float64 {
	if value <= 0 {
		return 1
	}
	if value <= 10 {
		return math.Ceil(value)
	}
	if value <= 100 {
		return math.Ceil(value/5.0) * 5.0
	}
	return math.Ceil(value/50.0) * 50.0
}

// niceAxisRangeXY computes a pleasant [minV, maxV] with ~20% headroom above
// the data max, matching PowerPoint's auto-axis for XY/bubble charts.
// For non-negative data the min is always 0.
func niceAxisRangeXY(values []float64) (float64, float64) {
	if len(values) == 0 {
		return 0, 1
	}
	rawMin, rawMax := minMax(values)
	if rawMin >= 0 {
		// Add 20% padding before rounding to a nice step — same result as PowerPoint auto-axis.
		padded := rawMax * 1.2
		step := niceStep(padded)
		if step <= 0 {
			step = 1
		}
		return 0, math.Ceil(padded/step) * step
	}
	// Fall back to symmetric range for mixed/negative data.
	return niceAxisRange(values)
}

// niceAxisRange computes a pleasant [minV, maxV] that covers all values.
// Non-negative data → min pinned to 0. Non-positive → max pinned to 0.
// Mixed → symmetric around 0.
func niceAxisRange(values []float64) (float64, float64) {
	if len(values) == 0 {
		return 0, 1
	}
	rawMin, rawMax := minMax(values)
	if rawMin >= 0 {
		return 0, niceAxisMax(rawMax)
	}
	if rawMax <= 0 {
		return -niceAxisMax(-rawMin), 0
	}
	absMax := math.Max(math.Abs(rawMin), math.Abs(rawMax))
	nice := niceAxisMax(absMax)
	return -nice, nice
}

// drawChartLegend draws a colour-coded legend inside r at the given position.
func drawChartLegend(pdf *gopdf.GoPdf, r chartRect, position string, entries []legendEntry) {
	if len(entries) == 0 {
		return
	}
	var x, y float64
	switch position {
	case "l":
		x, y = r.x+4, r.y+r.h*0.25
	case "t":
		x, y = r.x+r.w/2-50, r.y+20
	case "b":
		x, y = r.x+r.w/2-50, r.y+r.h-28
	default: // "r"
		x, y = r.x+r.w-98, r.y+r.h*0.25
	}
	for i, entry := range entries {
		ey := y + float64(i)*18
		pdf.SetFillColor(entry.R, entry.G, entry.B)
		pdf.RectFromUpperLeftWithStyle(x, ey, 12, 8, "F")
		pdf.SetTextColor(40, 40, 40)
		name := entry.Name
		if name == "" {
			name = "Series " + strconv.Itoa(i+1)
		}
		pdf.SetX(x + 16)
		pdf.SetY(ey - 1)
		_ = pdf.Cell(nil, name)
	}
	pdf.SetTextColor(0, 0, 0)
}

// drawBarDataLabel draws a numeric value label near a bar top or line point.
func drawBarDataLabel(pdf *gopdf.GoPdf, cx, labelY, value float64) {
	label := strconv.FormatFloat(value, 'f', 1, 64)
	label = strings.TrimSuffix(label, ".0")
	pdf.SetTextColor(60, 60, 60)
	pdf.SetX(cx - float64(len(label))*3)
	pdf.SetY(labelY)
	_ = pdf.Cell(nil, label)
	pdf.SetTextColor(0, 0, 0)
}

// drawPieSliceLabel draws a label at radius*1.2 from the pie centre.
func drawPieSliceLabel(pdf *gopdf.GoPdf, cx, cy, radius, midAngle float64, text string) {
	labelRadius := radius * 1.2
	lx := cx + math.Cos(midAngle)*labelRadius
	ly := cy + math.Sin(midAngle)*labelRadius
	pdf.SetTextColor(40, 40, 40)
	pdf.SetX(lx - float64(len(text))*3)
	pdf.SetY(ly - 4)
	_ = pdf.Cell(nil, text)
	pdf.SetTextColor(0, 0, 0)
}

// catmullRomPoints returns a smooth curve through pts using Catmull-Rom spline interpolation.
func catmullRomPoints(pts []gopdf.Point, segsPerInterval int) []gopdf.Point {
	n := len(pts)
	if n < 2 {
		return pts
	}
	if segsPerInterval < 4 {
		segsPerInterval = 8
	}
	ext := make([]gopdf.Point, 0, n+2)
	ext = append(ext, gopdf.Point{X: 2*pts[0].X - pts[1].X, Y: 2*pts[0].Y - pts[1].Y})
	ext = append(ext, pts...)
	ext = append(ext, gopdf.Point{X: 2*pts[n-1].X - pts[n-2].X, Y: 2*pts[n-1].Y - pts[n-2].Y})

	out := make([]gopdf.Point, 0, (n-1)*segsPerInterval+1)
	for i := 1; i < len(ext)-2; i++ {
		p0, p1, p2, p3 := ext[i-1], ext[i], ext[i+1], ext[i+2]
		for s := range segsPerInterval {
			t := float64(s) / float64(segsPerInterval)
			t2, t3 := t*t, t*t*t
			out = append(out, gopdf.Point{
				X: 0.5 * ((2 * p1.X) + (-p0.X+p2.X)*t + (2*p0.X-5*p1.X+4*p2.X-p3.X)*t2 + (-p0.X+3*p1.X-3*p2.X+p3.X)*t3),
				Y: 0.5 * ((2 * p1.Y) + (-p0.Y+p2.Y)*t + (2*p0.Y-5*p1.Y+4*p2.Y-p3.Y)*t2 + (-p0.Y+3*p1.Y-3*p2.Y+p3.Y)*t3),
			})
		}
	}
	out = append(out, pts[n-1])
	return out
}

// drawFilledCircle draws a filled polygon approximating a circle.
func drawFilledCircle(pdf *gopdf.GoPdf, cx, cy, r float64, colR, colG, colB uint8) {
	const steps = 16
	pts := make([]gopdf.Point, 0, steps)
	for i := range steps {
		angle := 2 * math.Pi * float64(i) / steps
		pts = append(pts, gopdf.Point{
			X: cx + math.Cos(angle)*r,
			Y: cy + math.Sin(angle)*r,
		})
	}
	pdf.SetFillColor(colR, colG, colB)
	pdf.Polygon(pts, "F")
}
