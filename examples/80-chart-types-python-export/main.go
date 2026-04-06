// Package main demonstrates all supported chart types using the Go pptx API.
// It is a direct Go translation of main.py in this directory.
//
// Run from the repository root:
//
//	go run ./examples/80-chart-types-python-export/
package main

import (
	"fmt"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/export"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	outputDir = "examples/output"
	pptxPath  = outputDir + "/80_chart_types_go_export.pptx"
	pdfPath   = outputDir + "/80_chart_types_go_export.pdf"
)

// chartBounds mirrors Python CHART_BOUNDS = (Inches(0.8), Inches(1.4), Inches(8.5), Inches(4.6)).
var (
	chartX = styling.Inches(0.8)
	chartY = styling.Inches(1.4)
	chartW = styling.Inches(8.5)
	chartH = styling.Inches(4.6)

	cats = []string{"Q1", "Q2", "Q3", "Q4"}
	vals = []float64{14, 21, 18, 27}

	hlcCats = []string{"D1", "D2", "D3", "D4"}
	high    = []float64{16, 23, 20, 29}
	low     = []float64{12, 18, 15, 24}
	closeV  = []float64{14, 21, 18, 27}
	openV   = []float64{13, 20, 17, 26}

	xVals = []float64{1, 2, 3, 4}
	yVals = []float64{14, 21, 18, 27}
	sizes = []float64{14, 21, 18, 27}
)

// chartTypeInfo holds a display name and the slide builder for each chart kind.
// Ordered alphabetically by internal value to match the Python sorted(set(...)) order.
type chartTypeInfo struct {
	name    string
	builder func() pptx.SlideContent
}

func titleOnlySlide(title string) pptx.SlideContent {
	s := pptx.NewSlide(title)
	s.Layout = pptx.SlideLayoutTitleOnly
	return s
}

