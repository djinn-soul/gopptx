package pptx

import "github.com/djinn09/goppt/internal/pptxxml"

func slideAreaChartSpec(slide SlideContent) (*pptxxml.ChartSpec, bool) {
	if slide.Area == nil {
		return nil, false
	}
	chart := slide.Area
	return &pptxxml.ChartSpec{
		Kind:                  pptxxml.ChartKindArea,
		Title:                 chart.Title,
		TitleOverlay:          chart.TitleOverlay,
		Categories:            copyStringSlice(chart.Categories),
		Values:                copyFloat64Slice(chart.Values),
		X:                     chart.X,
		Y:                     chart.Y,
		CX:                    chart.CX,
		CY:                    chart.CY,
		Color:                 normalizeHexColor(chart.AreaColor),
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
		Grouping:              "standard",
	}, true
}

func slideAreaStackedChartSpec(slide SlideContent) (*pptxxml.ChartSpec, bool) {
	if slide.AreaStacked == nil {
		return nil, false
	}
	chart := slide.AreaStacked
	return &pptxxml.ChartSpec{
		Kind:                  pptxxml.ChartKindAreaStacked,
		Title:                 chart.Title,
		TitleOverlay:          chart.TitleOverlay,
		Categories:            copyStringSlice(chart.Categories),
		Values:                copyFloat64Slice(chart.Values),
		X:                     chart.X,
		Y:                     chart.Y,
		CX:                    chart.CX,
		CY:                    chart.CY,
		Color:                 normalizeHexColor(chart.AreaColor),
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
		Grouping:              "stacked",
	}, true
}

func slideAreaStacked100ChartSpec(slide SlideContent) (*pptxxml.ChartSpec, bool) {
	if slide.AreaStacked100 == nil {
		return nil, false
	}
	chart := slide.AreaStacked100
	return &pptxxml.ChartSpec{
		Kind:                  pptxxml.ChartKindAreaStacked100,
		Title:                 chart.Title,
		TitleOverlay:          chart.TitleOverlay,
		Categories:            copyStringSlice(chart.Categories),
		Values:                copyFloat64Slice(chart.Values),
		X:                     chart.X,
		Y:                     chart.Y,
		CX:                    chart.CX,
		CY:                    chart.CY,
		Color:                 normalizeHexColor(chart.AreaColor),
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
		Grouping:              "percentStacked",
	}, true
}
