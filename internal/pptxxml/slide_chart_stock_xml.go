package pptxxml

import (
	"strconv"
	"strings"
)

const chartSeriesCapacity = 4

func stockChartPartXML(chart *ChartSpec) string {
	seriesParts := make([]string, 0, chartSeriesCapacity)
	index := 0
	if chart.Kind == ChartKindStockOHLC {
		seriesParts = append(seriesParts, chartSeriesWithValues(chart, "Open", chart.OpenValues, index))
		index++
	}
	seriesParts = append(seriesParts, chartSeriesWithValues(chart, "High", chart.HighValues, index))
	index++
	seriesParts = append(seriesParts, chartSeriesWithValues(chart, "Low", chart.LowValues, index))
	index++
	seriesParts = append(seriesParts, chartSeriesWithValues(chart, "Close", chart.CloseValues, index))
	visuals := stockVisualsXML(chart)
	labels := chartDataLabelsXML(chart)
	var plot strings.Builder
	plot.WriteString(`
<c:stockChart>`)
	plot.WriteString(strings.Join(seriesParts, ""))
	plot.WriteString(`
`)
	plot.WriteString(visuals)
	plot.WriteString(`
`)
	plot.WriteString(labels)
	plot.WriteString(`
<c:axId val="48650112"/>
<c:axId val="48672768"/>
</c:stockChart>
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

func chartSeriesWithValues(chart *ChartSpec, seriesName string, values []float64, idx int) string {
	var b strings.Builder
	b.WriteString(`
<c:ser>
<c:idx val="`)
	b.WriteString(strconv.Itoa(idx))
	b.WriteString(`"/>
<c:order val="`)
	b.WriteString(strconv.Itoa(idx))
	b.WriteString(`"/>
<c:tx><c:v>`)
	b.WriteString(Escape(seriesName))
	b.WriteString(`</c:v></c:tx>
<c:spPr>
<a:ln><a:noFill/></a:ln>
</c:spPr>
<c:marker><c:symbol val="none"/></c:marker>
<c:cat><c:strLit>
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
<c:val><c:numLit>
<c:formatCode>General</c:formatCode>`)
	b.WriteString(`
<c:ptCount val="`)
	b.WriteString(strconv.Itoa(len(values)))
	b.WriteString(`"/>`)
	for i, value := range values {
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

func stockVisualsXML(chart *ChartSpec) string {
	var b strings.Builder
	b.WriteString(`
<c:hiLowLines/>`)
	if chart.Kind == ChartKindStockOHLC {
		b.WriteString(`
<c:upDownBars>
<c:gapWidth val="150"/>
<c:upBars><c:spPr><a:solidFill><a:srgbClr val="00B050"/></a:solidFill><a:ln><a:noFill/></a:ln></c:spPr></c:upBars>
<c:downBars><c:spPr><a:solidFill><a:srgbClr val="C0504D"/></a:solidFill><a:ln><a:noFill/></a:ln></c:spPr></c:downBars>
</c:upDownBars>`)
	}
	return b.String()
}
