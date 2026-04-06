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
)

const outputDir = "examples/output"

func main() {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		fmt.Fprintf(os.Stderr, "error: create output dir: %v\n", err)
		os.Exit(1)
	}

	data := newChartDemoData()
	slides := buildSlides(data)

	pptxPath := outputDir + "/80_chart_types_go_export.pptx"
	pdfPath := outputDir + "/80_chart_types_go_export.pdf"
	if err := writePresentation(pptxPath, pdfPath, slides); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
