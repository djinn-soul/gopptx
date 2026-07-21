package shape

import (
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func applyParsedShapeConnector(ps *ParsedShapeProperties, s *shapeXML) {
	if !shapes.IsConnectorType(ps.Type) && s.NvCxnSpPr.CNvPr.ID == 0 {
		return
	}
	info := &common.ConnectorInfo{
		FlipH: parseBoolAttr(s.SpPr.Xfrm.FlipH) || parseBoolAttr(s.Xfrm.FlipH),
		FlipV: parseBoolAttr(s.SpPr.Xfrm.FlipV) || parseBoolAttr(s.Xfrm.FlipV),
	}
	if st := s.NvCxnSpPr.CNvCxnSpPr.StCxn; st != nil {
		info.StartShapeID = st.ID
		info.StartSiteIndex = st.Idx
	}
	if end := s.NvCxnSpPr.CNvCxnSpPr.EndCxn; end != nil {
		info.EndShapeID = end.ID
		info.EndSiteIndex = end.Idx
	}
	ps.Connector = info
}

func parseBoolAttr(value *string) bool {
	if value == nil {
		return false
	}
	switch *value {
	case "1", boolTrueLiteral, "TRUE", "True":
		return true
	default:
		return false
	}
}
