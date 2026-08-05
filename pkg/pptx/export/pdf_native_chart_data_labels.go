//nolint:mnd // Data-label placement uses tuned drawing offsets in points.
package export

import (
	"math"
	"strconv"
	"strings"

	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

// Data labels in PowerPoint are more than the bare number: c:dLbls decides which
// of series name, category name, value and percentage appear, in that order,
// separated by commas, and c:dLblPos decides where the label sits relative to
// its bar or point. This file turns those settings into label text and an
// anchor offset; the renderers supply the geometry.

const (
	// Gap between a bar or point and a label drawn outside it.
	dataLabelOutsideGapPt = 5.0
	// Inset used when the label is drawn inside the bar.
	dataLabelInsideGapPt = 2.0

	dataLabelSeparator = ", "

	percentScalePercent = 100.0

	// chartGeneralNumberFormat is Excel's unformatted number format, which
	// prints a value with as few digits as it needs.
	chartGeneralNumberFormat = "General"

	// Legend key swatch drawn in front of a label when c:showLegendKey is set.
	dataLabelLegendKeySizePt = 5.0
	dataLabelLegendKeyGapPt  = 2.5

	// Horizontal clearance for a label placed left or right of a data point.
	dataLabelPointSideGapPt = 6.0

	// A label never wraps narrower than this, so a chart with many categories
	// does not break every label onto one character per line.
	minDataLabelWrapWidthPt = 36.0
)

// chartDataLabel describes one label to draw.
type chartDataLabel struct {
	category   string
	seriesName string
	value      float64
	// total is the sum used for a percentage label. Zero suppresses it.
	total float64
}

// chartDataLabelText composes the label text for one data point. With no custom
// c:dLbls the label is the value alone, which is what PowerPoint shows when a
// chart simply switches data labels on.
func chartDataLabelText(opts chartSeriesOpts, label chartDataLabel) string {
	settings := opts.dataLabels
	if !settings.UseCustom {
		// Pie and doughnut charts ask for the category name through their own
		// flag, which predates the full settings struct.
		if opts.showCatName && label.category != "" {
			return label.category + dataLabelSeparator + formatChartValue(label.value, opts.valueFormat)
		}
		return formatChartValue(label.value, opts.valueFormat)
	}

	parts := make([]string, 0, 4)
	if settings.ShowSeriesName && label.seriesName != "" {
		parts = append(parts, label.seriesName)
	}
	if settings.ShowCategory && label.category != "" {
		parts = append(parts, label.category)
	}
	if settings.ShowValue {
		parts = append(parts, formatChartValue(label.value, opts.valueFormat))
	}
	if settings.ShowPercent && label.total != 0 {
		parts = append(parts, formatChartPercent(label.value/label.total*percentScalePercent))
	}
	return strings.Join(parts, dataLabelSeparator)
}

// formatChartValue renders a data-point value under an Excel-style number
// format. It covers the formats PowerPoint's own chart UI offers: a fixed number
// of decimals, percentages, currency and thousands grouping.
func formatChartValue(v float64, format string) string {
	trimmed := strings.TrimSpace(format)
	if trimmed == "" || trimmed == chartGeneralNumberFormat {
		// General prints the value as written, so the shortest representation
		// that round-trips is the right one: rounding to a fixed two decimals
		// here turned a 12.345 data point into 12.35.
		return strconv.FormatFloat(v, 'f', -1, 64)
	}

	decimals := formatDecimalPlaces(trimmed)
	scaled := v
	if strings.Contains(trimmed, "%") {
		scaled = v * percentScalePercent
	}
	out := strconv.FormatFloat(scaled, 'f', decimals, 64)
	if strings.Contains(trimmed, "#,#") || strings.Contains(trimmed, "#,0") {
		out = groupThousands(out)
	}
	if strings.HasPrefix(trimmed, "$") {
		out = "$" + out
	}
	if strings.Contains(trimmed, "%") {
		out += "%"
	}
	return out
}

func formatChartPercent(pct float64) string {
	return trimTrailingZeros(strconv.FormatFloat(pct, 'f', 1, 64)) + "%"
}

// formatDecimalPlaces counts the digit placeholders after the decimal point in
// an Excel number format, e.g. "0.00%" -> 2.
func formatDecimalPlaces(format string) int {
	_, decimals, hasDecimals := strings.Cut(format, ".")
	if !hasDecimals {
		return 0
	}
	count := 0
	for _, r := range decimals {
		if r != '0' && r != '#' {
			break
		}
		count++
	}
	return count
}

// groupThousands inserts thousands separators into the integer part of an
// already-formatted decimal string.
func groupThousands(value string) string {
	sign := ""
	if strings.HasPrefix(value, "-") {
		sign, value = "-", value[1:]
	}
	intPart, frac, hasFrac := strings.Cut(value, ".")
	var b strings.Builder
	for i, digit := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(digit)
	}
	out := sign + b.String()
	if hasFrac {
		out += "." + frac
	}
	return out
}

