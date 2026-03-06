package pptxxml

import (
	"strconv"
	"strings"
)

func bubbleChartPartXML(chart *ChartSpec) string {
	series := chartBubbleSeriesXML(chart)
	labels := chartDataLabelsXML(chart)
	var plot strings.Builder
	plot.WriteString(`
<c:bubbleChart>
<c:varyColors val="0"/>`)
	plot.WriteString(series)
	plot.WriteString(`
`)
	plot.WriteString(labels)
	plot.WriteString(`
<c:bubbleScale val="`)
	plot.WriteString(strconv.Itoa(chart.BubbleScale))
	plot.WriteString(`"/>
<c:axId val="48650112"/>
<c:axId val="48672768"/>
</c:bubbleChart>
`)
	plot.WriteString(scatterAxesXML(chart))
	return chartPartEnvelope(
		chart.Title,
		chart.TitleOverlay,
		chart.ShowLegend,
		chart.LegendPosition,
		chart.LegendOverlay,
		plot.String(),
	)
}

func chartBubbleSeriesXML(chart *ChartSpec) string {
	var b strings.Builder
	b.WriteString(`
<c:ser>
<c:idx val="0"/>
<c:order val="0"/>
<c:tx><c:v>` + Escape(chart.SeriesName) + `</c:v></c:tx>
<c:spPr><a:solidFill><a:srgbClr val="` + Escape(chart.Color) + `"/></a:solidFill></c:spPr>
<c:xVal><c:numLit>`)

	b.WriteString(`
<c:formatCode>General</c:formatCode>
<c:ptCount val="`)
	b.WriteString(strconv.Itoa(len(chart.XValues)))
	b.WriteString(`"/>`)
	for i, value := range chart.XValues {
		b.WriteString(`
<c:pt idx="`)
		b.WriteString(strconv.Itoa(i))
		b.WriteString(`"><c:v>`)
		b.WriteString(strconv.FormatFloat(value, 'f', 6, 64))
		b.WriteString(`</c:v></c:pt>`)
	}
	b.WriteString(`
</c:numLit></c:xVal>
<c:yVal><c:numLit>`)

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
</c:numLit></c:yVal>
<c:bubbleSize><c:numLit>`)

	b.WriteString(`
<c:formatCode>General</c:formatCode>
<c:ptCount val="`)
	b.WriteString(strconv.Itoa(len(chart.BubbleSizes)))
	b.WriteString(`"/>`)
	for i, value := range chart.BubbleSizes {
		b.WriteString(`
<c:pt idx="`)
		b.WriteString(strconv.Itoa(i))
		b.WriteString(`"><c:v>`)
		b.WriteString(strconv.FormatFloat(value, 'f', 6, 64))
		b.WriteString(`</c:v></c:pt>`)
	}
	b.WriteString(`
</c:numLit></c:bubbleSize>
</c:ser>`)
	return b.String()
}
