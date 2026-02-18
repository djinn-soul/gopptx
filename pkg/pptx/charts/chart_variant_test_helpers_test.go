package charts_test

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func chartXMLForSlide(t *testing.T, slide pptx.SlideContent) string {
	t.Helper()
	data, err := pptx.CreateWithSlides("Demo", []pptx.SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	return readZipFile(t, zr, "ppt/charts/chart1.xml")
}

func readZipFile(t *testing.T, zr *zip.Reader, name string) string {
	t.Helper()
	for _, f := range zr.File {
		if f.Name != name {
			continue
		}
		r, err := f.Open()
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}
		defer func() { _ = r.Close() }()
		buf := new(bytes.Buffer)
		if _, readErr := buf.ReadFrom(r); readErr != nil {
			t.Fatalf("read %s: %v", name, readErr)
		}
		return buf.String()
	}
	t.Fatalf("file %s not found in zip", name)
	return ""
}

func zipHasFile(zr *zip.Reader, name string) bool {
	for _, f := range zr.File {
		if f.Name == name {
			return true
		}
	}
	return false
}

func assertXMLContainsAll(t *testing.T, xml string, checks []string) {
	t.Helper()
	for _, needle := range checks {
		if !strings.Contains(xml, needle) {
			t.Fatalf("expected %q in chart XML", needle)
		}
	}
}
