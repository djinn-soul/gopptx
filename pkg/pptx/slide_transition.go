package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/transitions"
)

type (
	// SlideTransition is the extensibility contract for slide transitions.
	SlideTransition = transitions.SlideTransition
	// TransitionType is the built-in transition enum.
	TransitionType = transitions.TransitionType
	// TransitionDirection defines the direction of a transition.
	TransitionDirection = transitions.TransitionDirection
	// TransitionOrientation defines the orientation of a transition.
	TransitionOrientation = transitions.TransitionOrientation
	// TransitionOptions provides advanced configuration for a slide transition.
	TransitionOptions = transitions.TransitionOptions
)

const (
	TransitionNone       = transitions.TransitionNone
	TransitionCut        = transitions.TransitionCut
	TransitionFade       = transitions.TransitionFade
	TransitionPush       = transitions.TransitionPush
	TransitionWipe       = transitions.TransitionWipe
	TransitionSplit      = transitions.TransitionSplit
	TransitionReveal     = transitions.TransitionReveal
	TransitionCover      = transitions.TransitionCover
	TransitionZoom       = transitions.TransitionZoom
	TransitionRandomBars = transitions.TransitionRandomBars
	TransitionShape      = transitions.TransitionShape
	TransitionUncover    = transitions.TransitionUncover
	TransitionFlash      = transitions.TransitionFlash
	TransitionStrips     = transitions.TransitionStrips
	TransitionBlinds     = transitions.TransitionBlinds
	TransitionClock      = transitions.TransitionClock
	TransitionRipple     = transitions.TransitionRipple
	TransitionHoneycomb  = transitions.TransitionHoneycomb
	TransitionGlitter    = transitions.TransitionGlitter
	TransitionVortex     = transitions.TransitionVortex
	TransitionShred      = transitions.TransitionShred
	TransitionSwitch     = transitions.TransitionSwitch
	TransitionFlip       = transitions.TransitionFlip
	TransitionGallery    = transitions.TransitionGallery
	TransitionCube       = transitions.TransitionCube
	TransitionDoors      = transitions.TransitionDoors
	TransitionBox        = transitions.TransitionBox
	TransitionRandom     = transitions.TransitionRandom
)

const (
	TransitionDirIn        = transitions.TransitionDirIn
	TransitionDirOut       = transitions.TransitionDirOut
	TransitionDirUp        = transitions.TransitionDirUp
	TransitionDirDown      = transitions.TransitionDirDown
	TransitionDirLeft      = transitions.TransitionDirLeft
	TransitionDirRight     = transitions.TransitionDirRight
	TransitionDirUpLeft    = transitions.TransitionDirUpLeft
	TransitionDirUpRight   = transitions.TransitionDirUpRight
	TransitionDirDownLeft  = transitions.TransitionDirDownLeft
	TransitionDirDownRight = transitions.TransitionDirDownRight
)

const (
	TransitionOrientHorizontal = transitions.TransitionOrientHorizontal
	TransitionOrientVertical   = transitions.TransitionOrientVertical
)

func ValidateSlideTransition(s elements.SlideContent, index int) error {
	return elements.ValidateSlideTransition(s, index)
}
