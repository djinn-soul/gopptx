package shape

import (
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

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
	if s.SpPr.EffectLst.OuterShdw == nil {
		inherit := false
		ps.Shadow = &common.ShapeShadow{Inherit: &inherit}
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
