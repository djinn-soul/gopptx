package export

import (
	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// renderBarAndLineCharts renders bar, line, and area chart types from a slide.
func renderBarAndLineCharts(pdf *gopdf.GoPdf, slide elements.SlideContent) {
	if slide.Chart != nil {
		c := slide.Chart
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
	if slide.BarHorizontal != nil {
		c := slide.BarHorizontal
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
	if slide.BarStacked != nil {
		c := slide.BarStacked
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
	if slide.BarStacked100 != nil {
		c := slide.BarStacked100
		renderBarLike(pdf, c.Title,
			chartRectFromLength(c.X.Emu(), c.Y.Emu(), c.CX.Emu(), c.CY.Emu()),
			fullBars(len(c.Values)), c.Categories, false,
			chartSeriesOpts{
				color: c.BarColor,
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
	if slide.Line != nil {
		c := slide.Line
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
	if slide.LineMarkers != nil {
		c := slide.LineMarkers
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
	if slide.LineStacked != nil {
		c := slide.LineStacked
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
	if slide.Area != nil {
		c := slide.Area
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
	if slide.AreaStacked != nil {
		c := slide.AreaStacked
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
	if slide.AreaStacked100 != nil {
		c := slide.AreaStacked100
		renderAreaLike(pdf, c.Title,
			chartRectFromLength(c.X.Emu(), c.Y.Emu(), c.CX.Emu(), c.CY.Emu()),
			fullBars(len(c.Values)), c.Categories,
			chartSeriesOpts{
				color: c.AreaColor,
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
}

// renderOtherCharts renders pie, scatter, radar, stock, and combo chart types from a slide.
func renderOtherCharts(pdf *gopdf.GoPdf, slide elements.SlideContent) {
	if slide.Pie != nil {
		c := slide.Pie
		renderPieLike(pdf, c.Title,
			chartRectFromLength(c.X.Emu(), c.Y.Emu(), c.CX.Emu(), c.CY.Emu()),
			c.Values, 0, c.Categories,
			chartSeriesOpts{
				showDataLabels: c.ShowDataLabels,
				showLegend:     c.ShowLegend, legendPosition: c.LegendPosition, seriesName: c.SeriesName,
				titleOverlay: c.TitleOverlay,
			},
		)
	}
	if slide.Doughnut != nil {
		c := slide.Doughnut
		renderPieLike(pdf, c.Title,
			chartRectFromLength(c.X.Emu(), c.Y.Emu(), c.CX.Emu(), c.CY.Emu()),
			c.Values, c.HoleSize, c.Categories,
			chartSeriesOpts{
				showDataLabels: c.ShowDataLabels,
				showLegend:     c.ShowLegend, legendPosition: c.LegendPosition, seriesName: c.SeriesName,
				titleOverlay: c.TitleOverlay,
			},
		)
	}
	if slide.Scatter != nil {
		c := slide.Scatter
		renderScatterLike(pdf, c.Title,
			chartRectFromLength(c.X.Emu(), c.Y.Emu(), c.CX.Emu(), c.CY.Emu()),
			c.XValues, c.YValues, nil,
			chartSeriesOpts{
				color: c.LineColor, scatterStyle: c.ScatterStyle,
				showLegend: c.ShowLegend, legendPosition: c.LegendPosition, seriesName: c.SeriesName,
				showDataLabels: c.ShowDataLabels,
				catAxisTitle:   c.CategoryAxisTitle, valAxisTitle: c.ValueAxisTitle,
				showMajorGridlines: c.ShowMajorGridlines,
				titleOverlay:       c.TitleOverlay,
				valueFormat:        c.ValueFormat,
			},
		)
	}
	if slide.Bubble != nil {
		c := slide.Bubble
		renderScatterLike(pdf, c.Title,
			chartRectFromLength(c.X, c.Y, c.CX, c.CY),
			c.XValues, c.YValues, c.BubbleSizes,
			chartSeriesOpts{
				color:      c.LineColor,
				showLegend: c.ShowLegend, legendPosition: c.LegendPosition, seriesName: c.SeriesName,
				catAxisTitle: c.CategoryAxisTitle, valAxisTitle: c.ValueAxisTitle,
				showMajorGridlines: c.ShowMajorGridlines,
				titleOverlay:       c.TitleOverlay,
				valueFormat:        c.ValueFormat,
				bubbleScale:        c.BubbleScale,
			},
		)
	}
	if slide.Radar != nil {
		c := slide.Radar
		renderRadarLike(pdf, c.Title,
			chartRectFromLength(c.X.Emu(), c.Y.Emu(), c.CX.Emu(), c.CY.Emu()),
			c.Values, c.Categories, false,
			chartSeriesOpts{
				color:      c.LineColor,
				showLegend: c.ShowLegend, legendPosition: c.LegendPosition, seriesName: c.SeriesName,
				titleOverlay: c.TitleOverlay,
			},
		)
	}
	if slide.RadarFilled != nil {
		c := slide.RadarFilled
		renderRadarLike(pdf, c.Title,
			chartRectFromLength(c.X.Emu(), c.Y.Emu(), c.CX.Emu(), c.CY.Emu()),
			c.Values, c.Categories, true,
			chartSeriesOpts{
				color:      c.LineColor,
				showLegend: c.ShowLegend, legendPosition: c.LegendPosition, seriesName: c.SeriesName,
				titleOverlay: c.TitleOverlay,
			},
		)
	}
	if slide.StockHLC != nil {
		c := slide.StockHLC
		renderStockLike(pdf, c.Title,
			chartRectFromLength(c.X.Emu(), c.Y.Emu(), c.CX.Emu(), c.CY.Emu()),
			nil, c.HighValues, c.LowValues, c.CloseValues,
			c.Categories,
			chartSeriesOpts{
				showMajorGridlines: c.ShowMajorGridlines,
				catAxisTitle:       c.CategoryAxisTitle, valAxisTitle: c.ValueAxisTitle,
				showLegend: c.ShowLegend, legendPosition: c.LegendPosition,
				titleOverlay: c.TitleOverlay,
				valueFormat:  c.ValueFormat,
				minValue:     c.MinValue, maxValue: c.MaxValue,
			},
		)
	}
	if slide.StockOHLC != nil {
		c := slide.StockOHLC
		renderStockLike(pdf, c.Title,
			chartRectFromLength(c.X.Emu(), c.Y.Emu(), c.CX.Emu(), c.CY.Emu()),
			c.OpenValues, c.HighValues, c.LowValues, c.CloseValues,
			c.Categories,
			chartSeriesOpts{
				showMajorGridlines: c.ShowMajorGridlines,
				catAxisTitle:       c.CategoryAxisTitle, valAxisTitle: c.ValueAxisTitle,
				showLegend: c.ShowLegend, legendPosition: c.LegendPosition,
				titleOverlay: c.TitleOverlay,
				valueFormat:  c.ValueFormat,
				minValue:     c.MinValue, maxValue: c.MaxValue,
			},
		)
	}
	if slide.Combo != nil {
		c := slide.Combo
		renderComboLike(pdf, c.Title,
			chartRectFromLength(c.X.Emu(), c.Y.Emu(), c.CX.Emu(), c.CY.Emu()),
			c.BarSeries, c.LineSeries, c.Categories,
			chartSeriesOpts{
				showLegend: c.ShowLegend, legendPosition: c.LegendPosition,
				showMajorGridlines: c.ShowMajorGridlines,
				catAxisTitle:       c.CategoryAxisTitle, valAxisTitle: c.ValueAxisTitle,
				titleOverlay: c.TitleOverlay,
				valueFormat:  c.ValueFormat,
				minValue:     c.MinValue, maxValue: c.MaxValue,
			},
		)
	}
}
