// examples/69-chart-api demonstrates chart formatting options.
//
// Shows axis titles, legend position/overlay, data labels, gridlines,
// value ranges, tick label positions, axis crossing modes, and value format
// codes. Covers bar, line, scatter, area, radar, and combo chart types.
//
// Run with: go run ./examples/69-chart-api/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

const (
	outputDir  = "examples/output"
	outputFile = "69_chart_api.pptx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}

	slides := []pptx.SlideContent{
		buildBarFullSlide(months),
		buildStackedBarSlide(months),
		buildStacked100Slide(months),
		buildLineSmoothSlide(months),
		buildLineMarkersSlide(months),
		buildScatterSlide(),
		buildAreaSlide(months),
		buildRadarSlide(),
		buildLegendConstantsSlide(),
	}

	outputPath := filepath.Join(outputDir, outputFile)
	data, err := pptx.CreateWithSlides("Chart API Demo", slides)
	if err != nil {
		return fmt.Errorf("create presentation: %w", err)
	}
	if err := os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}

func buildBarFullSlide(months []string) pptx.SlideContent {
	barFull := pptx.NewBarChart(months, []float64{12, 19, 15, 22, 28, 24}).
		WithSeriesName("Monthly Sales").
		WithLegend(true).
		WithLegendPosition(pptx.LegendPositionBottom).
		WithLegendOverlay(false).
		WithTitleOverlay(false).
		WithDataLabels(true).
		WithMajorGridlines(true).
		WithCategoryMajorGridlines(false).
		WithAxisTitles("Month", "Units Sold").
		WithValueFormat("#,##0").
		WithValueRange(0, 35).
		WithValueAxisCrossBetween(pptx.ValueAxisCrossBetweenBetween).
		WithTickLabelPositions(
			charts.AxisTickLabelPositionNextTo, charts.AxisTickLabelPositionNextTo,
		).
		WithAxisCrosses(charts.AxisCrossesAutoZero, charts.AxisCrossesAutoZero)
	barFull.Title = "Monthly Sales (Bar Chart – Full Options)"
	barFull.BarColor = "4472C4"
	return pptx.NewSlide("Bar Chart – Full Formatting").WithBarChart(barFull)
}

func buildStackedBarSlide(months []string) pptx.SlideContent {
	stackedBar := pptx.NewBarStackedChart(months, []float64{10, 14, 12, 18, 22, 20}).
		WithSeriesName("Product A").
		WithLegend(true).
		WithAxisTitles("Month", "Total")
	stackedBar.Title = "Stacked Bar Chart"
	return pptx.NewSlide("Stacked Bar Chart").WithBarStackedChart(stackedBar)
}

func buildStacked100Slide(months []string) pptx.SlideContent {
	stacked100 := pptx.NewBarStacked100Chart(months, []float64{60, 55, 70, 65, 50, 75}).
		WithSeriesName("Product A (%)").
		WithLegend(true).
		WithAxisTitles("Month", "Percentage")
	stacked100.Title = "100% Stacked Bar Chart"
	return pptx.NewSlide("100% Stacked Bar Chart").WithBarStacked100Chart(stacked100)
}

func buildLineSmoothSlide(months []string) pptx.SlideContent {
	lineSmooth := pptx.NewLineChart(months, []float64{5, 7, 6, 10, 12, 9}).
		WithSeriesName("Revenue Trend").
		WithLegend(true).
		WithSmooth(true).
		WithDataLabels(false).
		WithAxisTitles("Month", "Revenue ($k)").
		WithMajorGridlines(true)
	lineSmooth.Title = "Line Chart (Smoothed)"
	return pptx.NewSlide("Line Chart Variants").WithLineChart(lineSmooth)
}

func buildLineMarkersSlide(months []string) pptx.SlideContent {
	lineMarkers := pptx.NewLineMarkersChart(months, []float64{8, 11, 9, 14, 16, 13}).
		WithSeriesName("Costs").
		WithLegend(true)
	lineMarkers.Title = "Line Chart with Markers"
	return pptx.NewSlide("Line with Markers").WithLineMarkersChart(lineMarkers)
}

func buildScatterSlide() pptx.SlideContent {
	scatterChart := pptx.NewScatterChart(
		[]float64{1, 2, 3, 4, 5, 6},
		[]float64{2.5, 3.1, 4.0, 3.7, 5.2, 6.1},
	)
	scatterChart.Title = "Scatter Chart"
	return pptx.NewSlide("Scatter Chart").WithScatterChart(scatterChart)
}

func buildAreaSlide(months []string) pptx.SlideContent {
	areaChart := pptx.NewAreaChart(months, []float64{10, 15, 13, 18, 22, 20}).
		WithSeriesName("Area Series").
		WithLegend(true).
		WithAxisTitles("Month", "Value")
	areaChart.Title = "Area Chart"
	return pptx.NewSlide("Area Chart").WithAreaChart(areaChart)
}

func buildRadarSlide() pptx.SlideContent {
	radarCategories := []string{"Speed", "Quality", "Cost", "Support", "Delivery"}
	radarChart := pptx.NewRadarChart(radarCategories, []float64{80, 90, 70, 85, 75}).
		WithSeriesName("Product A").
		WithLegend(true)
	radarChart.Title = "Radar Chart"
	return pptx.NewSlide("Radar Chart").WithRadarChart(radarChart)
}

func buildLegendConstantsSlide() pptx.SlideContent {
	return pptx.NewSlide("Legend Position Constants").
		AddBullet(fmt.Sprintf("LegendPositionRight  = %q", pptx.LegendPositionRight)).
		AddBullet(fmt.Sprintf("LegendPositionLeft   = %q", pptx.LegendPositionLeft)).
		AddBullet(fmt.Sprintf("LegendPositionTop    = %q", pptx.LegendPositionTop)).
		AddBullet(fmt.Sprintf("LegendPositionBottom = %q", pptx.LegendPositionBottom))
}
