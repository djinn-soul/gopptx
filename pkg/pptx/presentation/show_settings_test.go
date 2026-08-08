package presentation

import (
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func threeSlides() []elements.SlideContent {
	return []elements.SlideContent{
		elements.NewSlide("One"),
		elements.NewSlide("Two"),
		elements.NewSlide("Three"),
	}
}

// The show settings covered loop, mode, timings and animation; pen colour,
// laser colour, narration, media controls and the slide range had no
// representation at all.
func TestShowSettingsCoverThePlaybackOptions(t *testing.T) {
	meta := Metadata{Metadata: common.Metadata{
		ShowSettings: common.ShowSettings{
			Loop:              true,
			DisableNarration:  true,
			PenColor:          "#FF0000",
			LaserColor:        "00FF00",
			ShowMediaControls: true,
			RangeKind:         common.SlideRangeRange,
			RangeStart:        2,
			RangeEnd:          3,
		},
	}}

	parts := buildPackageParts(t, meta, threeSlides())
	// p:showPr belongs to CT_PresentationProperties, so it lives in presProps.
	presProps := parts["ppt/presProps.xml"]

	wants := []string{
		"<p:showPr",
		`loop="1"`,
		`showNarration="0"`,
		`<p:sldRg st="2" end="3"/>`,
		`<p:penClr><a:srgbClr val="FF0000"/></p:penClr>`,
		"p14:laserClr",
		"p14:showMediaCtrls",
	}
	for _, want := range wants {
		if !strings.Contains(presProps, want) {
			t.Errorf("presProps.xml is missing %s:\n%s", want, presProps)
		}
	}
	if strings.Contains(parts["ppt/presentation.xml"], "<p:showPr") {
		t.Error("showPr must not be written into presentation.xml")
	}
}

// Custom shows — named slide subsets — had no representation in the package.
func TestCustomShowsAreDeclaredAndSelectable(t *testing.T) {
	meta := Metadata{Metadata: common.Metadata{
		CustomShows: []common.CustomShow{
			{Name: "Exec summary", SlideIndices: []int{0, 2}},
		},
		ShowSettings: common.ShowSettings{
			RangeKind:      common.SlideRangeCustom,
			CustomShowName: "Exec summary",
		},
	}}

	parts := buildPackageParts(t, meta, threeSlides())
	presentation := parts["ppt/presentation.xml"]

	if !strings.Contains(presentation, `<p:custShow name="Exec summary" id="0">`) {
		t.Fatalf("custom show not declared: %s", presentation)
	}
	if got := strings.Count(presentation, "<p:sld r:id="); got != 2 {
		t.Fatalf("custom show lists %d slides, want 2", got)
	}
	if !strings.Contains(parts["ppt/presProps.xml"], `<p:custShow id="0"/>`) {
		t.Fatalf("the show does not select the custom show: %s", parts["ppt/presProps.xml"])
	}
}

func TestCustomShowRejectsAnOutOfRangeSlide(t *testing.T) {
	meta := Metadata{Metadata: common.Metadata{
		CustomShows: []common.CustomShow{{Name: "Broken", SlideIndices: []int{7}}},
	}}
	if _, err := convertCustomShows(meta.CustomShows, len(threeSlides()), 1); err == nil {
		t.Fatal("expected an error for a slide index past the end of the deck")
	}
}

// A deck that asks for nothing unusual still writes no showPr.
func TestDefaultShowSettingsWriteNoShowPr(t *testing.T) {
	parts := buildPackageParts(t, Metadata{}, threeSlides())
	if strings.Contains(parts["ppt/presProps.xml"], "<p:showPr") {
		t.Fatalf("default settings should not emit showPr: %s", parts["ppt/presProps.xml"])
	}
}
