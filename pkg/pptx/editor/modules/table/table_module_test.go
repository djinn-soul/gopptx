package table

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
)

func buildTestSlideAndFrame(t *testing.T, shapeID int) ([]byte, []byte) {
	t.Helper()
	spec := &pptxxml.TableSpec{
		X:          100,
		Y:          200,
		CX:         4000,
		CY:         2000,
		Rows:       [][]string{{"r0c0", "r0c1"}, {"r1c0", "r1c1"}},
		StyledRows: [][]pptxxml.TableCellSpec{{{}, {}}, {{}, {}}},
	}
	frame := []byte(pptxxml.RenderTable(spec, shapeID))
	slide := []byte(`<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">` +
		`<p:cSld><p:spTree>` + string(frame) + `</p:spTree></p:cSld></p:sld>`)
	return slide, frame
}

func TestFrameReadAndDimensions(t *testing.T) {
	slideContent, _ := buildTestSlideAndFrame(t, 7)
	start, end, frame, err := FindTableFrame(slideContent, 7)
	if err != nil {
		t.Fatalf("FindTableFrame failed: %v", err)
	}
	if start < 0 || end <= start || len(frame) == 0 {
		t.Fatalf("invalid frame bounds start=%d end=%d len=%d", start, end, len(frame))
	}
	if _, _, _, err = FindTableFrame(slideContent, 999); err == nil {
		t.Fatal("expected missing shape id error")
	}

	tableXML, err := ExtractTableXML(frame)
	if err != nil {
		t.Fatalf("ExtractTableXML failed: %v", err)
	}
	if !strings.Contains(string(tableXML), "<a:tbl") {
		t.Fatalf("expected <a:tbl> content, got: %s", string(tableXML))
	}

	parsed, err := ParseTable(frame)
	if err != nil {
		t.Fatalf("ParseTable failed: %v", err)
	}
	rows, cols := Dimensions(parsed)
	if rows != 2 || cols != 2 {
		t.Fatalf("unexpected dimensions rows=%d cols=%d", rows, cols)
	}
}

//nolint:gocyclo,cyclop // Table mutation coverage intentionally exercises a broad end-to-end scenario in one test.
func TestMutationsAndInfoProjection(t *testing.T) {
	_, frame := buildTestSlideAndFrame(t, 9)

	info, err := BuildTableInfo(frame)
	if err != nil {
		t.Fatalf("BuildTableInfo failed: %v", err)
	}
	tableMeta := info["table"].(map[string]any)
	if tableMeta["row_count"].(int) != 2 || tableMeta["col_count"].(int) != 2 {
		t.Fatalf("unexpected table counts: %+v", tableMeta)
	}

	updated, err := UpdateTableCellTextInFrame(frame, 1, 1, "changed")
	if err != nil {
		t.Fatalf("UpdateTableCellTextInFrame failed: %v", err)
	}
	info, err = BuildTableInfo(updated)
	if err != nil {
		t.Fatalf("BuildTableInfo(updated) failed: %v", err)
	}
	cells := info["table"].(map[string]any)["cells"].([]map[string]any)
	found := false
	for _, c := range cells {
		if c["row"].(int) == 1 && c["col"].(int) == 1 && c["text"].(string) == "changed" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected updated cell text in BuildTableInfo output")
	}

	updated, err = UpdateTableRowHeightInFrame(updated, 0, 1500)
	if err != nil {
		t.Fatalf("UpdateTableRowHeightInFrame failed: %v", err)
	}
	updated, err = UpdateTableColumnWidthInFrame(updated, 0, 2500)
	if err != nil {
		t.Fatalf("UpdateTableColumnWidthInFrame failed: %v", err)
	}
	info, err = BuildTableInfo(updated)
	if err != nil {
		t.Fatalf("BuildTableInfo(updated row/col) failed: %v", err)
	}
	rowHeights := info["table"].(map[string]any)["row_heights"].([]int64)
	colWidths := info["table"].(map[string]any)["column_widths"].([]int64)
	if rowHeights[0] != 1500 || colWidths[0] != 2500 {
		t.Fatalf("unexpected row/col metrics: rows=%v cols=%v", rowHeights, colWidths)
	}

	updated, err = UpdateTableFlagsInFrame(updated, map[string]any{"first_row": true, "band_col": false})
	if err != nil {
		t.Fatalf("UpdateTableFlagsInFrame failed: %v", err)
	}
	info, err = BuildTableInfo(updated)
	if err != nil {
		t.Fatalf("BuildTableInfo(updated flags) failed: %v", err)
	}
	meta := info["table"].(map[string]any)
	if !meta["first_row"].(bool) || meta["band_col"].(bool) {
		t.Fatalf("unexpected flag projection after update: %+v", meta)
	}

	merged, err := MergeCellsInFrame(updated, 0, 0, 1, 1)
	if err != nil {
		t.Fatalf("MergeCellsInFrame failed: %v", err)
	}
	info, err = BuildTableInfo(merged)
	if err != nil {
		t.Fatalf("BuildTableInfo(merged) failed: %v", err)
	}
	cells = info["table"].(map[string]any)["cells"].([]map[string]any)
	var origin, spanned map[string]any
	for _, c := range cells {
		if c["row"].(int) == 0 && c["col"].(int) == 0 {
			origin = c
		}
		if c["row"].(int) == 1 && c["col"].(int) == 1 {
			spanned = c
		}
	}
	if origin == nil || !origin["is_merge_origin"].(bool) {
		t.Fatalf("expected merge origin at [0,0], got %+v", origin)
	}
	if spanned == nil || !spanned["is_spanned"].(bool) {
		t.Fatalf("expected spanned cell at [1,1], got %+v", spanned)
	}

	split, err := SplitCellInFrame(merged, 0, 0)
	if err != nil {
		t.Fatalf("SplitCellInFrame failed: %v", err)
	}
	info, err = BuildTableInfo(split)
	if err != nil {
		t.Fatalf("BuildTableInfo(split) failed: %v", err)
	}
	cells = info["table"].(map[string]any)["cells"].([]map[string]any)
	for _, c := range cells {
		if c["is_merge_origin"].(bool) || c["is_spanned"].(bool) {
			t.Fatalf("expected no merged cells after split, found: %+v", c)
		}
	}
}

