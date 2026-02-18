package charts_test

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

func TestCreateWithSlidesEmbedsBarChart(t *testing.T) {
	chart := charts.NewBarChart(
		[]string{"Q1", "Q2", "Q3"},
		[]float64{10, 22, 15},
	).WithTitle("Quarterly Revenue")

	slides := []pptx.SlideContent{
		pptx.NewSlide("Chart Slide").WithBarChart(chart),
	}

	data, err := pptx.CreateWithSlides("Demo", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	if !zipHasFile(zr, "ppt/charts/chart1.xml") {
		t.Fatalf("expected chart part file")
	}

	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")
	if !strings.Contains(slideXML, `c:chart`) || !strings.Contains(slideXML, `r:id="rId2"`) {
		t.Fatalf("expected chart frame reference in slide XML")
	}

	relsXML := readZipFile(t, zr, "ppt/slides/_rels/slide1.xml.rels")
	if !strings.Contains(relsXML, `relationships/chart"`) {
		t.Fatalf("expected chart relationship type in slide rels")
	}
	if !strings.Contains(relsXML, `Target="../charts/chart1.xml"`) {
		t.Fatalf("expected chart target in slide rels")
	}

	chartXML := readZipFile(t, zr, "ppt/charts/chart1.xml")
	checks := []string{`<c:barChart>`, `Quarterly Revenue`, `<c:v>Q1</c:v>`, `<c:v>22.000000</c:v>`}
	for _, needle := range checks {
		if !strings.Contains(chartXML, needle) {
			t.Fatalf("expected %q in chart XML", needle)
		}
	}
}

func TestCreateWithSlidesRejectsInvalidBarChart(t *testing.T) {
	chart := charts.NewBarChart(
		[]string{"Q1", "Q2"},
		[]float64{4},
	)
	slides := []pptx.SlideContent{
		pptx.NewSlide("Broken Chart").WithBarChart(chart),
	}

	_, err := pptx.CreateWithSlides("Demo", slides)
	if err == nil {
		t.Fatalf("expected chart validation error")
	}
}

func TestCreateWithSlidesEmbedsLineChart(t *testing.T) {
	chart := charts.NewLineChart(
		[]string{"Jan", "Feb", "Mar", "Apr"},
		[]float64{8, 12, 10, 16},
	).WithTitle("Monthly Trend")

	slides := []pptx.SlideContent{
		pptx.NewSlide("Line Chart Slide").WithLineChart(chart),
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
		`<c:lineChart>`,
		`Monthly Trend`,
		`<c:v>Jan</c:v>`,
		`<c:v>16.000000</c:v>`,
	}
	for _, needle := range checks {
		if !strings.Contains(chartXML, needle) {
			t.Fatalf("expected %q in chart XML", needle)
		}
	}
}

func TestCreateWithSlidesRejectsInvalidLineChart(t *testing.T) {
	chart := charts.NewLineChart(
		[]string{"A", "B", "C"},
		[]float64{0, 0, 0},
	)
	slides := []pptx.SlideContent{
		pptx.NewSlide("Broken Line").WithLineChart(chart),
	}

	_, err := pptx.CreateWithSlides("Demo", slides)
	if err == nil {
		t.Fatalf("expected line chart validation error")
	}
}

func TestCreateWithSlidesSupportsMultipleChartKindsOnOneSlide(t *testing.T) {
	bar := charts.NewBarChart([]string{"A"}, []float64{1})
	line := charts.NewLineChart([]string{"A"}, []float64{1})
	slides := []pptx.SlideContent{
		{
			Title: "Mixed",
			Chart: &bar,
			Line:  &line,
		},
	}

	_, err := pptx.CreateWithSlides("Demo", slides)
	if err != nil {
		t.Fatalf("unexpected validation error for multiple chart kinds on one slide: %v", err)
	}
}

func TestCreateWithSlidesImageAndChartRelationshipOrder(t *testing.T) {
	tmpDir := t.TempDir()
	imgPath := tmpDir + string(filepath.Separator) + "sample.png"
	if err := os.WriteFile(imgPath, testutil.TinyPNG(), 0o600); err != nil {
		t.Fatalf("write image: %v", err)
	}

	chart := charts.NewBarChart([]string{"A", "B"}, []float64{1, 2})
	slides := []pptx.SlideContent{
		pptx.NewSlide("Mixed").AddImage(pptx.NewImage(imgPath, 10, 10, 1000000, 1000000)).WithBarChart(chart),
	}
	data, err := pptx.CreateWithSlides("Demo", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	relsXML := readZipFile(t, zr, "ppt/slides/_rels/slide1.xml.rels")
	if !strings.Contains(
		relsXML,
		`Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image"`,
	) {
		t.Fatalf("expected image rel at rId2")
	}
	if !strings.Contains(
		relsXML,
		`Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/chart"`,
	) {
		t.Fatalf("expected chart rel at rId3")
	}

	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")
	if !strings.Contains(slideXML, `c:chart`) || !strings.Contains(slideXML, `r:id="rId3"`) {
		t.Fatalf("expected chart frame to reference rId3")
	}
}

func TestCreateWithSlidesBarChartStyleOptions(t *testing.T) {
	chart := charts.NewBarChart(
		[]string{"Q1", "Q2"},
		[]float64{3, 7},
	).
		WithTitle("Styled Bar").
		WithSeriesName("Revenue").
		WithLegend(true).
		WithAxisTitles("Quarter", "USD").
		WithMajorGridlines(false)

	slides := []pptx.SlideContent{
		pptx.NewSlide("Styled").WithBarChart(chart),
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
		`<c:v>Revenue</c:v>`,
		`<c:legendPos val="r"/>`,
		`<a:t>Quarter</a:t>`,
		`<a:t>USD</a:t>`,
	}
	for _, needle := range checks {
		if !strings.Contains(chartXML, needle) {
			t.Fatalf("expected %q in chart XML", needle)
		}
	}
	if strings.Contains(chartXML, "<c:majorGridlines/>") {
		t.Fatalf("did not expect major gridlines when disabled")
	}
}

func TestCreateWithSlidesLineChartSmoothOption(t *testing.T) {
	chart := charts.NewLineChart(
		[]string{"W1", "W2", "W3"},
		[]float64{1, 2, 3},
	).WithSmooth(true)

	slides := []pptx.SlideContent{
		pptx.NewSlide("Smooth").WithLineChart(chart),
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
	if !strings.Contains(chartXML, `<c:smooth val="1"/>`) {
		t.Fatalf("expected smooth line flag in chart XML")
	}
}

func TestCreateWithSlidesRejectsEmptySeriesName(t *testing.T) {
	chart := charts.NewBarChart(
		[]string{"A"},
		[]float64{1},
	).WithSeriesName("   ")
	slides := []pptx.SlideContent{
		pptx.NewSlide("Broken").WithBarChart(chart),
	}

	_, err := pptx.CreateWithSlides("Demo", slides)
	if err == nil {
		t.Fatalf("expected validation error for empty series name")
	}
}

func TestLineChartSupportsNegativeValues(t *testing.T) {
	chart := charts.NewLineChart(
		[]string{"Jan", "Feb", "Mar"},
		[]float64{10, -5, 15},
	)
	err := chart.Validate(1)
	if err != nil {
		t.Fatalf("expected LineChart to support negative values: %v", err)
	}
}

func TestBarChartRejectsNegativeValues(t *testing.T) {
	chart := charts.NewBarChart(
		[]string{"Q1", "Q2"},
		[]float64{10, -5},
	)
	err := chart.Validate(1)
	if err == nil {
		t.Fatalf("expected BarChart to reject negative values")
	}
}
