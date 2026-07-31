package editor

import (
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestGetShapesIncludesRecursiveGroupChildren(t *testing.T) {
	basePath := writeDeckFixture(t, "group-shape-tree-test.pptx", []elements.SlideContent{
		elements.NewSlide("Group tree"),
	})
	presentation, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = presentation.Close() }()

	first := addTreeTestShape(t, presentation, "rect", 100, 100)
	second := addTreeTestShape(t, presentation, "ellipse", 600, 100)
	inner, err := presentation.AddGroupShape(0, []int{first, second})
	if err != nil {
		t.Fatalf("add inner group: %v", err)
	}
	third := addTreeTestShape(t, presentation, "triangle", 1100, 100)
	outer, err := presentation.AddGroupShape(0, []int{inner, third})
	if err != nil {
		t.Fatalf("add outer group: %v", err)
	}

	shapes, err := presentation.GetShapes(0)
	if err != nil {
		t.Fatalf("get shapes: %v", err)
	}
	group := findCommonShapeByID(shapes, outer)
	if group == nil || len(group.Shapes) != 2 {
		t.Fatalf("expected outer group with two children, got %#v", group)
	}
	if group.Shapes[0].ID != inner || len(group.Shapes[0].Shapes) != 2 {
		t.Fatalf("expected recursive inner group, got %#v", group.Shapes[0])
	}
	got := []int{
		group.Shapes[0].Shapes[0].ID,
		group.Shapes[0].Shapes[1].ID,
		group.Shapes[1].ID,
	}
	want := []int{first, second, third}
	for idx := range want {
		if got[idx] != want[idx] {
			t.Fatalf("leaf order mismatch: got %v want %v", got, want)
		}
	}
}

func addTreeTestShape(
	t *testing.T,
	presentation *PresentationEditor,
	shapeType string,
	x float64,
	y float64,
) int {
	t.Helper()
	id, err := presentation.AddShape(0, shapeType, x, y, 400, 300)
	if err != nil {
		t.Fatalf("add %s: %v", shapeType, err)
	}
	return id
}

func findCommonShapeByID(shapes []common.Shape, id int) *common.Shape {
	for idx := range shapes {
		if shapes[idx].ID == id {
			return &shapes[idx]
		}
	}
	return nil
}
