package pptxxml

import (
	"fmt"
	"strings"
)

func comboChartPartXML(chart *ChartSpec) string {
	barSeries := comboSeriesXML(chart, chart.BarSeries, 0)
	lineSeries := comboSeriesXML(chart, chart.LineSeries, len(chart.BarSeries))
	labels := chartDataLabelsXML(chart.ShowDataLabels)

	return chartPartEnvelope(
		chart.Title,
		chart.TitleOverlay,
		chart.ShowLegend,
		chart.LegendPosition,
		chart.LegendOverlay,
		fmt.Sprintf(`
<c:barChart>
<c:barDir val="col"/>
<c:grouping val="clustered"/>%s
<c:axId val="48650112"/>
<c:axId val="48672768"/>
</c:barChart>
<c:lineChart>
<c:grouping val="standard"/>%s
%s
<c:axId val="48650112"/>
<c:axId val="48672768"/>
</c:lineChart>
%s`, barSeries, lineSeries, labels, chartAxesXML(chart)),
	)
}

func comboSeriesXML(chart *ChartSpec, series []ChartSeries, start int) string {
	var out strings.Builder
	for i := range series {
		idx := start + i
		out.WriteString(fmt.Sprintf(`
<c:ser>
<c:idx val="%d"/>
<c:order val="%d"/>
<c:tx><c:v>%s</c:v></c:tx>
<c:cat><c:strLit>
<c:ptCount val="%d"/>`, idx, idx, Escape(series[i].Name), len(chart.Categories)))
		for j, category := range chart.Categories {
			out.WriteString(fmt.Sprintf(`
<c:pt idx="%d"><c:v>%s</c:v></c:pt>`, j, Escape(category)))
		}
		out.WriteString(`
</c:strLit></c:cat>
<c:val><c:numLit>
<c:formatCode>General</c:formatCode>`)
		out.WriteString(fmt.Sprintf(`
<c:ptCount val="%d"/>`, len(series[i].Values)))
		for j, value := range series[i].Values {
			out.WriteString(fmt.Sprintf(`
<c:pt idx="%d"><c:v>%.6f</c:v></c:pt>`, j, value))
		}
		out.WriteString(`
</c:numLit></c:val>
</c:ser>`)
	}
	return out.String()
}
