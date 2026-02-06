package pptx

import (
	"archive/zip"
	"fmt"

	"github.com/vegito/goppt/internal/pptxxml"
)

type chartPart struct {
	slideIndex int
	partNumber int
	spec       pptxxml.ChartSpec
}

func buildChartParts(slides []SlideContent) []chartPart {
	out := make([]chartPart, 0)
	for i, slide := range slides {
		spec, ok := slideChartSpec(slide)
		if !ok {
			continue
		}
		out = append(out, chartPart{
			slideIndex: i,
			partNumber: len(out) + 1,
			spec:       *spec,
		})
	}
	return out
}

func chartPartBySlide(parts []chartPart) map[int]chartPart {
	bySlide := make(map[int]chartPart, len(parts))
	for _, part := range parts {
		bySlide[part.slideIndex] = part
	}
	return bySlide
}

func writeChartFiles(zw *zip.Writer, parts []chartPart) error {
	for _, part := range parts {
		path := fmt.Sprintf("ppt/charts/chart%d.xml", part.partNumber)
		content := pptxxml.ChartPartXML(&part.spec)
		if err := writeFile(zw, path, content); err != nil {
			return err
		}
	}
	return nil
}

func slideChartSpec(slide SlideContent) (*pptxxml.ChartSpec, bool) {
	if slide.Chart != nil {
		categories := make([]string, len(slide.Chart.Categories))
		copy(categories, slide.Chart.Categories)
		values := make([]float64, len(slide.Chart.Values))
		copy(values, slide.Chart.Values)
		return &pptxxml.ChartSpec{
			Kind:               pptxxml.ChartKindBar,
			Title:              slide.Chart.Title,
			Categories:         categories,
			Values:             values,
			X:                  slide.Chart.X,
			Y:                  slide.Chart.Y,
			CX:                 slide.Chart.CX,
			CY:                 slide.Chart.CY,
			Color:              normalizeHexColor(slide.Chart.BarColor),
			SeriesName:         slide.Chart.SeriesName,
			ShowLegend:         slide.Chart.ShowLegend,
			LegendPosition:     slide.Chart.LegendPosition,
			ShowDataLabels:     slide.Chart.ShowDataLabels,
			ShowMajorGridlines: slide.Chart.ShowMajorGridlines,
			CategoryAxisTitle:  slide.Chart.CategoryAxisTitle,
			ValueAxisTitle:     slide.Chart.ValueAxisTitle,
			ValueFormat:        slide.Chart.ValueFormat,
			MinValue:           copyFloat64Pointer(slide.Chart.MinValue),
			MaxValue:           copyFloat64Pointer(slide.Chart.MaxValue),
		}, true
	}
	if slide.Line != nil {
		categories := make([]string, len(slide.Line.Categories))
		copy(categories, slide.Line.Categories)
		values := make([]float64, len(slide.Line.Values))
		copy(values, slide.Line.Values)
		return &pptxxml.ChartSpec{
			Kind:               pptxxml.ChartKindLine,
			Title:              slide.Line.Title,
			Categories:         categories,
			Values:             values,
			X:                  slide.Line.X,
			Y:                  slide.Line.Y,
			CX:                 slide.Line.CX,
			CY:                 slide.Line.CY,
			Color:              normalizeHexColor(slide.Line.LineColor),
			SeriesName:         slide.Line.SeriesName,
			ShowLegend:         slide.Line.ShowLegend,
			LegendPosition:     slide.Line.LegendPosition,
			ShowDataLabels:     slide.Line.ShowDataLabels,
			ShowMajorGridlines: slide.Line.ShowMajorGridlines,
			CategoryAxisTitle:  slide.Line.CategoryAxisTitle,
			ValueAxisTitle:     slide.Line.ValueAxisTitle,
			ValueFormat:        slide.Line.ValueFormat,
			MinValue:           copyFloat64Pointer(slide.Line.MinValue),
			MaxValue:           copyFloat64Pointer(slide.Line.MaxValue),
			Smooth:             slide.Line.Smooth,
		}, true
	}
	if slide.Pie != nil {
		categories := make([]string, len(slide.Pie.Categories))
		copy(categories, slide.Pie.Categories)
		values := make([]float64, len(slide.Pie.Values))
		copy(values, slide.Pie.Values)
		return &pptxxml.ChartSpec{
			Kind:           pptxxml.ChartKindPie,
			Title:          slide.Pie.Title,
			Categories:     categories,
			Values:         values,
			X:              slide.Pie.X,
			Y:              slide.Pie.Y,
			CX:             slide.Pie.CX,
			CY:             slide.Pie.CY,
			SeriesName:     slide.Pie.SeriesName,
			ShowLegend:     slide.Pie.ShowLegend,
			LegendPosition: slide.Pie.LegendPosition,
			ShowDataLabels: slide.Pie.ShowDataLabels,
		}, true
	}
	if slide.Dough != nil {
		categories := make([]string, len(slide.Dough.Categories))
		copy(categories, slide.Dough.Categories)
		values := make([]float64, len(slide.Dough.Values))
		copy(values, slide.Dough.Values)
		return &pptxxml.ChartSpec{
			Kind:           pptxxml.ChartKindDoughnut,
			Title:          slide.Dough.Title,
			Categories:     categories,
			Values:         values,
			X:              slide.Dough.X,
			Y:              slide.Dough.Y,
			CX:             slide.Dough.CX,
			CY:             slide.Dough.CY,
			SeriesName:     slide.Dough.SeriesName,
			ShowLegend:     slide.Dough.ShowLegend,
			LegendPosition: slide.Dough.LegendPosition,
			ShowDataLabels: slide.Dough.ShowDataLabels,
			HoleSize:       slide.Dough.HoleSize,
		}, true
	}
	return nil, false
}

func copyFloat64Pointer(value *float64) *float64 {
	if value == nil {
		return nil
	}
	v := *value
	return &v
}
