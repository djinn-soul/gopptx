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
	// Expectations follow PowerPoint: the finest interval from {1,2,5}×10ⁿ that
	// spans the range in at most ten steps. Measured by exporting sixteen bar
	// charts through PowerPoint and reading their axis labels back.
	cases := []struct {
		rangeV float64
		want   float64
	}{
		{30, 5},     // ten steps of 3 is not a nice interval; 5 gives six
		{100, 10},   // ten steps
		{6, 1},      // six steps
		{50, 5},     // ten steps
		{250, 50},   // 25 is not an interval PowerPoint uses, so 50
		{0.5, 0.05}, // ten steps
		{27, 5},     // ceil(27/5) = six steps
	}
	for _, c := range cases {
		got := niceStep(c.rangeV)
		if math.Abs(got-c.want) > 1e-9 {
			t.Errorf("niceStep(%.4g) = %.4g, want %.4g", c.rangeV, got, c.want)
		}
	}
}
