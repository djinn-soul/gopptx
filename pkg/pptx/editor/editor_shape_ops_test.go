package editor

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestClearShapesRemovesAllShapes(t *testing.T) {
	basePath := writeDeckFixture(t, "shape-clear-all.pptx", []elements.SlideContent{
		elements.NewSlide("Shapes"),
	})

	ed, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	if _, err := ed.AddShape(0, "rect", 100, 100, 600, 400); err != nil {
		t.Fatalf("add shape 1: %v", err)
	}
	if _, err := ed.AddShape(0, "ellipse", 900, 100, 600, 400); err != nil {
		t.Fatalf("add shape 2: %v", err)
	}

	before, err := ed.GetShapes(0)
	if err != nil {
		t.Fatalf("get shapes before clear: %v", err)
	}
	if len(before) < 2 {
		t.Fatalf("expected at least 2 shapes before clear, got %d", len(before))
	}

	if err := ed.ClearShapes(0); err != nil {
		t.Fatalf("clear shapes: %v", err)
	}

	after, err := ed.GetShapes(0)
	if err != nil {
		t.Fatalf("get shapes after clear: %v", err)
	}
	if len(after) != 0 {
		t.Fatalf("expected zero shapes after clear, got %d", len(after))
	}
}
