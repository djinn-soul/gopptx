package export

import (
	"github.com/signintech/gopdf"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

func renderPieChart(pdf *gopdf.GoPdf, c *charts.PieChart) {
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

func renderDoughnutChart(pdf *gopdf.GoPdf, c *charts.DoughnutChart) {
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

func renderScatterChart(pdf *gopdf.GoPdf, c *charts.ScatterChart) {
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

func renderBubbleChart(pdf *gopdf.GoPdf, c *charts.BubbleChart) {
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

func renderRadarChart(pdf *gopdf.GoPdf, c *charts.RadarChart) {
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

func renderRadarFilledChart(pdf *gopdf.GoPdf, c *charts.RadarFilledChart) {
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

func renderStockHLCChart(pdf *gopdf.GoPdf, c *charts.StockHLCChart) {
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

func renderStockOHLCChart(pdf *gopdf.GoPdf, c *charts.StockOHLCChart) {
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

func renderComboChart(pdf *gopdf.GoPdf, c *charts.ComboChart) {
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
