package pptxxml

import (
	"strconv"
	"strings"
)

type chartDataLabelDefaults struct {
	showLegendKey  bool
	showValue      bool
	showCategory   bool
	showSeriesName bool
	showPercent    bool
	showBubbleSize bool
}

func chartDataLabelsXML(chart *ChartSpec) string {
	return chartDataLabelsWithDefaults(chart, chartDataLabelDefaults{
		showValue: true,
	})
}

func chartDataLabelsWithDefaults(chart *ChartSpec, defaults chartDataLabelDefaults) string {
	if !chart.ShowDataLabels {
		return ""
	}
	position := normalizedDataLabelPosition(chart.DataLabelPosition)
	showLegendKey := resolvedDataLabelBool(chart.DataLabelShowLegendKey, defaults.showLegendKey)
	showValue := resolvedDataLabelBool(chart.DataLabelShowValue, defaults.showValue)
	showCategory := resolvedDataLabelBool(chart.DataLabelShowCategoryName, defaults.showCategory)
	showSeriesName := resolvedDataLabelBool(chart.DataLabelShowSeriesName, defaults.showSeriesName)
	showPercent := resolvedDataLabelBool(chart.DataLabelShowPercent, defaults.showPercent)
	showBubbleSize := resolvedDataLabelBool(chart.DataLabelShowBubbleSize, defaults.showBubbleSize)

	var b strings.Builder
	b.WriteString(`
<c:dLbls>
`)
	if position != "" {
		b.WriteString(`<c:dLblPos val="`)
		b.WriteString(position)
		b.WriteString(`"/>
`)
	}
	if showLegendKey {
		b.WriteString(`<c:showLegendKey val="1"/>
`)
	}
	if showValue {
		b.WriteString(`<c:showVal val="1"/>
`)
	}
	if showCategory {
		b.WriteString(`<c:showCatName val="1"/>
`)
	}
	if showSeriesName {
		b.WriteString(`<c:showSerName val="1"/>
`)
	}
	if showPercent {
		b.WriteString(`<c:showPercent val="1"/>
`)
	}
	if showBubbleSize {
		b.WriteString(`<c:showBubbleSize val="1"/>
`)
	}
	b.WriteString(`</c:dLbls>`)
	return b.String()
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
	return buildChartAxesXML(chart, "catAx", false, true)
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

func normalizedAxisTickLabelPosition(pos string) string {
	switch strings.ToLower(strings.TrimSpace(pos)) {
	case "low", "high", "none":
		return strings.ToLower(strings.TrimSpace(pos))
	default:
		return "nextTo"
	}
}

func normalizedAxisCrosses(mode string) string {
	switch strings.TrimSpace(mode) {
	case "min", "max":
		return strings.TrimSpace(mode)
	default:
		return "autoZero"
	}
}

func normalizedDataLabelPosition(position string) string {
	switch strings.TrimSpace(position) {
	case "ctr", "inEnd", "inBase", "outEnd", "bestFit", "l", "r", "t", "b":
		return strings.TrimSpace(position)
	default:
		return ""
	}
}

func resolvedDataLabelBool(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}
	return *value
}
