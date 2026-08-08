//nolint:mnd // Chart SVG uses a fixed internal coordinate space and fixed paddings.
package export

import (
	"fmt"
	"html"
	"io"
	"math"
	"strings"
)

// Charts are drawn into their own SVG in a fixed 640x360 coordinate space and
// scaled by the page's CSS, so a chart reads the same whatever size the slide
// gave it. Kinds with no simple plot are written out as their data instead —
// see htmlChartTable — which is also the accessible form of any chart.
const (
	chartSVGWidth    = 640.0
	chartSVGHeight   = 360.0
	chartPadLeft     = 48.0
	chartPadRight    = 16.0
	chartPadTop      = 32.0
	chartPadBottom   = 44.0
	chartSeriesGap   = 4.0
	chartLabelSizePt = 11.0
)

// chartSeriesPalette is the accent sequence a multi-series chart cycles.
var chartSeriesPalette = [...]string{ //nolint:gochecknoglobals // Immutable palette shared by every chart.
	"#4472C4", "#ED7D31", "#A5A5A5", "#FFC000", "#5B9BD5", "#70AD47",
}

func chartSeriesColor(index int) string {
	return chartSeriesPalette[index%len(chartSeriesPalette)]
}

func renderChartsToWriter(w io.Writer, slideCharts []htmlChart) error {
	for _, chart := range slideCharts {
		if err := renderChartToWriter(w, chart); err != nil {
			return err
		}
	}
	return nil
}

func renderChartToWriter(w io.Writer, chart htmlChart) error {
	if !chart.hasData() {
		return nil
	}
	if err := writeString(w, "<figure class=\"slide-chart\">\n"); err != nil {
		return err
	}
	if chart.Title != "" {
		if _, err := fmt.Fprintf(w, "<figcaption>%s</figcaption>\n", html.EscapeString(chart.Title)); err != nil {
			return err
		}
	}
	body := chartSVG(chart)
	if body == "" {
		body = chartDataTable(chart)
	}
	if err := writeString(w, body); err != nil {
		return err
	}
	return writeString(w, "</figure>\n")
}

// chartSVG draws the chart, or returns empty for a kind with no simple plot.
func chartSVG(chart htmlChart) string {
	switch chart.Kind {
	case htmlChartBar, htmlChartBarStacked:
		return wrapChartSVG(chart, barChartBody(chart))
	case htmlChartLine, htmlChartArea:
		return wrapChartSVG(chart, lineChartBody(chart))
	case htmlChartPie:
		return wrapChartSVG(chart, pieChartBody(chart))
	case htmlChartTable:
		return ""
	default:
		return ""
	}
}

func wrapChartSVG(chart htmlChart, body string) string {
	if body == "" {
		return ""
	}
	var sb strings.Builder
	fmt.Fprintf(&sb,
		`<svg class="chart-svg" viewBox="0 0 %.0f %.0f" preserveAspectRatio="xMidYMid meet"`+
			` role="img" xmlns="http://www.w3.org/2000/svg">`+"\n",
		chartSVGWidth, chartSVGHeight)
	label := chart.Title
	if label == "" {
		label = "chart"
	}
	fmt.Fprintf(&sb, "<title>%s</title>\n", html.EscapeString(label))
	sb.WriteString(body)
	sb.WriteString("</svg>\n")
	return sb.String()
}

// plotArea is the rectangle inside the axis labels.
func plotArea() (float64, float64, float64, float64) {
	return chartPadLeft,
		chartPadTop,
		chartSVGWidth - chartPadLeft - chartPadRight,
		chartSVGHeight - chartPadTop - chartPadBottom
}

func barChartBody(chart htmlChart) string {
	count := chart.categoryCount()
	maximum := chart.maxValue()
	if count == 0 || maximum <= 0 {
		return ""
	}
	x, y, w, h := plotArea()
	var sb strings.Builder
	sb.WriteString(chartAxes(x, y, w, h))

	slot := w / float64(count)
	stacked := chart.Kind == htmlChartBarStacked
	for index := range count {
		if stacked {
			sb.WriteString(stackedColumn(chart, index, x+float64(index)*slot, y, slot, h, maximum))
		} else {
			sb.WriteString(groupedColumn(chart, index, x+float64(index)*slot, y, slot, h, maximum))
		}
		sb.WriteString(categoryLabelText(chart.categoryLabel(index), x+float64(index)*slot+slot/2, y+h+16))
	}
	sb.WriteString(chartLegend(chart))
	return sb.String()
}

