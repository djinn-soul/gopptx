package charts

import "strings"

func (c BubbleChart) WithSeriesName(name string) BubbleChart {
	c.SeriesName = strings.TrimSpace(name)
	return c
}

func (c BubbleChart) WithLegend(show bool) BubbleChart {
	c.ShowLegend = show
	return c
}

func (c BubbleChart) WithLegendPosition(position string) BubbleChart {
	c.LegendPosition = strings.ToLower(strings.TrimSpace(position))
	return c
}

func (c BubbleChart) WithTitleOverlay(overlay bool) BubbleChart {
	c.TitleOverlay = overlay
	return c
}

func (c BubbleChart) WithLegendOverlay(overlay bool) BubbleChart {
	c.LegendOverlay = overlay
	return c
}

func (c BubbleChart) WithDataLabels(show bool) BubbleChart {
	c.ShowDataLabels = show
	return c
}

// WithDataLabelPosition sets data-label position:
// ctr|inEnd|inBase|outEnd|bestFit|l|r|t|b.
func (c BubbleChart) WithDataLabelPosition(position string) BubbleChart {
	c.DataLabels.Position = strings.TrimSpace(position)
	return c
}

// WithDataLabelContent customizes data-label content fields.
func (c BubbleChart) WithDataLabelContent(
	showValue bool,
	showCategory bool,
	showSeriesName bool,
	showPercent bool,
	showBubbleSize bool,
) BubbleChart {
	c.DataLabels.UseCustom = true
	c.DataLabels.ShowValue = showValue
	c.DataLabels.ShowCategory = showCategory
	c.DataLabels.ShowSeriesName = showSeriesName
	c.DataLabels.ShowPercent = showPercent
	c.DataLabels.ShowBubbleSize = showBubbleSize
	return c
}

func (c BubbleChart) WithAxisTitles(xAxisTitle string, yAxisTitle string) BubbleChart {
	c.CategoryAxisTitle = strings.TrimSpace(xAxisTitle)
	c.ValueAxisTitle = strings.TrimSpace(yAxisTitle)
	return c
}

func (c BubbleChart) WithMajorGridlines(show bool) BubbleChart {
	c.ShowMajorGridlines = show
	return c
}

// WithCategoryMajorGridlines toggles primary x-axis major gridlines.
func (c BubbleChart) WithCategoryMajorGridlines(show bool) BubbleChart {
	c.ShowCategoryMajorGridlines = show
	return c
}

func (c BubbleChart) WithValueFormat(format string) BubbleChart {
	c.ValueFormat = strings.TrimSpace(format)
	return c
}

func (c BubbleChart) WithValueAxisCrossBetween(mode string) BubbleChart {
	c.ValueAxisCrossBetween = strings.TrimSpace(mode)
	return c
}

// WithTickLabelPositions sets x/y axis tick label positions:
// nextTo|low|high|none.
func (c BubbleChart) WithTickLabelPositions(xAxisPos string, yAxisPos string) BubbleChart {
	c.CategoryTickLabelPosition = strings.TrimSpace(xAxisPos)
	c.ValueTickLabelPosition = strings.TrimSpace(yAxisPos)
	return c
}

// WithAxisCrosses sets x/y axis crosses mode: autoZero|min|max.
func (c BubbleChart) WithAxisCrosses(xAxisCrosses string, yAxisCrosses string) BubbleChart {
	c.CategoryAxisCrosses = strings.TrimSpace(xAxisCrosses)
	c.ValueAxisCrosses = strings.TrimSpace(yAxisCrosses)
	return c
}

func (c BubbleChart) WithValueRange(minValue float64, maxValue float64) BubbleChart {
	c.MinValue = &minValue
	c.MaxValue = &maxValue
	return c
}

func (c BubbleChart) WithBubbleScale(scale int) BubbleChart {
	c.BubbleScale = scale
	return c
}
