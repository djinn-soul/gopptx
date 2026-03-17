package common

import "testing"

func TestNumericParsers(t *testing.T) {
	t.Run("ParseInt", func(t *testing.T) {
		got, ok := ParseInt(float64(12))
		if !ok || got != 12 {
			t.Fatalf("ParseInt(float64) = (%d,%v), want (12,true)", got, ok)
		}
		got, ok = ParseInt(int64(9))
		if !ok || got != 9 {
			t.Fatalf("ParseInt(int64) = (%d,%v), want (9,true)", got, ok)
		}
		if _, ok = ParseInt("x"); ok {
			t.Fatal("ParseInt(string) should fail")
		}
	})

	t.Run("ParseInt64", func(t *testing.T) {
		got, ok := ParseInt64(float64(42))
		if !ok || got != 42 {
			t.Fatalf("ParseInt64(float64) = (%d,%v), want (42,true)", got, ok)
		}
		if _, ok = ParseInt64("x"); ok {
			t.Fatal("ParseInt64(string) should fail")
		}
	})

	t.Run("ParseFloat64", func(t *testing.T) {
		got, ok := ParseFloat64(int32(7))
		if !ok || got != 7 {
			t.Fatalf("ParseFloat64(int32) = (%v,%v), want (7,true)", got, ok)
		}
		if _, ok = ParseFloat64("x"); ok {
			t.Fatal("ParseFloat64(string) should fail")
		}
	})
}

func TestSliceParsers(t *testing.T) {
	if vals, ok := ParseStringSlice([]any{"a", "b"}); !ok || len(vals) != 2 || vals[1] != "b" {
		t.Fatalf("ParseStringSlice valid failed: %v %v", vals, ok)
	}
	if _, ok := ParseStringSlice([]any{"a", 2}); ok {
		t.Fatal("ParseStringSlice should fail on non-string element")
	}

	if vals, ok := ParseFloat64Slice([]any{1, float32(2), int64(3)}); !ok || len(vals) != 3 || vals[2] != 3 {
		t.Fatalf("ParseFloat64Slice valid failed: %v %v", vals, ok)
	}
	if _, ok := ParseFloat64Slice([]any{1, "nope"}); ok {
		t.Fatal("ParseFloat64Slice should fail on invalid element")
	}

	if vals, ok := ParseIntSlice([]any{1, int64(2), float64(3)}); !ok || len(vals) != 3 || vals[0] != 1 {
		t.Fatalf("ParseIntSlice valid failed: %v %v", vals, ok)
	}
	if _, ok := ParseIntSlice([]any{1, "nope"}); ok {
		t.Fatal("ParseIntSlice should fail on invalid element")
	}
}

func TestParsePlaceholderTextStyle(t *testing.T) {
	opts := ParsePlaceholderTextStyle(map[string]any{
		"text_style": map[string]any{
			"size_pt": 11.0,
			"bold":    true,
			"italic":  false,
			"color":   "FF0000",
			"font":    "Calibri",
		},
	})
	if opts == nil || opts.TextStyle == nil {
		t.Fatal("expected text style options")
	}
	if opts.TextStyle.SizePt == nil || *opts.TextStyle.SizePt != 11 {
		t.Fatalf("unexpected SizePt: %+v", opts.TextStyle.SizePt)
	}
	if opts.TextStyle.Bold == nil || !*opts.TextStyle.Bold {
		t.Fatal("expected Bold=true")
	}
	if opts.TextStyle.Italic == nil || *opts.TextStyle.Italic {
		t.Fatal("expected Italic=false")
	}
	if opts.TextStyle.Color == nil || *opts.TextStyle.Color != "FF0000" {
		t.Fatal("expected Color=FF0000")
	}
	if opts.TextStyle.Font == nil || *opts.TextStyle.Font != "Calibri" {
		t.Fatal("expected Font=Calibri")
	}

	if got := ParsePlaceholderTextStyle(map[string]any{}); got != nil {
		t.Fatal("expected nil style options for missing text_style")
	}
}

func TestParsePlaceholderImageBounds(t *testing.T) {
	x, y, cx, cy, err := ParsePlaceholderImageBounds(map[string]any{
		"bounds": []any{1.0, 2.0, 3.0, 4.0},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if x != 1 || y != 2 || cx != 3 || cy != 4 {
		t.Fatalf("unexpected bounds: %v %v %v %v", x, y, cx, cy)
	}

	_, _, _, _, err = ParsePlaceholderImageBounds(map[string]any{
		"bounds": []any{1.0, 2.0, 3.0},
	})
	if err == nil {
		t.Fatal("expected bounds length error")
	}

	_, _, _, _, err = ParsePlaceholderImageBounds(map[string]any{
		"bounds": []any{1.0, "bad", 3.0, 4.0},
	})
	if err == nil {
		t.Fatal("expected bounds type error")
	}
}

func TestPlaceholderHelpers(t *testing.T) {
	ref := BuildPlaceholderImageRef("rId1", "media/image1.png", 1, 2, 3, 4)
	if ref.RelID != "rId1" || ref.Name != "media/image1.png" {
		t.Fatalf("unexpected image ref identity: %+v", ref)
	}
	if ref.X != 12700 || ref.Y != 25400 || ref.CX != 38100 || ref.CY != 50800 {
		t.Fatalf("unexpected point conversion: %+v", ref)
	}

	shapes := []PlaceholderShapeRef{
		{Index: 2, Type: "title"},
		{Index: 2, Type: "ctrTitle"},
		{Index: 3, Type: "body"},
	}
	idx, matches := FindPlaceholderShapeIndex(shapes, 2, "")
	if idx != 1 || matches != 2 {
		t.Fatalf("unexpected placeholder match-all result idx=%d matches=%d", idx, matches)
	}
	idx, matches = FindPlaceholderShapeIndex(shapes, 2, "title")
	if idx != 0 || matches != 1 {
		t.Fatalf("expected title filter to match only title, got idx=%d matches=%d", idx, matches)
	}
	idx, matches = FindPlaceholderShapeIndex(shapes, 2, "body")
	if idx != -1 || matches != 0 {
		t.Fatalf("unexpected mismatch result idx=%d matches=%d", idx, matches)
	}

	if got := ResolvePlaceholderType("title", "body"); got != "title" {
		t.Fatalf("ResolvePlaceholderType explicit failed: %q", got)
	}
	if got := ResolvePlaceholderType("", "body"); got != "body" {
		t.Fatalf("ResolvePlaceholderType fallback failed: %q", got)
	}

	spec := BuildPlaceholderOverrideSpec(
		3,
		"body",
		"hello",
		ref,
		ParsePlaceholderTextStyle(map[string]any{"text_style": map[string]any{"bold": true}}),
	)
	if spec.Index != 3 || spec.Type != "body" || spec.Text != "hello" || spec.Image == nil {
		t.Fatalf("unexpected placeholder override spec: %+v", spec)
	}
	if spec.TextStyle == nil || spec.TextStyle.Bold == nil || !*spec.TextStyle.Bold {
		t.Fatalf("expected text style bold in override spec: %+v", spec.TextStyle)
	}
}
