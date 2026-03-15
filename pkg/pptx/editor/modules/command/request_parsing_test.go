package command

import (
	"strings"
	"testing"
)

func testParseSlideIndex(payload map[string]any) (int, bool) {
	v, ok := payload["slide_index"].(int)
	return v, ok
}

func testParseIntField(payload map[string]any, key string) (int, bool) {
	v, ok := payload[key].(int)
	return v, ok
}

func testParseInt64Field(payload map[string]any, key string) (int64, bool) {
	v, ok := payload[key].(int64)
	return v, ok
}

func testParseStringField(payload map[string]any, key string) (string, bool) {
	v, ok := payload[key].(string)
	return v, ok
}

func testOptionalStringField(payload map[string]any, key string) string {
	v, _ := payload[key].(string)
	return v
}

func testOptionalStringSliceField(payload map[string]any, key string) ([]string, bool) {
	v, ok := payload[key].([]string)
	return v, ok
}

func testParseIntSliceField(payload map[string]any, key string) ([]int, bool) {
	v, ok := payload[key].([]int)
	return v, ok
}

func testParseStringSliceField(payload map[string]any, key string) ([]string, bool) {
	v, ok := payload[key].([]string)
	return v, ok
}

func testParseFloat64SliceField(payload map[string]any, key string) ([]float64, bool) {
	v, ok := payload[key].([]float64)
	return v, ok
}

func testOptionalIntField(payload map[string]any, key string) (int, bool) {
	v, ok := payload[key].(int)
	return v, ok
}

func testOptionalInt64Field(payload map[string]any, key string) (int64, bool) {
	v, ok := payload[key].(int64)
	return v, ok
}

func TestValidationErrorAndOptionalPayloadDecode(t *testing.T) {
	err := NewValidationError("bad_request", "failed")
	if err.Code != "bad_request" || err.Error() != "failed" {
		t.Fatalf("unexpected validation error: %+v", err)
	}

	type target struct {
		Name string `json:"name"`
	}
	decoded := target{}
	if derr := DecodeOptionalPayloadValue(map[string]any{"data": map[string]any{"name": "ok"}}, "data", &decoded); derr != nil {
		t.Fatalf("DecodeOptionalPayloadValue failed: %v", derr)
	}
	if decoded.Name != "ok" {
		t.Fatalf("expected decoded name ok, got %q", decoded.Name)
	}
	if derr := DecodeOptionalPayloadValue(map[string]any{}, "data", &decoded); derr != nil {
		t.Fatalf("missing optional field should not error: %v", derr)
	}

	decoded = target{}
	derr := DecodeOptionalPayloadValue(map[string]any{"data": map[string]any{"name": 123}}, "data", &decoded)
	if derr == nil || !strings.Contains(derr.Error(), "invalid data payload") {
		t.Fatalf("expected payload decode error, got: %v", derr)
	}
}

