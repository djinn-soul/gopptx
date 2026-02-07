package pptxxml

import (
	"fmt"
	"strings"
)

// ChartPartXML renders a chart part (`ppt/charts/chartN.xml`).
func ChartPartXML(chart *ChartSpec) string {
	if chart.Kind == ChartKindBar || chart.Kind == ChartKindBarHorizontal ||
		chart.Kind == ChartKindBarStacked || chart.Kind == ChartKindBarStacked100 {
		return barChartPartXML(chart)
	}
	if chart.Kind == ChartKindLine || chart.Kind == ChartKindLineMarkers || chart.Kind == ChartKindLineStacked {
		return lineChartPartXML(chart)
	}
	if chart.Kind == ChartKindBubble {
		return bubbleChartPartXML(chart)
	}
	if chart.Kind == ChartKindScatter {
		return scatterChartPartXML(chart)
	}
	if chart.Kind == ChartKindArea || chart.Kind == ChartKindAreaStacked || chart.Kind == ChartKindAreaStacked100 {
		return areaChartPartXML(chart)
	}
	if chart.Kind == ChartKindPie {
		return pieChartPartXML(chart)
	}
	if chart.Kind == ChartKindDoughnut {
		return doughnutChartPartXML(chart)
	}
	if chart.Kind == ChartKindRadar || chart.Kind == ChartKindRadarFilled {
		return radarChartPartXML(chart)
	}
	if chart.Kind == ChartKindStockHLC || chart.Kind == ChartKindStockOHLC {
		return stockChartPartXML(chart)
	}
	if chart.Kind == ChartKindCombo {
		return comboChartPartXML(chart)
	}
	panic(fmt.Sprintf("unsupported chart kind: %s", chart.Kind))
}

func barChartPartXML(chart *ChartSpec) string {
	series := chartSeriesXML(chart)
	labels := chartDataLabelsXML(chart.ShowDataLabels)
	return chartPartEnvelope(chart.Title, chart.ShowLegend, chart.LegendPosition, fmt.Sprintf(`
<c:barChart>
<c:barDir val="%s"/>
<c:grouping val="%s"/>
<c:varyColors val="0"/>%s
%s
<c:axId val="48650112"/>
<c:axId val="48672768"/>
</c:barChart>
%s`, Escape(chart.BarDir), Escape(chart.Grouping), series, labels, chartAxesXML(chart)))
}

func lineChartPartXML(chart *ChartSpec) string {
	series := chartLineSeriesXML(chart)
	labels := chartDataLabelsXML(chart.ShowDataLabels)
	smooth := "0"
	if chart.Smooth {
		smooth = "1"
	}
	return chartPartEnvelope(chart.Title, chart.ShowLegend, chart.LegendPosition, fmt.Sprintf(`
<c:lineChart>
<c:grouping val="%s"/>
<c:varyColors val="0"/>%s
%s
<c:smooth val="%s"/>
<c:axId val="48650112"/>
<c:axId val="48672768"/>
</c:lineChart>
%s`, Escape(chart.Grouping), series, labels, smooth, chartAxesXML(chart)))
}

func chartPartEnvelope(title string, showLegend bool, legendPosition string, plotXML string) string {
	legend := ""
	if showLegend {
		legendPos := normalizedLegendPosition(legendPosition)
		legend = `
<c:legend>
<c:legendPos val="` + legendPos + `"/>
<c:overlay val="0"/>
</c:legend>`
	}
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<c:chartSpace xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
<c:lang val="en-US"/>
<c:chart>
<c:title>
<c:tx><c:rich><a:bodyPr/><a:lstStyle/><a:p><a:r><a:rPr lang="en-US"/><a:t>%s</a:t></a:r></a:p></c:rich></c:tx>
<c:overlay val="0"/>
</c:title>
<c:autoTitleDeleted val="0"/>
<c:plotArea>
<c:layout/>%s
</c:plotArea>
%s
<c:plotVisOnly val="1"/>
</c:chart>
</c:chartSpace>`, Escape(title), plotXML, legend)
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

func chartLineSeriesXML(chart *ChartSpec) string {
	base := chartSeriesXML(chart)
	if !chart.ShowMarkers {
		return base
	}
	return strings.Replace(
		base,
		"</c:ser>",
		"<c:marker><c:symbol val=\"circle\"/></c:marker></c:ser>",
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
	majorGrid := ""
	if chart.ShowMajorGridlines {
		majorGrid = "<c:majorGridlines/>"
	}

	return fmt.Sprintf(`
<c:catAx>
<c:axId val="48650112"/>
<c:scaling><c:orientation val="minMax"/></c:scaling>
<c:delete val="0"/>
<c:axPos val="b"/>
%s
<c:tickLblPos val="nextTo"/>
<c:crossAx val="48672768"/>
<c:crosses val="autoZero"/>
<c:auto val="1"/>
<c:lblAlgn val="ctr"/>
<c:lblOffset val="100"/>
</c:catAx>
<c:valAx>
<c:axId val="48672768"/>
%s
<c:delete val="0"/>
<c:axPos val="l"/>
%s
%s
%s
<c:tickLblPos val="nextTo"/>
<c:crossAx val="48650112"/>
<c:crosses val="autoZero"/>
<c:crossBetween val="between"/>
</c:valAx>`,
		categoryAxisTitle,
		valueScaling,
		valueAxisTitle,
		valueFormat,
		majorGrid,
	)
}

func chartAxisTitleXML(title string) string {
	trimmed := strings.TrimSpace(title)
	if trimmed == "" {
		return ""
	}
	return fmt.Sprintf(`
<c:title>
<c:tx><c:rich><a:bodyPr/><a:lstStyle/><a:p><a:r><a:rPr lang="en-US"/><a:t>%s</a:t></a:r></a:p></c:rich></c:tx>
<c:overlay val="0"/>
</c:title>`, Escape(trimmed))
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
	return fmt.Sprintf(`<c:numFmt formatCode="%s" sourceLinked="%s"/>`, Escape(trimmed), sourceLinked)
}

func valueAxisScalingXML(minValue *float64, maxValue *float64) string {
	minXML := ""
	maxXML := ""
	if minValue != nil {
		minXML = fmt.Sprintf(`<c:min val="%.6f"/>`, *minValue)
	}
	if maxValue != nil {
		maxXML = fmt.Sprintf(`<c:max val="%.6f"/>`, *maxValue)
	}
	return `<c:scaling><c:orientation val="minMax"/>` + minXML + maxXML + `</c:scaling>`
}
