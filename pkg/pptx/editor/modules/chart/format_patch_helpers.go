package chart

import "strings"

func firstChartBlockBounds(xml string) (int, int) {
	chartTags := []string{
		"barChart",
		"lineChart",
		"areaChart",
		"pieChart",
		"doughnutChart",
		"scatterChart",
		"bubbleChart",
		"radarChart",
		"stockChart",
	}
	for _, tag := range chartTags {
		startTag := "<c:" + tag + ">"
		start := strings.Index(xml, startTag)
		if start < 0 {
			continue
		}
		endTag := chartElementClosePrefix + tag + ">"
		relEnd := strings.Index(xml[start:], endTag)
		if relEnd < 0 {
			continue
		}
		return start, start + relEnd + len(endTag)
	}
	return -1, -1
}

func isLegendPosition(position string) bool {
	switch strings.ToLower(strings.TrimSpace(position)) {
	case "r", "l", "t", "b", "tr", "right", "left", "top", "bottom", "top_right", "topright":
		return true
	default:
		return false
	}
}

func isDataLabelPosition(position string) bool {
	norm := normalizeDataLabelPosition(position, false)
	switch norm {
	case dataLabelPositionCenter, dataLabelPositionInsideEnd, dataLabelPositionInsideBase,
		dataLabelPositionOutsideEnd, dataLabelPositionBestFit, "l", "r", "t", "b":
		return true
	default:
		return false
	}
}
