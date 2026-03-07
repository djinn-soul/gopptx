package chart

import (
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

var (
	reChartStyle = regexp.MustCompile(`<c:style val="(\d+)"`)
	reSerBlock   = regexp.MustCompile(`(?s)<c:ser>.*?</c:ser>`)
	reAText      = regexp.MustCompile(`(?s)<a:t>(.*?)</a:t>`)
	reNumValue   = regexp.MustCompile(`(?s)<c:v>(-?\d+(?:\.\d+)?)</c:v>`)
	reTickLblPos = regexp.MustCompile(`<c:tickLblPos val="([^"]+)"`)
	reCrosses    = regexp.MustCompile(`<c:crosses val="([^"]+)"`)
)

const expectedSingleGroupMatch = 2

// ExtractChartState parses the chart XML cache into a traversal-friendly snapshot.
func ExtractChartState(chartXML []byte) common.ChartState {
	xml := string(chartXML)
	state := common.ChartState{
		CategoryAx: buildAxisState(xml, []string{"catAx", "dateAx"}),
		ValueAx:    buildAxisState(xml, []string{"valAx"}),
		Series:     parseSeriesState(xml),
	}
	if match := reChartStyle.FindStringSubmatch(xml); len(match) == expectedSingleGroupMatch {
		if style, err := strconv.Atoi(match[1]); err == nil {
			state.ChartStyle = &style
		}
	}
	return state
}

func buildAxisState(xml string, tags []string) common.ChartAxisState {
	state := common.ChartAxisState{}
	for _, tag := range tags {
		startTag := "<c:" + tag + ">"
		endTag := "</c:" + tag + ">"
		start := strings.Index(xml, startTag)
		if start < 0 {
			continue
		}
		endRel := strings.Index(xml[start:], endTag)
		if endRel < 0 {
			continue
		}
		block := xml[start : start+endRel+len(endTag)]
		state.Present = true
		if match := reTickLblPos.FindStringSubmatch(block); len(match) == expectedSingleGroupMatch {
			state.TickLabelPos = strings.TrimSpace(match[1])
		}
		if match := reCrosses.FindStringSubmatch(block); len(match) == expectedSingleGroupMatch {
			state.Crosses = strings.TrimSpace(match[1])
		}
		state.MajorGridline = strings.Contains(block, "<c:majorGridlines")
		break
	}
	return state
}

func parseSeriesState(xml string) []common.ChartSeriesData {
	matches := reSerBlock.FindAllString(xml, -1)
	out := make([]common.ChartSeriesData, 0, len(matches))
	for _, block := range matches {
		series := common.ChartSeriesData{}
		if nameMatch := reAText.FindStringSubmatch(block); len(nameMatch) == expectedSingleGroupMatch {
			name := strings.TrimSpace(nameMatch[1])
			series.Name = &name
		}
		values := make([]float64, 0)
		for _, valueMatch := range reNumValue.FindAllStringSubmatch(block, -1) {
			if len(valueMatch) != expectedSingleGroupMatch {
				continue
			}
			v, err := strconv.ParseFloat(valueMatch[1], 64)
			if err != nil {
				continue
			}
			values = append(values, v)
		}
		if len(values) > 0 {
			series.Values = values
		}
		out = append(out, series)
	}
	return out
}
