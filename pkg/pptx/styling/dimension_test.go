package styling_test

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestRatioResolvesAgainstTheReference(t *testing.T) {
	half := styling.Ratio(0.5)
	if got := half.Resolve(styling.SlideWidth16x9); got != styling.SlideWidth16x9/2 {
		t.Fatalf("half of a 16:9 slide = %d, want %d", got, styling.SlideWidth16x9/2)
	}
	if !half.IsRelative() {
		t.Error("a ratio should be relative")
	}
}

func TestPercentOfMatchesRatio(t *testing.T) {
	if styling.PercentOf(25).Resolve(styling.SlideWidth4x3) != styling.Ratio(0.25).Resolve(styling.SlideWidth4x3) {
		t.Fatal("PercentOf(25) should equal Ratio(0.25)")
	}
}

func TestAbsoluteIgnoresTheReference(t *testing.T) {
	fixed := styling.Absolute(styling.Inches(2))
	if fixed.IsRelative() {
		t.Error("an absolute dimension should not be relative")
	}
	if got := fixed.Resolve(styling.SlideWidth16x9); got != styling.Inches(2) {
		t.Fatalf("absolute resolve = %d, want %d", got, styling.Inches(2))
	}
}

func TestClampedKeepsRatiosOnTheSlide(t *testing.T) {
	if got := styling.Ratio(1.5).Clamped().Resolve(styling.SlideWidth4x3); got != styling.SlideWidth4x3 {
		t.Fatalf("clamped overflow = %d, want the full width", got)
	}
	// An unclamped ratio is left alone: overflowing the slide can be deliberate.
	if got := styling.Ratio(1.5).Resolve(styling.SlideWidth4x3); got <= styling.SlideWidth4x3 {
		t.Fatal("an unclamped ratio should be allowed past the slide edge")
	}
}

func TestFlexSizeAndCentering(t *testing.T) {
	size := styling.FlexSize{CX: styling.Ratio(0.5), CY: styling.Absolute(styling.Inches(1))}
	cx, cy := size.Resolve(styling.SlideWidth16x9, styling.SlideHeight16x9)
	if cx != styling.SlideWidth16x9/2 || cy != styling.Inches(1) {
		t.Fatalf("resolved size = (%d, %d)", cx, cy)
	}

	x, y := styling.CenterFlex(size, styling.SlideWidth16x9, styling.SlideHeight16x9)
	if x != (styling.SlideWidth16x9-cx)/2 || y != (styling.SlideHeight16x9-cy)/2 {
		t.Fatalf("centred position = (%d, %d)", x, y)
	}
}

func TestFlexPositionResolvesPerAxis(t *testing.T) {
	pos := styling.FlexPosition{X: styling.PercentOf(10), Y: styling.Absolute(styling.Inches(1))}
	x, y := pos.Resolve(styling.SlideWidth4x3, styling.SlideHeight4x3)
	if x != styling.SlideWidth4x3/10 || y != styling.Inches(1) {
		t.Fatalf("resolved position = (%d, %d)", x, y)
	}
}
