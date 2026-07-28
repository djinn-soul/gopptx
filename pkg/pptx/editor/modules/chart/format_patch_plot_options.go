package chart

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

var (
	reDataLabelNumberFormat = regexp.MustCompile(`<c:numFmt [^>]*/>`)
	reChartGrouping         = regexp.MustCompile(`<c:grouping val="[^"]*"/>`)
	reChartGapWidth         = regexp.MustCompile(`<c:gapWidth val="[^"]*"/>`)
	reChartOverlap          = regexp.MustCompile(`<c:overlap val="[^"]*"/>`)
	reSeriesInvert          = regexp.MustCompile(`<c:invertIfNegative val="[^"]*"/>`)
	reSeriesValBlock        = regexp.MustCompile(`(?s)<c:val>.*?</c:val>`)
	reSeriesValuePoint      = regexp.MustCompile(`(?s)<c:pt idx="(\d+)"[^>]*>\s*<c:v>(-?[\d.eE+]+)</c:v>`)
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
		// A linked label takes its format from the source cells and ignores the
		// code, so asking for a code means unlinking unless the caller says
		// otherwise.
		sourceLinked = false
	}
	if linked != nil {
		sourceLinked = *linked
	}
	node := `<c:numFmt formatCode="` + xmlEscape(formatCode) + `" sourceLinked="` + boolToOneZero(sourceLinked) + `"/>`
	if reDataLabelNumberFormat.MatchString(match) {
		match = reDataLabelNumberFormat.ReplaceAllLiteralString(match, node)
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
		block = patchChartOptionBefore(
			block, reChartGrouping, "grouping", *req.ChartGrouping,
			[]string{"<c:varyColors", "<c:ser>", seriesDataLabelsTag, "<c:gapWidth", "<c:overlap"},
		)
	}
	if req.GapWidth != nil {
		block = patchChartOptionBefore(
			block, reChartGapWidth, "gapWidth", strconv.Itoa(*req.GapWidth),
			[]string{"<c:overlap", "<c:serLines", "<c:axId", chartExtensionListTagPrefix, chartElementClosePrefix},
		)
	}
	if req.Overlap != nil {
		block = patchChartOptionBefore(
			block, reChartOverlap, "overlap", strconv.Itoa(*req.Overlap),
			[]string{"<c:serLines", "<c:axId", chartExtensionListTagPrefix, chartElementClosePrefix},
		)
	}
	return xml[:start] + block + xml[end:]
}

func validateChartSeriesInverts(inverts []common.ChartSeriesInvert) error {
	for _, invert := range inverts {
		if invert.SeriesIndex < 0 {
			return errors.New("series_invert_if_negative series_index must not be negative")
		}
		if invert.NegativeFillColor == nil {
			continue
		}
		if _, err := editorshape.NormalizeHexColor(*invert.NegativeFillColor); err != nil {
			return errors.New("series_invert_if_negative negative_fill_color: " + err.Error())
		}
	}
	return nil
}

// patchChartSeriesInvert writes the series-level c:invertIfNegative flag and,
// when a negative fill colour is supplied, a c:dPt for each negative point.
//
// The flag on its own makes PowerPoint paint negative points with the fill's
// inverse — white against a default fill — so the explicit per-point fill is
// what actually makes them readable.
func patchChartSeriesInvert(xml string, inverts []common.ChartSeriesInvert) string {
	if len(inverts) == 0 {
		return xml
	}
	bySeries := map[int]common.ChartSeriesInvert{}
	for _, invert := range inverts {
		bySeries[invert.SeriesIndex] = invert
	}

	seriesIndex := -1
	return reSerBlocks.ReplaceAllStringFunc(xml, func(ser string) string {
		seriesIndex++
		invert, ok := bySeries[seriesIndex]
		if !ok {
			return ser
		}
		ser = writeSeriesInvertFlag(ser, invert.InvertIfNegative)
		if invert.NegativeFillColor == nil {
			return ser
		}
		return applySeriesDataPoints(ser, negativePointFormatting(ser, invert), false)
	})
}

// writeSeriesInvertFlag replaces or inserts c:invertIfNegative. CT_BarSer puts
// it after c:spPr and before c:pictureOptions, c:dPt, and c:dLbls.
func writeSeriesInvertFlag(ser string, enabled bool) string {
	node := `<c:invertIfNegative val="` + boolToOneZero(enabled) + `"/>`
	if reSeriesInvert.MatchString(ser) {
		return reSeriesInvert.ReplaceAllLiteralString(ser, node)
	}
	for _, anchor := range []string{
		"<c:pictureOptions", "<c:dPt>", seriesDataLabelsTag, seriesTrendlineTag, seriesErrorBarsTag,
		seriesCategoryTag, seriesXValuesTag, seriesValuesTag,
	} {
		if index := strings.Index(ser, anchor); index >= 0 {
			return ser[:index] + node + ser[index:]
		}
	}
	if index := strings.LastIndex(ser, seriesCloseTag); index >= 0 {
		return ser[:index] + node + ser[index:]
	}
	return ser
}

// negativePointFormatting builds one data point per negative value. Each one
// turns its own inversion off, so the explicit fill is what gets drawn.
func negativePointFormatting(ser string, invert common.ChartSeriesInvert) []common.ChartDataPoint {
	pointInvert := false
	points := make([]common.ChartDataPoint, 0)
	for index, value := range seriesValuesByIndex(ser) {
		if value >= 0 {
			continue
		}
		points = append(points, common.ChartDataPoint{
			SeriesIndex:      invert.SeriesIndex,
			PointIndex:       index,
			FillColor:        invert.NegativeFillColor,
			InvertIfNegative: &pointInvert,
		})
	}
	return points
}

// seriesValuesByIndex reads the cached c:val numbers of one series.
func seriesValuesByIndex(ser string) map[int]float64 {
	values := map[int]float64{}
	block := reSeriesValBlock.FindString(ser)
	if block == "" {
		return values
	}
	for _, match := range reSeriesValuePoint.FindAllStringSubmatch(block, -1) {
		index, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}
		value, err := strconv.ParseFloat(match[2], 64)
		if err != nil {
			continue
		}
		values[index] = value
	}
	return values
}

func patchChartOptionBefore(
	block string,
	re *regexp.Regexp,
	tag string,
	value string,
	anchors []string,
) string {
	node := `<c:` + tag + ` val="` + value + `"/>`
	if re.MatchString(block) {
		return re.ReplaceAllString(block, node)
	}
	for _, anchor := range anchors {
		if insertAt := strings.Index(block, anchor); insertAt >= 0 {
			return block[:insertAt] + node + block[insertAt:]
		}
	}
	return block
}
