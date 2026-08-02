package slidesmeta

import (
	"errors"
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

var (
	ErrUnsupportedChartType = errors.New("unsupported chart type")
	ErrUnknownThemeName     = errors.New("unknown theme name")
)

const (
	chartTypeBar      = "bar"
	chartTypeLine     = "line"
	chartTypeScatter  = "scatter"
	chartTypeArea     = "area"
	chartTypePie      = "pie"
	chartTypePie3D    = "pie3D"
	chartTypeDoughnut = "doughnut"
	chartTypeBubble   = "bubble"
	chartTypeRadar    = "radar"

	chartTypeBarHorizontal  = "barHorizontal"
	chartTypeBarStacked     = "barStacked"
	chartTypeBarStacked100  = "barStacked100"
	chartTypeLineMarkers    = "lineMarkers"
	chartTypeLineStacked    = "lineStacked"
	chartTypeAreaStacked    = "areaStacked"
	chartTypeAreaStacked100 = "areaStacked100"
	chartTypeRadarFilled    = "radarFilled"
	chartTypeStockHLC       = "stockHLC"
	chartTypeStockOHLC      = "stockOHLC"
	chartTypeCombo          = "combo"
	chartTypeColumn3D       = "column3D"
	chartTypeBar3D          = "bar3D"
	chartTypeLine3D         = "line3D"
	chartTypeArea3D         = "area3D"

	chartTypeStockDelta  = 2
	defaultStockOpenDiff = 1
)

//nolint:funlen
func BuildChartDefinition(request editorcommand.AddChartRequest) (charts.ChartDefinition, error) {
	switch canonicalChartType(request.ChartType) {
	case chartTypeBar:
		return withBounds(charts.NewBarChart(request.Categories, request.Values).WithTitle(request.Title), request), nil
	case chartTypeBarHorizontal:
		return withBounds(
			charts.NewBarHorizontalChart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case chartTypeBarStacked:
		return withBounds(
			charts.NewBarStackedChart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case chartTypeBarStacked100:
		return withBounds(
			charts.NewBarStacked100Chart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case chartTypeLine:
		return withBounds(
			charts.NewLineChart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case chartTypeLineMarkers:
		return withBounds(
			charts.NewLineMarkersChart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case chartTypeLineStacked:
		return withBounds(
			charts.NewLineStackedChart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case chartTypeScatter:
		xValues, yValues := syntheticXYSeries(request.Values)
		return withBounds(charts.NewScatterChart(xValues, yValues).WithTitle(request.Title), request), nil
	case chartTypeArea:
		return withBounds(
			charts.NewAreaChart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case chartTypeAreaStacked:
		return withBounds(
			charts.NewAreaStackedChart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case chartTypeAreaStacked100:
		return withBounds(
			charts.NewAreaStacked100Chart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case chartTypePie:
		return withBounds(charts.NewPieChart(request.Categories, request.Values).WithTitle(request.Title), request), nil
	case chartTypePie3D:
		return withBounds(
			charts.NewPie3DChart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case chartTypeColumn3D:
		return withBounds(
			charts.NewColumn3DChart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case chartTypeBar3D:
		return withBounds(
			charts.NewBar3DChart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case chartTypeLine3D:
		return withBounds(
			charts.NewLine3DChart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case chartTypeArea3D:
		return withBounds(
			charts.NewArea3DChart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case chartTypeDoughnut:
		return withBounds(
			charts.NewDoughnutChart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case chartTypeBubble:
		xValues, yValues := syntheticXYSeries(request.Values)
		sizes := syntheticBubbleSizes(request.Values)
		return withBubbleBounds(charts.NewBubbleChart(xValues, yValues, sizes).WithTitle(request.Title), request), nil
	case chartTypeRadar:
		return withBounds(
			charts.NewRadarChart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case chartTypeRadarFilled:
		return withBounds(
			charts.NewRadarFilledChart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case chartTypeStockHLC:
		high, low, closeVals := syntheticStockTriplet(request.Values)
		return withBounds(
			charts.NewStockHLCChart(request.Categories, high, low, closeVals).WithTitle(request.Title),
			request,
		), nil
	case chartTypeStockOHLC:
		openVals, high, low, closeVals := syntheticStockQuad(request.Values)
		return withBounds(
			charts.NewStockOHLCChart(request.Categories, openVals, high, low, closeVals).WithTitle(request.Title),
			request,
		), nil
	case chartTypeCombo:
		barSeries := make([]charts.Series, len(request.BarSeries))
		for i, s := range request.BarSeries {
			barSeries[i] = charts.Series{Name: s.Name, Values: s.Values}
		}
		lineSeries := make([]charts.Series, len(request.LineSeries))
		for i, s := range request.LineSeries {
			lineSeries[i] = charts.Series{
				Name:          s.Name,
				Values:        s.Values,
				SecondaryAxis: s.SecondaryAxis,
			}
		}
		combo := charts.NewComboChart(request.Categories, barSeries, lineSeries).
			WithTitle(request.Title).
			WithSecondaryAxis(request.SecondaryAxis)
		if request.SecondaryValueAxisTitle != "" {
			combo = combo.WithSecondaryAxisTitle(request.SecondaryValueAxisTitle)
		}
		return withBounds(combo, request), nil
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedChartType, request.ChartType)
	}
}

// chartTypeAliases maps every accepted spelling of a chart type to its
// canonical token. A table keeps the alias list flat and additive: a new chart
// kind adds rows here instead of growing a switch past the statement limit.
//
//nolint:gochecknoglobals // immutable lookup table
var chartTypeAliases = map[string]string{
	"column": chartTypeBar, chartTypeBar: chartTypeBar,
	"barhorizontal": chartTypeBarHorizontal, "bar_horizontal": chartTypeBarHorizontal,
	"bar-horizontal": chartTypeBarHorizontal,
	"barstacked":     chartTypeBarStacked, "bar_stacked": chartTypeBarStacked,
	"bar-stacked":   chartTypeBarStacked,
	"barstacked100": chartTypeBarStacked100, "bar_stacked_100": chartTypeBarStacked100,
	"bar-stacked-100": chartTypeBarStacked100,
	chartTypeLine:     chartTypeLine,
	"linemarkers":     chartTypeLineMarkers, "line_markers": chartTypeLineMarkers,
	"line-markers": chartTypeLineMarkers,
	"linestacked":  chartTypeLineStacked, "line_stacked": chartTypeLineStacked,
	"line-stacked":   chartTypeLineStacked,
	chartTypeScatter: chartTypeScatter,
	chartTypeArea:    chartTypeArea,
	"areastacked":    chartTypeAreaStacked, "area_stacked": chartTypeAreaStacked,
	"area-stacked":   chartTypeAreaStacked,
	"areastacked100": chartTypeAreaStacked100, "area_stacked_100": chartTypeAreaStacked100,
	"area-stacked-100": chartTypeAreaStacked100,
	chartTypePie:       chartTypePie,
	"pie3d":            chartTypePie3D, "three_d_pie": chartTypePie3D, "three-d-pie": chartTypePie3D,
	"column3d": chartTypeColumn3D, "three_d_column": chartTypeColumn3D,
	"three-d-column": chartTypeColumn3D,
	"bar3d":          chartTypeBar3D, "three_d_bar": chartTypeBar3D, "three-d-bar": chartTypeBar3D,
	"line3d": chartTypeLine3D, "three_d_line": chartTypeLine3D, "three-d-line": chartTypeLine3D,
	"area3d": chartTypeArea3D, "three_d_area": chartTypeArea3D, "three-d-area": chartTypeArea3D,
	chartTypeDoughnut: chartTypeDoughnut,
	chartTypeBubble:   chartTypeBubble,
	chartTypeRadar:    chartTypeRadar,
	"radarfilled":     chartTypeRadarFilled, "radar_filled": chartTypeRadarFilled,
	"radar-filled": chartTypeRadarFilled,
	"stockhlc":     chartTypeStockHLC, "stock_hlc": chartTypeStockHLC, "stock-hlc": chartTypeStockHLC,
	"stockohlc": chartTypeStockOHLC, "stock_ohlc": chartTypeStockOHLC,
	"stock-ohlc":   chartTypeStockOHLC,
	chartTypeCombo: chartTypeCombo,
}

// canonicalChartType resolves any accepted spelling to its canonical token,
// returning the input unchanged when it is not a known alias.
func canonicalChartType(value string) string {
	if canonical, ok := chartTypeAliases[strings.ToLower(value)]; ok {
		return canonical
	}
	return value
}

func withBounds[T interface {
	Size(cx styling.Length, cy styling.Length) T
	Position(x styling.Length, y styling.Length) T
}](chart T, request editorcommand.AddChartRequest) T {
	if request.W > 0 {
		return chart.Size(styling.Emu(request.W), styling.Emu(request.H)).
			Position(styling.Emu(request.X), styling.Emu(request.Y))
	}
	return chart
}

func syntheticXYSeries(values []float64) ([]float64, []float64) {
	if len(values) == 0 {
		return []float64{1, 2, 3}, []float64{1, 2, 3}
	}
	xValues := make([]float64, len(values))
	yValues := make([]float64, len(values))
	for i, value := range values {
		xValues[i] = float64(i + 1)
		yValues[i] = value
	}
	return xValues, yValues
}

func syntheticBubbleSizes(values []float64) []float64 {
	if len(values) == 0 {
		return []float64{8, 10, 12}
	}
	sizes := make([]float64, len(values))
	for i, value := range values {
		if value <= 0 {
			sizes[i] = 6
			continue
		}
		sizes[i] = value
	}
	return sizes
}

func syntheticStockTriplet(values []float64) ([]float64, []float64, []float64) {
	if len(values) == 0 {
		values = []float64{10, 12, 11}
	}
	high := make([]float64, len(values))
	low := make([]float64, len(values))
	closeVals := make([]float64, len(values))
	for i, value := range values {
		high[i] = value + chartTypeStockDelta
		low[i] = value - chartTypeStockDelta
		closeVals[i] = value
	}
	return high, low, closeVals
}

func syntheticStockQuad(values []float64) ([]float64, []float64, []float64, []float64) {
	if len(values) == 0 {
		values = []float64{10, 12, 11}
	}
	openVals := make([]float64, len(values))
	high := make([]float64, len(values))
	low := make([]float64, len(values))
	closeVals := make([]float64, len(values))
	for i, value := range values {
		openVals[i] = value - defaultStockOpenDiff
		high[i] = value + chartTypeStockDelta
		low[i] = value - chartTypeStockDelta
		closeVals[i] = value
	}
	return openVals, high, low, closeVals
}

func withBubbleBounds(chart charts.BubbleChart, request editorcommand.AddChartRequest) charts.BubbleChart {
	if request.W > 0 {
		return chart.Size(request.W, request.H).Position(request.X, request.Y)
	}
	return chart
}
