package shape

import (
	"encoding/xml"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func ParseShapeProperties(content []byte) (ParsedShapeProperties, error) {
	var s shapeXML
	// Freeform path segments are matched with `,any`, and reading back which
	// command each one is requires an untagged XMLName field.
	//nolint:musttag // Ordered path segments require an untagged XMLName.
	if err := xml.Unmarshal(content, &s); err != nil {
		return ParsedShapeProperties{}, err
	}
	ps := ParsedShapeProperties{PhIndex: -1}
	applyParsedShapeFill(&ps, &s)
	applyParsedShapeLine(&ps, &s)
	applyParsedShapeEffects(&ps, &s)
	applyParsedShapeIdentity(&ps, &s)
	applyParsedShapeGeometry(&ps, &s)
	applyParsedShapeFreeform(&ps, &s)
	applyParsedShapePictureFill(&ps, &s)
	applyParsedShapeTransform(&ps, &s)
	applyParsedShapeConnector(&ps, &s)
	applyParsedShapeTextFrame(&ps, &s)
	applyParsedShapeText(&ps, &s)
	return ps, nil
}
func applyParsedShapeGeometry(ps *ParsedShapeProperties, s *shapeXML) {
	if s.SpPr.PrstGeom == nil || s.SpPr.PrstGeom.Prst == "" {
		return
	}
	ps.Type = s.SpPr.PrstGeom.Prst
	if s.SpPr.PrstGeom.AvLst == nil || len(s.SpPr.PrstGeom.AvLst.Gd) == 0 {
		return
	}
	adjustments := make([]common.ShapeAdjustment, 0, len(s.SpPr.PrstGeom.AvLst.Gd))
	for _, gd := range s.SpPr.PrstGeom.AvLst.Gd {
		if gd.Name == "" || gd.Fmla == "" {
			continue
		}
		adjustments = append(adjustments, common.ShapeAdjustment{
			Name:    gd.Name,
			Formula: gd.Fmla,
		})
	}
	if len(adjustments) > 0 {
		ps.Adjustments = adjustments
	}
}
func applyParsedShapeFill(ps *ParsedShapeProperties, s *shapeXML) {
	if s.SpPr.NoFill != nil {
		background := true
		ps.Fill = &common.ShapeFill{Background: &background}
	}
	if s.SpPr.SolidFill != nil && s.SpPr.SolidFill.SrgbClr.Val != "" {
		fillColor := s.SpPr.SolidFill.SrgbClr.Val
		fill := &common.ShapeFill{Solid: &fillColor}
		if transparency, ok := parseSolidFillTransparency(s.SpPr.SolidFill); ok {
			fill.Transparency = &transparency
		}
		ps.Fill = fill
	}
	if s.SpPr.GradFill != nil {
		ps.Fill = &common.ShapeFill{Gradient: parseGradientFill(s.SpPr.GradFill)}
	}
	if s.SpPr.PattFill != nil {
		ps.Fill = &common.ShapeFill{Pattern: parsePatternFill(s.SpPr.PattFill)}
	}
}
func parseSolidFillTransparency(src *solidFillXML) (float64, bool) {
	if src == nil || src.SrgbClr.Alpha == nil || src.SrgbClr.Alpha.Val == nil {
		return 0, false
	}
	return alphaToTransparency(*src.SrgbClr.Alpha.Val), true
}

// alphaToTransparency converts an OOXML <a:alpha> value into the fraction of
// transparency the API reports, which is the inverse of opacity.
func alphaToTransparency(alpha int) float64 {
	alpha = max(alpha, 0)
	alpha = min(alpha, ooxmlPercentScale)
	return 1.0 - (float64(alpha) / float64(ooxmlPercentScale))
}
func parseGradientFill(src *gradientFillXML) *common.GradientFill {
	grad := &common.GradientFill{}
	if src.Lin != nil && src.Lin.Ang != nil {
		angle := float64(*src.Lin.Ang) / rotationDegreeToOOXML
		grad.AngleDeg = &angle
	}
	for _, gs := range src.GsLst.Gs {
		if gs.SrgbClr == nil || gs.SrgbClr.Val == "" {
			continue
		}
		stop := common.GradientStop{Color: gs.SrgbClr.Val}
		if gs.SrgbClr.Alpha != nil && gs.SrgbClr.Alpha.Val != nil {
			transparency := alphaToTransparency(*gs.SrgbClr.Alpha.Val)
			stop.Transparency = &transparency
		}
		if gs.Pos != nil {
			pos := float64(*gs.Pos) / gradientPositionScale
			stop.PositionPct = &pos
		}
		grad.Stops = append(grad.Stops, stop)
	}
	return grad
}
func parsePatternFill(src *patternFillXML) *common.PatternedFill {
	pattern := &common.PatternedFill{}
	if src.Prst != nil {
		pattern.Preset = src.Prst
	}
	if color, ok := parseColorRef(src.FgClr); ok {
		pattern.FgColor = &color
	}
	if color, ok := parseColorRef(src.BgClr); ok {
		pattern.BgColor = &color
	}
	return pattern
}
func applyParsedShapeLine(ps *ParsedShapeProperties, s *shapeXML) {
	if s.SpPr.Ln == nil {
		return
	}
	line := &common.ShapeLine{
		StartArrow:       nil,
		StartArrowWidth:  nil,
		StartArrowLength: nil,
		EndArrow:         nil,
		EndArrowWidth:    nil,
		EndArrowLength:   nil,
	}
	if color := parseLineColor(s); color != nil {
		lineColor := *color
		line.Color = &lineColor
	}
	if s.SpPr.Ln.W != nil {
		line.WidthEmu = s.SpPr.Ln.W
	}
	if s.SpPr.Ln.PrstDash != nil && s.SpPr.Ln.PrstDash.Val != "" {
		dash := s.SpPr.Ln.PrstDash.Val
		line.DashStyle = &dash
	}
	if s.SpPr.Ln.HeadEnd != nil {
		line.StartArrow, line.StartArrowWidth, line.StartArrowLength = parseLineEnd(
			s.SpPr.Ln.HeadEnd.Type, s.SpPr.Ln.HeadEnd.W, s.SpPr.Ln.HeadEnd.Len,
		)
	}
	if s.SpPr.Ln.TailEnd != nil {
		line.EndArrow, line.EndArrowWidth, line.EndArrowLength = parseLineEnd(
			s.SpPr.Ln.TailEnd.Type, s.SpPr.Ln.TailEnd.W, s.SpPr.Ln.TailEnd.Len,
		)
	}
	if hasLineStyle(line) {
		ps.Line = line
	}
}

func applyParsedShapeIdentity(ps *ParsedShapeProperties, s *shapeXML) {
	ps.Hidden = parseBoolAttr(s.NvSpPr.CNvPr.Hidden) || parseBoolAttr(s.NvPicPr.CNvPr.Hidden) ||
		parseBoolAttr(s.NvCxnSpPr.CNvPr.Hidden) || parseBoolAttr(s.NvGrpSpPr.CNvPr.Hidden) ||
		parseBoolAttr(s.NvGraphicFramePr.CNvPr.Hidden)
	switch {
	case s.NvSpPr.CNvPr.ID != 0:
		ps.ID = s.NvSpPr.CNvPr.ID
		ps.Name = s.NvSpPr.CNvPr.Name
		applyPlaceholderInfo(ps, s.NvSpPr.NvPr.Ph)
	case s.NvPicPr.CNvPr.ID != 0:
		ps.ID = s.NvPicPr.CNvPr.ID
		ps.Name = s.NvPicPr.CNvPr.Name
		applyPlaceholderInfo(ps, s.NvPicPr.NvPr.Ph)
	case s.NvCxnSpPr.CNvPr.ID != 0:
		ps.ID = s.NvCxnSpPr.CNvPr.ID
		ps.Name = s.NvCxnSpPr.CNvPr.Name
	case s.NvGrpSpPr.CNvPr.ID != 0:
		ps.ID = s.NvGrpSpPr.CNvPr.ID
		ps.Name = s.NvGrpSpPr.CNvPr.Name
	case s.NvGraphicFramePr.CNvPr.ID != 0:
		ps.ID = s.NvGraphicFramePr.CNvPr.ID
		ps.Name = s.NvGraphicFramePr.CNvPr.Name
	}
}
func applyPlaceholderInfo(ps *ParsedShapeProperties, ph *struct {
	Idx  *int   `xml:"idx,attr"`
	Type string `xml:"type,attr"`
}) {
	if ph == nil {
		return
	}
	ps.PhType = ph.Type
	if ph.Idx != nil {
		ps.PhIndex = *ph.Idx
		return
	}
	ps.PhIndex = 0
}
func applyParsedShapeTransform(ps *ParsedShapeProperties, s *shapeXML) {
	// Flip lives on whichever <a:xfrm> the shape carries, and is read for every
	// shape kind so a value written through ShapeUpdate can be read back.
	ps.FlipH = parseBoolAttr(s.SpPr.Xfrm.FlipH) || parseBoolAttr(s.Xfrm.FlipH) ||
		parseBoolAttr(s.GrpSpPr.Xfrm.FlipH)
	ps.FlipV = parseBoolAttr(s.SpPr.Xfrm.FlipV) || parseBoolAttr(s.Xfrm.FlipV) ||
		parseBoolAttr(s.GrpSpPr.Xfrm.FlipV)

	if s.SpPr.Xfrm.Ext.Cx != 0 || s.SpPr.Xfrm.Ext.Cy != 0 || s.SpPr.Xfrm.Off.X != 0 || s.SpPr.Xfrm.Off.Y != 0 {
		ps.X = s.SpPr.Xfrm.Off.X
		ps.Y = s.SpPr.Xfrm.Off.Y
		ps.W = s.SpPr.Xfrm.Ext.Cx
		ps.H = s.SpPr.Xfrm.Ext.Cy
		if s.SpPr.Xfrm.Rot != nil {
			rotation := float64(*s.SpPr.Xfrm.Rot) / rotationDegreeToOOXML
			ps.Rotation = &rotation
		}
		return
	}
	if s.Xfrm.Ext.Cx != 0 || s.Xfrm.Ext.Cy != 0 || s.Xfrm.Off.X != 0 || s.Xfrm.Off.Y != 0 {
		ps.X = s.Xfrm.Off.X
		ps.Y = s.Xfrm.Off.Y
		ps.W = s.Xfrm.Ext.Cx
		ps.H = s.Xfrm.Ext.Cy
		if s.Xfrm.Rot != nil {
			rotation := float64(*s.Xfrm.Rot) / rotationDegreeToOOXML
			ps.Rotation = &rotation
		}
		return
	}
	ps.X = s.GrpSpPr.Xfrm.Off.X
	ps.Y = s.GrpSpPr.Xfrm.Off.Y
	ps.W = s.GrpSpPr.Xfrm.Ext.Cx
	ps.H = s.GrpSpPr.Xfrm.Ext.Cy
	applyParsedGroupChildSpace(ps, s)
}

// applyParsedGroupChildSpace records a group's child coordinate space, which is
// what a child's raw x/y/cx/cy is expressed in. Reporting a child's own numbers
// as slide coordinates is the mismatch upstream issue #925 reports: the
// children look right relative to each other but wrong on the slide.
func applyParsedGroupChildSpace(ps *ParsedShapeProperties, s *shapeXML) {
	chOff := s.GrpSpPr.Xfrm.ChOff
	chExt := s.GrpSpPr.Xfrm.ChExt
	if chOff == nil && chExt == nil {
		return
	}
	space := &common.GroupChildSpace{}
	if chOff != nil {
		space.OffsetX = chOff.X
		space.OffsetY = chOff.Y
	}
	if chExt != nil {
		space.ExtentCx = chExt.Cx
		space.ExtentCy = chExt.Cy
	}
	ps.GroupChild = space
}
