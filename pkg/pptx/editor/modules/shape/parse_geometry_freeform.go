package shape

import (
	"strconv"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// applyParsedShapeFreeform reads <a:custGeom> into the shape read model. A
// preset geometry and a custom geometry are mutually exclusive in the schema,
// so a shape carrying custGeom leaves ps.Type as the element name.
func applyParsedShapeFreeform(ps *ParsedShapeProperties, s *shapeXML) {
	if s.SpPr.CustGeom == nil {
		return
	}
	paths := make([]common.FreeformPath, 0, len(s.SpPr.CustGeom.PathLst.Path))
	for _, p := range s.SpPr.CustGeom.PathLst.Path {
		paths = append(paths, common.FreeformPath{
			W:        p.W,
			H:        p.H,
			Fill:     p.Fill,
			Stroke:   parseOptionalXMLBool(p.Stroke),
			Segments: parseFreeformSegments(p.Segments),
		})
	}
	ps.Freeform = &common.FreeformGeometry{Paths: paths}
}

func parseFreeformSegments(src []pathSegmentXML) []common.FreeformSegment {
	segments := make([]common.FreeformSegment, 0, len(src))
	for _, seg := range src {
		kind, ok := freeformSegmentKind(seg.XMLName.Local)
		if !ok {
			continue
		}
		out := common.FreeformSegment{Type: kind}
		for _, pt := range seg.Pt {
			out.Points = append(out.Points, parseFreeformPoint(pt))
		}
		if kind == common.FreeformSegmentArcTo {
			applyFreeformArc(&out, seg)
		}
		segments = append(segments, out)
	}
	return segments
}

func freeformSegmentKind(local string) (string, bool) {
	switch local {
	case common.FreeformSegmentMoveTo,
		common.FreeformSegmentLineTo,
		common.FreeformSegmentCubicBezTo,
		common.FreeformSegmentQuadBezTo,
		common.FreeformSegmentArcTo,
		common.FreeformSegmentClose:
		return local, true
	default:
		return "", false
	}
}

func parseFreeformPoint(pt ptXML) common.FreeformPoint {
	out := common.FreeformPoint{}
	if value, err := strconv.Atoi(pt.X); err == nil {
		out.X = value
	} else if pt.X != "" {
		out.XExpr = pt.X
	}
	if value, err := strconv.Atoi(pt.Y); err == nil {
		out.Y = value
	} else if pt.Y != "" {
		out.YExpr = pt.Y
	}
	return out
}

func applyFreeformArc(out *common.FreeformSegment, seg pathSegmentXML) {
	if value, err := strconv.Atoi(seg.WR); err == nil {
		out.WidthRadius = &value
	}
	if value, err := strconv.Atoi(seg.HR); err == nil {
		out.HeightRadius = &value
	}
	if degrees, ok := parseOOXMLAngle(seg.StAng); ok {
		out.StartAngle = &degrees
	}
	if degrees, ok := parseOOXMLAngle(seg.SwAng); ok {
		out.SwingAngle = &degrees
	}
}

func parseOOXMLAngle(raw string) (float64, bool) {
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, false
	}
	return float64(value) / rotationDegreeToOOXML, true
}

func parseOptionalXMLBool(raw *string) *bool {
	if raw == nil {
		return nil
	}
	value := parseBoolAttr(raw)
	return &value
}

// applyParsedShapePictureFill reads an <a:blipFill> used as a shape fill. The
// relationship id is reported raw; resolving it to a part path needs the slide
// relationships and happens in the editor layer.
func applyParsedShapePictureFill(ps *ParsedShapeProperties, s *shapeXML) {
	src := s.SpPr.BlipFill
	if src == nil {
		return
	}
	picture := &common.PictureFill{}
	if src.Blip != nil {
		picture.RelID = src.Blip.Embed
		picture.External = src.Blip.Link
	}
	switch {
	case src.Stretch != nil:
		picture.Mode = "stretch"
	case src.Tile != nil:
		picture.Mode = "tile"
	}
	if src.SrcRect != nil {
		picture.Crop = &common.PictureFillCrop{
			Left:   srcRectFraction(src.SrcRect.L),
			Top:    srcRectFraction(src.SrcRect.T),
			Right:  srcRectFraction(src.SrcRect.R),
			Bottom: srcRectFraction(src.SrcRect.B),
		}
	}
	if ps.Fill == nil {
		ps.Fill = &common.ShapeFill{}
	}
	ps.Fill.Picture = picture
}

func srcRectFraction(value *int) float64 {
	if value == nil {
		return 0
	}
	return float64(*value) / srcRectScale
}
