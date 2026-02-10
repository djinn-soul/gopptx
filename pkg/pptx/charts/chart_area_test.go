package charts_test

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

func TestCreateWithSlidesEmbedsAreaChart(t *testing.T) {
	chart := charts.NewAreaChart(
		[]string{"Jan", "Feb", "Mar"},
		[]float64{0.2, 0.35, 0.5},
	).
		WithTitle("Adoption").
		WithAreaColor("92D050").
		WithSeriesName("Conversion").
		WithLegend(true).
		WithLegendPosition("l").
		WithDataLabels(true).
		WithAxisTitles("Month", "Rate").
		WithValueFormat("0.0%").
		WithValueRange(0, 1)

	slides := []pptx.SlideContent{
		pptx.NewSlide("Area").WithAreaChart(chart),
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
		`<c:areaChart>`,
		`<c:grouping val="standard"/>`,
		`<a:srgbClr val="92D050"/>`,
		`<c:v>Conversion</c:v>`,
		`<c:legendPos val="l"/>`,
		`<c:showVal val="1"/>`,
		`<a:t>Month</a:t>`,
		`<a:t>Rate</a:t>`,
		`<c:numFmt formatCode="0.0%" sourceLinked="0"/>`,
	}
	for _, needle := range checks {
		if !strings.Contains(chartXML, needle) {
			t.Fatalf("expected %q in chart XML", needle)
		}
	}
}

func TestCreateWithSlidesRejectsInvalidAreaColor(t *testing.T) {
	chart := charts.NewAreaChart(
		[]string{"A"},
		[]float64{1},
	).WithAreaColor("GHIJKL")

	slides := []pptx.SlideContent{
		pptx.NewSlide("Bad").WithAreaChart(chart),
	}

	_, err := pptx.CreateWithSlides("Demo", slides)
	if err == nil {
		t.Fatalf("expected area chart color validation error")
	}
}
