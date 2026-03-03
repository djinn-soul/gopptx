package editor

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

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
