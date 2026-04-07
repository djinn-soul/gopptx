package export

import (
	"math"
	"testing"
)

func TestHorizontalBarGeometry_MixedRange(t *testing.T) {
	plotX, plotW := 10.0, 100.0
	minV, maxV := -50.0, 50.0

	xPos, wPos := horizontalBarGeometry(25, minV, maxV, plotX, plotW)
	if math.Abs(xPos-60) > 1e-9 {
		t.Fatalf("positive x: got %.2f want 60.00", xPos)
	}
	if math.Abs(wPos-25) > 1e-9 {
		t.Fatalf("positive width: got %.2f want 25.00", wPos)
	}

	xNeg, wNeg := horizontalBarGeometry(-25, minV, maxV, plotX, plotW)
	if math.Abs(xNeg-35) > 1e-9 {
		t.Fatalf("negative x: got %.2f want 35.00", xNeg)
	}
	if math.Abs(wNeg-25) > 1e-9 {
		t.Fatalf("negative width: got %.2f want 25.00", wNeg)
	}
}

func TestHorizontalBarGeometry_AllNegative(t *testing.T) {
	plotX, plotW := 0.0, 80.0
	minV, maxV := -100.0, 0.0

	x, w := horizontalBarGeometry(-25, minV, maxV, plotX, plotW)
	if math.Abs(x-60) > 1e-9 {
		t.Fatalf("x: got %.2f want 60.00", x)
	}
	if math.Abs(w-20) > 1e-9 {
		t.Fatalf("width: got %.2f want 20.00", w)
	}
}

func TestNiceStep(t *testing.T) {
	cases := []struct {
		rangeV float64
		want   float64
	}{
		{30, 5},    // 0-30 → step=5 gives 6 ticks → match
		{100, 20},  // 0-100 → step=20 gives 5 ticks
		{6, 1},     // 0-6 → step=1 gives 6 ticks
		{50, 10},   // 0-50 → step=10 gives 5 ticks
		{250, 50},  // 0-250 → step=50 gives 5 ticks
		{0.5, 0.1}, // 0-0.5 → step=0.1 gives 5 ticks
		{27, 5},    // 0-27 → step=5 gives ceil(27/5)=6 ticks
	}
	for _, c := range cases {
		got := niceStep(c.rangeV)
		if math.Abs(got-c.want) > 1e-9 {
			t.Errorf("niceStep(%.4g) = %.4g, want %.4g", c.rangeV, got, c.want)
		}
	}
}
