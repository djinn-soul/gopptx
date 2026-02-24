package pptxxml

import (
	"strconv"
	"strings"
)

func scatterChartPartXML(chart *ChartSpec) string {
	series := chartScatterSeriesXML(chart)
	labels := chartDataLabelsXML(chart.ShowDataLabels)
	var plot strings.Builder
	plot.WriteString(`
<c:scatterChart>
<c:scatterStyle val="`)
	plot.WriteString(Escape(chart.ScatterStyle))
	plot.WriteString(`"/>`)
	plot.WriteString(series)
	plot.WriteString(`
`)
	plot.WriteString(labels)
	plot.WriteString(`
<c:axId val="48650112"/>
<c:axId val="48672768"/>
</c:scatterChart>
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

func chartScatterSeriesXML(chart *ChartSpec) string {
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
</c:ser>`)

	return b.String()
}

func scatterAxesXML(chart *ChartSpec) string {
	xAxisTitle := chartAxisTitleXML(chart.CategoryAxisTitle)
	yAxisTitle := chartAxisTitleXML(chart.ValueAxisTitle)
	yScaling := valueAxisScalingXML(chart.MinValue, chart.MaxValue)
	yFormat := chartValueFormatXML(chart.ValueFormat)
	crossBetween := normalizedValueAxisCrossBetween(chart.ValueAxisCrossBetween)
	majorGrid := ""
	if chart.ShowMajorGridlines {
		majorGrid = "<c:majorGridlines/>"
	}

	var b strings.Builder
	b.WriteString(`
<c:valAx>
<c:axId val="48650112"/>
<c:scaling><c:orientation val="minMax"/></c:scaling>
<c:delete val="0"/>
<c:axPos val="b"/>`)
	b.WriteString(xAxisTitle)
	b.WriteString(`
<c:numFmt formatCode="General" sourceLinked="1"/>
<c:tickLblPos val="nextTo"/>
<c:crossAx val="48672768"/>
<c:crosses val="autoZero"/>
</c:valAx>
<c:valAx>
<c:axId val="48672768"/>`)
	b.WriteString(yScaling)
	b.WriteString(`
<c:delete val="0"/>
<c:axPos val="l"/>`)
	b.WriteString(majorGrid)
	b.WriteString(yAxisTitle)
	b.WriteString(yFormat)
	b.WriteString(`
<c:tickLblPos val="nextTo"/>
<c:crossAx val="48650112"/>
<c:crosses val="autoZero"/>
<c:crossBetween val="`)
	b.WriteString(crossBetween)
	b.WriteString(`"/>
</c:valAx>`)
	return b.String()
}
