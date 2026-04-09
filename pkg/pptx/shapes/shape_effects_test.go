package shapes

import "testing"

func TestToXMLShapeSpec_MapsShapeEffects(t *testing.T) {
	shape := NewShape(ShapeTypeRectangle, 0, 0, 1, 1).
		WithShadow(true).
		WithGlowSpec(&ShapeGlow{Color: "AABBCC", RadiusEmu: 1200}).
		WithBlurSpec(&ShapeBlur{RadiusEmu: 1300}).
		WithSoftEdgeSpec(&ShapeSoftEdge{RadiusEmu: 1400}).
		WithReflectionSpec(&ShapeReflection{BlurEmu: 1500, DistanceEmu: 1600})

	spec := toXMLShapeSpec(shape, nil)
	if spec.Effects == nil {
		t.Fatal("expected effects in XML spec")
	}
	if !spec.Effects.Shadow || !spec.Effects.Glow || !spec.Effects.SoftEdges || !spec.Effects.Reflection {
		t.Fatalf("unexpected effects mapping: %+v", *spec.Effects)
	}
	if spec.Effects.GlowSpec == nil || spec.Effects.GlowSpec.Color != "AABBCC" || spec.Effects.GlowSpec.RadiusEmu != 1200 {
		t.Fatalf("unexpected glow mapping: %+v", spec.Effects.GlowSpec)
	}
	if spec.Effects.BlurSpec == nil || spec.Effects.BlurSpec.RadiusEmu != 1300 {
		t.Fatalf("unexpected blur mapping: %+v", spec.Effects.BlurSpec)
	}
	if spec.Effects.SoftEdgeSpec == nil || spec.Effects.SoftEdgeSpec.RadiusEmu != 1400 {
		t.Fatalf("unexpected soft edge mapping: %+v", spec.Effects.SoftEdgeSpec)
	}
	if spec.Effects.ReflectionSpec == nil || spec.Effects.ReflectionSpec.BlurEmu != 1500 ||
		spec.Effects.ReflectionSpec.DistanceEmu != 1600 {
		t.Fatalf("unexpected reflection mapping: %+v", spec.Effects.ReflectionSpec)
	}
}
