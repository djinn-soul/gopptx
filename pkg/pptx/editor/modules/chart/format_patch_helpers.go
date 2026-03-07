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
		endTag := "</c:" + tag + ">"
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
	case "r", "l", "t", "b":
		return true
	default:
		return false
	}
}

func isDataLabelPosition(position string) bool {
	switch strings.TrimSpace(position) {
	case "ctr", "inEnd", "inBase", "outEnd", "bestFit", "l", "r", "t", "b":
		return true
	default:
		return false
	}
}
