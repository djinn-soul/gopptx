package chart

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

var (
	reDataLabelNumberFormat = regexp.MustCompile(`<c:numFmt [^>]*/>`)
	reChartGrouping         = regexp.MustCompile(`<c:grouping val="[^"]*"/>`)
	reChartGapWidth         = regexp.MustCompile(`<c:gapWidth val="[^"]*"/>`)
	reChartOverlap          = regexp.MustCompile(`<c:overlap val="[^"]*"/>`)
)

func validateChartPlotOptions(req common.ChartFormatUpdate) error {
	if req.DataLabelNumberFormat != nil && strings.TrimSpace(*req.DataLabelNumberFormat) == "" {
		return errors.New("data_label_number_format must not be empty")
	}
	if req.ChartGrouping != nil && !isChartGrouping(*req.ChartGrouping) {
		return errors.New("chart_grouping must be one of clustered,stacked,percentStacked,standard")
	}
	if req.GapWidth != nil && (*req.GapWidth < 0 || *req.GapWidth > 500) {
		return errors.New("gap_width must be between 0 and 500")
	}
	if req.Overlap != nil && (*req.Overlap < -100 || *req.Overlap > 100) {
		return errors.New("overlap must be between -100 and 100")
	}
	return nil
}

func isChartGrouping(value string) bool {
	switch strings.TrimSpace(value) {
	case "clustered", "stacked", "percentStacked", "standard":
		return true
	default:
		return false
	}
}

func patchChartDataLabelNumberFormat(xml string, format *string, linked *bool) string {
	if format == nil && linked == nil {
		return xml
	}
	match := reDataLabelsBlock.FindString(xml)
	if match == "" {
		xml = insertDefaultDataLabels(xml)
		match = reDataLabelsBlock.FindString(xml)
	}
	if match == "" {
		return xml
	}
	formatCode, sourceLinked := defaultChartNumberFormat, true
	if current := reDataLabelNumberFormat.FindString(match); current != "" {
		formatCode = attributeValue(current, "formatCode", formatCode)
		sourceLinked = attributeValue(current, "sourceLinked", "1") == "1"
	}
	if format != nil {
		formatCode = *format
	}
	if linked != nil {
		sourceLinked = *linked
	}
	node := `<c:numFmt formatCode="` + xmlEscape(formatCode) + `" sourceLinked="` + boolToOneZero(sourceLinked) + `"/>`
	if reDataLabelNumberFormat.MatchString(match) {
		match = reDataLabelNumberFormat.ReplaceAllString(match, node)
	} else {
		match = strings.Replace(match, "<c:dLbls>", "<c:dLbls>"+node, 1)
	}
	return strings.Replace(xml, reDataLabelsBlock.FindString(xml), match, 1)
}

func patchChartPlotOptions(xml string, req common.ChartFormatUpdate) string {
	if req.ChartGrouping == nil && req.GapWidth == nil && req.Overlap == nil {
		return xml
	}
	start, end := firstChartBlockBounds(xml)
	if start < 0 || end <= start {
		return xml
	}
	block := xml[start:end]
	if req.ChartGrouping != nil {
		block = patchChartOption(block, reChartGrouping, "grouping", *req.ChartGrouping)
	}
	if req.GapWidth != nil {
		block = patchChartOption(block, reChartGapWidth, "gapWidth", strconv.Itoa(*req.GapWidth))
	}
	if req.Overlap != nil {
		block = patchChartOption(block, reChartOverlap, "overlap", strconv.Itoa(*req.Overlap))
	}
	return xml[:start] + block + xml[end:]
}

func patchChartOption(block string, re *regexp.Regexp, tag, value string) string {
	node := `<c:` + tag + ` val="` + value + `"/>`
	if re.MatchString(block) {
		return re.ReplaceAllString(block, node)
	}
	insertAt := strings.Index(block, "<c:ser>")
	if insertAt >= 0 {
		return block[:insertAt] + node + block[insertAt:]
	}
	return strings.Replace(block, ">", ">"+node, 1)
}
