package chart

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
	strLitPattern      = regexp.MustCompile(`(?s)<c:strLit>.*?</c:strLit>`)
	numLitPattern      = regexp.MustCompile(`(?s)<c:numLit>.*?</c:numLit>`)
	multiLvlCache      = regexp.MustCompile(`(?s)<c:multiLvlStrCache>.*?</c:multiLvlStrCache>`)
	multiLvlLit        = regexp.MustCompile(`(?s)<c:multiLvlStrLit>.*?</c:multiLvlStrLit>`)
	formatCodePattern  = regexp.MustCompile(`(?s)<c:formatCode>(.*?)</c:formatCode>`)
)

const (
	firstSeriesValueColumnOffset = 2
	scatterColumnsPerSeries      = 2
	bubbleColumnsPerSeries       = 3
	bubbleSizeColumnOffset       = 2
)

func PatchChartDataCache(chartXML []byte, kind Kind, req common.ChartDataUpdate) ([]byte, error) {
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
		case KindScatter:
			updated, err = patchScatterSeries(i, series[i], req.Series[i], false)
		case KindBubble:
			updated, err = patchScatterSeries(i, series[i], req.Series[i], true)
		default:
			updated, err = patchCategorySeries(i, series[i], req.Categories, req.MultiLevelCategories, req.Series[i])
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
	multiLevelCategories [][]string,
	data common.ChartSeriesData,
) (string, error) {
	if len(multiLevelCategories) > 0 {
		return patchMultiLevelCategorySeries(seriesIdx, seriesXML, multiLevelCategories, data)
	}

	cats := categories
	if len(data.Categories) > 0 {
		cats = data.Categories
	}
	if len(cats) != len(data.Values) {
		return "", fmt.Errorf("series %d category/value length mismatch", seriesIdx)
	}

	out, err := replaceFieldContent(seriesXML, "cat", sheetRange("A", len(cats)), cats, nil, nil)
	if err != nil {
		return "", fmt.Errorf("series %d categories: %w", seriesIdx, err)
	}

	valueCol := ColumnName(seriesIdx + firstSeriesValueColumnOffset)
	out, err = replaceFieldContent(out, "val", sheetRange(valueCol, len(data.Values)), nil, data.Values, nil)
	if err != nil {
		return "", fmt.Errorf("series %d values: %w", seriesIdx, err)
	}
	return out, nil
}

func patchMultiLevelCategorySeries(
	seriesIdx int,
	seriesXML string,
	multiLevelCategories [][]string,
	data common.ChartSeriesData,
) (string, error) {
	if len(multiLevelCategories) == 0 {
		return "", fmt.Errorf("series %d requires multi-level categories", seriesIdx)
	}
	leafCount := len(multiLevelCategories[0])
	if leafCount == 0 {
		return "", fmt.Errorf("series %d multi-level categories are empty", seriesIdx)
	}
	for lvl := 1; lvl < len(multiLevelCategories); lvl++ {
		if len(multiLevelCategories[lvl]) != leafCount {
			return "", fmt.Errorf("series %d multi-level categories have inconsistent lengths", seriesIdx)
		}
	}
	if len(data.Values) != leafCount {
		return "", fmt.Errorf("series %d category/value length mismatch", seriesIdx)
	}

	out, err := replaceFieldContent(
		seriesXML,
		"cat",
		sheetRangeAcrossColumns(1, len(multiLevelCategories), leafCount),
		nil,
		nil,
		multiLevelCategories,
	)
	if err != nil {
		return "", fmt.Errorf("series %d categories: %w", seriesIdx, err)
	}

	valueCol := ColumnName(seriesIdx + len(multiLevelCategories) + 1)
	out, err = replaceFieldContent(out, "val", sheetRange(valueCol, len(data.Values)), nil, data.Values, nil)
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
	xCol := ColumnName(baseCol)
	yCol := ColumnName(baseCol + 1)

	out, err := replaceFieldContent(seriesXML, "xVal", sheetRange(xCol, len(data.XValues)), nil, data.XValues, nil)
	if err != nil {
		return "", fmt.Errorf("series %d x values: %w", seriesIdx, err)
	}
	out, err = replaceFieldContent(out, "yVal", sheetRange(yCol, len(data.YValues)), nil, data.YValues, nil)
	if err != nil {
		return "", fmt.Errorf("series %d y values: %w", seriesIdx, err)
	}
	if !bubble {
		return out, nil
	}

	sizeCol := ColumnName(baseCol + bubbleSizeColumnOffset)
	out, err = replaceFieldContent(out, "bubbleSize", sheetRange(sizeCol, len(data.Sizes)), nil, data.Sizes, nil)
	if err != nil {
		return "", fmt.Errorf("series %d bubble sizes: %w", seriesIdx, err)
	}
	return out, nil
}

