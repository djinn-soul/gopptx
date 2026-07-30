package chart

import (
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// Readback for the series-level formatting written by the patches alongside.

// parseSeriesFormatState reads back the formatting of every series that carries
// any, skipping series with neither a c:spPr nor a c:marker of their own.
func parseSeriesFormatState(xml string) []common.ChartSeriesFormat {
	formats := make([]common.ChartSeriesFormat, 0)
	for index, ser := range reSerBlocks.FindAllString(xml, -1) {
		format := common.ChartSeriesFormat{SeriesIndex: index}
		filled := parseSeriesShapeState(ser, &format)
		filled = parseSeriesMarkerState(ser, &format) || filled
		if match := reSeriesSmooth.FindString(ser); match != "" {
			smooth := strings.Contains(match, `val="1"`)
			format.Smooth = &smooth
			filled = true
		}
		if filled {
			formats = append(formats, format)
		}
	}
	if len(formats) == 0 {
		return nil
	}
	return formats
}

func parseSeriesShapeState(ser string, format *common.ChartSeriesFormat) bool {
	spPr, _ := seriesChildBlock(ser, reChartLinesShapeProps)
	if spPr == "" {
		return false
	}
	format.FillColor = existingMarkerColor(spPr, reChartLineColor, false)
	if line := parseChartLineFormat(spPr); line != nil {
		format.LineColor, format.LineWidthEMU = line.Color, line.WidthEMU
		format.LineDash, format.NoLine = line.Dash, line.None
	}
	return format.FillColor != nil || format.LineColor != nil ||
		format.LineWidthEMU != nil || format.LineDash != nil || format.NoLine != nil
}

func parseSeriesMarkerState(ser string, format *common.ChartSeriesFormat) bool {
	marker, _ := seriesChildBlock(ser, reSeriesMarker)
	if marker == "" {
		return false
	}
	format.MarkerSymbol, format.MarkerSize = markerCarryOver(marker, nil, nil)
	format.MarkerFillColor = existingMarkerColor(marker, reChartLineColor, false)
	format.MarkerLineColor = existingMarkerColor(marker, reChartLineColor, true)
	return format.MarkerSymbol != nil || format.MarkerSize != nil ||
		format.MarkerFillColor != nil || format.MarkerLineColor != nil
}

// markerCarryOver keeps the symbol and size a marker already had, so setting
// only its colour does not reset it to PowerPoint's automatic marker.
func markerCarryOver(current string, symbol *string, size *int) (*string, *int) {
	if symbol == nil {
		if match := reDataPointMarkerSym.FindStringSubmatch(current); len(match) == 2 {
			existing := match[1]
			symbol = &existing
		}
	}
	if size == nil {
		if match := reDataPointMarkerSize.FindStringSubmatch(current); len(match) == 2 {
			if existing, err := strconv.Atoi(match[1]); err == nil {
				size = &existing
			}
		}
	}
	return symbol, size
}

// existingMarkerColor reads the current fill or line colour out of a c:spPr, so
// rewriting the element keeps the half the caller did not set.
func existingMarkerColor(current string, re *regexp.Regexp, fromLine bool) *string {
	segment := current
	if line := reChartLineBlock.FindString(current); fromLine {
		if line == "" {
			return nil
		}
		segment = line
	} else if line != "" {
		segment = strings.Replace(current, line, "", 1)
	}
	match := re.FindStringSubmatch(segment)
	if len(match) != 2 {
		return nil
	}
	color := match[1]
	return &color
}
