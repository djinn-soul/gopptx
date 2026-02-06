package pptxxml

import (
	"fmt"
	"strings"
)

func scatterChartPartXML(chart *ChartSpec) string {
	series := chartScatterSeriesXML(chart)
	labels := chartDataLabelsXML(chart.ShowDataLabels)
	return chartPartEnvelope(chart.Title, chart.ShowLegend, chart.LegendPosition, fmt.Sprintf(`
<c:scatterChart>
<c:scatterStyle val="%s"/>%s
%s
<c:axId val="48650112"/>
<c:axId val="48672768"/>
</c:scatterChart>
%s`, Escape(chart.ScatterStyle), series, labels, scatterAxesXML(chart)))
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

	b.WriteString(fmt.Sprintf(`
<c:formatCode>General</c:formatCode>
<c:ptCount val="%d"/>`, len(chart.XValues)))
	for i, value := range chart.XValues {
		b.WriteString(fmt.Sprintf(`
<c:pt idx="%d"><c:v>%.6f</c:v></c:pt>`, i, value))
	}
	b.WriteString(`
</c:numLit></c:xVal>
<c:yVal><c:numLit>`)

	b.WriteString(fmt.Sprintf(`
<c:formatCode>General</c:formatCode>
<c:ptCount val="%d"/>`, len(chart.Values)))
	for i, value := range chart.Values {
		b.WriteString(fmt.Sprintf(`
<c:pt idx="%d"><c:v>%.6f</c:v></c:pt>`, i, value))
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
	majorGrid := ""
	if chart.ShowMajorGridlines {
		majorGrid = "<c:majorGridlines/>"
	}

	return fmt.Sprintf(`
<c:valAx>
<c:axId val="48650112"/>
<c:scaling><c:orientation val="minMax"/></c:scaling>
<c:delete val="0"/>
<c:axPos val="b"/>
%s
<c:numFmt formatCode="General" sourceLinked="1"/>
<c:tickLblPos val="nextTo"/>
<c:crossAx val="48672768"/>
<c:crosses val="autoZero"/>
</c:valAx>
<c:valAx>
<c:axId val="48672768"/>
%s
<c:delete val="0"/>
<c:axPos val="l"/>
%s
%s
%s
<c:tickLblPos val="nextTo"/>
<c:crossAx val="48650112"/>
<c:crosses val="autoZero"/>
</c:valAx>`,
		xAxisTitle,
		yScaling,
		yAxisTitle,
		yFormat,
		majorGrid,
	)
}
