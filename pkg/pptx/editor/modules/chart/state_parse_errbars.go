package chart

import (
	"html"
	"regexp"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

var (
	reErrBarType     = regexp.MustCompile(`<c:errBarType val="([^"]+)"`)
	reErrValType     = regexp.MustCompile(`<c:errValType val="([^"]+)"`)
	reErrDir         = regexp.MustCompile(`<c:errDir val="([^"]+)"`)
	reErrNoEndCap    = regexp.MustCompile(`<c:noEndCap val="([^"]+)"`)
	reErrValue       = regexp.MustCompile(`<c:val val="([^"]+)"`)
	reErrPlusFormula = regexp.MustCompile(`(?s)<c:plus>.*?<c:f>(.*?)</c:f>.*?</c:plus>`)
	reErrMinusFormla = regexp.MustCompile(`(?s)<c:minus>.*?<c:f>(.*?)</c:f>.*?</c:minus>`)
)

// parseErrorBarState reads every c:errBars back into the update payload shape,
// so a chart snapshot round-trips what was written.
func parseErrorBarState(xml string) []common.ChartErrorBars {
	out := make([]common.ChartErrorBars, 0)
	for seriesIndex, ser := range reSerBlocks.FindAllString(xml, -1) {
		for _, block := range reErrBarsBlock.FindAllString(ser, -1) {
			out = append(out, parseErrorBarBlock(block, seriesIndex))
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func parseErrorBarBlock(block string, seriesIndex int) common.ChartErrorBars {
	bars := common.ChartErrorBars{SeriesIndex: seriesIndex}
	if match := reErrBarType.FindStringSubmatch(block); len(match) == expectedSingleGroupMatch {
		bars.BarType = strings.TrimSpace(match[1])
	}
	if match := reErrValType.FindStringSubmatch(block); len(match) == expectedSingleGroupMatch {
		bars.ValueType = strings.TrimSpace(match[1])
	}
	bars.Direction = trendlineStringValue(block, reErrDir)
	bars.NoEndCap = trendlineBoolValue(block, reErrNoEndCap)
	bars.Value = axisFloatValue(block, reErrValue)
	bars.PlusReference = errorBarFormula(block, reErrPlusFormula)
	bars.MinusReference = errorBarFormula(block, reErrMinusFormla)
	bars.LineColor = trendlineStringValue(block, reTrendlineColor)
	return bars
}

func errorBarFormula(block string, re *regexp.Regexp) *string {
	match := re.FindStringSubmatch(block)
	if len(match) != expectedSingleGroupMatch {
		return nil
	}
	value := html.UnescapeString(strings.TrimSpace(match[1]))
	return &value
}
