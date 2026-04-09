package editor

import (
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const (
	alignCenter            = "center"
	alignCtr               = "ctr"
	alignJust              = "just"
	alignDist              = "dist"
	bulletStyleNone        = "none"
	bulletStyleBullet      = "bullet"
	bulletStyleCustom      = "custom"
	bulletStyleNumber      = "number"
	bulletStyleLetterLower = "letter_lower"
	bulletStyleLetterUpper = "letter_upper"
	bulletStyleRomanLower  = "roman_lower"
	bulletStyleRomanUpper  = "roman_upper"
	paragraphPtUnit        = 100
	paragraphPct           = 1000
	defaultBodyPrIn        = 457200
	textFrameWrapSquare    = "square"
	textFrameWrapNone      = "none"
	textFrameAutoFitShape  = "spAutoFit"
	textFrameAutoFitNormal = "normAutoFit"
	textFrameAutoFitNone   = "none"
)

func editorTextFrameToSpec(tf *common.TextFrame) (pptxxml.TextFrameSpec, error) {
	spec := pptxxml.TextFrameSpec{
		MarginLeft:   defaultBodyPrIn,
		MarginRight:  defaultBodyPrIn,
		MarginTop:    defaultBodyPrIn,
		MarginBottom: defaultBodyPrIn,
		Wrap:         textFrameWrapSquare,
		Anchor:       alignCtr,
		AutoFit:      textFrameAutoFitShape,
	}
	if tf == nil {
		return spec, nil
	}
	applyTextFrameMargins(&spec, tf)
	if err := applyTextFrameWrapAnchorAndOrientation(&spec, tf); err != nil {
		return pptxxml.TextFrameSpec{}, err
	}
	if err := applyTextFrameColumnsAndRotation(&spec, tf); err != nil {
		return pptxxml.TextFrameSpec{}, err
	}
	applyTextFrameAutoFit(&spec, tf)
	return spec, nil
}

func applyTextFrameMargins(spec *pptxxml.TextFrameSpec, tf *common.TextFrame) {
	if tf.MarginTop != nil {
		spec.MarginTop = int64(*tf.MarginTop)
	}
	if tf.MarginBottom != nil {
		spec.MarginBottom = int64(*tf.MarginBottom)
	}
	if tf.MarginLeft != nil {
		spec.MarginLeft = int64(*tf.MarginLeft)
	}
	if tf.MarginRight != nil {
		spec.MarginRight = int64(*tf.MarginRight)
	}
}

func applyTextFrameWrapAnchorAndOrientation(spec *pptxxml.TextFrameSpec, tf *common.TextFrame) error {
	if tf.WordWrap != nil {
		if *tf.WordWrap {
			spec.Wrap = textFrameWrapSquare
		} else {
			spec.Wrap = textFrameWrapNone
		}
	}
	if tf.VerticalAlign != nil && *tf.VerticalAlign != "" {
		anchor, err := normalizeTextFrameVerticalAlign(*tf.VerticalAlign)
		if err != nil {
			return err
		}
		spec.Anchor = anchor
	}
	if tf.Orientation == nil || *tf.Orientation == "" {
		return nil
	}
	orientation, err := normalizeTextFrameOrientation(*tf.Orientation)
	if err != nil {
		return err
	}
	spec.Orientation = orientation
	return nil
}

func applyTextFrameColumnsAndRotation(spec *pptxxml.TextFrameSpec, tf *common.TextFrame) error {
	if tf.Columns != nil {
		if *tf.Columns < minTextFrameColumns {
			return fmt.Errorf("text_frame.columns must be >= %d", minTextFrameColumns)
		}
		spec.NumCol = *tf.Columns
	}
	if tf.Rotation == nil {
		return nil
	}
	rotation, err := normalizeTextFrameRotation(*tf.Rotation)
	if err != nil {
		return err
	}
	spec.Rotation = &rotation
	return nil
}

func applyTextFrameAutoFit(spec *pptxxml.TextFrameSpec, tf *common.TextFrame) {
	if tf.AutoFitType != nil {
		switch strings.ToLower(strings.TrimSpace(*tf.AutoFitType)) {
		case "normal":
			spec.AutoFit = textFrameAutoFitNormal
		case "shape":
			spec.AutoFit = textFrameAutoFitShape
		case bulletStyleNone:
			spec.AutoFit = textFrameAutoFitNone
		default:
			spec.AutoFit = textFrameAutoFitShape
		}
		return
	}
	if tf.AutoFit == nil {
		return
	}
	if *tf.AutoFit {
		spec.AutoFit = textFrameAutoFitShape
	} else {
		spec.AutoFit = textFrameAutoFitNone
	}
}
