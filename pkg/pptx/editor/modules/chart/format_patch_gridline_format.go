package chart

import (
	"regexp"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// Gridline styling, upstream #984: the on/off patches above write a bare
// c:majorGridlines or c:minorGridlines, and PowerPoint then draws it in the
// theme's default grey. The colour, width and dash of that line live in the
// gridline's own c:spPr, which is what these patches write.

var (
	reMajorGridlinesBlock = regexp.MustCompile(`(?s)<c:majorGridlines(?:\s*/>|>.*?</c:majorGridlines>)`)
	reMinorGridlinesBlock = regexp.MustCompile(`(?s)<c:minorGridlines(?:\s*/>|>.*?</c:minorGridlines>)`)
)

func validateGridlineFormats(req common.ChartFormatUpdate) error {
	for _, pair := range []struct {
		value *common.ChartLineFormat
		name  string
	}{
		{req.CategoryAxisMajorGridFormat, "category_axis_major_gridline_format"},
		{req.CategoryAxisMinorGridFormat, "category_axis_minor_gridline_format"},
		{req.ValueAxisMajorGridFormat, "value_axis_major_gridline_format"},
		{req.ValueAxisMinorGridFormat, "value_axis_minor_gridline_format"},
	} {
		if err := validateChartLineFormat(pair.name, pair.value); err != nil {
			return err
		}
	}
	return nil
}

// patchAxisGridlineFormats writes every requested gridline style. The category
// axis is patched on both of the elements it can be, as elsewhere here.
func patchAxisGridlineFormats(xml string, req common.ChartFormatUpdate) string {
	categoryTags := []string{axisTagCategory, axisTagDate}
	valueTags := []string{axisTagValue}
	for _, target := range []struct {
		format   *common.ChartLineFormat
		axisTags []string
		major    bool
	}{
		{req.CategoryAxisMajorGridFormat, categoryTags, true},
		{req.CategoryAxisMinorGridFormat, categoryTags, false},
		{req.ValueAxisMajorGridFormat, valueTags, true},
		{req.ValueAxisMinorGridFormat, valueTags, false},
	} {
		if target.format == nil {
			continue
		}
		for _, axisTag := range target.axisTags {
			xml = patchAxisGridlineFormat(xml, axisTag, target.major, target.format)
		}
	}
	return xml
}

func patchAxisGridlineFormat(
	xml string,
	axisTag string,
	major bool,
	format *common.ChartLineFormat,
) string {
	rendered := buildChartLine(format)
	if rendered == "" {
		return xml
	}
	tag, re := "minorGridlines", reMinorGridlinesBlock
	if major {
		tag, re = "majorGridlines", reMajorGridlinesBlock
	}
	return patchEachAxisBlock(xml, axisTag, func(block string) string {
		current := re.FindString(block)
		if current == "" {
			// Styling a gridline the axis does not have implies drawing it;
			// there is nowhere else to put the c:spPr.
			insertAt := axisNodeInsertIndex(block)
			if insertAt < 0 {
				return block
			}
			node := "<c:" + tag + "><c:spPr>" + rendered + "</c:spPr></c:" + tag + ">"
			return block[:insertAt] + node + block[insertAt:]
		}
		return strings.Replace(block, current, styleChartLinesElement(current, tag, rendered), 1)
	})
}

// parseAxisGridlineFormats reads back the styles of one axis block.
func parseAxisGridlineFormats(block string, state *common.ChartAxisState) {
	if major := reMajorGridlinesBlock.FindString(block); major != "" {
		state.MajorGridlineFormat = parseChartLineFormat(major)
	}
	if minor := reMinorGridlinesBlock.FindString(block); minor != "" {
		state.MinorGridlineFormat = parseChartLineFormat(minor)
	}
}
