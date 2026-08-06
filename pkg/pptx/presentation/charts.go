package presentation

import (
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	editorchart "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/chart"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

type ChartPart struct {
	slideIndex int
	partNumber int
	spec       pptxxml.ChartSpec
}

func ChartPartCount(parts []ChartPart) int {
	return len(parts)
}

func BuildChartParts(slides []elements.SlideContent) []ChartPart {
	out := make([]ChartPart, 0)
	for i, slide := range slides {
		spec, ok := slideChartSpec(slide)
		if ok {
			out = append(out, ChartPart{
				slideIndex: i,
				partNumber: len(out) + 1,
				spec:       *spec,
			})
		}

		for _, override := range slide.PlaceholderOverrides {
			if override.Chart != nil {
				out = append(out, ChartPart{
					slideIndex: i,
					partNumber: len(out) + 1,
					spec:       *override.Chart.ToChartSpec(),
				})
			}
		}
	}
	return out
}

func chartPartBySlide(parts []ChartPart) map[int][]ChartPart {
	bySlide := make(map[int][]ChartPart, len(parts))
	for _, part := range parts {
		bySlide[part.slideIndex] = append(bySlide[part.slideIndex], part)
	}
	return bySlide
}

// chartEmbeddingRelID is the only relationship a generated chart part carries,
// so it can be a fixed id.
const chartEmbeddingRelID = "rId1"

// ChartEmbeddingPartName is the workbook path for the nth generated chart.
func ChartEmbeddingPartName(partNumber int) string {
	return fmt.Sprintf("ppt/embeddings/Microsoft_Excel_Worksheet%d.xlsx", partNumber)
}

// hasEmbeddableData reports whether the chart's data is the category/value
// shape the embedded workbook holds. Scatter, bubble and combo charts are not,
// and keep their literal cached values with no workbook.
func (p ChartPart) hasEmbeddableData() bool {
	return len(p.spec.Categories) > 0 && len(p.spec.Categories) == len(p.spec.Values)
}

// ChartEmbeddingNumbers lists the chart parts that ship an embedded workbook,
// by part number. The package manifest needs it before the parts are written.
func ChartEmbeddingNumbers(parts []ChartPart) []int {
	numbers := make([]int, 0, len(parts))
	for _, part := range parts {
		if part.hasEmbeddableData() {
			numbers = append(numbers, part.partNumber)
		}
	}
	return numbers
}

// writeChartFiles writes each chart part together with the workbook it was
// built from: the .xlsx, the chart's .rels pointing at it, and the
// <c:externalData> reference inside the chart. Without all three, "Edit Data"
// on a generated chart finds no workbook to open.
func writeChartFiles(pw *pptxxml.PackageWriter, parts []ChartPart) error {
	for _, part := range parts {
		spec := part.spec

		if part.hasEmbeddableData() {
			workbook, err := editorchart.GenerateExcelForChart(spec.Categories, spec.Values)
			if err != nil {
				return fmt.Errorf("chart %d: generate embedded workbook: %w", part.partNumber, err)
			}
			spec.ExternalDataID = chartEmbeddingRelID
			pw.AddBinaryPart(ChartEmbeddingPartName(part.partNumber), workbook)
			pw.AddPart(
				fmt.Sprintf("ppt/charts/_rels/chart%d.xml.rels", part.partNumber),
				pptxxml.ChartRelationships(
					chartEmbeddingRelID,
					fmt.Sprintf("../embeddings/Microsoft_Excel_Worksheet%d.xlsx", part.partNumber),
				),
			)
		}

		pw.AddPart(
			fmt.Sprintf("ppt/charts/chart%d.xml", part.partNumber),
			pptxxml.ChartPartXML(&spec),
		)
	}
	return nil
}

func slideChartSpec(slide elements.SlideContent) (*pptxxml.ChartSpec, bool) {
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
	if slide.Pie3D != nil {
		return slide.Pie3D.ToChartSpec(), true
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

func slideChartKindDefined(slide elements.SlideContent) bool {
	_, ok := slideChartSpec(slide)
	return ok
}
