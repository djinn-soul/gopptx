package editor

import (
	"path/filepath"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestUngroupShapesRestoresChildShapes(t *testing.T) {
	path := writeDeckFixture(t, "ungroup-base.pptx", []elements.SlideContent{
		elements.NewSlide("Ungroup"),
	})

	editor, err := OpenPresentationEditor(path)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	shapeID1, err := editor.AddShape(0, "rect", 10, 10, 100, 60)
	if err != nil {
		t.Fatalf("add shape 1: %v", err)
	}
	shapeID2, err := editor.AddShape(0, "ellipse", 160, 10, 100, 60)
	if err != nil {
		t.Fatalf("add shape 2: %v", err)
	}

	groupID, err := editor.GroupShapes(0, []int{shapeID1, shapeID2})
	if err != nil {
		t.Fatalf("group shapes: %v", err)
	}

	ungroupFirstID, err := editor.UngroupShapes(0, groupID)
	if err != nil {
		t.Fatalf("ungroup shapes: %v", err)
	}
	if ungroupFirstID != shapeID1 && ungroupFirstID != shapeID2 {
		t.Fatalf("unexpected first ungrouped shape id: %d", ungroupFirstID)
	}

	shapes, err := editor.GetShapes(0)
	if err != nil {
		t.Fatalf("list shapes: %v", err)
	}

	foundShape1 := false
	foundShape2 := false
	foundGroup := false
	for _, shape := range shapes {
		switch shape.ID {
		case shapeID1:
			foundShape1 = true
		case shapeID2:
			foundShape2 = true
		case groupID:
			foundGroup = true
		}
	}

	if !foundShape1 || !foundShape2 {
		t.Fatalf("expected both child shapes after ungroup, got %+v", shapes)
	}
	if foundGroup {
		t.Fatalf("group shape should be removed after ungroup")
	}

	outPath := filepath.Join(t.TempDir(), "ungrouped.pptx")
	if saveErr := editor.Save(outPath); saveErr != nil {
		t.Fatalf("save edited deck: %v", saveErr)
	}
}
