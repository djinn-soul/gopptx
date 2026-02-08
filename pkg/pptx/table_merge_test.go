package pptx

import (
	"archive/zip"
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestCreateWithSlidesEmbedsTableMergeAttrs(t *testing.T) {
	table := NewTable([]int64{2000000, 2000000, 2000000}).
		AddStyledRow([]TableCell{
			NewTableCell("Merged").WithRowSpan(2).WithColSpan(2),
			NewTableCell(""),
			NewTableCell("R1C3"),
		}).
		AddStyledRow([]TableCell{
			NewTableCell(""),
			NewTableCell(""),
			NewTableCell("R2C3"),
		}).
		AddStyledRow([]TableCell{
			NewTableCell("R3C1"),
			NewTableCell("R3C2"),
			NewTableCell("R3C3"),
		})

	data, err := CreateWithSlides("Demo", []SlideContent{NewSlide("Merged Table").WithTable(table)})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")
	checks := []string{
		`<a:tc rowSpan="2" gridSpan="2">`,
		`hMerge="1"`,
		`vMerge="1"`,
		`<a:t>Merged</a:t>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in merged table XML", needle)
		}
	}
}

func TestCreateWithSlidesRejectsTableMergeOverlap(t *testing.T) {
	table := NewTable([]int64{2000000, 2000000, 2000000}).
		AddStyledRow([]TableCell{
			NewTableCell("Merged").WithColSpan(2),
			NewTableCell("invalid overlap"),
			NewTableCell("R1C3"),
		})

	_, err := CreateWithSlides("Demo", []SlideContent{NewSlide("Broken Merge").WithTable(table)})
	if err == nil {
		t.Fatalf("expected merge overlap validation error")
	}
	if !strings.Contains(err.Error(), "covered cells must be empty placeholders") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateWithSlidesRejectsTableMergeOutOfBounds(t *testing.T) {
	table := NewTable([]int64{2000000, 2000000}).
		AddStyledRow([]TableCell{
			NewTableCell("too-wide").WithColSpan(3),
			NewTableCell(""),
		})

	_, err := CreateWithSlides("Demo", []SlideContent{NewSlide("Broken Merge").WithTable(table)})
	if err == nil {
		t.Fatalf("expected merge bounds validation error")
	}
	if !strings.Contains(err.Error(), "merged span (1x3) exceeds table bounds") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTableMergeParityFixtureAgainstPptRs(t *testing.T) {
	fixture := rootTestdataPath("ppt_rs", "merged_cells.pptx")
	if _, err := os.Stat(fixture); err != nil {
		t.Skipf("missing ppt-rs merged-cells fixture: %v", err)
	}

	reference := fixtureSlideXML(t, "merged_cells.pptx", "ppt/slides/slide1.xml")
	ours := generatedSlideXML(t,
		NewSlide("Merged Cells").WithTable(
			NewTable([]int64{2000000, 2000000, 2000000}).
				AddStyledRow([]TableCell{
					NewTableCell("Merged Header").WithRowSpan(2).WithColSpan(2),
					NewTableCell(""),
					NewTableCell("C1"),
				}).
				AddStyledRow([]TableCell{
					NewTableCell(""),
					NewTableCell(""),
					NewTableCell("C2"),
				}).
				AddStyledRow([]TableCell{
					NewTableCell("R3C1"),
					NewTableCell("R3C2"),
					NewTableCell("R3C3"),
				}),
		),
	)

	tokens := []string{
		`rowSpan="2"`,
		`gridSpan="2"`,
		`hMerge="1"`,
		`vMerge="1"`,
	}
	assertContainsTokens(t, "ppt-rs merged-cells fixture", reference, tokens)
	assertContainsTokens(t, "gopptx merged-cells parity deck", ours, tokens)
}
