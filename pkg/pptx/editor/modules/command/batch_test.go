package command

import (
	"encoding/json"
	"errors"
	"testing"
)

func mustJSON(t *testing.T, v any) json.RawMessage {
	t.Helper()
	out, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	return out
}

func testBatchOptions() BatchOptions {
	return BatchOptions{
		BatchOp:       "batch_execute",
		UnknownOpCode: "UNKNOWN_OP",
		OpFailedCode:  "OP_FAILED",
	}
}

func TestWithBatchIndex(t *testing.T) {
	nilDetails := withBatchIndex(2, nil)
	if nilDetails["index"] != 2 {
		t.Fatalf("expected index in nil-details case, got %+v", nilDetails)
	}

	merged := withBatchIndex(3, map[string]any{"shape_id": 99})
	if merged["index"] != 3 || merged["shape_id"] != 99 {
		t.Fatalf("expected merged map details, got %+v", merged)
	}

	wrapped := withBatchIndex(4, "boom")
	if wrapped["index"] != 4 || wrapped["cause"] != "boom" {
		t.Fatalf("expected non-map details wrapped as cause, got %+v", wrapped)
	}
}

func TestHandleBatchExecute_NestedBatchRejected(t *testing.T) {
	payload := mustJSON(t, BatchPayload{
		Commands: []BatchCommand{
			{Op: "batch_execute", RequestID: "c1"},
		},
	})

	got, err := HandleBatchExecute(
		payload,
		func(string, json.RawMessage) (any, error, bool) {
			return nil, nil, false
		},
		func(error) (BridgeErrorView, bool) { return BridgeErrorView{}, false },
		testBatchOptions(),
	)
	if err != nil {
		t.Fatalf("HandleBatchExecute returned error: %v", err)
	}

	results := got.Results
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	item := results[0]
	if item.OK || item.Error == nil || item.Error.Code != "INVALID_BATCH_ITEM" {
		t.Fatalf("expected nested-batch INVALID_BATCH_ITEM result, got %+v", item)
	}
	details, ok := item.Error.Details.(map[string]int)
	if !ok || details["index"] != 0 {
		t.Fatalf("expected nested-batch index details, got %#v", item.Error.Details)
	}
}

func TestHandleBatchExecute_UnknownOp(t *testing.T) {
	payload := mustJSON(t, BatchPayload{
		Commands: []BatchCommand{
			{Op: "unknown_op", RequestID: "c1"},
		},
	})

	got, err := HandleBatchExecute(
		payload,
		func(string, json.RawMessage) (any, error, bool) {
			return nil, nil, false
		},
		func(error) (BridgeErrorView, bool) { return BridgeErrorView{}, false },
		testBatchOptions(),
	)
	if err != nil {
		t.Fatalf("HandleBatchExecute returned error: %v", err)
	}

	results := got.Results
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	item := results[0]
	if item.OK || item.Error == nil || item.Error.Code != "UNKNOWN_OP" {
		t.Fatalf("expected UNKNOWN_OP failure, got %+v", item)
	}
}

func TestHandleBatchExecute_BridgeErrorMappingAndDetailsMerge(t *testing.T) {
	bridgeErr := errors.New("bridge failed")
	payload := mustJSON(t, BatchPayload{
		Commands: []BatchCommand{
			{Op: "set_text", RequestID: "c1"},
		},
	})

	got, err := HandleBatchExecute(
		payload,
		func(op string, _ json.RawMessage) (any, error, bool) {
			if op != "set_text" {
				return nil, nil, false
			}
			return nil, bridgeErr, true
		},
		func(err error) (BridgeErrorView, bool) {
			if errors.Is(err, bridgeErr) {
				return BridgeErrorView{
					Code:    "INVALID_FIELD",
					Message: "ignored by HandleBatchExecute",
					Details: map[string]any{"shape_id": 12},
				}, true
			}
			return BridgeErrorView{}, false
		},
		testBatchOptions(),
	)
	if err != nil {
		t.Fatalf("HandleBatchExecute returned error: %v", err)
	}

	results := got.Results
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	item := results[0]
	if item.OK || item.Error == nil || item.Error.Code != "INVALID_FIELD" {
		t.Fatalf("expected mapped bridge error code, got %+v", item)
	}
	details, ok := item.Error.Details.(map[string]any)
	if !ok {
		t.Fatalf("expected map details, got %#v", item.Error.Details)
	}
	if details["index"] != 0 || details["shape_id"] != 12 {
		t.Fatalf("expected merged bridge details with index, got %+v", details)
	}
}

func TestHandleBatchExecute_PlainErrorFallbackAndStopOnError(t *testing.T) {
	payload := mustJSON(t, BatchPayload{
		StopOnError: true,
		Commands: []BatchCommand{
			{Op: "first_fail", RequestID: "c1"},
			{Op: "second_ok", RequestID: "c2"},
		},
	})

	got, err := HandleBatchExecute(
		payload,
		func(op string, _ json.RawMessage) (any, error, bool) {
			switch op {
			case "first_fail":
				return nil, errors.New("plain failure"), true
			case "second_ok":
				return map[string]any{"ran": true}, nil, true
			default:
				return nil, nil, false
			}
		},
		func(error) (BridgeErrorView, bool) { return BridgeErrorView{}, false },
		testBatchOptions(),
	)
	if err != nil {
		t.Fatalf("HandleBatchExecute returned error: %v", err)
	}

	results := got.Results
	if len(results) != 1 {
		t.Fatalf("expected stop_on_error to keep 1 result, got %d", len(results))
	}
	item := results[0]
	if item.OK || item.Error == nil || item.Error.Code != "OP_FAILED" {
		t.Fatalf("expected OP_FAILED fallback, got %+v", item)
	}
	details, ok := item.Error.Details.(map[string]any)
	if !ok || details["index"] != 0 {
		t.Fatalf("expected fallback details with index, got %#v", item.Error.Details)
	}
}

func TestHandleBatchExecute_SuccessResultPassThrough(t *testing.T) {
	payload := mustJSON(t, BatchPayload{
		Commands: []BatchCommand{
			{Op: "ok_op", RequestID: "c1"},
		},
	})
	expected := map[string]any{"shape_id": 77, "updated": true}

	got, err := HandleBatchExecute(
		payload,
		func(op string, _ json.RawMessage) (any, error, bool) {
			if op != "ok_op" {
				return nil, nil, false
			}
			return expected, nil, true
		},
		func(error) (BridgeErrorView, bool) { return BridgeErrorView{}, false },
		testBatchOptions(),
	)
	if err != nil {
		t.Fatalf("HandleBatchExecute returned error: %v", err)
	}

	results := got.Results
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	item := results[0]
	if !item.OK || item.Error != nil {
		t.Fatalf("expected successful result item, got %+v", item)
	}
	result, ok := item.Result.(map[string]any)
	if !ok || result["shape_id"] != 77 || result["updated"] != true {
		t.Fatalf("expected result pass-through payload, got %#v", item.Result)
	}
}
