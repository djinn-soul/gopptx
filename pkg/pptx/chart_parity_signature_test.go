package pptx

import (
	"strings"
	"testing"
)

func TestBubbleChartParitySignature(t *testing.T) {
	chart := NewBubbleChart(
		[]float64{1, 2, 3},
		[]float64{4, 5, 6},
		[]float64{10, 20, 30},
	).
		WithTitle("BubbleParity").
		WithSeriesName("Series A").
		WithBubbleScale(130)

	xml := chartXMLForSlide(t, NewSlide("S").WithBubbleChart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:bubbleChart>`,
		`<c:varyColors val="0"/>`,
		`<c:bubbleScale val="130"/>`,
		`<c:tx><c:v>Series A</c:v></c:tx>`,
		`<c:xVal><c:numLit>`,
		`<c:yVal><c:numLit>`,
		`<c:bubbleSize><c:numLit>`,
		`<c:axId val="48650112"/>`,
		`<c:axId val="48672768"/>`,
	})
}

func TestRadarChartParitySignature(t *testing.T) {
	chart := NewRadarChart(
		[]string{"A", "B", "C"},
		[]float64{2, 3, 4},
	).WithTitle("RadarParity")

	xml := chartXMLForSlide(t, NewSlide("S").WithRadarChart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:radarChart>`,
		`<c:radarStyle val="marker"/>`,
		`<c:ptCount val="3"/>`,
		`<c:axId val="48650112"/>`,
		`<c:axId val="48672768"/>`,
	})
}

func TestRadarFilledChartParitySignature(t *testing.T) {
	chart := NewRadarFilledChart(
		[]string{"A", "B", "C"},
		[]float64{3, 4, 5},
	).WithTitle("RadarFilledParity")

	xml := chartXMLForSlide(t, NewSlide("S").WithRadarFilledChart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:radarChart>`,
		`<c:radarStyle val="filled"/>`,
		`<c:ptCount val="3"/>`,
	})
}

func TestStockHLCParitySignature(t *testing.T) {
	chart := NewStockHLCChart(
		[]string{"D1", "D2", "D3"},
		[]float64{12, 13, 14},
		[]float64{8, 9, 10},
		[]float64{10, 11, 12},
	)
	xml := chartXMLForSlide(t, NewSlide("S").WithStockHLCChart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:stockChart>`,
		`<c:tx><c:v>High</c:v></c:tx>`,
		`<c:tx><c:v>Low</c:v></c:tx>`,
		`<c:tx><c:v>Close</c:v></c:tx>`,
		`<c:hiLowLines>`,
	})
	if strings.Contains(xml, `<c:upDownBars>`) {
		t.Fatalf("did not expect up/down bars for HLC")
	}
}

func TestStockOHLCParitySignature(t *testing.T) {
	chart := NewStockOHLCChart(
		[]string{"D1", "D2", "D3"},
		[]float64{9, 10, 11},
		[]float64{12, 13, 14},
		[]float64{8, 9, 10},
		[]float64{10, 11, 12},
	)
	xml := chartXMLForSlide(t, NewSlide("S").WithStockOHLCChart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:stockChart>`,
		`<c:tx><c:v>Open</c:v></c:tx>`,
		`<c:tx><c:v>High</c:v></c:tx>`,
		`<c:tx><c:v>Low</c:v></c:tx>`,
		`<c:tx><c:v>Close</c:v></c:tx>`,
		`<c:hiLowLines>`,
		`<c:upDownBars>`,
	})
}

func TestComboChartParitySignature(t *testing.T) {
	chart := NewComboChart(
		[]string{"Q1", "Q2", "Q3"},
		[]Series{
			{Name: "Bar A", Values: []float64{1, 2, 3}},
			{Name: "Bar B", Values: []float64{4, 5, 6}},
		},
		[]Series{
			{Name: "Line A", Values: []float64{7, 8, 9}},
		},
	)

	xml := chartXMLForSlide(t, NewSlide("S").WithComboChart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:barChart>`,
		`<c:lineChart>`,
		`<c:grouping val="clustered"/>`,
		`<c:grouping val="standard"/>`,
		`<c:idx val="0"/>`,
		`<c:idx val="1"/>`,
		`<c:idx val="2"/>`,
		`<c:tx><c:v>Bar A</c:v></c:tx>`,
		`<c:tx><c:v>Bar B</c:v></c:tx>`,
		`<c:tx><c:v>Line A</c:v></c:tx>`,
	})
}
