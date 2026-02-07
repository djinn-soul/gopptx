package pptxxml

import (
	"fmt"
	"strings"
)

func pieChartPartXML(chart *ChartSpec) string {
	series := chartPieSeriesXML(chart)
	labels := chartPieDataLabelsXML(chart.ShowDataLabels)
	return chartPartEnvelope(
		chart.Title,
		chart.TitleOverlay,
		chart.ShowLegend,
		chart.LegendPosition,
		chart.LegendOverlay,
		fmt.Sprintf(`
<c:pieChart>
<c:varyColors val="1"/>%s
%s
</c:pieChart>`, series, labels),
	)
}

func chartPieDataLabelsXML(show bool) string {
	if !show {
		return ""
	}
	return `
<c:dLbls>
<c:showCatName val="1"/>
<c:showPercent val="1"/>
</c:dLbls>`
}

func chartPieSeriesXML(chart *ChartSpec) string {
	seriesName := chart.SeriesName
	if strings.TrimSpace(seriesName) == "" {
		seriesName = "Series 1"
	}

	var b strings.Builder
	b.WriteString(`
<c:ser>
<c:idx val="0"/>
<c:order val="0"/>
<c:tx><c:v>` + Escape(seriesName) + `</c:v></c:tx>
<c:cat><c:strLit>`)

	b.WriteString(fmt.Sprintf(`
<c:ptCount val="%d"/>`, len(chart.Categories)))
	for i, category := range chart.Categories {
		b.WriteString(fmt.Sprintf(`
<c:pt idx="%d"><c:v>%s</c:v></c:pt>`, i, Escape(category)))
	}
	b.WriteString(`
</c:strLit></c:cat>
<c:val><c:numLit>`)

	b.WriteString(fmt.Sprintf(`
<c:formatCode>General</c:formatCode>
<c:ptCount val="%d"/>`, len(chart.Values)))
	for i, value := range chart.Values {
		b.WriteString(fmt.Sprintf(`
<c:pt idx="%d"><c:v>%.6f</c:v></c:pt>`, i, value))
	}
	b.WriteString(`
</c:numLit></c:val>
</c:ser>`)

	return b.String()
}
