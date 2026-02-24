package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

const outputDir = "examples/output"

func main() {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		fail("create output directory", err)
	}

	samples := map[string]func() ([]byte, error){
		"09_charts_family1_bar_line_area_variants.pptx": createFamily1Sample,
		"09_charts_family2_bubble.pptx":                 createFamily2Sample,
		"09_charts_family3_radar.pptx":                  createFamily3Sample,
		"09_charts_family4_stock_combo.pptx":            createFamily4Sample,
	}

	log.Println("Generating chart smoke samples...")
	for name, build := range samples {
		data, buildErr := build()
		if buildErr != nil {
			fail("build "+name, buildErr)
		}

		path := filepath.Join(outputDir, name)
		if err := os.WriteFile(path, data, 0o600); err != nil {
			fail("write "+path, err)
		}
		log.Printf("  wrote %s\n", path)
	}

	log.Println("Done.")
	log.Println("Manual checklist: scripts/chart_smoke_checklist.md")
}

func createFamily1Sample() ([]byte, error) {
	slides := []pptx.SlideContent{
		pptx.NewSlide("Bar Horizontal").WithBarHorizontalChart(
			pptx.NewBarHorizontalChart(
				[]string{"Q1", "Q2", "Q3"},
				[]float64{20, 30, 25},
			).WithTitle("BarHorizontal"),
		),
		pptx.NewSlide("Bar Stacked").WithBarStackedChart(
			pptx.NewBarStackedChart(
				[]string{"Q1", "Q2", "Q3"},
				[]float64{10, 14, 12},
			).WithTitle("BarStacked"),
		),
		pptx.NewSlide("Bar Stacked 100").WithBarStacked100Chart(
			pptx.NewBarStacked100Chart(
				[]string{"Q1", "Q2", "Q3"},
				[]float64{40, 55, 65},
			).WithTitle("BarStacked100"),
		),
		pptx.NewSlide("Line Markers").WithLineMarkersChart(
			pptx.NewLineMarkersChart(
				[]string{"Jan", "Feb", "Mar"},
				[]float64{8, 11, 10},
			).WithTitle("LineMarkers"),
		),
		pptx.NewSlide("Line Stacked").WithLineStackedChart(
			pptx.NewLineStackedChart(
				[]string{"Jan", "Feb", "Mar"},
				[]float64{5, 9, 7},
			).WithTitle("LineStacked"),
		),
		pptx.NewSlide("Area Stacked").WithAreaStackedChart(
			pptx.NewAreaStackedChart(
				[]string{"W1", "W2", "W3"},
				[]float64{6, 7, 8},
			).WithTitle("AreaStacked"),
		),
		pptx.NewSlide("Area Stacked 100").WithAreaStacked100Chart(
			pptx.NewAreaStacked100Chart(
				[]string{"W1", "W2", "W3"},
				[]float64{35, 50, 45},
			).WithTitle("AreaStacked100"),
		),
	}
	return pptx.CreateWithSlides("Family 1 Smoke", slides)
}

func createFamily2Sample() ([]byte, error) {
	chart := pptx.NewBubbleChart(
		[]float64{1, 2, 3, 4},
		[]float64{4, 3, 6, 5},
		[]float64{15, 25, 20, 35},
	).
		WithTitle("Bubble").
		WithBubbleScale(110).
		WithLegend(true)
	slides := []pptx.SlideContent{
		pptx.NewSlide("Bubble").WithBubbleChart(chart),
	}
	return pptx.CreateWithSlides("Family 2 Smoke", slides)
}

func createFamily3Sample() ([]byte, error) {
	slides := []pptx.SlideContent{
		pptx.NewSlide("Radar").WithRadarChart(
			pptx.NewRadarChart(
				[]string{"Speed", "Quality", "Cost", "Risk"},
				[]float64{7, 8, 5, 6},
			).WithTitle("Radar"),
		),
		pptx.NewSlide("Radar Filled").WithRadarFilledChart(
			pptx.NewRadarFilledChart(
				[]string{"Speed", "Quality", "Cost", "Risk"},
				[]float64{6, 9, 4, 7},
			).WithTitle("RadarFilled"),
		),
	}
	return pptx.CreateWithSlides("Family 3 Smoke", slides)
}

func createFamily4Sample() ([]byte, error) {
	slides := []pptx.SlideContent{
		pptx.NewSlide("Stock HLC").WithStockHLCChart(
			pptx.NewStockHLCChart(
				[]string{"D1", "D2", "D3"},
				[]float64{12, 13, 14},
				[]float64{8, 9, 10},
				[]float64{10, 11, 12},
			).WithTitle("StockHLC"),
		),
		pptx.NewSlide("Stock OHLC").WithStockOHLCChart(
			pptx.NewStockOHLCChart(
				[]string{"D1", "D2", "D3"},
				[]float64{9, 10, 11},
				[]float64{12, 13, 14},
				[]float64{8, 9, 10},
				[]float64{10, 11, 12},
			).WithTitle("StockOHLC"),
		),
		pptx.NewSlide("Combo").WithComboChart(
			pptx.NewComboChart(
				[]string{"Q1", "Q2", "Q3"},
				[]pptx.Series{
					{Name: "Bar A", Values: []float64{10, 12, 11}},
				},
				[]pptx.Series{
					{Name: "Line A", Values: []float64{9, 13, 12}},
				},
			).WithTitle("Combo"),
		),
	}
	return pptx.CreateWithSlides("Family 4 Smoke", slides)
}

func fail(step string, err error) {
	fmt.Fprintf(os.Stderr, "error: %s: %v\n", step, err)
	os.Exit(1)
}
