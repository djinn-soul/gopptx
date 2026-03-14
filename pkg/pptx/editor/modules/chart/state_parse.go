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
	reScene3D    = regexp.MustCompile(`(?s)<a:scene3d\b.*?</a:scene3d>`)
	reCamera     = regexp.MustCompile(`<a:camera\b[^>]*prst="([^"]+)"[^>]*>`)
	reCameraFOV  = regexp.MustCompile(`<a:camera\b[^>]*fov="([^"]+)"[^>]*>`)
	reLightRig   = regexp.MustCompile(`<a:lightRig\b[^>]*rig="([^"]+)"[^>]*>`)
	reLightDir   = regexp.MustCompile(`<a:lightRig\b[^>]*dir="([^"]+)"[^>]*>`)
	reLightRev   = regexp.MustCompile(`<a:lightRig\b[^>]*rev="([^"]+)"[^>]*>`)
)

const expectedSingleGroupMatch = 2

// ExtractChartState parses the chart XML cache into a traversal-friendly snapshot.
func ExtractChartState(chartXML []byte) common.ChartState {
	xml := string(chartXML)
	state := common.ChartState{
		CategoryAx: buildAxisState(xml, []string{"catAx", "dateAx"}),
		ValueAx:    buildAxisState(xml, []string{"valAx"}),
		Series:     parseSeriesState(xml),
		Scene3D:    parseScene3DState(xml),
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
		state.MinorGridline = strings.Contains(block, "<c:minorGridlines")
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
