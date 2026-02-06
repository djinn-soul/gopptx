package pptx

import "strings"

// WithSeriesName sets the single-series label in chart legends.
func (c AreaChart) WithSeriesName(name string) AreaChart {
	c.SeriesName = strings.TrimSpace(name)
	return c
}

// WithLegend toggles legend visibility.
func (c AreaChart) WithLegend(show bool) AreaChart {
	c.ShowLegend = show
	return c
}

// WithLegendPosition sets legend position as r/l/t/b.
func (c AreaChart) WithLegendPosition(position string) AreaChart {
	c.LegendPosition = strings.ToLower(strings.TrimSpace(position))
	return c
}

// WithDataLabels toggles value labels on chart points.
func (c AreaChart) WithDataLabels(show bool) AreaChart {
	c.ShowDataLabels = show
	return c
}

// WithAxisTitles sets optional category/value axis titles.
func (c AreaChart) WithAxisTitles(categoryAxisTitle string, valueAxisTitle string) AreaChart {
	c.CategoryAxisTitle = strings.TrimSpace(categoryAxisTitle)
	c.ValueAxisTitle = strings.TrimSpace(valueAxisTitle)
	return c
}

// WithMajorGridlines toggles value-axis major gridlines.
func (c AreaChart) WithMajorGridlines(show bool) AreaChart {
	c.ShowMajorGridlines = show
	return c
}

// WithValueFormat sets the value-axis number format code.
func (c AreaChart) WithValueFormat(format string) AreaChart {
	c.ValueFormat = strings.TrimSpace(format)
	return c
}

// WithValueRange sets optional value-axis min/max.
func (c AreaChart) WithValueRange(min float64, max float64) AreaChart {
	c.MinValue = &min
	c.MaxValue = &max
	return c
}
