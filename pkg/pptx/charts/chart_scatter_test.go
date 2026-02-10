package charts_test

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

func TestCreateWithSlidesEmbedsScatterChart(t *testing.T) {
	chart := charts.NewScatterChart(
		[]float64{1, 2, 3},
		[]float64{2.5, 3.0, 4.2},
	).
		WithTitle("Experiment").
		WithLineColor("1F4E78").
		WithSeriesName("Run A").
		WithScatterStyle(charts.ScatterStyleSmoothMarker).
		WithLegend(true).
		WithLegendPosition("t").
		WithDataLabels(true).
		WithAxisTitles("Time", "Value").
		WithValueFormat("0.00").
		WithValueRange(0, 5)

	slides := []pptx.SlideContent{
		pptx.NewSlide("Scatter").WithScatterChart(chart),
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
		`<c:scatterChart>`,
		`<c:scatterStyle val="smoothMarker"/>`,
		`<c:xVal><c:numLit>`,
		`<c:yVal><c:numLit>`,
		`<a:srgbClr val="1F4E78"/>`,
		`<c:v>Run A</c:v>`,
		`<c:legendPos val="t"/>`,
		`<c:numFmt formatCode="0.00" sourceLinked="0"/>`,
		`<c:min val="0.000000"/>`,
		`<c:max val="5.000000"/>`,
	}
	for _, needle := range checks {
		if !strings.Contains(chartXML, needle) {
			t.Fatalf("expected %q in chart XML", needle)
		}
	}
}

func TestCreateWithSlidesRejectsScatterLengthMismatch(t *testing.T) {
	chart := charts.NewScatterChart(
		[]float64{1, 2},
		[]float64{3},
	)
	slides := []pptx.SlideContent{
		pptx.NewSlide("Bad").WithScatterChart(chart),
	}

	_, err := pptx.CreateWithSlides("Demo", slides)
	if err == nil {
		t.Fatalf("expected scatter chart length validation error")
	}
}

func TestCreateWithSlidesRejectsInvalidScatterStyle(t *testing.T) {
	chart := charts.NewScatterChart(
		[]float64{1},
		[]float64{2},
	).WithScatterStyle("curved")
	slides := []pptx.SlideContent{
		pptx.NewSlide("Bad").WithScatterChart(chart),
	}

	_, err := pptx.CreateWithSlides("Demo", slides)
	if err == nil {
		t.Fatalf("expected scatter chart style validation error")
	}
}
