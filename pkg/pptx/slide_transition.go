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
	// TransitionPush applies a push transition.
	TransitionPush TransitionType = "push"
	// TransitionWipe applies a wipe transition.
	TransitionWipe TransitionType = "wipe"
	// TransitionSplit applies a split transition.
	TransitionSplit TransitionType = "split"
	// TransitionReveal applies a reveal transition.
	TransitionReveal TransitionType = "reveal"
	// TransitionCover applies a cover transition.
	TransitionCover TransitionType = "cover"
	// TransitionZoom applies a zoom transition.
	TransitionZoom TransitionType = "zoom"
	// TransitionRandomBars applies a random bars transition.
	TransitionRandomBars TransitionType = "randomBar"
	// TransitionShape applies a shape transition.
	TransitionShape TransitionType = "circle"
	// TransitionUncover applies an uncover transition.
	TransitionUncover TransitionType = "pull"
	// TransitionFlash applies a flash transition.
	TransitionFlash TransitionType = "flash"
	// TransitionStrips applies a strips transition.
	TransitionStrips TransitionType = "strips"
	// TransitionBlinds applies a blinds transition.
	TransitionBlinds TransitionType = "blinds"
	// TransitionClock applies a clock transition.
	TransitionClock TransitionType = "wheel"
	// TransitionRipple applies a ripple transition.
	TransitionRipple TransitionType = "ripple"
	// TransitionHoneycomb applies a honeycomb transition.
	TransitionHoneycomb TransitionType = "honeycomb"
	// TransitionGlitter applies a glitter transition.
	TransitionGlitter TransitionType = "glitter"
	// TransitionVortex applies a vortex transition.
	TransitionVortex TransitionType = "vortex"
	// TransitionShred applies a shred transition.
	TransitionShred TransitionType = "shred"
	// TransitionSwitch applies a switch transition.
	TransitionSwitch TransitionType = "switch"
	// TransitionFlip applies a flip transition.
	TransitionFlip TransitionType = "flip"
	// TransitionGallery applies a gallery transition.
	TransitionGallery TransitionType = "gallery"
	// TransitionCube applies a cube transition.
	TransitionCube TransitionType = "cube"
	// TransitionDoors applies a doors transition.
	TransitionDoors TransitionType = "doors"
	// TransitionBox applies a box transition.
	TransitionBox TransitionType = "box"
	// TransitionRandom applies a random transition.
	TransitionRandom TransitionType = "random"
)

// TransitionDirection defines the direction of a transition.
type TransitionDirection string

const (
	TransitionDirIn        TransitionDirection = "in"
	TransitionDirOut       TransitionDirection = "out"
	TransitionDirUp        TransitionDirection = "u"
	TransitionDirDown      TransitionDirection = "d"
	TransitionDirLeft      TransitionDirection = "l"
	TransitionDirRight     TransitionDirection = "r"
	TransitionDirUpLeft    TransitionDirection = "ul"
	TransitionDirUpRight   TransitionDirection = "ur"
	TransitionDirDownLeft  TransitionDirection = "dl"
	TransitionDirDownRight TransitionDirection = "dr"
)

// TransitionOptions provides advanced configuration for a slide transition.
type TransitionOptions struct {
	Type                  TransitionType
	DurationMS            uint32
	Direction             TransitionDirection
	DisableAdvanceOnClick bool
	AdvanceAfterMS        uint32 // 0 means no auto-advance
}

// Validate checks whether the transition options are valid.
func (o TransitionOptions) Validate() error {
	if err := o.Type.Validate(); err != nil {
		return err
	}
	return nil
}

// XML renders the transition XML fragment with options.
func (o TransitionOptions) XML() string {
	if o.Type == TransitionNone || o.Type == TransitionCut || o.Type == "" {
		return ""
	}

	advClick := ""
	if o.DisableAdvanceOnClick {
		advClick = ` advClick="0"`
	}
	advTm := ""
	if o.AdvanceAfterMS > 0 {
		advTm = fmt.Sprintf(` advTm="%d"`, o.AdvanceAfterMS)
	}

	dir := string(o.Direction)
	if dir == "" {
		// Defaults aligned with ppt-rs
		switch o.Type {
		case TransitionPush, TransitionWipe, TransitionReveal, TransitionCover:
			dir = "r"
		case TransitionZoom:
			dir = "in"
		case TransitionSplit:
			// Split in ppt-rs uses orient and dir
			return fmt.Sprintf(`<p:transition%s%s><p:split dir="out" orient="horz"/></p:transition>`, advClick, advTm)
		default:
			return fmt.Sprintf(`<p:transition%s%s><p:%s/></p:transition>`, advClick, advTm, string(o.Type))
		}
	}

	return fmt.Sprintf(`<p:transition%s%s><p:%s dir="%s"/></p:transition>`, advClick, advTm, string(o.Type), dir)
}

// WithTransition sets the transition behavior for a slide.
func (s SlideContent) WithTransition(transition SlideTransition) SlideContent {
	s.Transition = transition
	return s
}

// WithTransitionOptions sets advanced transition behavior for a slide.
func (s SlideContent) WithTransitionOptions(options TransitionOptions) SlideContent {
	s.Transition = options
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
		TransitionZoom,
		TransitionRandomBars,
		TransitionShape,
		TransitionUncover,
		TransitionFlash,
		TransitionStrips,
		TransitionBlinds,
		TransitionClock,
		TransitionRipple,
		TransitionHoneycomb,
		TransitionGlitter,
		TransitionVortex,
		TransitionShred,
		TransitionSwitch,
		TransitionFlip,
		TransitionGallery,
		TransitionCube,
		TransitionDoors,
		TransitionBox,
		TransitionRandom:
		return nil
	default:
		return fmt.Errorf("unsupported transition type: %q", string(t))
	}
}

// XML renders the transition XML fragment.
func (t TransitionType) XML() string {
	switch t {
	case TransitionFade, TransitionFlash, TransitionRandomBars, TransitionRandom,
		TransitionRipple, TransitionHoneycomb, TransitionGlitter, TransitionVortex,
		TransitionShred, TransitionSwitch, TransitionFlip, TransitionGallery,
		TransitionCube, TransitionDoors, TransitionBox:
		return fmt.Sprintf(`<p:transition><p:%s/></p:transition>`, string(t))
	case TransitionPush, TransitionWipe, TransitionReveal, TransitionCover:
		return fmt.Sprintf(`<p:transition><p:%s dir="r"/></p:transition>`, string(t))
	case TransitionSplit:
		return `<p:transition><p:split dir="out" orient="horz"/></p:transition>`
	case TransitionZoom:
		return `<p:transition><p:zoom dir="in"/></p:transition>`
	case TransitionShape:
		return `<p:transition><p:circle/></p:transition>`
	case TransitionUncover:
		return `<p:transition><p:pull dir="l"/></p:transition>`
	case TransitionClock:
		return `<p:transition><p:wheel/></p:transition>`
	case TransitionBlinds:
		return `<p:transition><p:blinds/></p:transition>`
	case TransitionStrips:
		return `<p:transition><p:strips/></p:transition>`
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
