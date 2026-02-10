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

func (c BubbleChart) WithAxisTitles(xAxisTitle string, yAxisTitle string) BubbleChart {
	c.CategoryAxisTitle = strings.TrimSpace(xAxisTitle)
	c.ValueAxisTitle = strings.TrimSpace(yAxisTitle)
	return c
}

func (c BubbleChart) WithMajorGridlines(show bool) BubbleChart {
	c.ShowMajorGridlines = show
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

func (c BubbleChart) WithValueRange(min float64, max float64) BubbleChart {
	c.MinValue = &min
	c.MaxValue = &max
	return c
}

func (c BubbleChart) WithBubbleScale(scale int) BubbleChart {
	c.BubbleScale = scale
	return c
}
