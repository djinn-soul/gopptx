package export

import (
	"archive/zip"
	"html"
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
	Color      string // hex RGB from <a:srgbClr> inside the series spPr
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
	reGraphicFrame  = regexp.MustCompile(`(?s)<p:graphicFrame\b.*?</p:graphicFrame>`)
	reChartRelID    = regexp.MustCompile(`\br:id="([^"]+)"`)
	reSeriesBlock   = regexp.MustCompile(`(?s)<c:ser\b.*?</c:ser>`)
	reTextValue     = regexp.MustCompile(`(?s)<a:t>(.*?)</a:t>|<c:v>(.*?)</c:v>`)
	rePointValue    = regexp.MustCompile(`(?s)<c:pt\b[^>]*>.*?<c:v>(.*?)</c:v>.*?</c:pt>`)
	reSeriesSrgbClr = regexp.MustCompile(`<a:srgbClr val="([0-9A-Fa-f]{6})"`)
)

const (
	minChartRelIDMatchGroups = 2
	minPointMatchGroups      = 2
	minTextMatchGroups       = 3
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
	raw := string(slideXML)
	frames := reGraphicFrame.FindAllString(raw, -1)
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
		if len(idMatch) < minChartRelIDMatchGroups {
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
	raw := string(chartXML)
	result := parsedChart{
		Kind:  detectChartKind(raw),
		Title: firstText(raw, "<c:title", "</c:title>"),
	}
	blocks := reSeriesBlock.FindAllString(raw, -1)
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
	// Extract fill color from series shape properties.
	if m := reSeriesSrgbClr.FindStringSubmatch(block); len(m) >= 2 { //nolint:mnd
		series.Color = m[1]
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

func extractTextPoints(block string) []string {
	matches := rePointValue.FindAllStringSubmatch(block, -1)
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) < minPointMatchGroups {
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
	if len(match) < minTextMatchGroups {
		return ""
	}
	if match[1] != "" {
		return htmlEntityDecode(strings.TrimSpace(match[1]))
	}
	return htmlEntityDecode(strings.TrimSpace(match[2]))
}

func htmlEntityDecode(value string) string {
	return html.UnescapeString(value)
}