func replaceFieldContent(
	seriesXML string,
	fieldTag string,
	formula string,
	strVals []string,
	numVals []float64,
	multiLevelVals [][]string,
) (string, error) {
	fieldPattern := regexp.MustCompile(`(?s)<c:` + fieldTag + `>.*?</c:` + fieldTag + `>`)
	field := fieldPattern.FindString(seriesXML)
	if field == "" {
		return "", fmt.Errorf("missing %s node", fieldTag)
	}

	updatedField := applyFieldFormula(field, formula)
	var err error
	switch {
	case len(multiLevelVals) > 0:
		updatedField, err = applyMultiLevelValues(fieldTag, updatedField, multiLevelVals)
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
		existing := numCachePattern.FindString(field)
		formatCode := extractFormatCode(existing)
		return numCachePattern.ReplaceAllString(field, buildNumberData("numCache", formatCode, numeric)), nil
	case numLitPattern.MatchString(field):
		numeric, err := convertStringsToFloats(strVals, "numeric literal")
		if err != nil {
			return "", err
		}
		existing := numLitPattern.FindString(field)
		formatCode := extractFormatCode(existing)
		return numLitPattern.ReplaceAllString(field, buildNumberData("numLit", formatCode, numeric)), nil
	default:
		return "", fmt.Errorf("missing data node for %s", fieldTag)
	}
}

func applyMultiLevelValues(fieldTag string, field string, vals [][]string) (string, error) {
	switch {
	case multiLvlCache.MatchString(field):
		return multiLvlCache.ReplaceAllString(field, buildMultiLevelData("multiLvlStrCache", vals)), nil
	case multiLvlLit.MatchString(field):
		return multiLvlLit.ReplaceAllString(field, buildMultiLevelData("multiLvlStrLit", vals)), nil
	default:
		return "", fmt.Errorf("missing multi-level category node for %s", fieldTag)
	}
}

func applyNumericValues(fieldTag string, field string, numVals []float64) (string, error) {
	switch {
	case numCachePattern.MatchString(field):
		existing := numCachePattern.FindString(field)
		formatCode := extractFormatCode(existing)
		return numCachePattern.ReplaceAllString(field, buildNumberData("numCache", formatCode, numVals)), nil
	case numLitPattern.MatchString(field):
		existing := numLitPattern.FindString(field)
		formatCode := extractFormatCode(existing)
		return numLitPattern.ReplaceAllString(field, buildNumberData("numLit", formatCode, numVals)), nil
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

func buildNumberData(tag string, formatCode string, vals []float64) string {
	var b strings.Builder
	b.WriteString("<c:" + tag + ">")
	b.WriteString("<c:formatCode>" + common.XMLEscape(formatCode) + "</c:formatCode>")
	b.WriteString(fmt.Sprintf("<c:ptCount val=\"%d\"/>", len(vals)))
	for i, v := range vals {
		b.WriteString(fmt.Sprintf("<c:pt idx=\"%d\"><c:v>%s</c:v></c:pt>", i, strconv.FormatFloat(v, 'f', -1, 64)))
	}
	b.WriteString("</c:" + tag + ">")
	return b.String()
}

func buildMultiLevelData(tag string, levels [][]string) string {
	leafCount := len(levels[0])
	var b strings.Builder
	b.WriteString("<c:" + tag + ">")
	b.WriteString(fmt.Sprintf("<c:ptCount val=\"%d\"/>", leafCount))
	for _, lvl := range levels {
		b.WriteString("<c:lvl>")
		for i, v := range lvl {
			b.WriteString(fmt.Sprintf("<c:pt idx=\"%d\"><c:v>%s</c:v></c:pt>", i, common.XMLEscape(v)))
		}
		b.WriteString("</c:lvl>")
	}
	b.WriteString("</c:" + tag + ">")
	return b.String()
}

func extractFormatCode(node string) string {
	matches := formatCodePattern.FindStringSubmatch(node)
	if len(matches) > 1 && strings.TrimSpace(matches[1]) != "" {
		return matches[1]
	}
	return "General"
}

func sheetRange(col string, n int) string {
	return fmt.Sprintf("Sheet1!$%s$2:$%s$%d", col, col, n+1)
}

func sheetRangeAcrossColumns(startCol int, endCol int, n int) string {
	start := ColumnName(startCol)
	end := ColumnName(endCol)
	return fmt.Sprintf("Sheet1!$%s$2:$%s$%d", start, end, n+1)
}
