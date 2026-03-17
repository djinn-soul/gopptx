package command

import (
	"strings"
	"testing"
)

func TestParsePlaceholderTableSpec_MissingTable(t *testing.T) {
	spec, ok, err := ParsePlaceholderTableSpec(map[string]any{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ok {
		t.Fatal("expected ok=false when table is missing")
	}
	if spec != nil {
		t.Fatalf("expected nil spec, got %+v", spec)
	}
}

func TestParsePlaceholderTableSpec_InvalidTypes(t *testing.T) {
	tests := []struct {
		name    string
		payload map[string]any
		wantErr string
	}{
		{
			name:    "table not object",
			payload: map[string]any{"table": "bad"},
			wantErr: "table must be an object",
		},
		{
			name:    "rows wrong type",
			payload: map[string]any{"table": map[string]any{"rows": "bad"}},
			wantErr: "table.rows must be a non-empty 2D array",
		},
		{
			name:    "row wrong type",
			payload: map[string]any{"table": map[string]any{"rows": []any{"bad"}}},
			wantErr: "table.rows[0] must be a non-empty array",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			spec, ok, err := ParsePlaceholderTableSpec(tc.payload)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tc.wantErr, err)
			}
			if ok {
				t.Fatal("expected ok=false on invalid payload")
			}
			if spec != nil {
				t.Fatalf("expected nil spec on invalid payload, got %+v", spec)
			}
		})
	}
}

func TestParsePlaceholderTableSpec_RequiredRowsErrors(t *testing.T) {
	tests := []struct {
		name    string
		payload map[string]any
		wantErr string
	}{
		{
			name:    "missing rows",
			payload: map[string]any{"table": map[string]any{}},
			wantErr: "table.rows is required",
		},
		{
			name:    "rows empty",
			payload: map[string]any{"table": map[string]any{"rows": []any{}}},
			wantErr: "table.rows must be a non-empty 2D array",
		},
		{
			name:    "row empty",
			payload: map[string]any{"table": map[string]any{"rows": []any{[]any{}}}},
			wantErr: "table.rows[0] must be a non-empty array",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := ParsePlaceholderTableSpec(tc.payload)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestParsePlaceholderTableSpec_CellCoercionAndExtraction(t *testing.T) {
	spec, ok, err := ParsePlaceholderTableSpec(map[string]any{
		"table": map[string]any{
			"rows": []any{
				[]any{"title", nil, 42, true, 3.5},
			},
			"alt_text":   "sample alt",
			"decorative": true,
		},
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if !ok {
		t.Fatal("expected ok=true")
	}
	if spec == nil {
		t.Fatal("expected non-nil spec")
	}
	if len(spec.Rows) != 1 || len(spec.Rows[0]) != 5 {
		t.Fatalf("unexpected row shape: %+v", spec.Rows)
	}
	if spec.Rows[0][0] != "title" {
		t.Fatalf("expected string value to pass through, got %q", spec.Rows[0][0])
	}
	if spec.Rows[0][1] != "" {
		t.Fatalf("expected nil cell to become empty string, got %q", spec.Rows[0][1])
	}
	if spec.Rows[0][2] != "42" || spec.Rows[0][3] != "true" || spec.Rows[0][4] != "3.5" {
		t.Fatalf("expected non-string coercion via fmt.Sprint, got %+v", spec.Rows[0])
	}
	if spec.AltText != "sample alt" {
		t.Fatalf("expected alt_text extraction, got %q", spec.AltText)
	}
	if !spec.IsDecorative {
		t.Fatal("expected decorative extraction to set IsDecorative=true")
	}
}
