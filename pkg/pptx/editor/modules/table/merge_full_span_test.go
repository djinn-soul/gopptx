package table

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
)

func buildFrameWithGrid(t *testing.T, shapeID int, rows [][]string) []byte {
	t.Helper()
	styled := make([][]pptxxml.TableCellSpec, len(rows))
	for i, row := range rows {
		styled[i] = make([]pptxxml.TableCellSpec, len(row))
	}
	spec := &pptxxml.TableSpec{
		X:          100,
		Y:          200,
		CX:         4000,
		CY:         2000,
		Rows:       rows,
		StyledRows: styled,
	}
	return []byte(pptxxml.RenderTable(spec, shapeID))
}

// Issue #636: PowerPoint drops the redundant rows when a merge spans every column,
// so a vMerge-only result renders as unmerged.
func TestMergeFullWidthRemovesMergedRows(t *testing.T) {
	frame := buildFrameWithGrid(t, 20, [][]string{{"a"}, {"b"}, {"c"}, {"d"}})

	before, err := ParseTable(frame)
	if err != nil {
		t.Fatalf("ParseTable failed: %v", err)
	}
	wantHeight := before.Rows[1].Height + before.Rows[2].Height

	updated, err := MergeCellsInFrame(frame, 1, 0, 2, 0)
	if err != nil {
		t.Fatalf("MergeCellsInFrame failed: %v", err)
	}

	parsed, err := ParseTable(updated)
	if err != nil {
		t.Fatalf("ParseTable(updated) failed: %v", err)
	}
	rows, cols := Dimensions(parsed)
	if rows != 3 || cols != 1 {
		t.Fatalf("expected 3x1 table after full-width merge, got %dx%d", rows, cols)
	}
	if strings.Contains(string(updated), `vMerge=`) {
		t.Errorf("expected no vMerge placeholders left, got: %s", string(updated))
	}
	if parsed.Rows[1].Cells[0].RowSpan != 0 {
		t.Errorf("expected rowSpan cleared on origin cell, got %d", parsed.Rows[1].Cells[0].RowSpan)
	}
	if wantHeight > 0 && parsed.Rows[1].Height != wantHeight {
		t.Errorf("expected merged row height %d, got %d", wantHeight, parsed.Rows[1].Height)
	}
}

// Merging the whole table collapses to one row; the columns stay folded under the
// origin cell's gridSpan, which PowerPoint renders correctly.
func TestMergeWholeTableCollapsesToSingleRow(t *testing.T) {
	frame := buildFrameWithGrid(t, 22, [][]string{{"a", "b"}, {"c", "d"}})

	updated, err := MergeCellsInFrame(frame, 0, 0, 1, 1)
	if err != nil {
		t.Fatalf("MergeCellsInFrame failed: %v", err)
	}
	parsed, err := ParseTable(updated)
	if err != nil {
		t.Fatalf("ParseTable(updated) failed: %v", err)
	}
	rows, cols := Dimensions(parsed)
	if rows != 1 || cols != 2 {
		t.Fatalf("expected 1x2 table after whole-table merge, got %dx%d", rows, cols)
	}
	if parsed.Rows[0].Cells[0].GridSpan != 2 {
		t.Errorf("expected gridSpan=2 retained on origin cell, got %d", parsed.Rows[0].Cells[0].GridSpan)
	}
}

// A vertical merge that leaves a column outside the range renders correctly in
// PowerPoint, so it keeps its vMerge placeholders.
func TestPartialWidthMergeKeepsSpanAttributes(t *testing.T) {
	frame := buildFrameWithGrid(t, 23, [][]string{{"a", "b", "c"}, {"d", "e", "f"}})

	updated, err := MergeCellsInFrame(frame, 0, 0, 1, 0)
	if err != nil {
		t.Fatalf("MergeCellsInFrame failed: %v", err)
	}
	parsed, err := ParseTable(updated)
	if err != nil {
		t.Fatalf("ParseTable(updated) failed: %v", err)
	}
	rows, cols := Dimensions(parsed)
	if rows != 2 || cols != 3 {
		t.Fatalf("expected 2x3 table unchanged, got %dx%d", rows, cols)
	}
	if parsed.Rows[0].Cells[0].RowSpan != 2 {
		t.Errorf("expected rowSpan=2 on origin cell, got %d", parsed.Rows[0].Cells[0].RowSpan)
	}
	if !strings.Contains(string(updated), `vMerge="1"`) {
		t.Errorf("expected vMerge placeholder retained for partial-width merge")
	}
}

// A horizontal merge across every column of a single-row table renders correctly in
// PowerPoint, so no columns are removed.
func TestSingleRowFullWidthMergeKeepsColumns(t *testing.T) {
	frame := buildFrameWithGrid(t, 21, [][]string{{"a", "b", "c"}})

	updated, err := MergeCellsInFrame(frame, 0, 0, 0, 2)
	if err != nil {
		t.Fatalf("MergeCellsInFrame failed: %v", err)
	}
	parsed, err := ParseTable(updated)
	if err != nil {
		t.Fatalf("ParseTable(updated) failed: %v", err)
	}
	rows, cols := Dimensions(parsed)
	if rows != 1 || cols != 3 {
		t.Fatalf("expected 1x3 table unchanged, got %dx%d", rows, cols)
	}
	if parsed.Rows[0].Cells[0].GridSpan != 3 {
		t.Errorf("expected gridSpan=3 on origin cell, got %d", parsed.Rows[0].Cells[0].GridSpan)
	}
}

func TestBuildTableInfoReportsStyleID(t *testing.T) {
	_, frame := buildTestSlideAndFrame(t, 24)

	info, err := BuildTableInfo(frame)
	if err != nil {
		t.Fatalf("BuildTableInfo failed: %v", err)
	}
	tbl, _ := info["table"].(map[string]any)
	if _, ok := tbl["style_id"]; !ok {
		t.Fatal("expected style_id key in table info")
	}

	const guid = "{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}"
	styled, err := SetTableStyleInFrame(frame, guid)
	if err != nil {
		t.Fatalf("SetTableStyleInFrame failed: %v", err)
	}
	info2, err := BuildTableInfo(styled)
	if err != nil {
		t.Fatalf("BuildTableInfo(styled) failed: %v", err)
	}
	tbl2, _ := info2["table"].(map[string]any)
	if got := tbl2["style_id"]; got != guid {
		t.Errorf("expected style_id %q, got %v", guid, got)
	}
}
