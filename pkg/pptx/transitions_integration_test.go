package pptx_test

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
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
func (c customTransition) XML() string     { return c.xml }

func TestCreateWithSlidesRendersRepresentativeTransitions(t *testing.T) {
	cases := []struct {
		name       string
		transition pptx.SlideTransition
		expectXML  string
	}{
		{name: "fade", transition: pptx.TransitionFade, expectXML: `<p:transition><p:fade/></p:transition>`},
		{name: "push", transition: pptx.TransitionPush, expectXML: `<p:transition><p:push dir="r"/></p:transition>`},
		{name: "split", transition: pptx.TransitionSplit, expectXML: `<p:transition><p:split dir="out" orient="horz"/></p:transition>`},
		{name: "zoom", transition: pptx.TransitionZoom, expectXML: `<p:transition><p:zoom dir="in"/></p:transition>`},
		{name: "none", transition: pptx.TransitionNone, expectXML: ``},
		{name: "cut", transition: pptx.TransitionCut, expectXML: ``},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			slideXML := transitionSlideXML(t, pptx.NewSlide("Transition").WithTransition(tc.transition))
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
	_, err := pptx.CreateWithSlides("Demo", []pptx.SlideContent{
		pptx.NewSlide("Bad").WithTransition(pptx.TransitionType("spin")),
	})
	if err == nil {
		t.Fatalf("expected unknown transition validation error")
	}
	if !strings.Contains(err.Error(), `unsupported transition type`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSlideTransitionSupportsCustomImplementation(t *testing.T) {
	xml := transitionSlideXML(t, pptx.NewSlide("Custom").WithTransition(customTransition{
		xml: `<p:transition><p:wheel spokes="1"/></p:transition>`,
	}))
	if !strings.Contains(xml, `<p:transition><p:wheel spokes="1"/></p:transition>`) {
		t.Fatalf("expected custom transition XML in slide")
	}
}

func TestSlideTransitionRejectsCustomValidationError(t *testing.T) {
	_, err := pptx.CreateWithSlides("Demo", []pptx.SlideContent{
		pptx.NewSlide("Bad").WithTransition(customTransition{validateErr: errors.New("bad transition")}),
	})
	if err == nil {
		t.Fatalf("expected custom transition validation error")
	}
	if !strings.Contains(err.Error(), "bad transition") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSlideTransitionRejectsMalformedCustomXML(t *testing.T) {
	_, err := pptx.CreateWithSlides("Demo", []pptx.SlideContent{
		pptx.NewSlide("Bad").WithTransition(customTransition{xml: `<p:fade/>`}),
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
		transition := pptx.TransitionType(entry.Transition)
		slideXML := transitionSlideXML(t, pptx.NewSlide("Parity").WithTransition(transition))

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

func transitionSlideXML(t *testing.T, slide pptx.SlideContent) string {
	t.Helper()
	data, err := pptx.CreateWithSlides("Transition Demo", []pptx.SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	return testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")
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
		options   pptx.TransitionOptions
		expectXML string
	}{
		{
			name: "push right auto advance",
			options: pptx.TransitionOptions{
				Type:                  pptx.TransitionPush,
				Direction:             pptx.TransitionDirRight,
				DisableAdvanceOnClick: true,
				AdvanceAfterMS:        2000,
				DurationMS:            1500,
			},
			expectXML: `<p:transition advClick="0" advTm="2000" dur="1500"><p:push dir="r"/></p:transition>`,
		},
		{
			name: "fade through black",
			options: pptx.TransitionOptions{
				Type:    pptx.TransitionFade,
				ThruBlk: true,
			},
			expectXML: `<p:transition><p:fade thruBlk="1"/></p:transition>`,
		},
		{
			name: "split vertical in",
			options: pptx.TransitionOptions{
				Type:        pptx.TransitionSplit,
				Orientation: pptx.TransitionOrientVertical,
				Direction:   pptx.TransitionDirIn,
			},
			expectXML: `<p:transition><p:split dir="in" orient="vert"/></p:transition>`,
		},
		{
			name: "wheel spokes",
			options: pptx.TransitionOptions{
				Type:       pptx.TransitionClock,
				SpokeCount: 8,
			},
			expectXML: `<p:transition><p:wheel spokes="8"/></p:transition>`,
		},
		{
			name: "blinds vertical",
			options: pptx.TransitionOptions{
				Type:        pptx.TransitionBlinds,
				Orientation: pptx.TransitionOrientVertical,
			},
			expectXML: `<p:transition><p:blinds orient="vert"/></p:transition>`,
		},
		{
			name: "randomBar vertical",
			options: pptx.TransitionOptions{
				Type:        pptx.TransitionRandomBars,
				Orientation: pptx.TransitionOrientVertical,
			},
			expectXML: `<p:transition><p:randomBar orient="vert"/></p:transition>`,
		},
		{
			name: "wipe up",
			options: pptx.TransitionOptions{
				Type:      pptx.TransitionWipe,
				Direction: pptx.TransitionDirUp,
			},
			expectXML: `<p:transition><p:wipe dir="u"/></p:transition>`,
		},
		{
			name: "uncover up-left",
			options: pptx.TransitionOptions{
				Type:      pptx.TransitionUncover,
				Direction: pptx.TransitionDirUpLeft,
			},
			expectXML: `<p:transition><p:pull dir="lu"/></p:transition>`,
		},
		{
			name: "strips up-right",
			options: pptx.TransitionOptions{
				Type:      pptx.TransitionStrips,
				Direction: pptx.TransitionDirUpRight,
			},
			expectXML: `<p:transition><p:strips dir="ru"/></p:transition>`,
		},
		{
			name: "fade advance on click",
			options: pptx.TransitionOptions{
				Type: pptx.TransitionFade,
			},
			expectXML: `<p:transition><p:fade/></p:transition>`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			slideXML := transitionSlideXML(t, pptx.NewSlide("Options").WithTransitionOptions(tc.options))
			if !strings.Contains(slideXML, tc.expectXML) {
				t.Errorf("expected transition XML %s not found in: %s", tc.expectXML, slideXML)
			}
		})
	}
}

func TestTransitionValidation(t *testing.T) {
	cases := []struct {
		name    string
		options pptx.TransitionOptions
		wantErr string
	}{
		{
			name: "invalid push direction",
			options: pptx.TransitionOptions{
				Type:      pptx.TransitionPush,
				Direction: "invalid",
			},
			wantErr: `invalid direction "invalid" for transition "push" (expected u|d|l|r)`,
		},
		{
			name: "invalid zoom direction",
			options: pptx.TransitionOptions{
				Type:      pptx.TransitionZoom,
				Direction: pptx.TransitionDirUp,
			},
			wantErr: `invalid direction "u" for transition "zoom" (expected in|out)`,
		},
		{
			name: "strips valid direction",
			options: pptx.TransitionOptions{
				Type:      pptx.TransitionStrips,
				Direction: pptx.TransitionDirUpLeft,
			},
			wantErr: "",
		},
		{
			name: "strips invalid direction",
			options: pptx.TransitionOptions{
				Type:      pptx.TransitionStrips,
				Direction: pptx.TransitionDirUp,
			},
			wantErr: `invalid direction "u" for transition "strips" (expected ul|ur|dl|dr)`,
		},
		{
			name: "fade no support for direction",
			options: pptx.TransitionOptions{
				Type:      pptx.TransitionFade,
				Direction: pptx.TransitionDirUp,
			},
			wantErr: `transition "fade" does not support direction`,
		},
		{
			name: "invalid wheel spoke count for non-wheel",
			options: pptx.TransitionOptions{
				Type:       pptx.TransitionFade,
				SpokeCount: 4,
			},
			wantErr: `transition "fade" does not support spoke count`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.options.Validate()
			if tc.wantErr == "" {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error %q, got nil", tc.wantErr)
				} else if !strings.Contains(err.Error(), tc.wantErr) {
					t.Errorf("expected error %q, got %q", tc.wantErr, err.Error())
				}
			}
		})
	}
}
