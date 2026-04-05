package export

import (
	"strconv"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	parsedChartKindBar = "bar"
	parsedChartKindPie = "pie"
	seriesIndex0       = 0
	seriesIndex1       = 1
	seriesIndex2       = 2
	seriesIndex3       = 3
)

//nolint:funlen
func applyParsedCharts(slide *elements.SlideContent, chartList []parsedChart) {
	for _, pc := range chartList {
		px := styling.Emu(pc.X)
		py := styling.Emu(pc.Y)
		pw := styling.Emu(pc.CX)
		ph := styling.Emu(pc.CY)
		title := nonEmpty(pc.Title, "Chart")
		cats := categoriesForChart(pc)
		vals := valuesForChart(pc)
		seriesColor := ""
		if len(pc.Series) > 0 {
			seriesColor = pc.Series[0].Color
		}
		switch pc.Kind {
		case parsedChartKindBar:
			c := charts.NewBarChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
			if seriesColor != "" {
				c.BarColor = seriesColor
			}
			c.AltText = pc.AltText
			c.IsDecorative = pc.IsDecorative
			slide.Chart = &c
		case "barHorizontal":
			c := charts.NewBarHorizontalChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
			if seriesColor != "" {
				c.BarColor = seriesColor
			}
			c.AltText = pc.AltText
			c.IsDecorative = pc.IsDecorative
			slide.BarHorizontal = &c
		case "barStacked":
			c := charts.NewBarStackedChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
			if seriesColor != "" {
				c.BarColor = seriesColor
			}
			c.AltText = pc.AltText
			c.IsDecorative = pc.IsDecorative
			slide.BarStacked = &c
		case "barStacked100":
			c := charts.NewBarStacked100Chart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
			if seriesColor != "" {
				c.BarColor = seriesColor
			}
			c.AltText = pc.AltText
			c.IsDecorative = pc.IsDecorative
			slide.BarStacked100 = &c
		case "line":
			c := charts.NewLineChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
			if seriesColor != "" {
				c.LineColor = seriesColor
			}
			c.AltText = pc.AltText
			c.IsDecorative = pc.IsDecorative
			slide.Line = &c
		case "lineMarkers":
			c := charts.NewLineMarkersChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
			if seriesColor != "" {
				c.LineColor = seriesColor
			}
			c.AltText = pc.AltText
			c.IsDecorative = pc.IsDecorative
			slide.LineMarkers = &c
		case "lineStacked":
			c := charts.NewLineStackedChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
			if seriesColor != "" {
				c.LineColor = seriesColor
			}
			c.AltText = pc.AltText
			c.IsDecorative = pc.IsDecorative
			slide.LineStacked = &c
		case "scatter":
			xs, ys := scatterForChart(pc)
			c := charts.NewScatterChart(xs, ys).WithTitle(title).Position(px, py).Size(pw, ph)
			if seriesColor != "" {
				c.LineColor = seriesColor
			}
			c.AltText = pc.AltText
			c.IsDecorative = pc.IsDecorative
			slide.Scatter = &c
		case "area":
			c := charts.NewAreaChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
			if seriesColor != "" {
				c.AreaColor = seriesColor
			}
			c.AltText = pc.AltText
			c.IsDecorative = pc.IsDecorative
			slide.Area = &c
		case "areaStacked":
			c := charts.NewAreaStackedChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
			if seriesColor != "" {
				c.AreaColor = seriesColor
			}
			c.AltText = pc.AltText
			c.IsDecorative = pc.IsDecorative
			slide.AreaStacked = &c
		case "areaStacked100":
			c := charts.NewAreaStacked100Chart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
			if seriesColor != "" {
				c.AreaColor = seriesColor
			}
			c.AltText = pc.AltText
			c.IsDecorative = pc.IsDecorative
			slide.AreaStacked100 = &c
		case parsedChartKindPie:
			c := charts.NewPieChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
			c.AltText = pc.AltText
			c.IsDecorative = pc.IsDecorative
			slide.Pie = &c
		case "doughnut":
			c := charts.NewDoughnutChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
			c.AltText = pc.AltText
			c.IsDecorative = pc.IsDecorative
			slide.Doughnut = &c
		case "bubble":
			xs, ys, sizes := bubbleForChart(pc)
			c := charts.NewBubbleChart(xs, ys, sizes).
				WithTitle(title).
				Position(int64(px), int64(py)).
				Size(int64(pw), int64(ph))
			c.AltText = pc.AltText
			c.IsDecorative = pc.IsDecorative
			slide.Bubble = &c
		case "radar":
			c := charts.NewRadarChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
			c.AltText = pc.AltText
			c.IsDecorative = pc.IsDecorative
			slide.Radar = &c
		case "radarFilled":
			c := charts.NewRadarFilledChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
			c.AltText = pc.AltText
			c.IsDecorative = pc.IsDecorative
			slide.RadarFilled = &c
		case "stockHLC":
			high, low, closeVals := stockTriplet(pc)
			c := charts.NewStockHLCChart(cats, high, low, closeVals).WithTitle(title).Position(px, py).Size(pw, ph)
			c.AltText = pc.AltText
			c.IsDecorative = pc.IsDecorative
			slide.StockHLC = &c
		case "stockOHLC":
			openVals, high, low, closeVals := stockQuad(pc)
			c := charts.NewStockOHLCChart(cats, openVals, high, low, closeVals).
				WithTitle(title).
				Position(px, py).
				Size(pw, ph)
			c.AltText = pc.AltText
			c.IsDecorative = pc.IsDecorative
			slide.StockOHLC = &c
		case "combo":
			barSeries, lineSeries := comboSeries(pc)
			c := charts.NewComboChart(cats, barSeries, lineSeries).WithTitle(title).Position(px, py).Size(pw, ph)
			c.AltText = pc.AltText
			c.IsDecorative = pc.IsDecorative
			slide.Combo = &c
		}
	}
}

