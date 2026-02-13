package editor

import (
	"archive/zip"
	"bytes"
	"testing"
)

func TestGenerateExcelForChart(t *testing.T) {
	categories := []string{"Cat 1", "Cat 2"}
	values := []float64{10.5, 20.0}

	data, err := generateExcelForChart(categories, values)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify it's a valid zip
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("invalid zip archive: %v", err)
	}

	requiredFiles := []string{
		"[Content_Types].xml",
		"_rels/.rels",
		"xl/workbook.xml",
		"xl/_rels/workbook.xml.rels",
		"xl/styles.xml",
		"xl/worksheets/sheet1.xml",
	}

	for _, req := range requiredFiles {
		found := false
		for _, f := range zr.File {
			if f.Name == req {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing required file in xlsx: %s", req)
		}
	}
}
