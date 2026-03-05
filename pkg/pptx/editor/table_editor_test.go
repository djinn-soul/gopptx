package editor

import (
	"regexp"
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestMergeTableCellsValidatesBoundsAndOrder(t *testing.T) {
	e := newTableEditorFixture()
	shapeID, err := e.AddTable(0, 3, 3, 0, 0, 0, 0)
	if err != nil {
		t.Fatalf("AddTable failed: %v", err)
	}

	if err := e.MergeTableCells(0, shapeID, 2, 2, 1, 1); err == nil {
		t.Fatalf("expected error for unordered merge coordinates")
	}
	if err := e.MergeTableCells(0, shapeID, -1, 0, 0, 0); err == nil {
		t.Fatalf("expected error for negative merge coordinates")
	}
	if err := e.MergeTableCells(0, shapeID, 0, 0, 3, 1); err == nil {
		t.Fatalf("expected error for out-of-bounds merge row")
	}
}

func TestSplitTableCellOnlyAffectsTargetedMerge(t *testing.T) {
	e := newTableEditorFixture()
	shapeID, err := e.AddTable(0, 4, 4, 0, 0, 0, 0)
	if err != nil {
		t.Fatalf("AddTable failed: %v", err)
	}
	if err := e.MergeTableCells(0, shapeID, 0, 0, 1, 1); err != nil {
		t.Fatalf("first merge failed: %v", err)
	}
	if err := e.MergeTableCells(0, shapeID, 2, 2, 3, 3); err != nil {
		t.Fatalf("second merge failed: %v", err)
	}
	if err := e.SplitTableCell(0, shapeID, 0, 0); err != nil {
		t.Fatalf("split failed: %v", err)
	}

	tbl, err := e.GetTable(0, shapeID)
	if err != nil {
		t.Fatalf("GetTable failed: %v", err)
	}
	meta := tbl["table"].(map[string]any)
	cell00 := findTableCell(meta["cells"].([]map[string]any), 0, 0)
	if cell00 == nil || cell00["is_merge_origin"].(bool) {
		t.Fatalf("expected [0,0] to no longer be merge origin")
	}
	cell11 := findTableCell(meta["cells"].([]map[string]any), 1, 1)
	if cell11 == nil || cell11["is_spanned"].(bool) {
		t.Fatalf("expected [1,1] to no longer be spanned")
	}
	cell22 := findTableCell(meta["cells"].([]map[string]any), 2, 2)
	if cell22 == nil || !cell22["is_merge_origin"].(bool) {
		t.Fatalf("expected [2,2] merge origin to remain")
	}
	cell33 := findTableCell(meta["cells"].([]map[string]any), 3, 3)
	if cell33 == nil || !cell33["is_spanned"].(bool) {
		t.Fatalf("expected [3,3] spanned cell to remain")
	}
}

func TestSplitTableCellSupportsLargeSpans(t *testing.T) {
	e := newTableEditorFixture()
	shapeID, err := e.AddTable(0, 1, 6, 0, 0, 0, 0)
	if err != nil {
		t.Fatalf("AddTable failed: %v", err)
	}
	if err := e.MergeTableCells(0, shapeID, 0, 0, 0, 5); err != nil {
		t.Fatalf("merge failed: %v", err)
	}
	if err := e.SplitTableCell(0, shapeID, 0, 0); err != nil {
		t.Fatalf("split failed: %v", err)
	}

	tbl, err := e.GetTable(0, shapeID)
	if err != nil {
		t.Fatalf("GetTable failed: %v", err)
	}
	meta := tbl["table"].(map[string]any)
	cells := meta["cells"].([]map[string]any)
	for col := range 6 {
		c := findTableCell(cells, 0, col)
		if c == nil {
			t.Fatalf("missing cell [0,%d]", col)
		}
		if c["row_span"].(int) != 1 || c["col_span"].(int) != 1 {
			t.Fatalf("expected [0,%d] span 1x1, got %v x %v", col, c["row_span"], c["col_span"])
		}
		if c["is_spanned"].(bool) || c["is_merge_origin"].(bool) {
			t.Fatalf("expected [0,%d] unmerged after split", col)
		}
	}
}

func TestSetTableStyleSupportsSelfClosingTblPr(t *testing.T) {
	e := newTableEditorFixture()
	shapeID, err := e.AddTable(0, 2, 2, 0, 0, 0, 0)
	if err != nil {
		t.Fatalf("AddTable failed: %v", err)
	}

	partPath := e.slides[0].Part
	slideXML, ok := e.parts.Get(partPath)
	if !ok {
		t.Fatal("expected slide content")
	}
	tblPrPattern := regexp.MustCompile(`(?s)<a:tblPr\b([^>]*)>.*?</a:tblPr>`)
	updatedSlideXML := tblPrPattern.ReplaceAllString(string(slideXML), `<a:tblPr$1/>`)
	if updatedSlideXML == string(slideXML) {
		t.Fatal("expected table XML fixture replacement to apply")
	}
	e.parts.Set(partPath, []byte(updatedSlideXML))

	styleGUID := "{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}"
	if err := e.SetTableStyle(0, shapeID, styleGUID); err != nil {
		t.Fatalf("SetTableStyle failed for self-closing tblPr: %v", err)
	}
	slideAfter, ok := e.parts.Get(partPath)
	if !ok {
		t.Fatal("expected updated slide content")
	}
	if !strings.Contains(
		string(slideAfter),
		`<a:tblPr firstRow="1" bandRow="1"><a:tableStyleId>{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}</a:tableStyleId></a:tblPr>`,
	) {
		t.Fatalf("expected style id inserted into expanded tblPr, got: %s", string(slideAfter))
	}
}

func TestGetTableIncludesTraversalViews(t *testing.T) {
	e := newTableEditorFixture()
	shapeID, err := e.AddTable(0, 2, 3, 0, 0, 0, 0)
	if err != nil {
		t.Fatalf("AddTable failed: %v", err)
	}
	if err := e.UpdateTableCellText(0, shapeID, 0, 0, "r0c0"); err != nil {
		t.Fatalf("UpdateTableCellText [0,0] failed: %v", err)
	}
	if err := e.UpdateTableCellText(0, shapeID, 1, 2, "r1c2"); err != nil {
		t.Fatalf("UpdateTableCellText [1,2] failed: %v", err)
	}
	if err := e.MergeTableCells(0, shapeID, 0, 1, 1, 2); err != nil {
		t.Fatalf("MergeTableCells failed: %v", err)
	}

	tableInfo, err := e.GetTable(0, shapeID)
	if err != nil {
		t.Fatalf("GetTable failed: %v", err)
	}
	meta := tableInfo["table"].(map[string]any)

	rows := meta["rows"].([]map[string]any)
	cols := meta["columns"].([]map[string]any)
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows in traversal view, got %d", len(rows))
	}
	if len(cols) != 3 {
		t.Fatalf("expected 3 columns in traversal view, got %d", len(cols))
	}

	row0cells := rows[0]["cells"].([]map[string]any)
	col2cells := cols[2]["cells"].([]map[string]any)
	if len(row0cells) != 3 {
		t.Fatalf("expected 3 cells in row[0], got %d", len(row0cells))
	}
	if len(col2cells) != 2 {
		t.Fatalf("expected 2 cells in column[2], got %d", len(col2cells))
	}

	origin := findTableCell(row0cells, 0, 1)
	if origin == nil || !origin["is_merge_origin"].(bool) {
		t.Fatalf("expected [0,1] to be merge origin in row traversal")
	}
	spanned := findTableCell(col2cells, 1, 2)
	if spanned == nil || !spanned["is_spanned"].(bool) {
		t.Fatalf("expected [1,2] to be spanned in column traversal")
	}
}

