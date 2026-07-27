package chart

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const defaultChartNumberFormat = "General"

var (
	reChartTitleBlock = regexp.MustCompile(`(?s)<c:title>.*?</c:title>`)
	reTitleText       = regexp.MustCompile(`(?s)<a:t>.*?</a:t>`)
	reAutoTitleDelete = regexp.MustCompile(`<c:autoTitleDeleted val="[^"]*"/>`)
	reOverlay         = regexp.MustCompile(`<c:overlay val="[^"]*"/>`)
	reLegendBlock     = regexp.MustCompile(`(?s)<c:legend>.*?</c:legend>`)
	reLegendPos       = regexp.MustCompile(`<c:legendPos val="[^"]*"/>`)
	reDataLabelsBlock = regexp.MustCompile(`(?s)<c:dLbls>.*?</c:dLbls>`)
	reDataLabelPos    = regexp.MustCompile(`<c:dLblPos val="[^"]*"/>`)
	rePlotVisOnly     = regexp.MustCompile(`<c:plotVisOnly val="[^"]*"/>`)

	reChartTitleLayout = regexp.MustCompile(`(?s)<c:layout>.*?</c:layout>|<c:layout/>`)
	reLayoutX          = regexp.MustCompile(`<c:x val="([^"]*)"/>`)
	reLayoutY          = regexp.MustCompile(`<c:y val="([^"]*)"/>`)
)

func ValidateChartFormatUpdate(req common.ChartFormatUpdate) error {
	if req.LegendPosition != nil && !isLegendPosition(*req.LegendPosition) {
		return errors.New("legend_position must be one of r,l,t,b")
	}
	if req.DataLabelPosition != nil && !isDataLabelPosition(*req.DataLabelPosition) {
		return errors.New("data_label_position must be one of ctr,inEnd,inBase,outEnd,bestFit,l,r,t,b")
	}
	if err := validateAxisTickLabelPosition("category_axis_tick_label_pos", req.CategoryAxisTickLabelPos); err != nil {
		return err
	}
	if err := validateAxisTickLabelPosition("value_axis_tick_label_pos", req.ValueAxisTickLabelPos); err != nil {
		return err
	}
	if err := validateAxisCrosses("category_axis_crosses", req.CategoryAxisCrosses); err != nil {
		return err
	}
	if err := validateAxisCrosses("value_axis_crosses", req.ValueAxisCrosses); err != nil {
		return err
	}
	if err := validateAxisDetails(req); err != nil {
		return err
	}
	if err := validateChartPlotOptions(req); err != nil {
		return err
	}
	if err := validateScene3DUpdate(req); err != nil {
		return err
	}
	return nil
}

