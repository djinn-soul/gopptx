package shapes

import "testing"

func TestToXMLShapeSpec_MapsShapeEffects(t *testing.T) {
	shape := NewShape(ShapeTypeRectangle, 0, 0, 1, 1).
		WithShadow(true).
		WithGlow(true).
		WithSoftEdges(true).
		WithReflection(true)

	spec := toXMLShapeSpec(shape, nil)
	if spec.Effects == nil {
		t.Fatal("expected effects in XML spec")
	}
	if !spec.Effects.Shadow || !spec.Effects.Glow || !spec.Effects.SoftEdges || !spec.Effects.Reflection {
		t.Fatalf("unexpected effects mapping: %+v", *spec.Effects)
	}
}
