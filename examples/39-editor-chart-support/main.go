package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
)

func main() {
	const outputDir = "examples/output"
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	tmpDir, tempErr := os.MkdirTemp("", "gopptx-example-39-*")
	if tempErr != nil {
		log.Fatalf("Failed to create temp directory: %v", tempErr)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	inputFile := filepath.Join(tmpDir, "39_editor_chart_support_input.pptx")
	outputFile := filepath.Join(outputDir, "39_editor_chart_support_output.pptx")

	// 1. Create a minimal valid PPTX file
	log.Printf("Generating minimal PPTX: %s...\n", inputFile)
	if err := createMinimalPPTX(inputFile); err != nil {
		log.Fatalf("Failed to create minimal PPTX: %v", err)
	}

	// 2. Open it
	log.Printf("Opening %s...\n", inputFile)
	ppt, err := editor.OpenPresentationEditor(inputFile)
	if err != nil {
		log.Fatalf("Failed to open presentation: %v", err)
	}
	defer func() { _ = ppt.Close() }()

	if ppt.SlideCount() == 0 {
		log.Fatal("Input presentation has no slides")
	}

	// 3. Add a Bar Chart to Slide 1
	log.Println("Adding Bar Chart to Slide 1...")
	barChart := charts.NewBarChart(
		[]string{"Q1", "Q2", "Q3", "Q4"},
		[]float64{100, 200, 150, 300},
	).WithTitle("Quarterly Sales")

	if err := ppt.AddChart(0, barChart); err != nil {
		log.Fatalf("Failed to add bar chart: %v", err)
	}

	// 4. Add a Line Chart to Slide 1
	log.Println("Adding Line Chart to Slide 1...")
	lineChart := charts.NewLineChart(
		[]string{"Jan", "Feb", "Mar"},
		[]float64{5, 10, 8},
	).WithTitle("Monthly Growth")

	// Offset it
	lineChart = lineChart.Position(914400*5, 1800000)

	if err := ppt.AddChart(0, lineChart); err != nil {
		log.Fatalf("Failed to add line chart: %v", err)
	}

	// 5. Save
	log.Printf("Saving to %s...\n", outputFile)
	if err := ppt.Save(outputFile); err != nil {
		log.Fatalf("Failed to save: %v", err)
	}

	log.Println("Done! Smoke test passed.")
}

func createMinimalPPTX(filename string) error {
	data, err := pptx.CreateWithSlides("Editor Chart Support", []pptx.SlideContent{
		pptx.NewSlide("Chart Playground"),
	})
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0o644)
}
