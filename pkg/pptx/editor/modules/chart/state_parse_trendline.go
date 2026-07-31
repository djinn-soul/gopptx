package chart

import (
	"html"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

var (
	reTrendlineType      = regexp.MustCompile(`<c:trendlineType val="([^"]+)"`)
	reTrendlineName      = regexp.MustCompile(`(?s)<c:name>(.*?)</c:name>`)
	reTrendlineOrder     = regexp.MustCompile(`<c:order val="([^"]+)"`)
	reTrendlinePeriod    = regexp.MustCompile(`<c:period val="([^"]+)"`)
	reTrendlineForward   = regexp.MustCompile(`<c:forward val="([^"]+)"`)
	reTrendlineBackward  = regexp.MustCompile(`<c:backward val="([^"]+)"`)
	reTrendlineIntercept = regexp.MustCompile(`<c:intercept val="([^"]+)"`)
	reTrendlineRSqr      = regexp.MustCompile(`<c:dispRSqr val="([^"]+)"`)
	reTrendlineEq        = regexp.MustCompile(`<c:dispEq val="([^"]+)"`)
	reTrendlineColor     = regexp.MustCompile(`<a:srgbClr val="([^"]+)"`)
	reTrendlineWidth     = regexp.MustCompile(`<a:ln\b[^>]*\bw="([^"]+)"`)
	reTrendlineDash      = regexp.MustCompile(`<a:prstDash val="([^"]+)"`)
)

// parseTrendlineState reads every c:trendline back into the update payload
// shape, so a chart snapshot round-trips what was written.
func parseTrendlineState(xml string) []common.ChartTrendline {
	out := make([]common.ChartTrendline, 0)
	for seriesIndex, ser := range reSerBlocks.FindAllString(xml, -1) {
		for _, block := range reTrendlineBlock.FindAllString(ser, -1) {
			out = append(out, parseTrendlineBlock(block, seriesIndex))
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func parseTrendlineBlock(block string, seriesIndex int) common.ChartTrendline {
	trendline := common.ChartTrendline{SeriesIndex: seriesIndex}
	if match := reTrendlineType.FindStringSubmatch(block); len(match) == expectedSingleGroupMatch {
		trendline.Type = strings.TrimSpace(match[1])
	}
	if match := reTrendlineName.FindStringSubmatch(block); len(match) == expectedSingleGroupMatch {
		name := html.UnescapeString(match[1])
		trendline.Name = &name
	}
	trendline.Order = trendlineIntValue(block, reTrendlineOrder)
	trendline.Period = trendlineIntValue(block, reTrendlinePeriod)
	trendline.Forward = axisFloatValue(block, reTrendlineForward)
	trendline.Backward = axisFloatValue(block, reTrendlineBackward)
	trendline.Intercept = axisFloatValue(block, reTrendlineIntercept)
	trendline.DisplayRSquared = trendlineBoolValue(block, reTrendlineRSqr)
	trendline.DisplayEquation = trendlineBoolValue(block, reTrendlineEq)
	trendline.LineColor = trendlineStringValue(block, reTrendlineColor)
	trendline.LineWidthEMU = trendlineIntValue(block, reTrendlineWidth)
	trendline.LineDash = trendlineStringValue(block, reTrendlineDash)
	return trendline
}

func trendlineIntValue(block string, re *regexp.Regexp) *int {
	match := re.FindStringSubmatch(block)
	if len(match) != expectedSingleGroupMatch {
		return nil
	}
	value, err := strconv.Atoi(strings.TrimSpace(match[1]))
	if err != nil {
		return nil
	}
	return &value
}

func trendlineBoolValue(block string, re *regexp.Regexp) *bool {
	match := re.FindStringSubmatch(block)
	if len(match) != expectedSingleGroupMatch {
		return nil
	}
	switch strings.ToLower(strings.TrimSpace(match[1])) {
	case "1", "true":
		value := true
		return &value
	case "0", "false":
		value := false
		return &value
	default:
		return nil
	}
}

func trendlineStringValue(block string, re *regexp.Regexp) *string {
	match := re.FindStringSubmatch(block)
	if len(match) != expectedSingleGroupMatch {
		return nil
	}
	value := strings.TrimSpace(match[1])
	return &value
}
