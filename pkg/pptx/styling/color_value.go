package styling

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ColorValue is a colour with an alpha channel, and the arithmetic that goes
// with one. The package has always had palette constants, but a caller wanting
// a lighter shade of a theme colour, or a blend of two, had to do the hex maths
// themselves.
//
// The zero value is opaque black.
type ColorValue struct {
	R, G, B uint8
	// A is the opacity: 0 is fully transparent, 255 fully opaque. The zero
	// value would be invisible, so a ColorValue built by hand should use RGB or
	// RGBA rather than a struct literal.
	A uint8
}

const (
	// channelMax is the largest value one 8-bit colour channel holds.
	channelMax     = 255
	fullyOpaque    = channelMax
	hexShortDigits = 3
	hexColorDigits = 6
	hexAlphaDigits = 8

	// Bit positions of the packed channels, most significant first.
	shiftFirst  = 24
	shiftSecond = 16
	shiftThird  = 8
	percentMax  = 100.0
)

// RGB builds an opaque colour.
func RGB(r, g, b uint8) ColorValue {
	return ColorValue{R: r, G: g, B: b, A: fullyOpaque}
}

// RGBA builds a colour with an explicit alpha.
func RGBA(r, g, b, a uint8) ColorValue {
	return ColorValue{R: r, G: g, B: b, A: a}
}

// FromHex parses "RRGGBB", "#RRGGBB", "RGB" or "RRGGBBAA". An unparsable
// string yields opaque black and a false second result, so a caller can tell a
// real black from a typo.
func FromHex(hex string) (ColorValue, bool) {
	// The colour normaliser lives in pkg/pptx/common, which imports this
	// package, so the short-form expansion is repeated here rather than
	// creating a cycle.
	clean := strings.TrimPrefix(strings.TrimSpace(hex), "#")
	if len(clean) == hexShortDigits {
		clean = string([]byte{
			clean[0], clean[0],
			clean[1], clean[1],
			clean[2], clean[2],
		})
	}
	if len(clean) != hexColorDigits && len(clean) != hexAlphaDigits {
		return RGB(0, 0, 0), false
	}

	value, err := strconv.ParseUint(clean, 16, 32)
	if err != nil {
		return RGB(0, 0, 0), false
	}
	//nolint:gosec // each shift is masked down to a byte by the conversion
	if len(clean) == hexColorDigits {
		return RGB(
			uint8(value>>shiftSecond),
			uint8(value>>shiftThird),
			uint8(value),
		), true
	}
	//nolint:gosec // each shift is masked down to a byte by the conversion
	return RGBA(
		uint8(value>>shiftFirst),
		uint8(value>>shiftSecond),
		uint8(value>>shiftThird),
		uint8(value),
	), true
}

// MustFromHex is FromHex for constants and palette entries, where a bad value
// is a programming error rather than input. It falls back to opaque black.
func MustFromHex(hex string) ColorValue {
	color, _ := FromHex(hex)
	return color
}

// Hex renders the colour as the six upper-case digits OOXML wants, dropping
// alpha — which OOXML carries as a separate <a:alpha> element, not in the hex.
func (c ColorValue) Hex() string {
	return fmt.Sprintf("%02X%02X%02X", c.R, c.G, c.B)
}

// HexAlpha renders the colour as eight digits, alpha last.
func (c ColorValue) HexAlpha() string {
	return fmt.Sprintf("%02X%02X%02X%02X", c.R, c.G, c.B, c.A)
}

// Transparency is the alpha expressed the way the OOXML writers take it: 0 is
// opaque, 1 fully transparent.
func (c ColorValue) Transparency() float64 {
	return 1 - float64(c.A)/fullyOpaque
}

// Lighter moves the colour toward white by the given fraction (0–1).
func (c ColorValue) Lighter(amount float64) ColorValue {
	return c.Mix(RGBA(channelMax, channelMax, channelMax, c.A), clampUnit(amount))
}

// Darker moves the colour toward black by the given fraction (0–1).
func (c ColorValue) Darker(amount float64) ColorValue {
	return c.Mix(RGBA(0, 0, 0, c.A), clampUnit(amount))
}

// Opacity returns the colour at the given opacity (0–1), leaving its channels
// alone.
func (c ColorValue) Opacity(opacity float64) ColorValue {
	c.A = uint8(math.Round(clampUnit(opacity) * fullyOpaque))
	return c
}

// Transparent is the colour at zero opacity.
func (c ColorValue) Transparent() ColorValue {
	return c.Opacity(0)
}

// Mix blends toward other by the given weight (0 keeps c, 1 returns other).
func (c ColorValue) Mix(other ColorValue, weight float64) ColorValue {
	w := clampUnit(weight)
	blend := func(a, b uint8) uint8 {
		return uint8(math.Round(float64(a)*(1-w) + float64(b)*w))
	}
	return ColorValue{
		R: blend(c.R, other.R),
		G: blend(c.G, other.G),
		B: blend(c.B, other.B),
		A: blend(c.A, other.A),
	}
}

// Grayscale converts to grey using the Rec. 601 luma weights, which is what
// PowerPoint's own greyscale view uses.
func (c ColorValue) Grayscale() ColorValue {
	//nolint:mnd // Rec. 601 luma coefficients
	luma := 0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)
	value := uint8(math.Round(luma))
	return ColorValue{R: value, G: value, B: value, A: c.A}
}

// Invert flips each channel, leaving alpha alone.
func (c ColorValue) Invert() ColorValue {
	return ColorValue{R: channelMax - c.R, G: channelMax - c.G, B: channelMax - c.B, A: c.A}
}

// Luminance is the relative luminance used by the WCAG contrast formula.
func (c ColorValue) Luminance() float64 {
	channel := func(v uint8) float64 {
		s := float64(v) / fullyOpaque
		//nolint:mnd // WCAG sRGB linearisation constants
		if s <= 0.03928 {
			return s / 12.92
		}
		//nolint:mnd // WCAG sRGB linearisation constants
		return math.Pow((s+0.055)/1.055, 2.4)
	}
	// WCAG luminance coefficients.
	return 0.2126*channel(c.R) + 0.7152*channel(c.G) + 0.0722*channel(c.B)
}

// ContrastRatio is the WCAG contrast between two colours, from 1 (identical) to
// 21 (black on white). Text needs 4.5 to pass AA, 3 at large sizes.
func (c ColorValue) ContrastRatio(other ColorValue) float64 {
	a, b := c.Luminance(), other.Luminance()
	if a < b {
		a, b = b, a
	}
	//nolint:mnd // WCAG contrast formula
	return (a + 0.05) / (b + 0.05)
}

// ReadableTextColor returns black or white, whichever reads better on this
// colour as a background.
func (c ColorValue) ReadableTextColor() ColorValue {
	black, white := RGB(0, 0, 0), RGB(channelMax, channelMax, channelMax)
	if c.ContrastRatio(black) >= c.ContrastRatio(white) {
		return black
	}
	return white
}

// Percent scales a value stated in percent into the 0–1 fraction the methods
// above take, so a caller can write Lighter(Percent(20)).
func Percent(value float64) float64 {
	return value / percentMax
}

func clampUnit(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}
