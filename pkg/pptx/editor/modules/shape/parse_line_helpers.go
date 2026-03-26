package shape

import common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"

func parseLineColor(s *shapeXML) *string {
	if s.SpPr.Ln.SolidFill == nil || s.SpPr.Ln.SolidFill.SrgbClr.Val == "" {
		return nil
	}
	color := s.SpPr.Ln.SolidFill.SrgbClr.Val
	return &color
}

func parseLineEnd(rawType, rawWidth, rawLength *string) (*string, *string, *string) {
	var arrow, width, length *string
	if rawType != nil && *rawType != "" {
		v := *rawType
		arrow = &v
	}
	if rawWidth != nil && *rawWidth != "" {
		v := *rawWidth
		width = &v
	}
	if rawLength != nil && *rawLength != "" {
		v := *rawLength
		length = &v
	}
	return arrow, width, length
}

func hasLineStyle(line *common.ShapeLine) bool {
	return line.Color != nil || line.WidthEmu != nil || line.DashStyle != nil ||
		line.StartArrow != nil || line.StartArrowWidth != nil || line.StartArrowLength != nil ||
		line.EndArrow != nil || line.EndArrowWidth != nil || line.EndArrowLength != nil
}
