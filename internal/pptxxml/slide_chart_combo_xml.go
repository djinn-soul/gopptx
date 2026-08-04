package pptxxml

import (
	"strconv"
	"strings"
)

func comboChartPartXML(chart *ChartSpec) string {
	barSeries := comboSeriesListXML(chart, chart.BarSeries, 0)
	primaryLine, secondaryLine := splitComboLineSeries(chart)
	labels := chartDataLabelsXML(chart)
	var plot strings.Builder
	plot.WriteString(`
<c:barChart>
<c:barDir val="col"/>
<c:grouping val="clustered"/>`)
	plot.WriteString(barSeries)
	plot.WriteString(primaryAxisIDRefsXML())
	plot.WriteString(`
</c:barChart>`)
	// A plot binds to exactly one axis pair, so line series that stay on the
	// primary scale and line series moved to the secondary one need separate
	// <c:lineChart> elements. Emitting one plot for all of them would drag the
	// unmarked series onto the secondary scale with the marked ones.
	if len(primaryLine) == 0 && len(secondaryLine) == 0 {
		// A combo with no line series still emitted an empty line plot before
		// the split; keep that shape rather than changing existing output.
		axisRefs := primaryAxisIDRefsXML()
		if chart.SecondaryAxis {
			axisRefs = secondaryAxisIDRefsXML()
		}
		plot.WriteString(comboLinePlotXML(chart, nil, labels, axisRefs))
	}
	if len(primaryLine) > 0 {
		plot.WriteString(comboLinePlotXML(chart, primaryLine, labels, primaryAxisIDRefsXML()))
	}
	if len(secondaryLine) > 0 {
		plot.WriteString(comboLinePlotXML(chart, secondaryLine, labels, secondaryAxisIDRefsXML()))
	}
	plot.WriteString("\n")
	plot.WriteString(chartAxesAndDataTableXML(chart))

	return chartPartEnvelope(
		chart.Title,
		chart.TitleOverlay,
		chart.ShowLegend,
		chart.LegendPosition,
		chart.LegendOverlay,
		plot.String(),
	)
}

// comboLineSeriesRef is one line series with the c:idx/c:order it carries
// across the whole chart, so splitting the series into two plots does not
// renumber them.
type comboLineSeriesRef struct {
	series ChartSeries
	index  int
}

// splitComboLineSeries divides the line series by which value axis they draw
// against. Without a secondary axis every series stays primary; with a
// chart-wide secondary axis and no per-series marks, every line series moves,
// which is what a caller enabling the axis alone asks for.
func splitComboLineSeries(chart *ChartSpec) ([]comboLineSeriesRef, []comboLineSeriesRef) {
	var primary, secondary []comboLineSeriesRef
	marked := false
	for i := range chart.LineSeries {
		if chart.LineSeries[i].SecondaryAxis {
			marked = true
			break
		}
	}
	start := len(chart.BarSeries)
	for i := range chart.LineSeries {
		ref := comboLineSeriesRef{series: chart.LineSeries[i], index: start + i}
		onSecondary := chart.SecondaryAxis &&
			(chart.LineSeries[i].SecondaryAxis || !marked)
		if onSecondary {
			secondary = append(secondary, ref)
		} else {
			primary = append(primary, ref)
		}
	}
	return primary, secondary
}

func comboLinePlotXML(
	chart *ChartSpec,
	series []comboLineSeriesRef,
	labels string,
	axisRefs string,
) string {
	var plot strings.Builder
	plot.WriteString(`
<c:lineChart>
<c:grouping val="standard"/>`)
	for _, ref := range series {
		plot.WriteString(comboSeriesXML(chart, ref.series, ref.index))
	}
	plot.WriteString(`
`)
	plot.WriteString(labels)
	plot.WriteString(axisRefs)
	plot.WriteString(`
</c:lineChart>`)
	return plot.String()
}

func comboSeriesListXML(chart *ChartSpec, series []ChartSeries, start int) string {
	var out strings.Builder
	for i, s := range series {
		out.WriteString(comboSeriesXML(chart, s, start+i))
	}
	return out.String()
}

func comboSeriesXML(chart *ChartSpec, series ChartSeries, idx int) string {
	var out strings.Builder
	out.WriteString(`
<c:ser>
<c:idx val="`)
	out.WriteString(strconv.Itoa(idx))
	out.WriteString(`"/>
<c:order val="`)
	out.WriteString(strconv.Itoa(idx))
	out.WriteString(`"/>
<c:tx><c:v>`)
	out.WriteString(Escape(series.Name))
	out.WriteString(`</c:v></c:tx>
<c:cat><c:strLit>
<c:ptCount val="`)
	out.WriteString(strconv.Itoa(len(chart.Categories)))
	out.WriteString(`"/>`)
	for j, category := range chart.Categories {
		out.WriteString(`
<c:pt idx="`)
		out.WriteString(strconv.Itoa(j))
		out.WriteString(`"><c:v>`)
		out.WriteString(Escape(category))
		out.WriteString(`</c:v></c:pt>`)
	}
	out.WriteString(`
</c:strLit></c:cat>
<c:val><c:numLit>
<c:formatCode>General</c:formatCode>`)
	writeNumericPoints(&out, series.Values)
	out.WriteString(`
</c:numLit></c:val>
</c:ser>`)
	return out.String()
}