func trimTrailingZeros(value string) string {
	if !strings.Contains(value, ".") {
		return value
	}
	value = strings.TrimRight(value, "0")
	return strings.TrimSuffix(value, ".")
}

// barLabelY places a vertical bar's label according to c:dLblPos. barTop and
// barBottom bound the drawn bar; the returned y is the label's anchor line.
// PowerPoint's default for a clustered bar is outEnd, which is what an unset
// position falls back to.
func barLabelY(opts chartSeriesOpts, barTop, barBottom float64, negative bool) float64 {
	outsideEnd, insideEnd := barTop-dataLabelOutsideGapPt, barTop+dataLabelInsideGapPt
	insideBase := barBottom - dataLabelOutsideGapPt
	if negative {
		outsideEnd, insideEnd = barBottom+dataLabelInsideGapPt, barBottom-dataLabelOutsideGapPt
		insideBase = barTop + dataLabelInsideGapPt
	}
	switch dataLabelPosition(opts) {
	case charts.DataLabelPositionCenter:
		return (barTop + barBottom) / 2
	case charts.DataLabelPositionInsideEnd:
		return insideEnd
	case charts.DataLabelPositionInsideBase:
		return insideBase
	default:
		return outsideEnd
	}
}

// dataLabelPosition returns the requested c:dLblPos, or "" when the chart never
// set one.
func dataLabelPosition(opts chartSeriesOpts) string {
	if !opts.dataLabels.UseCustom {
		return ""
	}
	return strings.TrimSpace(opts.dataLabels.Position)
}

// chartDataLabelAt builds the label for the point at index i from the category
// list and total the renderer stored on opts.
func chartDataLabelAt(opts chartSeriesOpts, i int, value float64) chartDataLabel {
	label := chartDataLabel{
		seriesName: opts.seriesName,
		value:      value,
		total:      opts.valueTotal,
	}
	if i >= 0 && i < len(opts.categories) {
		label.category = opts.categories[i]
	}
	return label
}

// withChartLabelData returns opts carrying the category list and value total the
// data-label helpers need.
func withChartLabelData(opts chartSeriesOpts, categories []string, values []float64) chartSeriesOpts {
	opts.categories = categories
	opts.valueTotal = sumChartValues(values)
	return opts
}

// withChartPlotArea records the solved plot rect, which data labels need to wrap
// against and to stay inside. It is known only after the layout is solved, so it
// is attached separately from the series data.
func withChartPlotArea(opts chartSeriesOpts, x, y, w, h float64) chartSeriesOpts {
	opts.plot = chartRect{x: x, y: y, w: w, h: h}
	opts.hasPlot = true
	// PowerPoint wraps a data label within its category's share of the plot.
	slots := max(len(opts.categories), 1)
	opts.labelWrapWidth = math.Max(w/float64(slots), minDataLabelWrapWidthPt)
	return opts
}

// pieSliceLabelText is the label on one pie or doughnut wedge. A chart that
// states its own c:dLbls gets the full composition; without one the wedge keeps
// PowerPoint's pie defaults, which show the category name when the chart asked
// for it and the slice's share of the whole otherwise.
func pieSliceLabelText(
	opts chartSeriesOpts,
	categories []string,
	i int,
	value, total, fraction float64,
) string {
	if opts.dataLabels.UseCustom {
		label := chartDataLabelAt(withChartLabelData(opts, categories, nil), i, value)
		label.total = total
		return chartDataLabelText(opts, label)
	}
	if opts.showCatName && i < len(categories) && categories[i] != "" {
		return categories[i]
	}
	return formatChartPercent(fraction * percentScalePercent)
}

