package pptx

import "github.com/djinn09/goppt/internal/pptxxml"

func slideScatterChartSpec(slide SlideContent) (*pptxxml.ChartSpec, bool) {
	if slide.Scatter == nil {
		return nil, false
	}
	chart := slide.Scatter
	return &pptxxml.ChartSpec{
		Kind:               pptxxml.ChartKindScatter,
		Title:              chart.Title,
		XValues:            copyFloat64Slice(chart.XValues),
		Values:             copyFloat64Slice(chart.YValues),
		X:                  chart.X,
		Y:                  chart.Y,
		CX:                 chart.CX,
		CY:                 chart.CY,
		Color:              normalizeHexColor(chart.LineColor),
		SeriesName:         chart.SeriesName,
		ScatterStyle:       chart.ScatterStyle,
		ShowLegend:         chart.ShowLegend,
		LegendPosition:     chart.LegendPosition,
		ShowDataLabels:     chart.ShowDataLabels,
		ShowMajorGridlines: chart.ShowMajorGridlines,
		CategoryAxisTitle:  chart.CategoryAxisTitle,
		ValueAxisTitle:     chart.ValueAxisTitle,
		ValueFormat:        chart.ValueFormat,
		MinValue:           copyFloat64Pointer(chart.MinValue),
		MaxValue:           copyFloat64Pointer(chart.MaxValue),
	}, true
}

func slidePieChartSpec(slide SlideContent) (*pptxxml.ChartSpec, bool) {
	if slide.Pie == nil {
		return nil, false
	}
	chart := slide.Pie
	return &pptxxml.ChartSpec{
		Kind:           pptxxml.ChartKindPie,
		Title:          chart.Title,
		Categories:     copyStringSlice(chart.Categories),
		Values:         copyFloat64Slice(chart.Values),
		X:              chart.X,
		Y:              chart.Y,
		CX:             chart.CX,
		CY:             chart.CY,
		SeriesName:     chart.SeriesName,
		ShowLegend:     chart.ShowLegend,
		LegendPosition: chart.LegendPosition,
		ShowDataLabels: chart.ShowDataLabels,
	}, true
}

func slideDoughnutChartSpec(slide SlideContent) (*pptxxml.ChartSpec, bool) {
	if slide.Dough == nil {
		return nil, false
	}
	chart := slide.Dough
	return &pptxxml.ChartSpec{
		Kind:           pptxxml.ChartKindDoughnut,
		Title:          chart.Title,
		Categories:     copyStringSlice(chart.Categories),
		Values:         copyFloat64Slice(chart.Values),
		X:              chart.X,
		Y:              chart.Y,
		CX:             chart.CX,
		CY:             chart.CY,
		SeriesName:     chart.SeriesName,
		ShowLegend:     chart.ShowLegend,
		LegendPosition: chart.LegendPosition,
		ShowDataLabels: chart.ShowDataLabels,
		HoleSize:       chart.HoleSize,
	}, true
}

func slideBubbleChartSpec(slide SlideContent) (*pptxxml.ChartSpec, bool) {
	if slide.Bubble == nil {
		return nil, false
	}
	chart := slide.Bubble
	return &pptxxml.ChartSpec{
		Kind:               pptxxml.ChartKindBubble,
		Title:              chart.Title,
		XValues:            copyFloat64Slice(chart.XValues),
		Values:             copyFloat64Slice(chart.YValues),
		BubbleSizes:        copyFloat64Slice(chart.BubbleSizes),
		X:                  chart.X,
		Y:                  chart.Y,
		CX:                 chart.CX,
		CY:                 chart.CY,
		Color:              normalizeHexColor(chart.LineColor),
		SeriesName:         chart.SeriesName,
		ShowLegend:         chart.ShowLegend,
		LegendPosition:     chart.LegendPosition,
		ShowDataLabels:     chart.ShowDataLabels,
		ShowMajorGridlines: chart.ShowMajorGridlines,
		CategoryAxisTitle:  chart.CategoryAxisTitle,
		ValueAxisTitle:     chart.ValueAxisTitle,
		ValueFormat:        chart.ValueFormat,
		MinValue:           copyFloat64Pointer(chart.MinValue),
		MaxValue:           copyFloat64Pointer(chart.MaxValue),
		BubbleScale:        chart.BubbleScale,
	}, true
}
