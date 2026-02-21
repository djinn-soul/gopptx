package editor

import (
	"encoding/json"
	"testing"
)

func TestPayloadValidator_RequireInt(t *testing.T) {
	tests := []struct {
		name     string
		payload  map[string]any
		field    string
		want     int
		wantOk   bool
		wantErrs int
	}{
		{
			name:     "missing field",
			payload:  map[string]any{},
			field:    "index",
			want:     0,
			wantOk:   false,
			wantErrs: 1,
		},
		{
			name:     "float64 value",
			payload:  map[string]any{"index": float64(42)},
			field:    "index",
			want:     42,
			wantOk:   true,
			wantErrs: 0,
		},
		{
			name:     "int value",
			payload:  map[string]any{"index": 42},
			field:    "index",
			want:     42,
			wantOk:   true,
			wantErrs: 0,
		},
		{
			name:     "string value - invalid",
			payload:  map[string]any{"index": "42"},
			field:    "index",
			want:     0,
			wantOk:   false,
			wantErrs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewPayloadValidator()
			got, ok := v.RequireInt(tt.payload, tt.field)
			if ok != tt.wantOk {
				t.Errorf("RequireInt() ok = %v, want %v", ok, tt.wantOk)
			}
			if got != tt.want {
				t.Errorf("RequireInt() got = %v, want %v", got, tt.want)
			}
			if len(v.errors) != tt.wantErrs {
				t.Errorf("RequireInt() errors = %v, want %v", len(v.errors), tt.wantErrs)
			}
		})
	}
}

func TestPayloadValidator_RequireString(t *testing.T) {
	tests := []struct {
		name     string
		payload  map[string]any
		field    string
		want     string
		wantOk   bool
		wantErrs int
	}{
		{
			name:     "missing field",
			payload:  map[string]any{},
			field:    "title",
			want:     "",
			wantOk:   false,
			wantErrs: 1,
		},
		{
			name:     "valid string",
			payload:  map[string]any{"title": "Hello"},
			field:    "title",
			want:     "Hello",
			wantOk:   true,
			wantErrs: 0,
		},
		{
			name:     "int value - invalid",
			payload:  map[string]any{"title": 42},
			field:    "title",
			want:     "",
			wantOk:   false,
			wantErrs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewPayloadValidator()
			got, ok := v.RequireString(tt.payload, tt.field)
			if ok != tt.wantOk {
				t.Errorf("RequireString() ok = %v, want %v", ok, tt.wantOk)
			}
			if got != tt.want {
				t.Errorf("RequireString() got = %v, want %v", got, tt.want)
			}
			if len(v.errors) != tt.wantErrs {
				t.Errorf("RequireString() errors = %v, want %v", len(v.errors), tt.wantErrs)
			}
		})
	}
}

func TestPayloadValidator_RequireStringSlice(t *testing.T) {
	tests := []struct {
		name     string
		payload  map[string]any
		field    string
		want     []string
		wantOk   bool
		wantErrs int
	}{
		{
			name:     "missing field",
			payload:  map[string]any{},
			field:    "items",
			want:     nil,
			wantOk:   false,
			wantErrs: 1,
		},
		{
			name:     "valid string array",
			payload:  map[string]any{"items": []any{"a", "b", "c"}},
			field:    "items",
			want:     []string{"a", "b", "c"},
			wantOk:   true,
			wantErrs: 0,
		},
		{
			name:     "mixed types - invalid",
			payload:  map[string]any{"items": []any{"a", 1, "c"}},
			field:    "items",
			want:     nil,
			wantOk:   false,
			wantErrs: 1,
		},
		{
			name:     "not an array",
			payload:  map[string]any{"items": "not-array"},
			field:    "items",
			want:     nil,
			wantOk:   false,
			wantErrs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewPayloadValidator()
			got, ok := v.RequireStringSlice(tt.payload, tt.field)
			if ok != tt.wantOk {
				t.Errorf("RequireStringSlice() ok = %v, want %v", ok, tt.wantOk)
			}
			if tt.wantOk && len(got) != len(tt.want) {
				t.Errorf("RequireStringSlice() got = %v, want %v", got, tt.want)
			}
			if len(v.errors) != tt.wantErrs {
				t.Errorf("RequireStringSlice() errors = %v, want %v", len(v.errors), tt.wantErrs)
			}
		})
	}
}

