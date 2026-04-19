package chart

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func PlaceholderChartDefinition(
	chartMap map[string]any,
	chartType string,
	title string,
	x, y, w, h int64,
) (charts.ChartDefinition, error) {
	normalizedType := strings.ToLower(chartType)
	switch normalizedType {
	case "scatter":
		return buildScatterChartDefinition(chartMap, title, x, y, w, h)
	case "bubble":
		return buildBubbleChartDefinition(chartMap, title, x, y, w, h)
	}
	categories, values, err := requireCategorySeries(chartMap)
	if err != nil {
		return nil, err
	}
	return buildCategoryChartDefinition(normalizedType, categories, values, title, x, y, w, h)
}

func buildScatterChartDefinition(
	chartMap map[string]any,
	title string,
	x, y, w, h int64,
) (charts.ChartDefinition, error) {
	xValues, yValues, err := requireXYSeries(chartMap)
	if err != nil {
		return nil, err
	}
	c := charts.NewScatterChart(xValues, yValues).WithTitle(title)
	return applyStyledBounds(c, x, y, w, h), nil
}

func buildBubbleChartDefinition(
	chartMap map[string]any,
	title string,
	x, y, w, h int64,
) (charts.ChartDefinition, error) {
	xValues, yValues, sizes, err := requireBubbleSeries(chartMap)
	if err != nil {
		return nil, err
	}
	c := charts.NewBubbleChart(xValues, yValues, sizes).WithTitle(title)
	if w > 0 && h > 0 {
		c = c.Size(w, h).Position(x, y)
	}
	return c, nil
}

func buildCategoryChartDefinition(
	chartType string,
	categories []string,
	values []float64,
	title string,
	x, y, w, h int64,
) (charts.ChartDefinition, error) {
	switch chartType {
	case "bar":
		return applyStyledBounds(charts.NewBarChart(categories, values).WithTitle(title), x, y, w, h), nil
	case "bar_horizontal", "bar-horizontal":
		return applyStyledBounds(charts.NewBarHorizontalChart(categories, values).WithTitle(title), x, y, w, h), nil
	case "bar_stacked", "bar-stacked":
		return applyStyledBounds(charts.NewBarStackedChart(categories, values).WithTitle(title), x, y, w, h), nil
	case "bar_stacked_100", "bar-stacked-100":
		return applyStyledBounds(charts.NewBarStacked100Chart(categories, values).WithTitle(title), x, y, w, h), nil
	case "line":
		return applyStyledBounds(charts.NewLineChart(categories, values).WithTitle(title), x, y, w, h), nil
	case "line_markers", "line-markers":
		return applyStyledBounds(charts.NewLineMarkersChart(categories, values).WithTitle(title), x, y, w, h), nil
	case "line_stacked", "line-stacked":
		return applyStyledBounds(charts.NewLineStackedChart(categories, values).WithTitle(title), x, y, w, h), nil
	case "area":
		return applyStyledBounds(charts.NewAreaChart(categories, values).WithTitle(title), x, y, w, h), nil
	case "area_stacked", "area-stacked":
		return applyStyledBounds(charts.NewAreaStackedChart(categories, values).WithTitle(title), x, y, w, h), nil
	case "area_stacked_100", "area-stacked-100":
		return applyStyledBounds(charts.NewAreaStacked100Chart(categories, values).WithTitle(title), x, y, w, h), nil
	case "pie":
		return applyStyledBounds(charts.NewPieChart(categories, values).WithTitle(title), x, y, w, h), nil
	case "doughnut":
		return applyStyledBounds(charts.NewDoughnutChart(categories, values).WithTitle(title), x, y, w, h), nil
	case "radar":
		return applyStyledBounds(charts.NewRadarChart(categories, values).WithTitle(title), x, y, w, h), nil
	case "radar_filled", "radar-filled":
		return applyStyledBounds(charts.NewRadarFilledChart(categories, values).WithTitle(title), x, y, w, h), nil
	default:
		return nil, fmt.Errorf("unsupported chart type: %q", chartType)
	}
}

func requireCategorySeries(chartMap map[string]any) ([]string, []float64, error) {
	categories, ok := parseStringSlice(chartMap["categories"])
	if !ok {
		return nil, nil, errors.New("field categories must be an array of strings")
	}
	values, ok := parseFloat64Slice(chartMap["values"])
	if !ok {
		return nil, nil, errors.New("field values must be an array of numbers")
	}
	if len(categories) != len(values) {
		return nil, nil, errors.New("field categories and values must have equal length")
	}
	return categories, values, nil
}

func requireXYSeries(chartMap map[string]any) ([]float64, []float64, error) {
	xValues, ok := parseFloat64Slice(chartMap["x_values"])
	if !ok {
		return nil, nil, errors.New("field x_values must be an array of numbers")
	}
	yValues, ok := parseFloat64Slice(chartMap["y_values"])
	if !ok {
		return nil, nil, errors.New("field y_values must be an array of numbers")
	}
	if len(xValues) != len(yValues) {
		return nil, nil, errors.New("field x_values and y_values must have equal length")
	}
	return xValues, yValues, nil
}

func requireBubbleSeries(chartMap map[string]any) ([]float64, []float64, []float64, error) {
	xValues, yValues, err := requireXYSeries(chartMap)
	if err != nil {
		return nil, nil, nil, err
	}
	sizes, ok := parseFloat64Slice(chartMap["sizes"])
	if !ok {
		return nil, nil, nil, errors.New("field sizes must be an array of numbers")
	}
	if len(sizes) != len(xValues) {
		return nil, nil, nil, errors.New("field sizes must match x_values/y_values length")
	}
	return xValues, yValues, sizes, nil
}

func parseStringSlice(raw any) ([]string, bool) {
	values, ok := raw.([]any)
	if !ok {
		return nil, false
	}
	out := make([]string, 0, len(values))
	for _, v := range values {
		s, ok := v.(string)
		if !ok {
			return nil, false
		}
		out = append(out, s)
	}
	return out, true
}

func parseFloat64Slice(raw any) ([]float64, bool) {
	values, ok := raw.([]any)
	if !ok {
		return nil, false
	}
	out := make([]float64, 0, len(values))
	for _, v := range values {
		num, ok := v.(float64)
		if !ok {
			return nil, false
		}
		if math.IsNaN(num) || math.IsInf(num, 0) {
			return nil, false
		}
		out = append(out, num)
	}
	return out, true
}

func applyStyledBounds[T interface {
	Size(cx styling.Length, cy styling.Length) T
	Position(x styling.Length, y styling.Length) T
}](chart T, x, y, w, h int64) T {
	if w > 0 && h > 0 {
		return chart.Size(styling.Emu(w), styling.Emu(h)).Position(styling.Emu(x), styling.Emu(y))
	}
	return chart
}
