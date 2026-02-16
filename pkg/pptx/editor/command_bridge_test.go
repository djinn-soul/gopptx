package editor

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestCommandShapeOps(t *testing.T) {
	basePath := writeDeckFixture(t, "bridge-shape-test.pptx", []elements.SlideContent{
		elements.NewSlide("Shape Test").AddBullet("body"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = e.Close() }()

	// 1. Add Shape
	addReq := `{"api_version":1,"request_id":"r1","op":"add_shape","payload":{"slide_index":0,"type":"rect","x":100,"y":100,"w":1000,"h":500}}`
	resp := ExecuteCommand(e, addReq)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("add_shape failed: %s", resp)
	}

	var addOut struct {
		Result struct {
			ShapeID int `json:"shape_id"`
		} `json:"result"`
	}
	if unmarshalErr := json.Unmarshal([]byte(resp), &addOut); unmarshalErr != nil {
		t.Fatalf("unmarshal add response: %v", unmarshalErr)
	}
	shapeID := addOut.Result.ShapeID
	if shapeID == 0 {
		t.Fatalf("expected valid shape_id, got 0")
	}

	// 2. List Shapes — should contain the new shape
	listReq := `{"api_version":1,"request_id":"r2","op":"list_shapes","payload":{"slide_index":0}}`
	resp = ExecuteCommand(e, listReq)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("list_shapes failed: %s", resp)
	}

	// 3. Update Shape
	updateReq := fmt.Sprintf(
		`{"api_version":1,"request_id":"r3","op":"update_shape","payload":{"slide_index":0,"shape_id":%d,"updates":{"text":"Updated"}}}`,
		shapeID,
	)
	resp = ExecuteCommand(e, updateReq)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("update_shape failed: %s", resp)
	}

	// 4. Remove Shape
	removeReq := fmt.Sprintf(
		`{"api_version":1,"request_id":"r4","op":"remove_shape","payload":{"slide_index":0,"shape_id":%d}}`,
		shapeID,
	)
	resp = ExecuteCommand(e, removeReq)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("remove_shape failed: %s", resp)
	}
}

func TestCommandNotesOps(t *testing.T) {
	basePath := writeDeckFixture(t, "bridge-notes-test.pptx", []elements.SlideContent{
		elements.NewSlide("Notes Test").AddBullet("body"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = e.Close() }()

	// 1. Set Notes
	setReq := `{"api_version":1,"request_id":"n1","op":"set_notes","payload":{"slide_index":0,"text":"Speaker notes here"}}`
	resp := ExecuteCommand(e, setReq)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("set_notes failed: %s", resp)
	}

	// 2. Get Notes
	getReq := `{"api_version":1,"request_id":"n2","op":"get_notes","payload":{"slide_index":0}}`
	resp = ExecuteCommand(e, getReq)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("get_notes failed: %s", resp)
	}
	if !strings.Contains(resp, "Speaker notes here") {
		t.Fatalf("get_notes missing expected text: %s", resp)
	}

	// 3. Update Notes
	setReq2 := `{"api_version":1,"request_id":"n3","op":"set_notes","payload":{"slide_index":0,"text":"Updated notes"}}`
	resp = ExecuteCommand(e, setReq2)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("set_notes update failed: %s", resp)
	}

	// 4. Verify Update
	resp = ExecuteCommand(e, getReq)
	if !strings.Contains(resp, "Updated notes") {
		t.Fatalf("get_notes mismatch after update: %s", resp)
	}
}
