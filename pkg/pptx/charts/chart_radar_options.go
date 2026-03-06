package charts

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func (c RadarChart) WithSeriesName(name string) RadarChart {
	c.SeriesName = strings.TrimSpace(name)
	return c
}

func (c RadarChart) WithLegend(show bool) RadarChart {
	c.ShowLegend = show
	return c
}

func (c RadarChart) WithLegendPosition(position string) RadarChart {
	c.LegendPosition = strings.ToLower(strings.TrimSpace(position))
	return c
}

func (c RadarChart) WithTitleOverlay(overlay bool) RadarChart {
	c.TitleOverlay = overlay
	return c
}

func (c RadarChart) WithLegendOverlay(overlay bool) RadarChart {
	c.LegendOverlay = overlay
	return c
}

func (c RadarChart) WithDataLabels(show bool) RadarChart {
	c.ShowDataLabels = show
	return c
}

// WithDataLabelPosition sets data-label position:
// ctr|inEnd|inBase|outEnd|bestFit|l|r|t|b.
func (c RadarChart) WithDataLabelPosition(position string) RadarChart {
	c.DataLabels.Position = strings.TrimSpace(position)
	return c
}

// WithDataLabelContent customizes data-label content fields.
func (c RadarChart) WithDataLabelContent(
	showValue bool,
	showCategory bool,
	showSeriesName bool,
	showPercent bool,
) RadarChart {
	c.DataLabels.UseCustom = true
	c.DataLabels.ShowValue = showValue
	c.DataLabels.ShowCategory = showCategory
	c.DataLabels.ShowSeriesName = showSeriesName
	c.DataLabels.ShowPercent = showPercent
	return c
}

func (c RadarChart) WithAxisTitles(categoryAxisTitle string, valueAxisTitle string) RadarChart {
	c.CategoryAxisTitle = strings.TrimSpace(categoryAxisTitle)
	c.ValueAxisTitle = strings.TrimSpace(valueAxisTitle)
	return c
}

func (c RadarChart) WithMajorGridlines(show bool) RadarChart {
	c.ShowMajorGridlines = show
	return c
}

func (c RadarChart) WithValueFormat(format string) RadarChart {
	c.ValueFormat = strings.TrimSpace(format)
	return c
}

func (c RadarChart) WithValueAxisCrossBetween(mode string) RadarChart {
	c.ValueAxisCrossBetween = strings.TrimSpace(mode)
	return c
}

func (c RadarChart) WithValueRange(minValue float64, maxValue float64) RadarChart {
	c.MinValue = &minValue
	c.MaxValue = &maxValue
	return c
}

func (c RadarFilledChart) Position(x styling.Length, y styling.Length) RadarFilledChart {
	c.RadarChart = c.RadarChart.Position(x, y)
	return c
}

func (c RadarFilledChart) Size(cx styling.Length, cy styling.Length) RadarFilledChart {
	c.RadarChart = c.RadarChart.Size(cx, cy)
	return c
}

func (c RadarFilledChart) WithTitle(title string) RadarFilledChart {
	c.RadarChart = c.RadarChart.WithTitle(title)
	return c
}

func (c RadarFilledChart) WithLineColor(color string) RadarFilledChart {
	c.RadarChart = c.RadarChart.WithLineColor(color)
	return c
}

func (c RadarFilledChart) WithSeriesName(name string) RadarFilledChart {
	c.RadarChart = c.RadarChart.WithSeriesName(name)
	return c
}

func (c RadarFilledChart) WithLegend(show bool) RadarFilledChart {
	c.RadarChart = c.RadarChart.WithLegend(show)
	return c
}

func (c RadarFilledChart) WithLegendPosition(position string) RadarFilledChart {
	c.RadarChart = c.RadarChart.WithLegendPosition(position)
	return c
}

func (c RadarFilledChart) WithTitleOverlay(overlay bool) RadarFilledChart {
	c.RadarChart = c.RadarChart.WithTitleOverlay(overlay)
	return c
}

func (c RadarFilledChart) WithLegendOverlay(overlay bool) RadarFilledChart {
	c.RadarChart = c.RadarChart.WithLegendOverlay(overlay)
	return c
}

func (c RadarFilledChart) WithDataLabels(show bool) RadarFilledChart {
	c.RadarChart = c.RadarChart.WithDataLabels(show)
	return c
}

func (c RadarFilledChart) WithDataLabelPosition(position string) RadarFilledChart {
	c.RadarChart = c.RadarChart.WithDataLabelPosition(position)
	return c
}

func (c RadarFilledChart) WithDataLabelContent(
	showValue bool,
	showCategory bool,
	showSeriesName bool,
	showPercent bool,
) RadarFilledChart {
	c.RadarChart = c.RadarChart.WithDataLabelContent(showValue, showCategory, showSeriesName, showPercent)
	return c
}

func (c RadarFilledChart) WithAxisTitles(categoryAxisTitle string, valueAxisTitle string) RadarFilledChart {
	c.RadarChart = c.RadarChart.WithAxisTitles(categoryAxisTitle, valueAxisTitle)
	return c
}

func (c RadarFilledChart) WithMajorGridlines(show bool) RadarFilledChart {
	c.RadarChart = c.RadarChart.WithMajorGridlines(show)
	return c
}

func (c RadarFilledChart) WithValueFormat(format string) RadarFilledChart {
	c.RadarChart = c.RadarChart.WithValueFormat(format)
	return c
}

func (c RadarFilledChart) WithValueAxisCrossBetween(mode string) RadarFilledChart {
	c.RadarChart = c.RadarChart.WithValueAxisCrossBetween(mode)
	return c
}

func (c RadarFilledChart) WithValueRange(minValue float64, maxValue float64) RadarFilledChart {
	c.RadarChart = c.RadarChart.WithValueRange(minValue, maxValue)
	return c
}
