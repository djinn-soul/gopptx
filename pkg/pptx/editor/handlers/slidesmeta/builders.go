package slidesmeta

import (
	"errors"
	"fmt"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

var (
	ErrUnsupportedChartType = errors.New("unsupported chart type")
	ErrUnknownThemeName     = errors.New("unknown theme name")
)

func BuildChartDefinition(request editorcommand.AddChartRequest) (charts.ChartDefinition, error) {
	switch canonicalChartType(request.ChartType) {
	case "bar":
		return withBounds(charts.NewBarChart(request.Categories, request.Values).WithTitle(request.Title), request), nil
	case "barHorizontal":
		return withBounds(charts.NewBarHorizontalChart(request.Categories, request.Values).WithTitle(request.Title), request), nil
	case "barStacked":
		return withBounds(charts.NewBarStackedChart(request.Categories, request.Values).WithTitle(request.Title), request), nil
	case "barStacked100":
		return withBounds(charts.NewBarStacked100Chart(request.Categories, request.Values).WithTitle(request.Title), request), nil
	case "line":
		return withBounds(charts.NewLineChart(request.Categories, request.Values).WithTitle(request.Title), request), nil
	case "lineMarkers":
		return withBounds(charts.NewLineMarkersChart(request.Categories, request.Values).WithTitle(request.Title), request), nil
	case "lineStacked":
		return withBounds(charts.NewLineStackedChart(request.Categories, request.Values).WithTitle(request.Title), request), nil
	case "scatter":
		xValues, yValues := syntheticXYSeries(request.Values)
		return withBounds(charts.NewScatterChart(xValues, yValues).WithTitle(request.Title), request), nil
	case "area":
		return withBounds(charts.NewAreaChart(request.Categories, request.Values).WithTitle(request.Title), request), nil
	case "areaStacked":
		return withBounds(charts.NewAreaStackedChart(request.Categories, request.Values).WithTitle(request.Title), request), nil
	case "areaStacked100":
		return withBounds(charts.NewAreaStacked100Chart(request.Categories, request.Values).WithTitle(request.Title), request), nil
	case "pie":
		return withBounds(charts.NewPieChart(request.Categories, request.Values).WithTitle(request.Title), request), nil
	case "doughnut":
		return withBounds(charts.NewDoughnutChart(request.Categories, request.Values).WithTitle(request.Title), request), nil
	case "bubble":
		xValues, yValues := syntheticXYSeries(request.Values)
		sizes := syntheticBubbleSizes(request.Values)
		return withBubbleBounds(charts.NewBubbleChart(xValues, yValues, sizes).WithTitle(request.Title), request), nil
	case "radar":
		return withBounds(charts.NewRadarChart(request.Categories, request.Values).WithTitle(request.Title), request), nil
	case "radarFilled":
		return withBounds(charts.NewRadarFilledChart(request.Categories, request.Values).WithTitle(request.Title), request), nil
	case "stockHLC":
		high, low, closeVals := syntheticStockTriplet(request.Values)
		return withBounds(charts.NewStockHLCChart(request.Categories, high, low, closeVals).WithTitle(request.Title), request), nil
	case "stockOHLC":
		openVals, high, low, closeVals := syntheticStockQuad(request.Values)
		return withBounds(charts.NewStockOHLCChart(request.Categories, openVals, high, low, closeVals).WithTitle(request.Title), request), nil
	case "combo":
		barSeries := make([]charts.Series, len(request.BarSeries))
		for i, s := range request.BarSeries {
			barSeries[i] = charts.Series{Name: s.Name, Values: s.Values}
		}
		lineSeries := make([]charts.Series, len(request.LineSeries))
		for i, s := range request.LineSeries {
			lineSeries[i] = charts.Series{Name: s.Name, Values: s.Values}
		}
		return withBounds(charts.NewComboChart(request.Categories, barSeries, lineSeries).WithTitle(request.Title), request), nil
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedChartType, request.ChartType)
	}
}

func canonicalChartType(value string) string {
	switch value {
	case "bar", "BAR", "column", "COLUMN":
		return "bar"
	case "barHorizontal", "bar_horizontal", "bar-horizontal", "BAR_HORIZONTAL":
		return "barHorizontal"
	case "barStacked", "bar_stacked", "bar-stacked", "BAR_STACKED":
		return "barStacked"
	case "barStacked100", "bar_stacked_100", "bar-stacked-100", "BAR_STACKED_100":
		return "barStacked100"
	case "line", "LINE":
		return "line"
	case "lineMarkers", "line_markers", "line-markers", "LINE_MARKERS":
		return "lineMarkers"
	case "lineStacked", "line_stacked", "line-stacked", "LINE_STACKED":
		return "lineStacked"
	case "scatter", "SCATTER":
		return "scatter"
	case "area", "AREA":
		return "area"
	case "areaStacked", "area_stacked", "area-stacked", "AREA_STACKED":
		return "areaStacked"
	case "areaStacked100", "area_stacked_100", "area-stacked-100", "AREA_STACKED_100":
		return "areaStacked100"
	case "pie", "PIE":
		return "pie"
	case "doughnut", "DOUGHNUT":
		return "doughnut"
	case "bubble", "BUBBLE":
		return "bubble"
	case "radar", "RADAR":
		return "radar"
	case "radarFilled", "radar_filled", "radar-filled", "RADAR_FILLED":
		return "radarFilled"
	case "stockHLC", "stock_hlc", "stock-hlc", "STOCK_HLC":
		return "stockHLC"
	case "stockOHLC", "stock_ohlc", "stock-ohlc", "STOCK_OHLC":
		return "stockOHLC"
	case "combo", "COMBO":
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
		high[i] = value + 2
		low[i] = value - 2
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
		openVals[i] = value - 1
		high[i] = value + 2
		low[i] = value - 2
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
