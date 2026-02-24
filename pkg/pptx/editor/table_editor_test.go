package editor

import (
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
