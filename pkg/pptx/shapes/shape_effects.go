package shapes

// ShapeEffects configures visual effects for one shape.
type ShapeEffects struct {
	Shadow         bool
	Glow           bool
	SoftEdges      bool
	Reflection     bool
	GlowSpec       *ShapeGlow
	BlurSpec       *ShapeBlur
	SoftEdgeSpec   *ShapeSoftEdge
	ReflectionSpec *ShapeReflection
}

// WithShadow enables or disables shape outer shadow.
func (s Shape) WithShadow(enabled bool) Shape {
	if s.Effects == nil {
		s.Effects = &ShapeEffects{}
	}
	s.Effects.Shadow = enabled
	return s
}

// WithGlow enables or disables shape glow.
func (s Shape) WithGlow(enabled bool) Shape {
	if s.Effects == nil {
		s.Effects = &ShapeEffects{}
	}
	s.Effects.Glow = enabled
	return s
}

// WithSoftEdges enables or disables shape soft edges.
func (s Shape) WithSoftEdges(enabled bool) Shape {
	if s.Effects == nil {
		s.Effects = &ShapeEffects{}
	}
	s.Effects.SoftEdges = enabled
	return s
}

// WithReflection enables or disables shape reflection.
func (s Shape) WithReflection(enabled bool) Shape {
	if s.Effects == nil {
		s.Effects = &ShapeEffects{}
	}
	s.Effects.Reflection = enabled
	return s
}
