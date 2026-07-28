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
		"<c:numFmt", axisMajorTickTag, axisMinorTickTag, axisTickLabelTagPrefix,
		axisShapePropsTag, axisTextPropsTag, axisCrossAxTag, axisCrossesTag, chartElementClosePrefix,
	})
}

func insertAxisNumberFormat(block string, node string) string {
	return insertBeforeFirst(block, node, []string{
		axisMajorTickTag, axisMinorTickTag, axisTickLabelTagPrefix, axisShapePropsTag,
		axisTextPropsTag, axisCrossAxTag, axisCrossesTag, chartElementClosePrefix,
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

func insertAxisLabelAlignment(block string, node string) string {
	return insertBeforeFirstOrClose(block, node, "catAx", []string{
		"<c:lblOffset", "<c:tickLblSkip", "<c:tickMarkSkip", "<c:noMultiLvlLbl",
		chartExtensionListTagPrefix,
	})
}

func insertAxisTickMarkSkip(block string, node string) string {
	return insertBeforeFirstOrClose(block, node, "catAx", []string{
		"<c:noMultiLvlLbl", chartExtensionListTagPrefix,
	})
}

// insertBeforeFirstOrClose splices node before the first matching anchor, or
// before the block's own closing tag when no later sibling exists. The closing
// tag is matched explicitly because a bare "</c:" would find a nested child.
func insertBeforeFirstOrClose(block string, node string, blockTag string, anchors []string) string {
	for _, anchor := range anchors {
		if index := strings.Index(block, anchor); index >= 0 {
			return block[:index] + node + block[index:]
		}
	}
	closeTag := chartElementClosePrefix + blockTag + ">"
	if index := strings.LastIndex(block, closeTag); index >= 0 {
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
