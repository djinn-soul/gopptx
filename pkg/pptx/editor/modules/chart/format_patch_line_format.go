package chart

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

// Shared c:spPr line handling for the chart elements whose only formatting is a
// line: gridlines, series lines, and the outline of a data label.

var (
	reChartLineBlock = regexp.MustCompile(`(?s)<a:ln\b[^>]*(?:/>|>.*?</a:ln>)`)
	reChartLineWidth = regexp.MustCompile(`<a:ln\b[^>]*\bw="([^"]+)"`)
	reChartLineDash  = regexp.MustCompile(`<a:prstDash val="([^"]+)"`)
	reChartLineColor = regexp.MustCompile(`(?s)<a:solidFill>\s*<a:srgbClr val="([^"]+)"`)
	reChartLineNone  = regexp.MustCompile(`<a:noFill\s*/>`)

	reChartLinesShapeProps = regexp.MustCompile(`(?s)<c:spPr>.*?</c:spPr>|<c:spPr\s*/>`)
)

func validateChartLineFormat(field string, line *common.ChartLineFormat) error {
	if line == nil {
		return nil
	}
	if line.Color != nil {
		if _, err := editorshape.NormalizeHexColor(*line.Color); err != nil {
			return fmt.Errorf("%s color: %w", field, err)
		}
	}
	if line.Dash != nil {
		if _, err := editorshape.NormalizeLineDashStyle(*line.Dash); err != nil {
			return fmt.Errorf("%s dash: %w", field, err)
		}
	}
	if line.WidthEMU != nil && *line.WidthEMU < 0 {
		return fmt.Errorf("%s width_emu must not be negative, got %d", field, *line.WidthEMU)
	}
	return nil
}

// buildChartLine renders one a:ln. CT_LineProperties carries the width as an
// attribute and orders its children fill then dash; None wins over the rest,
// because a line with no fill is how PowerPoint hides one of these elements.
func buildChartLine(line *common.ChartLineFormat) string {
	if line == nil {
		return ""
	}
	if line.None != nil && *line.None {
		return "<a:ln><a:noFill/></a:ln>"
	}
	if line.Color == nil && line.WidthEMU == nil && line.Dash == nil {
		return ""
	}
	attrs := ""
	if line.WidthEMU != nil {
		attrs = ` w="` + strconv.Itoa(*line.WidthEMU) + `"`
	}
	var b strings.Builder
	b.WriteString("<a:ln" + attrs + ">")
	if color := normalizedColor(line.Color); color != "" {
		b.WriteString(`<a:solidFill><a:srgbClr val="` + color + `"/></a:solidFill>`)
	}
	if line.Dash != nil {
		if dash, err := editorshape.NormalizeLineDashStyle(*line.Dash); err == nil {
			b.WriteString(`<a:prstDash val="` + dash + `"/>`)
		}
	}
	b.WriteString("</a:ln>")
	return b.String()
}

// setShapePropertiesLine replaces the a:ln of an existing c:spPr, keeping the
// fill and effects around it. CT_ShapeProperties puts the line after the fill
// and before the effects, so a fresh line is spliced in at the first of those.
func setShapePropertiesLine(spPr string, rendered string) string {
	if rendered == "" {
		return spPr
	}
	if reChartLineBlock.MatchString(spPr) {
		return reChartLineBlock.ReplaceAllLiteralString(spPr, rendered)
	}
	for _, anchor := range []string{"<a:effectLst", "<a:scene3d", "<a:sp3d", "</c:spPr>"} {
		if index := strings.Index(spPr, anchor); index >= 0 {
			return spPr[:index] + rendered + spPr[index:]
		}
	}
	return spPr
}

// styleChartLinesElement puts a line into a CT_ChartLines element — a gridline
// or the series lines. c:spPr is its only child, so an empty element is
// expanded rather than appended to, and an existing c:spPr keeps what it has.
func styleChartLinesElement(block string, tag string, rendered string) string {
	openTag := "<c:" + tag + ">"
	if strings.HasSuffix(block, "/>") {
		return openTag + "<c:spPr>" + rendered + "</c:spPr></c:" + tag + ">"
	}
	if current := reChartLinesShapeProps.FindString(block); current != "" {
		if strings.HasSuffix(current, "/>") {
			return strings.Replace(block, current, "<c:spPr>"+rendered+"</c:spPr>", 1)
		}
		return strings.Replace(block, current, setShapePropertiesLine(current, rendered), 1)
	}
	return strings.Replace(block, openTag, openTag+"<c:spPr>"+rendered+"</c:spPr>", 1)
}

// parseChartLineFormat reads back the line style this API models. It returns
// nil when the element carries no explicit line at all.
func parseChartLineFormat(spPr string) *common.ChartLineFormat {
	line := reChartLineBlock.FindString(spPr)
	if line == "" {
		return nil
	}
	format := &common.ChartLineFormat{}
	if reChartLineNone.MatchString(line) {
		none := true
		format.None = &none
		return format
	}
	if match := reChartLineWidth.FindStringSubmatch(line); len(match) == 2 {
		if width, err := strconv.Atoi(match[1]); err == nil {
			format.WidthEMU = &width
		}
	}
	if match := reChartLineColor.FindStringSubmatch(line); len(match) == 2 {
		color := match[1]
		format.Color = &color
	}
	if match := reChartLineDash.FindStringSubmatch(line); len(match) == 2 {
		dash := match[1]
		format.Dash = &dash
	}
	if format.Color == nil && format.WidthEMU == nil && format.Dash == nil {
		return nil
	}
	return format
}

// chartFillNode renders the solid fill or explicit no-fill of a chart element.
func chartFillNode(color *string, noFill *bool) string {
	if noFill != nil && *noFill {
		return "<a:noFill/>"
	}
	if hex := normalizedColor(color); hex != "" {
		return `<a:solidFill><a:srgbClr val="` + hex + `"/></a:solidFill>`
	}
	return ""
}
