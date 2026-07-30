package elements_test

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// Per-object rules are independent, so a slide with several defects must report
// them all: otherwise fixing a deck takes one validate pass per problem.
func TestValidateAllReportsEveryObjectDefect(t *testing.T) {
	slide := elements.SlideContent{
		Title:  "t",
		Layout: elements.SlideLayoutTitleAndContent,
		Shapes: []shapes.Shape{
			shapes.NewShape(shapes.ShapeTypeRectangle, 0, 0, 0, 100),
			shapes.NewShape(shapes.ShapeTypeRectangle, 0, 0, 100, 0),
			shapes.NewShape(shapes.ShapeTypeRectangle, 0, 0, -100, 100),
		},
	}

	errs := slide.ValidateAll(1)
	if len(errs) != 3 {
		t.Fatalf("got %d errors, want 3: %v", len(errs), errs)
	}
	for i, err := range errs {
		if !strings.Contains(err.Error(), "size must be > 0") {
			t.Errorf("error %d = %v, want a size complaint", i, err)
		}
	}
}

func TestValidateAllIsEmptyForAValidSlide(t *testing.T) {
	slide := elements.SlideContent{
		Title:  "t",
		Layout: elements.SlideLayoutTitleAndContent,
		Shapes: []shapes.Shape{
			shapes.NewShape(shapes.ShapeTypeRectangle, -50, -50, 100, 100),
		},
	}

	if errs := slide.ValidateAll(1); len(errs) != 0 {
		t.Fatalf("got %v, want no errors", errs)
	}
}

// Validate keeps its single-error shape for callers that only need a verdict.
func TestValidateStillReturnsTheFirstError(t *testing.T) {
	slide := elements.SlideContent{
		Title:  "t",
		Layout: elements.SlideLayoutTitleAndContent,
		Shapes: []shapes.Shape{shapes.NewShape(shapes.ShapeTypeRectangle, 0, 0, 0, 100)},
	}

	if err := slide.Validate(1); err == nil {
		t.Fatal("expected an error")
	}
}
