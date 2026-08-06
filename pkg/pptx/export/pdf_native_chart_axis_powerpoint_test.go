package export

import (
	"math"
	"testing"
)

// TestValueAxisMatchesPowerPoint pins the auto value axis to PowerPoint's own.
//
// Each row is a bar chart whose data tops out at dataMax, exported to PDF
// through PowerPoint on a full-slide plot; wantMax and wantStep are the axis
// maximum and label interval read back out of that PDF.
func TestValueAxisMatchesPowerPoint(t *testing.T) {
	cases := []struct {
		dataMax  float64
		wantMax  float64
		wantStep float64
	}{
		{1, 1.2, 0.2},
		{3, 3.5, 0.5},
		{9, 10, 1},
		{12, 14, 2},
		{18, 20, 2},
		{27, 30, 5},
		{30, 35, 5},
		{42, 45, 5},
		{55, 60, 10},
		{78, 90, 10},
		{120, 140, 20},
		{250, 300, 50},
		{375, 400, 50},
		{500, 600, 100},
		{900, 1000, 100},
		{1234, 1400, 200},
	}
	for _, c := range cases {
		minV, maxV := niceAxisRange([]float64{c.dataMax / 2, c.dataMax})
		if minV != 0 {
			t.Errorf("dataMax %g: axis min = %g, want 0", c.dataMax, minV)
		}
		if math.Abs(maxV-c.wantMax) > 1e-9 {
			t.Errorf("dataMax %g: axis max = %g, want %g", c.dataMax, maxV, c.wantMax)
		}
		if step := chartAxisStep(maxV-minV, chartAxisTickDensity{}); math.Abs(step-c.wantStep) > 1e-9 {
			t.Errorf("dataMax %g: step = %g, want %g", c.dataMax, step, c.wantStep)
		}
	}
}
