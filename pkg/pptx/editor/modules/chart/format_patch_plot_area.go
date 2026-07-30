package chart

import (
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func patchPlotAreaLine(xml string, line *common.ChartLineFormat) string {
	if line == nil {
		return xml
	}
	rendered := buildChartLine(line)
	if rendered == "" {
		return xml
	}
	start, end := plotAreaBounds(xml)
	if start < 0 || end <= start {
		return xml
	}
	block := xml[start:end]
	spIndex := plotAreaShapePropertiesIndex(block)
	if spIndex < 0 {
		patched := block + "<c:spPr>" + rendered + "</c:spPr>"
		return xml[:start] + patched + xml[end:]
	}
	current := reChartLinesShapeProps.FindString(block[spIndex:])
	if current == "" {
		return xml
	}
	var replacement string
	if strings.HasSuffix(current, "/>") {
		replacement = "<c:spPr>" + rendered + "</c:spPr>"
	} else {
		replacement = setShapePropertiesLine(current, rendered)
	}
	patched := block[:spIndex] +
		strings.Replace(block[spIndex:], current, replacement, 1)
	return xml[:start] + patched + xml[end:]
}

func parsePlotAreaLine(xml string) *common.ChartLineFormat {
	start, end := plotAreaBounds(xml)
	if start < 0 || end <= start {
		return nil
	}
	block := xml[start:end]
	spIndex := plotAreaShapePropertiesIndex(block)
	if spIndex < 0 {
		return nil
	}
	spPr := reChartLinesShapeProps.FindString(block[spIndex:])
	if spPr == "" {
		return nil
	}
	return parseChartLineFormat(spPr)
}
