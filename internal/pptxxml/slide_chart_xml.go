package pptxxml

import (
	"fmt"
	"strings"
)

// ChartSpec describes one chart written as a chart part.
type ChartSpec struct {
	Kind       string
	Title      string
	Categories []string
	Values     []float64
	X          int64
	Y          int64
	CX         int64
	CY         int64
	Color      string
}

const (
	ChartKindBar  = "bar"
	ChartKindLine = "line"
)

// ChartPartXML renders a chart part (`ppt/charts/chartN.xml`).
func ChartPartXML(chart *ChartSpec) string {
	if chart.Kind == ChartKindLine {
		return lineChartPartXML(chart)
	}
	return barChartPartXML(chart)
}

func barChartPartXML(chart *ChartSpec) string {
	series := chartSeriesXML(chart)
	return chartPartEnvelope(chart.Title, fmt.Sprintf(`
<c:barChart>
<c:barDir val="col"/>
<c:grouping val="clustered"/>
<c:varyColors val="0"/>%s
<c:axId val="48650112"/>
<c:axId val="48672768"/>
</c:barChart>
%s`, series, chartAxesXML()))
}

func lineChartPartXML(chart *ChartSpec) string {
	series := chartSeriesXML(chart)
	return chartPartEnvelope(chart.Title, fmt.Sprintf(`
<c:lineChart>
<c:grouping val="standard"/>
<c:varyColors val="0"/>%s
<c:smooth val="0"/>
<c:axId val="48650112"/>
<c:axId val="48672768"/>
</c:lineChart>
%s`, series, chartAxesXML()))
}

func chartPartEnvelope(title string, plotXML string) string {
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
<c:plotVisOnly val="1"/>
</c:chart>
</c:chartSpace>`, Escape(title), plotXML)
}

func chartSeriesXML(chart *ChartSpec) string {
	var b strings.Builder
	b.WriteString(`
<c:ser>
<c:idx val="0"/>
<c:order val="0"/>
<c:tx><c:v>Series 1</c:v></c:tx>
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

func chartAxesXML() string {
	return `
<c:catAx>
<c:axId val="48650112"/>
<c:scaling><c:orientation val="minMax"/></c:scaling>
<c:delete val="0"/>
<c:axPos val="b"/>
<c:tickLblPos val="nextTo"/>
<c:crossAx val="48672768"/>
<c:crosses val="autoZero"/>
<c:auto val="1"/>
<c:lblAlgn val="ctr"/>
<c:lblOffset val="100"/>
</c:catAx>
<c:valAx>
<c:axId val="48672768"/>
<c:scaling><c:orientation val="minMax"/></c:scaling>
<c:delete val="0"/>
<c:axPos val="l"/>
<c:numFmt formatCode="General" sourceLinked="1"/>
<c:majorGridlines/>
<c:tickLblPos val="nextTo"/>
<c:crossAx val="48650112"/>
<c:crosses val="autoZero"/>
<c:crossBetween val="between"/>
</c:valAx>`
}