func PatchChartFormatting(chartXML []byte, req common.ChartFormatUpdate) ([]byte, error) {
	if err := ValidateChartFormatUpdate(req); err != nil {
		return nil, err
	}

	updated := string(chartXML)
	if err := validateAxisScaleAgainstXML(updated, req); err != nil {
		return nil, err
	}
	if req.ShowTitle != nil || req.Title != nil || req.TitleOverlay != nil ||
		req.TitleX != nil || req.TitleY != nil {
		var err error
		updated, err = patchChartTitle(
			updated, req.ShowTitle, req.Title, req.TitleOverlay, req.TitleX, req.TitleY,
		)
		if err != nil {
			return nil, err
		}
	}
	updated = patchPlotVisibleOnly(updated, req.PlotVisibleOnly)
	updated = patchChartLegend(updated, req.ShowLegend, req.LegendPosition, req.LegendOverlay)
	updated = patchChartDataLabels(updated, req)
	updated = patchDataLabelOffsets(updated, req.DataLabelOffsets)
	updated = patchChartDataLabelNumberFormat(updated, req.DataLabelNumberFormat, req.DataLabelFormatLinked)
	updated = patchChartPlotOptions(updated, req)
	updated = patchAxisTickLabelPosition(updated, "catAx", req.CategoryAxisTickLabelPos)
	updated = patchAxisTickLabelPosition(updated, "dateAx", req.CategoryAxisTickLabelPos)
	updated = patchAxisTickLabelPosition(updated, "valAx", req.ValueAxisTickLabelPos)
	updated = patchAxisMajorGridlines(updated, "catAx", req.CategoryAxisMajorGrid)
	updated = patchAxisMajorGridlines(updated, "dateAx", req.CategoryAxisMajorGrid)
	updated = patchAxisMajorGridlines(updated, "valAx", req.ValueAxisMajorGrid)
	updated = patchAxisMinorGridlines(updated, "catAx", req.CategoryAxisMinorGrid)
	updated = patchAxisMinorGridlines(updated, "dateAx", req.CategoryAxisMinorGrid)
	updated = patchAxisMinorGridlines(updated, "valAx", req.ValueAxisMinorGrid)
	updated = patchAxisCrosses(updated, "catAx", req.CategoryAxisCrosses)
	updated = patchAxisCrosses(updated, "dateAx", req.CategoryAxisCrosses)
	updated = patchAxisCrosses(updated, "valAx", req.ValueAxisCrosses)
	updated = patchAxisDetails(updated, req)
	updated = patchChartScene3D(updated, req)
	return []byte(updated), nil
}

func patchChartTitle(
	xml string,
	show *bool,
	title *string,
	overlay *bool,
	titleX, titleY *float64,
) (string, error) {
	match := reChartTitleBlock.FindString(xml)
	if show != nil && !*show {
		xml = strings.Replace(xml, match, "", 1)
		return setAutoTitleDeleted(xml, true), nil
	}
	if match == "" {
		if show != nil || title != nil || overlay != nil || titleX != nil || titleY != nil {
			inserted := insertChartTitleBlock(xml, title, overlay)
			return patchChartTitleLayoutIn(inserted, titleX, titleY), nil
		}
		return "", errors.New("chart title block not found")
	}
	block := patchExistingChartTitleBlock(match, title, overlay)
	block = patchChartTitleLayout(block, titleX, titleY)
	xml = strings.Replace(xml, match, block, 1)
	return setAutoTitleDeleted(xml, false), nil
}

func patchChartTitleLayoutIn(xml string, titleX, titleY *float64) string {
	if titleX == nil && titleY == nil {
		return xml
	}
	match := reChartTitleBlock.FindString(xml)
	if match == "" {
		return xml
	}
	return strings.Replace(xml, match, patchChartTitleLayout(match, titleX, titleY), 1)
}

// patchChartTitleLayout writes the manual title position. CT_Title orders its
// children tx, layout, overlay, spPr, txPr, so the layout is spliced after
// </c:tx> rather than appended.
func patchChartTitleLayout(block string, titleX, titleY *float64) string {
	if titleX == nil && titleY == nil {
		return block
	}

	existing := reChartTitleLayout.FindString(block)
	x, y := 0.0, 0.0
	if existing != "" {
		x, y = parseManualLayoutXY(existing)
	}
	if titleX != nil {
		x = *titleX
	}
	if titleY != nil {
		y = *titleY
	}

	layout := `<c:layout><c:manualLayout>` +
		`<c:xMode val="edge"/><c:yMode val="edge"/>` +
		`<c:x val="` + formatChartFraction(x) + `"/>` +
		`<c:y val="` + formatChartFraction(y) + `"/>` +
		`</c:manualLayout></c:layout>`

	if existing != "" {
		return strings.Replace(block, existing, layout, 1)
	}
	if strings.Contains(block, "</c:tx>") {
		return strings.Replace(block, "</c:tx>", "</c:tx>"+layout, 1)
	}
	return strings.Replace(block, "<c:title>", "<c:title>"+layout, 1)
}

