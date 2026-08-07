package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// Relative dimensions and colour arithmetic, re-exported so a caller working
// through the top-level package does not have to reach into pkg/pptx/styling.

type (
	// Dimension is a length that may be relative to the slide.
	Dimension = styling.Dimension
	// FlexPosition is a point stated in dimensions.
	FlexPosition = styling.FlexPosition
	// FlexSize is an extent stated in dimensions.
	FlexSize = styling.FlexSize
	// ColorValue is a colour with an alpha channel and colour arithmetic.
	ColorValue = styling.ColorValue
)

// Slide margins.
const (
	Margin      = styling.Margin
	MarginSmall = styling.MarginSmall
	MarginLarge = styling.MarginLarge
)

// Absolute wraps a fixed length as a Dimension.
func Absolute(length styling.Length) Dimension { return styling.Absolute(length) }

// Ratio states a fraction of the slide: 0.5 is half.
func Ratio(fraction float64) Dimension { return styling.Ratio(fraction) }

// PercentOf states a percentage of the slide: 50 is half.
func PercentOf(percent float64) Dimension { return styling.PercentOf(percent) }

// CenterFlex resolves a flexible size and returns the position that centres it.
func CenterFlex(size FlexSize, slideWidth, slideHeight styling.Length) (styling.Length, styling.Length) {
	return styling.CenterFlex(size, slideWidth, slideHeight)
}

// RGB builds an opaque colour.
func RGB(r, g, b uint8) ColorValue { return styling.RGB(r, g, b) }

// RGBA builds a colour with an explicit alpha.
func RGBA(r, g, b, a uint8) ColorValue { return styling.RGBA(r, g, b, a) }

// ColorFromHex parses a hex colour, reporting whether it was understood.
func ColorFromHex(hex string) (ColorValue, bool) { return styling.FromHex(hex) }

// MustColorFromHex parses a hex colour, falling back to opaque black.
func MustColorFromHex(hex string) ColorValue { return styling.MustFromHex(hex) }
