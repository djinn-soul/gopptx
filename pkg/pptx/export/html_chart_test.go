package export

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestHTMLPlotsABarChart(t *testing.T) {
	slide := elements.NewSlide("Revenue").WithBarChart(charts.BarChart{
		Title:      "Revenue by quarter",
		Categories: []string{"Q1", "Q2", "Q3"},
		Values:     []float64{10, 25, 15},
	})

	out := HTML("Deck", []elements.SlideContent{slide})
	if !strings.Contains(out, "Revenue by quarter") {
		t.Error("the chart title is missing")
	}
	if !strings.Contains(out, `class="chart-svg"`) {
		t.Error("no chart SVG was emitted")
	}
	if strings.Count(out, "<rect") < 3 {
		t.Errorf("got %d rects, want one bar per category", strings.Count(out, "<rect"))
	}
	for _, label := range []string{"Q1", "Q2", "Q3"} {
		if !strings.Contains(out, ">"+label+"<") {
			t.Errorf("category %q is not labelled", label)
		}
	}
}

func TestHTMLPlotsLineAndPieCharts(t *testing.T) {
	// A slide carries one chart: the With*Chart builders clear the others.
	line := elements.NewSlide("Trend").WithLineChart(charts.LineChart{
		Categories: []string{"Jan", "Feb"},
		Values:     []float64{4, 8},
	})
	pie := elements.NewSlide("Split").WithPieChart(charts.PieChart{
		Categories: []string{"A", "B"},
		Values:     []float64{60, 40},
	})

	out := HTML("Deck", []elements.SlideContent{line, pie})
	if !strings.Contains(out, "<polyline") {
		t.Error("the line chart drew no line")
	}
	if !strings.Contains(out, "<path d=\"M ") {
		t.Error("the pie chart drew no slices")
	}
}

func TestHTMLTabulatesChartsWithNoSimplePlot(t *testing.T) {
	slide := elements.NewSlide("Scatter").WithScatterChart(charts.ScatterChart{
		Title:   "Correlation",
		XValues: []float64{1, 2, 3},
		YValues: []float64{4, 5, 6},
	})

	out := HTML("Deck", []elements.SlideContent{slide})
	if !strings.Contains(out, `class="chart-data"`) {
		t.Error("the scatter chart was neither plotted nor tabulated")
	}
	if !strings.Contains(out, "Correlation") {
		t.Error("the chart title is missing")
	}
	for _, value := range []string{"4", "5", "6"} {
		if !strings.Contains(out, ">"+value+"<") {
			t.Errorf("value %q is missing from the table", value)
		}
	}
}

func TestSlideHTMLChartsCoversEveryChartField(t *testing.T) {
	categories := []string{"A", "B"}
	values := []float64{1, 2}
	slide := elements.NewSlide("All charts")
	slide.Chart = &charts.BarChart{Categories: categories, Values: values}
	slide.Line = &charts.LineChart{Categories: categories, Values: values}
	slide.Pie = &charts.PieChart{Categories: categories, Values: values}
	slide.Area = &charts.AreaChart{Categories: categories, Values: values}
	slide.Doughnut = &charts.DoughnutChart{Categories: categories, Values: values}
	slide.Radar = &charts.RadarChart{Categories: categories, Values: values}
	slide.StockHLC = &charts.StockHLCChart{
		Categories:  categories,
		HighValues:  values,
		LowValues:   values,
		CloseValues: values,
	}

	got := slideHTMLCharts(slide)
	if len(got) != 7 {
		t.Fatalf("got %d charts, want 7 — one per populated field", len(got))
	}
	for _, chart := range got {
		if !chart.hasData() {
			t.Errorf("chart %+v came through with no data", chart)
		}
	}
}

func TestEmptyChartIsNotEmitted(t *testing.T) {
	slide := elements.NewSlide("Empty")
	slide.Chart = &charts.BarChart{Categories: []string{"A"}}

	if got := slideHTMLCharts(slide); len(got) != 0 {
		t.Errorf("got %d charts for an empty one, want 0", len(got))
	}
}

func TestHTMLDrawsConnectors(t *testing.T) {
	connector := shapes.NewStraightConnector(
		styling.Inches(1), styling.Inches(1), styling.Inches(3), styling.Inches(2),
	).WithArrows(shapes.ArrowTypeNone, shapes.ArrowTypeTriangle)
	connector.Label = "depends on"

	slide := elements.NewSlide("Wiring").AddConnector(connector)
	out := HTML("Deck", []elements.SlideContent{slide})

	if !strings.Contains(out, "<path d=\"M ") {
		t.Error("the connector line is missing")
	}
	if !strings.Contains(out, "<polygon") {
		t.Error("the arrowhead is missing")
	}
	if !strings.Contains(out, "depends on") {
		t.Error("the connector label is missing")
	}
}

func TestElbowConnectorRoutesAroundACorner(t *testing.T) {
	elbow := shapes.NewElbowConnector(
		styling.Inches(1), styling.Inches(1), styling.Inches(3), styling.Inches(2),
	)
	path := connectorPath(elbow, 100, 100, 300, 200)
	if strings.Count(path, "L") != 3 {
		t.Errorf("got %q, want a three-segment dog-leg", path)
	}
}

func TestDecorativeConnectorIsSkipped(t *testing.T) {
	connector := shapes.NewStraightConnector(0, 0, styling.Inches(1), styling.Inches(1))
	connector.IsDecorative = true
	if got := renderConnector(connector); got != "" {
		t.Errorf("got %q, want nothing for a decorative connector", got)
	}
}
