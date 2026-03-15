package slide

import (
	"errors"
	"strings"
	"testing"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func TestSlideNumberingHelpers(t *testing.T) {
	slides := []common.EditorSlideRef{
		{SlideID: 256, Part: "ppt/slides/slide1.xml"},
		{SlideID: 300, Part: "ppt/slides/slide4.xml"},
	}
	if got := NextSlideID(slides); got != 301 {
		t.Fatalf("NextSlideID=%d, want 301", got)
	}
	if got := NextSlidePartNumber(slides); got != 5 {
		t.Fatalf("NextSlidePartNumber=%d, want 5", got)
	}
	if got := NextSlidePartNumber(nil); got != 1 {
		t.Fatalf("NextSlidePartNumber(nil)=%d, want 1", got)
	}
}

func TestNextRelationshipIDNum(t *testing.T) {
	slides := []common.EditorSlideRef{{RelID: "rId7"}, {RelID: "rId3"}}
	nonSlide := []common.EditorRelationship{{ID: "rId10"}, {ID: "custom"}}
	if got := NextRelationshipIDNum(slides, nonSlide); got != 11 {
		t.Fatalf("NextRelationshipIDNum=%d, want 11", got)
	}
}

func TestSectionDataHelpers(t *testing.T) {
	slides := []common.EditorSlideRef{{SlideID: 1001}, {SlideID: 1002}, {SlideID: 1003}}
	ids, err := BuildSectionSlideIDs(slides, []int{0, 2})
	if err != nil {
		t.Fatalf("BuildSectionSlideIDs failed: %v", err)
	}
	if len(ids) != 2 || ids[0] != 1001 || ids[1] != 1003 {
		t.Fatalf("unexpected section ids: %+v", ids)
	}
	if _, err = BuildSectionSlideIDs(slides, []int{-1}); err == nil {
		t.Fatal("expected out-of-range index error")
	}

	current := []SectionData{{Name: "A", GUID: "g-1", SlideIDs: []int64{1001}}}
	next, err := AddSectionData(
		current,
		"B",
		[]int64{1002},
		func() (string, error) { return "g-2", nil },
	)
	if err != nil || len(next) != 2 || next[1].Name != "B" {
		t.Fatalf("AddSectionData failed: next=%+v err=%v", next, err)
	}
	if _, err = AddSectionData(current, "", nil, func() (string, error) { return "g-2", nil }); err == nil {
		t.Fatal("expected empty section name error")
	}
	if _, err = AddSectionData(current, "B", nil, func() (string, error) { return "", errors.New("boom") }); err == nil {
		t.Fatal("expected guid generation error")
	}

	removed, err := RemoveSectionData(next, "A")
	if err != nil || len(removed) != 1 || removed[0].Name != "B" {
		t.Fatalf("RemoveSectionData failed: removed=%+v err=%v", removed, err)
	}
	if _, err = RemoveSectionData(next, "Z"); err == nil {
		t.Fatal("expected missing section error")
	}

	renamed, err := RenameSectionData(next, "B", "Renamed")
	if err != nil || renamed[1].Name != "Renamed" {
		t.Fatalf("RenameSectionData failed: renamed=%+v err=%v", renamed, err)
	}
	if _, err = RenameSectionData(next, "B", ""); err == nil {
		t.Fatal("expected empty new name error")
	}
	if _, err = RenameSectionData(next, "missing", "X"); err == nil {
		t.Fatal("expected missing old name error")
	}
}

func TestTransitionHelpers(t *testing.T) {
	raw := RawSlideTransition{XMLValue: " <p:transition advTm=\"1000\"></p:transition> "}
	if err := raw.Validate(); err != nil {
		t.Fatalf("RawSlideTransition.Validate failed: %v", err)
	}
	if got := raw.XML(); got != `<p:transition advTm="1000"></p:transition>` {
		t.Fatalf("RawSlideTransition.XML unexpected: %q", got)
	}
	if err := (RawSlideTransition{XMLValue: "<p:bad></p:bad>"}).Validate(); err == nil {
		t.Fatal("expected transition validation error")
	}

	xml := `<p:sld><p:cSld/><p:transition advTm="500"/><p:timing/></p:sld>`
	if got := ExtractSlideTransitionXML(xml); got != "" {
		t.Fatalf("self-closing transition should not match explicit close-tag parser, got %q", got)
	}
	xml = `<p:sld><p:cSld/><p:transition advTm="500"></p:transition><p:timing/></p:sld>`
	if got := ExtractSlideTransitionXML(xml); got == "" {
		t.Fatal("expected ExtractSlideTransitionXML to return transition block")
	}

	base := elements.NewSlide("t")
	preserved := PreserveExistingSlideTransition(func(string) ([]byte, bool) {
		return []byte(xml), true
	}, "ppt/slides/slide1.xml", base)
	if preserved.Transition == nil {
		t.Fatal("expected transition to be preserved from existing slide XML")
	}
	if _, ok := preserved.Transition.(RawSlideTransition); !ok {
		t.Fatalf("expected RawSlideTransition type, got %T", preserved.Transition)
	}

	alreadySet := base.WithTransition(RawSlideTransition{XMLValue: `<p:transition></p:transition>`})
	unchanged := PreserveExistingSlideTransition(func(string) ([]byte, bool) {
		return []byte(`<p:transition advTm="300"></p:transition>`), true
	}, "ppt/slides/slide1.xml", alreadySet)
	if _, ok := unchanged.Transition.(RawSlideTransition); !ok {
		t.Fatalf("expected existing transition preserved, got %T", unchanged.Transition)
	}
}

func TestRelationshipsRenderingHelpers(t *testing.T) {
	nonSlide := []common.EditorRelationship{
		{ID: "rId1", Type: common.RelTypeSlideMaster, Target: "slideMasters/slideMaster1.xml"},
	}
	slides := []common.EditorSlideRef{
		{RelID: "rId2", Target: "slides/slide1.xml"},
		{RelID: "rId4", Target: "slides/slide2.xml"},
	}
	xml, err := RenderPresentationRelsXML(nonSlide, slides, true, true)
	if err != nil {
		t.Fatalf("RenderPresentationRelsXML failed: %v", err)
	}
	if !strings.Contains(xml, common.RelTypeSectionList) {
		t.Fatalf("expected sectionList relationship in XML: %s", xml)
	}
	if !strings.Contains(xml, "vbaProject.bin") {
		t.Fatalf("expected vba project relationship in XML: %s", xml)
	}
	if !strings.Contains(xml, `Id="rId2"`) || !strings.Contains(xml, `Id="rId4"`) {
		t.Fatalf("expected slide relationships in XML: %s", xml)
	}

	_, err = RenderPresentationRelsXML(
		[]common.EditorRelationship{{ID: "", Type: "x"}},
		slides,
		false,
		false,
	)
	if err == nil {
		t.Fatal("expected empty relationship id error")
	}
	_, err = RenderPresentationRelsXML([]common.EditorRelationship{
		{ID: "rId1", Type: "x"},
		{ID: "rId1", Type: "y"},
	}, slides, false, false)
	if err == nil {
		t.Fatal("expected duplicate relationship id error")
	}
}
