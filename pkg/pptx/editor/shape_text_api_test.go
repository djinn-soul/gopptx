package editor

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestShapeTextAPI_RunLifecycle(t *testing.T) {
	basePath := writeDeckFixture(t, "shape-text-api.pptx", []elements.SlideContent{
		elements.NewSlide("Text API"),
	})

	ed, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	shapeID, err := ed.AddShape(0, "rect", 120, 120, 2000, 1000)
	if err != nil {
		t.Fatalf("add shape: %v", err)
	}

	bold := true
	initialRuns := []common.TextRun{{Text: "Hello"}, {Text: "World", Bold: &bold}}
	if err := ed.SetShapeRuns(0, shapeID, initialRuns); err != nil {
		t.Fatalf("set shape runs: %v", err)
	}

	runs, err := ed.GetShapeRuns(0, shapeID)
	if err != nil {
		t.Fatalf("get shape runs: %v", err)
	}
	if len(runs) != 2 || runs[0].Text != "Hello" || runs[1].Text != "World" {
		t.Fatalf("unexpected runs after set: %#v", runs)
	}
	if runs[1].Bold == nil || !*runs[1].Bold {
		t.Fatalf("expected bold to be preserved in second run")
	}

	if err := ed.UpdateRunText(0, shapeID, 1, "Go"); err != nil {
		t.Fatalf("update run text: %v", err)
	}
	if err := ed.AppendShapeRun(0, shapeID, common.TextRun{Text: "API"}); err != nil {
		t.Fatalf("append shape run: %v", err)
	}

	state, err := ed.GetShapeTextState(0, shapeID)
	if err != nil {
		t.Fatalf("get shape text state: %v", err)
	}
	if len(state.Runs) != 3 {
		t.Fatalf("expected 3 runs, got %d", len(state.Runs))
	}
	if state.Runs[1].Text != "Go" || state.Runs[2].Text != "API" {
		t.Fatalf("unexpected updated runs: %#v", state.Runs)
	}

	partPath := ed.slides[0].Part
	slideBytes, ok := ed.parts.Get(partPath)
	if !ok {
		t.Fatalf("missing slide part %q", partPath)
	}
	slideXML := string(slideBytes)
	for _, token := range []string{"<a:t>Hello</a:t>", "<a:t>Go</a:t>", "<a:t>API</a:t>"} {
		if !strings.Contains(slideXML, token) {
			t.Fatalf("expected token %q in slide xml", token)
		}
	}
}