func TestPayloadValidator_IndexBounds(t *testing.T) {
	tests := []struct {
		name     string
		index    int
		min      int
		max      int
		field    string
		wantOk   bool
		wantErrs int
	}{
		{
			name:     "valid index",
			index:    2,
			min:      0,
			max:      5,
			field:    "slide_index",
			wantOk:   true,
			wantErrs: 0,
		},
		{
			name:     "negative index",
			index:    -1,
			min:      0,
			max:      5,
			field:    "slide_index",
			wantOk:   false,
			wantErrs: 1,
		},
		{
			name:     "index at max",
			index:    5,
			min:      0,
			max:      5,
			field:    "slide_index",
			wantOk:   false,
			wantErrs: 1,
		},
		{
			name:     "index above max",
			index:    10,
			min:      0,
			max:      5,
			field:    "slide_index",
			wantOk:   false,
			wantErrs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewPayloadValidator()
			ok := v.IndexBounds(tt.index, tt.min, tt.max, tt.field)
			if ok != tt.wantOk {
				t.Errorf("IndexBounds() ok = %v, want %v", ok, tt.wantOk)
			}
			if len(v.errors) != tt.wantErrs {
				t.Errorf("IndexBounds() errors = %v, want %v", len(v.errors), tt.wantErrs)
			}
		})
	}
}

func TestPayloadValidator_OptionalFields(t *testing.T) {
	v := NewPayloadValidator()
	payload := map[string]any{
		"optional_string": "value",
		"optional_int":    42,
		"optional_bool":   true,
	}

	// Test optional string
	s := v.OptionalString(payload, "optional_string")
	if s != "value" {
		t.Errorf("OptionalString() = %v, want 'value'", s)
	}

	// Test missing optional string
	s = v.OptionalString(payload, "missing_string")
	if s != "" {
		t.Errorf("OptionalString() for missing = %v, want ''", s)
	}

	// Test optional int
	i, ok := v.OptionalInt(payload, "optional_int")
	if !ok || i != 42 {
		t.Errorf("OptionalInt() = %v, %v, want 42, true", i, ok)
	}

	// Test missing optional int
	i, ok = v.OptionalInt(payload, "missing_int")
	if ok || i != 0 {
		t.Errorf("OptionalInt() for missing = %v, %v, want 0, false", i, ok)
	}

	// Test optional bool
	b, ok := v.OptionalBool(payload, "optional_bool")
	if !ok || !b {
		t.Errorf("OptionalBool() = %v, %v, want true, true", b, ok)
	}

	// Test missing optional bool
	b, ok = v.OptionalBool(payload, "missing_bool")
	if ok || b {
		t.Errorf("OptionalBool() for missing = %v, %v, want false, false", b, ok)
	}

	if v.HasErrors() {
		t.Errorf("Unexpected errors: %v", v.errors)
	}
}

func TestPayloadValidator_CombinedValidation(t *testing.T) {
	// Simulate validating a typical slide operation payload
	payload := map[string]any{
		"slide_index": 0,
		"title":       "New Title",
		"bullets":     []any{"point 1", "point 2"},
	}

	v := NewPayloadValidator()

	// Validate required fields
	slideIndex, ok := v.RequireInt(payload, "slide_index")
	if !ok {
		t.Fatalf("slide_index validation failed")
	}
	v.IndexBounds(slideIndex, 0, 10, "slide_index")

	// Validate optional fields
	title := v.OptionalString(payload, "title")
	if title != "New Title" {
		t.Errorf("title = %v, want 'New Title'", title)
	}

	bullets, ok := v.RequireStringSlice(payload, "bullets")
	if !ok {
		t.Fatalf("bullets validation failed")
	}
	if len(bullets) != 2 {
		t.Errorf("bullets length = %v, want 2", len(bullets))
	}

	if v.HasErrors() {
		t.Errorf("Unexpected errors: %v", v.errors)
	}
}

func TestParseRawPayload(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid JSON",
			input:   `{"slide_index": 0, "title": "Test"}`,
			wantErr: false,
		},
		{
			name:    "empty payload",
			input:   ``,
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			input:   `{invalid}`,
			wantErr: true,
		},
		{
			name:    "null payload",
			input:   `null`,
			wantErr: false, // null parses to nil map
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseRawPayload(json.RawMessage(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRawPayload() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.input != `null` && result == nil {
				t.Errorf("ParseRawPayload() result is nil for valid input")
			}
		})
	}
}

func TestBridgeError(t *testing.T) {
	err := NewBridgeError(ErrCodeInvalidPayload, "test error")
	if err.Code != ErrCodeInvalidPayload {
		t.Errorf("Code = %v, want %v", err.Code, ErrCodeInvalidPayload)
	}
	if err.Message != "test error" {
		t.Errorf("Message = %v, want 'test error'", err.Message)
	}
	if err.Error() != "test error" {
		t.Errorf("Error() = %v, want 'test error'", err.Error())
	}

	errWithDetails := NewBridgeErrorWithDetails(ErrCodeInvalidIndex, "index out of bounds", map[string]int{"index": 5})
	if errWithDetails.Details == nil {
		t.Error("Details should not be nil")
	}
}
