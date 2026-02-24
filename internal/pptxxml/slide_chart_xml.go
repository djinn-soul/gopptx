package pptxxml

import (
	"fmt"
	"strconv"
	"strings"
)

// ChartPartXML renders a chart part (`ppt/charts/chartN.xml`).
func ChartPartXML(chart *ChartSpec) string {
	return string(RenderChart(chart))
}

// RenderChart renders a chart part to bytes.
func RenderChart(chart *ChartSpec) []byte {
	if chart.Kind == ChartKindBar || chart.Kind == ChartKindBarHorizontal ||
		chart.Kind == ChartKindBarStacked || chart.Kind == ChartKindBarStacked100 {
		return []byte(barChartPartXML(chart))
	}
	if chart.Kind == ChartKindLine || chart.Kind == ChartKindLineMarkers || chart.Kind == ChartKindLineStacked {
		return []byte(lineChartPartXML(chart))
	}
	if chart.Kind == ChartKindBubble {
		return []byte(bubbleChartPartXML(chart))
	}
	if chart.Kind == ChartKindScatter {
		return []byte(scatterChartPartXML(chart))
	}
	if chart.Kind == ChartKindArea || chart.Kind == ChartKindAreaStacked || chart.Kind == ChartKindAreaStacked100 {
		return []byte(areaChartPartXML(chart))
	}
	if chart.Kind == ChartKindPie {
		return []byte(pieChartPartXML(chart))
	}
	if chart.Kind == ChartKindDoughnut {
		return []byte(doughnutChartPartXML(chart))
	}
	if chart.Kind == ChartKindRadar || chart.Kind == ChartKindRadarFilled {
		return []byte(radarChartPartXML(chart))
	}
	if chart.Kind == ChartKindStockHLC || chart.Kind == ChartKindStockOHLC {
		return []byte(stockChartPartXML(chart))
	}
	if chart.Kind == ChartKindCombo {
		return []byte(comboChartPartXML(chart))
	}
	panic(fmt.Sprintf("unsupported chart kind: %s", chart.Kind))
}

func barChartPartXML(chart *ChartSpec) string {
	series := chartSeriesXML(chart)
	labels := chartDataLabelsXML(chart.ShowDataLabels)
	return chartPartEnvelope(
		chart.Title,
		chart.TitleOverlay,
		chart.ShowLegend,
		chart.LegendPosition,
		chart.LegendOverlay,
		fmt.Sprintf(`
<c:barChart>
<c:barDir val="%s"/>
<c:grouping val="%s"/>
<c:varyColors val="0"/>%s
%s
<c:axId val="48650112"/>
<c:axId val="48672768"/>
</c:barChart>
%s`, Escape(chart.BarDir), Escape(chart.Grouping), series, labels, chartAxesXML(chart)),
	)
}

func lineChartPartXML(chart *ChartSpec) string {
	series := chartLineSeriesXML(chart)
	labels := chartDataLabelsXML(chart.ShowDataLabels)
	smooth := "0"
	if chart.Smooth {
		smooth = "1"
	}
	return chartPartEnvelope(
		chart.Title,
		chart.TitleOverlay,
		chart.ShowLegend,
		chart.LegendPosition,
		chart.LegendOverlay,
		fmt.Sprintf(`
<c:lineChart>
<c:grouping val="%s"/>
<c:varyColors val="0"/>%s
%s
<c:smooth val="%s"/>
<c:axId val="48650112"/>
<c:axId val="48672768"/>
</c:lineChart>
%s`, Escape(chart.Grouping), series, labels, smooth, chartAxesXML(chart)),
	)
}

func chartPartEnvelope(
	title string,
	titleOverlay bool,
	showLegend bool,
	legendPosition string,
	legendOverlay bool,
	plotXML string,
) string {
	legend := ""
	if showLegend {
		legendPos := normalizedLegendPosition(legendPosition)
		legend = `
<c:legend>
<c:legendPos val="` + legendPos + `"/>
<c:overlay val="` + boolToOneZero(legendOverlay) + `"/>
</c:legend>`
	}
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<c:chartSpace xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" `+
		`xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" `+
		`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
<c:lang val="en-US"/>
<c:chart>
<c:title>
<c:tx><c:rich><a:bodyPr/><a:lstStyle/><a:p><a:r><a:rPr lang="en-US"/><a:t>%s</a:t></a:r></a:p></c:rich></c:tx>
<c:overlay val="%s"/>
</c:title>
<c:autoTitleDeleted val="0"/>
<c:plotArea>
<c:layout/>%s
</c:plotArea>
%s
<c:plotVisOnly val="1"/>
</c:chart>
</c:chartSpace>`, Escape(title), boolToOneZero(titleOverlay), plotXML, legend)
}