// drawChartDataLabel draws one composed data label centred on cx at labelY. The
// label wraps when c:dLbls asked it to, is preceded by the series' legend key
// when c:showLegendKey is set, and is nudged back inside the plot area rather
// than being allowed to run off it.
func drawChartDataLabel(pdf *gopdf.GoPdf, opts chartSeriesOpts, label chartDataLabel, cx, labelY float64) {
	value := chartDataLabelText(opts, label)
	if value == "" {
		return
	}
	lines := dataLabelLines(pdf, opts, value)
	lineHeight := pdfLineHeight(chartLabelFontSize)
	keyWidth := dataLabelKeyWidth(opts)

	textWidth := 0.0
	for _, line := range lines {
		textWidth = math.Max(textWidth, chartLabelWidth(pdf, line, chartLabelFontSize))
	}
	boxW := textWidth + keyWidth
	boxH := lineHeight * float64(len(lines))
	boxX, boxY := clampDataLabelToPlot(opts, cx-boxW/2, labelY, boxW, boxH)

	if keyWidth > 0 {
		drawDataLabelLegendKey(pdf, opts, boxX, boxY+lineHeight/2)
	}
	pdf.SetTextColor(60, 60, 60)
	for i, line := range lines {
		drawChartLabel(
			pdf, line,
			boxX+keyWidth, boxY+lineHeight*float64(i)+lineHeight/2,
			chartLabelFontSize, chartTextLeft,
		)
	}
	pdf.SetTextColor(0, 0, 0)
}

// dataLabelLines splits a label for drawing. PowerPoint wraps data-label text by
// default and only stops when c:dLbls says wrap="0"; without a known plot area
// there is nothing to wrap against, so the label stays on one line.
func dataLabelLines(pdf *gopdf.GoPdf, opts chartSeriesOpts, value string) []string {
	wrapWidth := opts.labelWrapWidth
	if wrapWidth <= 0 || !dataLabelWraps(opts) {
		return []string{value}
	}
	setPDFTextFontWithHint(pdf, chartLabelFontSize, false, false, "")
	return wrapPDFTextWithMetrics(pdf, value, wrapWidth)
}

func dataLabelWraps(opts chartSeriesOpts) bool {
	if opts.dataLabels.WordWrap == nil {
		return true
	}
	return *opts.dataLabels.WordWrap
}

func dataLabelKeyWidth(opts chartSeriesOpts) float64 {
	if !opts.dataLabels.UseCustom || !opts.dataLabels.ShowLegendKey {
		return 0
	}
	return dataLabelLegendKeySizePt + dataLabelLegendKeyGapPt
}

// drawDataLabelLegendKey draws the small colour swatch that c:showLegendKey puts
// in front of the label text.
func drawDataLabelLegendKey(pdf *gopdf.GoPdf, opts chartSeriesOpts, x, middleY float64) {
	r, g, b := uint8(79), uint8(129), uint8(189)
	if opts.color != "" {
		r, g, b = hexToRGB(opts.color)
	}
	pdf.SetFillColor(r, g, b)
	pdf.RectFromUpperLeftWithStyle(
		x, middleY-dataLabelLegendKeySizePt/2,
		dataLabelLegendKeySizePt, dataLabelLegendKeySizePt, "F",
	)
}

// clampDataLabelToPlot keeps a label box inside the chart's plot area when one
// is known. A wide label on the first or last category used to hang off the
// side of the chart.
func clampDataLabelToPlot(opts chartSeriesOpts, x, y, w, h float64) (float64, float64) {
	if !opts.hasPlot {
		return x, y
	}
	plot := opts.plot
	if x < plot.x {
		x = plot.x
	}
	if right := plot.x + plot.w; x+w > right {
		x = math.Max(plot.x, right-w)
	}
	if y < plot.y {
		y = plot.y
	}
	if bottom := plot.y + plot.h; y+h > bottom {
		y = math.Max(plot.y, bottom-h)
	}
	return x, y
}

// dataLabelPointAnchor places a label around a data point (a line or scatter
// marker) according to c:dLblPos. The returned point is where the label's centre
// column starts; "t" is PowerPoint's default for a line series.
func dataLabelPointAnchor(opts chartSeriesOpts, x, y float64) (float64, float64) {
	switch dataLabelPosition(opts) {
	case charts.DataLabelPositionCenter:
		return x, y - pdfLineHeight(chartLabelFontSize)/2
	case charts.DataLabelPositionBottom:
		return x, y + dataLabelOutsideGapPt
	case charts.DataLabelPositionLeft:
		return x - dataLabelPointSideGapPt, y - pdfLineHeight(chartLabelFontSize)/2
	case charts.DataLabelPositionRight:
		return x + dataLabelPointSideGapPt, y - pdfLineHeight(chartLabelFontSize)/2
	default: // top, bestFit and anything a point series cannot honour
		return x, y - dataLabelOutsideGapPt - pdfLineHeight(chartLabelFontSize)
	}
}

// sumChartValues is the total a percentage data label is measured against.
// PowerPoint uses the sum of the series' absolute values.
func sumChartValues(values []float64) float64 {
	total := 0.0
	for _, v := range values {
		total += math.Abs(v)
	}
	return total
}
