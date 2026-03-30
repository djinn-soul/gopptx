// examples/68-chart-data-api demonstrates chart data builders.
//
// Shows how to create category chart data (bar, line, pie), XY chart data
// (scatter), and bubble chart data using the data-builder pattern.
// Also shows updating chart data on an existing presentation via the Presentation API.
//
// Run with: go run ./examples/68-chart-data-api/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

const (
	outputDir  = "examples/output"
	outputFile = "68_chart_data_api.pptx"
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

	outputPath := filepath.Join(outputDir, outputFile)
	if err := buildChartPresentation(outputPath); err != nil {
		return err
	}
	if err := updateChartData(outputPath); err != nil {
		return err
	}

	log.Printf("Generated %s\n", outputPath)
	return nil
}

func buildChartPresentation(outputPath string) error {
	categories := []string{"Q1", "Q2", "Q3", "Q4"}

	slides := []pptx.SlideContent{
		buildBarChartSlide(categories),
		buildLineChartSlide(categories),
		buildPieChartSlide(),
		buildHBarChartSlide(categories),
		buildDoughnutChartSlide(),
		buildMultiSeriesSlide(categories),
		buildXyDataSlide(),
		buildBubbleDataSlide(),
	}

	data, err := pptx.CreateWithSlides("Chart Data API Demo", slides)
	if err != nil {
		return fmt.Errorf("create presentation: %w", err)
	}
	if err := os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

func buildBarChartSlide(categories []string) pptx.SlideContent {
	barChart := pptx.NewBarChart(categories, []float64{42, 55, 61, 73}).
		WithSeriesName("Revenue ($k)").
		WithLegend(true).
		WithLegendPosition(pptx.LegendPositionBottom).
		WithDataLabels(true).
		WithMajorGridlines(true).
		WithAxisTitles("Quarter", "Revenue ($k)").
		WithValueFormat("$#,##0").
		WithValueRange(0, 100)
	barChart.Title = "Quarterly Revenue – Bar Chart"
	return pptx.NewSlide("Bar Chart").WithBarChart(barChart)
}

func buildLineChartSlide(categories []string) pptx.SlideContent {
	lineChart := pptx.NewLineChart(categories, []float64{30, 45, 55, 68}).
		WithSeriesName("Units Sold").
		WithLegend(true).
		WithSmooth(true).
		WithDataLabels(true).
		WithAxisTitles("Quarter", "Units")
	lineChart.Title = "Quarterly Sales – Line Chart"
	return pptx.NewSlide("Line Chart").WithLineChart(lineChart)
}

func buildPieChartSlide() pptx.SlideContent {
	pieCategories := []string{"Product A", "Product B", "Product C", "Product D"}
	pieChart := pptx.NewPieChart(pieCategories, []float64{35, 25, 20, 20}).
		WithSeriesName("Market Share").
		WithLegend(true).
		WithLegendPosition(pptx.LegendPositionRight).
		WithDataLabels(true)
	pieChart.Title = "Market Share – Pie Chart"
	return pptx.NewSlide("Pie Chart").WithPieChart(pieChart)
}

func buildHBarChartSlide(categories []string) pptx.SlideContent {
	hBarChart := pptx.NewBarHorizontalChart(categories, []float64{18, 24, 30, 41}).
		WithSeriesName("Costs ($k)").
		WithLegend(false).
		WithAxisTitles("", "Cost ($k)")
	hBarChart.Title = "Quarterly Costs – Horizontal Bar"
	return pptx.NewSlide("Horizontal Bar Chart").WithBarHorizontalChart(hBarChart)
}

func buildDoughnutChartSlide() pptx.SlideContent {
	pieCategories := []string{"Product A", "Product B", "Product C", "Product D"}
	doughnut := pptx.NewDoughnutChart(pieCategories, []float64{40, 30, 20, 10}).
		WithSeriesName("Revenue Share").
		WithLegend(true).
		WithHoleSize(50).
		WithDataLabels(true)
	doughnut.Title = "Revenue Share – Doughnut Chart"
	return pptx.NewSlide("Doughnut Chart").WithDoughnutChart(doughnut)
}

func buildMultiSeriesSlide(categories []string) pptx.SlideContent {
	multiCatData := pptx.NewCategoryChartData(categories).
		AddSeries("Product A", []float64{10, 15, 20, 25}).
		AddSeries("Product B", []float64{8, 12, 18, 22}).
		AddSeries("Product C", []float64{5, 8, 11, 15})

	return pptx.NewSlide("CategoryChartData Builder").
		AddBullet("CategoryChartData supports multi-series bar/line charts.").
		AddBullet(fmt.Sprintf("Categories: %v", multiCatData.Categories)).
		AddBullet(fmt.Sprintf("Series count: %d", len(multiCatData.Series)))
}

func buildXyDataSlide() pptx.SlideContent {
	xyData := pptx.NewXyChartData().
		AddSeries("Series A", []float64{1, 2, 3, 4, 5}, []float64{2.5, 4.0, 3.5, 5.0, 6.5}).
		AddSeries("Series B", []float64{1, 2, 3, 4, 5}, []float64{1.0, 2.5, 4.0, 3.0, 5.5})

	return pptx.NewSlide("XyChartData Builder").
		AddBullet("XyChartData is used for scatter charts.").
		AddBullet(fmt.Sprintf("Series count: %d", len(xyData.Series)))
}

func buildBubbleDataSlide() pptx.SlideContent {
	bubbleData := pptx.NewBubbleChartData().
		AddSeries(
			"Bubble Series A",
			[]float64{1, 2, 3, 4},
			[]float64{3, 1, 4, 2},
			[]float64{10, 20, 15, 25},
		)

	return pptx.NewSlide("BubbleChartData Builder").
		AddBullet("BubbleChartData is used for bubble charts.").
		AddBullet(fmt.Sprintf("Series count: %d", len(bubbleData.Series)))
}

func updateChartData(outputPath string) error {
	categories := []string{"Q1", "Q2", "Q3", "Q4"}

	// --- Part 2: Show UpdateChartDataByIndexFromBuilder using Presentation API ---
	prs, err := pptx.Open(outputPath)
	if err != nil {
		return fmt.Errorf("open for chart update: %w", err)
	}
	defer prs.Close()

	// Update the bar chart (slide 0, chart 0) with new data via builder.
	updatedBarData := pptx.NewCategoryChartData(categories).
		AddSeries("Updated Revenue", []float64{50, 62, 70, 85})

	if err := prs.UpdateChartDataByIndexFromBuilder(0, 0, updatedBarData); err != nil {
		// Non-fatal: bar chart may not have an embedded Excel workbook in this simple build.
		log.Printf(
			"Note: UpdateChartDataByIndexFromBuilder: %v (expected if no xlsx embed)\n", err,
		)
	}

	if saveErr := prs.Save(); saveErr != nil {
		return fmt.Errorf("save after chart update: %w", saveErr)
	}
	return nil
}
