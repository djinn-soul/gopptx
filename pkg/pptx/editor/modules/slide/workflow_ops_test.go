package slide

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestWorkflowAddRemoveMove(t *testing.T) {
	parts := map[string][]byte{}
	setPart := func(path string, content []byte) { parts[path] = content }
	deletePart := func(path string) { delete(parts, path) }

	state := AddSlideState{
		Slides:       nil,
		NextSlideNum: 1,
		NextRelIDNum: 2,
		NextSlideID:  256,
		SlideWidth:   100,
		SlideHeight:  200,
	}
	render := func(_ elements.SlideContent, _ string, _ int, _ int64, _ int64) (string, string, error) {
		return "<slide/>", "<rels/>", nil
	}
	added, idx, err := AddSlideToState(state, elements.NewSlide("S1"), render, setPart)
	if err != nil || idx != 0 || len(added.Slides) != 1 {
		t.Fatalf("AddSlideToState failed: idx=%d state=%+v err=%v", idx, added, err)
	}
	if _, ok := parts["ppt/slides/slide1.xml"]; !ok {
		t.Fatal("expected slide part to be created")
	}
	if _, ok := parts["ppt/slides/_rels/slide1.xml.rels"]; !ok {
		t.Fatal("expected slide rels part to be created")
	}

	removeState := RemoveSlideState{
		Slides: added.Slides,
		NotesInventory: map[string]string{
			"ppt/slides/slide1.xml": "ppt/notesSlides/notesSlide1.xml",
		},
	}
	removed, err := RemoveSlideAt(removeState, 0, deletePart)
	if err != nil || len(removed.Slides) != 0 {
		t.Fatalf("RemoveSlideAt failed: state=%+v err=%v", removed, err)
	}
	if _, ok := parts["ppt/slides/slide1.xml"]; ok {
		t.Fatal("expected slide part to be deleted")
	}
	if _, ok := parts["ppt/notesSlides/notesSlide1.xml"]; ok {
		t.Fatal("expected notes part to be deleted")
	}
	if _, err = RemoveSlideAt(removeState, 5, deletePart); err == nil {
		t.Fatal("expected remove out-of-range error")
	}

	moved, err := MoveSlideRefs(
		[]common.EditorSlideRef{{Part: "a"}, {Part: "b"}, {Part: "c"}},
		0,
		2,
	)
	if err != nil || moved[2].Part != "a" {
		t.Fatalf("MoveSlideRefs failed: moved=%+v err=%v", moved, err)
	}
	if _, err = MoveSlideRefs([]common.EditorSlideRef{{Part: "a"}}, 2, 0); err == nil {
		t.Fatal("expected from out-of-range error")
	}
}

func TestWorkflowDuplicateMergeAndValidation(t *testing.T) {
	parts := map[string][]byte{
		"ppt/slides/slide1.xml":            []byte(`<p:sld><a:t>Old</a:t></p:sld>`),
		"ppt/slides/_rels/slide1.xml.rels": []byte(`<Relationships/>`),
	}
	getPart := func(path string) ([]byte, bool) { v, ok := parts[path]; return v, ok }
	setPart := func(path string, content []byte) { parts[path] = content }

	dupState := DuplicateState{
		Slides: []common.EditorSlideRef{
			{SlideID: 256, RelID: "rId2", Part: "ppt/slides/slide1.xml", Title: "Slide 1"},
		},
		NextSlideNum: 2,
		NextRelIDNum: 3,
		NextSlideID:  300,
	}
	appendCopy := func(b []byte) []byte { return append(b, []byte(" copy")...) }
	cloneRels := func(src []byte, _ string) ([]byte, error) {
		return append([]byte{}, src...), nil
	}
	duplicated, insertIndex, err := DuplicateSlideInState(dupState, 0, 1, getPart, setPart, appendCopy, cloneRels)
	if err != nil || insertIndex != 1 || len(duplicated.Slides) != 2 {
		t.Fatalf("DuplicateSlideInState failed: idx=%d state=%+v err=%v", insertIndex, duplicated, err)
	}
	if !strings.Contains(duplicated.Slides[1].Title, "(Copy)") {
		t.Fatalf("expected duplicated title suffix, got %q", duplicated.Slides[1].Title)
	}
	if _, ok := parts["ppt/slides/slide2.xml"]; !ok {
		t.Fatal("expected duplicated slide part")
	}
	if _, _, err = DuplicateSlideInState(dupState, 3, 0, getPart, setPart, appendCopy, cloneRels); err == nil {
		t.Fatal("expected duplicate source out-of-range error")
	}
	if _, _, err = DuplicateSlideInState(dupState, 0, 5, getPart, setPart, appendCopy, cloneRels); err == nil {
		t.Fatal("expected duplicate destination out-of-range error")
	}
	_, _, err = DuplicateSlideInState(
		dupState,
		0,
		1,
		func(path string) ([]byte, bool) {
			if path == "ppt/slides/slide1.xml" {
				return []byte("<x/>"), true
			}
			return nil, false
		},
		setPart,
		appendCopy,
		cloneRels,
	)
	if err == nil {
		t.Fatal("expected missing source rels error")
	}

	mergeState := MergeState{
		Slides:       nil,
		NextSlideNum: 5,
		NextRelIDNum: 10,
		NextSlideID:  400,
	}
	sourceSlides := []common.EditorSlideRef{{Part: "ppt/slides/slideX.xml", Title: "Imported"}}
	sourceParts := map[string][]byte{
		"ppt/slides/slideX.xml":            []byte("<slide/>"),
		"ppt/slides/_rels/slideX.xml.rels": []byte("<rels/>"),
	}
	merged, err := MergeSlidesFromSource(
		mergeState,
		sourceSlides,
		func(path string) ([]byte, bool) { v, ok := sourceParts[path]; return v, ok },
		func(_ string, srcRels []byte, _ string) ([]byte, error) {
			return append([]byte{}, srcRels...), nil
		},
		setPart,
	)
	if err != nil || len(merged.Slides) != 1 || merged.Slides[0].Part != "ppt/slides/slide5.xml" {
		t.Fatalf("MergeSlidesFromSource failed: state=%+v err=%v", merged, err)
	}
	if _, err = MergeSlidesFromSource(
		mergeState,
		sourceSlides,
		func(string) ([]byte, bool) { return nil, false },
		func(string, []byte, string) ([]byte, error) { return nil, nil },
		setPart,
	); err == nil {
		t.Fatal("expected merge missing source parts error")
	}
	if _, err = MergeSlidesFromSource(
		mergeState,
		sourceSlides,
		func(path string) ([]byte, bool) { return sourceParts[path], true },
		func(string, []byte, string) ([]byte, error) { return nil, errors.New("clone fail") },
		setPart,
	); err == nil {
		t.Fatal("expected merge clone error")
	}

	if err := ValidateMergeEditorsNil(false, false); err != nil {
		t.Fatalf("ValidateMergeEditorsNil(false,false) unexpected error: %v", err)
	}
	if err := ValidateMergeEditorsNil(true, false); err == nil {
		t.Fatal("expected nil editor validation error")
	}
}

