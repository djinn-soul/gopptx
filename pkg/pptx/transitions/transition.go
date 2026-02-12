package transitions

import (
	"fmt"
	"strings"
)

// SlideTransition is the extensibility contract for slide transitions.
type SlideTransition interface {
	Validate() error
	XML() string
}

// TransitionType is the built-in transition enum.
type TransitionType string

const (
	TransitionNone       TransitionType = "none"
	TransitionCut        TransitionType = "cut"
	TransitionFade       TransitionType = "fade"
	TransitionPush       TransitionType = "push"
	TransitionWipe       TransitionType = "wipe"
	TransitionSplit      TransitionType = "split"
	TransitionReveal     TransitionType = "reveal"
	TransitionCover      TransitionType = "cover"
	TransitionZoom       TransitionType = "zoom"
	TransitionRandomBars TransitionType = "randomBar"
	TransitionShape      TransitionType = "circle"
	TransitionUncover    TransitionType = "pull"
	TransitionFlash      TransitionType = "flash"
	TransitionStrips     TransitionType = "strips"
	TransitionBlinds     TransitionType = "blinds"
	TransitionClock      TransitionType = "wheel"
	TransitionRipple     TransitionType = "ripple"
	TransitionHoneycomb  TransitionType = "honeycomb"
	TransitionGlitter    TransitionType = "glitter"
	TransitionVortex     TransitionType = "vortex"
	TransitionShred      TransitionType = "shred"
	TransitionSwitch     TransitionType = "switch"
	TransitionFlip       TransitionType = "flip"
	TransitionGallery    TransitionType = "gallery"
	TransitionCube       TransitionType = "cube"
	TransitionDoors      TransitionType = "doors"
	TransitionBox        TransitionType = "box"
	TransitionRandom     TransitionType = "random"
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
	TransitionDirUpLeft    TransitionDirection = "lu"
	TransitionDirUpRight   TransitionDirection = "ru"
	TransitionDirDownLeft  TransitionDirection = "ld"
	TransitionDirDownRight TransitionDirection = "rd"
)

// TransitionOrientation defines the orientation of a transition.
type TransitionOrientation string

const (
	TransitionOrientHorizontal TransitionOrientation = "horz"
	TransitionOrientVertical   TransitionOrientation = "vert"
)

// TransitionSound defines the audio configuration for a transition.
type TransitionSound struct {
	RelID string // Relationship ID for the audio file (required)
	Name  string // Display name (e.g., "Applause")
	Loop  bool   // Whether to loop the sound
}

// TransitionOptions provides advanced configuration for a slide transition.
type TransitionOptions struct {
	Type                  TransitionType
	DurationMS            uint32
	Direction             TransitionDirection
	Orientation           TransitionOrientation
	SpokeCount            uint32
	ThruBlk               bool
	Sound                 *TransitionSound
	DisableAdvanceOnClick bool
	AdvanceAfterMS        uint32
}

func (o TransitionOptions) Validate() error {
	if err := o.Type.Validate(); err != nil {
		return err
	}

	if o.Sound != nil {
		if strings.TrimSpace(o.Sound.RelID) == "" {
			return fmt.Errorf("transition sound requires a valid relationship ID")
		}
	}

	if o.Orientation != "" {
		switch o.Type {
		case TransitionSplit, TransitionBlinds, TransitionRandomBars:
			if o.Orientation != TransitionOrientHorizontal && o.Orientation != TransitionOrientVertical {
				return fmt.Errorf("invalid orientation %q for transition %q", o.Orientation, o.Type)
			}
		default:
			return fmt.Errorf("transition %q does not support orientation", o.Type)
		}
	}

	if o.SpokeCount > 0 && o.Type != TransitionClock {
		return fmt.Errorf("transition %q does not support spoke count", o.Type)
	}

	if o.Direction != "" {
		switch o.Type {
		case TransitionPush, TransitionWipe, TransitionReveal, TransitionCover:
			switch o.Direction {
			case TransitionDirUp, TransitionDirDown, TransitionDirLeft, TransitionDirRight:
			default:
				return fmt.Errorf("invalid direction %q for transition %q (expected u|d|l|r)", o.Direction, o.Type)
			}
		case TransitionZoom, TransitionSplit:
			switch o.Direction {
			case TransitionDirIn, TransitionDirOut:
			default:
				return fmt.Errorf("invalid direction %q for transition %q (expected in|out)", o.Direction, o.Type)
			}
		case TransitionUncover:
			switch o.Direction {
			case TransitionDirUp, TransitionDirDown, TransitionDirLeft, TransitionDirRight,
				TransitionDirUpLeft, TransitionDirUpRight, TransitionDirDownLeft, TransitionDirDownRight:
			default:
				return fmt.Errorf("invalid direction %q for transition %q (expected u|d|l|r|lu|ru|ld|rd)", o.Direction, o.Type)
			}
		case TransitionStrips:
			switch o.Direction {
			case TransitionDirUpLeft, TransitionDirUpRight, TransitionDirDownLeft, TransitionDirDownRight:
			default:
				return fmt.Errorf("invalid direction %q for transition %q (expected ul|ur|dl|dr)", o.Direction, o.Type)
			}
		default:
			return fmt.Errorf("transition %q does not support direction", o.Type)
		}
	}
	return nil
}

