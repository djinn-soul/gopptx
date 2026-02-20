package shapes_test

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestNewPieSlice(t *testing.T) {
	s := shapes.NewPieSlice(styling.Inches(1), styling.Inches(1), styling.Inches(2), styling.Inches(2), 0, 90)
	
	if s.Type != "pie" {
		t.Errorf("expected type pie, got %s", s.Type)
	}
	
	if len(s.Adjustments) != 2 {
		t.Fatalf("expected 2 adjustments, got %d", len(s.Adjustments))
	}
	
	if s.Adjustments[0].Name != "adj1" || s.Adjustments[0].Formula != "val 0" {
		t.Errorf("expected adj1 val 0, got %s %s", s.Adjustments[0].Name, s.Adjustments[0].Formula)
	}
	
	if s.Adjustments[1].Name != "adj2" || s.Adjustments[1].Formula != "val 5400000" {
		t.Errorf("expected adj2 val 5400000, got %s %s", s.Adjustments[1].Name, s.Adjustments[1].Formula)
	}
}

func TestNewArc(t *testing.T) {
	s := shapes.NewArc(styling.Inches(1), styling.Inches(1), styling.Inches(2), styling.Inches(2), 45, 180)
	
	if s.Type != "arc" {
		t.Errorf("expected type arc, got %s", s.Type)
	}
	
	if len(s.Adjustments) != 2 {
		t.Fatalf("expected 2 adjustments, got %d", len(s.Adjustments))
	}
	
	if s.Adjustments[0].Name != "adj1" || s.Adjustments[0].Formula != "val 2700000" {
		t.Errorf("expected adj1 val 2700000, got %s %s", s.Adjustments[0].Name, s.Adjustments[0].Formula)
	}
}
