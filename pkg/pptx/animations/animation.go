package animations

import (
	"errors"
	"fmt"
	"strings"
)

// AnimationEffect defines the type of animation effect.
type AnimationEffect string

const (
	// AnimationEntranceAppear starts the entrance-effects group.
	classEntr                                    = "entr"
	AnimationEntranceAppear      AnimationEffect = "entr_appear"
	AnimationEntranceFade        AnimationEffect = "entr_fade"
	AnimationEntranceFlyIn       AnimationEffect = "entr_flyIn"
	AnimationEntranceFloat       AnimationEffect = "entr_float"
	AnimationEntranceSplit       AnimationEffect = "entr_split"
	AnimationEntranceWipe        AnimationEffect = "entr_wipe"
	AnimationEntranceShape       AnimationEffect = "entr_shape"
	AnimationEntranceWheel       AnimationEffect = "entr_wheel"
	AnimationEntranceRandomBars  AnimationEffect = "entr_randomBars"
	AnimationEntranceGrowAndTurn AnimationEffect = "entr_growAndTurn"
	AnimationEntranceZoom        AnimationEffect = "entr_zoom"
	AnimationEntranceSwivel      AnimationEffect = "entr_swivel"
	AnimationEntranceBounce      AnimationEffect = "entr_bounce"

	// AnimationExitDisappear starts the exit-effects group.
	AnimationExitDisappear AnimationEffect = "exit_disappear"
	AnimationExitFadeOut   AnimationEffect = "exit_fade"
	AnimationExitFlyOut    AnimationEffect = "exit_flyOut"
	AnimationExitFloatOut  AnimationEffect = "exit_float"

	// AnimationEmphasisPulse starts the emphasis-effects group.
	AnimationEmphasisPulse        AnimationEffect = "emph_pulse"
	AnimationEmphasisColorPulse   AnimationEffect = "emph_colorPulse"
	AnimationEmphasisTeeter       AnimationEffect = "emph_teeter"
	AnimationEmphasisSpin         AnimationEffect = "emph_spin"
	AnimationEmphasisGrowShrink   AnimationEffect = "emph_growShrink"
	AnimationEmphasisDesaturate   AnimationEffect = "emph_desaturate"
	AnimationEmphasisDarken       AnimationEffect = "emph_darken"
	AnimationEmphasisLighten      AnimationEffect = "emph_lighten"
	AnimationEmphasisTransparency AnimationEffect = "emph_transparency"
	AnimationEmphasisObjectColor  AnimationEffect = "emph_objectColor"

	// AnimationPathLines starts the motion-path effects group.
	AnimationPathLines  AnimationEffect = "path_lines"
	AnimationPathArcs   AnimationEffect = "path_arcs"
	AnimationPathTurns  AnimationEffect = "path_turns"
	AnimationPathShapes AnimationEffect = "path_shapes"
	AnimationPathLoops  AnimationEffect = "path_loops"
	AnimationPathCustom AnimationEffect = "path_custom"
)

// AnimationTrigger defines when the animation starts.
type AnimationTrigger string

const (
	AnimationOnClick       AnimationTrigger = "onClick"
	AnimationWithPrevious  AnimationTrigger = "withPrev"
	AnimationAfterPrevious AnimationTrigger = "afterPrev"
)

// AnimationDirection defines the direction or subtype of the animation.
type AnimationDirection string

const (
	AnimationDirIn        AnimationDirection = "in"
	AnimationDirOut       AnimationDirection = "out"
	AnimationDirUp        AnimationDirection = "u"
	AnimationDirDown      AnimationDirection = "d"
	AnimationDirLeft      AnimationDirection = "l"
	AnimationDirRight     AnimationDirection = "r"
	AnimationDirUpLeft    AnimationDirection = "ul"
	AnimationDirUpRight   AnimationDirection = "ur"
	AnimationDirDownLeft  AnimationDirection = "dl"
	AnimationDirDownRight AnimationDirection = "dr"
)

// Animation represents a single animation effect on a slide object.
type Animation struct {
	// ShapeIndex is the 1-based index of the target element.
	// This index refers to a unified list of slide elements in the following order:
	// 1. Custom Shapes (s.Shapes)
	// 2. Connectors (s.Connectors)
	// 3. Images (s.Images)
	// 4. Placeholder Overrides (s.PlaceholderOverrides)
	ShapeIndex  int
	Effect      AnimationEffect
	Trigger     AnimationTrigger
	Direction   AnimationDirection
	DurationMS  uint32
	DelayMS     uint32
	RepeatCount uint32
	AutoReverse bool
}

const defaultAnimationDurationMS = 500