func chartSeriesXML(chart *ChartSpec) string {
	seriesName := chart.SeriesName
	var b strings.Builder
	b.WriteString(`
<c:ser>
<c:idx val="0"/>
<c:order val="0"/>
<c:tx><c:v>` + Escape(seriesName) + `</c:v></c:tx>
<c:spPr><a:solidFill><a:srgbClr val="` + Escape(chart.Color) + `"/></a:solidFill></c:spPr>
<c:cat><c:strLit>`)

	b.WriteString(`
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
<c:val><c:numLit>`)

	b.WriteString(`
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

func chartLineSeriesXML(chart *ChartSpec) string {
	base := chartSeriesXML(chart)
	if !chart.ShowMarkers {
		return base
	}
	return strings.Replace(
		base,
		"</c:spPr>",
		"</c:spPr><c:marker><c:symbol val=\"circle\"/></c:marker>",
		1,
	)
}

func chartDataLabelsXML(show bool) string {
	if !show {
		return ""
	}
	return `
<c:dLbls>
<c:showVal val="1"/>
</c:dLbls>`
}

func normalizedLegendPosition(pos string) string {
	switch strings.ToLower(strings.TrimSpace(pos)) {
	case "l", "t", "b":
		return strings.ToLower(strings.TrimSpace(pos))
	default:
		return "r"
	}
}

func chartAxesXML(chart *ChartSpec) string {
	categoryAxisTitle := chartAxisTitleXML(chart.CategoryAxisTitle)
	valueAxisTitle := chartAxisTitleXML(chart.ValueAxisTitle)
	valueScaling := valueAxisScalingXML(chart.MinValue, chart.MaxValue)
	valueFormat := chartValueFormatXML(chart.ValueFormat)
	crossBetween := normalizedValueAxisCrossBetween(chart.ValueAxisCrossBetween)
	majorGrid := ""
	if chart.ShowMajorGridlines {
		majorGrid = "<c:majorGridlines/>"
	}

	var b strings.Builder
	b.WriteString(`
<c:catAx>
<c:axId val="48650112"/>
<c:scaling><c:orientation val="minMax"/></c:scaling>
<c:delete val="0"/>
<c:axPos val="b"/>`)
	b.WriteString(categoryAxisTitle)
	b.WriteString(`
<c:tickLblPos val="nextTo"/>
<c:crossAx val="48672768"/>
<c:crosses val="autoZero"/>
<c:auto val="1"/>
<c:lblAlgn val="ctr"/>
<c:lblOffset val="100"/>
</c:catAx>
<c:valAx>
<c:axId val="48672768"/>`)
	b.WriteString(valueScaling)
	b.WriteString(`
<c:delete val="0"/>
<c:axPos val="l"/>`)
	b.WriteString(majorGrid)
	b.WriteString(valueAxisTitle)
	b.WriteString(valueFormat)
	b.WriteString(`
<c:tickLblPos val="nextTo"/>
<c:crossAx val="48650112"/>
<c:crosses val="autoZero"/>
<c:crossBetween val="`)
	b.WriteString(crossBetween)
	b.WriteString(`"/>
</c:valAx>`)
	return b.String()
}

func chartAxisTitleXML(title string) string {
	trimmed := strings.TrimSpace(title)
	if trimmed == "" {
		return ""
	}
	var b strings.Builder
	b.WriteString(`
<c:title>
<c:tx><c:rich><a:bodyPr/><a:lstStyle/><a:p><a:r><a:rPr lang="en-US"/><a:t>`)
	b.WriteString(Escape(trimmed))
	b.WriteString(`</a:t></a:r></a:p></c:rich></c:tx>
<c:overlay val="0"/>
</c:title>`)
	return b.String()
}

func chartValueFormatXML(format string) string {
	trimmed := strings.TrimSpace(format)
	if trimmed == "" {
		trimmed = "General"
	}

	sourceLinked := "0"
	if strings.EqualFold(trimmed, "General") {
		sourceLinked = "1"
	}
	var b strings.Builder
	b.WriteString(`<c:numFmt formatCode="`)
	b.WriteString(Escape(trimmed))
	b.WriteString(`" sourceLinked="`)
	b.WriteString(sourceLinked)
	b.WriteString(`"/>`)
	return b.String()
}

func valueAxisScalingXML(minValue *float64, maxValue *float64) string {
	var b strings.Builder
	b.WriteString(`<c:scaling><c:orientation val="minMax"/>`)
	if minValue != nil {
		b.WriteString(`<c:min val="`)
		b.WriteString(strconv.FormatFloat(*minValue, 'f', 6, 64))
		b.WriteString(`"/>`)
	}
	if maxValue != nil {
		b.WriteString(`<c:max val="`)
		b.WriteString(strconv.FormatFloat(*maxValue, 'f', 6, 64))
		b.WriteString(`"/>`)
	}
	b.WriteString(`</c:scaling>`)
	return b.String()
}

func normalizedValueAxisCrossBetween(mode string) string {
	switch strings.TrimSpace(mode) {
	case "midCat":
		return "midCat"
	default:
		return "between"
	}
}

func boolToOneZero(value bool) string {
	if value {
		return "1"
	}
	return "0"
}
