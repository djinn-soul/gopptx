package pptxxml

import (
	"strings"
	"testing"
)

// The blank sentinel means "no reading", so the point is left out. Formatting it
// straight through wrote the literal NaN into <c:v> and PowerPoint had to repair
// the file.
func TestScatterChartOmitsBlankPoints(t *testing.T) {
	chart := &ChartSpec{
		SeriesName: "S",
		Color:      "4472C4",
		XValues:    []float64{1, BlankValue(), 3},
		Values:     []float64{10, 20, BlankValue()},
	}

	xml := chartScatterSeriesXML(chart)

	if strings.Contains(xml, "NaN") {
		t.Fatalf("a blank must not be written as NaN:\n%s", xml)
	}
	if strings.Count(xml, `<c:ptCount val="3"/>`) != 2 {
		t.Errorf("both caches should still count every index:\n%s", xml)
	}
	if strings.Contains(xml, `<c:pt idx="1"><c:v>`) && !strings.Contains(xml, "20.000000") {
		t.Error("the blank x point should be skipped while the y point at idx 1 stays")
	}
}

func TestBubbleChartOmitsBlankPoints(t *testing.T) {
	chart := &ChartSpec{
		SeriesName:  "S",
		Color:       "4472C4",
		XValues:     []float64{1, 2},
		Values:      []float64{10, BlankValue()},
		BubbleSizes: []float64{5, 6},
	}

	xml := chartBubbleSeriesXML(chart)

	if strings.Contains(xml, "NaN") {
		t.Fatalf("a blank must not be written as NaN:\n%s", xml)
	}
	if !strings.Contains(xml, "6.000000") {
		t.Errorf("the bubble size cache should still be written:\n%s", xml)
	}
}