// NewAnimation creates a new animation with default settings (500ms duration, OnClick).
func NewAnimation(shapeIndex int, effect AnimationEffect) Animation {
	return Animation{
		ShapeIndex: shapeIndex,
		Effect:     effect,
		Trigger:    AnimationOnClick,
		DurationMS: defaultAnimationDurationMS,
	}
}

// WithTrigger sets the trigger for an animation.
func (a Animation) WithTrigger(trigger AnimationTrigger) Animation {
	a.Trigger = trigger
	return a
}

// WithDuration sets the duration of the animation in milliseconds.
func (a Animation) WithDuration(durationMS uint32) Animation {
	a.DurationMS = durationMS
	return a
}

// WithDelay sets the delay before the animation starts in milliseconds.
func (a Animation) WithDelay(delayMS uint32) Animation {
	a.DelayMS = delayMS
	return a
}

// WithRepeat sets the number of times the animation repeats.
func (a Animation) WithRepeat(count uint32) Animation {
	a.RepeatCount = count
	return a
}

// WithAutoReverse enables or disables auto-reverse for the animation.
func (a Animation) WithAutoReverse(autoReverse bool) Animation {
	a.AutoReverse = autoReverse
	return a
}

func (a Animation) Validate() error {
	if a.ShapeIndex <= 0 {
		return errors.New("animation shape index must be greater than zero")
	}
	if err := a.Effect.Validate(); err != nil {
		return err
	}
	if err := a.Trigger.Validate(); err != nil {
		return err
	}
	return nil
}

func (e AnimationEffect) Validate() error {
	switch e {
	case AnimationEntranceAppear, AnimationEntranceFade, AnimationEntranceFlyIn, AnimationEntranceFloat,
		AnimationEntranceSplit, AnimationEntranceWipe, AnimationEntranceShape, AnimationEntranceWheel,
		AnimationEntranceRandomBars, AnimationEntranceGrowAndTurn, AnimationEntranceZoom, AnimationEntranceSwivel,
		AnimationEntranceBounce, AnimationExitDisappear, AnimationExitFadeOut, AnimationExitFlyOut,
		AnimationExitFloatOut, AnimationEmphasisPulse, AnimationEmphasisColorPulse, AnimationEmphasisTeeter,
		AnimationEmphasisSpin, AnimationEmphasisGrowShrink, AnimationEmphasisDesaturate, AnimationEmphasisDarken,
		AnimationEmphasisLighten, AnimationEmphasisTransparency, AnimationEmphasisObjectColor, AnimationPathLines,
		AnimationPathArcs, AnimationPathTurns, AnimationPathShapes, AnimationPathLoops, AnimationPathCustom:
		return nil
	default:
		return fmt.Errorf("unsupported animation effect: %q", string(e))
	}
}

func (t AnimationTrigger) Validate() error {
	switch t {
	case AnimationOnClick, AnimationWithPrevious, AnimationAfterPrevious:
		return nil
	default:
		return fmt.Errorf("unsupported animation trigger: %q", string(t))
	}
}

func (a Animation) PresetID() uint32 {
	// NOTE: This switch statement avoids high-frequency map allocation. Do not convert to a map lookup.
	switch a.PresetClass() {
	case "entr", "exit":
		return a.presetIDEntranceExit()
	case "emph":
		return a.presetIDEmphasis()
	case "path":
		return a.presetIDPath()
	}
	return 0
}

//nolint:mnd // Preset IDs are from OOXML spec
func (a Animation) presetIDEntranceExit() uint32 {
	switch a.Effect {
	case AnimationEntranceAppear, AnimationExitDisappear:
		return 1
	case AnimationEntranceFade, AnimationExitFadeOut:
		return 10
	case AnimationEntranceFlyIn, AnimationExitFlyOut:
		return 2
	case AnimationEntranceFloat, AnimationExitFloatOut:
		return 14
	case AnimationEntranceSplit:
		return 16
	case AnimationEntranceWipe:
		return 22
	case AnimationEntranceShape:
		return 17
	case AnimationEntranceWheel:
		return 21
	case AnimationEntranceRandomBars:
		return 15
	case AnimationEntranceGrowAndTurn:
		return 26
	case AnimationEntranceZoom:
		return 23
	case AnimationEntranceSwivel:
		return 19
	case AnimationEntranceBounce:
		return 25
	default:
		return 0
	}
}

