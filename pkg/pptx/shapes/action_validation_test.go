package shapes_test

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestShapeValidate_ClickAction(t *testing.T) {
	s := shapes.NewShape("rect", 0, 0, styling.Inches(1), styling.Inches(1)).
		WithClickAction(action.NewHyperlink(action.HyperlinkURL("")))

	err := s.Validate(0, 0)
	if err == nil {
		t.Fatal("expected error for empty click action URL")
	}
	if !strings.Contains(err.Error(), "click action") {
		t.Fatalf("expected 'click action' in error, got: %v", err)
	}
}

func TestShapeValidate_HoverAction(t *testing.T) {
	s := shapes.NewShape("rect", 0, 0, styling.Inches(1), styling.Inches(1)).
		WithHoverAction(action.NewHyperlink(action.HyperlinkURL("")))

	err := s.Validate(0, 0)
	if err == nil {
		t.Fatal("expected error for empty hover action URL")
	}
	if !strings.Contains(err.Error(), "hover action") {
		t.Fatalf("expected 'hover action' in error, got: %v", err)
	}
}

func TestShapeValidate_ValidClickAndHover(t *testing.T) {
	s := shapes.NewShape("rect", 0, 0, styling.Inches(1), styling.Inches(1)).
		WithClickAction(action.NewHyperlink(action.HyperlinkURL("https://example.com"))).
		WithHoverAction(action.NewHyperlink(action.HyperlinkNextSlide()))

	err := s.Validate(0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestShapeValidate_LegacyHyperlinkFallback(t *testing.T) {
	// When ClickAction is nil but Hyperlink is set, Validate should still check Hyperlink
	s := shapes.NewShape("rect", 0, 0, styling.Inches(1), styling.Inches(1))
	link := action.NewHyperlink(action.HyperlinkURL(""))
	s.Hyperlink = &link

	err := s.Validate(0, 0)
	if err == nil {
		t.Fatal("expected error for empty legacy hyperlink URL")
	}
	if !strings.Contains(err.Error(), "hyperlink") {
		t.Fatalf("expected 'hyperlink' in error, got: %v", err)
	}
}
