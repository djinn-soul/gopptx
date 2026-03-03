package editor

import (
	"fmt"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestAddGroupShapeWithMembers(t *testing.T) {
	basePath := writeDeckFixture(t, "group-shape-members-test.pptx", []elements.SlideContent{
		elements.NewSlide("Group Test").AddBullet("body"),
	})
	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	shapeA, err := editor.AddShape(0, "rect", 100, 100, 400, 300)
	if err != nil {
		t.Fatalf("add shape A: %v", err)
	}
	shapeB, err := editor.AddShape(0, "ellipse", 650, 100, 500, 300)
	if err != nil {
		t.Fatalf("add shape B: %v", err)
	}

	groupID, err := editor.AddGroupShape(0, []int{shapeA, shapeB})
	if err != nil {
		t.Fatalf("add group shape: %v", err)
	}
	if groupID <= 0 {
		t.Fatalf("expected group ID > 0, got %d", groupID)
	}

	partPath := editor.slides[0].Part
	slideXML, ok := editor.parts.Get(partPath)
	if !ok {
		t.Fatalf("missing slide part %q", partPath)
	}

	xmlText := string(slideXML)
	if !strings.Contains(xmlText, "<p:grpSp") {
		t.Fatalf("expected group shape node in slide xml")
	}
	if !strings.Contains(xmlText, fmt.Sprintf(`name="Group %d"`, groupID)) {
		t.Fatalf("expected group name for id %d", groupID)
	}
	if !strings.Contains(xmlText, fmt.Sprintf(`<p:cNvPr id="%d" name="rect %d"`, shapeA, shapeA)) {
		t.Fatalf("expected shape A nested in group")
	}
	if !strings.Contains(xmlText, fmt.Sprintf(`<p:cNvPr id="%d" name="ellipse %d"`, shapeB, shapeB)) {
		t.Fatalf("expected shape B nested in group")
	}
}
