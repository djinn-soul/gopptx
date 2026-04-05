//nolint:mnd // Chart helper math uses tuned visual constants for native PDF fidelity.
package export

import (
	"math"
	"strconv"
	"strings"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

// legendEntry describes one entry in a chart legend (series name + color).
type legendEntry struct {
	Name    string
	R, G, B uint8
}

// chartSeriesOpts holds optional rendering hints for chart renderers.
type chartSeriesOpts struct {
	color              string   // hex color override; empty = use renderer default
	minValue           *float64 // axis min override
	maxValue           *float64 // axis max override
	showLegend         bool
	legendPosition     string // "r","l","t","b"
	seriesName         string
	showDataLabels     bool
	catAxisTitle       string
	valAxisTitle       string
	scatterStyle       string // "marker" | "lineMarker" | "smoothMarker"
	smooth             bool   // draw line chart with Catmull-Rom smooth curves
	showMajorGridlines bool   // draw horizontal value-axis gridlines
	showCatGridlines   bool   // draw vertical category-axis gridlines (horizontal charts)
	titleOverlay       bool   // title overlaps plot area; don't reserve top padding
	valueFormat        string // Excel-style number format ("General", "0%", "$#,##0", …)
	bubbleScale        int    // bubble size scale percent (1–300; 0 = use renderer default)
}

// categoryLabel returns the label for index i, falling back to "Q<i+1>".
func categoryLabel(categories []string, i int) string {
	if i < len(categories) && categories[i] != "" {
		return categories[i]
	}
	return "Q" + strconv.Itoa(i+1)
}

// formatTickValue formats a numeric axis-tick value using the given Excel-style format.
// Supported: "General"/empty → rounded integer; contains "%" → append "%";
// starts with "$" → prepend "$". Anything else falls back to rounded integer.
func formatTickValue(v float64, format string) string {
	if format == "" || format == "General" {
		return strconv.Itoa(int(math.Round(v)))
	}
	if strings.Contains(format, "%") {
		return strconv.Itoa(int(math.Round(v))) + "%"
	}
	if strings.HasPrefix(format, "$") {
		return "$" + strconv.FormatFloat(v, 'f', 0, 64)
	}
	return strconv.Itoa(int(math.Round(v)))
}

func pieColor(i int) (uint8, uint8, uint8) {
	palette := [][3]uint8{
		{79, 129, 189},
		{192, 80, 77},
		{155, 187, 89},
		{128, 100, 162},
		{75, 172, 198},
		{247, 150, 70},
	}
	c := palette[i%len(palette)]
	return c[0], c[1], c[2]
}

func drawWedge(pdf *gopdf.GoPdf, cx, cy, radius, start, end float64, r, g, b uint8) {
	pdf.SetFillColor(r, g, b)
	pts := []gopdf.Point{{X: cx, Y: cy}}
	steps := max(8, int(math.Ceil((end-start)/(math.Pi/18))))
	for i := 0; i <= steps; i++ {
		t := start + (end-start)*float64(i)/float64(steps)
		pts = append(pts, gopdf.Point{X: cx + radius*math.Cos(t), Y: cy + radius*math.Sin(t)})
	}
	pdf.Polygon(pts, "F")
}

func maxFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	v := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] > v {
			v = values[i]
		}
	}
	return v
}

func minFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	v := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] < v {
			v = values[i]
		}
	}
	return v
}

func sumFloat(values []float64) float64 {
	total := 0.0
	for _, v := range values {
		total += v
	}
	return total
}

func minMax(values []float64) (float64, float64) {
	if len(values) == 0 {
		return 0, 0
	}
	minV, maxV := values[0], values[0]
	for i := 1; i < len(values); i++ {
		if values[i] < minV {
			minV = values[i]
		}
		if values[i] > maxV {
			maxV = values[i]
		}
	}
	return minV, maxV
}

