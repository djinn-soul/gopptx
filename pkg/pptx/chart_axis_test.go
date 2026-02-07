package pptx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestCreateWithSlidesLineChartValueFormatAndRange(t *testing.T) {
	chart := NewLineChart(
		[]string{"P1", "P2", "P3"},
		[]float64{10.25, 20.5, 15.75},
	).
		WithValueFormat("#,##0.00").
		WithValueRange(0, 100)

	slides := []SlideContent{
		NewSlide("Line").WithLineChart(chart),
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
	chart := NewBarChart(
		[]string{"A", "B"},
		[]float64{1, 2},
	).WithValueRange(10, 10)

	slides := []SlideContent{
		NewSlide("Bad").WithBarChart(chart),
	}

	_, err := CreateWithSlides("Demo", slides)
	if err == nil {
		t.Fatalf("expected value-range validation error")
	}
}

func TestCreateWithSlidesRejectsEmptyValueFormat(t *testing.T) {
	chart := NewBarChart(
		[]string{"A"},
		[]float64{1},
	).WithValueFormat("   ")

	slides := []SlideContent{
		NewSlide("Bad").WithBarChart(chart),
	}

	_, err := CreateWithSlides("Demo", slides)
	if err == nil {
		t.Fatalf("expected value-format validation error")
	}
}

func TestCreateWithSlidesValueAxisCrossBetweenMidCategory(t *testing.T) {
	chart := NewBarChart(
		[]string{"P1", "P2", "P3"},
		[]float64{1, 2, 3},
	).WithValueAxisCrossBetween(ValueAxisCrossBetweenMidCategory)

	slides := []SlideContent{
		NewSlide("Bar").WithBarChart(chart),
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
	if !strings.Contains(chartXML, `<c:crossBetween val="midCat"/>`) {
		t.Fatalf("expected midCat crossBetween in chart XML")
	}
}

func TestCreateWithSlidesRejectsInvalidValueAxisCrossBetween(t *testing.T) {
	chart := NewBarChart(
		[]string{"A"},
		[]float64{1},
	).WithValueAxisCrossBetween("middle")

	slides := []SlideContent{
		NewSlide("Bad").WithBarChart(chart),
	}

	_, err := CreateWithSlides("Demo", slides)
	if err == nil {
		t.Fatalf("expected crossBetween validation error")
	}
}
