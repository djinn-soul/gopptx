package pptx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestCreateWithSlidesEmbedsStockHLCChart(t *testing.T) {
	chart := NewStockHLCChart(
		[]string{"D1", "D2"},
		[]float64{10, 11},
		[]float64{7, 8},
		[]float64{9, 10},
	)
	xml := chartXMLForSlide(t, NewSlide("S").WithStockHLCChart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:stockChart>`,
		`<c:v>High</c:v>`,
		`<c:v>Low</c:v>`,
		`<c:v>Close</c:v>`,
	})
}

func TestCreateWithSlidesRejectsInvalidStockHLCChart(t *testing.T) {
	chart := NewStockHLCChart(
		[]string{"D1", "D2"},
		[]float64{10, 11},
		[]float64{7},
		[]float64{9, 10},
	)
	_, err := CreateWithSlides("Demo", []SlideContent{NewSlide("S").WithStockHLCChart(chart)})
	if err == nil {
		t.Fatalf("expected stock HLC validation error")
	}
}

func TestCreateWithSlidesEmbedsStockOHLCChart(t *testing.T) {
	chart := NewStockOHLCChart(
		[]string{"D1", "D2"},
		[]float64{8, 9},
		[]float64{10, 11},
		[]float64{7, 8},
		[]float64{9, 10},
	)
	xml := chartXMLForSlide(t, NewSlide("S").WithStockOHLCChart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:stockChart>`,
		`<c:v>Open</c:v>`,
		`<c:v>High</c:v>`,
		`<c:v>Low</c:v>`,
		`<c:v>Close</c:v>`,
	})
}

func TestCreateWithSlidesRejectsInvalidStockOHLCChart(t *testing.T) {
	chart := NewStockOHLCChart(
		[]string{"D1", "D2"},
		[]float64{8},
		[]float64{10, 11},
		[]float64{7, 8},
		[]float64{9, 10},
	)
	_, err := CreateWithSlides("Demo", []SlideContent{NewSlide("S").WithStockOHLCChart(chart)})
	if err == nil {
		t.Fatalf("expected stock OHLC validation error")
	}
}

func TestCreateWithSlidesEmbedsComboChart(t *testing.T) {
	chart := NewComboChart(
		[]string{"Q1", "Q2"},
		[]Series{{Name: "Bars", Values: []float64{1, 2}}},
		[]Series{{Name: "Line", Values: []float64{2, 3}}},
	)
	xml := chartXMLForSlide(t, NewSlide("S").WithComboChart(chart))
	assertXMLContainsAll(t, xml, []string{
		`<c:barChart>`,
		`<c:lineChart>`,
		`<c:v>Bars</c:v>`,
		`<c:v>Line</c:v>`,
	})
}

func TestCreateWithSlidesRejectsInvalidComboChart(t *testing.T) {
	chart := NewComboChart(
		[]string{"Q1", "Q2"},
		[]Series{{Name: "Bars", Values: []float64{1, 2}}},
		nil,
	)
	_, err := CreateWithSlides("Demo", []SlideContent{NewSlide("S").WithComboChart(chart)})
	if err == nil {
		t.Fatalf("expected combo validation error")
	}
}

func TestCreateWithSlidesEmbedsComboChartRelationshipsAndContentType(t *testing.T) {
	chart := NewComboChart(
		[]string{"Q1", "Q2"},
		[]Series{{Name: "Bars", Values: []float64{1, 2}}},
		[]Series{{Name: "Line", Values: []float64{2, 3}}},
	)
	data, err := CreateWithSlides("Demo", []SlideContent{NewSlide("S").WithComboChart(chart)})
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

	relsXML := readZipFile(t, zr, "ppt/slides/_rels/slide1.xml.rels")
	if !strings.Contains(relsXML, `relationships/chart"`) {
		t.Fatalf("expected chart relationship type in slide rels")
	}
	if !strings.Contains(relsXML, `Target="../charts/chart1.xml"`) {
		t.Fatalf("expected chart target in slide rels")
	}

	contentTypesXML := readZipFile(t, zr, "[Content_Types].xml")
	if !strings.Contains(contentTypesXML, `PartName="/ppt/charts/chart1.xml"`) {
		t.Fatalf("expected chart part content type override")
	}
	if !strings.Contains(contentTypesXML, `application/vnd.openxmlformats-officedocument.drawingml.chart+xml`) {
		t.Fatalf("expected chart content type")
	}
}
