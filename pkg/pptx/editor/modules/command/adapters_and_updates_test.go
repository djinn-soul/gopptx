package command

import (
	"encoding/base64"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func TestMediaInsertAdaptersAndSpecs(t *testing.T) {
	placement := MediaPlacement{SlideIndex: 2, X: 1, Y: 2, W: 3, H: 4}

	videoCalled := false
	videoAdapter := AdaptVideoBinaryInsert(
		func(slideIndex int, videoData, posterData []byte, mimeType string, x, y, w, h float64) (int, error) {
			videoCalled = true
			if slideIndex != 2 || mimeType != "video/mp4" || len(videoData) != 1 ||
				len(posterData) != 1 {
				t.Fatalf("unexpected video adapter args")
			}
			if x != 1 || y != 2 || w != 3 || h != 4 {
				t.Fatalf("unexpected placement args")
			}
			return 101, nil
		},
	)
	id, err := videoAdapter(placement, "video/mp4", []byte{1}, []byte{2})
	if err != nil || id != 101 || !videoCalled {
		t.Fatalf("video adapter failed: id=%d err=%v called=%v", id, err, videoCalled)
	}

	audioPlainCalled := false
	audioIconCalled := false
	audioWithIcon := AdaptAudioBinaryInsertWithOptionalIcon(
		func(int, []byte, string, float64, float64, float64, float64) (int, error) {
			audioPlainCalled = true
			return 201, nil
		},
		func(int, []byte, []byte, string, float64, float64, float64, float64) (int, error) {
			audioIconCalled = true
			return 202, nil
		},
	)
	id, err = audioWithIcon(placement, "audio/mpeg", []byte{1}, nil)
	if err != nil || id != 201 || !audioPlainCalled || audioIconCalled {
		t.Fatalf(
			"audio adapter (plain) failed: id=%d err=%v plain=%v icon=%v",
			id,
			err,
			audioPlainCalled,
			audioIconCalled,
		)
	}
	audioPlainCalled = false
	audioIconCalled = false
	id, err = audioWithIcon(placement, "audio/mpeg", []byte{1}, []byte{9})
	if err != nil || id != 202 || audioPlainCalled || !audioIconCalled {
		t.Fatalf(
			"audio adapter (icon) failed: id=%d err=%v plain=%v icon=%v",
			id,
			err,
			audioPlainCalled,
			audioIconCalled,
		)
	}

	spec := NewVideoInsertSpec(4096, nil, nil)
	if spec.PrimaryLabel != "video" || spec.SecondaryLabel != "poster" ||
		spec.PrimaryMaxLen != 4096 {
		t.Fatalf("unexpected video spec: %+v", spec)
	}
	spec = NewAudioInsertSpec(2048, nil, nil)
	if spec.PrimaryLabel != "audio" || spec.SecondaryLabel != "icon" {
		t.Fatalf("unexpected audio spec: %+v", spec)
	}
	spec = NewOLEInsertSpec(1024, nil, nil)
	if spec.MetaKey != "prog_id" || spec.PrimaryLabel != "object" {
		t.Fatalf("unexpected ole spec: %+v", spec)
	}
}

func TestShapeUpdateAndAddHelpers(t *testing.T) {
	updates, has, err := ParseOptionalShapeUpdates(map[string]any{
		"text": "hello",
		"properties": map[string]any{
			"x": 10,
		},
	})
	if err != nil || !has || updates.Text == nil || *updates.Text != "hello" || updates.X == nil ||
		*updates.X != 10 {
		t.Fatalf("ParseOptionalShapeUpdates failed: updates=%+v has=%v err=%v", updates, has, err)
	}
	if _, _, err = ParseOptionalShapeUpdates(map[string]any{"text": 99}); err == nil {
		t.Fatal("expected text type validation error")
	}
	if _, has, err = ParseOptionalShapeUpdates(map[string]any{}); err != nil || has {
		t.Fatalf("empty payload should produce no updates: has=%v err=%v", has, err)
	}

	dst := map[string]any{"keep": 1}
	CopyShapeUpdateFields(map[string]any{"text": "x", "runs": []any{}, "unknown": true}, dst)
	if _, ok := dst["text"]; !ok {
		t.Fatalf("CopyShapeUpdateFields should copy known fields: %+v", dst)
	}
	if _, ok := dst["unknown"]; ok {
		t.Fatalf("CopyShapeUpdateFields should ignore unknown fields: %+v", dst)
	}

	req, ok := ParseAddShapeBase(
		map[string]any{
			"slide_index": 1,
			"type":        "rect",
			"x":           1.0,
			"y":           2.0,
			"w":           3.0,
			"h":           4.0,
			"text":        "t",
		},
		testParseSlideIndex,
		testParseStringField,
		func(payload map[string]any, key string) (float64, bool) {
			v, ok := payload[key].(float64)
			return v, ok
		},
		testOptionalStringField,
	)
	if !ok || req.ShapeType != "rect" || req.Text != "t" {
		t.Fatalf("ParseAddShapeBase failed: req=%+v ok=%v", req, ok)
	}
	req.Properties = common.ShapeUpdate{X: &[]int{99}[0]}
	shapeUpdate, hasUpdate := BuildShapeUpdateForAdd(req)
	if !hasUpdate || shapeUpdate.Text == nil || shapeUpdate.X == nil {
		t.Fatalf("BuildShapeUpdateForAdd failed: update=%+v has=%v", shapeUpdate, hasUpdate)
	}
	if !HasAnyUpdate(shapeUpdate) {
		t.Fatal("HasAnyUpdate should detect non-nil fields")
	}

	shapeID, err := ExecuteAddShapeRequest(
		req,
		func(int, string, float64, float64, float64, float64) (int, error) { return 55, nil },
		func(int, int, common.ShapeUpdate) error { return nil },
	)
	if err != nil || shapeID != 55 {
		t.Fatalf("ExecuteAddShapeRequest failed: shapeID=%d err=%v", shapeID, err)
	}
}

func TestContentRequestHelpers(t *testing.T) {
	if _, ok := ParseFindReplaceRequest(map[string]any{"find": "a", "replace": "b"}, testParseStringField); !ok {
		t.Fatal("ParseFindReplaceRequest should succeed")
	}
	if _, ok := ParseAuthorAddRequest(map[string]any{"name": "A", "initials": "AA"}, testParseStringField); !ok {
		t.Fatal("ParseAuthorAddRequest should succeed")
	}
	if _, ok := ParseCommentAddRequest(
		map[string]any{"slide_index": 1, "author_id": int64(2), "text": "x", "x": int64(3), "y": int64(4)},
		testParseSlideIndex,
		testParseInt64Field,
		testParseStringField,
	); !ok {
		t.Fatal("ParseCommentAddRequest should succeed")
	}
	if _, ok := ParseCommentRemoveRequest(
		map[string]any{"slide_index": 1, "author_id": int64(2), "author_index": 0},
		testParseSlideIndex,
		testParseInt64Field,
		testParseIntField,
	); !ok {
		t.Fatal("ParseCommentRemoveRequest should succeed")
	}
	if _, ok := ParseSetModifyPasswordRequest(map[string]any{"password": "pw"}, testParseStringField); !ok {
		t.Fatal("ParseSetModifyPasswordRequest should succeed")
	}
	if final, ok := ParseSetMarkAsFinalRequest(
		map[string]any{"final": true},
		func(payload map[string]any, key string) (bool, bool) {
			v, ok := payload[key].(bool)
			return v, ok
		},
	); !ok || !final {
		t.Fatalf("ParseSetMarkAsFinalRequest failed: final=%v ok=%v", final, ok)
	}

	var errs []string
	addErr := func(code, msg string) { errs = append(errs, code+":"+msg) }
	customReq := ParseCustomXMLAddRequest(
		map[string]any{
			"root_element": "root",
			"namespace":    "urn:x",
			"properties":   map[string]any{"a": "1", "b": 2},
		},
		testOptionalStringField,
		addErr,
		"MISSING",
		"INVALID",
	)
	if customReq.RootElement != "root" || customReq.Properties["a"] != "1" {
		t.Fatalf("ParseCustomXMLAddRequest failed: %+v", customReq)
	}
	if len(errs) == 0 {
		t.Fatal("expected validation errors for non-string property value")
	}

	b64 := base64.StdEncoding.EncodeToString([]byte("abc"))
	data, present, err := DecodeRequiredBase64Field(
		map[string]any{"blob": b64},
		testParseStringField,
		"blob",
		"bad",
	)
	if err != nil || !present || string(data) != "abc" {
		t.Fatalf(
			"DecodeRequiredBase64Field success failed: data=%q present=%v err=%v",
			string(data),
			present,
			err,
		)
	}
	_, present, err = DecodeRequiredBase64Field(
		map[string]any{},
		testParseStringField,
		"blob",
		"bad",
	)
	if err != nil || present {
		t.Fatalf("DecodeRequiredBase64Field missing field failed: present=%v err=%v", present, err)
	}
	_, present, err = DecodeRequiredBase64Field(
		map[string]any{"blob": "%%%invalid"},
		testParseStringField,
		"blob",
		"bad",
	)
	if err == nil || !present {
		t.Fatalf(
			"DecodeRequiredBase64Field invalid data should error with present=true: present=%v err=%v",
			present,
			err,
		)
	}
}
