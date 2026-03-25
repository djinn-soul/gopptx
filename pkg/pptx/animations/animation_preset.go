package animations

import "strings"

func (a Animation) PresetID() uint32 {
	switch a.PresetClass() {
	case classEntr, classExit:
		return a.presetIDEntranceExit()
	case classEmph:
		return a.presetIDEmphasis()
	case classPath:
		return a.presetIDPath()
	}
	return 0
}

func (a Animation) presetIDEntranceExit() uint32 {
	switch a.Effect {
	case AnimationEntranceAppear, AnimationExitDisappear:
		return presetIDAppear
	case AnimationEntranceFade, AnimationExitFadeOut:
		return presetIDFade
	case AnimationEntranceFlyIn, AnimationExitFlyOut:
		return presetIDFly
	case AnimationEntranceFloat, AnimationExitFloatOut:
		return presetIDFloat
	case AnimationEntranceSplit:
		return presetIDSplit
	case AnimationEntranceWipe:
		return presetIDWipe
	case AnimationEntranceShape:
		return presetIDShape
	case AnimationEntranceWheel:
		return presetIDWheel
	case AnimationEntranceRandomBars:
		return presetIDRandomBars
	case AnimationEntranceGrowAndTurn:
		return presetIDGrowAndTurn
	case AnimationEntranceZoom:
		return presetIDZoom
	case AnimationEntranceSwivel:
		return presetIDSwivel
	case AnimationEntranceBounce:
		return presetIDBounce
	default:
		return 0
	}
}

func (a Animation) presetIDEmphasis() uint32 {
	switch a.Effect {
	case AnimationEmphasisPulse:
		return presetIDEmphasisPulse
	case AnimationEmphasisColorPulse:
		return presetIDColorPulse
	case AnimationEmphasisTeeter:
		return presetIDTeeter
	case AnimationEmphasisSpin:
		return presetIDSpin
	case AnimationEmphasisGrowShrink:
		return presetIDGrowShrink
	case AnimationEmphasisDesaturate:
		return presetIDDesaturate
	case AnimationEmphasisDarken:
		return presetIDDarken
	case AnimationEmphasisLighten:
		return presetIDLighten
	case AnimationEmphasisTransparency:
		return presetIDTransparency
	case AnimationEmphasisObjectColor:
		return presetIDObjectColor
	default:
		return 0
	}
}

func (a Animation) presetIDPath() uint32 {
	switch a.Effect {
	case AnimationPathLines:
		return presetIDPathLines
	case AnimationPathArcs:
		return presetIDPathArcs
	case AnimationPathTurns:
		return presetIDPathTurns
	case AnimationPathShapes:
		return presetIDPathShapes
	case AnimationPathLoops:
		return presetIDPathLoops
	case AnimationPathCustom:
		return presetIDPathCustom
	default:
		return 0
	}
}

func (a Animation) PresetClass() string {
	switch {
	case strings.HasPrefix(string(a.Effect), "entr_"):
		return classEntr
	case strings.HasPrefix(string(a.Effect), "exit_"):
		return classExit
	case strings.HasPrefix(string(a.Effect), "emph_"):
		return classEmph
	case strings.HasPrefix(string(a.Effect), "path_"):
		return classPath
	default:
		return classEntr
	}
}

func (a Animation) PresetSubtype() int {
	switch a.Effect {
	case AnimationEntranceFlyIn, AnimationExitFlyOut:
		switch a.Direction {
		case AnimationDirDown:
			return presetSubtypeFromTop
		case AnimationDirLeft:
			return presetSubtypeFromRight
		case AnimationDirRight:
			return presetSubtypeFromLeft
		case AnimationDirUp:
			return presetSubtypeFromBottom
		case AnimationDirDownLeft:
			return presetSubtypeFromTopRight
		case AnimationDirDownRight:
			return presetSubtypeFromTopLeft
		case AnimationDirUpLeft:
			return presetSubtypeFromBottomRight
		case AnimationDirUpRight:
			return presetSubtypeFromBottomLeft
		default:
			return presetSubtypeFromBottom
		}
	case AnimationEntranceWipe, AnimationEntranceFloat:
		switch a.Direction {
		case AnimationDirUp:
			return presetSubtypeFromLeft
		case AnimationDirDown:
			return presetSubtypeFromBottom
		case AnimationDirLeft:
			return presetSubtypeFromTop
		case AnimationDirRight:
			return presetSubtypeFromRight
		default:
			return presetSubtypeFromLeft
		}
	case AnimationEntranceSplit:
		switch a.Direction {
		case AnimationDirIn:
			return presetSubtypeFromTop
		case AnimationDirOut:
			return presetSubtypeFromRight
		default:
			return presetSubtypeFromRight
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
