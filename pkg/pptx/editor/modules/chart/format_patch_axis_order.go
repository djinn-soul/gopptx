package chart

import "strings"

func axisScalingInsertIndex(block string) int {
	if axIDEnd := strings.Index(block, "/>"); axIDEnd >= 0 {
		return axIDEnd + len("/>")
	}
	return -1
}

func insertAxisTitle(block string, node string) string {
	return insertBeforeFirst(block, node, []string{
		"<c:numFmt", "<c:majorTickMark", "<c:minorTickMark", axisTickLabelTagPrefix,
		"<c:spPr", "<c:txPr", "<c:crossAx", "<c:crosses", chartElementClosePrefix,
	})
}

func insertAxisNumberFormat(block string, node string) string {
	return insertBeforeFirst(block, node, []string{
		"<c:majorTickMark", "<c:minorTickMark", axisTickLabelTagPrefix, "<c:spPr",
		"<c:txPr", "<c:crossAx", "<c:crosses", chartElementClosePrefix,
	})
}

func insertAxisUnit(block string, node string) string {
	for _, anchor := range []string{"<c:dispUnits", chartExtensionListTagPrefix} {
		if index := strings.Index(block, anchor); index >= 0 {
			return block[:index] + node + block[index:]
		}
	}
	if index := strings.LastIndex(block, chartElementClosePrefix); index >= 0 {
		return block[:index] + node + block[index:]
	}
	return block
}

func insertBeforeFirst(block string, node string, anchors []string) string {
	for _, anchor := range anchors {
		if index := strings.Index(block, anchor); index >= 0 {
			return block[:index] + node + block[index:]
		}
	}
	return block
}
