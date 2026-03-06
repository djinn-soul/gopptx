package charts_test

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

func TestCreateWithSlidesBarDataLabelFormatting(t *testing.T) {
	chart := charts.NewBarChart([]string{"A", "B"}, []float64{10, 20}).
		WithDataLabels(true).
		WithDataLabelPosition(charts.DataLabelPositionOutsideEnd).
		WithDataLabelContent(true, false, false, false)

	xml := chartXMLForSlide(t, pptx.NewSlide("S").WithBarChart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:dLblPos val="outEnd"/>`,
		`<c:showVal val="1"/>`,
	})
	if strings.Contains(xml, `<c:showPercent val="1"/>`) {
		t.Fatalf("expected bar labels to omit showPercent when disabled")
	}
	if strings.Contains(xml, `<c:showCatName val="1"/>`) {
		t.Fatalf("expected bar labels to omit showCatName when disabled")
	}
}

func TestCreateWithSlidesPieDataLabelDefaults(t *testing.T) {
	chart := charts.NewPieChart([]string{"A", "B"}, []float64{25, 75}).WithDataLabels(true)
	xml := chartXMLForSlide(t, pptx.NewSlide("S").WithPieChart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:showCatName val="1"/>`,
		`<c:showPercent val="1"/>`,
	})
	if strings.Contains(xml, `<c:showVal val="1"/>`) {
		t.Fatalf("expected pie default labels to omit showVal")
	}
}

func TestCreateWithSlidesPieDataLabelCustomFormatting(t *testing.T) {
	chart := charts.NewPieChart([]string{"A", "B"}, []float64{25, 75}).
		WithDataLabels(true).
		WithDataLabelPosition(charts.DataLabelPositionCenter).
		WithDataLabelContent(true, false, false, false)

	xml := chartXMLForSlide(t, pptx.NewSlide("S").WithPieChart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:dLblPos val="ctr"/>`,
		`<c:showVal val="1"/>`,
	})
	if strings.Contains(xml, `<c:showPercent val="1"/>`) {
		t.Fatalf("expected pie labels to omit showPercent when disabled")
	}
	if strings.Contains(xml, `<c:showCatName val="1"/>`) {
		t.Fatalf("expected pie labels to omit showCatName when disabled")
	}
}

func TestCreateWithSlidesRejectsInvalidDataLabelPosition(t *testing.T) {
	chart := charts.NewBarChart([]string{"A", "B"}, []float64{1, 2}).
		WithDataLabels(true).
		WithDataLabelPosition("diagonal")

	_, err := pptx.CreateWithSlides("Demo", []pptx.SlideContent{pptx.NewSlide("S").WithBarChart(chart)})
	if err == nil {
		t.Fatalf("expected data-label position validation error")
	}
	if !strings.Contains(err.Error(), "data-label position") {
		t.Fatalf("expected data-label position error, got: %v", err)
	}
}

func TestCreateWithSlidesRadarDataLabelFormatting(t *testing.T) {
	chart := charts.NewRadarChart([]string{"A", "B"}, []float64{2, 3}).
		WithDataLabels(true).
		WithDataLabelPosition(charts.DataLabelPositionBestFit).
		WithDataLabelContent(false, true, false, false)

	xml := chartXMLForSlide(t, pptx.NewSlide("S").WithRadarChart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:dLblPos val="bestFit"/>`,
		`<c:showCatName val="1"/>`,
	})
	if strings.Contains(xml, `<c:showVal val="1"/>`) {
		t.Fatalf("expected radar labels to omit showVal when disabled")
	}
}

func TestCreateWithSlidesComboDataLabelFormatting(t *testing.T) {
	chart := charts.NewComboChart(
		[]string{"Q1", "Q2"},
		[]charts.Series{{Name: "Bars", Values: []float64{10, 20}}},
		[]charts.Series{{Name: "Line", Values: []float64{15, 25}}},
	).WithDataLabels(true).
		WithDataLabelPosition(charts.DataLabelPositionInsideEnd).
		WithDataLabelContent(true, false, false, false)

	xml := chartXMLForSlide(t, pptx.NewSlide("S").WithComboChart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:dLblPos val="inEnd"/>`,
		`<c:showVal val="1"/>`,
	})
}
