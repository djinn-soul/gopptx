package charts

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

// WithTitleOverlay toggles title overlay on chart plot area.
func (c AreaChart) WithTitleOverlay(overlay bool) AreaChart {
	c.TitleOverlay = overlay
	return c
}

// WithLegendOverlay toggles legend overlay on chart plot area.
func (c AreaChart) WithLegendOverlay(overlay bool) AreaChart {
	c.LegendOverlay = overlay
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

// WithValueAxisCrossBetween sets value-axis crossing mode: between|midCat.
func (c AreaChart) WithValueAxisCrossBetween(mode string) AreaChart {
	c.ValueAxisCrossBetween = strings.TrimSpace(mode)
	return c
}

// WithValueRange sets optional value-axis min/max.
func (c AreaChart) WithValueRange(minValue float64, maxValue float64) AreaChart {
	c.MinValue = &minValue
	c.MaxValue = &maxValue
	return c
}
