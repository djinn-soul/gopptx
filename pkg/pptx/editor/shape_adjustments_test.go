package editor

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// Upstream #1017: an elbow connector routes through the point its adj1 names.
// Adjustments were readable but not settable, so the only way to straighten a
// connector was to drag its yellow handle by hand, once per connector.
func TestSetShapeAdjustmentsWritesGuides(t *testing.T) {
	ed, shapeID := newAdjustmentFixture(t, "shape-adjustments.pptx")
	defer func() { _ = ed.Close() }()

	err := ed.SetShapeAdjustments(0, shapeID, []common.ShapeAdjustmentValue{
		{Name: "adj1", Value: 0.25},
	})
	if err != nil {
		t.Fatalf("set adjustments: %v", err)
	}

	shapes, err := ed.GetShapes(0)
	if err != nil {
		t.Fatalf("get shapes: %v", err)
	}
	// PowerPoint refuses a file whose round2SameRect carries only adj1, so the
	// preset's whole guide set is written: the requested value, and the preset
	// default for the guide the caller did not name.
	adjustments := adjustmentsFor(t, shapes, shapeID)
	if len(adjustments) != 2 {
		t.Fatalf("expected the preset's complete guide set, got %+v", adjustments)
	}
	if adjustments[0].Name != "adj1" || adjustments[0].Formula != "val 25000" {
		t.Fatalf("unexpected adj1: %+v", adjustments[0])
	}
	if adjustments[1].Name != "adj2" || adjustments[1].Formula != "val 0" {
		t.Fatalf("expected adj2 to fall back to the preset default, got %+v", adjustments[1])
	}
}

// A preset this package has no guide table for is refused, rather than written
// with a guessed guide set that PowerPoint would reject.
func TestSetShapeAdjustmentsRejectsUnknownPreset(t *testing.T) {
	basePath := writeDeckFixture(t, "shape-adjustments-unknown.pptx", []elements.SlideContent{
		elements.NewSlide("Adjustments"),
	})
	ed, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	shapeID, err := ed.AddShape(0, "rect", 100, 100, 600, 400)
	if err != nil {
		t.Fatalf("add shape: %v", err)
	}
	err = ed.SetShapeAdjustments(0, shapeID, []common.ShapeAdjustmentValue{{Name: "adj1", Value: 0.5}})
	if err == nil {
		t.Fatalf("expected a rect with no adjustment guides to be refused")
	}
}

// A guide the preset does not define is refused for the same reason.
func TestSetShapeAdjustmentsRejectsUnknownGuide(t *testing.T) {
	ed, shapeID := newAdjustmentFixture(t, "shape-adjustments-guide.pptx")
	defer func() { _ = ed.Close() }()

	err := ed.SetShapeAdjustments(0, shapeID, []common.ShapeAdjustmentValue{
		{Name: "adj9", Value: 0.5},
	})
	if err == nil {
		t.Fatalf("expected an unknown guide name to be refused")
	}
}

// Setting one handle must not disturb the others.
func TestSetShapeAdjustmentsMergesWithExisting(t *testing.T) {
	ed, shapeID := newAdjustmentFixture(t, "shape-adjustments-merge.pptx")
	defer func() { _ = ed.Close() }()

	initial := []common.ShapeAdjustmentValue{
		{Name: "adj1", Value: 0.25},
		{Name: "adj2", Value: 0.75},
	}
	if err := ed.SetShapeAdjustments(0, shapeID, initial); err != nil {
		t.Fatalf("set initial adjustments: %v", err)
	}

	err := ed.SetShapeAdjustments(0, shapeID, []common.ShapeAdjustmentValue{
		{Name: "adj2", Value: 0.1},
	})
	if err != nil {
		t.Fatalf("set second adjustment: %v", err)
	}

	shapes, err := ed.GetShapes(0)
	if err != nil {
		t.Fatalf("get shapes: %v", err)
	}
	adjustments := adjustmentsFor(t, shapes, shapeID)
	if len(adjustments) != 2 {
		t.Fatalf("expected both adjustments to survive, got %+v", adjustments)
	}
	if adjustments[0].Name != "adj1" || adjustments[0].Formula != "val 25000" {
		t.Fatalf("the untouched adjustment changed: %+v", adjustments[0])
	}
	if adjustments[1].Name != "adj2" || adjustments[1].Formula != "val 10000" {
		t.Fatalf("the updated adjustment is wrong: %+v", adjustments[1])
	}
}

// A raw guide formula is written verbatim, for the cases a fraction cannot
// express.
func TestSetShapeAdjustmentsAcceptsRawFormula(t *testing.T) {
	ed, shapeID := newAdjustmentFixture(t, "shape-adjustments-formula.pptx")
	defer func() { _ = ed.Close() }()

	err := ed.SetShapeAdjustments(0, shapeID, []common.ShapeAdjustmentValue{
		{Name: "adj1", Formula: "val 12345"},
	})
	if err != nil {
		t.Fatalf("set adjustments: %v", err)
	}

	content, ok := ed.parts.Get("ppt/slides/slide1.xml")
	if !ok {
		t.Fatalf("slide part missing")
	}
	if !strings.Contains(string(content), `<a:gd name="adj1" fmla="val 12345"/>`) {
		t.Fatalf("expected the raw formula in the slide XML:\n%s", content)
	}
}

func TestSetShapeAdjustmentsRejectsUnknownShape(t *testing.T) {
	ed, _ := newAdjustmentFixture(t, "shape-adjustments-missing.pptx")
	defer func() { _ = ed.Close() }()

	err := ed.SetShapeAdjustments(0, 9999, []common.ShapeAdjustmentValue{
		{Name: "adj1", Value: 0.5},
	})
	if err == nil {
		t.Fatalf("expected an error for a shape that is not on the slide")
	}
}

func newAdjustmentFixture(t *testing.T, name string) (*PresentationEditor, int) {
	t.Helper()

	basePath := writeDeckFixture(t, name, []elements.SlideContent{
		elements.NewSlide("Adjustments"),
	})
	ed, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	shapeID, err := ed.AddShape(0, "round2SameRect", 100, 100, 600, 400)
	if err != nil {
		_ = ed.Close()
		t.Fatalf("add shape: %v", err)
	}
	return ed, shapeID
}

func adjustmentsFor(t *testing.T, shapes []common.Shape, shapeID int) []common.ShapeAdjustment {
	t.Helper()
	for _, shape := range shapes {
		if shape.ID == shapeID {
			return shape.Adjustments
		}
	}
	t.Fatalf("shape %d not found", shapeID)
	return nil
}
