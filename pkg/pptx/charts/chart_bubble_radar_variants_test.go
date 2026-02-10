package charts_test

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

func TestCreateWithSlidesEmbedsBubbleChart(t *testing.T) {
	chart := charts.NewBubbleChart(
		[]float64{1, 2},
		[]float64{3, 4},
		[]float64{10, 20},
	).WithBubbleScale(120)
	xml := chartXMLForSlide(t, pptx.NewSlide("S").WithBubbleChart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:bubbleChart>`,
		`<c:bubbleScale val="120"/>`,
		`<c:bubbleSize><c:numLit>`,
	})
}

func TestCreateWithSlidesRejectsInvalidBubbleChart(t *testing.T) {
	chart := charts.NewBubbleChart(
		[]float64{1, 2},
		[]float64{3},
		[]float64{10, 20},
	)
	_, err := pptx.CreateWithSlides("Demo", []pptx.SlideContent{pptx.NewSlide("S").WithBubbleChart(chart)})
	if err == nil {
		t.Fatalf("expected bubble validation error")
	}
}

func TestCreateWithSlidesEmbedsRadarChart(t *testing.T) {
	chart := charts.NewRadarChart([]string{"A", "B"}, []float64{1, 2})
	xml := chartXMLForSlide(t, pptx.NewSlide("S").WithRadarChart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:radarChart>`,
		`<c:radarStyle val="marker"/>`,
	})
}

func TestCreateWithSlidesRejectsInvalidRadarChart(t *testing.T) {
	chart := charts.NewRadarChart([]string{"A"}, []float64{1})
	chart.RadarStyle = "weird"
	_, err := pptx.CreateWithSlides("Demo", []pptx.SlideContent{pptx.NewSlide("S").WithRadarChart(chart)})
	if err == nil {
		t.Fatalf("expected radar validation error")
	}
}

func TestCreateWithSlidesEmbedsRadarFilledChart(t *testing.T) {
	chart := charts.NewRadarFilledChart([]string{"A", "B"}, []float64{1, 2})
	xml := chartXMLForSlide(t, pptx.NewSlide("S").WithRadarFilledChart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:radarChart>`,
		`<c:radarStyle val="filled"/>`,
	})
}

func TestCreateWithSlidesRejectsInvalidRadarFilledChart(t *testing.T) {
	chart := charts.NewRadarFilledChart([]string{"A"}, []float64{1}).WithLegendPosition("middle")
	_, err := pptx.CreateWithSlides("Demo", []pptx.SlideContent{pptx.NewSlide("S").WithRadarFilledChart(chart)})
	if err == nil {
		t.Fatalf("expected radar filled validation error")
	}
}
