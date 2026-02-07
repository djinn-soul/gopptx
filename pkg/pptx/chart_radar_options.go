package pptx

import "strings"

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

func (c RadarChart) WithDataLabels(show bool) RadarChart {
	c.ShowDataLabels = show
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

func (c RadarChart) WithValueRange(min float64, max float64) RadarChart {
	c.MinValue = &min
	c.MaxValue = &max
	return c
}

func (c RadarFilledChart) Position(x int64, y int64) RadarFilledChart {
	c.RadarChart = c.RadarChart.Position(x, y)
	return c
}

func (c RadarFilledChart) Size(cx int64, cy int64) RadarFilledChart {
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

func (c RadarFilledChart) WithDataLabels(show bool) RadarFilledChart {
	c.RadarChart = c.RadarChart.WithDataLabels(show)
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

func (c RadarFilledChart) WithValueRange(min float64, max float64) RadarFilledChart {
	c.RadarChart = c.RadarChart.WithValueRange(min, max)
	return c
}