func parseManualLayoutXY(layout string) (float64, float64) {
	var x, y float64
	if m := reLayoutX.FindStringSubmatch(layout); len(m) > 1 {
		x, _ = strconv.ParseFloat(m[1], 64)
	}
	if m := reLayoutY.FindStringSubmatch(layout); len(m) > 1 {
		y, _ = strconv.ParseFloat(m[1], 64)
	}
	return x, y
}

func formatChartFraction(v float64) string {
	return strconv.FormatFloat(v, 'g', -1, 64)
}

func insertChartTitleBlock(xml string, title *string, overlay *bool) string {
	titleText := "Chart"
	if title != nil {
		titleText = *title
	}
	overlayVal := "0"
	if overlay != nil {
		overlayVal = boolToOneZero(*overlay)
	}
	titleBlock := `<c:title><c:tx><c:rich><a:bodyPr/><a:lstStyle/><a:p><a:r><a:rPr lang="en-US"/><a:t>` +
		xmlEscape(titleText) +
		`</a:t></a:r></a:p></c:rich></c:tx><c:overlay val="` + overlayVal + `"/></c:title>`
	if strings.Contains(xml, "<c:autoTitleDeleted") {
		xml = strings.Replace(xml, "<c:autoTitleDeleted", titleBlock+"<c:autoTitleDeleted", 1)
	} else {
		xml = strings.Replace(xml, "<c:plotArea>", titleBlock+"<c:plotArea>", 1)
	}
	return setAutoTitleDeleted(xml, false)
}

func patchExistingChartTitleBlock(match string, title *string, overlay *bool) string {
	block := match
	if title != nil {
		block = patchChartTitleText(block, *title)
	}
	if overlay != nil {
		block = patchChartTitleOverlay(block, *overlay)
	}
	return block
}

func patchChartTitleText(block string, title string) string {
	escaped := `<a:t>` + xmlEscape(title) + `</a:t>`
	if reTitleText.MatchString(block) {
		return reTitleText.ReplaceAllLiteralString(block, escaped)
	}
	return strings.Replace(
		block,
		"</c:tx>",
		`<c:rich><a:bodyPr/><a:lstStyle/><a:p><a:r><a:rPr lang="en-US"/>`+escaped+`</a:r></a:p></c:rich></c:tx>`,
		1,
	)
}

func patchChartTitleOverlay(block string, overlay bool) string {
	overlayNode := `<c:overlay val="` + boolToOneZero(overlay) + `"/>`
	if reOverlay.MatchString(block) {
		return reOverlay.ReplaceAllString(block, overlayNode)
	}
	return strings.Replace(block, "</c:title>", overlayNode+"</c:title>", 1)
}

func setAutoTitleDeleted(xml string, deleted bool) string {
	node := `<c:autoTitleDeleted val="` + boolToOneZero(deleted) + `"/>`
	if reAutoTitleDelete.MatchString(xml) {
		return reAutoTitleDelete.ReplaceAllString(xml, node)
	}
	return strings.Replace(xml, "<c:plotArea>", node+"<c:plotArea>", 1)
}

func patchChartLegend(xml string, show *bool, position *string, overlay *bool) string {
	match := reLegendBlock.FindString(xml)
	if show != nil && !*show {
		if match == "" {
			return xml
		}
		return strings.Replace(xml, match, "", 1)
	}

	if match == "" && (show != nil || position != nil || overlay != nil) {
		legendPos := "r"
		if position != nil {
			legendPos = strings.TrimSpace(*position)
		}
		overlayVal := "0"
		if overlay != nil {
			overlayVal = boolToOneZero(*overlay)
		}
		legend := `<c:legend><c:legendPos val="` + legendPos + `"/><c:overlay val="` + overlayVal + `"/></c:legend>`
		return strings.Replace(xml, "<c:plotVisOnly", legend+"<c:plotVisOnly", 1)
	}
	if match == "" {
		return xml
	}

	block := match
	if position != nil {
		node := `<c:legendPos val="` + strings.TrimSpace(*position) + `"/>`
		if reLegendPos.MatchString(block) {
			block = reLegendPos.ReplaceAllString(block, node)
		} else {
			block = strings.Replace(block, "<c:legend>", "<c:legend>"+node, 1)
		}
	}
	if overlay != nil {
		overlayNode := `<c:overlay val="` + boolToOneZero(*overlay) + `"/>`
		if reOverlay.MatchString(block) {
			block = reOverlay.ReplaceAllString(block, overlayNode)
		} else {
			block = strings.Replace(block, "</c:legend>", overlayNode+"</c:legend>", 1)
		}
	}
	return strings.Replace(xml, match, block, 1)
}

