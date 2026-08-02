package shape

import common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"

const (
	bodyPrAutoFitNormal = "normal"
	bodyPrAutoFitShape  = "shape"
	bodyPrAutoFitNone   = "none"
)

func applyParsedShapeTextFrame(ps *ParsedShapeProperties, s *shapeXML) {
	if s.TxBody.BodyPr == nil {
		return
	}
	bodyPr := s.TxBody.BodyPr
	frame := &common.TextFrame{}

	if bodyPr.TIns != nil {
		frame.MarginTop = bodyPr.TIns
	}
	if bodyPr.BIns != nil {
		frame.MarginBottom = bodyPr.BIns
	}
	if bodyPr.LeftInset != nil {
		frame.MarginLeft = bodyPr.LeftInset
	}
	if bodyPr.RIns != nil {
		frame.MarginRight = bodyPr.RIns
	}
	if bodyPr.Wrap != nil {
		wrap := *bodyPr.Wrap != "none"
		frame.WordWrap = &wrap
	}
	if bodyPr.Anchor != nil && *bodyPr.Anchor != "" {
		frame.VerticalAlign = bodyPr.Anchor
	}
	if bodyPr.Vert != nil && *bodyPr.Vert != "" {
		frame.Orientation = bodyPr.Vert
	}
	if bodyPr.NumCol != nil {
		frame.Columns = bodyPr.NumCol
	}
	if bodyPr.Rot != nil {
		rotation := float64(*bodyPr.Rot) / rotationDegreeToOOXML
		frame.Rotation = &rotation
	}

	autoFitType, autoFitBool := parseBodyPrAutoFit(bodyPr)
	if autoFitType != nil {
		frame.AutoFitType = autoFitType
	}
	if autoFitBool != nil {
		frame.AutoFit = autoFitBool
	}
	applyParsedShrinkAmounts(frame, bodyPr)

	if hasTextFrameProps(frame) {
		ps.TextFrame = frame
	}
}

// autofitPercentScale converts OOXML's thousandths of a percent to percent.
const autofitPercentScale = 1000.0

// applyParsedShrinkAmounts reports the fontScale and lnSpcReduction of a
// normAutofit frame in percent, so a caller can read back what PowerPoint (or
// an earlier edit) computed instead of guessing.
func applyParsedShrinkAmounts(frame *common.TextFrame, bodyPr *bodyPrXML) {
	norm := bodyPr.NormAutoFit
	if norm == nil {
		norm = bodyPr.NormAutofit
	}
	if norm == nil {
		return
	}
	if norm.FontScale != nil {
		scale := float64(*norm.FontScale) / autofitPercentScale
		frame.FontScale = &scale
	}
	if norm.LnSpcReduction != nil {
		reduction := float64(*norm.LnSpcReduction) / autofitPercentScale
		frame.LineSpaceReduction = &reduction
	}
}

func parseBodyPrAutoFit(bodyPr *bodyPrXML) (*string, *bool) {
	switch {
	case bodyPr.NormAutoFit != nil || bodyPr.NormAutofit != nil:
		autoFitType := bodyPrAutoFitNormal
		autoFit := true
		return &autoFitType, &autoFit
	case bodyPr.SpAutoFit != nil:
		autoFitType := bodyPrAutoFitShape
		autoFit := true
		return &autoFitType, &autoFit
	case bodyPr.NoAutofit != nil:
		autoFitType := bodyPrAutoFitNone
		autoFit := false
		return &autoFitType, &autoFit
	default:
		return nil, nil
	}
}

func hasTextFrameProps(frame *common.TextFrame) bool {
	return frame.MarginTop != nil ||
		frame.MarginBottom != nil ||
		frame.MarginLeft != nil ||
		frame.MarginRight != nil ||
		frame.WordWrap != nil ||
		frame.AutoFit != nil ||
		frame.AutoFitType != nil ||
		frame.VerticalAlign != nil ||
		frame.Orientation != nil ||
		frame.Columns != nil ||
		frame.Rotation != nil ||
		frame.FontScale != nil ||
		frame.LineSpaceReduction != nil
}
