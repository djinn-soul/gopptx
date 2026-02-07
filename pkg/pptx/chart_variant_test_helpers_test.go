package pptx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func chartXMLForSlide(t *testing.T, slide SlideContent) string {
	t.Helper()
	data, err := CreateWithSlides("Demo", []SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	return readZipFile(t, zr, "ppt/charts/chart1.xml")
}

func assertXMLContainsAll(t *testing.T, xml string, checks []string) {
	t.Helper()
	for _, needle := range checks {
		if !strings.Contains(xml, needle) {
			t.Fatalf("expected %q in chart XML", needle)
		}
	}
}
