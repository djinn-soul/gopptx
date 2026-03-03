package editor

import (
	"encoding/json"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// TestBridgeErrorCodes tests that the bridge returns proper error codes
// for various failure conditions.

func TestBridgeErrorCode_InvalidJSON(t *testing.T) {
	basePath := writeDeckFixture(t, "error-invalid-json.pptx", []elements.SlideContent{
		elements.NewSlide("Test"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer e.Close()

	response := ExecuteCommand(e, `{invalid json}`)
	var resp ResponseEnvelope
	if err := json.Unmarshal([]byte(response), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.OK {
		t.Fatal("expected ok=false for invalid JSON")
	}
	if resp.Error == nil {
		t.Fatal("expected error detail")
	}
	if resp.Error.Code != ErrCodeInvalidJSON {
		t.Errorf("expected code %s, got %s", ErrCodeInvalidJSON, resp.Error.Code)
	}
}

func TestBridgeErrorCode_UnknownOp(t *testing.T) {
	basePath := writeDeckFixture(t, "error-unknown-op.pptx", []elements.SlideContent{
		elements.NewSlide("Test"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer e.Close()

	req := `{"api_version": 1, "request_id": "test", "op": "nonexistent_op", "payload": {}}`
	response := ExecuteCommand(e, req)

	var resp ResponseEnvelope
	if err := json.Unmarshal([]byte(response), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.OK {
		t.Fatal("expected ok=false for unknown op")
	}
	if resp.Error == nil {
		t.Fatal("expected error detail")
	}
	if resp.Error.Code != ErrCodeUnknownOp {
		t.Errorf("expected code %s, got %s", ErrCodeUnknownOp, resp.Error.Code)
	}
}

func TestBridgeErrorCode_MissingField(t *testing.T) {
	basePath := writeDeckFixture(t, "error-missing-field.pptx", []elements.SlideContent{
		elements.NewSlide("Test"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer e.Close()

	// remove_slide requires "index" field
	req := `{"api_version": 1, "request_id": "test", "op": "remove_slide", "payload": {}}`
	response := ExecuteCommand(e, req)

	var resp ResponseEnvelope
	if err := json.Unmarshal([]byte(response), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.OK {
		t.Fatal("expected ok=false for missing field")
	}
	if resp.Error == nil {
		t.Fatal("expected error detail")
	}
	if resp.Error.Code != ErrCodeMissingField {
		t.Errorf("expected code %s, got %s", ErrCodeMissingField, resp.Error.Code)
	}
}

func TestBridgeErrorCode_InvalidIndex(t *testing.T) {
	basePath := writeDeckFixture(t, "error-invalid-index.pptx", []elements.SlideContent{
		elements.NewSlide("Test"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer e.Close()

	// Index 999 is out of bounds
	req := `{"api_version": 1, "request_id": "test", "op": "remove_slide", "payload": {"index": 999}}`
	response := ExecuteCommand(e, req)

	var resp ResponseEnvelope
	if err := json.Unmarshal([]byte(response), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.OK {
		t.Fatal("expected ok=false for invalid index")
	}
	if resp.Error == nil {
		t.Fatal("expected error detail")
	}
	if resp.Error.Code != ErrCodeInvalidIndex {
		t.Errorf("expected code %s, got %s", ErrCodeInvalidIndex, resp.Error.Code)
	}
}

func TestBridgeErrorCode_InvalidType(t *testing.T) {
	basePath := writeDeckFixture(t, "error-invalid-type.pptx", []elements.SlideContent{
		elements.NewSlide("Test"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer e.Close()

	// "index" should be int, not string
	req := `{"api_version": 1, "request_id": "test", "op": "remove_slide", "payload": {"index": "not-a-number"}}`
	response := ExecuteCommand(e, req)

	var resp ResponseEnvelope
	if err := json.Unmarshal([]byte(response), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.OK {
		t.Fatal("expected ok=false for invalid type")
	}
	if resp.Error == nil {
		t.Fatal("expected error detail")
	}
	if resp.Error.Code != ErrCodeInvalidType {
		t.Errorf("expected code %s, got %s", ErrCodeInvalidType, resp.Error.Code)
	}
}

func TestBridgeErrorCode_InvalidValue(t *testing.T) {
	basePath := writeDeckFixture(t, "error-invalid-value.pptx", []elements.SlideContent{
		elements.NewSlide("Test"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer e.Close()

	// Unknown theme name
	req := `{"api_version": 1, "request_id": "test", "op": "apply_theme", "payload": {"theme_name": "nonexistent_theme"}}`
	response := ExecuteCommand(e, req)

	var resp ResponseEnvelope
	if err := json.Unmarshal([]byte(response), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.OK {
		t.Fatal("expected ok=false for invalid value")
	}
	if resp.Error == nil {
		t.Fatal("expected error detail")
	}
	if resp.Error.Code != ErrCodeInvalidValue {
		t.Errorf("expected code %s, got %s", ErrCodeInvalidValue, resp.Error.Code)
	}
}

func TestBridgeErrorCode_UnsupportedVersion(t *testing.T) {
	basePath := writeDeckFixture(t, "error-unsupported-version.pptx", []elements.SlideContent{
		elements.NewSlide("Test"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer e.Close()

	req := `{"api_version": 99, "request_id": "test", "op": "slide_count", "payload": {}}`
	response := ExecuteCommand(e, req)

	var resp ResponseEnvelope
	if err := json.Unmarshal([]byte(response), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.OK {
		t.Fatal("expected ok=false for unsupported version")
	}
	if resp.Error == nil {
		t.Fatal("expected error detail")
	}
	if resp.Error.Code != ErrCodeUnsupportedVer {
		t.Errorf("expected code %s, got %s", ErrCodeUnsupportedVer, resp.Error.Code)
	}
}

func TestBridgeErrorCode_ShapeInvalidSlideIndex(t *testing.T) {
	basePath := writeDeckFixture(t, "error-shape-invalid-slide.pptx", []elements.SlideContent{
		elements.NewSlide("Test"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer e.Close()

	// Invalid slide_index for list_shapes
	req := `{"api_version": 1, "request_id": "test", "op": "list_shapes", "payload": {"slide_index": 999}}`
	response := ExecuteCommand(e, req)

	var resp ResponseEnvelope
	if err := json.Unmarshal([]byte(response), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.OK {
		t.Fatal("expected ok=false for invalid slide index")
	}
	if resp.Error == nil {
		t.Fatal("expected error detail")
	}
	if resp.Error.Code != ErrCodeInvalidIndex {
		t.Errorf("expected code %s, got %s", ErrCodeInvalidIndex, resp.Error.Code)
	}
}

func TestBridgeErrorCode_NotesInvalidSlideIndex(t *testing.T) {
	basePath := writeDeckFixture(t, "error-notes-invalid-slide.pptx", []elements.SlideContent{
		elements.NewSlide("Test"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer e.Close()

	// Invalid slide_index for get_notes
	req := `{"api_version": 1, "request_id": "test", "op": "get_notes", "payload": {"slide_index": 999}}`
	response := ExecuteCommand(e, req)

	var resp ResponseEnvelope
	if err := json.Unmarshal([]byte(response), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.OK {
		t.Fatal("expected ok=false for invalid slide index")
	}
	if resp.Error == nil {
		t.Fatal("expected error detail")
	}
	if resp.Error.Code != ErrCodeInvalidIndex {
		t.Errorf("expected code %s, got %s", ErrCodeInvalidIndex, resp.Error.Code)
	}
}

func TestBridgeErrorCode_AddShapeMissingFields(t *testing.T) {
	basePath := writeDeckFixture(t, "error-add-shape-missing.pptx", []elements.SlideContent{
		elements.NewSlide("Test"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer e.Close()

	// Missing required fields for add_shape
	req := `{"api_version": 1, "request_id": "test", "op": "add_shape", "payload": {"slide_index": 0}}`
	response := ExecuteCommand(e, req)

	var resp ResponseEnvelope
	if err := json.Unmarshal([]byte(response), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.OK {
		t.Fatal("expected ok=false for missing fields")
	}
	if resp.Error == nil {
		t.Fatal("expected error detail")
	}
	if resp.Error.Code != ErrCodeMissingField {
		t.Errorf("expected code %s, got %s", ErrCodeMissingField, resp.Error.Code)
	}
}

func TestBridgeErrorCode_BatchWithInvalidPayload(t *testing.T) {
	basePath := writeDeckFixture(t, "error-batch-invalid.pptx", []elements.SlideContent{
		elements.NewSlide("Test"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer e.Close()

	// Batch with invalid payload in one command (stop_on_error: false means batch returns ok=true)
	req := `{
		"api_version": 1,
		"request_id": "test",
		"op": "batch_execute",
		"payload": {
			"commands": [
				{"op": "remove_slide", "payload": {}}
			],
			"stop_on_error": false
		}
	}`
	response := ExecuteCommand(e, req)

	var resp ResponseEnvelope
	if err := json.Unmarshal([]byte(response), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	// Batch with stop_on_error: false returns ok=true even if individual commands fail
	// Check that results contain the error
	results, ok := resp.Result.(map[string]any)
	if !ok {
		t.Fatal("expected result to be a map")
	}
	batchResults, ok := results["results"].([]any)
	if !ok || len(batchResults) == 0 {
		t.Fatal("expected batch results")
	}

	firstResult, ok := batchResults[0].(map[string]any)
	if !ok {
		t.Fatal("expected first result to be a map")
	}
	if firstResult["ok"] == true {
		t.Error("expected first batch result to have ok=false")
	}
	// Check that error code is present
	if firstResult["error"] == nil {
		t.Fatal("expected error in first batch result")
	}
	errDetail, ok := firstResult["error"].(map[string]any)
	if !ok {
		t.Fatal("expected error to be a map")
	}
	if errDetail["code"] != ErrCodeMissingField {
		t.Errorf("expected error code %s, got %v", ErrCodeMissingField, errDetail["code"])
	}
}

func TestBridgeErrorCode_DetailsField(t *testing.T) {
	basePath := writeDeckFixture(t, "error-details-field.pptx", []elements.SlideContent{
		elements.NewSlide("Test"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer e.Close()

	// Invalid index should include details about the validation failure
	req := `{"api_version": 1, "request_id": "test", "op": "remove_slide", "payload": {"index": 999}}`
	response := ExecuteCommand(e, req)

	var resp ResponseEnvelope
	if err := json.Unmarshal([]byte(response), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error detail")
	}

	// Details should contain validation errors
	if resp.Error.Details == nil {
		t.Error("expected error details to be populated for validation failure")
	}
}

func TestBridgeErrorCode_BuildFreeformInvalidPointsType(t *testing.T) {
	basePath := writeDeckFixture(t, "error-freeform-points-type.pptx", []elements.SlideContent{
		elements.NewSlide("Test"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer e.Close()

	req := `{"api_version": 1, "request_id": "test", "op": "build_freeform", "payload": {"slide_index": 0, "points": "bad"}}`
	response := ExecuteCommand(e, req)

	var resp ResponseEnvelope
	if err := json.Unmarshal([]byte(response), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.OK {
		t.Fatal("expected ok=false for invalid points type")
	}
	if resp.Error == nil {
		t.Fatal("expected error detail")
	}
	if resp.Error.Code != ErrCodeInvalidType {
		t.Errorf("expected code %s, got %s", ErrCodeInvalidType, resp.Error.Code)
	}
}

func TestBridgeErrorCode_BuildFreeformCloseInvalidType(t *testing.T) {
	basePath := writeDeckFixture(t, "error-freeform-close-type.pptx", []elements.SlideContent{
		elements.NewSlide("Test"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer e.Close()

	req := `{"api_version": 1, "request_id": "test", "op": "build_freeform", "payload": {"slide_index": 0, "points": [[0,0],[10,10]], "close": "yes"}}`
	response := ExecuteCommand(e, req)

	var resp ResponseEnvelope
	if err := json.Unmarshal([]byte(response), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.OK {
		t.Fatal("expected ok=false for invalid close type")
	}
	if resp.Error == nil {
		t.Fatal("expected error detail")
	}
	if resp.Error.Code != ErrCodeInvalidType {
		t.Errorf("expected code %s, got %s", ErrCodeInvalidType, resp.Error.Code)
	}
}
