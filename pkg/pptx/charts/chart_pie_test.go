package charts_test

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

func TestCreateWithSlidesEmbedsPieChart(t *testing.T) {
	chart := charts.NewPieChart(
		[]string{"Product A", "Product B", "Product C"},
		[]float64{25, 35, 40},
	).
		WithTitle("Market Share").
		WithSeriesName("FY2026").
		WithLegend(true).
		WithLegendPosition("b").
		WithDataLabels(true)

	slides := []pptx.SlideContent{
		pptx.NewSlide("Pie").WithPieChart(chart),
	}

	data, err := pptx.CreateWithSlides("Demo", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	chartXML := readZipFile(t, zr, "ppt/charts/chart1.xml")
	checks := []string{
		`<c:pieChart>`,
		`<c:varyColors val="1"/>`,
		`Market Share`,
		`<c:v>Product A</c:v>`,
		`<c:v>40.000000</c:v>`,
		`<c:v>FY2026</c:v>`,
		`<c:legendPos val="b"/>`,
		`<c:dLbls>`,
		`<c:showPercent val="1"/>`,
	}
	for _, needle := range checks {
		if !strings.Contains(chartXML, needle) {
			t.Fatalf("expected %q in chart XML", needle)
		}
	}
}

func TestCreateWithSlidesRejectsInvalidPieChart(t *testing.T) {
	chart := charts.NewPieChart(
		[]string{"A", "B"},
		[]float64{10},
	)

	slides := []pptx.SlideContent{
		pptx.NewSlide("Bad").WithPieChart(chart),
	}

	_, err := pptx.CreateWithSlides("Demo", slides)
	if err == nil {
		t.Fatalf("expected pie chart validation error")
	}
}

func TestCreateWithSlidesRejectsEmptyPieSeriesName(t *testing.T) {
	chart := charts.NewPieChart(
		[]string{"A", "B"},
		[]float64{1, 2},
	).WithSeriesName(" ")

	slides := []pptx.SlideContent{
		pptx.NewSlide("Bad").WithPieChart(chart),
	}

	_, err := pptx.CreateWithSlides("Demo", slides)
	if err == nil {
		t.Fatalf("expected pie series-name validation error")
	}
}
