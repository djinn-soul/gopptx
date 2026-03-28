package pptxxml

import (
	"fmt"
	"strconv"
	"strings"
)

func radarChartPartXML(chart *ChartSpec) string {
	series := radarSeriesXML(chart)
	labels := chartDataLabelsXML(chart)
	return chartPartEnvelope(
		chart.Title,
		chart.TitleOverlay,
		chart.ShowLegend,
		chart.LegendPosition,
		chart.LegendOverlay,
		fmt.Sprintf(`
<c:radarChart>
<c:radarStyle val="%s"/>
<c:varyColors val="0"/>%s
%s
<c:axId val="48650112"/>
<c:axId val="48672768"/>
</c:radarChart>
%s`, Escape(chart.RadarStyle), series, labels, radarAxesXML(chart)),
	)
}

func radarSeriesXML(chart *ChartSpec) string {
	var b strings.Builder
	b.WriteString(`
<c:ser>
<c:idx val="0"/>
<c:order val="0"/>
<c:tx><c:v>`)
	b.WriteString(Escape(chart.SeriesName))
	b.WriteString(`</c:v></c:tx>
<c:spPr><a:ln><a:solidFill><a:srgbClr val="`)
	b.WriteString(Escape(chart.Color))
	b.WriteString(`"/></a:solidFill></a:ln></c:spPr>
<c:marker><c:symbol val="none"/></c:marker>
<c:cat><c:strLit>
<c:ptCount val="`)
	b.WriteString(strconv.Itoa(len(chart.Categories)))
	b.WriteString(`"/>`)
	for i, category := range chart.Categories {
		b.WriteString(`
<c:pt idx="`)
		b.WriteString(strconv.Itoa(i))
		b.WriteString(`"><c:v>`)
		b.WriteString(Escape(category))
		b.WriteString(`</c:v></c:pt>`)
	}
	b.WriteString(`
</c:strLit></c:cat>
<c:val><c:numLit>
<c:formatCode>General</c:formatCode>
<c:ptCount val="`)
	b.WriteString(strconv.Itoa(len(chart.Values)))
	b.WriteString(`"/>`)
	for i, value := range chart.Values {
		b.WriteString(`
<c:pt idx="`)
		b.WriteString(strconv.Itoa(i))
		b.WriteString(`"><c:v>`)
		b.WriteString(strconv.FormatFloat(value, 'f', 6, 64))
		b.WriteString(`</c:v></c:pt>`)
	}
	b.WriteString(`
</c:numLit></c:val>
</c:ser>`)
	return b.String()
}

func radarAxesXML(chart *ChartSpec) string {
	axes := chartAxesXML(chart)
	axes = strings.ReplaceAll(
		axes,
		`<c:tickLblPos val="nextTo"/>`,
		`<c:majorTickMark val="none"/><c:minorTickMark val="none"/><c:tickLblPos val="nextTo"/>`,
	)
	return strings.Replace(
		axes,
		`<c:crossBetween val="between"/>`,
		`<c:majorTickMark val="none"/><c:minorTickMark val="none"/><c:crossBetween val="between"/>`,
		1,
	)
}
