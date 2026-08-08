package pptxxml

import (
	"fmt"
	"strings"
)

const chartMajorGridlinesXML = "<c:majorGridlines/>"

// chartAxesAndDataTableXML renders the axes followed by the optional data
// table. CT_PlotArea puts <c:dTable> after the axes, so the two are emitted
// together to keep that order in every chart kind that supports a data table.
func chartAxesAndDataTableXML(chart *ChartSpec) string {
	return chartAxesXML(chart) + chartDataTableXML(chart)
}

// chartDataTableXML renders <c:dTable>, the grid of plotted values PowerPoint
// draws beneath the plot area.
func chartDataTableXML(chart *ChartSpec) string {
	if !chart.ShowDataTable {
		return ""
	}
	flag := func(v *bool) string {
		if v == nil {
			return "1"
		}
		return boolToOneZero(*v)
	}
	return "\n<c:dTable>" +
		`<c:showHorzBorder val="` + flag(chart.DataTableShowHorzBorder) + `"/>` +
		`<c:showVertBorder val="` + flag(chart.DataTableShowVertBorder) + `"/>` +
		`<c:showOutline val="` + flag(chart.DataTableShowOutline) + `"/>` +
		`<c:showKeys val="` + flag(chart.DataTableShowLegendKeys) + `"/>` +
		"</c:dTable>"
}

// ChartPartXML renders a chart part (`ppt/charts/chartN.xml`).
func ChartPartXML(chart *ChartSpec) string {
	return string(RenderChart(chart))
}

// RenderChart renders a chart part to bytes.
func RenderChart(chart *ChartSpec) []byte {
	return withExternalData(
		withDisplayBlanksAs(renderChartBody(chart), chart.DisplayBlanksAs),
		chart.ExternalDataID,
	)
}

// withExternalData names the embedded workbook the chart was built from.
// Without it PowerPoint's "Edit Data" has nothing to open, even when the
// package carries the .xlsx part. CT_ChartSpace puts <c:externalData> after
// <c:chart>, so it goes just before the closing tag.
func withExternalData(chartXML []byte, relID string) []byte {
	if relID == "" {
		return chartXML
	}
	const closeTag = "</c:chartSpace>"
	node := `<c:externalData r:id="` + Escape(relID) + `"><c:autoUpdate val="0"/></c:externalData>` + "\n"
	return []byte(strings.Replace(string(chartXML), closeTag, node+closeTag, 1))
}

// withDisplayBlanksAs inserts <c:dispBlanksAs> after <c:plotVisOnly>, the only
// position the schema allows it in. Every chart kind goes through here, so a
// blank in any of them is drawn the same way.
func withDisplayBlanksAs(chartXML []byte, displayBlanksAs string) []byte {
	node := displayBlanksAsXML(displayBlanksAs)
	if node == "" {
		return chartXML
	}
	return []byte(strings.Replace(
		string(chartXML), plotVisOnlyElement, plotVisOnlyElement+node, 1,
	))
}

func renderChartBody(chart *ChartSpec) []byte {
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
	if chart.Kind == ChartKindThreeDPie {
		return []byte(pie3DChartPartXML(chart))
	}
	if chart.Kind == ChartKindThreeDColumn || chart.Kind == ChartKindThreeDBar {
		return []byte(bar3DChartPartXML(chart))
	}
	if chart.Kind == ChartKindThreeDLine {
		return []byte(line3DChartPartXML(chart))
	}
	if chart.Kind == ChartKindThreeDArea {
		return []byte(area3DChartPartXML(chart))
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
	panic("unsupported chart kind: " + chart.Kind)
}

func barChartPartXML(chart *ChartSpec) string {
	series := chartSeriesXML(chart)
	labels := chartDataLabelsXML(chart)
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
%s%s
</c:barChart>
%s`, Escape(chart.BarDir), Escape(chart.Grouping), series, labels, primaryAxisIDRefsXML(), chartAxesAndDataTableXML(chart)),
	)
}

func lineChartPartXML(chart *ChartSpec) string {
	series := chartLineSeriesXML(chart)
	labels := chartDataLabelsXML(chart)
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
<c:smooth val="%s"/>%s
</c:lineChart>
%s`, Escape(chart.Grouping), series, labels, smooth, primaryAxisIDRefsXML(), chartAxesAndDataTableXML(chart)),
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
	// date1904 and roundedCorners are the two CT_ChartSpace defaults PowerPoint
	// always writes. roundedCorners is the visible one: with the element absent
	// PowerPoint rounds the chart frame.
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<c:chartSpace xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" `+
		`xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" `+
		`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
<c:date1904 val="0"/>
<c:lang val="en-US"/>
<c:roundedCorners val="0"/>
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
`+plotVisOnlyElement+`
<c:showDLblsOverMax val="0"/>
</c:chart>
</c:chartSpace>`, Escape(title), boolToOneZero(titleOverlay), plotXML, legend)
}

func chartSeriesXML(chart *ChartSpec) string {
	bound := chartIsWorkbookBound(chart)
	var b strings.Builder
	b.WriteString(`
<c:ser>
<c:idx val="0"/>
<c:order val="0"/>`)
	writeChartSeriesName(&b, chart.SeriesName, bound)
	b.WriteString(`
<c:spPr><a:solidFill><a:srgbClr val="` + Escape(chart.Color) + `"/></a:solidFill></c:spPr>`)
	writeChartCategories(&b, chart.Categories, bound)
	writeChartValues(&b, chart.Values, bound)
	b.WriteString(`
</c:ser>`)
	return b.String()
}

func chartLineSeriesXML(chart *ChartSpec) string {
	base := chartSeriesXML(chart)
	if chart.ShowMarkers {
		return strings.Replace(
			base,
			"</c:spPr>",
			"</c:spPr><c:marker><c:symbol val=\"circle\"/></c:marker>",
			1,
		)
	}
	return strings.Replace(
		base,
		"</c:spPr>",
		"</c:spPr><c:marker><c:symbol val=\"none\"/></c:marker>",
		1,
	)
}
