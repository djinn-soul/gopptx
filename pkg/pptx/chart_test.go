package pptx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestCreateWithSlidesEmbedsBarChart(t *testing.T) {
	chart := NewBarChart(
		[]string{"Q1", "Q2", "Q3"},
		[]float64{10, 22, 15},
	).WithTitle("Quarterly Revenue")

	slides := []SlideContent{
		NewSlide("Chart Slide").WithBarChart(chart),
	}

	data, err := CreateWithSlides("Demo", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")
	checks := []string{
		`name="Bar Chart Plot"`,
		`Quarterly Revenue`,
		`<a:t>Q1</a:t>`,
		`<a:t>22.0</a:t>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in slide XML", needle)
		}
	}
}

func TestCreateWithSlidesRejectsInvalidBarChart(t *testing.T) {
	chart := NewBarChart(
		[]string{"Q1", "Q2"},
		[]float64{4},
	)
	slides := []SlideContent{
		NewSlide("Broken Chart").WithBarChart(chart),
	}

	_, err := CreateWithSlides("Demo", slides)
	if err == nil {
		t.Fatalf("expected chart validation error")
	}
}
