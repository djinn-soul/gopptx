package slidesmeta

import (
	"errors"
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

var (
	ErrUnsupportedChartType = errors.New("unsupported chart type")
	ErrUnknownThemeName     = errors.New("unknown theme name")
)

const (
	chartTypeBar         = "bar"
	chartTypeLine        = "line"
	chartTypeScatter     = "scatter"
	chartTypeArea        = "area"
	chartTypePie         = "pie"
	chartTypeDoughnut    = "doughnut"
	chartTypeBubble      = "bubble"
	chartTypeRadar       = "radar"
	chartTypeStockDelta  = 2
	defaultStockOpenDiff = 1
)

//nolint:funlen
func BuildChartDefinition(request editorcommand.AddChartRequest) (charts.ChartDefinition, error) {
	switch canonicalChartType(request.ChartType) {
	case chartTypeBar:
		return withBounds(charts.NewBarChart(request.Categories, request.Values).WithTitle(request.Title), request), nil
	case "barHorizontal":
		return withBounds(
			charts.NewBarHorizontalChart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case "barStacked":
		return withBounds(
			charts.NewBarStackedChart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case "barStacked100":
		return withBounds(
			charts.NewBarStacked100Chart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case chartTypeLine:
		return withBounds(
			charts.NewLineChart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case "lineMarkers":
		return withBounds(
			charts.NewLineMarkersChart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case "lineStacked":
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
	case "areaStacked":
		return withBounds(
			charts.NewAreaStackedChart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case "areaStacked100":
		return withBounds(
			charts.NewAreaStacked100Chart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case chartTypePie:
		return withBounds(charts.NewPieChart(request.Categories, request.Values).WithTitle(request.Title), request), nil
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
	case "radarFilled":
		return withBounds(
			charts.NewRadarFilledChart(request.Categories, request.Values).WithTitle(request.Title),
			request,
		), nil
	case "stockHLC":
		high, low, closeVals := syntheticStockTriplet(request.Values)
		return withBounds(
			charts.NewStockHLCChart(request.Categories, high, low, closeVals).WithTitle(request.Title),
			request,
		), nil
	case "stockOHLC":
		openVals, high, low, closeVals := syntheticStockQuad(request.Values)
		return withBounds(
			charts.NewStockOHLCChart(request.Categories, openVals, high, low, closeVals).WithTitle(request.Title),
			request,
		), nil
	case "combo":
		barSeries := make([]charts.Series, len(request.BarSeries))
		for i, s := range request.BarSeries {
			barSeries[i] = charts.Series{Name: s.Name, Values: s.Values}
		}
		lineSeries := make([]charts.Series, len(request.LineSeries))
		for i, s := range request.LineSeries {
			lineSeries[i] = charts.Series{Name: s.Name, Values: s.Values}
		}
		return withBounds(
			charts.NewComboChart(request.Categories, barSeries, lineSeries).WithTitle(request.Title),
			request,
		), nil
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedChartType, request.ChartType)
	}
}

func canonicalChartType(value string) string {
	switch strings.ToLower(value) {
	case chartTypeBar, "column":
		return chartTypeBar
	case "barhorizontal", "bar_horizontal", "bar-horizontal":
		return "barHorizontal"
	case "barstacked", "bar_stacked", "bar-stacked":
		return "barStacked"
	case "barstacked100", "bar_stacked_100", "bar-stacked-100":
		return "barStacked100"
	case chartTypeLine:
		return chartTypeLine
	case "linemarkers", "line_markers", "line-markers":
		return "lineMarkers"
	case "linestacked", "line_stacked", "line-stacked":
		return "lineStacked"
	case chartTypeScatter:
		return chartTypeScatter
	case chartTypeArea:
		return chartTypeArea
	case "areastacked", "area_stacked", "area-stacked":
		return "areaStacked"
	case "areastacked100", "area_stacked_100", "area-stacked-100":
		return "areaStacked100"
	case chartTypePie:
		return chartTypePie
	case chartTypeDoughnut:
		return chartTypeDoughnut
	case chartTypeBubble:
		return chartTypeBubble
	case chartTypeRadar:
		return chartTypeRadar
	case "radarfilled", "radar_filled", "radar-filled":
		return "radarFilled"
	case "stockhlc", "stock_hlc", "stock-hlc":
		return "stockHLC"
	case "stockohlc", "stock_ohlc", "stock-ohlc":
		return "stockOHLC"
	case "combo":
		return "combo"
	default:
		return value
	}
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

func BuildSlideContent(request editorcommand.UpdateSlideRequest, currentTitle string) elements.SlideContent {
	title := request.Title
	if title == "" {
		title = currentTitle
	}

	slide := elements.NewSlide(title)
	if request.Layout != "" {
		slide = slide.WithLayout(request.Layout)
	}
	for _, bullet := range request.Bullets {
		slide = slide.AddBullet(bullet)
	}
	return slide
}

func ResolveThemeByName(name string) (styling.Theme, error) {
	switch name {
	case "Corporate":
		return styling.ThemeCorporate, nil
	case "Modern":
		return styling.ThemeModern, nil
	case "Vibrant":
		return styling.ThemeVibrant, nil
	case "Dark":
		return styling.ThemeDark, nil
	case "Nature":
		return styling.ThemeNature, nil
	case "Tech":
		return styling.ThemeTech, nil
	case "Carbon":
		return styling.ThemeCarbon, nil
	default:
		return styling.Theme{}, fmt.Errorf("%w: %q", ErrUnknownThemeName, name)
	}
}
