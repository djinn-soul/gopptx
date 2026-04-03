package shape

import (
	"encoding/xml"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func ParseShapeProperties(content []byte) (ParsedShapeProperties, error) {
	var s shapeXML
	if err := xml.Unmarshal(content, &s); err != nil {
		return ParsedShapeProperties{}, err
	}
	ps := ParsedShapeProperties{PhIndex: -1}
	applyParsedShapeFill(&ps, &s)
	applyParsedShapeLine(&ps, &s)
	applyParsedShapeEffects(&ps, &s)
	applyParsedShapeIdentity(&ps, &s)
	applyParsedShapeGeometry(&ps, &s)
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
	alpha := *src.SrgbClr.Alpha.Val
	alpha = max(alpha, 0)
	alpha = min(alpha, ooxmlPercentScale)
	return 1.0 - (float64(alpha) / float64(ooxmlPercentScale)), true
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

func applyParsedShapeEffects(ps *ParsedShapeProperties, s *shapeXML) {
	if s.SpPr.EffectLst == nil {
		return
	}
	applyParsedShadow(ps, s)
	if s.SpPr.EffectLst.Glow != nil {
		ps.Glow = parseGlowEffect(s)
	}
	if s.SpPr.EffectLst.Blur != nil {
		ps.Blur = parseBlurEffect(s)
	}
	if s.SpPr.EffectLst.SoftEdge != nil {
		ps.SoftEdge = parseSoftEdgeEffect(s)
	}
	if s.SpPr.EffectLst.Reflection != nil {
		ps.Reflection = parseReflectionEffect(s)
	}
}
func applyParsedShadow(ps *ParsedShapeProperties, s *shapeXML) {
	if s.SpPr.EffectLst.OuterShdw == nil &&
		s.SpPr.EffectLst.Glow == nil &&
		s.SpPr.EffectLst.Blur == nil &&
		s.SpPr.EffectLst.SoftEdge == nil &&
		s.SpPr.EffectLst.Reflection == nil {
		inherit := false
		ps.Shadow = &common.ShapeShadow{Inherit: &inherit}
		return
	}
	if s.SpPr.EffectLst.OuterShdw == nil {
		return
	}
	outer := s.SpPr.EffectLst.OuterShdw
	shadow := &common.ShapeShadow{}
	if outer.SrgbClr != nil && outer.SrgbClr.Val != "" {
		color := outer.SrgbClr.Val
		shadow.Color = &color
	}
	if outer.BlurRad != nil {
		shadow.BlurEmu = outer.BlurRad
	}
	if outer.Dist != nil {
		shadow.DistanceEmu = outer.Dist
	}
	if outer.Dir != nil {
		angle := float64(*outer.Dir) / rotationDegreeToOOXML
		shadow.AngleDeg = &angle
	}
	ps.Shadow = shadow
}
func parseGlowEffect(s *shapeXML) *common.ShapeGlow {
	glow := &common.ShapeGlow{}
	if s.SpPr.EffectLst.Glow.SrgbClr != nil && s.SpPr.EffectLst.Glow.SrgbClr.Val != "" {
		color := s.SpPr.EffectLst.Glow.SrgbClr.Val
		glow.Color = &color
	}
	if s.SpPr.EffectLst.Glow.Rad != nil {
		glow.RadiusEmu = s.SpPr.EffectLst.Glow.Rad
	}
	return glow
}
func parseBlurEffect(s *shapeXML) *common.ShapeBlur {
	blur := &common.ShapeBlur{}
	if s.SpPr.EffectLst.Blur.Rad != nil {
		blur.RadiusEmu = s.SpPr.EffectLst.Blur.Rad
	}
	return blur
}
func parseSoftEdgeEffect(s *shapeXML) *common.ShapeSoftEdge {
	softEdge := &common.ShapeSoftEdge{}
	if s.SpPr.EffectLst.SoftEdge.Rad != nil {
		softEdge.RadiusEmu = s.SpPr.EffectLst.SoftEdge.Rad
	}
	return softEdge
}
func parseReflectionEffect(s *shapeXML) *common.ShapeReflection {
	reflection := &common.ShapeReflection{}
	if s.SpPr.EffectLst.Reflection.BlurRad != nil {
		reflection.BlurEmu = s.SpPr.EffectLst.Reflection.BlurRad
	}
	if s.SpPr.EffectLst.Reflection.Dist != nil {
		reflection.DistanceEmu = s.SpPr.EffectLst.Reflection.Dist
	}
	return reflection
}
func applyParsedShapeIdentity(ps *ParsedShapeProperties, s *shapeXML) {
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
}
