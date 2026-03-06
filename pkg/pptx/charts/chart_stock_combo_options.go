package charts

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func (c StockHLCChart) Position(x styling.Length, y styling.Length) StockHLCChart {
	c.X = x
	c.Y = y
	return c
}

func (c StockHLCChart) Size(cx styling.Length, cy styling.Length) StockHLCChart {
	c.CX = cx
	c.CY = cy
	return c
}

func (c StockHLCChart) WithTitle(title string) StockHLCChart {
	c.Title = title
	return c
}

func (c StockHLCChart) WithLegend(show bool) StockHLCChart {
	c.ShowLegend = show
	return c
}

func (c StockHLCChart) WithLegendPosition(position string) StockHLCChart {
	c.LegendPosition = strings.ToLower(strings.TrimSpace(position))
	return c
}

func (c StockHLCChart) WithTitleOverlay(overlay bool) StockHLCChart {
	c.TitleOverlay = overlay
	return c
}

func (c StockHLCChart) WithLegendOverlay(overlay bool) StockHLCChart {
	c.LegendOverlay = overlay
	return c
}

func (c StockHLCChart) WithDataLabels(show bool) StockHLCChart {
	c.ShowDataLabels = show
	return c
}

func (c StockHLCChart) WithDataLabelPosition(position string) StockHLCChart {
	c.DataLabels.Position = strings.TrimSpace(position)
	return c
}

func (c StockHLCChart) WithDataLabelContent(
	showValue bool,
	showCategory bool,
	showSeriesName bool,
	showPercent bool,
) StockHLCChart {
	c.DataLabels.UseCustom = true
	c.DataLabels.ShowValue = showValue
	c.DataLabels.ShowCategory = showCategory
	c.DataLabels.ShowSeriesName = showSeriesName
	c.DataLabels.ShowPercent = showPercent
	return c
}

func (c StockHLCChart) WithAxisTitles(categoryAxisTitle string, valueAxisTitle string) StockHLCChart {
	c.CategoryAxisTitle = strings.TrimSpace(categoryAxisTitle)
	c.ValueAxisTitle = strings.TrimSpace(valueAxisTitle)
	return c
}

func (c StockHLCChart) WithMajorGridlines(show bool) StockHLCChart {
	c.ShowMajorGridlines = show
	return c
}

func (c StockHLCChart) WithValueFormat(format string) StockHLCChart {
	c.ValueFormat = strings.TrimSpace(format)
	return c
}

func (c StockHLCChart) WithValueAxisCrossBetween(mode string) StockHLCChart {
	c.ValueAxisCrossBetween = strings.TrimSpace(mode)
	return c
}

func (c StockHLCChart) WithValueRange(minValue float64, maxValue float64) StockHLCChart {
	c.MinValue = &minValue
	c.MaxValue = &maxValue
	return c
}

func (c StockOHLCChart) Position(x styling.Length, y styling.Length) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.Position(x, y)
	return c
}

func (c StockOHLCChart) Size(cx styling.Length, cy styling.Length) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.Size(cx, cy)
	return c
}

func (c StockOHLCChart) WithTitle(title string) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithTitle(title)
	return c
}

func (c StockOHLCChart) WithLegend(show bool) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithLegend(show)
	return c
}

func (c StockOHLCChart) WithLegendPosition(position string) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithLegendPosition(position)
	return c
}

func (c StockOHLCChart) WithTitleOverlay(overlay bool) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithTitleOverlay(overlay)
	return c
}

func (c StockOHLCChart) WithLegendOverlay(overlay bool) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithLegendOverlay(overlay)
	return c
}

func (c StockOHLCChart) WithDataLabels(show bool) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithDataLabels(show)
	return c
}

func (c StockOHLCChart) WithDataLabelPosition(position string) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithDataLabelPosition(position)
	return c
}

func (c StockOHLCChart) WithDataLabelContent(
	showValue bool,
	showCategory bool,
	showSeriesName bool,
	showPercent bool,
) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithDataLabelContent(showValue, showCategory, showSeriesName, showPercent)
	return c
}

func (c StockOHLCChart) WithAxisTitles(categoryAxisTitle string, valueAxisTitle string) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithAxisTitles(categoryAxisTitle, valueAxisTitle)
	return c
}

func (c StockOHLCChart) WithMajorGridlines(show bool) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithMajorGridlines(show)
	return c
}

func (c StockOHLCChart) WithValueFormat(format string) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithValueFormat(format)
	return c
}

func (c StockOHLCChart) WithValueAxisCrossBetween(mode string) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithValueAxisCrossBetween(mode)
	return c
}

func (c StockOHLCChart) WithValueRange(minValue float64, maxValue float64) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithValueRange(minValue, maxValue)
	return c
}

func (c ComboChart) Position(x styling.Length, y styling.Length) ComboChart {
	c.X = x
	c.Y = y
	return c
}

func (c ComboChart) Size(cx styling.Length, cy styling.Length) ComboChart {
	c.CX = cx
	c.CY = cy
	return c
}

func (c ComboChart) WithTitle(title string) ComboChart {
	c.Title = title
	return c
}

func (c ComboChart) WithLegend(show bool) ComboChart {
	c.ShowLegend = show
	return c
}

func (c ComboChart) WithLegendPosition(position string) ComboChart {
	c.LegendPosition = strings.ToLower(strings.TrimSpace(position))
	return c
}

func (c ComboChart) WithTitleOverlay(overlay bool) ComboChart {
	c.TitleOverlay = overlay
	return c
}

func (c ComboChart) WithLegendOverlay(overlay bool) ComboChart {
	c.LegendOverlay = overlay
	return c
}

func (c ComboChart) WithDataLabels(show bool) ComboChart {
	c.ShowDataLabels = show
	return c
}

func (c ComboChart) WithDataLabelPosition(position string) ComboChart {
	c.DataLabels.Position = strings.TrimSpace(position)
	return c
}

func (c ComboChart) WithDataLabelContent(
	showValue bool,
	showCategory bool,
	showSeriesName bool,
	showPercent bool,
) ComboChart {
	c.DataLabels.UseCustom = true
	c.DataLabels.ShowValue = showValue
	c.DataLabels.ShowCategory = showCategory
	c.DataLabels.ShowSeriesName = showSeriesName
	c.DataLabels.ShowPercent = showPercent
	return c
}

func (c ComboChart) WithAxisTitles(categoryAxisTitle string, valueAxisTitle string) ComboChart {
	c.CategoryAxisTitle = strings.TrimSpace(categoryAxisTitle)
	c.ValueAxisTitle = strings.TrimSpace(valueAxisTitle)
	return c
}

func (c ComboChart) WithMajorGridlines(show bool) ComboChart {
	c.ShowMajorGridlines = show
	return c
}

func (c ComboChart) WithValueFormat(format string) ComboChart {
	c.ValueFormat = strings.TrimSpace(format)
	return c
}

func (c ComboChart) WithValueAxisCrossBetween(mode string) ComboChart {
	c.ValueAxisCrossBetween = strings.TrimSpace(mode)
	return c
}

func (c ComboChart) WithValueRange(minValue float64, maxValue float64) ComboChart {
	c.MinValue = &minValue
	c.MaxValue = &maxValue
	return c
}
