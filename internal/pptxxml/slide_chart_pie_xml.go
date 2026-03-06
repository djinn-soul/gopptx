package pptxxml

import (
	"strconv"
	"strings"
)

func pieChartPartXML(chart *ChartSpec) string {
	series := chartPieSeriesXML(chart)
	labels := chartPieDataLabelsXML(chart)
	var plot strings.Builder
	plot.WriteString(`
<c:pieChart>
<c:varyColors val="1"/>`)
	plot.WriteString(series)
	plot.WriteString(`
`)
	plot.WriteString(labels)
	plot.WriteString(`
</c:pieChart>`)
	return chartPartEnvelope(
		chart.Title,
		chart.TitleOverlay,
		chart.ShowLegend,
		chart.LegendPosition,
		chart.LegendOverlay,
		plot.String(),
	)
}

func chartPieDataLabelsXML(chart *ChartSpec) string {
	return chartDataLabelsWithDefaults(chart, chartDataLabelDefaults{
		showCategory: true,
		showPercent:  true,
	})
}

func chartPieSeriesXML(chart *ChartSpec) string {
	seriesName := chart.SeriesName
	if strings.TrimSpace(seriesName) == "" {
		seriesName = "Series 1"
	}

	var b strings.Builder
	b.WriteString(`
<c:ser>
<c:idx val="0"/>
<c:order val="0"/>
<c:tx><c:v>` + Escape(seriesName) + `</c:v></c:tx>
<c:cat><c:strLit>`)

	b.WriteString(`
<c:ptCount val="`)
	b.WriteString(strconv.Itoa(len(chart.Categories)))
	b.WriteString(`"/>`)
	for i, category := range chart.Categories {
		b.WriteString(`
<c:pt idx="`)
		b.WriteString(strconv.Itoa(i))
		b.WriteString(`"><c:v>`)
		b.WriteString(Escape(category))
		b.WriteString(`</c:v></c:pt>`)
	}
	b.WriteString(`
</c:strLit></c:cat>
<c:val><c:numLit>`)

	b.WriteString(`
<c:formatCode>General</c:formatCode>
<c:ptCount val="`)
	b.WriteString(strconv.Itoa(len(chart.Values)))
	b.WriteString(`"/>`)
	for i, value := range chart.Values {
		b.WriteString(`
<c:pt idx="`)
		b.WriteString(strconv.Itoa(i))
		b.WriteString(`"><c:v>`)
		b.WriteString(strconv.FormatFloat(value, 'f', 6, 64))
		b.WriteString(`</c:v></c:pt>`)
	}
	b.WriteString(`
</c:numLit></c:val>
</c:ser>`)

	return b.String()
}