func TestSlideAndTableRequestParsing(t *testing.T) {
	addReq := ParseAddSlideRequest(
		map[string]any{"title": "T", "layout": "L", "bullets": []string{"a", "b"}},
		testOptionalStringField,
		testOptionalStringSliceField,
	)
	if addReq.Title != "T" || addReq.Layout != "L" || len(addReq.Bullets) != 2 {
		t.Fatalf("unexpected add slide request: %+v", addReq)
	}

	if _, ok := ParseSlideIndexOnlyRequest(map[string]any{"index": 2}, testParseIntField); !ok {
		t.Fatal("ParseSlideIndexOnlyRequest should succeed")
	}
	if _, ok := ParseMoveSlideRequest(map[string]any{"from": 0, "to": 3}, testParseIntField); !ok {
		t.Fatal("ParseMoveSlideRequest should succeed")
	}
	dup, ok := ParseDuplicateSlideRequest(
		map[string]any{"index": 1, "insert_at": 5},
		testParseIntField,
		testOptionalIntField,
	)
	if !ok || dup.InsertAt != 5 {
		t.Fatalf("unexpected duplicate request: %+v ok=%v", dup, ok)
	}

	tableReq, ok := ParseTableAddRequest(
		map[string]any{
			"slide_index": 1,
			"rows":        2,
			"cols":        3,
			"x":           int64(10),
			"y":           int64(20),
			"cx":          int64(30),
			"cy":          int64(40),
		},
		testParseSlideIndex,
		testParseIntField,
		testParseInt64Field,
	)
	if !ok || tableReq.Rows != 2 || tableReq.CY != 40 {
		t.Fatalf("unexpected table add request: %+v ok=%v", tableReq, ok)
	}

	if err := ValidateTableDimensions(0, 2, 10); err == nil {
		t.Fatal("expected rows range error")
	}
	if err := ValidateTableDimensions(2, 11, 10); err == nil {
		t.Fatal("expected cols range error")
	}
	if err := ValidateTableDimensions(2, 3, 10); err != nil {
		t.Fatalf("expected valid dimensions: %v", err)
	}

	shapeReq, ok := ParseTableShapeRequest(map[string]any{"slide_index": 1, "shape_id": 7}, testParseIntField)
	if !ok || shapeReq.ShapeID != 7 {
		t.Fatalf("unexpected table shape request: %+v ok=%v", shapeReq, ok)
	}
	if _, ok = ParseTableCellRangeRequest(
		map[string]any{"slide_index": 1, "shape_id": 7, "row1": 0, "col1": 0, "row2": 1, "col2": 1},
		testParseIntField,
	); !ok {
		t.Fatal("ParseTableCellRangeRequest should succeed")
	}
	if _, ok = ParseTableCellRequest(
		map[string]any{"slide_index": 1, "shape_id": 7, "row": 0, "col": 1},
		testParseIntField,
	); !ok {
		t.Fatal("ParseTableCellRequest should succeed")
	}

	styleReq, ok := ParseTableStyleRequest(
		map[string]any{"slide_index": 1, "shape_id": 7, "style_guid": "{GUID}"},
		testParseIntField,
		testParseStringField,
	)
	if !ok || styleReq.StyleGUID != "{GUID}" {
		t.Fatalf("unexpected table style request: %+v ok=%v", styleReq, ok)
	}
}

func TestObjectAndTextUpdateParsing(t *testing.T) {
	updates, err := ParseRequiredObjectField(
		map[string]any{"updates": map[string]any{"text": "hello"}},
		"updates",
		"missing",
		"invalid type",
	)
	if err != nil || updates["text"] != "hello" {
		t.Fatalf("unexpected updates parse result: %+v err=%v", updates, err)
	}
	if _, err = ParseRequiredObjectField(map[string]any{}, "updates", "missing", "invalid"); err == nil {
		t.Fatal("expected missing object field error")
	}
	if _, err = ParseRequiredObjectField(map[string]any{"updates": 12}, "updates", "missing", "invalid"); err == nil {
		t.Fatal("expected object type error")
	}

	text, ok, err := ParseOptionalTextUpdate(map[string]any{"text": "hi"})
	if err != nil || !ok || text != "hi" {
		t.Fatalf("unexpected optional text update parse: text=%q ok=%v err=%v", text, ok, err)
	}
	_, _, err = ParseOptionalTextUpdate(map[string]any{"text": 7})
	if err == nil {
		t.Fatal("expected invalid text update type error")
	}
}

