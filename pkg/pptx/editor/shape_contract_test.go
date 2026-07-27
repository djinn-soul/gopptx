package editor

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// OOXML types cx/cy as ST_PositiveCoordinate, so the edit path must reject what
// the create path already rejects, at the call that caused it.
func TestAddShapeRejectsInvalidExtents(t *testing.T) {
	for _, tc := range []struct {
		name string
		w, h float64
	}{
		{"negative width", -300, 400},
		{"negative height", 300, -400},
		{"zero size", 0, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			e := newTableEditorFixture()
			if _, err := e.AddShape(0, "rect", 100, 100, tc.w, tc.h); err == nil {
				t.Error("expected an error, got nil")
			}
		})
	}
}

func TestAddShapeAcceptsNegativePosition(t *testing.T) {
	e := newTableEditorFixture()
	// Valid OOXML: places the shape partly off the slide.
	if _, err := e.AddShape(0, "rect", -500, -500, 300, 400); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestUpdateShapeRejectsInvalidExtents(t *testing.T) {
	e := newTableEditorFixture()
	id, err := e.AddShape(0, "rect", 100, 100, 300, 400)
	if err != nil {
		t.Fatalf("AddShape: %v", err)
	}

	for _, tc := range []struct {
		name    string
		updates common.ShapeUpdate
	}{
		{"negative width", common.ShapeUpdate{W: intPtr(-500)}},
		{"negative height", common.ShapeUpdate{H: intPtr(-500)}},
		{"zero width", common.ShapeUpdate{W: intPtr(0)}},
		{"zero height", common.ShapeUpdate{H: intPtr(0)}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := e.UpdateShape(0, id, tc.updates); err == nil {
				t.Error("expected an error, got nil")
			}
		})
	}
}

// PowerPoint does not reuse <p:cNvPr id> values within a slide, and neither can
// this editor: the whole API is id-addressed, so a recycled id silently
// redirects a caller's edit to an unrelated shape.
func TestShapeIDsAreNotRecycledAfterRemoval(t *testing.T) {
	e := newTableEditorFixture()
	first, err := e.AddShape(0, "rect", 100, 100, 300, 400)
	if err != nil {
		t.Fatalf("AddShape: %v", err)
	}
	if err = e.RemoveShape(0, first); err != nil {
		t.Fatalf("RemoveShape: %v", err)
	}

	second, err := e.AddShape(0, "ellipse", 200, 200, 300, 400)
	if err != nil {
		t.Fatalf("AddShape after removal: %v", err)
	}
	if second == first {
		t.Fatalf("id %d was recycled after removal", second)
	}
	if err = e.UpdateShape(0, first, common.ShapeUpdate{X: intPtr(1)}); err == nil {
		t.Error("expected an error updating a removed shape id, got nil")
	}
}

func TestReserveShapeIDsAreNotHandedOutAgain(t *testing.T) {
	e := newTableEditorFixture()
	reserved, err := e.ReserveShapeIDs(0, 3)
	if err != nil {
		t.Fatalf("ReserveShapeIDs: %v", err)
	}

	added, err := e.AddShape(0, "rect", 0, 0, 100, 100)
	if err != nil {
		t.Fatalf("AddShape: %v", err)
	}
	for _, id := range reserved {
		if id == added {
			t.Fatalf("reserved id %d was handed out again", id)
		}
	}
}

// Flip is settable through ShapeUpdate, so it must also be readable: otherwise
// a read → modify → write round trip drops it.
func TestFlipRoundTripsThroughTheReadModel(t *testing.T) {
	e := newTableEditorFixture()
	id, err := e.AddShape(0, "rect", 100, 100, 300, 400)
	if err != nil {
		t.Fatalf("AddShape: %v", err)
	}
	if err = e.UpdateShape(0, id, common.ShapeUpdate{
		FlipH: boolPtr(true),
		FlipV: boolPtr(true),
	}); err != nil {
		t.Fatalf("UpdateShape: %v", err)
	}

	shapes, err := e.GetShapes(0)
	if err != nil {
		t.Fatalf("GetShapes: %v", err)
	}
	for _, shape := range shapes {
		if shape.ID != id {
			continue
		}
		if !shape.FlipH || !shape.FlipV {
			part, _ := e.parts.Get(e.slides[0].Part)
			t.Fatalf("flip_h=%v flip_v=%v, want both true; slide XML:\n%s",
				shape.FlipH, shape.FlipV, truncate(string(part), 600))
		}
		return
	}
	t.Fatalf("shape %d not found", id)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func TestShapeTextIsUnaffectedByExtentValidation(t *testing.T) {
	e := newTableEditorFixture()
	id, err := e.AddShape(0, "rect", 0, 0, 100, 100)
	if err != nil {
		t.Fatalf("AddShape: %v", err)
	}
	if err = e.UpdateShape(0, id, common.ShapeUpdate{W: intPtr(200)}); err != nil {
		t.Fatalf("UpdateShape with a valid width: %v", err)
	}
	part, _ := e.parts.Get(e.slides[0].Part)
	if !strings.Contains(string(part), `cx="200"`) {
		t.Errorf("width was not applied:\n%s", truncate(string(part), 600))
	}
}
