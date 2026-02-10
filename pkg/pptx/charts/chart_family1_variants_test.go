package charts_test

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

func TestCreateWithSlidesEmbedsBarHorizontalChart(t *testing.T) {
	chart := charts.NewBarHorizontalChart([]string{"A", "B"}, []float64{1, 2}).WithTitle("H")
	xml := chartXMLForSlide(t, pptx.NewSlide("S").WithBarHorizontalChart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:barChart>`,
		`<c:barDir val="bar"/>`,
		`<c:grouping val="clustered"/>`,
		`<a:t>H</a:t>`,
	})
}

func TestCreateWithSlidesRejectsInvalidBarHorizontalChart(t *testing.T) {
	chart := charts.NewBarHorizontalChart([]string{"A"}, []float64{1}).WithValueFormat(" ")
	_, err := pptx.CreateWithSlides("Demo", []pptx.SlideContent{pptx.NewSlide("S").WithBarHorizontalChart(chart)})
	if err == nil {
		t.Fatalf("expected bar horizontal validation error")
	}
}

func TestCreateWithSlidesEmbedsBarStackedChart(t *testing.T) {
	chart := charts.NewBarStackedChart([]string{"A", "B"}, []float64{1, 2})
	xml := chartXMLForSlide(t, pptx.NewSlide("S").WithBarStackedChart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:barDir val="col"/>`,
		`<c:grouping val="stacked"/>`,
	})
}

func TestCreateWithSlidesRejectsInvalidBarStackedChart(t *testing.T) {
	chart := charts.NewBarStackedChart([]string{"A"}, []float64{1}).WithBarColor("GGGGGG")
	_, err := pptx.CreateWithSlides("Demo", []pptx.SlideContent{pptx.NewSlide("S").WithBarStackedChart(chart)})
	if err == nil {
		t.Fatalf("expected bar stacked validation error")
	}
}

func TestCreateWithSlidesEmbedsBarStacked100Chart(t *testing.T) {
	chart := charts.NewBarStacked100Chart([]string{"A", "B"}, []float64{1, 2})
	xml := chartXMLForSlide(t, pptx.NewSlide("S").WithBarStacked100Chart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:barDir val="col"/>`,
		`<c:grouping val="percentStacked"/>`,
	})
}

func TestCreateWithSlidesRejectsInvalidBarStacked100Chart(t *testing.T) {
	chart := charts.NewBarStacked100Chart([]string{"A", "B"}, []float64{1, 2}).WithValueRange(10, 1)
	_, err := pptx.CreateWithSlides("Demo", []pptx.SlideContent{pptx.NewSlide("S").WithBarStacked100Chart(chart)})
	if err == nil {
		t.Fatalf("expected bar stacked 100 validation error")
	}
}

func TestCreateWithSlidesEmbedsLineMarkersChart(t *testing.T) {
	chart := charts.NewLineMarkersChart([]string{"A", "B"}, []float64{1, 2})
	xml := chartXMLForSlide(t, pptx.NewSlide("S").WithLineMarkersChart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:lineChart>`,
		`<c:grouping val="standard"/>`,
		`<c:marker><c:symbol val="circle"/></c:marker>`,
	})
}

func TestCreateWithSlidesRejectsInvalidLineMarkersChart(t *testing.T) {
	chart := charts.NewLineMarkersChart([]string{"A"}, []float64{1}).WithLegendPosition("x")
	_, err := pptx.CreateWithSlides("Demo", []pptx.SlideContent{pptx.NewSlide("S").WithLineMarkersChart(chart)})
	if err == nil {
		t.Fatalf("expected line markers validation error")
	}
}

func TestCreateWithSlidesEmbedsLineStackedChart(t *testing.T) {
	chart := charts.NewLineStackedChart([]string{"A", "B"}, []float64{1, 2})
	xml := chartXMLForSlide(t, pptx.NewSlide("S").WithLineStackedChart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:lineChart>`,
		`<c:grouping val="stacked"/>`,
	})
}

func TestCreateWithSlidesRejectsInvalidLineStackedChart(t *testing.T) {
	chart := charts.NewLineStackedChart([]string{"A"}, []float64{1}).WithSeriesName(" ")
	_, err := pptx.CreateWithSlides("Demo", []pptx.SlideContent{pptx.NewSlide("S").WithLineStackedChart(chart)})
	if err == nil {
		t.Fatalf("expected line stacked validation error")
	}
}

func TestCreateWithSlidesEmbedsAreaStackedChart(t *testing.T) {
	chart := charts.NewAreaStackedChart([]string{"A", "B"}, []float64{1, 2})
	xml := chartXMLForSlide(t, pptx.NewSlide("S").WithAreaStackedChart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:areaChart>`,
		`<c:grouping val="stacked"/>`,
	})
}

func TestCreateWithSlidesRejectsInvalidAreaStackedChart(t *testing.T) {
	chart := charts.NewAreaStackedChart([]string{"A"}, []float64{1}).WithAreaColor("XYZ123")
	_, err := pptx.CreateWithSlides("Demo", []pptx.SlideContent{pptx.NewSlide("S").WithAreaStackedChart(chart)})
	if err == nil {
		t.Fatalf("expected area stacked validation error")
	}
}

func TestCreateWithSlidesEmbedsAreaStacked100Chart(t *testing.T) {
	chart := charts.NewAreaStacked100Chart([]string{"A", "B"}, []float64{1, 2})
	xml := chartXMLForSlide(t, pptx.NewSlide("S").WithAreaStacked100Chart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:areaChart>`,
		`<c:grouping val="percentStacked"/>`,
	})
}

func TestCreateWithSlidesRejectsInvalidAreaStacked100Chart(t *testing.T) {
	chart := charts.NewAreaStacked100Chart([]string{"A"}, []float64{1}).WithValueFormat(" ")
	_, err := pptx.CreateWithSlides("Demo", []pptx.SlideContent{pptx.NewSlide("S").WithAreaStacked100Chart(chart)})
	if err == nil {
		t.Fatalf("expected area stacked 100 validation error")
	}
}
