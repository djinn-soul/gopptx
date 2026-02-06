package pptxxml

import "fmt"

func areaChartPartXML(chart *ChartSpec) string {
	series := chartSeriesXML(chart)
	labels := chartDataLabelsXML(chart.ShowDataLabels)
	return chartPartEnvelope(chart.Title, chart.ShowLegend, chart.LegendPosition, fmt.Sprintf(`
<c:areaChart>
<c:grouping val="standard"/>
<c:varyColors val="0"/>%s
%s
<c:axId val="48650112"/>
<c:axId val="48672768"/>
</c:areaChart>
%s`, series, labels, chartAxesXML(chart)))
}
