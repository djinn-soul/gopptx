package styling

// Dimension is a length that may be stated relative to the slide rather than in
// absolute units. Length is absolute-only — inches, centimetres, points, EMU —
// so placing a box across "half the slide" meant knowing the slide size at the
// call site and doing the arithmetic there.
//
// The zero Dimension is an absolute zero.
type Dimension struct {
	kind  dimensionKind
	value float64
}

type dimensionKind uint8

const (
	dimensionAbsolute dimensionKind = iota
	dimensionRatio
)

// Absolute wraps a fixed length.
func Absolute(length Length) Dimension {
	return Dimension{kind: dimensionAbsolute, value: float64(length)}
}

// Ratio states a fraction of the reference length: 0.5 is half the slide.
// Values outside 0–1 are kept, since a caller may deliberately overflow the
// slide, but see Clamped for the other behaviour.
func Ratio(fraction float64) Dimension {
	return Dimension{kind: dimensionRatio, value: fraction}
}

// PercentOf states a percentage of the reference length: 50 is half the slide.
func PercentOf(percent float64) Dimension {
	return Ratio(percent / percentMax)
}

// IsRelative reports whether resolving this dimension depends on the reference.
func (d Dimension) IsRelative() bool {
	return d.kind == dimensionRatio
}

// Resolve turns the dimension into EMU against a reference length — usually the
// slide width for x and cx, the slide height for y and cy.
func (d Dimension) Resolve(reference Length) Length {
	if d.kind == dimensionRatio {
		return clampToLength(d.value * float64(reference))
	}
	return clampToLength(d.value)
}

// Clamped keeps a ratio inside 0–1, for callers who want a percentage that
// cannot walk off the slide. Absolute dimensions are returned unchanged.
func (d Dimension) Clamped() Dimension {
	if d.kind != dimensionRatio {
		return d
	}
	return Ratio(clampUnit(d.value))
}

// FlexPosition is a point stated in dimensions, so a caller can say "centre
// horizontally, one inch down".
type FlexPosition struct {
	X Dimension
	Y Dimension
}

// FlexSize is an extent stated in dimensions.
type FlexSize struct {
	CX Dimension
	CY Dimension
}

// Resolve turns the position into EMU against a slide size.
func (p FlexPosition) Resolve(slideWidth, slideHeight Length) (Length, Length) {
	return p.X.Resolve(slideWidth), p.Y.Resolve(slideHeight)
}

// Resolve turns the size into EMU against a slide size.
func (s FlexSize) Resolve(slideWidth, slideHeight Length) (Length, Length) {
	return s.CX.Resolve(slideWidth), s.CY.Resolve(slideHeight)
}

// CenterFlex resolves the size against the slide and returns the position that
// centres it. Both axes are resolved first, so it works for absolute and
// relative extents alike.
func CenterFlex(size FlexSize, slideWidth, slideHeight Length) (Length, Length) {
	cx, cy := size.Resolve(slideWidth, slideHeight)
	return (slideWidth - cx) / 2, (slideHeight - cy) / 2
}
