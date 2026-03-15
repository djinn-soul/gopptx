package editor

import "testing"

func TestTableResizeOps(t *testing.T) {
	ed := newTableEditorFixture()
	sid := 0

	shapeID, err := ed.AddTable(sid, 2, 2, 100, 100, 4000, 2000)
	if err != nil {
		t.Fatalf("AddTable failed: %v", err)
	}

	err = ed.SetTableRowHeight(sid, shapeID, 0, 1500)
	if err != nil {
		t.Errorf("SetTableRowHeight failed: %v", err)
	}

	err = ed.SetTableColumnWidth(sid, shapeID, 0, 2500)
	if err != nil {
		t.Errorf("SetTableColumnWidth failed: %v", err)
	}

	// Read back to verify
	info, err := ed.GetTable(sid, shapeID)
	if err != nil {
		t.Fatalf("GetTable failed: %v", err)
	}

	tableInfo, ok := info["table"].(map[string]any)
	if !ok {
		t.Fatalf("expected table metadata in GetTable response")
	}

	if h, ok := tableInfo["row_heights"].([]int64); ok && len(h) > 0 {
		if h[0] != 1500 {
			t.Errorf("expected row height 1500, got %d", h[0])
		}
	} else {
		t.Errorf("failed to retrieve row height")
	}

	if w, ok := tableInfo["column_widths"].([]int64); ok && len(w) > 0 {
		if w[0] != 2500 {
			t.Errorf("expected col width 2500, got %d", w[0])
		}
	} else {
		t.Errorf("failed to retrieve col width")
	}
	
	// Error cases
	if err := ed.SetTableRowHeight(sid, 9999, 0, 1500); err == nil {
		t.Error("expected error for invalid shapeID")
	}
	if err := ed.SetTableColumnWidth(sid, 9999, 0, 2500); err == nil {
		t.Error("expected error for invalid shapeID")
	}
}
