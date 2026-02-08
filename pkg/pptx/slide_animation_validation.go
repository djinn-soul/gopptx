package pptx

import (
	"fmt"
)

func validateSlideAnimations(s SlideContent, slideIndex int) error {
	for i, anim := range s.Animations {
		if anim.ShapeIndex < 1 || anim.ShapeIndex > len(s.Shapes) {
			return fmt.Errorf("slide %d animation %d: shape index %d out of bounds (1-%d)",
				slideIndex, i+1, anim.ShapeIndex, len(s.Shapes))
		}

		if anim.Effect == "" {
			return fmt.Errorf("slide %d animation %d: effect cannot be empty", slideIndex, i+1)
		}
		if err := validateAnimationEffect(anim.Effect); err != nil {
			return fmt.Errorf("slide %d animation %d: %w", slideIndex, i+1, err)
		}
		if err := validateAnimationTrigger(anim.Trigger); err != nil {
			return fmt.Errorf("slide %d animation %d: %w", slideIndex, i+1, err)
		}
	}
	return nil
}

func validateAnimationEffect(effect AnimationEffect) error {
	switch effect {
	case AnimationEntranceAppear,
		AnimationEntranceFade,
		AnimationEntranceFlyIn,
		AnimationEntranceFloat,
		AnimationEntranceSplit,
		AnimationEntranceWipe,
		AnimationEntranceShape,
		AnimationEntranceWheel,
		AnimationEntranceRandomBars,
		AnimationEntranceGrowAndTurn,
		AnimationEntranceZoom,
		AnimationEntranceSwivel,
		AnimationEntranceBounce,
		AnimationExitDisappear,
		AnimationExitFadeOut,
		AnimationExitFlyOut,
		AnimationExitFloatOut,
		AnimationEmphasisPulse,
		AnimationEmphasisColorPulse,
		AnimationEmphasisTeeter,
		AnimationEmphasisSpin,
		AnimationEmphasisGrowShrink,
		AnimationEmphasisDesaturate,
		AnimationEmphasisDarken,
		AnimationEmphasisLighten,
		AnimationEmphasisTransparency,
		AnimationEmphasisObjectColor,
		AnimationPathLines,
		AnimationPathArcs,
		AnimationPathTurns,
		AnimationPathShapes,
		AnimationPathLoops,
		AnimationPathCustom:
		return nil
	default:
		return fmt.Errorf("unsupported animation effect %q", string(effect))
	}
}

func validateAnimationTrigger(trigger AnimationTrigger) error {
	switch trigger {
	case AnimationOnClick, AnimationWithPrevious, AnimationAfterPrevious:
		return nil
	default:
		return fmt.Errorf("unsupported animation trigger %q", string(trigger))
	}
}
