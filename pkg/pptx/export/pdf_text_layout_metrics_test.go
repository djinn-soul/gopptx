package export

import "testing"

func TestFontMetricProfileSelection(t *testing.T) {
	t.Parallel()

	calibri := fontMetricProfile("Calibri")
	if calibri.WidthFactor >= 1.0 {
		t.Fatalf("calibri width factor=%v want < 1.0", calibri.WidthFactor)
	}

	mono := fontMetricProfile("Consolas")
	if mono.KernPairFactor != 0 {
		t.Fatalf("monospace kern pair factor=%v want 0", mono.KernPairFactor)
	}
}

func TestKerningAdjustmentTightPairsAndSpaces(t *testing.T) {
	t.Parallel()

	profile := fontMetricProfile("Calibri")
	tight := kerningAdjustment("To", profile, 20)
	plain := kerningAdjustment("oo", profile, 20)
	if tight >= plain {
		t.Fatalf("tight-pair kerning=%v should be smaller than plain=%v", tight, plain)
	}

	withSpace := kerningAdjustment("a a", profile, 30)
	noSpace := kerningAdjustment("aaa", profile, 30)
	if withSpace >= noSpace {
		t.Fatalf("space kerning=%v should be smaller than no-space=%v", withSpace, noSpace)
	}
}
