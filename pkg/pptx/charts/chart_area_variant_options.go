package charts

func (c AreaStackedChart) Position(x int64, y int64) AreaStackedChart {
	c.AreaChart = c.AreaChart.Position(x, y)
	return c
}

func (c AreaStackedChart) Size(cx int64, cy int64) AreaStackedChart {
	c.AreaChart = c.AreaChart.Size(cx, cy)
	return c
}

func (c AreaStackedChart) WithTitle(title string) AreaStackedChart {
	c.AreaChart = c.AreaChart.WithTitle(title)
	return c
}

func (c AreaStackedChart) WithAreaColor(color string) AreaStackedChart {
	c.AreaChart = c.AreaChart.WithAreaColor(color)
	return c
}

func (c AreaStackedChart) WithSeriesName(name string) AreaStackedChart {
	c.AreaChart = c.AreaChart.WithSeriesName(name)
	return c
}

func (c AreaStackedChart) WithLegend(show bool) AreaStackedChart {
	c.AreaChart = c.AreaChart.WithLegend(show)
	return c
}

func (c AreaStackedChart) WithLegendPosition(position string) AreaStackedChart {
	c.AreaChart = c.AreaChart.WithLegendPosition(position)
	return c
}

func (c AreaStackedChart) WithTitleOverlay(overlay bool) AreaStackedChart {
	c.AreaChart = c.AreaChart.WithTitleOverlay(overlay)
	return c
}

func (c AreaStackedChart) WithLegendOverlay(overlay bool) AreaStackedChart {
	c.AreaChart = c.AreaChart.WithLegendOverlay(overlay)
	return c
}

func (c AreaStackedChart) WithDataLabels(show bool) AreaStackedChart {
	c.AreaChart = c.AreaChart.WithDataLabels(show)
	return c
}

func (c AreaStackedChart) WithAxisTitles(categoryAxisTitle string, valueAxisTitle string) AreaStackedChart {
	c.AreaChart = c.AreaChart.WithAxisTitles(categoryAxisTitle, valueAxisTitle)
	return c
}

func (c AreaStackedChart) WithMajorGridlines(show bool) AreaStackedChart {
	c.AreaChart = c.AreaChart.WithMajorGridlines(show)
	return c
}

func (c AreaStackedChart) WithValueFormat(format string) AreaStackedChart {
	c.AreaChart = c.AreaChart.WithValueFormat(format)
	return c
}

func (c AreaStackedChart) WithValueAxisCrossBetween(mode string) AreaStackedChart {
	c.AreaChart = c.AreaChart.WithValueAxisCrossBetween(mode)
	return c
}

func (c AreaStackedChart) WithValueRange(min float64, max float64) AreaStackedChart {
	c.AreaChart = c.AreaChart.WithValueRange(min, max)
	return c
}

func (c AreaStacked100Chart) Position(x int64, y int64) AreaStacked100Chart {
	c.AreaChart = c.AreaChart.Position(x, y)
	return c
}

func (c AreaStacked100Chart) Size(cx int64, cy int64) AreaStacked100Chart {
	c.AreaChart = c.AreaChart.Size(cx, cy)
	return c
}

func (c AreaStacked100Chart) WithTitle(title string) AreaStacked100Chart {
	c.AreaChart = c.AreaChart.WithTitle(title)
	return c
}

func (c AreaStacked100Chart) WithAreaColor(color string) AreaStacked100Chart {
	c.AreaChart = c.AreaChart.WithAreaColor(color)
	return c
}

func (c AreaStacked100Chart) WithSeriesName(name string) AreaStacked100Chart {
	c.AreaChart = c.AreaChart.WithSeriesName(name)
	return c
}

func (c AreaStacked100Chart) WithLegend(show bool) AreaStacked100Chart {
	c.AreaChart = c.AreaChart.WithLegend(show)
	return c
}

func (c AreaStacked100Chart) WithLegendPosition(position string) AreaStacked100Chart {
	c.AreaChart = c.AreaChart.WithLegendPosition(position)
	return c
}

func (c AreaStacked100Chart) WithTitleOverlay(overlay bool) AreaStacked100Chart {
	c.AreaChart = c.AreaChart.WithTitleOverlay(overlay)
	return c
}

func (c AreaStacked100Chart) WithLegendOverlay(overlay bool) AreaStacked100Chart {
	c.AreaChart = c.AreaChart.WithLegendOverlay(overlay)
	return c
}

func (c AreaStacked100Chart) WithDataLabels(show bool) AreaStacked100Chart {
	c.AreaChart = c.AreaChart.WithDataLabels(show)
	return c
}

func (c AreaStacked100Chart) WithAxisTitles(categoryAxisTitle string, valueAxisTitle string) AreaStacked100Chart {
	c.AreaChart = c.AreaChart.WithAxisTitles(categoryAxisTitle, valueAxisTitle)
	return c
}

func (c AreaStacked100Chart) WithMajorGridlines(show bool) AreaStacked100Chart {
	c.AreaChart = c.AreaChart.WithMajorGridlines(show)
	return c
}

func (c AreaStacked100Chart) WithValueFormat(format string) AreaStacked100Chart {
	c.AreaChart = c.AreaChart.WithValueFormat(format)
	return c
}

func (c AreaStacked100Chart) WithValueAxisCrossBetween(mode string) AreaStacked100Chart {
	c.AreaChart = c.AreaChart.WithValueAxisCrossBetween(mode)
	return c
}

func (c AreaStacked100Chart) WithValueRange(min float64, max float64) AreaStacked100Chart {
	c.AreaChart = c.AreaChart.WithValueRange(min, max)
	return c
}