func TestTableFlagsRoundTripInGetTable(t *testing.T) {
	e := newTableEditorFixture()
	shapeID, err := e.AddTable(0, 2, 2, 0, 0, 0, 0)
	if err != nil {
		t.Fatalf("AddTable failed: %v", err)
	}
	flags := map[string]any{
		"first_row": true,
		"first_col": true,
		"last_row":  true,
		"last_col":  true,
		"band_row":  false,
		"band_col":  false,
	}
	if err := e.UpdateTableFlags(0, shapeID, flags); err != nil {
		t.Fatalf("UpdateTableFlags failed: %v", err)
	}

	tableInfo, err := e.GetTable(0, shapeID)
	if err != nil {
		t.Fatalf("GetTable failed: %v", err)
	}
	meta := tableInfo["table"].(map[string]any)
	if !meta["first_row"].(bool) || !meta["first_col"].(bool) || !meta["last_row"].(bool) || !meta["last_col"].(bool) {
		t.Fatalf("expected first/last row/col flags to be true, got: %#v", meta)
	}
	if meta["band_row"].(bool) || meta["band_col"].(bool) {
		t.Fatalf("expected band row/col flags to be false, got: %#v", meta)
	}
}

func newTableEditorFixture() *PresentationEditor {
	ps := NewPartStore()
	ps.Set(
		"ppt/slides/slide1.xml",
		[]byte(
			`<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"><p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/></p:nvGrpSpPr><p:grpSpPr/></p:spTree></p:cSld></p:sld>`,
		),
	)
	return &PresentationEditor{
		parts: ps,
		slides: []common.EditorSlideRef{{
			SlideID: 256,
			Part:    "ppt/slides/slide1.xml",
		}},
	}
}

func findTableCell(cells []map[string]any, row, col int) map[string]any {
	for _, c := range cells {
		if c["row"].(int) == row && c["col"].(int) == col {
			return c
		}
	}
	return nil
}
