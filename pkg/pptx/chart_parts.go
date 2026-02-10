package pptx

import (
	"archive/zip"
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
)

type chartPart struct {
	slideIndex int
	partNumber int
	spec       pptxxml.ChartSpec
}

func buildChartParts(slides []SlideContent) []chartPart {
	out := make([]chartPart, 0)
	for i, slide := range slides {
		// Existing single chart
		spec, ok := slideChartSpec(slide)
		if ok {
			out = append(out, chartPart{
				slideIndex: i,
				partNumber: len(out) + 1,
				spec:       *spec,
			})
		}

		// Charts in placeholders
		for _, override := range slide.PlaceholderOverrides {
			if override.Chart != nil {
				out = append(out, chartPart{
					slideIndex: i,
					partNumber: len(out) + 1,
					spec:       *override.Chart.ToChartSpec(),
				})
			}
		}
	}
	return out
}

func chartPartBySlide(parts []chartPart) map[int][]chartPart {
	bySlide := make(map[int][]chartPart, len(parts))
	for _, part := range parts {
		bySlide[part.slideIndex] = append(bySlide[part.slideIndex], part)
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
		return slide.Chart.ToChartSpec(), true
	}
	if slide.BarHorizontal != nil {
		return slide.BarHorizontal.ToChartSpec(), true
	}
	if slide.BarStacked != nil {
		return slide.BarStacked.ToChartSpec(), true
	}
	if slide.BarStacked100 != nil {
		return slide.BarStacked100.ToChartSpec(), true
	}
	if slide.Line != nil {
		return slide.Line.ToChartSpec(), true
	}
	if slide.LineMarkers != nil {
		return slide.LineMarkers.ToChartSpec(), true
	}
	if slide.LineStacked != nil {
		return slide.LineStacked.ToChartSpec(), true
	}
	if slide.Scatter != nil {
		return slide.Scatter.ToChartSpec(), true
	}
	if slide.Area != nil {
		return slide.Area.ToChartSpec(), true
	}
	if slide.AreaStacked != nil {
		return slide.AreaStacked.ToChartSpec(), true
	}
	if slide.AreaStacked100 != nil {
		return slide.AreaStacked100.ToChartSpec(), true
	}
	if slide.Pie != nil {
		return slide.Pie.ToChartSpec(), true
	}
	if slide.Doughnut != nil {
		return slide.Doughnut.ToChartSpec(), true
	}
	if slide.Bubble != nil {
		return slide.Bubble.ToChartSpec(), true
	}
	if slide.Radar != nil {
		return slide.Radar.ToChartSpec(), true
	}
	if slide.RadarFilled != nil {
		return slide.RadarFilled.ToChartSpec(), true
	}
	if slide.StockHLC != nil {
		return slide.StockHLC.ToChartSpec(), true
	}
	if slide.StockOHLC != nil {
		return slide.StockOHLC.ToChartSpec(), true
	}
	if slide.Combo != nil {
		return slide.Combo.ToChartSpec(), true
	}
	return nil, false
}

func slideChartKindDefined(slide SlideContent) bool {
	_, ok := slideChartSpec(slide)
	return ok
}
