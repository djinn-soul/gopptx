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

// WithDataLabelPosition sets data-label position:
// ctr|inEnd|inBase|outEnd|bestFit|l|r|t|b.
func (c ScatterChart) WithDataLabelPosition(position string) ScatterChart {
	c.DataLabels.Position = strings.TrimSpace(position)
	return c
}

// WithDataLabelContent customizes data-label content fields.
func (c ScatterChart) WithDataLabelContent(
	showValue bool,
	showCategory bool,
	showSeriesName bool,
	showPercent bool,
) ScatterChart {
	c.DataLabels.UseCustom = true
	c.DataLabels.ShowValue = showValue
	c.DataLabels.ShowCategory = showCategory
	c.DataLabels.ShowSeriesName = showSeriesName
	c.DataLabels.ShowPercent = showPercent
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

// WithCategoryMajorGridlines toggles primary x-axis major gridlines.
func (c ScatterChart) WithCategoryMajorGridlines(show bool) ScatterChart {
	c.ShowCategoryMajorGridlines = show
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

// WithTickLabelPositions sets x/y axis tick label positions:
// nextTo|low|high|none.
func (c ScatterChart) WithTickLabelPositions(xAxisPos string, yAxisPos string) ScatterChart {
	c.CategoryTickLabelPosition = strings.TrimSpace(xAxisPos)
	c.ValueTickLabelPosition = strings.TrimSpace(yAxisPos)
	return c
}

// WithAxisCrosses sets x/y axis crosses mode: autoZero|min|max.
func (c ScatterChart) WithAxisCrosses(xAxisCrosses string, yAxisCrosses string) ScatterChart {
	c.CategoryAxisCrosses = strings.TrimSpace(xAxisCrosses)
	c.ValueAxisCrosses = strings.TrimSpace(yAxisCrosses)
	return c
}

// WithValueRange sets optional y-axis min/max.
func (c ScatterChart) WithValueRange(minValue float64, maxValue float64) ScatterChart {
	c.MinValue = &minValue
	c.MaxValue = &maxValue
	return c
}
