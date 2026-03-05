package shapes

import (
	"errors"
	"fmt"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
)

// FillType represents the type of shape fill.
type FillType string

const (
	// FillTypeSolid indicates a solid color fill.
	FillTypeSolid FillType = "solid"
	// FillTypeGradient indicates a gradient fill.
	FillTypeGradient FillType = "gradient"
	// FillTypePattern indicates a pattern fill.
	FillTypePattern FillType = "pattern"
	// FillTypeNoFill indicates no fill (transparent).
	FillTypeNoFill FillType = "noFill"
)

// PatternType represents predefined pattern types for shape fills.
type PatternType string

const (
	// PatternTypePct5 is 5% fill pattern.
	PatternTypePct5 PatternType = "pct5"
	// PatternTypePct10 is 10% fill pattern.
	PatternTypePct10 PatternType = "pct10"
	// PatternTypePct20 is 20% fill pattern.
	PatternTypePct20 PatternType = "pct20"
	// PatternTypePct25 is 25% fill pattern.
	PatternTypePct25 PatternType = "pct25"
	// PatternTypePct30 is 30% fill pattern.
	PatternTypePct30 PatternType = "pct30"
	// PatternTypePct40 is 40% fill pattern.
	PatternTypePct40 PatternType = "pct40"
	// PatternTypePct50 is 50% fill pattern.
	PatternTypePct50 PatternType = "pct50"
	// PatternTypePct60 is 60% fill pattern.
	PatternTypePct60 PatternType = "pct60"
	// PatternTypePct70 is 70% fill pattern.
	PatternTypePct70 PatternType = "pct70"
	// PatternTypePct75 is 75% fill pattern.
	PatternTypePct75 PatternType = "pct75"
	// PatternTypePct80 is 80% fill pattern.
	PatternTypePct80 PatternType = "pct80"
	// PatternTypePct90 is 90% fill pattern.
	PatternTypePct90 PatternType = "pct90"
	// PatternTypeHorz is horizontal lines pattern.
	PatternTypeHorz PatternType = "horz"
	// PatternTypeVert is vertical lines pattern.
	PatternTypeVert PatternType = "vert"
	// PatternTypeDiagCross is diagonal cross pattern.
	PatternTypeDiagCross PatternType = "diagCross"
	// PatternTypeDiagStripe is diagonal stripe pattern.
	PatternTypeDiagStripe PatternType = "diagStripe"
	// PatternTypeSmCheck is small checkerboard pattern.
	PatternTypeSmCheck PatternType = "smCheck"
	// PatternTypeDnDiag is down diagonal pattern.
	PatternTypeDnDiag PatternType = "dnDiag"
	// PatternTypeUpDiag is up diagonal pattern.
	PatternTypeUpDiag PatternType = "upDiag"
)

// IsValidPatternType returns true if the pattern type is valid.
func IsValidPatternType(t PatternType) bool {
	switch t {
	case PatternTypePct5, PatternTypePct10, PatternTypePct20, PatternTypePct25,
		PatternTypePct30, PatternTypePct40, PatternTypePct50, PatternTypePct60,
		PatternTypePct70, PatternTypePct75, PatternTypePct80, PatternTypePct90,
		PatternTypeHorz, PatternTypeVert, PatternTypeDiagCross, PatternTypeDiagStripe,
		PatternTypeSmCheck, PatternTypeDnDiag, PatternTypeUpDiag:
		return true
	}
	return false
}

// NormalizePatternType normalizes a pattern type string.
func NormalizePatternType(t string) PatternType {
	pt := PatternType(t)
	if IsValidPatternType(pt) {
		return pt
	}
	return PatternTypePct5
}

// SolidFill represents a solid color fill with transparency support.
type SolidFill struct {
	Color        string
	Transparency float64 // 0.0 = opaque, 1.0 = fully transparent
}

// PatternFill represents a pattern fill with foreground and background colors.
type PatternFill struct {
	Pattern PatternType
	FgColor string
	BgColor string
}

// RichShapeFill provides a unified interface for all fill types.
// It supports solid fills, gradient fills, pattern fills, and no-fill.
type RichShapeFill struct {
	Type         FillType
	Solid        *SolidFill
	Gradient     *ShapeGradientFill
	Pattern      *PatternFill
}

// NewSolidFill creates a new solid color fill.
func NewSolidFill(color string) *RichShapeFill {
	return &RichShapeFill{
		Type: FillTypeSolid,
		Solid: &SolidFill{
			Color:        common.NormalizeHexColor(color),
			Transparency: 0.0,
		},
	}
}

