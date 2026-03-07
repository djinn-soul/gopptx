package chart

import (
	"errors"
	"fmt"
	"regexp"

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
