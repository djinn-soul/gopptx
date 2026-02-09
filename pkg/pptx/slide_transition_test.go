package pptx

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type transitionParityFixtureEntry struct {
	Transition string `json:"transition"`
	XML        string `json:"xml"`
}

type customTransition struct {
	validateErr error
	xml         string
}

func (c customTransition) Validate() error { return c.validateErr }

func (c customTransition) XML() string { return c.xml }

func TestCreateWithSlidesRendersRepresentativeTransitions(t *testing.T) {
	cases := []struct {
		name       string
		transition SlideTransition
		expectXML  string
	}{
		{name: "fade", transition: TransitionFade, expectXML: `<p:transition><p:fade/></p:transition>`},
		{name: "push", transition: TransitionPush, expectXML: `<p:transition><p:push dir="r"/></p:transition>`},
		{name: "split", transition: TransitionSplit, expectXML: `<p:transition><p:split dir="out" orient="horz"/></p:transition>`},
		{name: "zoom", transition: TransitionZoom, expectXML: `<p:transition><p:zoom dir="in"/></p:transition>`},
		{name: "none", transition: TransitionNone, expectXML: ``},
		{name: "cut", transition: TransitionCut, expectXML: ``},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			slideXML := transitionSlideXML(t, NewSlide("Transition").WithTransition(tc.transition))
			if tc.expectXML == "" {
				if strings.Contains(slideXML, "<p:transition>") {
					t.Fatalf("did not expect transition XML for %s", tc.name)
				}
				return
			}
			if !strings.Contains(slideXML, tc.expectXML) {
				t.Fatalf("expected %q in slide XML", tc.expectXML)
			}
		})
	}
}

func TestCreateWithSlidesRejectsUnknownTransition(t *testing.T) {
	_, err := CreateWithSlides("Demo", []SlideContent{
		NewSlide("Bad").WithTransition(TransitionType("spin")),
	})
	if err == nil {
		t.Fatalf("expected unknown transition validation error")
	}
	if !strings.Contains(err.Error(), `unsupported transition type`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSlideTransitionSupportsCustomImplementation(t *testing.T) {
	xml := transitionSlideXML(t, NewSlide("Custom").WithTransition(customTransition{
		xml: `<p:transition><p:wheel spokes="1"/></p:transition>`,
	}))
	if !strings.Contains(xml, `<p:transition><p:wheel spokes="1"/></p:transition>`) {
		t.Fatalf("expected custom transition XML in slide")
	}
}

func TestSlideTransitionRejectsCustomValidationError(t *testing.T) {
	_, err := CreateWithSlides("Demo", []SlideContent{
		NewSlide("Bad").WithTransition(customTransition{validateErr: errors.New("bad transition")}),
	})
	if err == nil {
		t.Fatalf("expected custom transition validation error")
	}
	if !strings.Contains(err.Error(), "bad transition") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSlideTransitionRejectsMalformedCustomXML(t *testing.T) {
	_, err := CreateWithSlides("Demo", []SlideContent{
		NewSlide("Bad").WithTransition(customTransition{xml: `<p:fade/>`}),
	})
	if err == nil {
		t.Fatalf("expected malformed transition XML error")
	}
	if !strings.Contains(err.Error(), "transition XML must be wrapped") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTransitionParityFixturesAgainstPptRsFragments(t *testing.T) {
	entries := loadTransitionParityFixture(t)
	if len(entries) == 0 {
		t.Fatalf("transition parity fixture is empty")
	}
	for _, entry := range entries {
		transition := TransitionType(entry.Transition)
		slideXML := transitionSlideXML(t, NewSlide("Parity").WithTransition(transition))

		if entry.XML == "" {
			if strings.Contains(slideXML, "<p:transition>") {
				t.Fatalf("did not expect transition XML for %q", entry.Transition)
			}
			continue
		}
		if !strings.Contains(slideXML, entry.XML) {
			t.Fatalf("transition %q missing expected ppt-rs XML fragment %q", entry.Transition, entry.XML)
		}
	}
}

func transitionSlideXML(t *testing.T, slide SlideContent) string {
	t.Helper()
	data, err := CreateWithSlides("Transition Demo", []SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	return readZipFile(t, zr, "ppt/slides/slide1.xml")
}

func loadTransitionParityFixture(t *testing.T) []transitionParityFixtureEntry {
	t.Helper()
	path := filepath.Join("fixtures", "ppt_rs_transition_fragments.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read transition fixture %s: %v", path, err)
	}
	var entries []transitionParityFixtureEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("decode transition fixture %s: %v", path, err)
	}
	return entries
}

func TestTransitionOptions(t *testing.T) {
	cases := []struct {
		name      string
		options   TransitionOptions
		expectXML string
	}{
		{
			name: "push right auto advance",
			options: TransitionOptions{
				Type:                  TransitionPush,
				Direction:             TransitionDirRight,
				DisableAdvanceOnClick: true,
				AdvanceAfterMS:        2000,
				DurationMS:            1500,
			},
			expectXML: `<p:transition advClick="0" advTm="2000" dur="1500"><p:push dir="r"/></p:transition>`,
		},
		{
			name: "fade through black",
			options: TransitionOptions{
				Type:    TransitionFade,
				ThruBlk: true,
			},
			expectXML: `<p:transition><p:fade thruBlk="1"/></p:transition>`,
		},
		{
			name: "split vertical in",
			options: TransitionOptions{
				Type:        TransitionSplit,
				Orientation: TransitionOrientVertical,
				Direction:   TransitionDirIn,
			},
			expectXML: `<p:transition><p:split dir="in" orient="vert"/></p:transition>`,
		},
		{
			name: "wheel spokes",
			options: TransitionOptions{
				Type:       TransitionClock,
				SpokeCount: 8,
			},
			expectXML: `<p:transition><p:wheel spokes="8"/></p:transition>`,
		},
		{
			name: "blinds vertical",
			options: TransitionOptions{
				Type:        TransitionBlinds,
				Orientation: TransitionOrientVertical,
			},
			expectXML: `<p:transition><p:blinds orient="vert"/></p:transition>`,
		},
		{
			name: "randomBar vertical",
			options: TransitionOptions{
				Type:        TransitionRandomBars,
				Orientation: TransitionOrientVertical,
			},
			expectXML: `<p:transition><p:randomBar orient="vert"/></p:transition>`,
		},
		{
			name: "fade advance on click",
			options: TransitionOptions{
				Type: TransitionFade,
				// DisableAdvanceOnClick default is false (enabled)
			},
			expectXML: `<p:transition><p:fade/></p:transition>`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			slideXML := transitionSlideXML(t, NewSlide("Options").WithTransitionOptions(tc.options))
			if !strings.Contains(slideXML, tc.expectXML) {
				t.Errorf("expected transition XML %s not found in: %s", tc.expectXML, slideXML)
			}
		})
	}
}
