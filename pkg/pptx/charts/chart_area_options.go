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

// WithDataLabelPosition sets data-label position:
// ctr|inEnd|inBase|outEnd|bestFit|l|r|t|b.
func (c AreaChart) WithDataLabelPosition(position string) AreaChart {
	c.DataLabels.Position = strings.TrimSpace(position)
	return c
}

// WithDataLabelContent customizes data-label content fields.
func (c AreaChart) WithDataLabelContent(
	showValue bool,
	showCategory bool,
	showSeriesName bool,
	showPercent bool,
) AreaChart {
	c.DataLabels.UseCustom = true
	c.DataLabels.ShowValue = showValue
	c.DataLabels.ShowCategory = showCategory
	c.DataLabels.ShowSeriesName = showSeriesName
	c.DataLabels.ShowPercent = showPercent
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

// WithCategoryMajorGridlines toggles category-axis major gridlines.
func (c AreaChart) WithCategoryMajorGridlines(show bool) AreaChart {
	c.ShowCategoryMajorGridlines = show
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

// WithTickLabelPositions sets category/value axis tick label positions:
// nextTo|low|high|none.
func (c AreaChart) WithTickLabelPositions(categoryPos string, valuePos string) AreaChart {
	c.CategoryTickLabelPosition = strings.TrimSpace(categoryPos)
	c.ValueTickLabelPosition = strings.TrimSpace(valuePos)
	return c
}

// WithAxisCrosses sets category/value axis crosses mode: autoZero|min|max.
func (c AreaChart) WithAxisCrosses(categoryCrosses string, valueCrosses string) AreaChart {
	c.CategoryAxisCrosses = strings.TrimSpace(categoryCrosses)
	c.ValueAxisCrosses = strings.TrimSpace(valueCrosses)
	return c
}

// WithValueRange sets optional value-axis min/max.
func (c AreaChart) WithValueRange(minValue float64, maxValue float64) AreaChart {
	c.MinValue = &minValue
	c.MaxValue = &maxValue
	return c
}
