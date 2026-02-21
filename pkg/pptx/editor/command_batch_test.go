package editor

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestCommandBatchExecuteSuccess(t *testing.T) {
	basePath := writeDeckFixture(t, "batch-success.pptx", []elements.SlideContent{
		elements.NewSlide("A"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = e.Close() }()

	req := `{"api_version":1,"request_id":"b1","op":"batch_execute","payload":{"commands":[` +
		`{"op":"slide_count","payload":{}},` +
		`{"op":"set_slide_title","payload":{"slide_index":0,"title":"B"}}` +
		`]}}`
	resp := ExecuteCommand(e, req)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("batch_execute failed: %s", resp)
	}
	if e.Slides()[0].Title != "B" {
		t.Fatalf("expected title update from batch")
	}
}

func TestCommandBatchExecuteStopOnError(t *testing.T) {
	basePath := writeDeckFixture(t, "batch-stop.pptx", []elements.SlideContent{
		elements.NewSlide("A"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = e.Close() }()

	req := `{"api_version":1,"request_id":"b2","op":"batch_execute","payload":{"stop_on_error":true,"commands":[` +
		`{"op":"missing_op","payload":{}},` +
		`{"op":"set_slide_title","payload":{"slide_index":0,"title":"B"}}` +
		`]}}`
	resp := ExecuteCommand(e, req)
	var out struct {
		Result struct {
			Results []struct {
				OK bool `json:"ok"`
			} `json:"results"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(resp), &out); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if len(out.Result.Results) != 1 {
		t.Fatalf("expected one result due to stop_on_error, got %d", len(out.Result.Results))
	}
	if out.Result.Results[0].OK {
		t.Fatalf("expected first result to fail")
	}
	if e.Slides()[0].Title != "A" {
		t.Fatalf("expected second command not to run")
	}
}

func TestCommandBatchExecuteNestedRejected(t *testing.T) {
	basePath := writeDeckFixture(t, "batch-nested.pptx", []elements.SlideContent{
		elements.NewSlide("A"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = e.Close() }()

	req := `{"api_version":1,"request_id":"b3","op":"batch_execute","payload":{"commands":[` +
		`{"op":"batch_execute","payload":{"commands":[]}}` +
		`]}}`
	resp := ExecuteCommand(e, req)
	if !strings.Contains(resp, "nested batch_execute is not supported") {
		t.Fatalf("expected nested batch rejection, got: %s", resp)
	}
}

func TestCommandBatchExecuteUnknownOpContinuesWhenStopOnErrorFalse(t *testing.T) {
	basePath := writeDeckFixture(t, "batch-continue.pptx", []elements.SlideContent{
		elements.NewSlide("A"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = e.Close() }()

	req := `{"api_version":1,"request_id":"b4","op":"batch_execute","payload":{"commands":[` +
		`{"op":"missing_op","payload":{}},` +
		`{"op":"set_slide_title","payload":{"slide_index":0,"title":"B"}}` +
		`]}}`
	resp := ExecuteCommand(e, req)
	if !strings.Contains(resp, `"ok":true`) {
		t.Fatalf("batch_execute failed: %s", resp)
	}
	if e.Slides()[0].Title != "B" {
		t.Fatalf("expected second command to run when stop_on_error is false")
	}
}

func TestCommandBatchExecuteMixedResultsIncludePerItemDetails(t *testing.T) {
	basePath := writeDeckFixture(t, "batch-mixed-details.pptx", []elements.SlideContent{
		elements.NewSlide("A"),
	})
	e, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = e.Close() }()

	req := `{"api_version":1,"request_id":"b5","op":"batch_execute","payload":{"commands":[` +
		`{"op":"slide_count","payload":{}},` +
		`{"op":"missing_op","payload":{}},` +
		`{"op":"set_slide_title","payload":{"slide_index":0,"title":"C"}}` +
		`]}}`
	resp := ExecuteCommand(e, req)

	var out struct {
		OK     bool `json:"ok"`
		Result struct {
			Results []struct {
				OK    bool `json:"ok"`
				Op    string
				Error struct {
					Code    string `json:"code"`
					Message string `json:"message"`
					Details struct {
						Index int `json:"index"`
					} `json:"details"`
				} `json:"error"`
			} `json:"results"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(resp), &out); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if !out.OK {
		t.Fatalf("expected top-level OK for mixed batch, got %s", resp)
	}
	if len(out.Result.Results) != 3 {
		t.Fatalf("expected 3 batch results, got %d", len(out.Result.Results))
	}
	if !out.Result.Results[0].OK || out.Result.Results[0].Op != "slide_count" {
		t.Fatalf("unexpected first result ordering/content")
	}
	if out.Result.Results[1].OK || out.Result.Results[1].Error.Code != ErrCodeUnknownOp {
		t.Fatalf("expected second result unknown-op failure, got %+v", out.Result.Results[1])
	}
	if out.Result.Results[1].Error.Details.Index != 1 {
		t.Fatalf("expected second result details.index=1, got %d", out.Result.Results[1].Error.Details.Index)
	}
	if !out.Result.Results[2].OK || out.Result.Results[2].Op != "set_slide_title" {
		t.Fatalf("unexpected third result ordering/content")
	}
	if e.Slides()[0].Title != "C" {
		t.Fatalf("expected third command to execute after failure")
	}
}
