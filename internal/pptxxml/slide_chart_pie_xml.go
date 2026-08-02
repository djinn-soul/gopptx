package pptxxml

import (
	"strconv"
	"strings"
)

// defaultSeriesName labels a series the caller left unnamed.
const defaultSeriesName = "Series 1"

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

func pie3DChartPartXML(chart *ChartSpec) string {
	xml := pieChartPartXML(chart)
	xml = strings.Replace(xml, "<c:pieChart>", "<c:pie3DChart>", 1)
	xml = strings.Replace(xml, "</c:pieChart>", "</c:pie3DChart>", 1)
	const view3D = `<c:view3D><c:rotX val="30"/><c:rotY val="0"/>` +
		`<c:rAngAx val="1"/><c:perspective val="30"/></c:view3D>`
	return strings.Replace(xml, "<c:plotArea>", view3D+"<c:plotArea>", 1)
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
		seriesName = defaultSeriesName
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
<c:formatCode>General</c:formatCode>`)
	writeNumericPoints(&b, chart.Values)
	b.WriteString(`
</c:numLit></c:val>
</c:ser>`)

	return b.String()
}