func TestSectionAndChartRequestParsing(t *testing.T) {
	sectionReq, ok := ParseSectionAddRequest(
		map[string]any{"name": "S1", "slide_indices": []int{0, 2}},
		testParseStringField,
		testParseIntSliceField,
	)
	if !ok || sectionReq.Name != "S1" || len(sectionReq.SlideIndices) != 2 {
		t.Fatalf("unexpected section add request: %+v ok=%v", sectionReq, ok)
	}
	if _, ok = ParseSectionNameRequest(map[string]any{"name": "S1"}, testParseStringField); !ok {
		t.Fatal("ParseSectionNameRequest should succeed")
	}
	if _, ok = ParseSectionRenameRequest(
		map[string]any{"old_name": "S1", "new_name": "S2"},
		testParseStringField,
	); !ok {
		t.Fatal("ParseSectionRenameRequest should succeed")
	}
	if _, ok = ParseSlideSizeRequest(
		map[string]any{"width": int64(9144000), "height": int64(6858000)},
		testParseInt64Field,
	); !ok {
		t.Fatal("ParseSlideSizeRequest should succeed")
	}
	if _, ok = ParseSlideTitleRequest(
		map[string]any{"slide_index": 2, "title": "My Title"},
		testParseSlideIndex,
		testParseStringField,
	); !ok {
		t.Fatal("ParseSlideTitleRequest should succeed")
	}
	if _, ok = ParseMergeFromFileRequest(
		map[string]any{"path": "source.pptx"},
		testParseStringField,
	); !ok {
		t.Fatal("ParseMergeFromFileRequest should succeed")
	}

	core := ParseCorePropertiesRequest(
		map[string]any{"title": "Deck", "creator": "author"},
		testOptionalStringField,
	)
	if core.Title != "Deck" || core.Creator != "author" {
		t.Fatalf("unexpected core properties: %+v", core)
	}

	updateReq, ok := ParseUpdateSlideRequest(
		map[string]any{"slide_index": 2, "title": "T", "layout": "L", "bullets": []string{"one"}},
		testParseSlideIndex,
		testOptionalStringField,
		testOptionalStringSliceField,
	)
	if !ok || updateReq.SlideIndex != 2 || len(updateReq.Bullets) != 1 {
		t.Fatalf("unexpected update request: %+v ok=%v", updateReq, ok)
	}

	addChartReq, ok := ParseAddChartRequest(
		map[string]any{
			"slide_index": 1,
			"chart_type":  "bar",
			"title":       "Sales",
			"categories":  []string{"Q1", "Q2"},
			"values":      []float64{10, 20},
			"x":           int64(1),
			"y":           int64(2),
			"w":           int64(3),
			"h":           int64(4),
		},
		testParseSlideIndex,
		testParseStringField,
		testOptionalStringField,
		testParseStringSliceField,
		testParseFloat64SliceField,
		testOptionalInt64Field,
	)
	if !ok || addChartReq.ChartType != "bar" || addChartReq.H != 4 {
		t.Fatalf("unexpected add chart request: %+v ok=%v", addChartReq, ok)
	}
}

func TestNotesPayloadParsers(t *testing.T) {
	if _, ok := ParseSlideIndexRequest(map[string]any{"slide_index": 1}, testParseSlideIndex); !ok {
		t.Fatal("ParseSlideIndexRequest should succeed")
	}
	if _, ok := ParseSlideShapeRequest(
		map[string]any{"slide_index": 1, "shape_id": 9},
		testParseSlideIndex,
		testParseIntField,
	); !ok {
		t.Fatal("ParseSlideShapeRequest should succeed")
	}
	if _, ok := ParseSlideShapeIDsRequest(
		map[string]any{"slide_index": 1, "shape_ids": []int{1, 2}},
		testParseSlideIndex,
		testParseIntSliceField,
	); !ok {
		t.Fatal("ParseSlideShapeIDsRequest should succeed")
	}

	if _, _, ok := ParseSetNotesRequest(
		map[string]any{"slide_index": 1, "text": "speaker note"},
		testParseSlideIndex,
		testParseStringField,
	); !ok {
		t.Fatal("ParseSetNotesRequest should succeed")
	}
	if _, ok := ParseSetNotesShapeTextRequest(
		map[string]any{"slide_index": 1, "shape_id": 3, "text": "hello"},
		testParseSlideIndex,
		testParseIntField,
		testParseStringField,
	); !ok {
		t.Fatal("ParseSetNotesShapeTextRequest should succeed")
	}

	minimal := BuildNotesResult("n1", false)
	if minimal["text"] != "n1" || minimal["notes_slide"] != nil {
		t.Fatalf("unexpected BuildNotesResult payload: %+v", minimal)
	}
	detailed := BuildNotesResultDetailed("n2", true, nil, nil)
	if detailed["text"] != "n2" {
		t.Fatalf("unexpected BuildNotesResultDetailed text: %+v", detailed)
	}
	if detailed["notes_slide"] == nil {
		t.Fatalf("expected notes_slide details when hasNotesSlide=true: %+v", detailed)
	}
}
