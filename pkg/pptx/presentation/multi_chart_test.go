package presentation

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// Setting two charts on one slide used to emit only the first: the generator
// asked slideChartSpec for one spec and stopped there.
func TestSlideCarriesMoreThanOneChart(t *testing.T) {
	bar := charts.NewBarChart([]string{"A", "B"}, []float64{1, 2})
	line := charts.NewLineChart([]string{"A", "B"}, []float64{3, 4}).
		Position(styling.Inches(5), styling.Inches(1))

	slide := elements.NewSlide("Two charts")
	slide.Chart = &bar
	slide.Line = &line

	parts := buildPackageParts(t, Metadata{}, []elements.SlideContent{slide})

	for _, want := range []string{"ppt/charts/chart1.xml", "ppt/charts/chart2.xml"} {
		if _, ok := parts[want]; !ok {
			t.Fatalf("package is missing %s", want)
		}
	}
	if got := strings.Count(parts["ppt/slides/slide1.xml"], "<p:graphicFrame>"); got != 2 {
		t.Fatalf("slide has %d graphic frames, want 2", got)
	}
	if got := strings.Count(parts["ppt/slides/_rels/slide1.xml.rels"], "relationships/chart"); got != 2 {
		t.Fatalf("slide rels declare %d charts, want 2", got)
	}
}

// The Charts overflow places charts of the same kind, which the typed fields
// cannot express.
func TestChartsOverflowIsEmitted(t *testing.T) {
	first := charts.NewPieChart([]string{"A"}, []float64{1})
	second := charts.NewPieChart([]string{"B"}, []float64{2}).
		Position(styling.Inches(5), styling.Inches(1))

	slide := elements.NewSlide("Two pies")
	slide.Pie = &first
	slide.Charts = []elements.ChartDefinition{second}

	parts := buildPackageParts(t, Metadata{}, []elements.SlideContent{slide})
	if _, ok := parts["ppt/charts/chart2.xml"]; !ok {
		t.Fatal("the overflow chart was not written")
	}
}
