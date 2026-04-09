package shapes

// ShapeGlow defines detailed glow settings for a shape.
type ShapeGlow struct {
	Color     string
	RadiusEmu int
}

// ShapeBlur defines detailed blur settings for a shape.
type ShapeBlur struct {
	RadiusEmu int
}

// ShapeSoftEdge defines detailed soft-edge settings for a shape.
type ShapeSoftEdge struct {
	RadiusEmu int
}

// ShapeReflection defines detailed reflection settings for a shape.
type ShapeReflection struct {
	BlurEmu     int
	DistanceEmu int
}

// WithGlowSpec sets detailed glow settings and enables the glow flag.
func (s Shape) WithGlowSpec(glow *ShapeGlow) Shape {
	if s.Effects == nil {
		s.Effects = &ShapeEffects{}
	}
	s.Effects.Glow = glow != nil
	s.Effects.GlowSpec = glow
	return s
}

// WithBlurSpec sets detailed blur settings.
func (s Shape) WithBlurSpec(blur *ShapeBlur) Shape {
	if s.Effects == nil {
		s.Effects = &ShapeEffects{}
	}
	s.Effects.BlurSpec = blur
	return s
}

// WithSoftEdgeSpec sets detailed soft-edge settings and enables the soft-edge flag.
func (s Shape) WithSoftEdgeSpec(softEdge *ShapeSoftEdge) Shape {
	if s.Effects == nil {
		s.Effects = &ShapeEffects{}
	}
	s.Effects.SoftEdges = softEdge != nil
	s.Effects.SoftEdgeSpec = softEdge
	return s
}

// WithReflectionSpec sets detailed reflection settings and enables the reflection flag.
func (s Shape) WithReflectionSpec(reflection *ShapeReflection) Shape {
	if s.Effects == nil {
		s.Effects = &ShapeEffects{}
	}
	s.Effects.Reflection = reflection != nil
	s.Effects.ReflectionSpec = reflection
	return s
}
