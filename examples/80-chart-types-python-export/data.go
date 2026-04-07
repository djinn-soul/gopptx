package main

import "github.com/djinn-soul/gopptx/pkg/pptx/styling"

type chartDemoData struct {
	chartX  styling.Length
	chartY  styling.Length
	chartW  styling.Length
	chartH  styling.Length
	cats    []string
	vals    []float64
	hlcCats []string
	high    []float64
	low     []float64
	closeV  []float64
	openV   []float64
	xVals   []float64
	yVals   []float64
	sizes   []float64
}

func newChartDemoData() chartDemoData {
	return chartDemoData{
		chartX:  styling.Inches(0.8),
		chartY:  styling.Inches(1.4),
		chartW:  styling.Inches(8.5),
		chartH:  styling.Inches(4.6),
		cats:    []string{"Q1", "Q2", "Q3", "Q4"},
		vals:    []float64{14, 21, 18, 27},
		hlcCats: []string{"D1", "D2", "D3", "D4"},
		high:    []float64{16, 23, 20, 29},
		low:     []float64{12, 18, 15, 24},
		closeV:  []float64{14, 21, 18, 27},
		openV:   []float64{13, 20, 17, 26},
		xVals:   []float64{1, 2, 3, 4},
		yVals:   []float64{14, 21, 18, 27},
		sizes:   []float64{14, 21, 18, 27},
	}
}