func TestWorkflowUpdateAndSetTitle(t *testing.T) {
	originalSlide := []byte(`
<p:sld xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
  <p:sp>
    <p:nvSpPr>
      <p:nvPr><p:ph type="title"/></p:nvPr>
    </p:nvSpPr>
    <p:txBody><a:p><a:r><a:t>Old Title</a:t></a:r></a:p></p:txBody>
  </p:sp>
 </p:sld>`)
	parts := map[string][]byte{
		"ppt/slides/slide3.xml":            originalSlide,
		"ppt/slides/_rels/slide3.xml.rels": []byte(`<Relationships/>`),
	}
	getPart := func(path string) ([]byte, bool) { v, ok := parts[path]; return v, ok }
	setPart := func(path string, content []byte) { parts[path] = content }
	parseRels := func([]byte) ([]common.EditorRelationship, error) {
		return []common.EditorRelationship{{Type: common.RelTypeSlideLayout, Target: "../slideLayouts/slideLayout1.xml"}}, nil
	}
	renderExisting := func(
		slide elements.SlideContent,
		_ string,
		_ int,
		_ string,
		_, _ int64,
	) (string, string, error) {
		return fmt.Sprintf("<slide title=%q/>", slide.Title), "<rels/>", nil
	}

	updateState := UpdateSlideState{
		Slides: []common.EditorSlideRef{{Part: "ppt/slides/slide3.xml", Title: "Old"}},
	}
	updated, err := UpdateSlideInState(
		updateState,
		0,
		elements.NewSlide("New"),
		getPart,
		setPart,
		parseRels,
		renderExisting,
	)
	if err != nil || updated.Slides[0].Title != "New" {
		t.Fatalf("UpdateSlideInState failed: state=%+v err=%v", updated, err)
	}
	if _, err = UpdateSlideInState(updateState, 2, elements.NewSlide("x"), getPart, setPart, parseRels, renderExisting); err == nil {
		t.Fatal("expected update index out-of-range error")
	}
	if _, err = UpdateSlideInState(updateState, 0, elements.NewSlide("x"), getPart, setPart,
		func([]byte) ([]common.EditorRelationship, error) { return nil, errors.New("bad rels") },
		renderExisting,
	); err == nil {
		t.Fatal("expected relationship parse error")
	}
	parts["ppt/slides/slide3.xml"] = originalSlide

	slides, err := SetSlideTitleInState(
		[]common.EditorSlideRef{{Part: "ppt/slides/slide3.xml", Title: "Old"}},
		0,
		"Replaced",
		getPart,
		setPart,
	)
	if err != nil || slides[0].Title != "Replaced" {
		t.Fatalf("SetSlideTitleInState failed: slides=%+v err=%v", slides, err)
	}
	if _, err = SetSlideTitleInState(slides, 3, "x", getPart, setPart); err == nil {
		t.Fatal("expected set title out-of-range error")
	}
	if _, err = SetSlideTitleInState(
		[]common.EditorSlideRef{{Part: "ppt/slides/missing.xml"}},
		0,
		"x",
		getPart,
		setPart,
	); err == nil {
		t.Fatal("expected set title missing part error")
	}
	if _, err = SetSlideTitleInState(
		[]common.EditorSlideRef{{Part: "ppt/slides/no-title.xml"}},
		0,
		"x",
		func(path string) ([]byte, bool) {
			if path == "ppt/slides/no-title.xml" {
				return []byte(`<p:sld><p:sp/></p:sld>`), true
			}
			return nil, false
		},
		setPart,
	); err == nil {
		t.Fatal("expected set title no text-run error")
	}
}
