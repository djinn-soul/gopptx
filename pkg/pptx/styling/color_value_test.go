package styling_test

import (
	"math"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestFromHexAcceptsTheSpellingsCallersUse(t *testing.T) {
	tests := map[string]string{
		"FF0000":   "FF0000",
		"#00ff00":  "00FF00",
		"00f":      "0000FF",
		"336699AA": "336699",
	}
	for input, want := range tests {
		color, ok := styling.FromHex(input)
		if !ok {
			t.Fatalf("FromHex(%q) reported failure", input)
		}
		if color.Hex() != want {
			t.Errorf("FromHex(%q).Hex() = %q, want %q", input, color.Hex(), want)
		}
	}

	if _, ok := styling.FromHex("nonsense"); ok {
		t.Error("FromHex should report failure on an unparsable value")
	}
}

func TestAlphaSurvivesTheRoundTrip(t *testing.T) {
	color, ok := styling.FromHex("112233CC")
	if !ok {
		t.Fatal("FromHex failed on an eight-digit value")
	}
	if color.HexAlpha() != "112233CC" {
		t.Fatalf("HexAlpha = %q, want 112233CC", color.HexAlpha())
	}
	if got := color.Transparency(); math.Abs(got-0.2) > 0.01 {
		t.Fatalf("Transparency = %v, want about 0.2", got)
	}
}

func TestLighterAndDarkerMoveTowardTheExtremes(t *testing.T) {
	mid := styling.RGB(128, 128, 128)
	if lighter := mid.Lighter(1); lighter.Hex() != "FFFFFF" {
		t.Errorf("Lighter(1) = %q, want FFFFFF", lighter.Hex())
	}
	if darker := mid.Darker(1); darker.Hex() != "000000" {
		t.Errorf("Darker(1) = %q, want 000000", darker.Hex())
	}
	if unchanged := mid.Lighter(0); unchanged.Hex() != mid.Hex() {
		t.Errorf("Lighter(0) = %q, want %q", unchanged.Hex(), mid.Hex())
	}
	// Out-of-range amounts clamp rather than overflow the channels.
	if clamped := mid.Lighter(5); clamped.Hex() != "FFFFFF" {
		t.Errorf("Lighter(5) = %q, want FFFFFF", clamped.Hex())
	}
}

func TestMixGrayscaleAndInvert(t *testing.T) {
	red, blue := styling.RGB(255, 0, 0), styling.RGB(0, 0, 255)
	if mixed := red.Mix(blue, 0.5); mixed.Hex() != "800080" {
		t.Errorf("Mix halfway = %q, want 800080", mixed.Hex())
	}
	if gray := styling.RGB(255, 255, 255).Grayscale(); gray.Hex() != "FFFFFF" {
		t.Errorf("Grayscale of white = %q, want FFFFFF", gray.Hex())
	}
	if inverted := styling.RGB(0, 0, 0).Invert(); inverted.Hex() != "FFFFFF" {
		t.Errorf("Invert of black = %q, want FFFFFF", inverted.Hex())
	}
}

func TestContrastAndReadableText(t *testing.T) {
	white, black := styling.RGB(255, 255, 255), styling.RGB(0, 0, 0)
	if got := white.ContrastRatio(black); math.Abs(got-21) > 0.1 {
		t.Fatalf("black on white contrast = %v, want 21", got)
	}
	if got := white.ReadableTextColor().Hex(); got != "000000" {
		t.Errorf("text on white = %q, want 000000", got)
	}
	if got := black.ReadableTextColor().Hex(); got != "FFFFFF" {
		t.Errorf("text on black = %q, want FFFFFF", got)
	}
}
