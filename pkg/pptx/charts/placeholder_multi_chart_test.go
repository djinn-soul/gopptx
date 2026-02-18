package charts_test

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

func TestPlaceholderMultiChart(t *testing.T) {
	chart1 := charts.NewBarChart([]string{"A", "B"}, []float64{10, 20}).WithTitle("Chart 1")
	chart2 := charts.NewBarChart([]string{"C", "D"}, []float64{30, 40}).WithTitle("Chart 2")

	slides := []pptx.SlideContent{
		pptx.NewSlide("Multi Chart").
			WithPlaceholderChartAs(1, "body", chart1).
			WithPlaceholderChartAs(2, "body", chart2),
	}

	data, err := pptx.CreateWithSlides("Multi Chart Demo", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides failed: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip.NewReader failed: %v", err)
	}

	verifySlideContent(t, zr)
	verifySlideRelationships(t, zr)
	verifyChartFilesExist(t, zr)
}

func verifySlideContent(t *testing.T, zr *zip.Reader) {
	foundSlide := false
	for _, f := range zr.File {
		if f.Name == "ppt/slides/slide1.xml" {
			foundSlide = true
			xml := readZipFile(t, zr, "ppt/slides/slide1.xml")

			expectedStrings := []string{
				`r:id="rId2"`, `r:id="rId3"`,
				"Placeholder Chart 1", "Placeholder Chart 2",
			}
			for _, s := range expectedStrings {
				if !strings.Contains(xml, s) {
					t.Errorf("missing %s in slide1.xml", s)
				}
			}
		}
	}
	if !foundSlide {
		t.Fatal("ppt/slides/slide1.xml not found")
	}
}

func verifySlideRelationships(t *testing.T, zr *zip.Reader) {
	foundRels := false
	for _, f := range zr.File {
		if f.Name == "ppt/slides/_rels/slide1.xml.rels" {
			foundRels = true
			xml := readZipFile(t, zr, "ppt/slides/_rels/slide1.xml.rels")

			if !strings.Contains(xml, `Id="rId2"`) || !strings.Contains(xml, `Target="../charts/chart1.xml"`) {
				t.Errorf("incorrect rel for chart 1: %s", xml)
			}
			if !strings.Contains(xml, `Id="rId3"`) || !strings.Contains(xml, `Target="../charts/chart2.xml"`) {
				t.Errorf("incorrect rel for chart 2: %s", xml)
			}
		}
	}
	if !foundRels {
		t.Fatal("ppt/slides/_rels/slide1.xml.rels not found")
	}
}

func verifyChartFilesExist(t *testing.T, zr *zip.Reader) {
	chart1Found, chart2Found := false, false
	for _, f := range zr.File {
		if f.Name == "ppt/charts/chart1.xml" {
			chart1Found = true
		}
		if f.Name == "ppt/charts/chart2.xml" {
			chart2Found = true
		}
	}
	if !chart1Found || !chart2Found {
		t.Errorf("missing chart files: chart1=%v, chart2=%v", chart1Found, chart2Found)
	}
}
