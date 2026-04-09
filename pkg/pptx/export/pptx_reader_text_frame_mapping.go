package export

import (
	"strings"

	editorcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func editorTextFrameToExportTextFrame(frame *editorcommon.TextFrame) *shapes.TextFrame {
	if frame == nil {
		return nil
	}
	tf := shapes.NewTextFrame()
	has := applyExportTextFrameMargins(&tf, frame)
	has = applyExportTextFrameAnchor(&tf, frame) || has
	has = applyExportTextFrameWrap(&tf, frame) || has
	has = applyExportTextFrameAutoFit(&tf, frame) || has
	has = applyExportTextFrameLayout(&tf, frame) || has
	if !has {
		return nil
	}
	return &tf
}

func applyExportTextFrameMargins(tf *shapes.TextFrame, frame *editorcommon.TextFrame) bool {
	has := false
	if frame.MarginLeft != nil {
		tf.MarginLeft = styling.Emu(int64(*frame.MarginLeft))
		has = true
	}
	if frame.MarginRight != nil {
		tf.MarginRight = styling.Emu(int64(*frame.MarginRight))
		has = true
	}
	if frame.MarginTop != nil {
		tf.MarginTop = styling.Emu(int64(*frame.MarginTop))
		has = true
	}
	if frame.MarginBottom != nil {
		tf.MarginBottom = styling.Emu(int64(*frame.MarginBottom))
		has = true
	}
	return has
}

func applyExportTextFrameAnchor(tf *shapes.TextFrame, frame *editorcommon.TextFrame) bool {
	if frame.VerticalAlign == nil {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(*frame.VerticalAlign)) {
	case "t", "top":
		tf.Anchor = shapes.TextAnchorTop
	case "b", "bottom":
		tf.Anchor = shapes.TextAnchorBottom
	case "ctr", "center", "middle":
		tf.Anchor = shapes.TextAnchorMiddle
	default:
		return false
	}
	return true
}

func applyExportTextFrameWrap(tf *shapes.TextFrame, frame *editorcommon.TextFrame) bool {
	if frame.WordWrap == nil {
		return false
	}
	if *frame.WordWrap {
		tf.Wrap = shapes.TextWrapSquare
	} else {
		tf.Wrap = shapes.TextWrapNone
	}
	return true
}

func applyExportTextFrameAutoFit(tf *shapes.TextFrame, frame *editorcommon.TextFrame) bool {
	if frame.AutoFitType != nil {
		switch strings.ToLower(strings.TrimSpace(*frame.AutoFitType)) {
		case string(shapes.TextAutoFitNone):
			tf.AutoFit = shapes.TextAutoFitNone
		case "normal":
			tf.AutoFit = shapes.TextAutoFitNormal
		case "shape":
			tf.AutoFit = shapes.TextAutoFitShape
		default:
			return false
		}
		return true
	}
	if frame.AutoFit == nil {
		return false
	}
	if *frame.AutoFit {
		tf.AutoFit = shapes.TextAutoFitShape
	} else {
		tf.AutoFit = shapes.TextAutoFitNone
	}
	return true
}

func applyExportTextFrameLayout(tf *shapes.TextFrame, frame *editorcommon.TextFrame) bool {
	has := false
	if frame.Rotation != nil {
		tf.RotationDeg = frame.Rotation
		has = true
	}
	if frame.Orientation != nil && *frame.Orientation != "" {
		tf.Orientation = strings.TrimSpace(*frame.Orientation)
		has = true
	}
	if frame.Columns != nil {
		tf.Columns = *frame.Columns
		has = true
	}
	return has
}
