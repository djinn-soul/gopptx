package pptx

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

func createPresentationWithBarChart(t *testing.T) string {
	t.Helper()
	return createPresentationWithChart(
		t,
		charts.NewBarChart([]string{"A", "B"}, []float64{1, 2}).WithTitle("Fixture Chart"),
	)
}

func createPresentationWithChart(t *testing.T, chartDef charts.ChartDefinition) string {
	t.Helper()

	tmpDir := t.TempDir()
	pptxPath := filepath.Join(tmpDir, "chart-api-test.pptx")

	data, err := Create("Chart API Test", 1)
	if err != nil {
		t.Fatalf("create base presentation: %v", err)
	}
	if err := os.WriteFile(pptxPath, data, 0o600); err != nil {
		t.Fatalf("write base presentation: %v", err)
	}

	ed, err := OpenPresentationEditor(pptxPath)
	if err != nil {
		t.Fatalf("open editor for fixture: %v", err)
	}
	if err := ed.AddChart(0, chartDef); err != nil {
		_ = ed.Close()
		t.Fatalf("add chart fixture: %v", err)
	}
	if err := ed.Save(pptxPath); err != nil {
		_ = ed.Close()
		t.Fatalf("save fixture chart: %v", err)
	}
	if err := ed.Close(); err != nil {
		t.Fatalf("close fixture editor: %v", err)
	}
	return pptxPath
}

const chartPartPath = "ppt/charts/chart1.xml"

func readZipEntry(t *testing.T, zipPath string) string {
	t.Helper()

	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatalf("open zip reader: %v", err)
	}
	defer zr.Close()

	for _, f := range zr.File {
		if f.Name != chartPartPath {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("open zip part %s: %v", chartPartPath, err)
		}
		defer rc.Close()
		data, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("read zip part %s: %v", chartPartPath, err)
		}
		return string(data)
	}

	t.Fatalf("zip part %s not found", chartPartPath)
	return ""
}

func intPtr(v int) *int {
	return &v
}
