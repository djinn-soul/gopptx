package pptxxml

import "fmt"

func doughnutChartPartXML(chart *ChartSpec) string {
	series := chartPieSeriesXML(chart)
	labels := chartPieDataLabelsXML(chart.ShowDataLabels)
	return chartPartEnvelope(chart.Title, chart.ShowLegend, chart.LegendPosition, fmt.Sprintf(`
<c:doughnutChart>
<c:varyColors val="1"/>
<c:holeSize val="%d"/>%s
%s
</c:doughnutChart>`, chart.HoleSize, series, labels))
}
