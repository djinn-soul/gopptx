package pptx

import "strings"

func (c StockHLCChart) Position(x int64, y int64) StockHLCChart {
	c.X = x
	c.Y = y
	return c
}

func (c StockHLCChart) Size(cx int64, cy int64) StockHLCChart {
	c.CX = cx
	c.CY = cy
	return c
}

func (c StockHLCChart) WithTitle(title string) StockHLCChart {
	c.Title = title
	return c
}

func (c StockHLCChart) WithLegend(show bool) StockHLCChart {
	c.ShowLegend = show
	return c
}

func (c StockHLCChart) WithLegendPosition(position string) StockHLCChart {
	c.LegendPosition = strings.ToLower(strings.TrimSpace(position))
	return c
}

func (c StockHLCChart) WithTitleOverlay(overlay bool) StockHLCChart {
	c.TitleOverlay = overlay
	return c
}

func (c StockHLCChart) WithLegendOverlay(overlay bool) StockHLCChart {
	c.LegendOverlay = overlay
	return c
}

func (c StockHLCChart) WithDataLabels(show bool) StockHLCChart {
	c.ShowDataLabels = show
	return c
}

func (c StockHLCChart) WithAxisTitles(categoryAxisTitle string, valueAxisTitle string) StockHLCChart {
	c.CategoryAxisTitle = strings.TrimSpace(categoryAxisTitle)
	c.ValueAxisTitle = strings.TrimSpace(valueAxisTitle)
	return c
}

func (c StockHLCChart) WithMajorGridlines(show bool) StockHLCChart {
	c.ShowMajorGridlines = show
	return c
}

func (c StockHLCChart) WithValueFormat(format string) StockHLCChart {
	c.ValueFormat = strings.TrimSpace(format)
	return c
}

func (c StockHLCChart) WithValueAxisCrossBetween(mode string) StockHLCChart {
	c.ValueAxisCrossBetween = strings.TrimSpace(mode)
	return c
}

func (c StockHLCChart) WithValueRange(min float64, max float64) StockHLCChart {
	c.MinValue = &min
	c.MaxValue = &max
	return c
}

func (c StockOHLCChart) Position(x int64, y int64) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.Position(x, y)
	return c
}

func (c StockOHLCChart) Size(cx int64, cy int64) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.Size(cx, cy)
	return c
}

func (c StockOHLCChart) WithTitle(title string) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithTitle(title)
	return c
}

func (c StockOHLCChart) WithLegend(show bool) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithLegend(show)
	return c
}

func (c StockOHLCChart) WithLegendPosition(position string) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithLegendPosition(position)
	return c
}

func (c StockOHLCChart) WithTitleOverlay(overlay bool) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithTitleOverlay(overlay)
	return c
}

func (c StockOHLCChart) WithLegendOverlay(overlay bool) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithLegendOverlay(overlay)
	return c
}

func (c StockOHLCChart) WithDataLabels(show bool) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithDataLabels(show)
	return c
}

func (c StockOHLCChart) WithAxisTitles(categoryAxisTitle string, valueAxisTitle string) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithAxisTitles(categoryAxisTitle, valueAxisTitle)
	return c
}

func (c StockOHLCChart) WithMajorGridlines(show bool) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithMajorGridlines(show)
	return c
}

func (c StockOHLCChart) WithValueFormat(format string) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithValueFormat(format)
	return c
}

func (c StockOHLCChart) WithValueAxisCrossBetween(mode string) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithValueAxisCrossBetween(mode)
	return c
}

func (c StockOHLCChart) WithValueRange(min float64, max float64) StockOHLCChart {
	c.StockHLCChart = c.StockHLCChart.WithValueRange(min, max)
	return c
}

func (c ComboChart) Position(x int64, y int64) ComboChart {
	c.X = x
	c.Y = y
	return c
}

func (c ComboChart) Size(cx int64, cy int64) ComboChart {
	c.CX = cx
	c.CY = cy
	return c
}

func (c ComboChart) WithTitle(title string) ComboChart {
	c.Title = title
	return c
}

func (c ComboChart) WithLegend(show bool) ComboChart {
	c.ShowLegend = show
	return c
}

func (c ComboChart) WithLegendPosition(position string) ComboChart {
	c.LegendPosition = strings.ToLower(strings.TrimSpace(position))
	return c
}

func (c ComboChart) WithTitleOverlay(overlay bool) ComboChart {
	c.TitleOverlay = overlay
	return c
}

func (c ComboChart) WithLegendOverlay(overlay bool) ComboChart {
	c.LegendOverlay = overlay
	return c
}

func (c ComboChart) WithDataLabels(show bool) ComboChart {
	c.ShowDataLabels = show
	return c
}

func (c ComboChart) WithAxisTitles(categoryAxisTitle string, valueAxisTitle string) ComboChart {
	c.CategoryAxisTitle = strings.TrimSpace(categoryAxisTitle)
	c.ValueAxisTitle = strings.TrimSpace(valueAxisTitle)
	return c
}

func (c ComboChart) WithMajorGridlines(show bool) ComboChart {
	c.ShowMajorGridlines = show
	return c
}

func (c ComboChart) WithValueFormat(format string) ComboChart {
	c.ValueFormat = strings.TrimSpace(format)
	return c
}

func (c ComboChart) WithValueAxisCrossBetween(mode string) ComboChart {
	c.ValueAxisCrossBetween = strings.TrimSpace(mode)
	return c
}

func (c ComboChart) WithValueRange(min float64, max float64) ComboChart {
	c.MinValue = &min
	c.MaxValue = &max
	return c
}
