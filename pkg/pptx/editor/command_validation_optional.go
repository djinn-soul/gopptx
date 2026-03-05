package editor

import editormodcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/common"

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
