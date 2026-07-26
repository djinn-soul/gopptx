// examples/16-chart-type-enum/main.go demonstrates the chart type enum API.
//
// pptx.NewChart takes an enums.XLChartType, so a chart type chosen at runtime -
// from config, a CLI flag, user input - does not need a switch over the
// nineteen NewBarChart / NewPieChart / ... constructors. The result satisfies
// pptx.Chart and goes onto a slide through SlideContent.WithChart.
//
// Run with: go run ./examples/16-chart-type-enum/main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/enums"
)

const (
	outputDir  = "examples/output"
	outputFile = "16_chart_type_enum.pptx"
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

	categories := []string{"Q1", "Q2", "Q3", "Q4"}
	values := []float64{12, 18, 15, 24}

	// The same data rendered as every category/value chart type, selected by
	// enum rather than by calling a different constructor for each one.
	chartTypes := []enums.XLChartType{
		enums.XLChartTypeBar,
		enums.XLChartTypeBarHorizontal,
		enums.XLChartTypeBarStacked,
		enums.XLChartTypeLine,
		enums.XLChartTypeLineMarkers,
		enums.XLChartTypeArea,
		enums.XLChartTypePie,
		enums.XLChartTypeDoughnut,
		enums.XLChartTypeRadar,
	}

	slides := make([]pptx.SlideContent, 0, len(chartTypes))
	for _, chartType := range chartTypes {
		chart, err := pptx.NewChart(chartType, categories, values)
		if err != nil {
			return fmt.Errorf("build %s chart: %w", chartType, err)
		}
		slides = append(slides, pptx.NewSlide("Revenue: "+string(chartType)).WithChart(chart))
	}

	// Charts whose data is not categories and values are rejected by NewChart,
	// with an error naming the constructor that takes the right data shape.
	if _, err := pptx.NewChart(enums.XLChartTypeScatter, categories, values); err != nil {
		log.Printf("scatter rejected as expected: %v\n", err)
	}

	// Those types are still reachable through their own constructors, and the
	// result goes onto a slide through the same WithChart call.
	scatter := pptx.NewScatterChart([]float64{1, 2, 3, 4}, []float64{4, 9, 6, 12}).
		WithTitle("Scatter")
	slides = append(slides, pptx.NewSlide("Revenue: scatter").WithChart(scatter))

	data, err := pptx.CreateWithSlides("Chart Type Enum Demo", slides)
	if err != nil {
		return fmt.Errorf("create presentation: %w", err)
	}

	outputPath := filepath.Join(outputDir, outputFile)
	if err := os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	log.Printf("Generated %s with %d chart slides\n", outputPath, len(slides))
	return nil
}