// chartPlotRect returns the inner (x,y,w,h) of the chart plot area, reserving space
// for axes, labels, and (optionally) a title. When titleOverlay is true the title
// overlaps the chart body, so minimal top padding is applied.
func chartPlotRect(r chartRect, titleOverlay bool) (float64, float64, float64, float64) {
	leftPad := math.Max(36, r.w*0.08)
	rightPad := math.Max(16, r.w*0.07)
	topPad := math.Max(24, r.h*0.12)
	if titleOverlay {
		topPad = 4
	}
	bottomPad := math.Max(26, r.h*0.12)
	return r.x + leftPad, r.y + topPad, r.w - leftPad - rightPad, r.h - topPad - bottomPad
}

func chartPlotRectHorizontal(r chartRect, titleOverlay bool) (float64, float64, float64, float64) {
	leftPad := math.Max(18, r.w*0.03)
	rightPad := math.Max(26, r.w*0.10)
	topPad := math.Max(24, r.h*0.12)
	if titleOverlay {
		topPad = 4
	}
	bottomPad := math.Max(26, r.h*0.12)
	return r.x + leftPad, r.y + topPad, r.w - leftPad - rightPad, r.h - topPad - bottomPad
}

// chartRectWithLegendMargin shrinks the chart rect to make room for the legend.
func chartRectWithLegendMargin(r chartRect, pos string) chartRect {
	const legendW = 110
	const legendH = 36
	switch pos {
	case "l":
		return chartRect{r.x + legendW, r.y, r.w - legendW, r.h}
	case "t":
		return chartRect{r.x, r.y + legendH, r.w, r.h - legendH}
	case "b":
		return chartRect{r.x, r.y, r.w, r.h - legendH}
	default: // "r"
		return chartRect{r.x, r.y, r.w - legendW, r.h}
	}
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
	for i := 1; i <= 5; i++ {
		yLine := y + h - (float64(i)/5.0)*h
		if showGridlines {
			pdf.SetStrokeColor(90, 90, 90)
			pdf.Line(x, yLine, x+w, yLine)
		}
		labelVal := minV + rangeV*float64(i)/5.0
		pdf.SetX(x - 28)
		pdf.SetY(yLine - 3)
		_ = pdf.Cell(nil, formatTickValue(labelVal, valueFormat))
	}
	pdf.SetX(x - 28)
	pdf.SetY(y + h - 3)
	_ = pdf.Cell(nil, formatTickValue(minV, valueFormat))
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
	const steps = 6
	for i := 1; i <= steps; i++ {
		xLine := x + (float64(i)/float64(steps))*w
		if showGridlines {
			pdf.SetStrokeColor(90, 90, 90)
			pdf.Line(xLine, y, xLine, y+h)
		}
	}
	if minV < 0 && maxV > 0 {
		zeroX := x + ((0-minV)/rangeV)*w
		pdf.SetStrokeColor(30, 30, 30)
		pdf.Line(zeroX, y, zeroX, y+h)
	}
	for i := 0; i <= steps; i++ {
		xTick := x + (float64(i)/float64(steps))*w
		tickValue := minV + rangeV*float64(i)/float64(steps)
		pdf.SetX(xTick - 4)
		pdf.SetY(y + h + 8)
		_ = pdf.Cell(nil, formatTickValue(tickValue, valueFormat))
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

// normalizePercentSeries converts values to their percentage of the total (0–100).
func normalizePercentSeries(values []float64) []float64 {
	if len(values) == 0 {
		return nil
	}
	total := sumFloat(values)
	if total <= 0 {
		return nil
	}
	out := make([]float64, len(values))
	for i, v := range values {
		out[i] = (v / total) * 100
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

// drawCategoryAxisTitle draws the category (X) axis title below the category labels.
func drawCategoryAxisTitle(pdf *gopdf.GoPdf, px, py, pw, ph float64, title string) {
	pdf.SetTextColor(60, 60, 60)
	pdf.SetX(px + pw/2 - float64(len(title))*3)
	pdf.SetY(py + ph + 26)
	_ = pdf.Cell(nil, title)
	pdf.SetTextColor(0, 0, 0)
}

// drawValueAxisTitle draws the value (Y) axis title rotated 90° to the left of the Y-axis.
func drawValueAxisTitle(pdf *gopdf.GoPdf, px, py, pw, ph float64, title string) {
	pdf.SetTextColor(60, 60, 60)
	titleX := px - 42
	titleY := py + ph/2
	pdf.Rotate(-90, titleX, titleY)
	pdf.SetX(titleX - float64(len(title))*3)
	pdf.SetY(titleY - 3)
	_ = pdf.Cell(nil, title)
	pdf.RotateReset()
	pdf.SetTextColor(0, 0, 0)
}

// drawBarDataLabel draws a numeric value label near a bar top or line point.
func drawBarDataLabel(pdf *gopdf.GoPdf, cx, labelY, value float64) {
	label := strconv.FormatFloat(value, 'f', 1, 64)
	if strings.HasSuffix(label, ".0") {
		label = label[:len(label)-2]
	}
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
				X: 0.5 * ((2*p1.X) + (-p0.X+p2.X)*t + (2*p0.X-5*p1.X+4*p2.X-p3.X)*t2 + (-p0.X+3*p1.X-3*p2.X+p3.X)*t3),
				Y: 0.5 * ((2*p1.Y) + (-p0.Y+p2.Y)*t + (2*p0.Y-5*p1.Y+4*p2.Y-p3.Y)*t2 + (-p0.Y+3*p1.Y-3*p2.Y+p3.Y)*t3),
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
	px, py, pw, ph := chartPlotRect(r, opts.titleOverlay)

	// Determine axis range from all OHLC values.
	all := append(append([]float64{}, highVals[:n]...), lowVals[:n]...)
	minV, maxV := minMax(all)
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

	drawChartFrame(pdf, px, py, pw, ph, minV, maxV, opts.showMajorGridlines, opts.valueFormat)
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

// renderComboLike renders a combination bar+line chart with multiple series per type.
//
//nolint:cyclop,funlen // Combo charts necessarily enumerate both bar and line series in one pass.
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
		plotR = chartRectWithLegendMargin(r, opts.legendPosition)
	}
	px, py, pw, ph := chartPlotRect(plotR, opts.titleOverlay)

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

	drawChartFrame(pdf, px, py, pw, ph, minV, maxV, opts.showMajorGridlines, opts.valueFormat)

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

	// Render all bar series as grouped bars.
	for si, bs := range barSeries {
		colR, colG, colB := pieColor(si)
		pdf.SetFillColor(colR, colG, colB)
		for ci, v := range bs.Values {
			if ci >= nCats {
				break
			}
			subSlot := catSlot / float64(max(1, nBarSeries))
			bw := math.Max(4, subSlot*0.85)
			bx := px + catSlot*float64(ci) + subSlot*float64(si) + (subSlot-bw)/2
			valueY := py + ph*(maxV-v)/rangeV
			barTop := math.Min(zeroY, valueY)
			barH := math.Abs(zeroY - valueY)
			if barH < 0.5 {
				barH = 0.5
			}
			pdf.RectFromUpperLeftWithStyle(bx, barTop, bw, barH, "F")
		}
	}

	// Render all line series overlaid.
	for si, ls := range lineSeries {
		colR, colG, colB := pieColor(si + nBarSeries)
		pdf.SetStrokeColor(colR, colG, colB)
		n := min(len(ls.Values), nCats)
		if n < 1 {
			continue
		}
		for i := 0; i < n; i++ {
			x := px + catSlot*(float64(i)+0.5)
			y := py + ph - ((ls.Values[i]-minV)/rangeV)*(ph-4)
			if i > 0 {
				prevX := px + catSlot*(float64(i-1)+0.5)
				prevY := py + ph - ((ls.Values[i-1]-minV)/rangeV)*(ph-4)
				pdf.Line(prevX, prevY, x, y)
			}
			drawFilledCircle(pdf, x, y, 2.5, colR, colG, colB)
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
