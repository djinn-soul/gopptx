package editor

import (
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// A text_frame update carries only the fields the caller set, so applying one
// must not reset the rest of the bodyPr to the render defaults.
func TestUpdateShapeTextFramePreservesUnmentionedFields(t *testing.T) {
	e := newTableEditorFixture()
	id, err := e.AddShape(0, "rect", 100, 100, 3000000, 2000000)
	if err != nil {
		t.Fatalf("AddShape: %v", err)
	}

	for _, update := range []common.ShapeUpdate{
		{TextFrame: &common.TextFrame{MarginLeft: intPtr(111111)}},
		{TextFrame: &common.TextFrame{MarginRight: intPtr(222222)}},
		{TextFrame: &common.TextFrame{MarginTop: intPtr(333333)}},
		{TextFrame: &common.TextFrame{MarginBottom: intPtr(444444)}},
	} {
		if err := e.UpdateShape(0, id, update); err != nil {
			t.Fatalf("UpdateShape: %v", err)
		}
	}

	data, ok := e.parts.Get("ppt/slides/slide1.xml")
	if !ok {
		t.Fatal("slide part missing")
	}
	xml := string(data)
	for _, want := range []string{
		`lIns="111111"`,
		`rIns="222222"`,
		`tIns="333333"`,
		`bIns="444444"`,
	} {
		if !strings.Contains(xml, want) {
			t.Errorf("expected %s in slide XML, got:\n%s", want, xml)
		}
	}
}

func TestMergeTextFrameOverlaysOnlySetFields(t *testing.T) {
	wrap := true
	anchor := "ctr"
	existing := &common.TextFrame{
		MarginLeft:    intPtr(10),
		MarginRight:   intPtr(20),
		WordWrap:      &wrap,
		VerticalAlign: &anchor,
	}
	patch := &common.TextFrame{MarginRight: intPtr(99)}

	merged := mergeTextFrame(existing, patch)

	if merged.MarginRight == nil || *merged.MarginRight != 99 {
		t.Errorf("expected patched MarginRight 99, got %v", merged.MarginRight)
	}
	if merged.MarginLeft == nil || *merged.MarginLeft != 10 {
		t.Errorf("expected preserved MarginLeft 10, got %v", merged.MarginLeft)
	}
	if merged.WordWrap == nil || !*merged.WordWrap {
		t.Error("expected preserved WordWrap")
	}
	if merged.VerticalAlign == nil || *merged.VerticalAlign != "ctr" {
		t.Errorf("expected preserved VerticalAlign, got %v", merged.VerticalAlign)
	}
	if existing.MarginRight == nil || *existing.MarginRight != 20 {
		t.Error("mergeTextFrame must not mutate the existing frame")
	}
}

func TestMergeTextFrameNilHandling(t *testing.T) {
	frame := &common.TextFrame{MarginLeft: intPtr(1)}

	if got := mergeTextFrame(frame, nil); got != frame {
		t.Error("nil patch must return the existing frame unchanged")
	}
	if got := mergeTextFrame(nil, frame); got != frame {
		t.Error("nil existing must return the patch")
	}
}
