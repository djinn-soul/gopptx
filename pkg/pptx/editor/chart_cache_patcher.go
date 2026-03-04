package editor

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editormodchart "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/chart"
)

var (
	chartSeriesPattern = regexp.MustCompile(`(?s)<c:ser>.*?</c:ser>`)
	xmlFormulaPattern  = regexp.MustCompile(`(?s)<c:f>.*?</c:f>`)
	strCachePattern    = regexp.MustCompile(`(?s)<c:strCache>.*?</c:strCache>`)
	numCachePattern    = regexp.MustCompile(`(?s)<c:numCache>.*?</c:numCache>`)
	strLitPattern      = regexp.MustCompile(`(?s)<c:strLit>.*?</c:strLit>`)
	numLitPattern      = regexp.MustCompile(`(?s)<c:numLit>.*?</c:numLit>`)
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

	valueCol := editormodchart.ColumnName(seriesIdx + firstSeriesValueColumnOffset)
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
	xCol := editormodchart.ColumnName(baseCol)
	yCol := editormodchart.ColumnName(baseCol + 1)

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

	sizeCol := editormodchart.ColumnName(baseCol + bubbleSizeColumnOffset)
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

	updatedField := applyFieldFormula(field, formula)
	var err error
	switch {
	case len(strVals) > 0:
		updatedField, err = applyStringValues(fieldTag, updatedField, strVals)
	case len(numVals) > 0:
		updatedField, err = applyNumericValues(fieldTag, updatedField, numVals)
	default:
		err = fmt.Errorf("no values provided for %s", fieldTag)
	}
	if err != nil {
		return "", err
	}

	return strings.Replace(seriesXML, field, updatedField, 1), nil
}

func applyFieldFormula(field string, formula string) string {
	formulaNode := xmlFormulaPattern.FindString(field)
	if formulaNode == "" {
		return field
	}
	return strings.Replace(field, formulaNode, "<c:f>"+common.XMLEscape(formula)+"</c:f>", 1)
}

func applyStringValues(fieldTag string, field string, strVals []string) (string, error) {
	switch {
	case strCachePattern.MatchString(field):
		return strCachePattern.ReplaceAllString(field, buildStringData("strCache", strVals)), nil
	case strLitPattern.MatchString(field):
		return strLitPattern.ReplaceAllString(field, buildStringData("strLit", strVals)), nil
	case numCachePattern.MatchString(field):
		numeric, err := convertStringsToFloats(strVals, "numeric cache")
		if err != nil {
			return "", err
		}
		return numCachePattern.ReplaceAllString(field, buildNumberData("numCache", numeric)), nil
	case numLitPattern.MatchString(field):
		numeric, err := convertStringsToFloats(strVals, "numeric literal")
		if err != nil {
			return "", err
		}
		return numLitPattern.ReplaceAllString(field, buildNumberData("numLit", numeric)), nil
	default:
		return "", fmt.Errorf("missing data node for %s", fieldTag)
	}
}

func applyNumericValues(fieldTag string, field string, numVals []float64) (string, error) {
	switch {
	case numCachePattern.MatchString(field):
		return numCachePattern.ReplaceAllString(field, buildNumberData("numCache", numVals)), nil
	case numLitPattern.MatchString(field):
		return numLitPattern.ReplaceAllString(field, buildNumberData("numLit", numVals)), nil
	default:
		return "", fmt.Errorf("missing numeric data node for %s", fieldTag)
	}
}

func convertStringsToFloats(strVals []string, dest string) ([]float64, error) {
	values := make([]float64, 0, len(strVals))
	for _, s := range strVals {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, fmt.Errorf("category %q cannot be represented in %s", s, dest)
		}
		values = append(values, f)
	}
	return values, nil
}

func buildStringData(tag string, vals []string) string {
	var b strings.Builder
	b.WriteString("<c:" + tag + ">")
	b.WriteString(fmt.Sprintf("<c:ptCount val=\"%d\"/>", len(vals)))
	for i, v := range vals {
		b.WriteString(fmt.Sprintf("<c:pt idx=\"%d\"><c:v>%s</c:v></c:pt>", i, common.XMLEscape(v)))
	}
	b.WriteString("</c:" + tag + ">")
	return b.String()
}

func buildNumberData(tag string, vals []float64) string {
	var b strings.Builder
	b.WriteString("<c:" + tag + ">")
	b.WriteString("<c:formatCode>General</c:formatCode>")
	b.WriteString(fmt.Sprintf("<c:ptCount val=\"%d\"/>", len(vals)))
	for i, v := range vals {
		b.WriteString(fmt.Sprintf("<c:pt idx=\"%d\"><c:v>%s</c:v></c:pt>", i, strconv.FormatFloat(v, 'f', -1, 64)))
	}
	b.WriteString("</c:" + tag + ">")
	return b.String()
}

func sheetRange(col string, n int) string {
	return fmt.Sprintf("Sheet1!$%s$2:$%s$%d", col, col, n+1)
}
