package editor

import (
	"encoding/json"
	"fmt"
)

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
