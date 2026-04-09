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

	if hasTextFrameProps(frame) {
		ps.TextFrame = frame
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
		frame.Rotation != nil
}
