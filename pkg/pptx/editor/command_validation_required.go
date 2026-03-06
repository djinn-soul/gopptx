package editor

import (
	"fmt"

	editormodcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/common"
)

// RequireInt validates that a field exists and is an integer.
func (v *PayloadValidator) RequireInt(payload map[string]any, field string) (int, bool) {
	val, ok := payload[field]
	if !ok {
		v.missingField(field)
		return 0, false
	}
	n, ok := editormodcommon.ParseInt(val)
	if !ok {
		v.invalidType(field, "an integer", val)
		return 0, false
	}
	return n, true
}

// RequireInt64 validates that a field exists and is an int64.
func (v *PayloadValidator) RequireInt64(payload map[string]any, field string) (int64, bool) {
	val, ok := payload[field]
	if !ok {
		v.missingField(field)
		return 0, false
	}
	n, ok := editormodcommon.ParseInt64(val)
	if !ok {
		v.invalidType(field, "an integer", val)
		return 0, false
	}
	return n, true
}

// RequireString validates that a field exists and is a string.
func (v *PayloadValidator) RequireString(payload map[string]any, field string) (string, bool) {
	val, ok := payload[field]
	if !ok {
		v.missingField(field)
		return "", false
	}
	s, ok := val.(string)
	if !ok {
		v.invalidType(field, "a string", val)
		return "", false
	}
	return s, true
}

// RequireFloat64 validates that a field exists and is a float64.
func (v *PayloadValidator) RequireFloat64(payload map[string]any, field string) (float64, bool) {
	val, ok := payload[field]
	if !ok {
		v.missingField(field)
		return 0, false
	}
	n, ok := editormodcommon.ParseFloat64(val)
	if !ok {
		v.invalidType(field, "a number", val)
		return 0, false
	}
	return n, true
}

// RequireStringSlice validates that a field exists and is a string slice.
func (v *PayloadValidator) RequireStringSlice(payload map[string]any, field string) ([]string, bool) {
	val, ok := payload[field]
	if !ok {
		v.missingField(field)
		return nil, false
	}
	values, ok := editormodcommon.ParseStringSlice(val)
	if !ok {
		v.invalidType(field, "an array of strings", val)
		return nil, false
	}
	return values, true
}

// RequireFloat64Slice validates that a field exists and is a float64 slice.
func (v *PayloadValidator) RequireFloat64Slice(payload map[string]any, field string) ([]float64, bool) {
	val, ok := payload[field]
	if !ok {
		v.missingField(field)
		return nil, false
	}
	values, ok := editormodcommon.ParseFloat64Slice(val)
	if !ok {
		v.invalidType(field, "an array of numbers", val)
		return nil, false
	}
	return values, true
}

// RequireIntSlice validates that a field exists and is an int slice.
func (v *PayloadValidator) RequireIntSlice(payload map[string]any, field string) ([]int, bool) {
	val, ok := payload[field]
	if !ok {
		v.missingField(field)
		return nil, false
	}
	values, ok := editormodcommon.ParseIntSlice(val)
	if !ok {
		v.invalidType(field, "an array of integers", val)
		return nil, false
	}
	return values, true
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

// OptionalString returns a string if present, or empty string.
func (v *PayloadValidator) OptionalString(payload map[string]any, field string) string {
	val, ok := payload[field]
	if !ok {
		return ""
	}
	s, ok := val.(string)
	if !ok {
		v.invalidType(field, "a string", val)
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
	n, ok := editormodcommon.ParseInt(val)
	if !ok {
		v.invalidType(field, "an integer", val)
		return 0, false
	}
	return n, true
}

// OptionalInt64 returns an int64 if present, or 0 with false.
func (v *PayloadValidator) OptionalInt64(payload map[string]any, field string) (int64, bool) {
	val, ok := payload[field]
	if !ok {
		return 0, false
	}
	n, ok := editormodcommon.ParseInt64(val)
	if !ok {
		v.invalidType(field, "an integer", val)
		return 0, false
	}
	return n, true
}

// OptionalFloat64 returns a float64 if present, or 0 with false.
func (v *PayloadValidator) OptionalFloat64(payload map[string]any, field string) (float64, bool) {
	val, ok := payload[field]
	if !ok {
		return 0, false
	}
	num, ok := editormodcommon.ParseFloat64(val)
	if !ok {
		v.invalidType(field, "a number", val)
		return 0, false
	}
	return num, true
}

// OptionalBool returns a bool if present, or false with false.
func (v *PayloadValidator) OptionalBool(payload map[string]any, field string) (bool, bool) {
	val, ok := payload[field]
	if !ok {
		return false, false
	}
	b, ok := val.(bool)
	if !ok {
		v.invalidType(field, "a boolean", val)
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
	values, ok := editormodcommon.ParseStringSlice(val)
	if !ok {
		v.invalidType(field, "an array of strings", val)
		return nil, false
	}
	return values, true
}

// OptionalFloat64Slice returns a float64 slice if present, or empty slice with false.
func (v *PayloadValidator) OptionalFloat64Slice(payload map[string]any, field string) ([]float64, bool) {
	val, ok := payload[field]
	if !ok {
		return nil, false
	}
	values, ok := editormodcommon.ParseFloat64Slice(val)
	if !ok {
		v.invalidType(field, "an array of numbers", val)
		return nil, false
	}
	return values, true
}