// allChartTypes returns all chart type entries in the same order as Python's
// sorted(set(ChartType.get_all().values())) — alphabetical by raw value.
func allChartTypes() []chartTypeInfo {
	return []chartTypeInfo{
		{
			name: "Area",
			builder: func() pptx.SlideContent {
				return titleOnlySlide("Area Chart").WithAreaChart(
					pptx.NewAreaChart(cats, vals).WithTitle("Area Demo").
						Position(chartX, chartY).Size(chartW, chartH),
				)
			},
		},
		{
			name: "Area Stacked",
			builder: func() pptx.SlideContent {
				return titleOnlySlide("Area Stacked Chart").WithAreaStackedChart(
					pptx.NewAreaStackedChart(cats, vals).WithTitle("Area Stacked Demo").
						Position(chartX, chartY).Size(chartW, chartH),
				)
			},
		},
		{
			name: "Area Stacked 100%",
			builder: func() pptx.SlideContent {
				return titleOnlySlide("Area Stacked 100% Chart").WithAreaStacked100Chart(
					pptx.NewAreaStacked100Chart(cats, vals).WithTitle("Area Stacked 100% Demo").
						Position(chartX, chartY).Size(chartW, chartH),
				)
			},
		},
		{
			name: "Column / Bar",
			builder: func() pptx.SlideContent {
				return titleOnlySlide("Column / Bar Chart").WithBarChart(
					pptx.NewBarChart(cats, vals).WithTitle("Column / Bar Demo").
						Position(chartX, chartY).Size(chartW, chartH),
				)
			},
		},
		{
			name: "Bar Horizontal",
			builder: func() pptx.SlideContent {
				return titleOnlySlide("Bar Horizontal Chart").WithBarHorizontalChart(
					pptx.NewBarHorizontalChart(cats, vals).WithTitle("Bar Horizontal Demo").
						Position(chartX, chartY).Size(chartW, chartH),
				)
			},
		},
		{
			name: "Bar Stacked",
			builder: func() pptx.SlideContent {
				return titleOnlySlide("Bar Stacked Chart").WithBarStackedChart(
					pptx.NewBarStackedChart(cats, vals).WithTitle("Bar Stacked Demo").
						Position(chartX, chartY).Size(chartW, chartH),
				)
			},
		},
		{
			name: "Bar Stacked 100%",
			builder: func() pptx.SlideContent {
				return titleOnlySlide("Bar Stacked 100% Chart").WithBarStacked100Chart(
					pptx.NewBarStacked100Chart(cats, vals).WithTitle("Bar Stacked 100% Demo").
						Position(chartX, chartY).Size(chartW, chartH),
				)
			},
		},
		{
			name: "Bubble",
			builder: func() pptx.SlideContent {
				return titleOnlySlide("Bubble Chart").WithBubbleChart(
					pptx.NewBubbleChart(xVals, yVals, sizes).WithTitle("Bubble Demo").
						Position(chartX.Emu(), chartY.Emu()).Size(chartW.Emu(), chartH.Emu()),
				)
			},
		},
		{
			name: "Doughnut",
			builder: func() pptx.SlideContent {
				return titleOnlySlide("Doughnut Chart").WithDoughnutChart(
					pptx.NewDoughnutChart(cats, vals).WithTitle("Doughnut Demo").
						Position(chartX, chartY).Size(chartW, chartH),
				)
			},
		},
		{
			name: "Line",
			builder: func() pptx.SlideContent {
				return titleOnlySlide("Line Chart").WithLineChart(
					pptx.NewLineChart(cats, vals).WithTitle("Line Demo").
						Position(chartX, chartY).Size(chartW, chartH),
				)
			},
		},
		{
			name: "Line with Markers",
			builder: func() pptx.SlideContent {
				return titleOnlySlide("Line with Markers Chart").WithLineMarkersChart(
					pptx.NewLineMarkersChart(cats, vals).WithTitle("Line with Markers Demo").
						Position(chartX, chartY).Size(chartW, chartH),
				)
			},
		},
		{
			name: "Line Stacked",
			builder: func() pptx.SlideContent {
				return titleOnlySlide("Line Stacked Chart").WithLineStackedChart(
					pptx.NewLineStackedChart(cats, vals).WithTitle("Line Stacked Demo").
						Position(chartX, chartY).Size(chartW, chartH),
				)
			},
		},
		{
			name: "Pie",
			builder: func() pptx.SlideContent {
				return titleOnlySlide("Pie Chart").WithPieChart(
					pptx.NewPieChart(cats, vals).WithTitle("Pie Demo").
						Position(chartX, chartY).Size(chartW, chartH),
				)
			},
		},
		{
			name: "Radar",
			builder: func() pptx.SlideContent {
				return titleOnlySlide("Radar Chart").WithRadarChart(
					pptx.NewRadarChart(cats, vals).WithTitle("Radar Demo").
						Position(chartX, chartY).Size(chartW, chartH),
				)
			},
		},
		{
			name: "Radar Filled",
			builder: func() pptx.SlideContent {
				return titleOnlySlide("Radar Filled Chart").WithRadarFilledChart(
					pptx.NewRadarFilledChart(cats, vals).WithTitle("Radar Filled Demo").
						Position(chartX, chartY).Size(chartW, chartH),
				)
			},
		},
		{
			name: "Scatter",
			builder: func() pptx.SlideContent {
				return titleOnlySlide("Scatter Chart").WithScatterChart(
					pptx.NewScatterChart(xVals, yVals).WithTitle("Scatter Demo").
						Position(chartX, chartY).Size(chartW, chartH),
				)
			},
		},
		{
			name: "Stock HLC",
			builder: func() pptx.SlideContent {
				return titleOnlySlide("Stock HLC Chart").WithStockHLCChart(
					pptx.NewStockHLCChart(hlcCats, high, low, closeV).WithTitle("Stock HLC Demo").
						Position(chartX, chartY).Size(chartW, chartH),
				)
			},
		},
		{
			name: "Stock OHLC",
			builder: func() pptx.SlideContent {
				return titleOnlySlide("Stock OHLC Chart").WithStockOHLCChart(
					pptx.NewStockOHLCChart(hlcCats, openV, high, low, closeV).WithTitle("Stock OHLC Demo").
						Position(chartX, chartY).Size(chartW, chartH),
				)
			},
		},
	}
}

func comboSlide() pptx.SlideContent {
	return titleOnlySlide("Combo Chart").WithComboChart(
		pptx.NewComboChart(
			cats,
			[]pptx.Series{{Name: "Revenue", Values: []float64{180, 220, 210, 260}}},
			[]pptx.Series{{Name: "Growth %", Values: []float64{8, 11, 10, 14}}},
		).WithTitle("Combo Demo").Position(chartX, chartY).Size(chartW, chartH),
	)
}

