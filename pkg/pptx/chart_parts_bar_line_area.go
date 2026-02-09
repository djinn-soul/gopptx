package pptx

import "github.com/djinn-soul/gopptx/internal/pptxxml"

func slideBarChartSpec(slide SlideContent) (*pptxxml.ChartSpec, bool) {
	if slide.Chart == nil {
		return nil, false
	}
	chart := slide.Chart
	return &pptxxml.ChartSpec{
		Kind:                  pptxxml.ChartKindBar,
		Title:                 chart.Title,
		TitleOverlay:          chart.TitleOverlay,
		Categories:            copyStringSlice(chart.Categories),
		Values:                copyFloat64Slice(chart.Values),
		X:                     chart.X,
		Y:                     chart.Y,
		CX:                    chart.CX,
		CY:                    chart.CY,
		Color:                 normalizeHexColor(chart.BarColor),
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
		BarDir:                "col",
		Grouping:              "clustered",
	}, true
}

func slideBarHorizontalChartSpec(slide SlideContent) (*pptxxml.ChartSpec, bool) {
	if slide.BarHorizontal == nil {
		return nil, false
	}
	chart := slide.BarHorizontal
	return &pptxxml.ChartSpec{
		Kind:                  pptxxml.ChartKindBarHorizontal,
		Title:                 chart.Title,
		TitleOverlay:          chart.TitleOverlay,
		Categories:            copyStringSlice(chart.Categories),
		Values:                copyFloat64Slice(chart.Values),
		X:                     chart.X,
		Y:                     chart.Y,
		CX:                    chart.CX,
		CY:                    chart.CY,
		Color:                 normalizeHexColor(chart.BarColor),
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
		BarDir:                "bar",
		Grouping:              "clustered",
	}, true
}

func slideBarStackedChartSpec(slide SlideContent) (*pptxxml.ChartSpec, bool) {
	if slide.BarStacked == nil {
		return nil, false
	}
	chart := slide.BarStacked
	return &pptxxml.ChartSpec{
		Kind:                  pptxxml.ChartKindBarStacked,
		Title:                 chart.Title,
		TitleOverlay:          chart.TitleOverlay,
		Categories:            copyStringSlice(chart.Categories),
		Values:                copyFloat64Slice(chart.Values),
		X:                     chart.X,
		Y:                     chart.Y,
		CX:                    chart.CX,
		CY:                    chart.CY,
		Color:                 normalizeHexColor(chart.BarColor),
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
		BarDir:                "col",
		Grouping:              "stacked",
	}, true
}

func slideBarStacked100ChartSpec(slide SlideContent) (*pptxxml.ChartSpec, bool) {
	if slide.BarStacked100 == nil {
		return nil, false
	}
	chart := slide.BarStacked100
	return &pptxxml.ChartSpec{
		Kind:                  pptxxml.ChartKindBarStacked100,
		Title:                 chart.Title,
		TitleOverlay:          chart.TitleOverlay,
		Categories:            copyStringSlice(chart.Categories),
		Values:                copyFloat64Slice(chart.Values),
		X:                     chart.X,
		Y:                     chart.Y,
		CX:                    chart.CX,
		CY:                    chart.CY,
		Color:                 normalizeHexColor(chart.BarColor),
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
		BarDir:                "col",
		Grouping:              "percentStacked",
	}, true
}

func slideLineChartSpec(slide SlideContent) (*pptxxml.ChartSpec, bool) {
	if slide.Line == nil {
		return nil, false
	}
	chart := slide.Line
	return &pptxxml.ChartSpec{
		Kind:                  pptxxml.ChartKindLine,
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
		Grouping:              "standard",
		Smooth:                chart.Smooth,
	}, true
}

func slideLineMarkersChartSpec(slide SlideContent) (*pptxxml.ChartSpec, bool) {
	if slide.LineMarkers == nil {
		return nil, false
	}
	chart := slide.LineMarkers
	return &pptxxml.ChartSpec{
		Kind:                  pptxxml.ChartKindLineMarkers,
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
		Grouping:              "standard",
		ShowMarkers:           true,
	}, true
}

func slideLineStackedChartSpec(slide SlideContent) (*pptxxml.ChartSpec, bool) {
	if slide.LineStacked == nil {
		return nil, false
	}
	chart := slide.LineStacked
	return &pptxxml.ChartSpec{
		Kind:                  pptxxml.ChartKindLineStacked,
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
		Grouping:              "stacked",
	}, true
}