//nolint:mnd // Preset IDs are from OOXML spec
func (a Animation) presetIDEmphasis() uint32 {
	switch a.Effect {
	case AnimationEmphasisPulse:
		return 31
	case AnimationEmphasisColorPulse:
		return 32
	case AnimationEmphasisTeeter:
		return 33
	case AnimationEmphasisSpin:
		return 34
	case AnimationEmphasisGrowShrink:
		return 35
	case AnimationEmphasisDesaturate:
		return 36
	case AnimationEmphasisDarken:
		return 37
	case AnimationEmphasisLighten:
		return 38
	case AnimationEmphasisTransparency:
		return 39
	case AnimationEmphasisObjectColor:
		return 40
	default:
		return 0
	}
}

//nolint:mnd // Preset IDs are from OOXML spec
func (a Animation) presetIDPath() uint32 {
	switch a.Effect {
	case AnimationPathLines:
		return 42
	case AnimationPathArcs:
		return 43
	case AnimationPathTurns:
		return 44
	case AnimationPathShapes:
		return 45
	case AnimationPathLoops:
		return 46
	case AnimationPathCustom:
		return 47
	default:
		return 0
	}
}

func (a Animation) PresetClass() string {
	switch {
	case strings.HasPrefix(string(a.Effect), "entr_"):
		return classEntr
	case strings.HasPrefix(string(a.Effect), "exit_"):
		return "exit"
	case strings.HasPrefix(string(a.Effect), "emph_"):
		return "emph"
	case strings.HasPrefix(string(a.Effect), "path_"):
		return "path"
	default:
		return classEntr
	}
}

func (a Animation) XML(seqID int, actualShapeID int) string {
	repeatAttr := ""
	if a.RepeatCount > 0 {
		const repeatMultiplier = 1000
		repeatAttr = fmt.Sprintf(` repeatCount="%d"`, a.RepeatCount*repeatMultiplier)
	}
	reverseAttr := ""
	if a.AutoReverse {
		reverseAttr = ` autoRev="1"`
	}

	return fmt.Sprintf(`
<p:par>
  <p:cTn id="%d" presetID="%d" presetClass="%s" presetSubtype="%d" fill="hold" nodeType="%s">
    <p:stCondLst>
      <p:cond delay="%d"/>
    </p:stCondLst>
    <p:childTnLst>
      <p:set>
        <p:cBhvr>
          <p:cTn id="%d" dur="%d" fill="hold"%s%s>
            <p:stCondLst><p:cond delay="0"/></p:stCondLst>
          </p:cTn>
          <p:tgtEl>
            <p:spTgt spid="%d"/>
          </p:tgtEl>
        </p:cBhvr>
      </p:set>
    </p:childTnLst>
  </p:cTn>
</p:par>`,
		seqID,
		a.PresetID(),
		a.PresetClass(),
		a.PresetSubtype(),
		a.NodeType(),
		a.DelayMS,
		seqID+1,
		a.DurationMS,
		repeatAttr,
		reverseAttr,
		actualShapeID,
	)
}

//nolint:mnd // Preset subtypes are from OOXML spec
func (a Animation) PresetSubtype() int {
	// Mapping based on MS-PPTX / OOXML standards for common effects.
	switch a.Effect {
	case AnimationEntranceFlyIn, AnimationExitFlyOut:
		switch a.Direction {
		case AnimationDirDown:
			return 1 // From Top
		case AnimationDirLeft:
			return 2 // From Right
		case AnimationDirRight:
			return 4 // From Left
		case AnimationDirUp:
			return 8 // From Bottom
		case AnimationDirDownLeft:
			return 3 // From Top-Right
		case AnimationDirDownRight:
			return 5 // From Top-Left
		case AnimationDirUpLeft:
			return 6 // From Bottom-Right
		case AnimationDirUpRight:
			return 7 // From Bottom-Left
		default:
			return 8 // Default "From Bottom"
		}
	case AnimationEntranceWipe, AnimationEntranceFloat:
		switch a.Direction {
		case AnimationDirUp:
			return 4
		case AnimationDirDown:
			return 8
		case AnimationDirLeft:
			return 1
		case AnimationDirRight:
			return 2
		default:
			return 4
		}
	case AnimationEntranceSplit:
		switch a.Direction {
		case AnimationDirIn:
			return 1
		case AnimationDirOut:
			return 2
		default:
			return 2
		}
	default:
		return 0
	}
}

func (a Animation) NodeType() string {
	switch a.Trigger {
	case AnimationOnClick:
		return "clickEffect"
	case AnimationWithPrevious:
		return "withEffect"
	case AnimationAfterPrevious:
		return "afterEffect"
	default:
		return "clickEffect"
	}
}

// AnimationDefinition is the interface for types that can be converted to an Animation.
type AnimationDefinition interface {
	ToAnimation() Animation
}

func (a Animation) ToAnimation() Animation {
	return a
}
