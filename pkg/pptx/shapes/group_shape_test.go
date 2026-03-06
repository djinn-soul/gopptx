package shapes

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestNewGroupShape(t *testing.T) {
	g := NewGroupShape(
		styling.Inches(1),
		styling.Inches(2),
		styling.Inches(3),
		styling.Inches(4),
	)

	if g.X != styling.Inches(1) {
		t.Errorf("expected X=1 inch, got %v", g.X)
	}
	if g.Y != styling.Inches(2) {
		t.Errorf("expected Y=2 inches, got %v", g.Y)
	}
	if g.CX != styling.Inches(3) {
		t.Errorf("expected CX=3 inches, got %v", g.CX)
	}
	if g.CY != styling.Inches(4) {
		t.Errorf("expected CY=4 inches, got %v", g.CY)
	}
	if len(g.Shapes) != 0 {
		t.Errorf("expected empty shapes, got %d", len(g.Shapes))
	}
}

func TestNewGroupShapeBounds(t *testing.T) {
	shapes := []Shape{
		NewShape(ShapeTypeRectangle, styling.Emu(0), styling.Emu(0), styling.Emu(100000), styling.Emu(100000)),
		NewShape(ShapeTypeRectangle, styling.Emu(200000), styling.Emu(200000), styling.Emu(100000), styling.Emu(100000)),
	}

	g := NewGroupShapeBounds(shapes)

	// Should calculate bounds from children: (0,0) to (300000, 300000)
	if g.X != styling.Emu(0) {
		t.Errorf("expected X=0, got %v", g.X)
	}
	if g.Y != styling.Emu(0) {
		t.Errorf("expected Y=0, got %v", g.Y)
	}
	if g.CX != styling.Emu(300000) {
		t.Errorf("expected CX=300000, got %v", g.CX)
	}
	if g.CY != styling.Emu(300000) {
		t.Errorf("expected CY=300000, got %v", g.CY)
	}
	if len(g.Shapes) != 2 {
		t.Errorf("expected 2 shapes, got %d", len(g.Shapes))
	}
}

func TestGroupShapeWithName(t *testing.T) {
	g := NewGroupShape(0, 0, 100, 100).WithName("MyGroup")

	if g.Name != "MyGroup" {
		t.Errorf("expected name 'MyGroup', got %q", g.Name)
	}
}

func TestGroupShapeWithAltText(t *testing.T) {
	g := NewGroupShape(0, 0, 100, 100).WithAltText("Accessible group")

	if g.AltText != "Accessible group" {
		t.Errorf("expected alt text 'Accessible group', got %q", g.AltText)
	}
}

func TestGroupShapeAddShape(t *testing.T) {
	shape1 := NewShape(ShapeTypeRectangle, 0, 0, 100, 100)
	shape2 := NewShape(ShapeTypeEllipse, 200, 200, 100, 100)

	g := NewGroupShape(0, 0, 400, 400).
		AddShape(shape1).
		AddShape(shape2)

	if len(g.Shapes) != 2 {
		t.Errorf("expected 2 shapes, got %d", len(g.Shapes))
	}
}

func TestGroupShapeAddShapes(t *testing.T) {
	shapes := []Shape{
		NewShape(ShapeTypeRectangle, 0, 0, 100, 100),
		NewShape(ShapeTypeEllipse, 200, 200, 100, 100),
		NewShape(ShapeTypeTriangle, 400, 400, 100, 100),
	}

	g := NewGroupShape(0, 0, 600, 600).AddShapes(shapes...)

	if len(g.Shapes) != 3 {
		t.Errorf("expected 3 shapes, got %d", len(g.Shapes))
	}
}

func TestGroupShapeValidate(t *testing.T) {
	g := NewGroupShape(0, 0, 100, 100)
	if err := g.Validate(); err != nil {
		t.Errorf("validation failed: %v", err)
	}
}

func TestShapeTypeGroup(t *testing.T) {
	if ShapeTypeGroup != "grpSp" {
		t.Errorf("expected ShapeTypeGroup='grpSp', got %q", ShapeTypeGroup)
	}
}
