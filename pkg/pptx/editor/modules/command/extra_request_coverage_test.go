package command

import (
	"encoding/base64"
	"strings"
	"testing"
)

func testParseFloatField(payload map[string]any, key string) (float64, bool) {
	v, ok := payload[key].(float64)
	return v, ok
}

func testOptionalBoolField(payload map[string]any, key string) (bool, bool) {
	v, ok := payload[key].(bool)
	return v, ok
}

func TestContentRequestParsers(t *testing.T) {
	if _, ok := ParseFindReplaceRequest(map[string]any{"find": "a", "replace": "b"}, testParseStringField); !ok {
		t.Fatal("ParseFindReplaceRequest should succeed")
	}
	if _, ok := ParseAuthorAddRequest(map[string]any{"name": "A", "initials": "AA"}, testParseStringField); !ok {
		t.Fatal("ParseAuthorAddRequest should succeed")
	}
	if _, ok := ParseCommentAddRequest(
		map[string]any{"slide_index": 1, "author_id": int64(2), "text": "Hi", "x": int64(10), "y": int64(20)},
		testParseSlideIndex,
		testParseInt64Field,
		testParseStringField,
	); !ok {
		t.Fatal("ParseCommentAddRequest should succeed")
	}
	if _, ok := ParseCommentRemoveRequest(
		map[string]any{"slide_index": 1, "author_id": int64(2), "author_index": 3},
		testParseSlideIndex,
		testParseInt64Field,
		testParseIntField,
	); !ok {
		t.Fatal("ParseCommentRemoveRequest should succeed")
	}
	if req, ok := ParseSetModifyPasswordRequest(map[string]any{"password": "pw"}, testParseStringField); !ok || req.Password != "pw" {
		t.Fatalf("ParseSetModifyPasswordRequest failed: req=%+v ok=%v", req, ok)
	}
	if final, ok := ParseSetMarkAsFinalRequest(map[string]any{"final": true}, testOptionalBoolField); !ok || !final {
		t.Fatalf("ParseSetMarkAsFinalRequest failed: final=%v ok=%v", final, ok)
	}
}

func TestCustomXMLAndBase64Parsing(t *testing.T) {
	errs := make([]string, 0)
	addErr := func(code, message string) {
		errs = append(errs, code+":"+message)
	}

	req := ParseCustomXMLAddRequest(
		map[string]any{
			"content":      "",
			"root_element": "Root",
			"namespace":    "urn:demo",
			"properties": map[string]any{
				"a": "1",
				"b": 2,
			},
		},
		testOptionalStringField,
		addErr,
		"MISSING_FIELD",
		"INVALID_TYPE",
	)
	if req.RootElement != "Root" || req.Namespace != "urn:demo" || req.Properties["a"] != "1" {
		t.Fatalf("unexpected ParseCustomXMLAddRequest result: %+v", req)
	}
	if len(errs) == 0 {
		t.Fatal("expected invalid non-string property to produce validation error")
	}

	errs = errs[:0]
	_ = ParseCustomXMLAddRequest(
		map[string]any{},
		testOptionalStringField,
		addErr,
		"MISSING_FIELD",
		"INVALID_TYPE",
	)
	if len(errs) == 0 || !strings.Contains(errs[0], "MISSING_FIELD") {
		t.Fatalf("expected missing-field validation error, got: %v", errs)
	}

	encoded := base64.StdEncoding.EncodeToString([]byte("abc"))
	data, ok, err := DecodeRequiredBase64Field(
		map[string]any{"data": encoded},
		testParseStringField,
		"data",
		"bad base64",
	)
	if err != nil || !ok || string(data) != "abc" {
		t.Fatalf("DecodeRequiredBase64Field failed: data=%q ok=%v err=%v", string(data), ok, err)
	}
	if _, ok, err = DecodeRequiredBase64Field(map[string]any{}, testParseStringField, "data", "bad"); err != nil || ok {
		t.Fatalf("missing base64 field should be ignored: ok=%v err=%v", ok, err)
	}
	if _, _, err = DecodeRequiredBase64Field(map[string]any{"data": "%%%bad"}, testParseStringField, "data", "bad"); err == nil {
		t.Fatal("expected invalid base64 error")
	}
}

