package pptx

import "github.com/djinn09/goppt/internal/pptxxml"

func slideRadarChartSpec(slide SlideContent) (*pptxxml.ChartSpec, bool) {
	if slide.Radar == nil {
		return nil, false
	}
	chart := slide.Radar
	return &pptxxml.ChartSpec{
		Kind:                  pptxxml.ChartKindRadar,
		Title:                 chart.Title,
		TitleOverlay:          chart.TitleOverlay,
		Categories:            copyStringSlice(chart.Categories),
		Values:                copyFloat64Slice(chart.Values),
		X:                     chart.X,
		Y:                     chart.Y,
		CX:                    chart.CX,
		CY:                    chart.CY,
		Color:                 normalizeHexColor(chart.LineColor),
		SeriesName:            chart.SeriesName,
		ShowLegend:            chart.ShowLegend,
		LegendPosition:        chart.LegendPosition,
		LegendOverlay:         chart.LegendOverlay,
		ShowDataLabels:        chart.ShowDataLabels,
		ShowMajorGridlines:    chart.ShowMajorGridlines,
		CategoryAxisTitle:     chart.CategoryAxisTitle,
		ValueAxisTitle:        chart.ValueAxisTitle,
		ValueFormat:           chart.ValueFormat,
		ValueAxisCrossBetween: chart.ValueAxisCrossBetween,
		MinValue:              copyFloat64Pointer(chart.MinValue),
		MaxValue:              copyFloat64Pointer(chart.MaxValue),
		RadarStyle:            chart.RadarStyle,
	}, true
}

func slideRadarFilledChartSpec(slide SlideContent) (*pptxxml.ChartSpec, bool) {
	if slide.RadarFilled == nil {
		return nil, false
	}
	chart := slide.RadarFilled
	return &pptxxml.ChartSpec{
		Kind:                  pptxxml.ChartKindRadarFilled,
		Title:                 chart.Title,
		TitleOverlay:          chart.TitleOverlay,
		Categories:            copyStringSlice(chart.Categories),
		Values:                copyFloat64Slice(chart.Values),
		X:                     chart.X,
		Y:                     chart.Y,
		CX:                    chart.CX,
		CY:                    chart.CY,
		Color:                 normalizeHexColor(chart.LineColor),
		SeriesName:            chart.SeriesName,
		ShowLegend:            chart.ShowLegend,
		LegendPosition:        chart.LegendPosition,
		LegendOverlay:         chart.LegendOverlay,
		ShowDataLabels:        chart.ShowDataLabels,
		ShowMajorGridlines:    chart.ShowMajorGridlines,
		CategoryAxisTitle:     chart.CategoryAxisTitle,
		ValueAxisTitle:        chart.ValueAxisTitle,
		ValueFormat:           chart.ValueFormat,
		ValueAxisCrossBetween: chart.ValueAxisCrossBetween,
		MinValue:              copyFloat64Pointer(chart.MinValue),
		MaxValue:              copyFloat64Pointer(chart.MaxValue),
		RadarStyle:            RadarStyleFilled,
	}, true
}

func slideStockHLCChartSpec(slide SlideContent) (*pptxxml.ChartSpec, bool) {
	if slide.StockHLC == nil {
		return nil, false
	}
	chart := slide.StockHLC
	return &pptxxml.ChartSpec{
		Kind:                  pptxxml.ChartKindStockHLC,
		Title:                 chart.Title,
		TitleOverlay:          chart.TitleOverlay,
		Categories:            copyStringSlice(chart.Categories),
		HighValues:            copyFloat64Slice(chart.HighValues),
		LowValues:             copyFloat64Slice(chart.LowValues),
		CloseValues:           copyFloat64Slice(chart.CloseValues),
		X:                     chart.X,
		Y:                     chart.Y,
		CX:                    chart.CX,
		CY:                    chart.CY,
		ShowLegend:            chart.ShowLegend,
		LegendPosition:        chart.LegendPosition,
		LegendOverlay:         chart.LegendOverlay,
		ShowDataLabels:        chart.ShowDataLabels,
		ShowMajorGridlines:    chart.ShowMajorGridlines,
		CategoryAxisTitle:     chart.CategoryAxisTitle,
		ValueAxisTitle:        chart.ValueAxisTitle,
		ValueFormat:           chart.ValueFormat,
		ValueAxisCrossBetween: chart.ValueAxisCrossBetween,
		MinValue:              copyFloat64Pointer(chart.MinValue),
		MaxValue:              copyFloat64Pointer(chart.MaxValue),
	}, true
}

func slideStockOHLCChartSpec(slide SlideContent) (*pptxxml.ChartSpec, bool) {
	if slide.StockOHLC == nil {
		return nil, false
	}
	chart := slide.StockOHLC
	return &pptxxml.ChartSpec{
		Kind:                  pptxxml.ChartKindStockOHLC,
		Title:                 chart.Title,
		TitleOverlay:          chart.TitleOverlay,
		Categories:            copyStringSlice(chart.Categories),
		OpenValues:            copyFloat64Slice(chart.OpenValues),
		HighValues:            copyFloat64Slice(chart.HighValues),
		LowValues:             copyFloat64Slice(chart.LowValues),
		CloseValues:           copyFloat64Slice(chart.CloseValues),
		X:                     chart.X,
		Y:                     chart.Y,
		CX:                    chart.CX,
		CY:                    chart.CY,
		ShowLegend:            chart.ShowLegend,
		LegendPosition:        chart.LegendPosition,
		LegendOverlay:         chart.LegendOverlay,
		ShowDataLabels:        chart.ShowDataLabels,
		ShowMajorGridlines:    chart.ShowMajorGridlines,
		CategoryAxisTitle:     chart.CategoryAxisTitle,
		ValueAxisTitle:        chart.ValueAxisTitle,
		ValueFormat:           chart.ValueFormat,
		ValueAxisCrossBetween: chart.ValueAxisCrossBetween,
		MinValue:              copyFloat64Pointer(chart.MinValue),
		MaxValue:              copyFloat64Pointer(chart.MaxValue),
	}, true
}

func slideComboChartSpec(slide SlideContent) (*pptxxml.ChartSpec, bool) {
	if slide.Combo == nil {
		return nil, false
	}
	chart := slide.Combo
	return &pptxxml.ChartSpec{
		Kind:                  pptxxml.ChartKindCombo,
		Title:                 chart.Title,
		TitleOverlay:          chart.TitleOverlay,
		Categories:            copyStringSlice(chart.Categories),
		BarSeries:             toXMLSeries(chart.BarSeries),
		LineSeries:            toXMLSeries(chart.LineSeries),
		X:                     chart.X,
		Y:                     chart.Y,
		CX:                    chart.CX,
		CY:                    chart.CY,
		ShowLegend:            chart.ShowLegend,
		LegendPosition:        chart.LegendPosition,
		LegendOverlay:         chart.LegendOverlay,
		ShowDataLabels:        chart.ShowDataLabels,
		ShowMajorGridlines:    chart.ShowMajorGridlines,
		CategoryAxisTitle:     chart.CategoryAxisTitle,
		ValueAxisTitle:        chart.ValueAxisTitle,
		ValueFormat:           chart.ValueFormat,
		ValueAxisCrossBetween: chart.ValueAxisCrossBetween,
		MinValue:              copyFloat64Pointer(chart.MinValue),
		MaxValue:              copyFloat64Pointer(chart.MaxValue),
	}, true
}
