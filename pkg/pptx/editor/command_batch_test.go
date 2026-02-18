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
