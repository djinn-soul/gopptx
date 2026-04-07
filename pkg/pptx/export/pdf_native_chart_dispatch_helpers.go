package export

import (
	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

func renderBarChart(pdf *gopdf.GoPdf, c *charts.BarChart) {
	renderBarLike(pdf, c.Title,
		chartRectFromLength(c.X.Emu(), c.Y.Emu(), c.CX.Emu(), c.CY.Emu()),
		c.Values, c.Categories, false,
		chartSeriesOpts{
			color: c.BarColor, minValue: c.MinValue, maxValue: c.MaxValue,
			showLegend: c.ShowLegend, legendPosition: c.LegendPosition, seriesName: c.SeriesName,
			showDataLabels: c.ShowDataLabels,
			catAxisTitle:   c.CategoryAxisTitle, valAxisTitle: c.ValueAxisTitle,
			showMajorGridlines: c.ShowMajorGridlines,
			showCatGridlines:   c.ShowCategoryMajorGridlines,
			titleOverlay:       c.TitleOverlay,
			valueFormat:        c.ValueFormat,
		},
	)
}

func renderBarHorizontalChart(pdf *gopdf.GoPdf, c *charts.BarHorizontalChart) {
	renderBarLike(pdf, c.Title,
		chartRectFromLength(c.X.Emu(), c.Y.Emu(), c.CX.Emu(), c.CY.Emu()),
		c.Values, c.Categories, true,
		chartSeriesOpts{
			color: c.BarColor, minValue: c.MinValue, maxValue: c.MaxValue,
			showLegend: c.ShowLegend, legendPosition: c.LegendPosition, seriesName: c.SeriesName,
			showDataLabels: c.ShowDataLabels,
			catAxisTitle:   c.CategoryAxisTitle, valAxisTitle: c.ValueAxisTitle,
			showMajorGridlines: c.ShowMajorGridlines,
			showCatGridlines:   c.ShowCategoryMajorGridlines,
			titleOverlay:       c.TitleOverlay,
			valueFormat:        c.ValueFormat,
		},
	)
}

func renderBarStackedChart(pdf *gopdf.GoPdf, c *charts.BarStackedChart) {
	renderBarLike(pdf, c.Title,
		chartRectFromLength(c.X.Emu(), c.Y.Emu(), c.CX.Emu(), c.CY.Emu()),
		c.Values, c.Categories, true,
		chartSeriesOpts{
			color: c.BarColor, minValue: c.MinValue, maxValue: c.MaxValue,
			showLegend: c.ShowLegend, legendPosition: c.LegendPosition, seriesName: c.SeriesName,
			showDataLabels: c.ShowDataLabels,
			catAxisTitle:   c.CategoryAxisTitle, valAxisTitle: c.ValueAxisTitle,
			showMajorGridlines: c.ShowMajorGridlines,
			showCatGridlines:   c.ShowCategoryMajorGridlines,
			titleOverlay:       c.TitleOverlay,
			valueFormat:        c.ValueFormat,
		},
	)
}

func renderBarStacked100Chart(pdf *gopdf.GoPdf, c *charts.BarStacked100Chart) {
	renderBarLike(pdf, c.Title,
		chartRectFromLength(c.X.Emu(), c.Y.Emu(), c.CX.Emu(), c.CY.Emu()),
		fullBars(len(c.Values)), c.Categories, true,
		chartSeriesOpts{
			color:      c.BarColor,
			showLegend: c.ShowLegend, legendPosition: c.LegendPosition, seriesName: c.SeriesName,
			showDataLabels: c.ShowDataLabels,
			catAxisTitle:   c.CategoryAxisTitle, valAxisTitle: c.ValueAxisTitle,
			showMajorGridlines: c.ShowMajorGridlines,
			showCatGridlines:   c.ShowCategoryMajorGridlines,
			titleOverlay:       c.TitleOverlay,
			valueFormat:        "0%",
		},
	)
}

func renderLineChart(pdf *gopdf.GoPdf, c *charts.LineChart) {
	renderLineLike(pdf, c.Title,
		chartRectFromLength(c.X.Emu(), c.Y.Emu(), c.CX.Emu(), c.CY.Emu()),
		c.Values, c.Categories, false,
		chartSeriesOpts{
			color: c.LineColor, minValue: c.MinValue, maxValue: c.MaxValue,
			showLegend: c.ShowLegend, legendPosition: c.LegendPosition, seriesName: c.SeriesName,
			showDataLabels: c.ShowDataLabels,
			catAxisTitle:   c.CategoryAxisTitle, valAxisTitle: c.ValueAxisTitle,
			showMajorGridlines: c.ShowMajorGridlines,
			showCatGridlines:   c.ShowCategoryMajorGridlines,
			titleOverlay:       c.TitleOverlay,
			valueFormat:        c.ValueFormat,
			smooth:             c.Smooth,
		},
	)
}

func renderLineMarkersChart(pdf *gopdf.GoPdf, c *charts.LineMarkersChart) {
	renderLineLike(pdf, c.Title,
		chartRectFromLength(c.X.Emu(), c.Y.Emu(), c.CX.Emu(), c.CY.Emu()),
		c.Values, c.Categories, true,
		chartSeriesOpts{
			color: c.LineColor, minValue: c.MinValue, maxValue: c.MaxValue,
			showLegend: c.ShowLegend, legendPosition: c.LegendPosition, seriesName: c.SeriesName,
			showDataLabels: c.ShowDataLabels,
			catAxisTitle:   c.CategoryAxisTitle, valAxisTitle: c.ValueAxisTitle,
			showMajorGridlines: c.ShowMajorGridlines,
			showCatGridlines:   c.ShowCategoryMajorGridlines,
			titleOverlay:       c.TitleOverlay,
			valueFormat:        c.ValueFormat,
			smooth:             c.Smooth,
		},
	)
}

func renderLineStackedChart(pdf *gopdf.GoPdf, c *charts.LineStackedChart) {
	renderLineLike(pdf, c.Title,
		chartRectFromLength(c.X.Emu(), c.Y.Emu(), c.CX.Emu(), c.CY.Emu()),
		c.Values, c.Categories, false,
		chartSeriesOpts{
			color: c.LineColor, minValue: c.MinValue, maxValue: c.MaxValue,
			showLegend: c.ShowLegend, legendPosition: c.LegendPosition, seriesName: c.SeriesName,
			showDataLabels: c.ShowDataLabels,
			catAxisTitle:   c.CategoryAxisTitle, valAxisTitle: c.ValueAxisTitle,
			showMajorGridlines: c.ShowMajorGridlines,
			showCatGridlines:   c.ShowCategoryMajorGridlines,
			titleOverlay:       c.TitleOverlay,
			valueFormat:        c.ValueFormat,
			smooth:             c.Smooth,
		},
	)
}

func renderAreaChart(pdf *gopdf.GoPdf, c *charts.AreaChart) {
	renderAreaLike(pdf, c.Title,
		chartRectFromLength(c.X.Emu(), c.Y.Emu(), c.CX.Emu(), c.CY.Emu()),
		c.Values, c.Categories,
		chartSeriesOpts{
			color: c.AreaColor, minValue: c.MinValue, maxValue: c.MaxValue,
			showLegend: c.ShowLegend, legendPosition: c.LegendPosition, seriesName: c.SeriesName,
			showDataLabels: c.ShowDataLabels,
			catAxisTitle:   c.CategoryAxisTitle, valAxisTitle: c.ValueAxisTitle,
			showMajorGridlines: c.ShowMajorGridlines,
			showCatGridlines:   c.ShowCategoryMajorGridlines,
			titleOverlay:       c.TitleOverlay,
			valueFormat:        c.ValueFormat,
		},
	)
}

func renderAreaStackedChart(pdf *gopdf.GoPdf, c *charts.AreaStackedChart) {
	renderAreaLike(pdf, c.Title,
		chartRectFromLength(c.X.Emu(), c.Y.Emu(), c.CX.Emu(), c.CY.Emu()),
		c.Values, c.Categories,
		chartSeriesOpts{
			color: c.AreaColor, minValue: c.MinValue, maxValue: c.MaxValue,
			showLegend: c.ShowLegend, legendPosition: c.LegendPosition, seriesName: c.SeriesName,
			showDataLabels: c.ShowDataLabels,
			catAxisTitle:   c.CategoryAxisTitle, valAxisTitle: c.ValueAxisTitle,
			showMajorGridlines: c.ShowMajorGridlines,
			showCatGridlines:   c.ShowCategoryMajorGridlines,
			titleOverlay:       c.TitleOverlay,
			valueFormat:        c.ValueFormat,
		},
	)
}

func renderAreaStacked100Chart(pdf *gopdf.GoPdf, c *charts.AreaStacked100Chart) {
	renderAreaLike(pdf, c.Title,
		chartRectFromLength(c.X.Emu(), c.Y.Emu(), c.CX.Emu(), c.CY.Emu()),
		fullBars(len(c.Values)), c.Categories,
		chartSeriesOpts{
			color:      c.AreaColor,
			showLegend: c.ShowLegend, legendPosition: c.LegendPosition, seriesName: c.SeriesName,
			showDataLabels: c.ShowDataLabels,
			catAxisTitle:   c.CategoryAxisTitle, valAxisTitle: c.ValueAxisTitle,
			showMajorGridlines: c.ShowMajorGridlines,
			showCatGridlines:   c.ShowCategoryMajorGridlines,
			titleOverlay:       c.TitleOverlay,
			valueFormat:        "0%",
		},
	)
}