// goChartConstructors lists every public Go constructor paired with its internal chart kind value.
// Mirrors Python's ChartType.get_all() — same 19 unique kinds + combo.
var goChartConstructors = []struct{ constructor, kind string }{
	{"pptx.NewBarChart", "bar"},           // COLUMN
	{"pptx.NewBarChart (BAR alias)", "bar"}, // BAR — same constructor, mirrors ChartType.BAR = ChartType.COLUMN
	{"pptx.NewBarHorizontalChart", "barHorizontal"},
	{"pptx.NewBarStackedChart", "barStacked"},
	{"pptx.NewBarStacked100Chart", "barStacked100"},
	{"pptx.NewLineChart", "line"},
	{"pptx.NewLineMarkersChart", "lineMarkers"},
	{"pptx.NewLineStackedChart", "lineStacked"},
	{"pptx.NewScatterChart", "scatter"},
	{"pptx.NewAreaChart", "area"},
	{"pptx.NewAreaStackedChart", "areaStacked"},
	{"pptx.NewAreaStacked100Chart", "areaStacked100"},
	{"pptx.NewPieChart", "pie"},
	{"pptx.NewDoughnutChart", "doughnut"},
	{"pptx.NewBubbleChart", "bubble"},
	{"pptx.NewRadarChart", "radar"},
	{"pptx.NewRadarFilledChart", "radarFilled"},
	{"pptx.NewStockHLCChart", "stockHLC"},
	{"pptx.NewStockOHLCChart", "stockOHLC"},
	{"pptx.NewComboChart", "combo"},
}

// referenceSlides mirrors Python add_chart_surface_reference — split across two slides to avoid overflow.
// Shows Go constructor names and their internal chart kind values, equivalent to
// Python's "ChartType.COLUMN = 'bar'" format.
func referenceSlides() []pptx.SlideContent {
	lines := make([]string, 0, len(goChartConstructors))
	for _, c := range goChartConstructors {
		lines = append(lines, fmt.Sprintf("%s = %q", c.constructor, c.kind))
	}
	half := len(lines) / 2

	uniqueKinds := make(map[string]struct{}, len(goChartConstructors))
	for _, c := range goChartConstructors {
		uniqueKinds[c.kind] = struct{}{}
	}
	s1 := pptx.NewSlide("Go Chart Surface (1/2)").
		AddBullet(fmt.Sprintf("Named constructors: %d", len(goChartConstructors))).
		AddBullet(fmt.Sprintf("Unique chart kinds: %d", len(uniqueKinds)))
	for _, l := range lines[:half] {
		s1 = s1.AddBullet(l)
	}

	s2 := pptx.NewSlide("Go Chart Surface (2/2)")
	for _, l := range lines[half:] {
		s2 = s2.AddBullet(l)
	}

	return []pptx.SlideContent{s1, s2}
}

func buildSlides() []pptx.SlideContent {
	chartTypes := allChartTypes()

	slides := []pptx.SlideContent{
		// Intro — mirrors Python add_bullet_slide
		pptx.NewSlide("Chart Types Demo").
			AddBullet("This deck demonstrates chart creation from Go API.").
			AddBullet("It outputs both PPTX and PDF artifacts."),
	}

	// One slide per chart type (mirrors Python for chart_type in CATEGORY_CHART_VALUES).
	for _, ct := range chartTypes {
		slides = append(slides, ct.builder())
	}

	// Combo slide.
	slides = append(slides, comboSlide())

	// Reference slides.
	slides = append(slides, referenceSlides()...)

	return slides
}

func main() {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		fmt.Fprintf(os.Stderr, "error: create output dir: %v\n", err)
		os.Exit(1)
	}

	const title = "Go Chart Types Export Demo"
	slides := buildSlides()

	// Save PPTX.
	data, err := pptx.CreateWithSlides(title, slides)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: create presentation: %v\n", err)
		os.Exit(1)
	}
	if err = os.WriteFile(pptxPath, data, 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "error: write pptx: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Saved PPTX: %s (%d slides)\n", pptxPath, len(slides))

	// Export PDF via native renderer — same driver as the Python main.py.
	opts := export.PDFOptions{Driver: export.PDFDriverNative}
	if err = export.PDFWithOptions(title, slides, pdfPath, opts); err != nil {
		fmt.Fprintf(os.Stderr, "error: export pdf: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Saved PDF:  %s\n", pdfPath)
}
