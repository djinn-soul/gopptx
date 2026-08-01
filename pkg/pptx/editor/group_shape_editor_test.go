package editor

import (
	"fmt"
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
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

func TestAddFreeformShapeWritesCustomGeometry(t *testing.T) {
	basePath := writeDeckFixture(t, "freeform-shape-test.pptx", []elements.SlideContent{
		elements.NewSlide("Freeform Test").AddBullet("body"),
	})
	editor, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = editor.Close() }()

	shapeID, err := editor.AddFreeformShape(0, []freeformPoint{
		{X: 100, Y: 100},
		{X: 600, Y: 100},
		{X: 600, Y: 400},
	}, true)
	if err != nil {
		t.Fatalf("add freeform shape: %v", err)
	}
	if shapeID <= 0 {
		t.Fatalf("expected shapeID > 0, got %d", shapeID)
	}

	partPath := editor.slides[0].Part
	slideXML, ok := editor.parts.Get(partPath)
	if !ok {
		t.Fatalf("missing slide part %q", partPath)
	}
	xmlText := string(slideXML)
	if !strings.Contains(xmlText, "<a:custGeom>") {
		t.Fatalf("expected custom geometry in freeform shape xml")
	}
	if !strings.Contains(xmlText, "<a:pathLst>") {
		t.Fatalf("expected path list in freeform shape xml")
	}
	if !strings.Contains(xmlText, "<a:close/>") {
		t.Fatalf("expected closed path marker in freeform shape xml")
	}
}

// The slide scan consumes a group as one node, so updating a shape inside one
// used to report the id as missing (Codex review).
func TestUpdateShapeInsideGroup(t *testing.T) {
	basePath := writeDeckFixture(t, "group-nested-update-test.pptx", []elements.SlideContent{
		elements.NewSlide("Group Update").AddBullet("body"),
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
	if _, err = editor.AddGroupShape(0, []int{shapeA, shapeB}); err != nil {
		t.Fatalf("add group shape: %v", err)
	}

	text := "nested text"
	if err = editor.UpdateShape(0, shapeB, common.ShapeUpdate{Text: &text}); err != nil {
		t.Fatalf("update nested shape: %v", err)
	}

	slideXML, ok := editor.parts.Get(editor.slides[0].Part)
	if !ok {
		t.Fatal("missing slide part")
	}
	xmlText := string(slideXML)
	if !strings.Contains(xmlText, text) {
		t.Fatalf("nested shape text was not written: %s", xmlText)
	}
	// The group wrapper and its other child must survive the rewrite.
	if !strings.Contains(xmlText, "<p:grpSp") {
		t.Fatal("group wrapper was lost")
	}
	if !strings.Contains(xmlText, fmt.Sprintf(`<p:cNvPr id="%d"`, shapeA)) {
		t.Fatal("sibling shape was lost")
	}
	if strings.Count(xmlText, fmt.Sprintf(`<p:cNvPr id="%d"`, shapeB)) != 1 {
		t.Fatalf("nested shape was duplicated: %s", xmlText)
	}

	missing := "no such shape"
	if err = editor.UpdateShape(0, 999999, common.ShapeUpdate{Text: &missing}); err == nil {
		t.Fatal("expected an error for an id that is in no group")
	}
}