func (o TransitionOptions) XML() string {
	if o.Type == TransitionNone || (o.Type == TransitionCut && o.Sound == nil) {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<p:transition`)
	if o.DisableAdvanceOnClick {
		b.WriteString(` advClick="0"`)
	}
	if o.AdvanceAfterMS > 0 {
		fmt.Fprintf(&b, ` advTm="%d"`, o.AdvanceAfterMS)
	}
	if o.DurationMS > 0 {
		fmt.Fprintf(&b, ` dur="%d"`, o.DurationMS)
	}
	b.WriteString(`>`)

	if o.Sound != nil {
		b.WriteString(`<p:sndAc><p:stSnd`)
		if o.Sound.Loop {
			b.WriteString(` loop="1"`)
		}
		b.WriteString(`><p:snd r:embed="` + escape(o.Sound.RelID) + `"`)
		if o.Sound.Name != "" {
			b.WriteString(` name="` + escape(o.Sound.Name) + `"`)
		}
		b.WriteString(`/>`)
		b.WriteString(`</p:stSnd></p:sndAc>`)
	}

	b.WriteString(`<p:`)
	b.WriteString(string(o.Type))

	if o.Direction != "" {
		fmt.Fprintf(&b, ` dir="%s"`, o.Direction)
	}
	if o.Orientation != "" {
		fmt.Fprintf(&b, ` orient="%s"`, o.Orientation)
	}
	if o.SpokeCount > 0 {
		fmt.Fprintf(&b, ` spokes="%d"`, o.SpokeCount)
	}
	if o.ThruBlk {
		b.WriteString(` thruBlk="1"`)
	}
	b.WriteString(`/></p:transition>`)
	return b.String()
}

func escape(value string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&apos;",
	)
	return replacer.Replace(value)
}

func (t TransitionType) Validate() error {
	switch t {
	case TransitionNone, TransitionCut, TransitionFade, TransitionPush, TransitionWipe,
		TransitionSplit, TransitionReveal, TransitionCover, TransitionZoom,
		TransitionRandomBars, TransitionShape, TransitionUncover, TransitionFlash,
		TransitionStrips, TransitionBlinds, TransitionClock, TransitionRipple,
		TransitionHoneycomb, TransitionGlitter, TransitionVortex, TransitionShred,
		TransitionSwitch, TransitionFlip, TransitionGallery, TransitionCube,
		TransitionDoors, TransitionBox, TransitionRandom:
		return nil
	default:
		return fmt.Errorf("unsupported transition type: %q", string(t))
	}
}

func (t TransitionType) XML() string {
	switch t {
	case TransitionNone, TransitionCut:
		// TransitionCut is the default and requires no XML unless options (like sound) are set.
		return ""
	case TransitionPush:
		return `<p:transition><p:push dir="r"/></p:transition>`
	case TransitionWipe:
		return `<p:transition><p:wipe dir="r"/></p:transition>`
	case TransitionSplit:
		return `<p:transition><p:split dir="out" orient="horz"/></p:transition>`
	case TransitionZoom:
		return `<p:transition><p:zoom dir="in"/></p:transition>`
	case TransitionFade:
		return `<p:transition><p:fade/></p:transition>`
	case TransitionReveal:
		return `<p:transition><p:reveal dir="r"/></p:transition>`
	case TransitionCover:
		return `<p:transition><p:cover dir="r"/></p:transition>`
	default:
		return fmt.Sprintf(`<p:transition><p:%s/></p:transition>`, t)
	}
}
