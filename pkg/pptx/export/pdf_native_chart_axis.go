package export

import "math"

// horizontalBarGeometry maps a value to a horizontal bar segment in plot space.
// It returns the left X and width so negative values extend left of the zero axis.
func horizontalBarGeometry(value, minV, maxV, plotX, plotW float64) (float64, float64) {
	rangeV := maxV - minV
	if rangeV <= 0 {
		rangeV = 1
	}
	zeroX := plotX + ((0-minV)/rangeV)*plotW
	valueX := plotX + ((value-minV)/rangeV)*plotW
	return math.Min(zeroX, valueX), math.Abs(valueX - zeroX)
}