func patchChartDataLabels(xml string, req common.ChartFormatUpdate) string {
	show := req.ShowDataLabels
	position := req.DataLabelPosition
	if show != nil && !*show {
		return reDataLabelsBlock.ReplaceAllString(xml, "")
	}

	hasLabels := reDataLabelsBlock.MatchString(xml)
	needLabels := (show != nil && *show) ||
		position != nil ||
		req.DataLabelShowLegendKey != nil ||
		req.DataLabelShowValue != nil ||
		req.DataLabelShowCategory != nil ||
		req.DataLabelShowSeriesName != nil ||
		req.DataLabelShowPercent != nil ||
		req.DataLabelShowBubbleSize != nil ||
		req.DataLabelWordWrap != nil
	if !hasLabels && needLabels {
		xml = insertDefaultDataLabels(xml)
	}
	isBarChart := strings.Contains(xml, "<c:barChart")
	return reDataLabelsBlock.ReplaceAllStringFunc(xml, func(block string) string {
		if position != nil {
			normPos := normalizeDataLabelPosition(*position, isBarChart)
			node := `<c:dLblPos val="` + normPos + `"/>`
			if reDataLabelPos.MatchString(block) {
				block = reDataLabelPos.ReplaceAllString(block, node)
			} else {
				block = strings.Replace(block, "<c:dLbls>", "<c:dLbls>"+node, 1)
			}
		}
		block = patchDataLabelBool(block, "showLegendKey", req.DataLabelShowLegendKey)
		block = patchDataLabelBool(block, "showVal", req.DataLabelShowValue)
		block = patchDataLabelBool(block, "showCatName", req.DataLabelShowCategory)
		block = patchDataLabelBool(block, "showSerName", req.DataLabelShowSeriesName)
		block = patchDataLabelBool(block, "showPercent", req.DataLabelShowPercent)
		block = patchDataLabelBool(block, "showBubbleSize", req.DataLabelShowBubbleSize)
		block = patchDataLabelWordWrap(block, req.DataLabelWordWrap)
		return block
	})
}

func patchDataLabelBool(block string, tag string, value *bool) string {
	if value == nil {
		return block
	}
	re := regexp.MustCompile(`<c:` + tag + ` val="[^"]*"/>`)
	if *value {
		node := `<c:` + tag + ` val="1"/>`
		if re.MatchString(block) {
			return re.ReplaceAllString(block, node)
		}
		return strings.Replace(block, "<c:dLbls>", "<c:dLbls>"+node, 1)
	}
	return re.ReplaceAllString(block, "")
}

func insertDefaultDataLabels(xml string) string {
	start, end := firstChartBlockRange(xml)
	if start < 0 || end <= start {
		return xml
	}
	chartBlock := xml[start:end]
	insertAt := strings.Index(chartBlock, "<c:axId")
	if insertAt < 0 {
		insertAt = strings.LastIndex(chartBlock, chartElementClosePrefix)
		if insertAt < 0 {
			return xml
		}
	}
	labels := `<c:dLbls><c:showVal val="1"/></c:dLbls>`
	patched := chartBlock[:insertAt] + labels + chartBlock[insertAt:]
	return xml[:start] + patched + xml[end:]
}

func firstChartBlockRange(xml string) (int, int) {
	return firstChartBlockBounds(xml)
}
