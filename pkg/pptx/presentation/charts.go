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
		for _, spec := range slideChartSpecs(slide) {
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

// slideChartSpecs returns every chart placed directly on the slide, in field
// order, followed by the Charts overflow. Setting two charts used to emit only
// the first: the generator asked for one spec and stopped.
func slideChartSpecs(slide elements.SlideContent) []*pptxxml.ChartSpec {
	typed := make([]elements.ChartDefinition, 0, 1+len(slide.Charts))
	typed = appendCartesianCharts(typed, slide)
	typed = appendRoundAndSeriesCharts(typed, slide)
	typed = append(typed, slide.Charts...)

	specs := make([]*pptxxml.ChartSpec, 0, len(typed))
	for _, chart := range typed {
		if chart == nil {
			continue
		}
		specs = append(specs, chart.ToChartSpec())
	}
	return specs
}

// appendCartesianCharts collects the bar, line, scatter and area fields.
func appendCartesianCharts(
	out []elements.ChartDefinition,
	slide elements.SlideContent,
) []elements.ChartDefinition {
	if slide.Chart != nil {
		out = append(out, *slide.Chart)
	}
	if slide.BarHorizontal != nil {
		out = append(out, *slide.BarHorizontal)
	}
	if slide.BarStacked != nil {
		out = append(out, *slide.BarStacked)
	}
	if slide.BarStacked100 != nil {
		out = append(out, *slide.BarStacked100)
	}
	if slide.Line != nil {
		out = append(out, *slide.Line)
	}
	if slide.LineMarkers != nil {
		out = append(out, *slide.LineMarkers)
	}
	if slide.LineStacked != nil {
		out = append(out, *slide.LineStacked)
	}
	if slide.Scatter != nil {
		out = append(out, *slide.Scatter)
	}
	if slide.Area != nil {
		out = append(out, *slide.Area)
	}
	if slide.AreaStacked != nil {
		out = append(out, *slide.AreaStacked)
	}
	if slide.AreaStacked100 != nil {
		out = append(out, *slide.AreaStacked100)
	}
	return out
}

// appendRoundAndSeriesCharts collects the pie family plus bubble, radar, stock
// and combo.
func appendRoundAndSeriesCharts(
	out []elements.ChartDefinition,
	slide elements.SlideContent,
) []elements.ChartDefinition {
	if slide.Pie != nil {
		out = append(out, *slide.Pie)
	}
	if slide.Pie3D != nil {
		out = append(out, *slide.Pie3D)
	}
	if slide.Doughnut != nil {
		out = append(out, *slide.Doughnut)
	}
	if slide.Bubble != nil {
		out = append(out, *slide.Bubble)
	}
	if slide.Radar != nil {
		out = append(out, *slide.Radar)
	}
	if slide.RadarFilled != nil {
		out = append(out, *slide.RadarFilled)
	}
	if slide.StockHLC != nil {
		out = append(out, *slide.StockHLC)
	}
	if slide.StockOHLC != nil {
		out = append(out, *slide.StockOHLC)
	}
	if slide.Combo != nil {
		out = append(out, *slide.Combo)
	}
	return out
}

// slideChartCount is how many chart parts the slide itself owns, which is where
// the placeholder charts start in the slide's chart list.
func slideChartCount(slide elements.SlideContent) int {
	return len(slideChartSpecs(slide))
}
