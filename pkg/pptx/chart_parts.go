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
			Kind:       pptxxml.ChartKindBar,
			Title:      slide.Chart.Title,
			Categories: categories,
			Values:     values,
			X:          slide.Chart.X,
			Y:          slide.Chart.Y,
			CX:         slide.Chart.CX,
			CY:         slide.Chart.CY,
			Color:      normalizeHexColor(slide.Chart.BarColor),
		}, true
	}
	if slide.Line != nil {
		categories := make([]string, len(slide.Line.Categories))
		copy(categories, slide.Line.Categories)
		values := make([]float64, len(slide.Line.Values))
		copy(values, slide.Line.Values)
		return &pptxxml.ChartSpec{
			Kind:       pptxxml.ChartKindLine,
			Title:      slide.Line.Title,
			Categories: categories,
			Values:     values,
			X:          slide.Line.X,
			Y:          slide.Line.Y,
			CX:         slide.Line.CX,
			CY:         slide.Line.CY,
			Color:      normalizeHexColor(slide.Line.LineColor),
		}, true
	}
	return nil, false
}
