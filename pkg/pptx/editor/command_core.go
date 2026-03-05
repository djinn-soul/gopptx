package editor

import (
	"encoding/json"
	"errors"
	"fmt"
)

// RequestEnvelope is the standard wrapper for all incoming commands.
type RequestEnvelope struct {
	APIVersion int             `json:"api_version"`
	Op         string          `json:"op"`
	Payload    json.RawMessage `json:"payload"`
	RequestID  string          `json:"request_id"`
}

// ResponseEnvelope is the standard wrapper for all outgoing responses.
type ResponseEnvelope struct {
	OK        bool         `json:"ok"`
	Result    any          `json:"result,omitempty"`
	Error     *ErrorDetail `json:"error,omitempty"`
	RequestID string       `json:"request_id,omitempty"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type commandHandler func(*PresentationEditor, json.RawMessage) (any, error)

// ExecuteCommand dispatches a JSON command to the appropriate editor method.
func ExecuteCommand(e *PresentationEditor, jsonInput string) string {
	var req RequestEnvelope
	if err := json.Unmarshal([]byte(jsonInput), &req); err != nil {
		return errorResponse(ErrCodeInvalidJSON, err.Error(), "")
	}

	if req.APIVersion != 1 {
		return errorResponse(
			ErrCodeUnsupportedVer,
			fmt.Sprintf("API version %d not supported", req.APIVersion),
			req.RequestID,
		)
	}

	handler, ok := commandHandlerFor(req.Op)
	if !ok {
		return errorResponse(ErrCodeUnknownOp, fmt.Sprintf("Operation %q not recognized", req.Op), req.RequestID)
	}

	result, err := handler(e, req.Payload)
	if err != nil {
		// Check if error is a BridgeError with specific code
		var bridgeErr *BridgeError
		if errors.As(err, &bridgeErr) {
			return errorResponseWithDetails(bridgeErr.Code, bridgeErr.Message, bridgeErr.Details, req.RequestID)
		}
		return errorResponse(ErrCodeOpFailed, err.Error(), req.RequestID)
	}

	resp := ResponseEnvelope{
		OK:        true,
		Result:    result,
		RequestID: req.RequestID,
	}
	out, marshalErr := json.Marshal(resp)
	if marshalErr != nil {
		return errorResponse(ErrCodeMarshalError, marshalErr.Error(), req.RequestID)
	}
	return string(out)
}

func errorResponse(code, message, reqID string) string {
	resp := ResponseEnvelope{
		OK:        false,
		RequestID: reqID,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
	out, err := json.Marshal(resp)
	if err != nil {
		return `{"ok": false, "error": {"code": "INTERNAL_ERROR", "message": "Failed to marshal error response"}}`
	}
	return string(out)
}

func errorResponseWithDetails(code, message string, details any, reqID string) string {
	resp := ResponseEnvelope{
		OK:        false,
		RequestID: reqID,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
	out, err := json.Marshal(resp)
	if err != nil {
		return `{"ok": false, "error": {"code": "INTERNAL_ERROR", "message": "Failed to marshal error response"}}`
	}
	return string(out)
}
