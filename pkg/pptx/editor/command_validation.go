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

// RequireInt validates that a field exists and is an integer.
func (v *PayloadValidator) RequireInt(payload map[string]any, field string) (int, bool) {
	val, ok := payload[field]
	if !ok {
		v.setCode(ErrCodeMissingField)
		v.errors = append(v.errors, fmt.Sprintf("missing required field: %s", field))
		return 0, false
	}
	switch n := val.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	default:
		v.setCode(ErrCodeInvalidType)
		v.errors = append(v.errors, fmt.Sprintf("field %s must be an integer, got %T", field, val))
		return 0, false
	}
}

// RequireInt64 validates that a field exists and is an int64.
func (v *PayloadValidator) RequireInt64(payload map[string]any, field string) (int64, bool) {
	val, ok := payload[field]
	if !ok {
		v.setCode(ErrCodeMissingField)
		v.errors = append(v.errors, fmt.Sprintf("missing required field: %s", field))
		return 0, false
	}
	switch n := val.(type) {
	case float64:
		return int64(n), true
	case int:
		return int64(n), true
	case int64:
		return n, true
	default:
		v.setCode(ErrCodeInvalidType)
		v.errors = append(v.errors, fmt.Sprintf("field %s must be an integer, got %T", field, val))
		return 0, false
	}
}

// RequireString validates that a field exists and is a string.
func (v *PayloadValidator) RequireString(payload map[string]any, field string) (string, bool) {
	val, ok := payload[field]
	if !ok {
		v.setCode(ErrCodeMissingField)
		v.errors = append(v.errors, fmt.Sprintf("missing required field: %s", field))
		return "", false
	}
	s, ok := val.(string)
	if !ok {
		v.setCode(ErrCodeInvalidType)
		v.errors = append(v.errors, fmt.Sprintf("field %s must be a string, got %T", field, val))
		return "", false
	}
	return s, true
}

// RequireFloat64 validates that a field exists and is a float64.
func (v *PayloadValidator) RequireFloat64(payload map[string]any, field string) (float64, bool) {
	val, ok := payload[field]
	if !ok {
		v.setCode(ErrCodeMissingField)
		v.errors = append(v.errors, fmt.Sprintf("missing required field: %s", field))
		return 0, false
	}
	switch n := val.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		v.setCode(ErrCodeInvalidType)
		v.errors = append(v.errors, fmt.Sprintf("field %s must be a number, got %T", field, val))
		return 0, false
	}
}

// RequireStringSlice validates that a field exists and is a string slice.
func (v *PayloadValidator) RequireStringSlice(payload map[string]any, field string) ([]string, bool) {
	val, ok := payload[field]
	if !ok {
		v.setCode(ErrCodeMissingField)
		v.errors = append(v.errors, fmt.Sprintf("missing required field: %s", field))
		return nil, false
	}
	arr, ok := val.([]any)
	if !ok {
		v.setCode(ErrCodeInvalidType)
		v.errors = append(v.errors, fmt.Sprintf("field %s must be an array, got %T", field, val))
		return nil, false
	}
	result := make([]string, 0, len(arr))
	for i, item := range arr {
		s, ok := item.(string)
		if !ok {
			v.setCode(ErrCodeInvalidType)
			v.errors = append(v.errors, fmt.Sprintf("field %s[%d] must be a string, got %T", field, i, item))
			return nil, false
		}
		result = append(result, s)
	}
	return result, true
}

// RequireFloat64Slice validates that a field exists and is a float64 slice.
func (v *PayloadValidator) RequireFloat64Slice(payload map[string]any, field string) ([]float64, bool) {
	val, ok := payload[field]
	if !ok {
		v.setCode(ErrCodeMissingField)
		v.errors = append(v.errors, fmt.Sprintf("missing required field: %s", field))
		return nil, false
	}
	arr, ok := val.([]any)
	if !ok {
		v.setCode(ErrCodeInvalidType)
		v.errors = append(v.errors, fmt.Sprintf("field %s must be an array, got %T", field, val))
		return nil, false
	}
	result := make([]float64, 0, len(arr))
	for i, item := range arr {
		switch n := item.(type) {
		case float64:
			result = append(result, n)
		case int:
			result = append(result, float64(n))
		case int64:
			result = append(result, float64(n))
		default:
			v.setCode(ErrCodeInvalidType)
			v.errors = append(v.errors, fmt.Sprintf("field %s[%d] must be a number, got %T", field, i, item))
			return nil, false
		}
	}
	return result, true
}

