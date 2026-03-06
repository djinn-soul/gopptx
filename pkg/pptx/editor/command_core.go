package editor

import (
	"encoding/json"
	"errors"
	"fmt"

	editormodcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
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

// BridgeError represents a structured error returned by the bridge API.
type BridgeError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// Error codes for the bridge API.
const (
	ErrCodeInvalidJSON    = "INVALID_JSON"
	ErrCodeUnsupportedVer = "UNSUPPORTED_VERSION"
	ErrCodeUnknownOp      = "UNKNOWN_OP"
	ErrCodeInvalidPayload = "INVALID_PAYLOAD"
	ErrCodeMissingField   = "MISSING_FIELD"
	ErrCodeInvalidIndex   = "INVALID_INDEX"
	ErrCodeInvalidType    = "INVALID_TYPE"
	ErrCodeOpFailed       = "OP_FAILED"
	ErrCodeMarshalError   = "MARSHAL_ERROR"
	ErrCodeInternalError  = "INTERNAL_ERROR"
	ErrCodeInvalidHandle  = "INVALID_HANDLE"
	ErrCodeInvalidValue   = "INVALID_VALUE"
)

// NewBridgeError creates a new BridgeError with the given code and message.
func NewBridgeError(code, message string) *BridgeError {
	return &BridgeError{Code: code, Message: message}
}

// NewBridgeErrorWithDetails creates a new BridgeError with details.
func NewBridgeErrorWithDetails(code, message string, details any) *BridgeError {
	return &BridgeError{Code: code, Message: message, Details: details}
}

func (e *BridgeError) Error() string {
	return e.Message
}

// PayloadValidator provides validation helpers for command payloads.
type PayloadValidator struct {
	errors []string
	code   string
}

// NewPayloadValidator creates a new validator.
func NewPayloadValidator() *PayloadValidator {
	return &PayloadValidator{errors: nil, code: ErrCodeInvalidPayload}
}

func (v *PayloadValidator) setCode(code string) {
	if v.code == ErrCodeInvalidPayload {
		v.code = code
	}
}

func (v *PayloadValidator) missingField(field string) {
	v.setCode(ErrCodeMissingField)
	v.errors = append(v.errors, fmt.Sprintf("missing required field: %s", field))
}

func (v *PayloadValidator) invalidType(field, expected string, value any) {
	v.setCode(ErrCodeInvalidType)
	v.errors = append(v.errors, fmt.Sprintf("field %s must be %s, got %T", field, expected, value))
}

// HasErrors returns true if any validation errors occurred.
func (v *PayloadValidator) HasErrors() bool {
	return len(v.errors) > 0
}

// Error returns a combined error message.
func (v *PayloadValidator) Error() error {
	if len(v.errors) == 0 {
		return nil
	}
	return NewBridgeErrorWithDetails(v.code, "payload validation failed", v.errors)
}

// ParseRawPayload parses [json.RawMessage] into a map for validation.
func ParseRawPayload(payload json.RawMessage) (map[string]any, error) {
	if len(payload) == 0 {
		return nil, NewBridgeError(ErrCodeInvalidPayload, "empty payload")
	}
	var result map[string]any
	if err := json.Unmarshal(payload, &result); err != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, fmt.Sprintf("invalid JSON payload: %v", err))
	}
	return result, nil
}

func requireSlideIndex(
	e *PresentationEditor,
	payload map[string]any,
	v *PayloadValidator,
) (int, bool) {
	slideIndex, ok := v.RequireInt(payload, "slide_index")
	if !ok {
		return 0, false
	}
	if !v.IndexBounds(slideIndex, 0, e.SlideCount(), "slide_index") {
		return 0, false
	}
	return slideIndex, true
}

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

func handleBatchExecute(e *PresentationEditor, payload json.RawMessage) (any, error) {
	result, err := editormodcommand.HandleBatchExecute(
		payload,
		func(op string) (func(json.RawMessage) (any, error), bool) {
			handler, ok := commandHandlerFor(op)
			if !ok {
				return nil, false
			}
			return func(itemPayload json.RawMessage) (any, error) {
				return handler(e, itemPayload)
			}, true
		},
		func(err error) (editormodcommand.BridgeErrorView, bool) {
			var bridgeErr *BridgeError
			if errors.As(err, &bridgeErr) {
				return editormodcommand.BridgeErrorView{
					Code:    bridgeErr.Code,
					Message: bridgeErr.Message,
					Details: bridgeErr.Details,
				}, true
			}
			return editormodcommand.BridgeErrorView{}, false
		},
		editormodcommand.BatchOptions{
			BatchOp:       OpBatchExecute,
			UnknownOpCode: ErrCodeUnknownOp,
			OpFailedCode:  ErrCodeOpFailed,
		},
	)
	if err != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	return result, nil
}

func commandHandlerFor(op string) (commandHandler, bool) {
	for _, lookup := range []func(string) (commandHandler, bool){
		commandHandlerForSlides,
		commandHandlerForLayoutMetadata,
		commandHandlerForContent,
		commandHandlerForCommentsShapes,
		commandHandlerForNotesTables,
	} {
		if h, ok := lookup(op); ok {
			return h, true
		}
	}
	return nil, false
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
