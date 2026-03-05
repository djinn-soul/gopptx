package shapes

import (
	"errors"
	"fmt"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
)

// ShadowType represents the type of shadow effect.
type ShadowType string

const (
	// ShadowTypeOuter is an outer shadow (dropshadow).
	ShadowTypeOuter ShadowType = "outer"
	// ShadowTypeInner is an inner shadow.
	ShadowTypeInner ShadowType = "inner"
	// ShadowTypePerspective is a perspective shadow.
	ShadowTypePerspective ShadowType = "perspective"
)

// IsValidShadowType returns true if the shadow type is valid.
func IsValidShadowType(t ShadowType) bool {
	switch t {
	case ShadowTypeOuter, ShadowTypeInner, ShadowTypePerspective:
		return true
	}
	return false
}

// ShadowAlignment represents the alignment of the shadow relative to the shape.
type ShadowAlignment string

const (
	// ShadowAlignTopLeft aligns shadow to top-left.
	ShadowAlignTopLeft ShadowAlignment = "tl"
	// ShadowAlignTop aligns shadow to top.
	ShadowAlignTop ShadowAlignment = "t"
	// ShadowAlignTopRight aligns shadow to top-right.
	ShadowAlignTopRight ShadowAlignment = "tr"
	// ShadowAlignLeft aligns shadow to left.
	ShadowAlignLeft ShadowAlignment = "l"
	// ShadowAlignCenter aligns shadow to center.
	ShadowAlignCenter ShadowAlignment = "ctr"
	// ShadowAlignRight aligns shadow to right.
	ShadowAlignRight ShadowAlignment = "r"
	// ShadowAlignBottomLeft aligns shadow to bottom-left.
	ShadowAlignBottomLeft ShadowAlignment = "bl"
	// ShadowAlignBottom aligns shadow to bottom.
	ShadowAlignBottom ShadowAlignment = "b"
	// ShadowAlignBottomRight aligns shadow to bottom-right.
	ShadowAlignBottomRight ShadowAlignment = "br"
)

// IsValidShadowAlignment returns true if the shadow alignment is valid.
func IsValidShadowAlignment(a ShadowAlignment) bool {
	switch a {
	case ShadowAlignTopLeft, ShadowAlignTop, ShadowAlignTopRight,
		ShadowAlignLeft, ShadowAlignCenter, ShadowAlignRight,
		ShadowAlignBottomLeft, ShadowAlignBottom, ShadowAlignBottomRight:
		return true
	}
	return false
}

// RichShapeShadow provides detailed control over shape shadow effects.
type RichShapeShadow struct {
	Type         ShadowType
	Color        string
	Transparency float64         // 0.0 = opaque, 1.0 = fully transparent
	BlurRadius   int             // Blur radius in EMU
	Distance     int             // Distance from shape in EMU
	Angle        float64         // Direction angle in degrees (0-360)
	Alignment    ShadowAlignment // Alignment relative to shape
	SkewX        float64         // Horizontal skew angle for perspective shadows
	SkewY        float64         // Vertical skew angle for perspective shadows
	ScaleX       float64         // Horizontal scale factor
	ScaleY       float64         // Vertical scale factor
	RotateWithShape bool         // Whether shadow rotates with shape
}

// NewRichShapeShadow creates a new outer shadow with default settings.
func NewRichShapeShadow() *RichShapeShadow {
	return &RichShapeShadow{
		Type:            ShadowTypeOuter,
		Color:           "000000",
		Transparency:    0.6,
		BlurRadius:      40000, // ~1mm blur
		Distance:        20000, // ~0.5mm distance
		Angle:           45,
		Alignment:       ShadowAlignBottomRight,
		ScaleX:          1.0,
		ScaleY:          1.0,
		RotateWithShape: true,
	}
}

// NewOuterShadow creates a new outer shadow with the specified color.
func NewOuterShadow(color string) *RichShapeShadow {
	return NewRichShapeShadow().WithColor(color)
}

// NewInnerShadow creates a new inner shadow with the specified color.
func NewInnerShadow(color string) *RichShapeShadow {
	return NewRichShapeShadow().
		WithType(ShadowTypeInner).
		WithColor(color)
}

// NewPerspectiveShadow creates a new perspective shadow with the specified color.
func NewPerspectiveShadow(color string) *RichShapeShadow {
	return NewRichShapeShadow().
		WithType(ShadowTypePerspective).
		WithColor(color)
}

// WithType sets the shadow type.
func (s *RichShapeShadow) WithType(t ShadowType) *RichShapeShadow {
	s.Type = t
	return s
}

// WithColor sets the shadow color.
func (s *RichShapeShadow) WithColor(color string) *RichShapeShadow {
	s.Color = common.NormalizeHexColor(color)
	return s
}

// WithTransparency sets the shadow transparency (0.0 to 1.0).
func (s *RichShapeShadow) WithTransparency(transparency float64) *RichShapeShadow {
	s.Transparency = transparency
	return s
}

// WithBlurRadius sets the shadow blur radius in EMU.
func (s *RichShapeShadow) WithBlurRadius(radius int) *RichShapeShadow {
	s.BlurRadius = radius
	return s
}

// WithDistance sets the shadow distance from the shape in EMU.
func (s *RichShapeShadow) WithDistance(distance int) *RichShapeShadow {
	s.Distance = distance
	return s
}

// WithAngle sets the shadow direction angle in degrees.
func (s *RichShapeShadow) WithAngle(angle float64) *RichShapeShadow {
	s.Angle = angle
	return s
}

// WithAlignment sets the shadow alignment relative to the shape.
func (s *RichShapeShadow) WithAlignment(alignment ShadowAlignment) *RichShapeShadow {
	s.Alignment = alignment
	return s
}

// WithSkew sets the skew angles for perspective shadows.
func (s *RichShapeShadow) WithSkew(skewX, skewY float64) *RichShapeShadow {
	s.SkewX = skewX
	s.SkewY = skewY
	return s
}

// WithScale sets the scale factors for the shadow.
func (s *RichShapeShadow) WithScale(scaleX, scaleY float64) *RichShapeShadow {
	s.ScaleX = scaleX
	s.ScaleY = scaleY
	return s
}

// WithRotateWithShape sets whether the shadow rotates with the shape.
func (s *RichShapeShadow) WithRotateWithShape(rotate bool) *RichShapeShadow {
	s.RotateWithShape = rotate
	return s
}

// Validate checks the shadow for validity.
func (s *RichShapeShadow) Validate() error {
	if s == nil {
		return nil
	}

	if !IsValidShadowType(s.Type) {
		return fmt.Errorf("invalid shadow type: %s", s.Type)
	}

	if !common.IsHexColor(s.Color) {
		return fmt.Errorf("invalid shadow color: %s", s.Color)
	}

	if s.Transparency < 0 || s.Transparency > 1 {
		return errors.New("transparency must be between 0.0 and 1.0")
	}

	if s.BlurRadius < 0 {
		return errors.New("blur radius cannot be negative")
	}

	if s.Distance < 0 {
		return errors.New("distance cannot be negative")
	}

	if s.Alignment != "" && !IsValidShadowAlignment(s.Alignment) {
		return fmt.Errorf("invalid shadow alignment: %s", s.Alignment)
	}

	return nil
}
