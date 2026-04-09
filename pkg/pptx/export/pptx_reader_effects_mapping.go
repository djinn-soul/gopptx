package export

import (
	editorcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func editorEffectsToExportEffects(es editorcommon.Shape) *shapes.ShapeEffects {
	effects := &shapes.ShapeEffects{
		Shadow:         editorShapeHasShadow(es.Shadow),
		Glow:           es.Glow != nil,
		GlowSpec:       editorGlowToExportGlow(es.Glow),
		BlurSpec:       editorBlurToExportBlur(es.Blur),
		SoftEdges:      es.SoftEdge != nil,
		SoftEdgeSpec:   editorSoftEdgeToExportSoftEdge(es.SoftEdge),
		Reflection:     es.Reflection != nil,
		ReflectionSpec: editorReflectionToExportReflection(es.Reflection),
	}
	if !effects.Shadow && !effects.Glow && !effects.SoftEdges && !effects.Reflection &&
		effects.GlowSpec == nil && effects.BlurSpec == nil && effects.SoftEdgeSpec == nil &&
		effects.ReflectionSpec == nil {
		return nil
	}
	return effects
}

func editorShapeHasShadow(shadow *editorcommon.ShapeShadow) bool {
	return shadow != nil &&
		(shadow.Inherit == nil || *shadow.Inherit || shadow.Color != nil ||
			shadow.BlurEmu != nil || shadow.DistanceEmu != nil || shadow.AngleDeg != nil)
}

func editorGlowToExportGlow(glow *editorcommon.ShapeGlow) *shapes.ShapeGlow {
	if glow == nil {
		return nil
	}
	exported := &shapes.ShapeGlow{}
	if glow.Color != nil {
		exported.Color = *glow.Color
	}
	if glow.RadiusEmu != nil {
		exported.RadiusEmu = *glow.RadiusEmu
	}
	return exported
}

func editorBlurToExportBlur(blur *editorcommon.ShapeBlur) *shapes.ShapeBlur {
	if blur == nil {
		return nil
	}
	exported := &shapes.ShapeBlur{}
	if blur.RadiusEmu != nil {
		exported.RadiusEmu = *blur.RadiusEmu
	}
	return exported
}

func editorSoftEdgeToExportSoftEdge(softEdge *editorcommon.ShapeSoftEdge) *shapes.ShapeSoftEdge {
	if softEdge == nil {
		return nil
	}
	exported := &shapes.ShapeSoftEdge{}
	if softEdge.RadiusEmu != nil {
		exported.RadiusEmu = *softEdge.RadiusEmu
	}
	return exported
}

func editorReflectionToExportReflection(reflection *editorcommon.ShapeReflection) *shapes.ShapeReflection {
	if reflection == nil {
		return nil
	}
	exported := &shapes.ShapeReflection{}
	if reflection.BlurEmu != nil {
		exported.BlurEmu = *reflection.BlurEmu
	}
	if reflection.DistanceEmu != nil {
		exported.DistanceEmu = *reflection.DistanceEmu
	}
	return exported
}