// NewNoFill creates a fill that represents "no fill" (transparent).
func NewNoFill() *RichShapeFill {
	return &RichShapeFill{
		Type: FillTypeNoFill,
	}
}

// NewPatternFill creates a new pattern fill with the specified pattern type.
func NewPatternFill(pattern PatternType) *RichShapeFill {
	return &RichShapeFill{
		Type: FillTypePattern,
		Pattern: &PatternFill{
			Pattern: pattern,
			FgColor: "000000",
			BgColor: "FFFFFF",
		},
	}
}

// WithSolid creates a solid fill from the current fill state.
// This is useful for switching fill types or setting initial solid fill.
func (f *RichShapeFill) WithSolid(color string) *RichShapeFill {
	f.Type = FillTypeSolid
	f.Solid = &SolidFill{
		Color:        common.NormalizeHexColor(color),
		Transparency: 0.0,
	}
	f.Gradient = nil
	f.Pattern = nil
	return f
}

// WithTransparency sets the transparency for a solid fill (0.0 to 1.0).
func (f *RichShapeFill) WithTransparency(transparency float64) *RichShapeFill {
	if f.Type == FillTypeSolid && f.Solid != nil {
		f.Solid.Transparency = transparency
	}
	return f
}

// WithGradient sets a gradient fill.
func (f *RichShapeFill) WithGradient(gradient ShapeGradientFill) *RichShapeFill {
	f.Type = FillTypeGradient
	f.Gradient = &gradient
	f.Solid = nil
	f.Pattern = nil
	return f
}

// WithPattern sets a pattern fill.
func (f *RichShapeFill) WithPattern(pattern PatternType) *RichShapeFill {
	f.Type = FillTypePattern
	f.Pattern = &PatternFill{
		Pattern: pattern,
		FgColor: "000000",
		BgColor: "FFFFFF",
	}
	f.Solid = nil
	f.Gradient = nil
	return f
}

// WithPatternColors sets the foreground and background colors for a pattern fill.
func (f *RichShapeFill) WithPatternColors(fgColor, bgColor string) *RichShapeFill {
	if f.Type == FillTypePattern && f.Pattern != nil {
		f.Pattern.FgColor = common.NormalizeHexColor(fgColor)
		f.Pattern.BgColor = common.NormalizeHexColor(bgColor)
	}
	return f
}

// Background marks this fill as "no fill" (transparent background).
func (f *RichShapeFill) Background() *RichShapeFill {
	f.Type = FillTypeNoFill
	f.Solid = nil
	f.Gradient = nil
	f.Pattern = nil
	return f
}

// Foreground marks this fill as a solid fill (default behavior).
func (f *RichShapeFill) Foreground() *RichShapeFill {
	if f.Type == FillTypeNoFill {
		f.Type = FillTypeSolid
		f.Solid = &SolidFill{
			Color:        "FFFFFF",
			Transparency: 0.0,
		}
	}
	return f
}

// Type returns the type of fill.
func (f *RichShapeFill) GetType() FillType {
	if f == nil {
		return FillTypeNoFill
	}
	return f.Type
}

// Validate checks the fill for validity.
func (f *RichShapeFill) Validate() error {
	if f == nil {
		return nil
	}

	switch f.Type {
	case FillTypeSolid:
		if f.Solid == nil {
			return errors.New("solid fill requires Solid to be set")
		}
		if !common.IsHexColor(f.Solid.Color) {
			return fmt.Errorf("invalid solid fill color: %s", f.Solid.Color)
		}
		if f.Solid.Transparency < 0 || f.Solid.Transparency > 1 {
			return errors.New("transparency must be between 0.0 and 1.0")
		}

	case FillTypeGradient:
		if f.Gradient == nil {
			return errors.New("gradient fill requires Gradient to be set")
		}
		if err := f.Gradient.Validate(); err != nil {
			return fmt.Errorf("invalid gradient: %w", err)
		}

	case FillTypePattern:
		if f.Pattern == nil {
			return errors.New("pattern fill requires Pattern to be set")
		}
		if !IsValidPatternType(f.Pattern.Pattern) {
			return fmt.Errorf("invalid pattern type: %s", f.Pattern.Pattern)
		}
		if !common.IsHexColor(f.Pattern.FgColor) {
			return fmt.Errorf("invalid pattern foreground color: %s", f.Pattern.FgColor)
		}
		if !common.IsHexColor(f.Pattern.BgColor) {
			return fmt.Errorf("invalid pattern background color: %s", f.Pattern.BgColor)
		}

	case FillTypeNoFill:
		// No validation needed for no-fill

	default:
		return fmt.Errorf("unknown fill type: %s", f.Type)
	}

	return nil
}