func TestShapeTextAPI_UpdateRunTextRejectsOutOfRange(t *testing.T) {
	basePath := writeDeckFixture(t, "shape-text-api-range.pptx", []elements.SlideContent{
		elements.NewSlide("Text API"),
	})

	ed, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	shapeID, err := ed.AddShape(0, "rect", 120, 120, 2000, 1000)
	if err != nil {
		t.Fatalf("add shape: %v", err)
	}

	if err := ed.SetShapeRuns(0, shapeID, []common.TextRun{{Text: "Only"}}); err != nil {
		t.Fatalf("set shape runs: %v", err)
	}

	err = ed.UpdateRunText(0, shapeID, 5, "Nope")
	if err == nil {
		t.Fatalf("expected out-of-range error")
	}
	if !strings.Contains(err.Error(), "run index 5 out of range") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestShapeTextAPI_UpdateSlideRunTexts(t *testing.T) {
	basePath := writeDeckFixture(t, "shape-text-api-bulk.pptx", []elements.SlideContent{
		elements.NewSlide("Text API"),
	})

	ed, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	firstShapeID, err := ed.AddShape(0, "rect", 120, 120, 2000, 1000)
	if err != nil {
		t.Fatalf("add first shape: %v", err)
	}
	secondShapeID, err := ed.AddShape(0, "rect", 120, 1320, 2000, 1000)
	if err != nil {
		t.Fatalf("add second shape: %v", err)
	}

	if err := ed.SetShapeRuns(0, firstShapeID, []common.TextRun{{Text: "One"}}); err != nil {
		t.Fatalf("set first shape runs: %v", err)
	}
	if err := ed.SetShapeRuns(0, secondShapeID, []common.TextRun{{Text: "Two"}}); err != nil {
		t.Fatalf("set second shape runs: %v", err)
	}

	err = ed.UpdateSlideRunTexts(0, []common.ShapeRunTextUpdate{
		{ShapeID: firstShapeID, RunIndex: 0, Text: "Alpha"},
		{ShapeID: secondShapeID, RunIndex: 0, Text: "Beta"},
	})
	if err != nil {
		t.Fatalf("bulk update run texts: %v", err)
	}

	firstRuns, err := ed.GetShapeRuns(0, firstShapeID)
	if err != nil {
		t.Fatalf("get first shape runs: %v", err)
	}
	secondRuns, err := ed.GetShapeRuns(0, secondShapeID)
	if err != nil {
		t.Fatalf("get second shape runs: %v", err)
	}
	if firstRuns[0].Text != "Alpha" || secondRuns[0].Text != "Beta" {
		t.Fatalf("unexpected bulk-updated runs: %#v %#v", firstRuns, secondRuns)
	}
}

func TestShapeTextAPI_SetSlideShapeRuns(t *testing.T) {
	basePath := writeDeckFixture(t, "shape-text-api-set-runs-bulk.pptx", []elements.SlideContent{
		elements.NewSlide("Text API"),
	})

	ed, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	firstShapeID, err := ed.AddShape(0, "rect", 120, 120, 2000, 1000)
	if err != nil {
		t.Fatalf("add first shape: %v", err)
	}
	secondShapeID, err := ed.AddShape(0, "rect", 120, 1320, 2000, 1000)
	if err != nil {
		t.Fatalf("add second shape: %v", err)
	}

	if err := ed.SetSlideShapeRuns(0, []common.ShapeRunsUpdate{
		{ShapeID: firstShapeID, Runs: []common.TextRun{{Text: "Alpha"}}},
		{ShapeID: secondShapeID, Runs: []common.TextRun{{Text: "Beta"}}},
	}); err != nil {
		t.Fatalf("set slide shape runs: %v", err)
	}

	firstRuns, err := ed.GetShapeRuns(0, firstShapeID)
	if err != nil {
		t.Fatalf("get first shape runs: %v", err)
	}
	secondRuns, err := ed.GetShapeRuns(0, secondShapeID)
	if err != nil {
		t.Fatalf("get second shape runs: %v", err)
	}
	if firstRuns[0].Text != "Alpha" || secondRuns[0].Text != "Beta" {
		t.Fatalf("unexpected bulk-set runs: %#v %#v", firstRuns, secondRuns)
	}
}

func TestShapeTextAPI_UpdateDeckRunTexts(t *testing.T) {
	basePath := writeDeckFixture(t, "shape-text-api-deck-bulk.pptx", []elements.SlideContent{
		elements.NewSlide("Text API"),
		elements.NewSlide("Text API 2"),
	})

	ed, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	firstShapeID, err := ed.AddShape(0, "rect", 120, 120, 2000, 1000)
	if err != nil {
		t.Fatalf("add first slide shape: %v", err)
	}
	secondShapeID, err := ed.AddShape(1, "rect", 120, 120, 2000, 1000)
	if err != nil {
		t.Fatalf("add second slide shape: %v", err)
	}
	if err := ed.SetShapeRuns(0, firstShapeID, []common.TextRun{{Text: "One"}}); err != nil {
		t.Fatalf("set first slide runs: %v", err)
	}
	if err := ed.SetShapeRuns(1, secondShapeID, []common.TextRun{{Text: "Two"}}); err != nil {
		t.Fatalf("set second slide runs: %v", err)
	}

	err = ed.UpdateDeckRunTexts([]common.SlideRunTextUpdates{
		{
			SlideIndex: 0,
			Updates: []common.ShapeRunTextUpdate{
				{ShapeID: firstShapeID, RunIndex: 0, Text: "Alpha"},
			},
		},
		{
			SlideIndex: 1,
			Updates: []common.ShapeRunTextUpdate{
				{ShapeID: secondShapeID, RunIndex: 0, Text: "Beta"},
			},
		},
	})
	if err != nil {
		t.Fatalf("deck bulk update run texts: %v", err)
	}

	firstRuns, err := ed.GetShapeRuns(0, firstShapeID)
	if err != nil {
		t.Fatalf("get first shape runs: %v", err)
	}
	secondRuns, err := ed.GetShapeRuns(1, secondShapeID)
	if err != nil {
		t.Fatalf("get second shape runs: %v", err)
	}
	if firstRuns[0].Text != "Alpha" || secondRuns[0].Text != "Beta" {
		t.Fatalf("unexpected deck bulk-updated runs: %#v %#v", firstRuns, secondRuns)
	}
}

func TestShapeTextAPI_UpdateDeckRunTexts_FailFastByInputOrder(t *testing.T) {
	basePath := writeDeckFixture(t, "shape-text-api-deck-bulk-failfast.pptx", []elements.SlideContent{
		elements.NewSlide("Text API"),
		elements.NewSlide("Text API 2"),
	})

	ed, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	validShapeID, err := ed.AddShape(1, "rect", 120, 120, 2000, 1000)
	if err != nil {
		t.Fatalf("add second slide shape: %v", err)
	}
	if err := ed.SetShapeRuns(1, validShapeID, []common.TextRun{{Text: "Two"}}); err != nil {
		t.Fatalf("set second slide runs: %v", err)
	}

	err = ed.UpdateDeckRunTexts([]common.SlideRunTextUpdates{
		{
			SlideIndex: 99, // invalid first item should fail before later slide updates
			Updates: []common.ShapeRunTextUpdate{
				{ShapeID: 1, RunIndex: 0, Text: "Nope"},
			},
		},
		{
			SlideIndex: 1,
			Updates: []common.ShapeRunTextUpdate{
				{ShapeID: validShapeID, RunIndex: 0, Text: "Beta"},
			},
		},
	})
	if err == nil {
		t.Fatalf("expected out-of-range error")
	}

	secondRuns, err := ed.GetShapeRuns(1, validShapeID)
	if err != nil {
		t.Fatalf("get second shape runs: %v", err)
	}
	if secondRuns[0].Text != "Two" {
		t.Fatalf("expected second slide runs unchanged on fail-fast, got %#v", secondRuns)
	}
}

func TestShapeTextAPI_GetShapeRunsRejectsMissingShape(t *testing.T) {
	basePath := writeDeckFixture(t, "shape-text-api-missing-shape.pptx", []elements.SlideContent{
		elements.NewSlide("Text API"),
	})

	ed, err := OpenPresentationEditor(basePath)
	if err != nil {
		t.Fatalf("open editor: %v", err)
	}
	defer func() { _ = ed.Close() }()

	_, err = ed.GetShapeRuns(0, 99999)
	if err == nil {
		t.Fatalf("expected missing-shape error")
	}
	if !strings.Contains(err.Error(), "shape id 99999 not found") {
		t.Fatalf("unexpected error: %v", err)
	}
}
