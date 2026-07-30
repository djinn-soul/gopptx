package chart

import (
	"errors"
	"regexp"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// Series lines, upstream #846: the c:serLines connectors PowerPoint draws
// between the segments of a stacked bar chart, and between the pie and the bar
// of a bar-of-pie chart. CT_BarChart and CT_OfPieChart are the only plots that
// accept the element; writing it anywhere else makes PowerPoint repair the file.

var reSeriesLinesBlock = regexp.MustCompile(`(?s)<c:serLines(?:\s*/>|>.*?</c:serLines>)`)

// seriesLinePlots are the plot elements whose schema has c:serLines.
//
//nolint:gochecknoglobals // Fixed element set, as elsewhere in this package.
var seriesLinePlots = []string{barChartPlotTag, "bar3DChart", "ofPieChart"}

func validateChartSeriesLines(lines *common.ChartSeriesLines) error {
	if lines == nil {
		return nil
	}
	if lines.Show == nil && lines.Line == nil {
		return errors.New("series_lines needs show or line")
	}
	return validateChartLineFormat("series_lines line", lines.Line)
}

// patchChartSeriesLines writes c:serLines into every plot that accepts it.
// A line style with no Show implies drawing the connectors, since there is
// nothing else for the style to apply to.
func patchChartSeriesLines(xml string, lines *common.ChartSeriesLines) string {
	if lines == nil {
		return xml
	}
	remove := lines.Show != nil && !*lines.Show
	rendered := buildChartLine(lines.Line)
	for _, plot := range seriesLinePlots {
		xml = patchPlotSeriesLines(xml, plot, remove, rendered)
	}
	return xml
}

func patchPlotSeriesLines(xml string, plot string, remove bool, rendered string) string {
	startTag, endTag := "<c:"+plot+">", chartElementClosePrefix+plot+">"
	start := strings.Index(xml, startTag)
	for start >= 0 {
		endRel := strings.Index(xml[start:], endTag)
		if endRel < 0 {
			return xml
		}
		end := start + endRel + len(endTag)
		block := applyPlotSeriesLines(xml[start:end], remove, rendered)
		xml = xml[:start] + block + xml[end:]
		next := strings.Index(xml[start+len(block):], startTag)
		if next < 0 {
			return xml
		}
		start += len(block) + next
	}
	return xml
}

func applyPlotSeriesLines(block string, remove bool, rendered string) string {
	if remove {
		return reSeriesLinesBlock.ReplaceAllLiteralString(block, "")
	}
	node := "<c:serLines/>"
	if rendered != "" {
		node = "<c:serLines><c:spPr>" + rendered + "</c:spPr></c:serLines>"
	}
	if current := reSeriesLinesBlock.FindString(block); current != "" {
		if rendered == "" {
			return block
		}
		return strings.Replace(block, current, styleChartLinesElement(current, "serLines", rendered), 1)
	}
	insertAt := seriesLinesInsertIndex(block)
	if insertAt < 0 {
		return block
	}
	return block[:insertAt] + node + block[insertAt:]
}

// seriesLinesInsertIndex finds the schema slot for c:serLines: CT_BarChart puts
// it after c:gapWidth and c:overlap and before the c:axId references, and
// CT_OfPieChart after c:secondPieSize.
func seriesLinesInsertIndex(block string) int {
	if index := strings.Index(block, "<c:axId"); index >= 0 {
		return index
	}
	if index := strings.Index(block, chartExtensionListTagPrefix); index >= 0 {
		return index
	}
	return strings.LastIndex(block, chartElementClosePrefix)
}

// parseChartSeriesLines reads back the connectors of the first plot that has
// them, so a caller can tell a styled set from a default one.
func parseChartSeriesLines(xml string) *common.ChartSeriesLines {
	block := reSeriesLinesBlock.FindString(xml)
	if block == "" {
		return nil
	}
	show := true
	return &common.ChartSeriesLines{Show: &show, Line: parseChartLineFormat(block)}
}
