package pptxxml

import (
	"strconv"
	"strings"
)

func comboChartPartXML(chart *ChartSpec) string {
	barSeries := comboSeriesXML(chart, chart.BarSeries, 0)
	lineSeries := comboSeriesXML(chart, chart.LineSeries, len(chart.BarSeries))
	labels := chartDataLabelsXML(chart)
	var plot strings.Builder
	plot.WriteString(`
<c:barChart>
<c:barDir val="col"/>
<c:grouping val="clustered"/>`)
	plot.WriteString(barSeries)
	plot.WriteString(`
<c:axId val="48650112"/>
<c:axId val="48672768"/>
</c:barChart>
<c:lineChart>
<c:grouping val="standard"/>`)
	plot.WriteString(lineSeries)
	plot.WriteString(`
`)
	plot.WriteString(labels)
	plot.WriteString(`
<c:axId val="48650112"/>
<c:axId val="48672768"/>
</c:lineChart>
`)
	plot.WriteString(chartAxesXML(chart))

	return chartPartEnvelope(
		chart.Title,
		chart.TitleOverlay,
		chart.ShowLegend,
		chart.LegendPosition,
		chart.LegendOverlay,
		plot.String(),
	)
}

func comboSeriesXML(chart *ChartSpec, series []ChartSeries, start int) string {
	var out strings.Builder
	for i := range series {
		idx := start + i
		out.WriteString(`
<c:ser>
<c:idx val="`)
		out.WriteString(strconv.Itoa(idx))
		out.WriteString(`"/>
<c:order val="`)
		out.WriteString(strconv.Itoa(idx))
		out.WriteString(`"/>
<c:tx><c:v>`)
		out.WriteString(Escape(series[i].Name))
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
		out.WriteString(`
<c:ptCount val="`)
		out.WriteString(strconv.Itoa(len(series[i].Values)))
		out.WriteString(`"/>`)
		for j, value := range series[i].Values {
			out.WriteString(`
<c:pt idx="`)
			out.WriteString(strconv.Itoa(j))
			out.WriteString(`"><c:v>`)
			out.WriteString(strconv.FormatFloat(value, 'f', 6, 64))
			out.WriteString(`</c:v></c:pt>`)
		}
		out.WriteString(`
</c:numLit></c:val>
</c:ser>`)
	}
	return out.String()
}
