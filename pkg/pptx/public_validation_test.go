package pptx_test

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func TestSlideContentValidate(t *testing.T) {
	s := pptx.NewSlide("Valid").AddBullet("Item")
	if err := s.Validate(1); err != nil {
		t.Errorf("expected no error for valid slide, got %v", err)
	}

	s2 := pptx.NewSlide("Invalid").AddBullet("")
	if err := s2.Validate(1); err == nil {
		t.Error("expected error for empty bullet, got nil")
	}
}

func TestTableValidate(t *testing.T) {
	tbl := pptx.NewTable([]pptx.Length{pptx.Emu(100), pptx.Emu(100)}).AddRow([]string{"A", "B"})
	if err := tbl.Validate(1); err != nil {
		t.Errorf("expected no error for valid table, got %v", err)
	}

	tbl2 := pptx.NewTable([]pptx.Length{pptx.Emu(100)}).AddRow([]string{"A", "B"})
	if err := tbl2.Validate(1); err == nil {
		t.Error("expected error for mismatched columns, got nil")
	}
}

func TestImageValidate(t *testing.T) {
	img := pptx.NewImage("path.png", 0, 0, 100, 100)
	if err := img.Validate(1, 1); err != nil {
		t.Errorf("expected no error for valid image, got %v", err)
	}

	img2 := pptx.Image{X: 0, Y: 0, CX: 100, CY: 100}
	if err := img2.Validate(1, 1); err == nil {
		t.Error("expected error for image with no source, got nil")
	}
}

func TestShapeValidate(t *testing.T) {
	sh := pptx.NewShape(pptx.ShapeTypeRectangle, 0, 0, 100, 100)
	if err := sh.Validate(1, 1); err != nil {
		t.Errorf("expected no error for valid shape, got %v", err)
	}

	sh2 := pptx.NewShape(pptx.ShapeTypeRectangle, -1, 0, 100, 100)
	if err := sh2.Validate(1, 1); err == nil {
		t.Error("expected error for shape with negative position, got nil")
	}
}
