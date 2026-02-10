package charts_test

import (
	"archive/zip"
	"bytes"
	"io"
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
			WithPlaceholderChart(1, chart1).
			WithPlaceholderChart(2, chart2),
	}

	data, err := pptx.CreateWithSlides("Multi Chart Demo", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides failed: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip.NewReader failed: %v", err)
	}

	// Verify slide XML contains both chart references
	foundSlide := false
	for _, f := range zr.File {
		if f.Name == "ppt/slides/slide1.xml" {
			foundSlide = true
			rc, err := f.Open()
			if err != nil {
				t.Fatalf("failed to open slide1.xml: %v", err)
			}
			content, err := io.ReadAll(rc)
			if err != nil {
				t.Fatalf("failed to read slide1.xml: %v", err)
			}
			if err := rc.Close(); err != nil {
				t.Errorf("failed to close rc: %v", err)
			}
			xml := string(content)

			if !strings.Contains(xml, `r:id="rId2"`) {
				t.Errorf("missing rId2 for chart 1 in slide1.xml")
			}
			if !strings.Contains(xml, `r:id="rId3"`) {
				t.Errorf("missing rId3 for chart 2 in slide1.xml")
			}
			if !strings.Contains(xml, "Placeholder Chart 1") {
				t.Errorf("missing name for placeholder chart 1")
			}
			if !strings.Contains(xml, "Placeholder Chart 2") {
				t.Errorf("missing name for placeholder chart 2")
			}
		}
	}
	if !foundSlide {
		t.Fatal("ppt/slides/slide1.xml not found")
	}

	// Verify relationships
	foundRels := false
	for _, f := range zr.File {
		if f.Name == "ppt/slides/_rels/slide1.xml.rels" {
			foundRels = true
			rc, err := f.Open()
			if err != nil {
				t.Fatalf("failed to open slide1.xml.rels: %v", err)
			}
			content, err := io.ReadAll(rc)
			if err != nil {
				t.Fatalf("failed to read slide1.xml.rels: %v", err)
			}
			if err := rc.Close(); err != nil {
				t.Errorf("failed to close rc: %v", err)
			}
			xml := string(content)

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

	// Verify both chart files exist
	chart1Found := false
	chart2Found := false
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
