package pptxxml

import "fmt"

func doughnutChartPartXML(chart *ChartSpec) string {
	series := chartPieSeriesXML(chart)
	labels := chartPieDataLabelsXML(chart.ShowDataLabels)
	return chartPartEnvelope(
		chart.Title,
		chart.TitleOverlay,
		chart.ShowLegend,
		chart.LegendPosition,
		chart.LegendOverlay,
		fmt.Sprintf(`
<c:doughnutChart>
<c:varyColors val="1"/>%s
%s
<c:holeSize val="%d"/>
</c:doughnutChart>`, series, labels, chart.HoleSize),
	)
}
