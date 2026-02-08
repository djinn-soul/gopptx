package pptx

import (
	"testing"
)

func TestSlideContentValidate(t *testing.T) {
	// Valid slide
	s := NewSlide("Valid").AddBullet("Item")
	if err := s.Validate(1); err != nil {
		t.Errorf("expected no error for valid slide, got %v", err)
	}

	// Invalid slide (empty bullet)
	s2 := NewSlide("Invalid").AddBullet("")
	if err := s2.Validate(1); err == nil {
		t.Error("expected error for empty bullet, got nil")
	}
}

func TestTableValidate(t *testing.T) {
	// Valid table
	tbl := NewTable([]int64{100, 100}).AddRow([]string{"A", "B"})
	if err := tbl.Validate(1); err != nil {
		t.Errorf("expected no error for valid table, got %v", err)
	}

	// Invalid table (mismatched columns)
	tbl2 := NewTable([]int64{100}).AddRow([]string{"A", "B"})
	if err := tbl2.Validate(1); err == nil {
		t.Error("expected error for mismatched columns, got nil")
	}
}

func TestImageValidate(t *testing.T) {
	// Valid image
	img := NewImage("path.png", 0, 0, 100, 100)
	if err := img.Validate(1, 1); err != nil {
		t.Errorf("expected no error for valid image, got %v", err)
	}

	// Invalid image (no source)
	img2 := Image{X: 0, Y: 0, CX: 100, CY: 100}
	if err := img2.Validate(1, 1); err == nil {
		t.Error("expected error for image with no source, got nil")
	}
}

func TestShapeValidate(t *testing.T) {
	// Valid shape
	sh := NewShape(ShapeTypeRectangle, 0, 0, 100, 100)
	if err := sh.Validate(1, 1); err != nil {
		t.Errorf("expected no error for valid shape, got %v", err)
	}

	// Invalid shape (negative position)
	sh2 := NewShape(ShapeTypeRectangle, -1, 0, 100, 100)
	if err := sh2.Validate(1, 1); err == nil {
		t.Error("expected error for shape with negative position, got nil")
	}
}