func TestParityAndShapeUpdateParsers(t *testing.T) {
	if req, ok := ParseGroupShapeRequest(
		map[string]any{"slide_index": 1, "shapes": []int{1, 2}},
		testParseIntField,
		testParseIntSliceField,
	); !ok || len(req.ShapeIDs) != 2 {
		t.Fatalf("ParseGroupShapeRequest with shapes failed: req=%+v ok=%v", req, ok)
	}
	if req, ok := ParseGroupShapeRequest(map[string]any{"slide_index": 2}, testParseIntField, testParseIntSliceField); !ok || len(req.ShapeIDs) != 0 {
		t.Fatalf("ParseGroupShapeRequest without shapes failed: req=%+v ok=%v", req, ok)
	}

	points, err := ParseFreeformPoints(map[string]any{"points": []any{[]any{1, 2}, []any{3.0, 4.0}}})
	if err != nil || len(points) != 2 || points[1].X != 3 {
		t.Fatalf("ParseFreeformPoints failed: points=%+v err=%v", points, err)
	}
	if _, err = ParseFreeformPoints(map[string]any{"points": []any{[]any{1, "bad"}}}); err == nil {
		t.Fatal("expected ParseFreeformPoints coordinate type error")
	}

	closeFlag, err := ParseOptionalCloseFlag(map[string]any{"close": true})
	if err != nil || !closeFlag {
		t.Fatalf("ParseOptionalCloseFlag failed: close=%v err=%v", closeFlag, err)
	}
	if _, err = ParseOptionalCloseFlag(map[string]any{"close": "yes"}); err == nil {
		t.Fatal("expected close flag type error")
	}

	if _, ok := ParseTextboxPlacementRequest(
		map[string]any{"slide_index": 1, "left": 1.0, "top": 2.0, "width": 3.0, "height": 4.0},
		testParseSlideIndex,
		testParseFloatField,
	); !ok {
		t.Fatal("ParseTextboxPlacementRequest should succeed")
	}
	if _, ok := ParseConnectorPlacementRequest(
		map[string]any{
			"slide_index":    1,
			"connector_type": "straight",
			"begin_x":        1.0, "begin_y": 2.0, "end_x": 3.0, "end_y": 4.0,
		},
		testParseSlideIndex,
		testParseStringField,
		testParseFloatField,
	); !ok {
		t.Fatal("ParseConnectorPlacementRequest should succeed")
	}

	updates, has, err := ParseOptionalShapeUpdates(map[string]any{
		"text":        "hello",
		"runs":        []map[string]any{{"text": "x"}},
		"text_frame":  map[string]any{"columns": 2},
		"click_action": map[string]any{"address": "https://example.com"},
	})
	if err != nil || !has || updates.Text == nil || *updates.Text != "hello" {
		t.Fatalf("ParseOptionalShapeUpdates failed: updates=%+v has=%v err=%v", updates, has, err)
	}
	if _, has, err = ParseOptionalShapeUpdates(map[string]any{}); err != nil || has {
		t.Fatalf("expected empty optional update payload: has=%v err=%v", has, err)
	}
	if _, _, err = ParseOptionalShapeUpdates(map[string]any{"text": 123}); err == nil {
		t.Fatal("expected invalid text type error")
	}

	dst := map[string]any{}
	CopyShapeUpdateFields(map[string]any{"text": "a", "runs": []any{"x"}, "unknown": 1}, dst)
	if _, ok := dst["text"]; !ok {
		t.Fatal("CopyShapeUpdateFields should copy known text field")
	}
	if _, ok := dst["unknown"]; ok {
		t.Fatal("CopyShapeUpdateFields should ignore unknown fields")
	}
}

func TestMediaInsertHelpersAndAdapters(t *testing.T) {
	raw := base64.StdEncoding.EncodeToString([]byte("bin"))
	decoded, err := DecodeOptionalBase64Field(raw, 100, "audio")
	if err != nil || string(decoded) != "bin" {
		t.Fatalf("DecodeOptionalBase64Field failed: decoded=%q err=%v", string(decoded), err)
	}
	if _, err = DecodeOptionalBase64Field(raw, 2, "audio"); err == nil {
		t.Fatal("expected max length error")
	}
	if _, err = DecodeOptionalBase64Field("%%%", 100, "audio"); err == nil {
		t.Fatal("expected invalid base64 error")
	}

	binaryCalled := false
	shapeID, err := InsertShapeFromBinaryOrPath(true, func() (int, error) {
		binaryCalled = true
		return 7, nil
	}, func() (int, error) { return 0, nil })
	if err != nil || !binaryCalled || shapeID != 7 {
		t.Fatalf("InsertShapeFromBinaryOrPath(binary) failed: id=%d called=%v err=%v", shapeID, binaryCalled, err)
	}

	placement := MediaPlacement{SlideIndex: 1, X: 1, Y: 2, W: 3, H: 4}
	videoAdapter := AdaptVideoBinaryInsert(func(idx int, video, poster []byte, mime string, x, y, w, h float64) (int, error) {
		if idx != 1 || mime != "video/mp4" || x != 1 || h != 4 {
			t.Fatalf("AdaptVideoBinaryInsert passed wrong args: idx=%d mime=%q x=%v h=%v", idx, mime, x, h)
		}
		return 9, nil
	})
	if _, err = videoAdapter(placement, "video/mp4", []byte{1}, []byte{2}); err != nil {
		t.Fatalf("AdaptVideoBinaryInsert execution failed: %v", err)
	}

	audioAdapter := AdaptAudioPathInsertWithOptionalIcon(
		func(int, string, string, float64, float64, float64, float64) (int, error) { return 10, nil },
		func(idx int, path, icon, mime string, _, _, _, _ float64) (int, error) {
			if idx != 1 || icon != "icon.png" || mime != "audio/mpeg" {
				t.Fatalf("AdaptAudioPathInsertWithOptionalIcon wrong args idx=%d icon=%q mime=%q", idx, icon, mime)
			}
			return 11, nil
		},
	)
	if _, err = audioAdapter(placement, "audio/mpeg", "audio.mp3", "icon.png"); err != nil {
		t.Fatalf("AdaptAudioPathInsertWithOptionalIcon failed: %v", err)
	}

	spec := NewVideoInsertSpec(100, func(MediaPlacement, string, []byte, []byte) (int, error) { return 20, nil }, func(MediaPlacement, string, string, string) (int, error) { return 21, nil })
	if spec.MetaKey != "mime_type" || spec.PrimaryLabel != "video" || spec.SecondaryLabel != "poster" {
		t.Fatalf("NewVideoInsertSpec unexpected defaults: %+v", spec)
	}
}
