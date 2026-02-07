package pptx

func (c LineMarkersChart) Position(x int64, y int64) LineMarkersChart {
	c.LineChart = c.LineChart.Position(x, y)
	return c
}

func (c LineMarkersChart) Size(cx int64, cy int64) LineMarkersChart {
	c.LineChart = c.LineChart.Size(cx, cy)
	return c
}

func (c LineMarkersChart) WithTitle(title string) LineMarkersChart {
	c.LineChart = c.LineChart.WithTitle(title)
	return c
}

func (c LineMarkersChart) WithLineColor(color string) LineMarkersChart {
	c.LineChart = c.LineChart.WithLineColor(color)
	return c
}

func (c LineMarkersChart) WithSeriesName(name string) LineMarkersChart {
	c.LineChart = c.LineChart.WithSeriesName(name)
	return c
}

func (c LineMarkersChart) WithLegend(show bool) LineMarkersChart {
	c.LineChart = c.LineChart.WithLegend(show)
	return c
}

func (c LineMarkersChart) WithLegendPosition(position string) LineMarkersChart {
	c.LineChart = c.LineChart.WithLegendPosition(position)
	return c
}

func (c LineMarkersChart) WithTitleOverlay(overlay bool) LineMarkersChart {
	c.LineChart = c.LineChart.WithTitleOverlay(overlay)
	return c
}

func (c LineMarkersChart) WithLegendOverlay(overlay bool) LineMarkersChart {
	c.LineChart = c.LineChart.WithLegendOverlay(overlay)
	return c
}

func (c LineMarkersChart) WithDataLabels(show bool) LineMarkersChart {
	c.LineChart = c.LineChart.WithDataLabels(show)
	return c
}

func (c LineMarkersChart) WithAxisTitles(categoryAxisTitle string, valueAxisTitle string) LineMarkersChart {
	c.LineChart = c.LineChart.WithAxisTitles(categoryAxisTitle, valueAxisTitle)
	return c
}

func (c LineMarkersChart) WithMajorGridlines(show bool) LineMarkersChart {
	c.LineChart = c.LineChart.WithMajorGridlines(show)
	return c
}

func (c LineMarkersChart) WithValueFormat(format string) LineMarkersChart {
	c.LineChart = c.LineChart.WithValueFormat(format)
	return c
}

func (c LineMarkersChart) WithValueAxisCrossBetween(mode string) LineMarkersChart {
	c.LineChart = c.LineChart.WithValueAxisCrossBetween(mode)
	return c
}

func (c LineMarkersChart) WithValueRange(min float64, max float64) LineMarkersChart {
	c.LineChart = c.LineChart.WithValueRange(min, max)
	return c
}

func (c LineStackedChart) Position(x int64, y int64) LineStackedChart {
	c.LineChart = c.LineChart.Position(x, y)
	return c
}

func (c LineStackedChart) Size(cx int64, cy int64) LineStackedChart {
	c.LineChart = c.LineChart.Size(cx, cy)
	return c
}

func (c LineStackedChart) WithTitle(title string) LineStackedChart {
	c.LineChart = c.LineChart.WithTitle(title)
	return c
}

func (c LineStackedChart) WithLineColor(color string) LineStackedChart {
	c.LineChart = c.LineChart.WithLineColor(color)
	return c
}

func (c LineStackedChart) WithSeriesName(name string) LineStackedChart {
	c.LineChart = c.LineChart.WithSeriesName(name)
	return c
}

func (c LineStackedChart) WithLegend(show bool) LineStackedChart {
	c.LineChart = c.LineChart.WithLegend(show)
	return c
}

func (c LineStackedChart) WithLegendPosition(position string) LineStackedChart {
	c.LineChart = c.LineChart.WithLegendPosition(position)
	return c
}

func (c LineStackedChart) WithTitleOverlay(overlay bool) LineStackedChart {
	c.LineChart = c.LineChart.WithTitleOverlay(overlay)
	return c
}

func (c LineStackedChart) WithLegendOverlay(overlay bool) LineStackedChart {
	c.LineChart = c.LineChart.WithLegendOverlay(overlay)
	return c
}

func (c LineStackedChart) WithDataLabels(show bool) LineStackedChart {
	c.LineChart = c.LineChart.WithDataLabels(show)
	return c
}

func (c LineStackedChart) WithAxisTitles(categoryAxisTitle string, valueAxisTitle string) LineStackedChart {
	c.LineChart = c.LineChart.WithAxisTitles(categoryAxisTitle, valueAxisTitle)
	return c
}

func (c LineStackedChart) WithMajorGridlines(show bool) LineStackedChart {
	c.LineChart = c.LineChart.WithMajorGridlines(show)
	return c
}

func (c LineStackedChart) WithValueFormat(format string) LineStackedChart {
	c.LineChart = c.LineChart.WithValueFormat(format)
	return c
}

func (c LineStackedChart) WithValueAxisCrossBetween(mode string) LineStackedChart {
	c.LineChart = c.LineChart.WithValueAxisCrossBetween(mode)
	return c
}

func (c LineStackedChart) WithValueRange(min float64, max float64) LineStackedChart {
	c.LineChart = c.LineChart.WithValueRange(min, max)
	return c
}
