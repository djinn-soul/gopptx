package pptxxml

import "fmt"

func radarChartPartXML(chart *ChartSpec) string {
	series := chartSeriesXML(chart)
	labels := chartDataLabelsXML(chart.ShowDataLabels)
	return chartPartEnvelope(chart.Title, chart.ShowLegend, chart.LegendPosition, fmt.Sprintf(`
<c:radarChart>
<c:radarStyle val="%s"/>%s
%s
<c:axId val="48650112"/>
<c:axId val="48672768"/>
</c:radarChart>
%s`, Escape(chart.RadarStyle), series, labels, chartAxesXML(chart)))
}
