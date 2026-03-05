package charts_test

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

func TestCreateWithSlidesLineChartValueFormatAndRange(t *testing.T) {
	chart := charts.NewLineChart(
		[]string{"P1", "P2", "P3"},
		[]float64{10.25, 20.5, 15.75},
	).
		WithValueFormat("#,##0.00").
		WithValueRange(0, 100)

	slides := []pptx.SlideContent{
		pptx.NewSlide("Line").WithLineChart(chart),
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
		`<c:numFmt formatCode="#,##0.00" sourceLinked="0"/>`,
		`<c:min val="0.000000"/>`,
		`<c:max val="100.000000"/>`,
	}
	for _, needle := range checks {
		if !strings.Contains(chartXML, needle) {
			t.Fatalf("expected %q in chart XML", needle)
		}
	}
}

func TestCreateWithSlidesRejectsInvalidValueRange(t *testing.T) {
	chart := charts.NewBarChart(
		[]string{"A", "B"},
		[]float64{1, 2},
	).WithValueRange(10, 10)

	slides := []pptx.SlideContent{
		pptx.NewSlide("Bad").WithBarChart(chart),
	}

	_, err := pptx.CreateWithSlides("Demo", slides)
	if err == nil {
		t.Fatalf("expected value-range validation error")
	}
}

func TestCreateWithSlidesRejectsEmptyValueFormat(t *testing.T) {
	chart := charts.NewBarChart(
		[]string{"A"},
		[]float64{1},
	).WithValueFormat("   ")

	slides := []pptx.SlideContent{
		pptx.NewSlide("Bad").WithBarChart(chart),
	}

	_, err := pptx.CreateWithSlides("Demo", slides)
	if err == nil {
		t.Fatalf("expected value-format validation error")
	}
}

func TestCreateWithSlidesValueAxisCrossBetweenMidCategory(t *testing.T) {
	chart := charts.NewBarChart(
		[]string{"P1", "P2", "P3"},
		[]float64{1, 2, 3},
	).WithValueAxisCrossBetween(charts.ValueAxisCrossBetweenMidCategory)

	slides := []pptx.SlideContent{
		pptx.NewSlide("Bar").WithBarChart(chart),
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
	if !strings.Contains(chartXML, `<c:crossBetween val="midCat"/>`) {
		t.Fatalf("expected midCat crossBetween in chart XML")
	}
}

func TestCreateWithSlidesRejectsInvalidValueAxisCrossBetween(t *testing.T) {
	chart := charts.NewBarChart(
		[]string{"A"},
		[]float64{1},
	).WithValueAxisCrossBetween("middle")

	slides := []pptx.SlideContent{
		pptx.NewSlide("Bad").WithBarChart(chart),
	}

	_, err := pptx.CreateWithSlides("Demo", slides)
	if err == nil {
		t.Fatalf("expected crossBetween validation error")
	}
}

func TestCreateWithSlidesAxisControlsRender(t *testing.T) {
	chart := charts.NewBarChart(
		[]string{"A", "B"},
		[]float64{1, 2},
	).
		WithCategoryMajorGridlines(true).
		WithTickLabelPositions(charts.AxisTickLabelPositionLow, charts.AxisTickLabelPositionHigh).
		WithAxisCrosses(charts.AxisCrossesMin, charts.AxisCrossesMax)

	slides := []pptx.SlideContent{
		pptx.NewSlide("Bar").WithBarChart(chart),
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
		`<c:tickLblPos val="low"/>`,
		`<c:tickLblPos val="high"/>`,
		`<c:crosses val="min"/>`,
		`<c:crosses val="max"/>`,
	}
	for _, needle := range checks {
		if !strings.Contains(chartXML, needle) {
			t.Fatalf("expected %q in chart XML", needle)
		}
	}
}

func TestCreateWithSlidesRejectsInvalidAxisTickLabelPosition(t *testing.T) {
	chart := charts.NewLineChart(
		[]string{"A"},
		[]float64{1},
	).WithTickLabelPositions("side", charts.AxisTickLabelPositionNextTo)

	slides := []pptx.SlideContent{
		pptx.NewSlide("Bad").WithLineChart(chart),
	}

	_, err := pptx.CreateWithSlides("Demo", slides)
	if err == nil {
		t.Fatalf("expected tick label position validation error")
	}
}

func TestCreateWithSlidesRejectsInvalidAxisCrosses(t *testing.T) {
	chart := charts.NewAreaChart(
		[]string{"A"},
		[]float64{1},
	).WithAxisCrosses("center", charts.AxisCrossesAutoZero)

	slides := []pptx.SlideContent{
		pptx.NewSlide("Bad").WithAreaChart(chart),
	}

	_, err := pptx.CreateWithSlides("Demo", slides)
	if err == nil {
		t.Fatalf("expected axis crosses validation error")
	}
}
