package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

type (
	// SlideTransition is the extensibility contract for slide transitions.
	SlideTransition = elements.SlideTransition
	// TransitionType is the built-in transition enum.
	TransitionType = elements.TransitionType
	// TransitionDirection defines the direction of a transition.
	TransitionDirection = elements.TransitionDirection
	// TransitionOrientation defines the orientation of a transition.
	TransitionOrientation = elements.TransitionOrientation
	// TransitionOptions provides advanced configuration for a slide transition.
	TransitionOptions = elements.TransitionOptions
)

const (
	TransitionNone       = elements.TransitionNone
	TransitionCut        = elements.TransitionCut
	TransitionFade       = elements.TransitionFade
	TransitionPush       = elements.TransitionPush
	TransitionWipe       = elements.TransitionWipe
	TransitionSplit      = elements.TransitionSplit
	TransitionReveal     = elements.TransitionReveal
	TransitionCover      = elements.TransitionCover
	TransitionZoom       = elements.TransitionZoom
	TransitionRandomBars = elements.TransitionRandomBars
	TransitionShape      = elements.TransitionShape
	TransitionUncover    = elements.TransitionUncover
	TransitionFlash      = elements.TransitionFlash
	TransitionStrips     = elements.TransitionStrips
	TransitionBlinds     = elements.TransitionBlinds
	TransitionClock      = elements.TransitionClock
	TransitionRipple     = elements.TransitionRipple
	TransitionHoneycomb  = elements.TransitionHoneycomb
	TransitionGlitter    = elements.TransitionGlitter
	TransitionVortex     = elements.TransitionVortex
	TransitionShred      = elements.TransitionShred
	TransitionSwitch     = elements.TransitionSwitch
	TransitionFlip       = elements.TransitionFlip
	TransitionGallery    = elements.TransitionGallery
	TransitionCube       = elements.TransitionCube
	TransitionDoors      = elements.TransitionDoors
	TransitionBox        = elements.TransitionBox
	TransitionRandom     = elements.TransitionRandom
)

const (
	TransitionDirIn        = elements.TransitionDirIn
	TransitionDirOut       = elements.TransitionDirOut
	TransitionDirUp        = elements.TransitionDirUp
	TransitionDirDown      = elements.TransitionDirDown
	TransitionDirLeft      = elements.TransitionDirLeft
	TransitionDirRight     = elements.TransitionDirRight
	TransitionDirUpLeft    = elements.TransitionDirUpLeft
	TransitionDirUpRight   = elements.TransitionDirUpRight
	TransitionDirDownLeft  = elements.TransitionDirDownLeft
	TransitionDirDownRight = elements.TransitionDirDownRight
)

const (
	TransitionOrientHorizontal = elements.TransitionOrientHorizontal
	TransitionOrientVertical   = elements.TransitionOrientVertical
)