func categoriesForChart(c parsedChart) []string {
	if len(c.Series) > 0 && len(c.Series[0].Categories) > 0 {
		return c.Series[0].Categories
	}
	vals := valuesForChart(c)
	out := make([]string, len(vals))
	for i := range vals {
		out[i] = "Item " + strconv.Itoa(i+1)
	}
	return out
}

func valuesForChart(c parsedChart) []float64 {
	if len(c.Series) > 0 && len(c.Series[0].Values) > 0 {
		return c.Series[0].Values
	}
	return []float64{1, 2, 3}
}

func scatterForChart(c parsedChart) ([]float64, []float64) {
	if len(c.Series) > 0 && len(c.Series[0].XValues) > 0 && len(c.Series[0].YValues) > 0 {
		return c.Series[0].XValues, c.Series[0].YValues
	}
	return []float64{1, 2, 3}, []float64{1, 2, 3}
}

func bubbleForChart(c parsedChart) ([]float64, []float64, []float64) {
	if len(c.Series) > 0 {
		s := c.Series[0]
		if len(s.XValues) > 0 && len(s.YValues) > 0 && len(s.Sizes) > 0 {
			return s.XValues, s.YValues, s.Sizes
		}
	}
	return []float64{1, 2, 3}, []float64{1, 2, 3}, []float64{12, 18, 24}
}

func stockTriplet(c parsedChart) ([]float64, []float64, []float64) {
	high := pickSeriesValues(c, seriesIndex0)
	low := pickSeriesValues(c, seriesIndex1)
	closeVals := pickSeriesValues(c, seriesIndex2)
	return high, low, closeVals
}

func stockQuad(c parsedChart) ([]float64, []float64, []float64, []float64) {
	openVals := pickSeriesValues(c, seriesIndex0)
	high := pickSeriesValues(c, seriesIndex1)
	low := pickSeriesValues(c, seriesIndex2)
	closeVals := pickSeriesValues(c, seriesIndex3)
	return openVals, high, low, closeVals
}

func comboSeries(c parsedChart) ([]charts.Series, []charts.Series) {
	barVals := pickSeriesValues(c, seriesIndex0)
	lineVals := pickSeriesValues(c, seriesIndex1)
	return []charts.Series{{Name: "Bar", Values: barVals}}, []charts.Series{{Name: "Line", Values: lineVals}}
}

func pickSeriesValues(c parsedChart, idx int) []float64 {
	if idx >= 0 && idx < len(c.Series) && len(c.Series[idx].Values) > 0 {
		return c.Series[idx].Values
	}
	return valuesForChart(c)
}

func nonEmpty(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
