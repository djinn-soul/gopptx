package editor

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBridgeErrorTaxonomy(t *testing.T) {
	e := &PresentationEditor{
		parts:  NewPartStore(),
		slides: nil, // 0 slides
	}

	tests := []struct {
		name     string
		input    string
		wantCode string
	}{
		{
			name:     "invalid json",
			input:    `{invalid}`,
			wantCode: ErrCodeInvalidJSON,
		},
		{
			name:     "unknown op",
			input:    `{"api_version":1, "op":"kaboom", "payload":{}}`,
			wantCode: ErrCodeUnknownOp,
		},
		{
			name:     "unsupported version",
			input:    `{"api_version":2, "op":"kaboom", "payload":{}}`,
			wantCode: ErrCodeUnsupportedVer,
		},
		{
			name:     "missing required field",
			input:    `{"api_version":1, "op":"remove_slide", "payload":{}}`,
			wantCode: ErrCodeMissingField,
		},
		{
			name:     "invalid type",
			input:    `{"api_version":1, "op":"remove_slide", "payload":{"index":"zero"}}`,
			wantCode: ErrCodeInvalidType,
		},
		{
			name:     "invalid index",
			input:    `{"api_version":1, "op":"remove_slide", "payload":{"index":5}}`,
			wantCode: ErrCodeInvalidIndex,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := ExecuteCommand(e, tt.input)
			var out ResponseEnvelope
			if err := json.Unmarshal([]byte(resp), &out); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}
			if out.OK {
				t.Fatalf("expected error, got OK")
			}
			if out.Error == nil {
				t.Fatalf("expected ErrorDetail, got nil")
			}
			if out.Error.Code != tt.wantCode {
				t.Errorf("Error.Code = %v, want %v (Message: %s)", out.Error.Code, tt.wantCode, out.Error.Message)
			}
		})
	}
}

func TestBridgeDetailedValidationErrors(t *testing.T) {
	e := &PresentationEditor{
		parts: NewPartStore(),
	}

	// add_shape requires multiple fields
	req := `{"api_version":1, "op":"add_shape", "payload":{"slide_index":0}}`
	resp := ExecuteCommand(e, req)

	if !strings.Contains(resp, "payload validation failed") {
		t.Errorf("expected validation failure message, got %s", resp)
	}
}
