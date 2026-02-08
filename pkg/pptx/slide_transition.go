package pptx

import (
	"fmt"
	"strings"
)

// SlideTransition is the extensibility contract for slide transitions.
//
// Implementations must return valid transition XML fragments for `XML`.
// Validation must fail fast for unsupported or invalid transition settings.
type SlideTransition interface {
	Validate() error
	XML() string
}

// TransitionType is the built-in transition enum aligned to ppt-rs transition kinds.
type TransitionType string

const (
	// TransitionNone applies no explicit transition node.
	TransitionNone TransitionType = "none"
	// TransitionCut applies an instant cut transition (default behavior, no node emitted).
	TransitionCut TransitionType = "cut"
	// TransitionFade applies a fade transition.
	TransitionFade TransitionType = "fade"
	// TransitionPush applies a right-direction push transition.
	TransitionPush TransitionType = "push"
	// TransitionWipe applies a right-direction wipe transition.
	TransitionWipe TransitionType = "wipe"
	// TransitionSplit applies a horizontal outward split transition.
	TransitionSplit TransitionType = "split"
	// TransitionReveal applies a right-direction reveal transition.
	TransitionReveal TransitionType = "reveal"
	// TransitionCover applies a right-direction cover transition.
	TransitionCover TransitionType = "cover"
	// TransitionZoom applies an inward zoom transition.
	TransitionZoom TransitionType = "zoom"
)

// WithTransition sets the transition behavior for a slide.
//
// Example:
//
//	slide := NewSlide("Next").WithTransition(TransitionPush)
func (s SlideContent) WithTransition(transition SlideTransition) SlideContent {
	s.Transition = transition
	return s
}

// Validate checks whether a transition type is supported.
func (t TransitionType) Validate() error {
	switch t {
	case TransitionNone,
		TransitionCut,
		TransitionFade,
		TransitionPush,
		TransitionWipe,
		TransitionSplit,
		TransitionReveal,
		TransitionCover,
		TransitionZoom:
		return nil
	default:
		return fmt.Errorf(
			"transition must be one of none|cut|fade|push|wipe|split|reveal|cover|zoom, got %q",
			string(t),
		)
	}
}

// XML renders the transition XML fragment.
func (t TransitionType) XML() string {
	switch t {
	case TransitionFade:
		return `<p:transition><p:fade/></p:transition>`
	case TransitionPush:
		return `<p:transition><p:push dir="r"/></p:transition>`
	case TransitionWipe:
		return `<p:transition><p:wipe dir="r"/></p:transition>`
	case TransitionSplit:
		return `<p:transition><p:split dir="out" orient="horz"/></p:transition>`
	case TransitionReveal:
		return `<p:transition><p:reveal dir="r"/></p:transition>`
	case TransitionCover:
		return `<p:transition><p:cover dir="r"/></p:transition>`
	case TransitionZoom:
		return `<p:transition><p:zoom dir="in"/></p:transition>`
	default:
		return ""
	}
}

func validateSlideTransition(s SlideContent, index int) error {
	if s.Transition == nil {
		return nil
	}
	if err := s.Transition.Validate(); err != nil {
		return fmt.Errorf("slide %d transition: %w", index, err)
	}
	xml := strings.TrimSpace(s.Transition.XML())
	if xml == "" {
		return nil
	}
	if !strings.HasPrefix(xml, "<p:transition") || !strings.HasSuffix(xml, "</p:transition>") {
		return fmt.Errorf("slide %d transition XML must be wrapped in <p:transition>...</p:transition>", index)
	}
	return nil
}

func slideTransitionXML(s SlideContent) string {
	if s.Transition == nil {
		return ""
	}
	return strings.TrimSpace(s.Transition.XML())
}