// RequireIntSlice validates that a field exists and is an int slice.
func (v *PayloadValidator) RequireIntSlice(payload map[string]any, field string) ([]int, bool) {
	val, ok := payload[field]
	if !ok {
		v.setCode(ErrCodeMissingField)
		v.errors = append(v.errors, fmt.Sprintf("missing required field: %s", field))
		return nil, false
	}
	arr, ok := val.([]any)
	if !ok {
		v.setCode(ErrCodeInvalidType)
		v.errors = append(v.errors, fmt.Sprintf("field %s must be an array, got %T", field, val))
		return nil, false
	}
	result := make([]int, 0, len(arr))
	for i, item := range arr {
		switch n := item.(type) {
		case float64:
			result = append(result, int(n))
		case int:
			result = append(result, n)
		case int64:
			result = append(result, int(n))
		default:
			v.setCode(ErrCodeInvalidType)
			v.errors = append(v.errors, fmt.Sprintf("field %s[%d] must be an integer, got %T", field, i, item))
			return nil, false
		}
	}
	return result, true
}

// OptionalString returns a string if present, or empty string.
func (v *PayloadValidator) OptionalString(payload map[string]any, field string) string {
	val, ok := payload[field]
	if !ok {
		return ""
	}
	s, ok := val.(string)
	if !ok {
		v.setCode(ErrCodeInvalidType)
		v.errors = append(v.errors, fmt.Sprintf("field %s must be a string, got %T", field, val))
		return ""
	}
	return s
}

// OptionalInt returns an int if present, or 0 with false.
func (v *PayloadValidator) OptionalInt(payload map[string]any, field string) (int, bool) {
	val, ok := payload[field]
	if !ok {
		return 0, false
	}
	switch n := val.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	default:
		v.setCode(ErrCodeInvalidType)
		v.errors = append(v.errors, fmt.Sprintf("field %s must be an integer, got %T", field, val))
		return 0, false
	}
}

// OptionalInt64 returns an int64 if present, or 0 with false.
func (v *PayloadValidator) OptionalInt64(payload map[string]any, field string) (int64, bool) {
	val, ok := payload[field]
	if !ok {
		return 0, false
	}
	switch n := val.(type) {
	case float64:
		return int64(n), true
	case int:
		return int64(n), true
	case int64:
		return n, true
	default:
		v.setCode(ErrCodeInvalidType)
		v.errors = append(v.errors, fmt.Sprintf("field %s must be an integer, got %T", field, val))
		return 0, false
	}
}

// OptionalBool returns a bool if present, or false with false.
func (v *PayloadValidator) OptionalBool(payload map[string]any, field string) (bool, bool) {
	val, ok := payload[field]
	if !ok {
		return false, false
	}
	b, ok := val.(bool)
	if !ok {
		v.setCode(ErrCodeInvalidType)
		v.errors = append(v.errors, fmt.Sprintf("field %s must be a boolean, got %T", field, val))
		return false, false
	}
	return b, true
}

// OptionalStringSlice returns a string slice if present, or empty slice with false.
func (v *PayloadValidator) OptionalStringSlice(payload map[string]any, field string) ([]string, bool) {
	val, ok := payload[field]
	if !ok {
		return nil, false
	}
	arr, ok := val.([]any)
	if !ok {
		v.setCode(ErrCodeInvalidType)
		v.errors = append(v.errors, fmt.Sprintf("field %s must be an array, got %T", field, val))
		return nil, false
	}
	result := make([]string, 0, len(arr))
	for i, item := range arr {
		s, ok := item.(string)
		if !ok {
			v.setCode(ErrCodeInvalidType)
			v.errors = append(v.errors, fmt.Sprintf("field %s[%d] must be a string, got %T", field, i, item))
			return nil, false
		}
		result = append(result, s)
	}
	return result, true
}

// OptionalFloat64Slice returns a float64 slice if present, or empty slice with false.
func (v *PayloadValidator) OptionalFloat64Slice(payload map[string]any, field string) ([]float64, bool) {
	val, ok := payload[field]
	if !ok {
		return nil, false
	}
	arr, ok := val.([]any)
	if !ok {
		v.setCode(ErrCodeInvalidType)
		v.errors = append(v.errors, fmt.Sprintf("field %s must be an array, got %T", field, val))
		return nil, false
	}
	result := make([]float64, 0, len(arr))
	for i, item := range arr {
		switch n := item.(type) {
		case float64:
			result = append(result, n)
		case int:
			result = append(result, float64(n))
		case int64:
			result = append(result, float64(n))
		default:
			v.setCode(ErrCodeInvalidType)
			v.errors = append(v.errors, fmt.Sprintf("field %s[%d] must be a number, got %T", field, i, item))
			return nil, false
		}
	}
	return result, true
}

// IndexBounds checks if an index is within valid bounds.
func (v *PayloadValidator) IndexBounds(index, minValue, maxValue int, field string) bool {
	if index < minValue || index >= maxValue {
		v.setCode(ErrCodeInvalidIndex)
		v.errors = append(
			v.errors,
			fmt.Sprintf("%s %d out of bounds [%d, %d)", field, index, minValue, maxValue),
		)
		return false
	}
	return true
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

// parseFloat is a internal helper to handle various numeric types from unmarshaled JSON.
func parseFloat(v any) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case float32:
		return float64(val), true
	}
	return 0, false
}
