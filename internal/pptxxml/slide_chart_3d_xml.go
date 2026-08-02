package pptxxml

import (
	"fmt"
	"strconv"
	"strings"
)

// chart3DSeriesAxID is the series (depth) axis a 3-D plot is drawn against.
// CT_Line3DChart requires all three axes; the bar and area 3-D groups accept
// two, so only the line group declares it.
const chart3DSeriesAxID = 48739584

// chartView3DXML is the camera CT_Chart carries for a 3-D plot. Without it
// PowerPoint draws the group flat.
const chartView3DXML = `<c:view3D><c:rotX val="15"/><c:rotY val="20"/>` +
	`<c:depthPercent val="100"/><c:rAngAx val="1"/></c:view3D>`

// withView3D inserts the 3-D camera before the plot area, where CT_Chart
// expects it.
func withView3D(xml string) string {
	return strings.Replace(xml, "<c:plotArea>", chartView3DXML+"<c:plotArea>", 1)
}

// bar3DChartPartXML renders a 3-D column or bar chart. CT_Bar3DChart orders
// barDir, grouping, varyColors, ser, dLbls, gapWidth, gapDepth, shape, axId.
func bar3DChartPartXML(chart *ChartSpec) string {
	series := chartSeriesXML(chart)
	labels := chartDataLabelsXML(chart)
	body := fmt.Sprintf(`
<c:bar3DChart>
<c:barDir val="%s"/>
<c:grouping val="%s"/>
<c:varyColors val="0"/>%s
%s
<c:shape val="box"/>%s
</c:bar3DChart>
%s`, Escape(chart.BarDir), Escape(chart.Grouping), series, labels,
		primaryAxisIDRefsXML(), chartAxesAndDataTableXML(chart))

	return withView3D(chartPartEnvelope(
		chart.Title,
		chart.TitleOverlay,
		chart.ShowLegend,
		chart.LegendPosition,
		chart.LegendOverlay,
		body,
	))
}

// area3DChartPartXML renders a 3-D area chart. CT_Area3DChart orders grouping,
// varyColors, ser, dLbls, dropLines, gapDepth, axId.
func area3DChartPartXML(chart *ChartSpec) string {
	series := chartSeriesXML(chart)
	labels := chartDataLabelsXML(chart)
	body := fmt.Sprintf(`
<c:area3DChart>
<c:grouping val="%s"/>
<c:varyColors val="0"/>%s
%s%s
</c:area3DChart>
%s`, Escape(chart.Grouping), series, labels,
		primaryAxisIDRefsXML(), chartAxesAndDataTableXML(chart))

	return withView3D(chartPartEnvelope(
		chart.Title,
		chart.TitleOverlay,
		chart.ShowLegend,
		chart.LegendPosition,
		chart.LegendOverlay,
		body,
	))
}

// line3DChartPartXML renders a 3-D line chart. CT_Line3DChart has no c:smooth
// child and requires three c:axId references, so it also emits the series axis
// the 2-D line group does not have.
func line3DChartPartXML(chart *ChartSpec) string {
	series := chartLineSeriesXML(chart)
	labels := chartDataLabelsXML(chart)
	axIDs := primaryAxisIDRefsXML() + `
<c:axId val="` + strconv.Itoa(chart3DSeriesAxID) + `"/>`
	body := fmt.Sprintf(`
<c:line3DChart>
<c:grouping val="%s"/>
<c:varyColors val="0"/>%s
%s%s
</c:line3DChart>
%s`, Escape(chart.Grouping), series, labels, axIDs,
		chartAxesXML(chart)+chartSeriesAxisXML()+chartDataTableXML(chart))

	return withView3D(chartPartEnvelope(
		chart.Title,
		chart.TitleOverlay,
		chart.ShowLegend,
		chart.LegendPosition,
		chart.LegendOverlay,
		body,
	))
}

// chartSeriesAxisXML declares the depth axis a 3-D line plot references.
func chartSeriesAxisXML() string {
	return `
<c:serAx>
<c:axId val="` + strconv.Itoa(chart3DSeriesAxID) + `"/>
<c:scaling><c:orientation val="minMax"/></c:scaling>
<c:delete val="0"/>
<c:axPos val="b"/>
<c:tickLblPos val="nextTo"/>
<c:crossAx val="` + strconv.Itoa(primaryValAxID) + `"/>
</c:serAx>`
}
