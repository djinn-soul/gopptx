package pptxxml

import "fmt"

func radarChartPartXML(chart *ChartSpec) string {
	series := chartSeriesXML(chart)
	labels := chartDataLabelsXML(chart)
	return chartPartEnvelope(
		chart.Title,
		chart.TitleOverlay,
		chart.ShowLegend,
		chart.LegendPosition,
		chart.LegendOverlay,
		fmt.Sprintf(`
<c:radarChart>
<c:radarStyle val="%s"/>%s
%s
<c:axId val="48650112"/>
<c:axId val="48672768"/>
</c:radarChart>
%s`, Escape(chart.RadarStyle), series, labels, chartAxesXML(chart)),
	)
}
