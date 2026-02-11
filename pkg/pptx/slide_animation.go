package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

type (
	// AnimationEffect defines the type of animation effect.
	AnimationEffect = elements.AnimationEffect
	// AnimationTrigger defines when the animation starts.
	AnimationTrigger = elements.AnimationTrigger
	// AnimationDirection defines the direction or subtype of the animation.
	AnimationDirection = elements.AnimationDirection
	// Animation represents a single animation effect on a slide object.
	Animation = elements.Animation
	// AnimationDefinition is the interface for types that can be converted to an Animation.
	AnimationDefinition = elements.AnimationDefinition
)

const (
	AnimationEntranceAppear      = elements.AnimationEntranceAppear
	AnimationEntranceFade        = elements.AnimationEntranceFade
	AnimationEntranceFlyIn       = elements.AnimationEntranceFlyIn
	AnimationEntranceFloat       = elements.AnimationEntranceFloat
	AnimationEntranceSplit       = elements.AnimationEntranceSplit
	AnimationEntranceWipe        = elements.AnimationEntranceWipe
	AnimationEntranceShape       = elements.AnimationEntranceShape
	AnimationEntranceWheel       = elements.AnimationEntranceWheel
	AnimationEntranceRandomBars  = elements.AnimationEntranceRandomBars
	AnimationEntranceGrowAndTurn = elements.AnimationEntranceGrowAndTurn
	AnimationEntranceZoom        = elements.AnimationEntranceZoom
	AnimationEntranceSwivel      = elements.AnimationEntranceSwivel
	AnimationEntranceBounce      = elements.AnimationEntranceBounce

	AnimationExitDisappear = elements.AnimationExitDisappear
	AnimationExitFadeOut   = elements.AnimationExitFadeOut
	AnimationExitFlyOut    = elements.AnimationExitFlyOut
	AnimationExitFloatOut  = elements.AnimationExitFloatOut

	AnimationEmphasisPulse                        = elements.AnimationEmphasisPulse
	AnimationEmphasisColorPulse                   = elements.AnimationEmphasisColorPulse
	AnimationEmphasisTeeter                       = elements.AnimationEmphasisTeeter
	AnimationEmphasisSpin                         = elements.AnimationEmphasisSpin
	AnimationEmphasisGrowShrink                   = elements.AnimationEmphasisGrowShrink
	AnimationEmphasisDesaturate   AnimationEffect = elements.AnimationEmphasisDesaturate
	AnimationEmphasisDarken       AnimationEffect = elements.AnimationEmphasisDarken
	AnimationEmphasisLighten      AnimationEffect = elements.AnimationEmphasisLighten
	AnimationEmphasisTransparency AnimationEffect = elements.AnimationEmphasisTransparency
	AnimationEmphasisObjectColor  AnimationEffect = elements.AnimationEmphasisObjectColor

	AnimationPathLines  = elements.AnimationPathLines
	AnimationPathArcs   = elements.AnimationPathArcs
	AnimationPathTurns  = elements.AnimationPathTurns
	AnimationPathShapes = elements.AnimationPathShapes
	AnimationPathLoops  = elements.AnimationPathLoops
	AnimationPathCustom = elements.AnimationPathCustom
)

const (
	AnimationOnClick       = elements.AnimationOnClick
	AnimationWithPrevious  = elements.AnimationWithPrevious
	AnimationAfterPrevious = elements.AnimationAfterPrevious
)

const (
	AnimationDirIn        = elements.AnimationDirIn
	AnimationDirOut       = elements.AnimationDirOut
	AnimationDirUp        = elements.AnimationDirUp
	AnimationDirDown      = elements.AnimationDirDown
	AnimationDirLeft      = elements.AnimationDirLeft
	AnimationDirRight     = elements.AnimationDirRight
	AnimationDirUpLeft    = elements.AnimationDirUpLeft
	AnimationDirUpRight   = elements.AnimationDirUpRight
	AnimationDirDownLeft  = elements.AnimationDirDownLeft
	AnimationDirDownRight = elements.AnimationDirDownRight
)

func NewAnimation(shapeIndex int, effect AnimationEffect) Animation {
	return elements.NewAnimation(shapeIndex, effect)
}

func slideAnimationsXML(s SlideContent, shapeIDs []int) string {
	return elements.SlideAnimationsXML(s, shapeIDs)
}
