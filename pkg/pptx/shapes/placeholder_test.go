package shapes

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestPlaceholderTargetMapping(t *testing.T) {
	// Simple tests for struct initialization and basic fields
	target := PlaceholderTarget{
		Type:  "body",
		Index: 1,
		Name:  "Test Body",
	}

	if target.Type != "body" {
		t.Errorf("expected type body, got %s", target.Type)
	}
	if target.Index != 1 {
		t.Errorf("expected index 1, got %d", target.Index)
	}
	if target.Name != "Test Body" {
		t.Errorf("expected name Test Body, got %s", target.Name)
	}
}

func TestPlaceholderTextStyleDefaults(t *testing.T) {
	style := PlaceholderTextStyle{}
	if style.Bold != nil {
		t.Error("expected Bold to be nil by default")
	}
	if style.SizePt != nil {
		t.Error("expected SizePt to be nil by default")
	}
}

func TestPlaceholderOverrideOptionsValidate(t *testing.T) {
	x := styling.Inches(1)
	y := styling.Inches(1)
	cx := styling.Inches(4)
	cy := styling.Inches(2)
	color := "FF0000"
	size := 18
	align := "ctr"

	valid := &PlaceholderOverrideOptions{
		X:  &x,
		Y:  &y,
		CX: &cx,
		CY: &cy,
		TextStyle: &PlaceholderTextStyle{
			Color:  &color,
			SizePt: &size,
			Align:  &align,
		},
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid override, got %v", err)
	}

	partial := &PlaceholderOverrideOptions{X: &x}
	if err := partial.Validate(); err == nil {
		t.Fatalf("expected geometry validation error for partial coordinates")
	}
}
