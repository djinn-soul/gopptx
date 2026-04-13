package transitions

import (
	"errors"
	"fmt"
	"strings"
)

func (o TransitionOptions) Validate() error {
	if err := o.Type.Validate(); err != nil {
		return err
	}

	if o.Sound != nil && strings.TrimSpace(o.Sound.RelID) == "" {
		return errors.New("transition sound requires a valid relationship ID")
	}

	if err := o.validateOrientation(); err != nil {
		return err
	}

	if o.SpokeCount > 0 && o.Type != TransitionClock {
		return fmt.Errorf("transition %q does not support spoke count", o.Type)
	}

	return o.validateDirection()
}

func (o TransitionOptions) validateOrientation() error {
	if o.Orientation == "" {
		return nil
	}
	switch o.Type {
	case TransitionSplit, TransitionBlinds, TransitionRandomBars:
		if o.Orientation != TransitionOrientHorizontal && o.Orientation != TransitionOrientVertical {
			return fmt.Errorf("invalid orientation %q for transition %q", o.Orientation, o.Type)
		}
		return nil
	default:
		return fmt.Errorf("transition %q does not support orientation", o.Type)
	}
}

func (o TransitionOptions) validateDirection() error {
	if o.Direction == "" {
		return nil
	}
	switch o.Type {
	case TransitionPush, TransitionWipe, TransitionReveal, TransitionCover:
		return validateSimpleDirection(o.Type, o.Direction)
	case TransitionZoom, TransitionSplit:
		return validateInOutDirection(o.Type, o.Direction)
	case TransitionUncover:
		return validateUncoverDirection(o.Direction)
	case TransitionStrips:
		return validateStripsDirection(o.Direction)
	case TransitionMorph:
		return o.validateMorph()
	default:
		return fmt.Errorf("transition %q does not support direction", o.Type)
	}
}

func (o TransitionOptions) validateMorph() error {
	if o.Direction != "" {
		return errors.New("morph transition does not support direction")
	}
	if o.Orientation != "" {
		return errors.New("morph transition does not support orientation")
	}
	switch o.MorphOption {
	case "", MorphOptionObject, MorphOptionWord, MorphOptionCharacter:
		return nil
	default:
		return fmt.Errorf("invalid morph option %q", o.MorphOption)
	}
}

func validateSimpleDirection(t TransitionType, d TransitionDirection) error {
	if d == TransitionDirUp || d == TransitionDirDown || d == TransitionDirLeft || d == TransitionDirRight {
		return nil
	}
	return fmt.Errorf("invalid direction %q for transition %q (expected u|d|l|r)", d, t)
}

func validateInOutDirection(t TransitionType, d TransitionDirection) error {
	if d == TransitionDirIn || d == TransitionDirOut {
		return nil
	}
	return fmt.Errorf("invalid direction %q for transition %q (expected in|out)", d, t)
}

func validateUncoverDirection(d TransitionDirection) error {
	switch d {
	case TransitionDirUp, TransitionDirDown, TransitionDirLeft, TransitionDirRight,
		TransitionDirUpLeft, TransitionDirUpRight, TransitionDirDownLeft, TransitionDirDownRight:
		return nil
	default:
		return fmt.Errorf("invalid direction %q for transition %q (expected u|d|l|r|lu|ru|ld|rd)", d, TransitionUncover)
	}
}

func validateStripsDirection(d TransitionDirection) error {
	if d == TransitionDirUpLeft || d == TransitionDirUpRight ||
		d == TransitionDirDownLeft || d == TransitionDirDownRight {
		return nil
	}
	return fmt.Errorf("invalid direction %q for transition %q (expected ul|ur|dl|dr)", d, TransitionStrips)
}

func (t TransitionType) Validate() error {
	switch t {
	case TransitionNone, TransitionCut, TransitionFade, TransitionPush, TransitionWipe,
		TransitionSplit, TransitionReveal, TransitionCover, TransitionZoom,
		TransitionRandomBars, TransitionShape, TransitionUncover, TransitionFlash,
		TransitionStrips, TransitionBlinds, TransitionClock, TransitionRipple,
		TransitionHoneycomb, TransitionGlitter, TransitionVortex, TransitionShred,
		TransitionSwitch, TransitionFlip, TransitionGallery, TransitionCube,
		TransitionDoors, TransitionBox, TransitionRandom, TransitionMorph:
		return nil
	default:
		return fmt.Errorf("unsupported transition type: %q", string(t))
	}
}
