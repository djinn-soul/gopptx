package pptxxml

import "strings"

func buildChartAxesXML(
	chart *ChartSpec,
	primaryAxisTag string,
	includePrimaryNumFmt bool,
	includePrimaryCategoryDefaults bool,
) string {
	primaryAxisTitle := chartAxisTitleXML(chart.CategoryAxisTitle)
	valueAxisTitle := chartAxisTitleXML(chart.ValueAxisTitle)
	valueScaling := valueAxisScalingXML(chart.MinValue, chart.MaxValue)
	valueFormat := chartValueFormatXML(chart.ValueFormat)
	crossBetween := normalizedValueAxisCrossBetween(chart.ValueAxisCrossBetween)
	valueMajorGrid := ""
	if chart.ShowMajorGridlines {
		valueMajorGrid = chartMajorGridlinesXML
	}
	primaryMajorGrid := ""
	if chart.ShowCategoryMajorGridlines {
		primaryMajorGrid = chartMajorGridlinesXML
	}
	primaryTickLblPos := normalizedAxisTickLabelPosition(chart.CategoryTickLabelPosition)
	valueTickLblPos := normalizedAxisTickLabelPosition(chart.ValueTickLabelPosition)
	primaryCrosses := normalizedAxisCrosses(chart.CategoryAxisCrosses)
	valueCrosses := normalizedAxisCrosses(chart.ValueAxisCrosses)

	var b strings.Builder
	b.WriteString(`
<c:`)
	b.WriteString(primaryAxisTag)
	b.WriteString(`>
<c:axId val="48650112"/>
<c:scaling><c:orientation val="minMax"/></c:scaling>
<c:delete val="0"/>
<c:axPos val="b"/>`)
	b.WriteString(primaryMajorGrid)
	b.WriteString(primaryAxisTitle)
	if includePrimaryNumFmt {
		b.WriteString(`
<c:numFmt formatCode="General" sourceLinked="1"/>`)
	}
	b.WriteString(`
<c:tickLblPos val="`)
	b.WriteString(primaryTickLblPos)
	b.WriteString(`"/>
<c:crossAx val="48672768"/>
<c:crosses val="`)
	b.WriteString(primaryCrosses)
	b.WriteString(`"/>`)
	if includePrimaryCategoryDefaults {
		b.WriteString(`
<c:auto val="1"/>
<c:lblAlgn val="ctr"/>
<c:lblOffset val="100"/>`)
	}
	b.WriteString(`
</c:`)
	b.WriteString(primaryAxisTag)
	b.WriteString(`>
<c:valAx>
<c:axId val="48672768"/>`)
	b.WriteString(valueScaling)
	b.WriteString(`
<c:delete val="0"/>
<c:axPos val="l"/>`)
	b.WriteString(valueMajorGrid)
	b.WriteString(valueAxisTitle)
	b.WriteString(valueFormat)
	b.WriteString(`
<c:tickLblPos val="`)
	b.WriteString(valueTickLblPos)
	b.WriteString(`"/>
<c:crossAx val="48650112"/>
<c:crosses val="`)
	b.WriteString(valueCrosses)
	b.WriteString(`"/>
<c:crossBetween val="`)
	b.WriteString(crossBetween)
	b.WriteString(`"/>
</c:valAx>`)
	return b.String()
}
