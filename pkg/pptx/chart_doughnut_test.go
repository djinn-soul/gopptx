package pptx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestCreateWithSlidesEmbedsDoughnutChart(t *testing.T) {
	chart := NewDoughnutChart(
		[]string{"North", "South", "West"},
		[]float64{30, 45, 25},
	).
		WithTitle("Regional Mix").
		WithSeriesName("FY2026").
		WithLegend(true).
		WithLegendPosition("t").
		WithDataLabels(true).
		WithHoleSize(60)

	slides := []SlideContent{
		NewSlide("Doughnut").WithDoughnutChart(chart),
	}

	data, err := CreateWithSlides("Demo", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	chartXML := readZipFile(t, zr, "ppt/charts/chart1.xml")
	checks := []string{
		`<c:doughnutChart>`,
		`<c:varyColors val="1"/>`,
		`<c:holeSize val="60"/>`,
		`Regional Mix`,
		`<c:v>North</c:v>`,
		`<c:v>45.000000</c:v>`,
		`<c:v>FY2026</c:v>`,
		`<c:legendPos val="t"/>`,
		`<c:dLbls>`,
		`<c:showPercent val="1"/>`,
	}
	for _, needle := range checks {
		if !strings.Contains(chartXML, needle) {
			t.Fatalf("expected %q in chart XML", needle)
		}
	}
}

func TestCreateWithSlidesRejectsInvalidDoughnutHoleSize(t *testing.T) {
	chart := NewDoughnutChart(
		[]string{"A", "B"},
		[]float64{1, 2},
	).WithHoleSize(5)

	slides := []SlideContent{
		NewSlide("Bad").WithDoughnutChart(chart),
	}

	_, err := CreateWithSlides("Demo", slides)
	if err == nil {
		t.Fatalf("expected doughnut chart hole-size validation error")
	}
}