func groupedColumn(chart htmlChart, index int, slotX, plotY, slot, plotH, maximum float64) string {
	var sb strings.Builder
	bars := max(len(chart.Series), 1)
	barW := (slot - chartSeriesGap*float64(bars+1)) / float64(bars)
	if barW <= 0 {
		barW = slot / float64(bars)
	}
	for s, series := range chart.Series {
		if index >= len(series.Values) {
			continue
		}
		value := series.Values[index]
		barH := plotH * (value / maximum)
		if barH < 0 {
			barH = 0
		}
		left := slotX + chartSeriesGap*float64(s+1) + float64(s)*barW
		fmt.Fprintf(&sb, `<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="%s"/>`+"\n",
			left, plotY+plotH-barH, barW, barH, chartSeriesColor(s))
	}
	return sb.String()
}

func stackedColumn(chart htmlChart, index int, slotX, plotY, slot, plotH, maximum float64) string {
	var sb strings.Builder
	barW := slot * 0.6
	left := slotX + (slot-barW)/2
	bottom := plotY + plotH
	for s, series := range chart.Series {
		if index >= len(series.Values) {
			continue
		}
		barH := plotH * (series.Values[index] / maximum)
		if barH < 0 {
			barH = 0
		}
		bottom -= barH
		fmt.Fprintf(&sb, `<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="%s"/>`+"\n",
			left, bottom, barW, barH, chartSeriesColor(s))
	}
	return sb.String()
}

func lineChartBody(chart htmlChart) string {
	count := chart.categoryCount()
	maximum := chart.maxValue()
	if count == 0 || maximum <= 0 {
		return ""
	}
	x, y, w, h := plotArea()
	var sb strings.Builder
	sb.WriteString(chartAxes(x, y, w, h))

	step := w
	if count > 1 {
		step = w / float64(count-1)
	}
	pointX := func(index int) float64 {
		if count == 1 {
			return x + w/2
		}
		return x + float64(index)*step
	}
	for s, series := range chart.Series {
		points := make([]string, 0, len(series.Values))
		for index, value := range series.Values {
			points = append(points, fmt.Sprintf("%.2f,%.2f", pointX(index), y+h-h*(value/maximum)))
		}
		if len(points) == 0 {
			continue
		}
		joined := strings.Join(points, " ")
		if chart.Kind == htmlChartArea {
			fmt.Fprintf(&sb, `<polygon points="%.2f,%.2f %s %.2f,%.2f" fill="%s" fill-opacity="0.35"/>`+"\n",
				pointX(0), y+h, joined, pointX(len(series.Values)-1), y+h, chartSeriesColor(s))
		}
		fmt.Fprintf(&sb, `<polyline points="%s" fill="none" stroke="%s" stroke-width="2"/>`+"\n",
			joined, chartSeriesColor(s))
	}
	for index := range count {
		sb.WriteString(categoryLabelText(chart.categoryLabel(index), pointX(index), y+h+16))
	}
	sb.WriteString(chartLegend(chart))
	return sb.String()
}

func pieChartBody(chart htmlChart) string {
	if len(chart.Series) == 0 {
		return ""
	}
	values := chart.Series[0].Values
	total := 0.0
	for _, value := range values {
		if value > 0 {
			total += value
		}
	}
	if total <= 0 {
		return ""
	}
	cx, cy := chartSVGWidth/2, chartSVGHeight/2
	radius := math.Min(chartSVGWidth, chartSVGHeight)/2 - chartPadTop

	var sb strings.Builder
	angle := -math.Pi / 2
	for i, value := range values {
		if value <= 0 {
			continue
		}
		sweep := 2 * math.Pi * (value / total)
		sb.WriteString(pieSlicePath(cx, cy, radius, angle, sweep, chartSeriesColor(i)))
		angle += sweep
	}
	sb.WriteString(chartLegendFromCategories(chart))
	return sb.String()
}

