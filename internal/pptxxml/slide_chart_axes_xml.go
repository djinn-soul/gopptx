package pptxxml

import (
	"strconv"
	"strings"
)

func buildChartAxesXML(
	chart *ChartSpec,
	primaryAxisTag string,
	includePrimaryNumFmt bool,
	includePrimaryCategoryDefaults bool,
) string {
	var b strings.Builder
	b.WriteString(primaryCategoryAxisXML(
		chart, primaryAxisTag, includePrimaryNumFmt, includePrimaryCategoryDefaults,
	))
	b.WriteString(primaryValueAxisXML(chart))
	if chart.SecondaryAxis {
		b.WriteString(secondaryChartAxesXML(chart, primaryAxisTag))
	}
	return b.String()
}

// primaryCategoryAxisXML renders the category (or date) axis every chart has.
func primaryCategoryAxisXML(
	chart *ChartSpec,
	primaryAxisTag string,
	includeNumFmt bool,
	includeCategoryDefaults bool,
) string {
	majorGrid := ""
	if chart.ShowCategoryMajorGridlines {
		majorGrid = chartMajorGridlinesXML
	}

	var b strings.Builder
	b.WriteString(`
<c:`)
	b.WriteString(primaryAxisTag)
	b.WriteString(`>
<c:axId val="`)
	b.WriteString(strconv.Itoa(primaryCatAxID))
	b.WriteString(`"/>
<c:scaling><c:orientation val="minMax"/></c:scaling>
<c:delete val="0"/>
<c:axPos val="b"/>`)
	b.WriteString(majorGrid)
	b.WriteString(chartAxisTitleXML(chart.CategoryAxisTitle))
	if includeNumFmt {
		b.WriteString(`
<c:numFmt formatCode="General" sourceLinked="1"/>`)
	}
	b.WriteString(`
<c:tickLblPos val="`)
	b.WriteString(normalizedAxisTickLabelPosition(chart.CategoryTickLabelPosition))
	b.WriteString(`"/>
<c:crossAx val="`)
	b.WriteString(strconv.Itoa(primaryValAxID))
	b.WriteString(`"/>
<c:crosses val="`)
	b.WriteString(normalizedAxisCrosses(chart.CategoryAxisCrosses))
	b.WriteString(`"/>`)
	if includeCategoryDefaults {
		b.WriteString(`
<c:auto val="1"/>
<c:lblAlgn val="ctr"/>
<c:lblOffset val="100"/>`)
	}
	b.WriteString(`
</c:`)
	b.WriteString(primaryAxisTag)
	b.WriteString(`>`)
	return b.String()
}

// primaryValueAxisXML renders the value axis on the left.
func primaryValueAxisXML(chart *ChartSpec) string {
	majorGrid := ""
	if chart.ShowMajorGridlines {
		majorGrid = chartMajorGridlinesXML
	}

	var b strings.Builder
	b.WriteString(`
<c:valAx>
<c:axId val="`)
	b.WriteString(strconv.Itoa(primaryValAxID))
	b.WriteString(`"/>`)
	b.WriteString(valueAxisScalingXML(chart.MinValue, chart.MaxValue))
	b.WriteString(`
<c:delete val="0"/>
<c:axPos val="l"/>`)
	b.WriteString(majorGrid)
	b.WriteString(chartAxisTitleXML(chart.ValueAxisTitle))
	b.WriteString(chartValueFormatXML(chart.ValueFormat))
	b.WriteString(`
<c:tickLblPos val="`)
	b.WriteString(normalizedAxisTickLabelPosition(chart.ValueTickLabelPosition))
	b.WriteString(`"/>
<c:crossAx val="`)
	b.WriteString(strconv.Itoa(primaryCatAxID))
	b.WriteString(`"/>
<c:crosses val="`)
	b.WriteString(normalizedAxisCrosses(chart.ValueAxisCrosses))
	b.WriteString(`"/>
<c:crossBetween val="`)
	b.WriteString(normalizedValueAxisCrossBetween(chart.ValueAxisCrossBetween))
	b.WriteString(`"/>
</c:valAx>`)
	return b.String()
}

// secondaryChartAxesXML renders the second value axis and the hidden category
// axis it crosses.
//
// A plot bound to the secondary pair needs both axes declared: the value axis
// is drawn on the right at the far end of its category axis, and that category
// axis is deleted so the chart does not show two identical category scales.
func secondaryChartAxesXML(chart *ChartSpec, primaryAxisTag string) string {
	title := chartAxisTitleXML(chart.SecondaryValueAxisTitle)
	scaling := valueAxisScalingXML(chart.SecondaryMinValue, chart.SecondaryMaxValue)
	format := chartValueFormatXML(chart.SecondaryValueFormat)
	tickLblPos := normalizedAxisTickLabelPosition(chart.SecondaryValueTickLabelPosition)
	crossBetween := normalizedValueAxisCrossBetween(chart.ValueAxisCrossBetween)
	majorGrid := ""
	if chart.ShowSecondaryMajorGridlines {
		majorGrid = chartMajorGridlinesXML
	}

	var b strings.Builder
	b.WriteString(`
<c:valAx>
<c:axId val="`)
	b.WriteString(strconv.Itoa(secondaryValAxID))
	b.WriteString(`"/>`)
	b.WriteString(scaling)
	b.WriteString(`
<c:delete val="0"/>
<c:axPos val="r"/>`)
	b.WriteString(majorGrid)
	b.WriteString(title)
	b.WriteString(format)
	b.WriteString(`
<c:tickLblPos val="`)
	b.WriteString(tickLblPos)
	b.WriteString(`"/>
<c:crossAx val="`)
	b.WriteString(strconv.Itoa(secondaryCatAxID))
	b.WriteString(`"/>
<c:crosses val="max"/>
<c:crossBetween val="`)
	b.WriteString(crossBetween)
	b.WriteString(`"/>
</c:valAx>
<c:`)
	b.WriteString(primaryAxisTag)
	b.WriteString(`>
<c:axId val="`)
	b.WriteString(strconv.Itoa(secondaryCatAxID))
	b.WriteString(`"/>
<c:scaling><c:orientation val="minMax"/></c:scaling>
<c:delete val="1"/>
<c:axPos val="b"/>
<c:crossAx val="`)
	b.WriteString(strconv.Itoa(secondaryValAxID))
	b.WriteString(`"/>
</c:`)
	b.WriteString(primaryAxisTag)
	b.WriteString(`>`)
	return b.String()
}
