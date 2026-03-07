package chart

import (
	"fmt"
	"strconv"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const excelColumnBase = 26

func categoryHeaders(req common.ChartDataUpdate) []string {
	if len(req.MultiLevelCategories) == 0 {
		return []string{"Category"}
	}
	headers := make([]string, 0, len(req.MultiLevelCategories))
	for i := range req.MultiLevelCategories {
		headers = append(headers, fmt.Sprintf("Category Level %d", i+1))
	}
	return headers
}

func buildCategoryRows(
	categories []string,
	multiLevel [][]string,
	series []common.ChartSeriesData,
) [][]string {
	rowCount := categoryRowCount(categories, multiLevel, series)
	rows := make([][]string, rowCount)
	for i := range rowCount {
		rows[i] = buildCategoryRow(i, categories, multiLevel, series)
	}
	return rows
}

func buildCategoryRow(
	rowIndex int,
	categories []string,
	multiLevel [][]string,
	series []common.ChartSeriesData,
) []string {
	row := make([]string, 0, categoryColumnCount(multiLevel)+len(series))
	row = append(row, rowCategoryValues(rowIndex, categories, multiLevel, series)...)
	for _, s := range series {
		row = append(row, formatSeriesValueAt(s.Values, rowIndex))
	}
	return row
}

func rowCategoryValues(
	rowIndex int,
	categories []string,
	multiLevel [][]string,
	series []common.ChartSeriesData,
) []string {
	if len(multiLevel) > 0 {
		return multiLevelCategoryValues(rowIndex, multiLevel)
	}
	return []string{singleLevelCategoryValue(rowIndex, categories, series)}
}

func multiLevelCategoryValues(rowIndex int, multiLevel [][]string) []string {
	values := make([]string, 0, len(multiLevel))
	for lvl := range multiLevel {
		cat := ""
		if len(multiLevel[lvl]) > rowIndex {
			cat = multiLevel[lvl][rowIndex]
		}
		values = append(values, cat)
	}
	return values
}

func singleLevelCategoryValue(rowIndex int, categories []string, series []common.ChartSeriesData) string {
	if len(categories) > rowIndex {
		return categories[rowIndex]
	}
	if len(series) > 0 && len(series[0].Categories) > rowIndex {
		return series[0].Categories[rowIndex]
	}
	return ""
}

func formatSeriesValueAt(values []float64, rowIndex int) string {
	if len(values) <= rowIndex {
		return ""
	}
	return strconv.FormatFloat(values[rowIndex], 'f', -1, 64)
}

func categoryRowCount(categories []string, multiLevel [][]string, series []common.ChartSeriesData) int {
	if len(multiLevel) > 0 {
		return len(multiLevel[0])
	}
	rowCount := len(categories)
	if rowCount == 0 && len(series) > 0 && len(series[0].Categories) > 0 {
		rowCount = len(series[0].Categories)
	}
	return rowCount
}

func categoryColumnCount(multiLevel [][]string) int {
	if len(multiLevel) == 0 {
		return 1
	}
	return len(multiLevel)
}

func buildScatterSheet(series []common.ChartSeriesData, withSizes bool) ([]string, [][]string) {
	headers, maxRows := scatterHeadersAndMaxRows(series, withSizes)
	rows := make([][]string, maxRows)
	for rowIdx := range maxRows {
		rows[rowIdx] = scatterRowValues(series, rowIdx, withSizes, len(headers))
	}
	return headers, rows
}

func scatterHeadersAndMaxRows(series []common.ChartSeriesData, withSizes bool) ([]string, int) {
	headers := make([]string, 0, len(series)*scatterHeaderColFactor)
	maxRows := 0
	for i, s := range series {
		headers = append(headers, fmt.Sprintf("X%d", i+1), fmt.Sprintf("Y%d", i+1))
		if withSizes {
			headers = append(headers, fmt.Sprintf("S%d", i+1))
		}
		if len(s.XValues) > maxRows {
			maxRows = len(s.XValues)
		}
	}
	return headers, maxRows
}

func scatterRowValues(series []common.ChartSeriesData, rowIdx int, withSizes bool, capacity int) []string {
	row := make([]string, 0, capacity)
	for _, s := range series {
		row = append(row, scatterValueAt(s.XValues, rowIdx))
		row = append(row, scatterValueAt(s.YValues, rowIdx))
		if withSizes {
			row = append(row, scatterValueAt(s.Sizes, rowIdx))
		}
	}
	return row
}

func scatterValueAt(values []float64, idx int) string {
	if len(values) <= idx {
		return ""
	}
	return strconv.FormatFloat(values[idx], 'f', -1, 64)
}

func ColumnName(n int) string {
	name := ""
	for n > 0 {
		n--
		name = string(rune('A'+(n%excelColumnBase))) + name
		n /= excelColumnBase
	}
	return name
}

func isNumberLiteral(s string) bool {
	if s == "" {
		return false
	}
	if _, err := strconv.ParseFloat(s, 64); err != nil {
		return false
	}
	return true
}
