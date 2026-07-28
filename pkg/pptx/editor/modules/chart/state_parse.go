package chart

import (
	"html"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

var (
	reChartStyle      = regexp.MustCompile(`<c:style val="(\d+)"`)
	reSerBlock        = regexp.MustCompile(`(?s)<c:ser>.*?</c:ser>`)
	reAText           = regexp.MustCompile(`(?s)<a:t>(.*?)</a:t>`)
	reNumValue        = regexp.MustCompile(`(?s)<c:v>(-?\d+(?:\.\d+)?)</c:v>`)
	reTickLblPos      = regexp.MustCompile(`<c:tickLblPos val="([^"]+)"`)
	reCrosses         = regexp.MustCompile(`<c:crosses val="([^"]+)"`)
	reStateAxisTitle  = regexp.MustCompile(`(?s)<c:title>.*?<a:t>(.*?)</a:t>.*?</c:title>`)
	reStateAxisMin    = regexp.MustCompile(`<c:min val="([^"]+)"`)
	reStateAxisMax    = regexp.MustCompile(`<c:max val="([^"]+)"`)
	reStateAxisMajor  = regexp.MustCompile(`<c:majorUnit val="([^"]+)"`)
	reStateAxisMinor  = regexp.MustCompile(`<c:minorUnit val="([^"]+)"`)
	reStateAxisNumFmt = regexp.MustCompile(`<c:numFmt\b[^>]*formatCode="([^"]*)"[^>]*sourceLinked="([^"]+)"[^>]*/>`)
	reStateTickSkip   = regexp.MustCompile(`<c:tickMarkSkip val="([^"]+)"`)
	reStateLblAlgn    = regexp.MustCompile(`<c:lblAlgn val="([^"]+)"`)
	reScene3D         = regexp.MustCompile(`(?s)<a:scene3d\b.*?</a:scene3d>`)
	reCamera          = regexp.MustCompile(`<a:camera\b[^>]*prst="([^"]+)"[^>]*>`)
	reCameraFOV       = regexp.MustCompile(`<a:camera\b[^>]*fov="([^"]+)"[^>]*>`)
	reLightRig        = regexp.MustCompile(`<a:lightRig\b[^>]*rig="([^"]+)"[^>]*>`)
	reLightDir        = regexp.MustCompile(`<a:lightRig\b[^>]*dir="([^"]+)"[^>]*>`)
	reLightRev        = regexp.MustCompile(`<a:lightRig\b[^>]*rev="([^"]+)"[^>]*>`)
)

const expectedAxisNumFmtMatch = 3

const expectedSingleGroupMatch = 2

// ExtractChartState parses the chart XML cache into a traversal-friendly snapshot.
func ExtractChartState(chartXML []byte) common.ChartState {
	xml := string(chartXML)
	state := common.ChartState{
		CategoryAx: buildAxisState(xml, []string{"catAx", "dateAx"}),
		ValueAx:    buildAxisState(xml, []string{"valAx"}),
		Series:     parseSeriesState(xml),
		Scene3D:    parseScene3DState(xml),
		DataLabels: parseDataLabelState(xml),
		Trendlines: parseTrendlineState(xml),
		ErrorBars:  parseErrorBarState(xml),
		DataPoints: parseDataPointState(xml),
		DataTable:  parseDataTableState(xml),
		// Per-label reads: a point's number format lives on its c:dLbl, not on
		// the chart-wide c:dLbls.
		DataLabelPoints: parseDataLabelPointState(xml),
	}
	if match := reChartStyle.FindStringSubmatch(xml); len(match) == expectedSingleGroupMatch {
		if style, err := strconv.Atoi(match[1]); err == nil {
			state.ChartStyle = &style
		}
	}
	return state
}

func parseScene3DState(xml string) common.ChartScene3DState {
	match := reScene3D.FindString(xml)
	if match == "" {
		return common.ChartScene3DState{}
	}
	state := common.ChartScene3DState{}
	if m := reCamera.FindStringSubmatch(match); len(m) == expectedSingleGroupMatch {
		state.CameraPreset = strings.TrimSpace(m[1])
	}
	if m := reCameraFOV.FindStringSubmatch(match); len(m) == expectedSingleGroupMatch {
		if fov, err := strconv.Atoi(strings.TrimSpace(m[1])); err == nil {
			state.CameraFieldOfView = fov
		}
	}
	if m := reLightRig.FindStringSubmatch(match); len(m) == expectedSingleGroupMatch {
		state.LightRig = strings.TrimSpace(m[1])
	}
	if m := reLightDir.FindStringSubmatch(match); len(m) == expectedSingleGroupMatch {
		state.LightDirection = strings.TrimSpace(m[1])
	}
	if m := reLightRev.FindStringSubmatch(match); len(m) == expectedSingleGroupMatch {
		state.LightRigRevolution = strings.TrimSpace(m[1]) == "1"
	}
	return state
}

func buildAxisState(xml string, tags []string) common.ChartAxisState {
	states := buildAxisStates(xml, tags)
	if len(states) == 0 {
		return common.ChartAxisState{}
	}
	return states[0]
}

func buildAxisStates(xml string, tags []string) []common.ChartAxisState {
	states := make([]common.ChartAxisState, 0)
	for _, tag := range tags {
		startTag := "<c:" + tag + ">"
		endTag := chartElementClosePrefix + tag + ">"
		remaining := xml
		for {
			start := strings.Index(remaining, startTag)
			if start < 0 {
				break
			}
			endRel := strings.Index(remaining[start:], endTag)
			if endRel < 0 {
				break
			}
			end := start + endRel + len(endTag)
			states = append(states, parseAxisStateBlock(remaining[start:end]))
			remaining = remaining[end:]
		}
	}
	return states
}

func parseAxisStateBlock(block string) common.ChartAxisState {
	state := common.ChartAxisState{Present: true}
	if match := reTickLblPos.FindStringSubmatch(block); len(match) == expectedSingleGroupMatch {
		state.TickLabelPos = strings.TrimSpace(match[1])
	}
	if match := reCrosses.FindStringSubmatch(block); len(match) == expectedSingleGroupMatch {
		state.Crosses = strings.TrimSpace(match[1])
	}
	if match := reStateAxisTitle.FindStringSubmatch(block); len(match) == expectedSingleGroupMatch {
		state.Title = strings.TrimSpace(html.UnescapeString(match[1]))
	}
	state.MinimumScale = axisFloatValue(block, reStateAxisMin)
	state.MaximumScale = axisFloatValue(block, reStateAxisMax)
	state.MajorUnit = axisFloatValue(block, reStateAxisMajor)
	state.MinorUnit = axisFloatValue(block, reStateAxisMinor)
	if match := reStateAxisNumFmt.FindStringSubmatch(block); len(match) == expectedAxisNumFmtMatch {
		state.NumberFormat = match[1]
		linked := strings.TrimSpace(match[2]) == "1"
		state.FormatLinked = &linked
	}
	if match := reStateTickSkip.FindStringSubmatch(block); len(match) == expectedSingleGroupMatch {
		if skip, err := strconv.Atoi(strings.TrimSpace(match[1])); err == nil {
			state.TickMarkSkip = &skip
		}
	}
	if match := reStateLblAlgn.FindStringSubmatch(block); len(match) == expectedSingleGroupMatch {
		state.LabelAlignment = strings.TrimSpace(match[1])
	}
	state.MajorGridline = strings.Contains(block, "<c:majorGridlines")
	state.MinorGridline = strings.Contains(block, "<c:minorGridlines")
	return state
}

func axisFloatValue(block string, re *regexp.Regexp) *float64 {
	match := re.FindStringSubmatch(block)
	if len(match) != expectedSingleGroupMatch {
		return nil
	}
	value, err := strconv.ParseFloat(strings.TrimSpace(match[1]), 64)
	if err != nil {
		return nil
	}
	return &value
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
