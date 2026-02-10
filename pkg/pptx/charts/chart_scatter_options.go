package charts

import "strings"

// WithSeriesName sets the single-series label in chart legends.
func (c ScatterChart) WithSeriesName(name string) ScatterChart {
	c.SeriesName = strings.TrimSpace(name)
	return c
}

// WithScatterStyle sets scatter style: marker, lineMarker, smoothMarker.
func (c ScatterChart) WithScatterStyle(style string) ScatterChart {
	c.ScatterStyle = strings.TrimSpace(style)
	return c
}

// WithLegend toggles legend visibility.
func (c ScatterChart) WithLegend(show bool) ScatterChart {
	c.ShowLegend = show
	return c
}

// WithLegendPosition sets legend position as r/l/t/b.
func (c ScatterChart) WithLegendPosition(position string) ScatterChart {
	c.LegendPosition = strings.ToLower(strings.TrimSpace(position))
	return c
}

// WithTitleOverlay toggles title overlay on chart plot area.
func (c ScatterChart) WithTitleOverlay(overlay bool) ScatterChart {
	c.TitleOverlay = overlay
	return c
}

// WithLegendOverlay toggles legend overlay on chart plot area.
func (c ScatterChart) WithLegendOverlay(overlay bool) ScatterChart {
	c.LegendOverlay = overlay
	return c
}

// WithDataLabels toggles value labels on chart points.
func (c ScatterChart) WithDataLabels(show bool) ScatterChart {
	c.ShowDataLabels = show
	return c
}

// WithAxisTitles sets optional x/y axis titles.
func (c ScatterChart) WithAxisTitles(xAxisTitle string, yAxisTitle string) ScatterChart {
	c.CategoryAxisTitle = strings.TrimSpace(xAxisTitle)
	c.ValueAxisTitle = strings.TrimSpace(yAxisTitle)
	return c
}

// WithMajorGridlines toggles primary y-axis major gridlines.
func (c ScatterChart) WithMajorGridlines(show bool) ScatterChart {
	c.ShowMajorGridlines = show
	return c
}

// WithValueFormat sets the value-axis number format code.
func (c ScatterChart) WithValueFormat(format string) ScatterChart {
	c.ValueFormat = strings.TrimSpace(format)
	return c
}

// WithValueAxisCrossBetween sets value-axis crossing mode: between|midCat.
func (c ScatterChart) WithValueAxisCrossBetween(mode string) ScatterChart {
	c.ValueAxisCrossBetween = strings.TrimSpace(mode)
	return c
}

// WithValueRange sets optional y-axis min/max.
func (c ScatterChart) WithValueRange(min float64, max float64) ScatterChart {
	c.MinValue = &min
	c.MaxValue = &max
	return c
}