func TestStyleAndAttributeHelpers(t *testing.T) {
	_, frame := buildTestSlideAndFrame(t, 11)
	withStyle, err := SetTableStyleInFrame(frame, "{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}")
	if err != nil {
		t.Fatalf("SetTableStyleInFrame failed: %v", err)
	}
	if !strings.Contains(string(withStyle), "<a:tableStyleId>{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}</a:tableStyleId>") {
		t.Fatalf("expected table style id tag in updated frame: %s", string(withStyle))
	}

	tag := []byte(`<a:tr h="100">`)
	tag = SetOrInsertAttr(tag, "h", "200")
	if !strings.Contains(string(tag), `h="200"`) {
		t.Fatalf("SetOrInsertAttr should replace existing attr: %s", string(tag))
	}
	tag = SetOrInsertAttr(tag, "x", "1")
	if !strings.Contains(string(tag), `x="1"`) {
		t.Fatalf("SetOrInsertAttr should insert missing attr: %s", string(tag))
	}

	tc := []byte(`<a:tc rowSpan="2" gridSpan="3"></a:tc>`)
	tc = SetTcAttr(tc, "rowSpan", "4")
	if !strings.Contains(string(tc), `rowSpan="4"`) {
		t.Fatalf("SetTcAttr did not set rowSpan: %s", string(tc))
	}
	tc = RemoveTcAttr(tc, "gridSpan")
	if strings.Contains(string(tc), "gridSpan=") {
		t.Fatalf("RemoveTcAttr did not remove gridSpan: %s", string(tc))
	}

	if !TruthyAttr("1") || !TruthyAttr("true") || TruthyAttr("0") {
		t.Fatal("TruthyAttr returned unexpected values")
	}
}

func TestMutationErrorPaths(t *testing.T) {
	_, frame := buildTestSlideAndFrame(t, 12)
	if _, err := UpdateTableRowHeightInFrame(frame, 10, 100); err == nil {
		t.Fatal("expected out-of-range row error")
	}
	if _, err := UpdateTableRowHeightInFrame(frame, 0, 0); err == nil {
		t.Fatal("expected invalid row height error")
	}
	if _, err := UpdateTableColumnWidthInFrame(frame, 10, 100); err == nil {
		t.Fatal("expected out-of-range column error")
	}
	if _, err := UpdateTableColumnWidthInFrame(frame, 0, 0); err == nil {
		t.Fatal("expected invalid column width error")
	}
	if _, err := MergeCellsInFrame(frame, 1, 1, 0, 0); err == nil {
		t.Fatal("expected invalid merge range ordering error")
	}
	if _, err := SplitCellInFrame(frame, 4, 4); err == nil {
		t.Fatal("expected out-of-range split error")
	}
}
