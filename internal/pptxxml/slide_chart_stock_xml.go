package pptxxml

import (
	"fmt"
	"strings"
)

func stockChartPartXML(chart *ChartSpec) string {
	seriesParts := make([]string, 0, 4)
	index := 0
	if chart.Kind == ChartKindStockOHLC {
		seriesParts = append(seriesParts, chartSeriesWithValues(chart, "Open", chart.OpenValues, index))
		index++
	}
	seriesParts = append(seriesParts, chartSeriesWithValues(chart, "High", chart.HighValues, index))
	index++
	seriesParts = append(seriesParts, chartSeriesWithValues(chart, "Low", chart.LowValues, index))
	index++
	seriesParts = append(seriesParts, chartSeriesWithValues(chart, "Close", chart.CloseValues, index))
	labels := chartDataLabelsXML(chart.ShowDataLabels)

	return chartPartEnvelope(chart.Title, chart.ShowLegend, chart.LegendPosition, fmt.Sprintf(`
<c:stockChart>%s
%s
<c:axId val="48650112"/>
<c:axId val="48672768"/>
</c:stockChart>
%s`, strings.Join(seriesParts, ""), labels, chartAxesXML(chart)))
}

func chartSeriesWithValues(chart *ChartSpec, seriesName string, values []float64, idx int) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`
<c:ser>
<c:idx val="%d"/>
<c:order val="%d"/>
<c:tx><c:v>%s</c:v></c:tx>
<c:cat><c:strLit>
<c:ptCount val="%d"/>`, idx, idx, Escape(seriesName), len(chart.Categories)))
	for i, category := range chart.Categories {
		b.WriteString(fmt.Sprintf(`
<c:pt idx="%d"><c:v>%s</c:v></c:pt>`, i, Escape(category)))
	}
	b.WriteString(`
</c:strLit></c:cat>
<c:val><c:numLit>
<c:formatCode>General</c:formatCode>`)
	b.WriteString(fmt.Sprintf(`
<c:ptCount val="%d"/>`, len(values)))
	for i, value := range values {
		b.WriteString(fmt.Sprintf(`
<c:pt idx="%d"><c:v>%.6f</c:v></c:pt>`, i, value))
	}
	b.WriteString(`
</c:numLit></c:val>
</c:ser>`)
	return b.String()
}
