package charts

import "strings"

// WithSeriesName sets the single-series label in chart legends.
func (c BarChart) WithSeriesName(name string) BarChart {
	c.SeriesName = strings.TrimSpace(name)
	return c
}

// WithLegend toggles legend visibility.
func (c BarChart) WithLegend(show bool) BarChart {
	c.ShowLegend = show
	return c
}

// WithLegendPosition sets legend position as r/l/t/b.
func (c BarChart) WithLegendPosition(position string) BarChart {
	c.LegendPosition = strings.ToLower(strings.TrimSpace(position))
	return c
}

// WithTitleOverlay toggles title overlay on chart plot area.
func (c BarChart) WithTitleOverlay(overlay bool) BarChart {
	c.TitleOverlay = overlay
	return c
}

// WithLegendOverlay toggles legend overlay on chart plot area.
func (c BarChart) WithLegendOverlay(overlay bool) BarChart {
	c.LegendOverlay = overlay
	return c
}

// WithDataLabels toggles value labels on chart points.
func (c BarChart) WithDataLabels(show bool) BarChart {
	c.ShowDataLabels = show
	return c
}

// WithAxisTitles sets optional category/value axis titles.
func (c BarChart) WithAxisTitles(categoryAxisTitle string, valueAxisTitle string) BarChart {
	c.CategoryAxisTitle = strings.TrimSpace(categoryAxisTitle)
	c.ValueAxisTitle = strings.TrimSpace(valueAxisTitle)
	return c
}

// WithMajorGridlines toggles value-axis major gridlines.
func (c BarChart) WithMajorGridlines(show bool) BarChart {
	c.ShowMajorGridlines = show
	return c
}

// WithValueFormat sets the value-axis number format code.
func (c BarChart) WithValueFormat(format string) BarChart {
	c.ValueFormat = strings.TrimSpace(format)
	return c
}

// WithValueAxisCrossBetween sets value-axis crossing mode: between|midCat.
func (c BarChart) WithValueAxisCrossBetween(mode string) BarChart {
	c.ValueAxisCrossBetween = strings.TrimSpace(mode)
	return c
}

// WithValueRange sets optional value-axis min/max.
func (c BarChart) WithValueRange(min float64, max float64) BarChart {
	c.MinValue = &min
	c.MaxValue = &max
	return c
}

// WithSeriesName sets the single-series label in chart legends.
func (c LineChart) WithSeriesName(name string) LineChart {
	c.SeriesName = strings.TrimSpace(name)
	return c
}

// WithLegend toggles legend visibility.
func (c LineChart) WithLegend(show bool) LineChart {
	c.ShowLegend = show
	return c
}

// WithLegendPosition sets legend position as r/l/t/b.
func (c LineChart) WithLegendPosition(position string) LineChart {
	c.LegendPosition = strings.ToLower(strings.TrimSpace(position))
	return c
}

// WithTitleOverlay toggles title overlay on chart plot area.
func (c LineChart) WithTitleOverlay(overlay bool) LineChart {
	c.TitleOverlay = overlay
	return c
}

// WithLegendOverlay toggles legend overlay on chart plot area.
func (c LineChart) WithLegendOverlay(overlay bool) LineChart {
	c.LegendOverlay = overlay
	return c
}

// WithDataLabels toggles value labels on chart points.
func (c LineChart) WithDataLabels(show bool) LineChart {
	c.ShowDataLabels = show
	return c
}

// WithAxisTitles sets optional category/value axis titles.
func (c LineChart) WithAxisTitles(categoryAxisTitle string, valueAxisTitle string) LineChart {
	c.CategoryAxisTitle = strings.TrimSpace(categoryAxisTitle)
	c.ValueAxisTitle = strings.TrimSpace(valueAxisTitle)
	return c
}

// WithMajorGridlines toggles value-axis major gridlines.
func (c LineChart) WithMajorGridlines(show bool) LineChart {
	c.ShowMajorGridlines = show
	return c
}

// WithValueFormat sets the value-axis number format code.
func (c LineChart) WithValueFormat(format string) LineChart {
	c.ValueFormat = strings.TrimSpace(format)
	return c
}

// WithValueAxisCrossBetween sets value-axis crossing mode: between|midCat.
func (c LineChart) WithValueAxisCrossBetween(mode string) LineChart {
	c.ValueAxisCrossBetween = strings.TrimSpace(mode)
	return c
}

// WithValueRange sets optional value-axis min/max.
func (c LineChart) WithValueRange(min float64, max float64) LineChart {
	c.MinValue = &min
	c.MaxValue = &max
	return c
}

// WithSmooth toggles line smoothing.
func (c LineChart) WithSmooth(smooth bool) LineChart {
	c.Smooth = smooth
	return c
}

// WithSeriesName sets the single-series label in chart legends.
func (c PieChart) WithSeriesName(name string) PieChart {
	c.SeriesName = strings.TrimSpace(name)
	return c
}

// WithLegend toggles legend visibility.
func (c PieChart) WithLegend(show bool) PieChart {
	c.ShowLegend = show
	return c
}

// WithLegendPosition sets legend position as r/l/t/b.
func (c PieChart) WithLegendPosition(position string) PieChart {
	c.LegendPosition = strings.ToLower(strings.TrimSpace(position))
	return c
}

// WithTitleOverlay toggles title overlay on chart plot area.
func (c PieChart) WithTitleOverlay(overlay bool) PieChart {
	c.TitleOverlay = overlay
	return c
}

// WithLegendOverlay toggles legend overlay on chart plot area.
func (c PieChart) WithLegendOverlay(overlay bool) PieChart {
	c.LegendOverlay = overlay
	return c
}

// WithDataLabels toggles data labels on pie slices.
func (c PieChart) WithDataLabels(show bool) PieChart {
	c.ShowDataLabels = show
	return c
}

// WithSeriesName sets the single-series label in chart legends.
func (c DoughnutChart) WithSeriesName(name string) DoughnutChart {
	c.SeriesName = strings.TrimSpace(name)
	return c
}

// WithLegend toggles legend visibility.
func (c DoughnutChart) WithLegend(show bool) DoughnutChart {
	c.ShowLegend = show
	return c
}

// WithLegendPosition sets legend position as r/l/t/b.
func (c DoughnutChart) WithLegendPosition(position string) DoughnutChart {
	c.LegendPosition = strings.ToLower(strings.TrimSpace(position))
	return c
}

// WithTitleOverlay toggles title overlay on chart plot area.
func (c DoughnutChart) WithTitleOverlay(overlay bool) DoughnutChart {
	c.TitleOverlay = overlay
	return c
}

// WithLegendOverlay toggles legend overlay on chart plot area.
func (c DoughnutChart) WithLegendOverlay(overlay bool) DoughnutChart {
	c.LegendOverlay = overlay
	return c
}

// WithDataLabels toggles data labels on doughnut slices.
func (c DoughnutChart) WithDataLabels(show bool) DoughnutChart {
	c.ShowDataLabels = show
	return c
}

// WithHoleSize sets the center hole size percentage (10-90).
func (c DoughnutChart) WithHoleSize(size int) DoughnutChart {
	c.HoleSize = size
	return c
}
