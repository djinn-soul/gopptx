package editor

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestAddTextboxes_BulkInsert(t *testing.T) {
	basePath := writeDeckFixture(t, "shape-textbox-batch.pptx", []elements.SlideContent{
		elements.NewSlide("Textbox Batch"),
	})

	ed, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	shapeIDs, err := ed.AddTextboxes(0, []common.TextboxInsert{
		{Left: 120, Top: 120, Width: 2000, Height: 800, Text: "One"},
		{Left: 120, Top: 1040, Width: 2000, Height: 800, Text: "Two"},
	})
	if err != nil {
		t.Fatalf("bulk add textboxes: %v", err)
	}
	if len(shapeIDs) != 2 || shapeIDs[0] <= 0 || shapeIDs[1] <= shapeIDs[0] {
		t.Fatalf("unexpected shape ids: %#v", shapeIDs)
	}

	for index, expected := range []string{"One", "Two"} {
		state, stateErr := ed.GetShapeTextState(0, shapeIDs[index])
		if stateErr != nil {
			t.Fatalf("get shape text state %d: %v", index, stateErr)
		}
		if state.Text != expected {
			t.Fatalf("expected text %q, got %q", expected, state.Text)
		}
	}

	partPath := ed.slides[0].Part
	slideXML, ok := ed.parts.Get(partPath)
	if !ok {
		t.Fatalf("missing slide part %q", partPath)
	}
	for _, token := range []string{"<a:t>One</a:t>", "<a:t>Two</a:t>"} {
		if !strings.Contains(string(slideXML), token) {
			t.Fatalf("expected token %q in slide xml", token)
		}
	}
}
