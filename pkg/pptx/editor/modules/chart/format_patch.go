package chart

import (
	"errors"
	"regexp"
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
	if err := ValidateAxisVisibility(req); err != nil {
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
	if err := validateChartTrendlines(
		req.Trendlines, req.AppendTrendlines, req.ClearTrendlineSeries,
	); err != nil {
		return err
	}
	if err := validateChartErrorBars(req.ErrorBars, req.ClearErrorBarSeries); err != nil {
		return err
	}
	if err := validateChartDataPoints(req.DataPoints, req.ClearDataPointSeries); err != nil {
		return err
	}
	if err := validateChartDataLabelPoints(req.DataLabelPoints); err != nil {
		return err
	}
	if err := validateChartDataLabelBox(req); err != nil {
		return err
	}
	return validateChartFormatLinesAndTables(req)
}

func validateChartFormatLinesAndTables(req common.ChartFormatUpdate) error {
	if err := validateGridlineFormats(req); err != nil {
		return err
	}
	if err := validateChartSeriesFormats(req.SeriesFormats); err != nil {
		return err
	}
	if err := validateChartSeriesLines(req.SeriesLines); err != nil {
		return err
	}
	if err := validateChartSeriesInverts(req.SeriesInverts); err != nil {
		return err
	}
	if err := validateChartDataTable(req.DataTable); err != nil {
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
	updated = patchChartDataLabelBox(updated, req)
	// After the chart-wide label settings: a per-label patch reads the series
	// flags it has to repeat, and must win over them.
	updated = patchChartDataLabelPoints(updated, req.DataLabelPoints)
	// Data points first: the invert patch derives more of them from the series
	// values, and both write into the same c:dPt run.
	updated = patchChartDataPoints(updated, req.DataPoints, req.ClearDataPointSeries)
	updated = patchChartSeriesInvert(updated, req.SeriesInverts)
	// After the per-point work: a series format writes the c:spPr and c:marker
	// that sit before the c:dPt run it must not disturb.
	updated = patchChartSeriesFormats(updated, req.SeriesFormats)
	updated = patchChartTrendlines(
		updated, req.Trendlines, req.AppendTrendlines, req.ClearTrendlineSeries,
	)
	updated = patchChartErrorBars(updated, req.ErrorBars, req.ClearErrorBarSeries)
	updated = patchChartPlotOptions(updated, req)
	// After gapWidth and overlap: CT_BarChart puts c:serLines behind both.
	updated = patchChartSeriesLines(updated, req.SeriesLines)
	updated = patchAxisTickLabelPosition(updated, "catAx", req.CategoryAxisTickLabelPos)
	updated = patchAxisTickLabelPosition(updated, "dateAx", req.CategoryAxisTickLabelPos)
	updated = patchAxisTickLabelPosition(updated, "valAx", req.ValueAxisTickLabelPos)
	updated = patchAxisMajorGridlines(updated, "catAx", req.CategoryAxisMajorGrid)
	updated = patchAxisMajorGridlines(updated, "dateAx", req.CategoryAxisMajorGrid)
	updated = patchAxisMajorGridlines(updated, "valAx", req.ValueAxisMajorGrid)
	updated = patchAxisMinorGridlines(updated, "catAx", req.CategoryAxisMinorGrid)
	updated = patchAxisMinorGridlines(updated, "dateAx", req.CategoryAxisMinorGrid)
	updated = patchAxisMinorGridlines(updated, "valAx", req.ValueAxisMinorGrid)
	// After the on/off patches: styling a gridline draws it, and both write the
	// same element.
	updated = patchAxisGridlineFormats(updated, req)
	updated = patchAxisCrosses(updated, "catAx", req.CategoryAxisCrosses)
	updated = patchAxisCrosses(updated, "dateAx", req.CategoryAxisCrosses)
	updated = patchAxisCrosses(updated, "valAx", req.ValueAxisCrosses)
	updated = patchAxisDetails(updated, req)
	updated = PatchAxisVisibility(updated, req)
	// After the axes: CT_PlotArea puts c:dTable behind them.
	updated = patchChartDataTable(updated, req.DataTable)
	updated = patchChartScene3D(updated, req)
	return []byte(updated), nil
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
		block = patchDataLabelFlags(block, req)
		block = patchDataLabelWordWrap(block, req.DataLabelWordWrap)
		return block
	})
}

// patchDataLabelFlags rewrites the display flags of one c:dLbls as a complete
// set in schema order. CT_DLbls puts them after numFmt, spPr, txPr and dLblPos
// and before c:separator, and a partial set makes PowerPoint fall back to its
// own defaults, which is how a per-point label loses its number format.
func patchDataLabelFlags(block string, req common.ChartFormatUpdate) string {
	overrides := map[string]*bool{
		flagShowLegendKey:  req.DataLabelShowLegendKey,
		flagShowValue:      req.DataLabelShowValue,
		flagShowCategory:   req.DataLabelShowCategory,
		flagShowSeriesName: req.DataLabelShowSeriesName,
		flagShowPercent:    req.DataLabelShowPercent,
		flagShowBubbleSize: req.DataLabelShowBubbleSize,
	}

	values := dataLabelFlagValues(block)
	present := len(values) > 0
	for name, override := range overrides {
		if override == nil {
			continue
		}
		values[name] = boolToOneZero(*override)
		present = true
	}
	if !present {
		return block
	}

	block = reDataLabelFlag.ReplaceAllLiteralString(block, "")
	var flags strings.Builder
	for _, name := range dataLabelFlagNames() {
		value, ok := values[name]
		if !ok {
			value = "0"
		}
		flags.WriteString(`<c:` + name + ` val="` + value + `"/>`)
	}

	if index := strings.Index(block, "<c:separator>"); index >= 0 {
		return block[:index] + flags.String() + block[index:]
	}
	return strings.Replace(block, "</c:dLbls>", flags.String()+"</c:dLbls>", 1)
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
