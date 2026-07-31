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

// CachePatchOptions tunes what PatchChartDataCache rewrites besides the caches
// themselves. The zero value repoints every c:f at the columns the new data
// occupies, which is what a full data replacement wants.
type CachePatchOptions struct {
	// KeepFormulas leaves every c:f exactly as authored, so a refresh of the
	// cached numbers does not disturb a workbook link the user set up.
	KeepFormulas bool
}

// PatchChartDataCache refreshes the cached values a chart draws from.
func PatchChartDataCache(
	chartXML []byte,
	kind Kind,
	req common.ChartDataUpdate,
	opts CachePatchOptions,
) ([]byte, error) {
	patched, err := patchChartDataCacheXML(chartXML, kind, req)
	if err != nil {
		return nil, err
	}
	if !opts.KeepFormulas {
		return patched, nil
	}
	return restoreChartFormulas(chartXML, patched), nil
}

func patchChartDataCacheXML(
	chartXML []byte,
	kind Kind,
	req common.ChartDataUpdate,
) ([]byte, error) {
	src := string(chartXML)
	series := chartSeriesPattern.FindAllString(src, -1)
	if len(series) == 0 {
		return nil, errors.New("chart has no series nodes")
	}

	plotted := PlottedSeries(req.Series)
	if len(series) != len(plotted) {
		return nil, fmt.Errorf(
			"series count mismatch: chart has %d, payload has %d plotted (%d total)",
			len(series), len(plotted), len(req.Series),
		)
	}

	for i := range series {
		var (
			updated string
			err     error
		)
		// The workbook keeps every series, hidden ones included, so the sheet
		// column comes from the payload index rather than the plot order.
		col := plotted[i].WorkbookIndex
		switch kind {
		case KindScatter:
			updated, err = patchScatterSeries(col, series[i], plotted[i].Data, false)
		case KindBubble:
			updated, err = patchScatterSeries(col, series[i], plotted[i].Data, true)
		default:
			updated, err = patchCategorySeries(
				col, series[i], req.Categories, req.MultiLevelCategories, plotted[i].Data,
			)
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

// restoreChartFormulas writes every c:f of the source back over the patched
// XML, series by series. Pairing them by position is exact because the patch
// only ever rewrites the text of a c:f that already exists — applyFieldFormula
// never adds or drops one — so both sides hold the same nodes in the same
// order.
func restoreChartFormulas(chartXML, patched []byte) []byte {
	originalSeries := chartSeriesPattern.FindAllString(string(chartXML), -1)
	seriesIndex := 0
	result := chartSeriesPattern.ReplaceAllStringFunc(string(patched), func(updated string) string {
		if seriesIndex >= len(originalSeries) {
			return updated
		}
		restored := restoreSeriesFormulas(originalSeries[seriesIndex], updated)
		seriesIndex++
		return restored
	})
	return []byte(result)
}

func restoreSeriesFormulas(original, updated string) string {
	formulas := xmlFormulaPattern.FindAllString(original, -1)
	index := 0
	return xmlFormulaPattern.ReplaceAllStringFunc(updated, func(current string) string {
		if index >= len(formulas) {
			return current
		}
		formula := formulas[index]
		index++
		return formula
	})
}

// PlottedSeriesRef pairs a series that the chart draws with its column in the
// embedded workbook.
type PlottedSeriesRef struct {
	// WorkbookIndex is the series' position in the payload, which is also its
	// column in the embedded sheet.
	WorkbookIndex int
	Data          common.ChartSeriesData
}

// PlottedSeries returns the series the chart draws, dropping the hidden ones.
//
// A hidden series still gets its column in the embedded workbook — the data is
// there for the user to see behind the chart, it just is not plotted (upstream
// issue #1043).
func PlottedSeries(all []common.ChartSeriesData) []PlottedSeriesRef {
	out := make([]PlottedSeriesRef, 0, len(all))
	for i, data := range all {
		if data.Hidden {
			continue
		}
		out = append(out, PlottedSeriesRef{WorkbookIndex: i, Data: data})
	}
	return out
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
