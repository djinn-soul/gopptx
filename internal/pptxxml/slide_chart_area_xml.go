package pptxxml

import "fmt"

func areaChartPartXML(chart *ChartSpec) string {
	series := chartSeriesXML(chart)
	labels := chartDataLabelsXML(chart.ShowDataLabels)
	return chartPartEnvelope(
		chart.Title,
		chart.TitleOverlay,
		chart.ShowLegend,
		chart.LegendPosition,
		chart.LegendOverlay,
		fmt.Sprintf(`
<c:areaChart>
<c:grouping val="%s"/>
<c:varyColors val="0"/>%s
%s
<c:axId val="48650112"/>
<c:axId val="48672768"/>
</c:areaChart>
%s`, Escape(chart.Grouping), series, labels, chartAxesXML(chart)),
	)
}
