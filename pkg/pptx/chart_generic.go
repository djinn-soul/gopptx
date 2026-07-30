package pptx

import (
	"fmt"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/enums"
)

// Chart is an alias for charts.Chart: the sealed interface satisfied by every
// chart type. Use it to hold a chart whose type is only known at runtime, and
// pass it to SlideContent.WithChart.
type Chart = charts.Chart

// NewChart builds a category/value chart of the requested type, so that a chart
// type chosen at runtime does not need a switch over the nineteen NewBarChart /
// NewPieChart / ... constructors.
//
// Scatter, bubble, stock and combo charts do not take categories and values;
// NewChart rejects them and they must be built through their dedicated
// constructors (NewScatterChart, NewBubbleChart, NewStockHLCChart,
// NewStockOHLCChart, NewComboChart), which take the data shape each one needs.
//
// The result is a Chart, which carries no fluent option methods. To set a title,
// legend, colors and the like, either use the concrete constructor directly or
// type-assert the result back to its concrete type.
func NewChart(chartType enums.XLChartType, categories []string, values []float64) (Chart, error) {
	switch chartType {
	case enums.XLChartTypeBar:
		return charts.NewBarChart(categories, values), nil
	case enums.XLChartTypeBarHorizontal:
		return charts.NewBarHorizontalChart(categories, values), nil
	case enums.XLChartTypeBarStacked:
		return charts.NewBarStackedChart(categories, values), nil
	case enums.XLChartTypeBarStacked100:
		return charts.NewBarStacked100Chart(categories, values), nil
	case enums.XLChartTypeLine:
		return charts.NewLineChart(categories, values), nil
	case enums.XLChartTypeLineMarkers:
		return charts.NewLineMarkersChart(categories, values), nil
	case enums.XLChartTypeLineStacked:
		return charts.NewLineStackedChart(categories, values), nil
	case enums.XLChartTypeArea:
		return charts.NewAreaChart(categories, values), nil
	case enums.XLChartTypeAreaStacked:
		return charts.NewAreaStackedChart(categories, values), nil
	case enums.XLChartTypeAreaStacked100:
		return charts.NewAreaStacked100Chart(categories, values), nil
	case enums.XLChartTypePie:
		return charts.NewPieChart(categories, values), nil
	case enums.XLChartTypeThreeDPie:
		return charts.NewPie3DChart(categories, values), nil
	case enums.XLChartTypeDoughnut:
		return charts.NewDoughnutChart(categories, values), nil
	case enums.XLChartTypeRadar:
		return charts.NewRadarChart(categories, values), nil
	case enums.XLChartTypeRadarFilled:
		return charts.NewRadarFilledChart(categories, values), nil
	case enums.XLChartTypeScatter:
		return nil, notCategoryChartError(chartType, "x and y values", "NewScatterChart")
	case enums.XLChartTypeBubble:
		return nil, notCategoryChartError(chartType, "x, y and size values", "NewBubbleChart")
	case enums.XLChartTypeStockHLC:
		return nil, notCategoryChartError(chartType, "high, low and close values", "NewStockHLCChart")
	case enums.XLChartTypeStockOHLC:
		return nil, notCategoryChartError(chartType, "open, high, low and close values", "NewStockOHLCChart")
	case enums.XLChartTypeCombo:
		return nil, notCategoryChartError(chartType, "bar and line series", "NewComboChart")
	default:
		return nil, fmt.Errorf("unknown chart type %q", string(chartType))
	}
}

// notCategoryChartError reports that chartType cannot be built from categories
// and values, and names both the data it needs and the constructor that takes it.
func notCategoryChartError(chartType enums.XLChartType, wants string, constructor string) error {
	return fmt.Errorf(
		"chart type %q takes %s, not categories and values: use %s",
		string(chartType), wants, constructor,
	)
}
