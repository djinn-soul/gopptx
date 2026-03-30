package export

import (
	"archive/zip"
	"regexp"
	"strconv"
	"strings"

	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

type parsedChartSeries struct {
	Name       string
	Categories []string
	Values     []float64
	XValues    []float64
	YValues    []float64
	Sizes      []float64
}

type parsedChart struct {
	ShapeID      int
	Kind         string
	Title        string
	X            int64
	Y            int64
	CX           int64
	CY           int64
	AltText      string
	IsDecorative bool
	Series       []parsedChartSeries
}

type chartFrameRef struct {
	ShapeID      int
	RelID        string
	X            int64
	Y            int64
	CX           int64
	CY           int64
	AltText      string
	IsDecorative bool
}

var (
	reGraphicFrame = regexp.MustCompile(`(?s)<p:graphicFrame\b.*?</p:graphicFrame>`)
	reChartRelID   = regexp.MustCompile(`\br:id="([^"]+)"`)
	reSeriesBlock  = regexp.MustCompile(`(?s)<c:ser\b.*?</c:ser>`)
	reTextValue    = regexp.MustCompile(`(?s)<a:t>(.*?)</a:t>|<c:v>(.*?)</c:v>`)
	rePointValue   = regexp.MustCompile(`(?s)<c:pt\b[^>]*>.*?<c:v>(.*?)</c:v>.*?</c:pt>`)
)

func extractSlideCharts(pptxPath string) ([][]parsedChart, error) {
	zr, err := zip.OpenReader(pptxPath)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	fileMap := make(map[string]*zip.File, len(zr.File))
	for _, f := range zr.File {
		fileMap[canonicalZipPath(f.Name)] = f
	}
	slideOrder := resolveSlideOrder(fileMap)
	out := make([][]parsedChart, len(slideOrder))
	for idx, slidePart := range slideOrder {
		slideXML := readZipBytes(fileMap, slidePart)
		if slideXML == nil {
			continue
		}
		frames := parseChartFrames(slideXML)
		if len(frames) == 0 {
			continue
		}
		rels := readZipRelationships(fileMap, slideRelsPath(slidePart))
		row := make([]parsedChart, 0, len(frames))
		for _, frame := range frames {
			target := rels[frame.RelID]
			if target == "" {
				continue
			}
			chartPath := resolveRelPath(slidePart, target)
			if chartPath == "" {
				continue
			}
			chartXML := readZipBytes(fileMap, chartPath)
			if chartXML == nil {
				continue
			}
			pc := parseChartPart(chartXML)
			pc.ShapeID = frame.ShapeID
			pc.X, pc.Y, pc.CX, pc.CY = frame.X, frame.Y, frame.CX, frame.CY
			pc.AltText = frame.AltText
			pc.IsDecorative = frame.IsDecorative
			row = append(row, pc)
		}
		out[idx] = row
	}
	return out, nil
}

func parseChartFrames(slideXML []byte) []chartFrameRef {
	xml := string(slideXML)
	frames := reGraphicFrame.FindAllString(xml, -1)
	out := make([]chartFrameRef, 0, len(frames))
	for _, frame := range frames {
		if !strings.Contains(frame, "<c:chart") && !strings.Contains(frame, "<cx:chart") {
			continue
		}
		props, err := editorshape.ParseShapeProperties([]byte(frame))
		if err != nil {
			continue
		}
		meta, err := editorshape.ParseShapeReaderMetadata([]byte(frame))
		if err != nil {
			continue
		}
		idMatch := reChartRelID.FindStringSubmatch(frame)
		if len(idMatch) < 2 {
			continue
		}
		out = append(out, chartFrameRef{
			ShapeID:      props.ID,
			RelID:        idMatch[1],
			X:            int64(props.X),
			Y:            int64(props.Y),
			CX:           int64(props.W),
			CY:           int64(props.H),
			AltText:      meta.AltText,
			IsDecorative: meta.IsDecorative,
		})
	}
	return out
}

func parseChartPart(chartXML []byte) parsedChart {
	xml := string(chartXML)
	result := parsedChart{
		Kind:  detectChartKind(xml),
		Title: firstText(xml, "<c:title", "</c:title>"),
	}
	blocks := reSeriesBlock.FindAllString(xml, -1)
	result.Series = make([]parsedChartSeries, 0, len(blocks))
	for _, block := range blocks {
		result.Series = append(result.Series, parseSeriesBlock(block))
	}
	return result
}

func parseSeriesBlock(block string) parsedChartSeries {
	series := parsedChartSeries{Name: firstText(block, "<c:tx", "</c:tx>")}
	if full := firstTagBlock(block, "c:cat"); full != "" {
		series.Categories = extractTextPoints(full)
	}
	if full := firstTagBlock(block, "c:val"); full != "" {
		series.Values = extractFloatPoints(full)
	}
	if full := firstTagBlock(block, "c:xVal"); full != "" {
		series.XValues = extractFloatPoints(full)
	}
	if full := firstTagBlock(block, "c:yVal"); full != "" {
		series.YValues = extractFloatPoints(full)
	}
	if full := firstTagBlock(block, "c:bubbleSize"); full != "" {
		series.Sizes = extractFloatPoints(full)
	}
	return series
}

func firstTagBlock(xml, tag string) string {
	start := strings.Index(xml, "<"+tag)
	if start < 0 {
		return ""
	}
	endTag := "</" + tag + ">"
	endRel := strings.Index(xml[start:], endTag)
	if endRel < 0 {
		return ""
	}
	return xml[start : start+endRel+len(endTag)]
}

func detectChartKind(xml string) string {
	switch {
	case strings.Contains(xml, "<c:stockChart"):
		if strings.Count(xml, "<c:ser") >= 4 {
			return "stockOHLC"
		}
		return "stockHLC"
	case strings.Contains(xml, "<c:barChart") && strings.Contains(xml, "<c:lineChart"):
		return "combo"
	case strings.Contains(xml, "<c:doughnutChart"):
		return "doughnut"
	case strings.Contains(xml, "<c:pieChart"):
		return "pie"
	case strings.Contains(xml, "<c:bubbleChart"):
		return "bubble"
	case strings.Contains(xml, "<c:scatterChart"):
		return "scatter"
	case strings.Contains(xml, "<c:radarChart"):
		if strings.Contains(xml, `radarStyle val="filled"`) {
			return "radarFilled"
		}
		return "radar"
	case strings.Contains(xml, "<c:areaChart"):
		if strings.Contains(xml, `grouping val="percentStacked"`) {
			return "areaStacked100"
		}
		if strings.Contains(xml, `grouping val="stacked"`) {
			return "areaStacked"
		}
		return "area"
	case strings.Contains(xml, "<c:lineChart"):
		if strings.Contains(xml, `grouping val="stacked"`) {
			return "lineStacked"
		}
		if strings.Contains(xml, "<c:marker") {
			return "lineMarkers"
		}
		return "line"
	case strings.Contains(xml, "<c:barChart"):
		if strings.Contains(xml, `grouping val="percentStacked"`) {
			return "barStacked100"
		}
		if strings.Contains(xml, `grouping val="stacked"`) {
			return "barStacked"
		}
		if strings.Contains(xml, `barDir val="bar"`) {
			return "barHorizontal"
		}
		return "bar"
	default:
		return "bar"
	}
}

func extractTextPoints(block string) []string {
	matches := rePointValue.FindAllStringSubmatch(block, -1)
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		out = append(out, htmlEntityDecode(strings.TrimSpace(m[1])))
	}
	return out
}

func extractFloatPoints(block string) []float64 {
	texts := extractTextPoints(block)
	out := make([]float64, 0, len(texts))
	for _, txt := range texts {
		if txt == "" {
			continue
		}
		n, err := strconv.ParseFloat(strings.TrimSpace(txt), 64)
		if err == nil {
			out = append(out, n)
		}
	}
	return out
}

func firstText(xml, startTag, endTag string) string {
	start := strings.Index(xml, startTag)
	if start < 0 {
		return ""
	}
	endRel := strings.Index(xml[start:], endTag)
	if endRel < 0 {
		return ""
	}
	segment := xml[start : start+endRel+len(endTag)]
	match := reTextValue.FindStringSubmatch(segment)
	if len(match) < 3 {
		return ""
	}
	if match[1] != "" {
		return htmlEntityDecode(strings.TrimSpace(match[1]))
	}
	return htmlEntityDecode(strings.TrimSpace(match[2]))
}

func htmlEntityDecode(value string) string {
	replacer := strings.NewReplacer(
		"&amp;", "&",
		"&lt;", "<",
		"&gt;", ">",
		"&quot;", `"`,
		"&apos;", "'",
	)
	return replacer.Replace(value)
}