// pieSlicePath draws one slice. A slice covering the whole circle is drawn as a
// circle, because an arc from a point back to itself has no path.
func pieSlicePath(cx, cy, radius, start, sweep float64, color string) string {
	if sweep >= 2*math.Pi-0.0001 {
		return fmt.Sprintf(`<circle cx="%.2f" cy="%.2f" r="%.2f" fill="%s"/>`+"\n", cx, cy, radius, color)
	}
	x1, y1 := cx+radius*math.Cos(start), cy+radius*math.Sin(start)
	x2, y2 := cx+radius*math.Cos(start+sweep), cy+radius*math.Sin(start+sweep)
	largeArc := 0
	if sweep > math.Pi {
		largeArc = 1
	}
	return fmt.Sprintf(
		`<path d="M %.2f %.2f L %.2f %.2f A %.2f %.2f 0 %d 1 %.2f %.2f Z" fill="%s"/>`+"\n",
		cx, cy, x1, y1, radius, radius, largeArc, x2, y2, color,
	)
}

func chartAxes(x, y, w, h float64) string {
	return fmt.Sprintf(
		`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#999" stroke-width="1"/>`+"\n"+
			`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#999" stroke-width="1"/>`+"\n",
		x, y, x, y+h,
		x, y+h, x+w, y+h,
	)
}

func categoryLabelText(label string, cx, y float64) string {
	if label == "" {
		return ""
	}
	return fmt.Sprintf(
		`<text x="%.2f" y="%.2f" font-size="%.0fpx" fill="#444" text-anchor="middle">%s</text>`+"\n",
		cx, y, chartLabelSizePt, html.EscapeString(label),
	)
}

// chartLegend names the series along the bottom, for a chart that has more than
// one to tell apart.
func chartLegend(chart htmlChart) string {
	if len(chart.Series) < 2 {
		return ""
	}
	names := make([]string, 0, len(chart.Series))
	for _, series := range chart.Series {
		names = append(names, series.Name)
	}
	return legendRow(names)
}

// chartLegendFromCategories names a pie's slices, which are its categories.
func chartLegendFromCategories(chart htmlChart) string {
	count := chart.categoryCount()
	if count == 0 {
		return ""
	}
	names := make([]string, 0, count)
	for index := range count {
		names = append(names, chart.categoryLabel(index))
	}
	return legendRow(names)
}

func legendRow(names []string) string {
	var sb strings.Builder
	x := chartPadLeft
	y := chartSVGHeight - 12
	for i, name := range names {
		if strings.TrimSpace(name) == "" {
			continue
		}
		fmt.Fprintf(&sb, `<rect x="%.2f" y="%.2f" width="10" height="10" fill="%s"/>`+"\n",
			x, y-9, chartSeriesColor(i))
		fmt.Fprintf(&sb,
			`<text x="%.2f" y="%.2f" font-size="%.0fpx" fill="#444">%s</text>`+"\n",
			x+14, y, chartLabelSizePt, html.EscapeString(name))
		x += 14 + float64(len(name))*6 + 18
	}
	return sb.String()
}

// chartDataTable writes a chart out as its numbers, for the kinds with no simple
// plot and for anything the plotter could not scale.
func chartDataTable(chart htmlChart) string {
	var sb strings.Builder
	sb.WriteString("<table class=\"chart-data\">\n<tr><th></th>")
	for _, series := range chart.Series {
		fmt.Fprintf(&sb, "<th>%s</th>", html.EscapeString(series.Name))
	}
	sb.WriteString("</tr>\n")
	for index := range chart.categoryCount() {
		fmt.Fprintf(&sb, "<tr><th>%s</th>", html.EscapeString(chart.categoryLabel(index)))
		for _, series := range chart.Series {
			cell := ""
			if index < len(series.Values) {
				cell = formatHTMLChartValue(series.Values[index])
			}
			fmt.Fprintf(&sb, "<td>%s</td>", cell)
		}
		sb.WriteString("</tr>\n")
	}
	sb.WriteString("</table>\n")
	return sb.String()
}
