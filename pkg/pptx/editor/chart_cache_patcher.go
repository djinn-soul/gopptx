package editor

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

var (
	chartSeriesPattern = regexp.MustCompile(`(?s)<c:ser>.*?</c:ser>`)
	xmlFormulaPattern  = regexp.MustCompile(`(?s)<c:f>.*?</c:f>`)
	strCachePattern    = regexp.MustCompile(`(?s)<c:strCache>.*?</c:strCache>`)
	numCachePattern    = regexp.MustCompile(`(?s)<c:numCache>.*?</c:numCache>`)
)

const (
	firstSeriesValueColumnOffset = 2
	scatterColumnsPerSeries      = 2
	bubbleColumnsPerSeries       = 3
	bubbleSizeColumnOffset       = 2
)

func patchChartDataCache(chartXML []byte, kind chartKind, req common.ChartDataUpdate) ([]byte, error) {
	src := string(chartXML)
	series := chartSeriesPattern.FindAllString(src, -1)
	if len(series) == 0 {
		return nil, errors.New("chart has no series nodes")
	}
	if len(series) != len(req.Series) {
		return nil, fmt.Errorf("series count mismatch: chart has %d, payload has %d", len(series), len(req.Series))
	}

	for i := range series {
		var (
			updated string
			err     error
		)
		switch kind {
		case chartKindScatter:
			updated, err = patchScatterSeries(i, series[i], req.Series[i], false)
		case chartKindBubble:
			updated, err = patchScatterSeries(i, series[i], req.Series[i], true)
		default:
			updated, err = patchCategorySeries(i, series[i], req.Categories, req.Series[i])
		}
		if err != nil {
			return nil, err
		}
		series[i] = updated
	}

	result := chartSeriesPattern.ReplaceAllStringFunc(src, func(_ string) string {
		if len(series) == 0 {
			return ""
		}
		out := series[0]
		series = series[1:]
		return out
	})
	return []byte(result), nil
}

func patchCategorySeries(
	seriesIdx int,
	seriesXML string,
	categories []string,
	data common.ChartSeriesData,
) (string, error) {
	cats := categories
	if len(data.Categories) > 0 {
		cats = data.Categories
	}
	if len(cats) != len(data.Values) {
		return "", fmt.Errorf("series %d category/value length mismatch", seriesIdx)
	}

	out, err := replaceFieldContent(seriesXML, "cat", sheetRange("A", len(cats)), cats, nil)
	if err != nil {
		return "", fmt.Errorf("series %d categories: %w", seriesIdx, err)
	}

	valueCol := columnName(seriesIdx + firstSeriesValueColumnOffset)
	out, err = replaceFieldContent(out, "val", sheetRange(valueCol, len(data.Values)), nil, data.Values)
	if err != nil {
		return "", fmt.Errorf("series %d values: %w", seriesIdx, err)
	}
	return out, nil
}

func patchScatterSeries(seriesIdx int, seriesXML string, data common.ChartSeriesData, bubble bool) (string, error) {
	baseCol := seriesIdx*scatterColumnsPerSeries + 1
	if bubble {
		baseCol = seriesIdx*bubbleColumnsPerSeries + 1
	}
	xCol := columnName(baseCol)
	yCol := columnName(baseCol + 1)

	out, err := replaceFieldContent(seriesXML, "xVal", sheetRange(xCol, len(data.XValues)), nil, data.XValues)
	if err != nil {
		return "", fmt.Errorf("series %d x values: %w", seriesIdx, err)
	}
	out, err = replaceFieldContent(out, "yVal", sheetRange(yCol, len(data.YValues)), nil, data.YValues)
	if err != nil {
		return "", fmt.Errorf("series %d y values: %w", seriesIdx, err)
	}
	if !bubble {
		return out, nil
	}

	sizeCol := columnName(baseCol + bubbleSizeColumnOffset)
	out, err = replaceFieldContent(out, "bubbleSize", sheetRange(sizeCol, len(data.Sizes)), nil, data.Sizes)
	if err != nil {
		return "", fmt.Errorf("series %d bubble sizes: %w", seriesIdx, err)
	}
	return out, nil
}

func replaceFieldContent(seriesXML, fieldTag, formula string, strVals []string, numVals []float64) (string, error) {
	fieldPattern := regexp.MustCompile(`(?s)<c:` + fieldTag + `>.*?</c:` + fieldTag + `>`)
	field := fieldPattern.FindString(seriesXML)
	if field == "" {
		return "", fmt.Errorf("missing %s node", fieldTag)
	}

	formulaNode := xmlFormulaPattern.FindString(field)
	if formulaNode == "" {
		return "", fmt.Errorf("missing formula node for %s", fieldTag)
	}
	updatedField := strings.Replace(field, formulaNode, "<c:f>"+common.XMLEscape(formula)+"</c:f>", 1)

	switch {
	case len(strVals) > 0:
		switch {
		case strCachePattern.MatchString(updatedField):
			updatedField = strCachePattern.ReplaceAllString(updatedField, buildStringCache(strVals))
		case numCachePattern.MatchString(updatedField):
			numValsFromCats := make([]float64, 0, len(strVals))
			for _, s := range strVals {
				f, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return "", fmt.Errorf("category %q cannot be represented in numeric cache", s)
				}
				numValsFromCats = append(numValsFromCats, f)
			}
			updatedField = numCachePattern.ReplaceAllString(updatedField, buildNumberCache(numValsFromCats))
		default:
			return "", fmt.Errorf("missing cache node for %s", fieldTag)
		}
	case len(numVals) > 0:
		if !numCachePattern.MatchString(updatedField) {
			return "", fmt.Errorf("missing numeric cache node for %s", fieldTag)
		}
		updatedField = numCachePattern.ReplaceAllString(updatedField, buildNumberCache(numVals))
	default:
		return "", fmt.Errorf("no values provided for %s", fieldTag)
	}

	return strings.Replace(seriesXML, field, updatedField, 1), nil
}

func buildStringCache(vals []string) string {
	var b strings.Builder
	b.WriteString("<c:strCache>")
	b.WriteString(fmt.Sprintf("<c:ptCount val=\"%d\"/>", len(vals)))
	for i, v := range vals {
		b.WriteString(fmt.Sprintf("<c:pt idx=\"%d\"><c:v>%s</c:v></c:pt>", i, common.XMLEscape(v)))
	}
	b.WriteString("</c:strCache>")
	return b.String()
}

func buildNumberCache(vals []float64) string {
	var b strings.Builder
	b.WriteString("<c:numCache>")
	b.WriteString(fmt.Sprintf("<c:ptCount val=\"%d\"/>", len(vals)))
	for i, v := range vals {
		b.WriteString(fmt.Sprintf("<c:pt idx=\"%d\"><c:v>%s</c:v></c:pt>", i, strconv.FormatFloat(v, 'f', -1, 64)))
	}
	b.WriteString("</c:numCache>")
	return b.String()
}

func sheetRange(col string, n int) string {
	return fmt.Sprintf("Sheet1!$%s$2:$%s$%d", col, col, n+1)
}
