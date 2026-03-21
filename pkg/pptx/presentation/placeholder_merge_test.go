package presentation

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestMergePlaceholderOverrides_MergesTextAndStyleForSameTarget(t *testing.T) {
	size := 30
	bold := true
	color := "0B5FA5"
	x := styling.Inches(1)
	y := styling.Inches(2)
	cx := styling.Inches(4.5)
	cy := styling.Inches(2.2)

	merged := mergePlaceholderOverrides([]shapes.PlaceholderContent{
		{Index: 1, Type: "body", Text: "Original body placeholder text"},
		{
			Index: 1,
			Type:  "body",
			Override: &shapes.PlaceholderOverrideOptions{
				X:  &x,
				Y:  &y,
				CX: &cx,
				CY: &cy,
				TextStyle: &shapes.PlaceholderTextStyle{
					SizePt: &size,
					Bold:   &bold,
					Color:  &color,
				},
			},
		},
	})

	if len(merged) != 1 {
		t.Fatalf("expected 1 merged override, got %d", len(merged))
	}
	if merged[0].Text != "Original body placeholder text" {
		t.Fatalf("expected merged text to be preserved, got %q", merged[0].Text)
	}
	if merged[0].Override == nil || merged[0].Override.TextStyle == nil {
		t.Fatal("expected merged override text style")
	}
	if merged[0].Override.TextStyle.SizePt == nil || *merged[0].Override.TextStyle.SizePt != size {
		t.Fatalf("expected merged size %d, got %#v", size, merged[0].Override.TextStyle.SizePt)
	}
}

func TestMergePlaceholderOverrides_EmptyTypeMergesWithTypedTarget(t *testing.T) {
	merged := mergePlaceholderOverrides([]shapes.PlaceholderContent{
		{Index: 1, Type: "body", Text: "Body text"},
		{Index: 1, Override: &shapes.PlaceholderOverrideOptions{}},
	})

	if len(merged) != 1 {
		t.Fatalf("expected merge for empty type target, got %d entries", len(merged))
	}
	if merged[0].Type != "body" {
		t.Fatalf("expected merged type body, got %q", merged[0].Type)
	}
	if merged[0].Text != "Body text" {
		t.Fatalf("expected merged text preservation, got %q", merged[0].Text)
	}
}
