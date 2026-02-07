package pptx

import (
	"archive/zip"
	"fmt"

	"github.com/djinn09/goppt/internal/pptxxml"
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

type slideChartSpecResolver func(slide SlideContent) (*pptxxml.ChartSpec, bool)

var slideChartSpecResolvers = []slideChartSpecResolver{
	slideBarChartSpec,
	slideBarHorizontalChartSpec,
	slideBarStackedChartSpec,
	slideBarStacked100ChartSpec,
	slideLineChartSpec,
	slideLineMarkersChartSpec,
	slideLineStackedChartSpec,
	slideScatterChartSpec,
	slideAreaChartSpec,
	slideAreaStackedChartSpec,
	slideAreaStacked100ChartSpec,
	slidePieChartSpec,
	slideDoughnutChartSpec,
	slideBubbleChartSpec,
	slideRadarChartSpec,
	slideRadarFilledChartSpec,
	slideStockHLCChartSpec,
	slideStockOHLCChartSpec,
	slideComboChartSpec,
}

func slideChartSpec(slide SlideContent) (*pptxxml.ChartSpec, bool) {
	for _, resolver := range slideChartSpecResolvers {
		spec, ok := resolver(slide)
		if ok {
			return spec, true
		}
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
