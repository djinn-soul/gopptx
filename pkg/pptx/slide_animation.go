package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/animations"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

type (
	// AnimationEffect defines the type of animation effect.
	AnimationEffect = animations.AnimationEffect
	// AnimationTrigger defines when the animation starts.
	AnimationTrigger = animations.AnimationTrigger
	// AnimationDirection defines the direction or subtype of the animation.
	AnimationDirection = animations.AnimationDirection
	// Animation represents a single animation effect on a slide object.
	Animation = animations.Animation
	// AnimationDefinition is the interface for types that can be converted to an Animation.
	AnimationDefinition = animations.AnimationDefinition
)

const (
	AnimationEntranceAppear      = animations.AnimationEntranceAppear
	AnimationEntranceFade        = animations.AnimationEntranceFade
	AnimationEntranceFlyIn       = animations.AnimationEntranceFlyIn
	AnimationEntranceFloat       = animations.AnimationEntranceFloat
	AnimationEntranceSplit       = animations.AnimationEntranceSplit
	AnimationEntranceWipe        = animations.AnimationEntranceWipe
	AnimationEntranceShape       = animations.AnimationEntranceShape
	AnimationEntranceWheel       = animations.AnimationEntranceWheel
	AnimationEntranceRandomBars  = animations.AnimationEntranceRandomBars
	AnimationEntranceGrowAndTurn = animations.AnimationEntranceGrowAndTurn
	AnimationEntranceZoom        = animations.AnimationEntranceZoom
	AnimationEntranceSwivel      = animations.AnimationEntranceSwivel
	AnimationEntranceBounce      = animations.AnimationEntranceBounce

	AnimationExitDisappear = animations.AnimationExitDisappear
	AnimationExitFadeOut   = animations.AnimationExitFadeOut
	AnimationExitFlyOut    = animations.AnimationExitFlyOut
	AnimationExitFloatOut  = animations.AnimationExitFloatOut

	AnimationEmphasisPulse                        = animations.AnimationEmphasisPulse
	AnimationEmphasisColorPulse                   = animations.AnimationEmphasisColorPulse
	AnimationEmphasisTeeter                       = animations.AnimationEmphasisTeeter
	AnimationEmphasisSpin                         = animations.AnimationEmphasisSpin
	AnimationEmphasisGrowShrink                   = animations.AnimationEmphasisGrowShrink
	AnimationEmphasisDesaturate   AnimationEffect = animations.AnimationEmphasisDesaturate
	AnimationEmphasisDarken       AnimationEffect = animations.AnimationEmphasisDarken
	AnimationEmphasisLighten      AnimationEffect = animations.AnimationEmphasisLighten
	AnimationEmphasisTransparency AnimationEffect = animations.AnimationEmphasisTransparency
	AnimationEmphasisObjectColor  AnimationEffect = animations.AnimationEmphasisObjectColor

	AnimationPathLines  = animations.AnimationPathLines
	AnimationPathArcs   = animations.AnimationPathArcs
	AnimationPathTurns  = animations.AnimationPathTurns
	AnimationPathShapes = animations.AnimationPathShapes
	AnimationPathLoops  = animations.AnimationPathLoops
	AnimationPathCustom = animations.AnimationPathCustom
)

const (
	AnimationOnClick       = animations.AnimationOnClick
	AnimationWithPrevious  = animations.AnimationWithPrevious
	AnimationAfterPrevious = animations.AnimationAfterPrevious
)

const (
	AnimationDirIn        = animations.AnimationDirIn
	AnimationDirOut       = animations.AnimationDirOut
	AnimationDirUp        = animations.AnimationDirUp
	AnimationDirDown      = animations.AnimationDirDown
	AnimationDirLeft      = animations.AnimationDirLeft
	AnimationDirRight     = animations.AnimationDirRight
	AnimationDirUpLeft    = animations.AnimationDirUpLeft
	AnimationDirUpRight   = animations.AnimationDirUpRight
	AnimationDirDownLeft  = animations.AnimationDirDownLeft
	AnimationDirDownRight = animations.AnimationDirDownRight
)

func NewAnimation(shapeIndex int, effect AnimationEffect) Animation {
	return animations.NewAnimation(shapeIndex, effect)
}

func slideAnimationsXML(s elements.SlideContent, shapeIDs []int) string {
	return elements.SlideAnimationsXML(s, shapeIDs)
}
