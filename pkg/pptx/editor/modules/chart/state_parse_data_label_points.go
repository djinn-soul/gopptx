package chart

import (
	"html"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// parseDataLabelPointState reads back every c:dLbl that carries its own number
// format, font or display flags, so a caller can read the label of one point
// rather than only the chart-wide settings (upstream #803).
func parseDataLabelPointState(xml string) []common.ChartDataLabelPoint {
	out := make([]common.ChartDataLabelPoint, 0)
	for seriesIndex, ser := range reSerBlocks.FindAllString(xml, -1) {
		for _, block := range dataLabelBlocks(ser) {
			point := parseDataLabelPoint(block)
			point.SeriesIndex = seriesIndex
			out = append(out, point)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// dataLabelBlocks returns the per-point <c:dLbl> elements of one series.
func dataLabelBlocks(ser string) []string {
	var blocks []string
	cursor := 0
	for {
		rel := strings.Index(ser[cursor:], "<c:dLbl>")
		if rel < 0 {
			return blocks
		}
		start := cursor + rel
		closeRel := strings.Index(ser[start:], "</c:dLbl>")
		if closeRel < 0 {
			return blocks
		}
		end := start + closeRel + len("</c:dLbl>")
		blocks = append(blocks, ser[start:end])
		cursor = end
	}
}

// parseDataLabelFont reads the font properties this API models out of a label's
// c:txPr.
func parseDataLabelFont(textPr string, point *common.ChartDataLabelPoint) {
	if textPr == "" {
		return
	}
	if size := labelAttribute(textPr, reDataLabelFontSize); size != "" {
		if value, err := strconv.Atoi(size); err == nil {
			points := value / pointsPerFontSizeUnit
			point.FontSizePt = &points
		}
	}
	if bold := labelAttribute(textPr, reDataLabelFontBold); bold != "" {
		value := bold == "1"
		point.FontBold = &value
	}
	if color := labelAttribute(textPr, reDataLabelFontColor); color != "" {
		point.FontColor = &color
	}
}

func parseDataLabelPoint(block string) common.ChartDataLabelPoint {
	point := common.ChartDataLabelPoint{}
	if index := trendlineIntValue(block, reDataLabelIdxValue); index != nil {
		point.PointIndex = *index
	}
	if numFmt := reDataLabelNumFmt.FindString(block); numFmt != "" {
		if format := labelAttribute(numFmt, reDataLabelFormatVal); format != "" {
			decoded := html.UnescapeString(format)
			point.NumberFormat = &decoded
			linked := labelAttribute(numFmt, reDataLabelLinkedVal) == "1"
			point.FormatLinked = &linked
		}
	}
	parseDataLabelFont(reDataLabelTextPr.FindString(block), &point)
	values := dataLabelFlagValues(block)
	for name, target := range map[string]**bool{
		flagShowLegendKey:  &point.ShowLegendKey,
		flagShowValue:      &point.ShowValue,
		flagShowCategory:   &point.ShowCategory,
		flagShowSeriesName: &point.ShowSeriesName,
		flagShowPercent:    &point.ShowPercent,
	} {
		value, ok := values[name]
		if !ok {
			continue
		}
		flag := value == "1"
		*target = &flag
	}
	if deleted := reDataLabelDelete.FindString(block); deleted != "" {
		value := !strings.Contains(deleted, `val="0"`)
		point.Delete = &value
	}
	return point
}
