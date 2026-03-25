package animations

import (
	"errors"
	"fmt"
)

// AnimationEffect defines the type of animation effect.
type AnimationEffect string

const (
	// AnimationEntranceAppear starts the entrance-effects group.
	classEntr                                    = "entr"
	classExit                                    = "exit"
	classEmph                                    = "emph"
	classPath                                    = "path"
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

const (
	defaultAnimationDurationMS   uint32 = 500
	presetIDAppear               uint32 = 1
	presetIDFly                  uint32 = 2
	presetIDFade                 uint32 = 10
	presetIDFloat                uint32 = 14
	presetIDRandomBars           uint32 = 15
	presetIDSplit                uint32 = 16
	presetIDShape                uint32 = 17
	presetIDSwivel               uint32 = 19
	presetIDWheel                uint32 = 21
	presetIDWipe                 uint32 = 22
	presetIDZoom                 uint32 = 23
	presetIDBounce               uint32 = 25
	presetIDGrowAndTurn          uint32 = 26
	presetIDEmphasisPulse        uint32 = 31
	presetIDColorPulse           uint32 = 32
	presetIDTeeter               uint32 = 33
	presetIDSpin                 uint32 = 34
	presetIDGrowShrink           uint32 = 35
	presetIDDesaturate           uint32 = 36
	presetIDDarken               uint32 = 37
	presetIDLighten              uint32 = 38
	presetIDTransparency         uint32 = 39
	presetIDObjectColor          uint32 = 40
	presetIDPathLines            uint32 = 42
	presetIDPathArcs             uint32 = 43
	presetIDPathTurns            uint32 = 44
	presetIDPathShapes           uint32 = 45
	presetIDPathLoops            uint32 = 46
	presetIDPathCustom           uint32 = 47
	presetSubtypeFromTop                = 1
	presetSubtypeFromRight              = 2
	presetSubtypeFromTopRight           = 3
	presetSubtypeFromLeft               = 4
	presetSubtypeFromTopLeft            = 5
	presetSubtypeFromBottomRight        = 6
	presetSubtypeFromBottomLeft         = 7
	presetSubtypeFromBottom             = 8
)

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

func (a Animation) XML(seqID int, actualShapeID int) string {
	repeatAttr := ""
	if a.RepeatCount > 0 {
		repeatAttr = fmt.Sprintf(` repeatCount="%d"`, a.RepeatCount*1000) //nolint:mnd // OOXML scale
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

// AnimationDefinition is the interface for types that can be converted to an Animation.
type AnimationDefinition interface {
	ToAnimation() Animation
}

func (a Animation) ToAnimation() Animation {
	return a
}
