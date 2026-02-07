package pptx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestCreateWithSlidesBarLegendPositionAndDataLabels(t *testing.T) {
	chart := NewBarChart(
		[]string{"Q1", "Q2"},
		[]float64{3, 7},
	).
		WithLegend(true).
		WithLegendPosition("l").
		WithDataLabels(true)

	slides := []SlideContent{
		NewSlide("Styled").WithBarChart(chart),
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
		`<c:legendPos val="l"/>`,
		`<c:dLbls>`,
		`<c:showVal val="1"/>`,
	}
	for _, needle := range checks {
		if !strings.Contains(chartXML, needle) {
			t.Fatalf("expected %q in chart XML", needle)
		}
	}
}

func TestCreateWithSlidesLineLegendPositionTop(t *testing.T) {
	chart := NewLineChart(
		[]string{"W1", "W2", "W3"},
		[]float64{1, 2, 3},
	).
		WithLegend(true).
		WithLegendPosition("t")

	slides := []SlideContent{
		NewSlide("Styled").WithLineChart(chart),
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
	if !strings.Contains(chartXML, `<c:legendPos val="t"/>`) {
		t.Fatalf("expected legend top position in chart XML")
	}
}

func TestCreateWithSlidesRejectsInvalidLegendPosition(t *testing.T) {
	chart := NewBarChart(
		[]string{"A"},
		[]float64{1},
	).WithLegendPosition("center")

	slides := []SlideContent{
		NewSlide("Bad").WithBarChart(chart),
	}
	_, err := CreateWithSlides("Demo", slides)
	if err == nil {
		t.Fatalf("expected legend position validation error")
	}
}

func TestCreateWithSlidesChartTitleAndLegendOverlay(t *testing.T) {
	chart := NewBarChart(
		[]string{"Q1", "Q2"},
		[]float64{3, 7},
	).
		WithLegend(true).
		WithTitleOverlay(true).
		WithLegendOverlay(true)

	slides := []SlideContent{
		NewSlide("Overlay").WithBarChart(chart),
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
	if strings.Count(chartXML, `<c:overlay val="1"/>`) != 2 {
		t.Fatalf("expected title and legend overlay values to be enabled")
	}
}
