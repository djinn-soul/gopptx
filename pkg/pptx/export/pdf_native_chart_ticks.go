//nolint:mnd // Tick math uses fixed decimal magnitudes and float-rounding scales.
package export

import "math"

// Axis ticks are produced in exactly one place. Three call sites used to walk
// the range independently — the two frame drawers and the inset measurer — and
// they had to agree exactly or the plot area would be sized for labels that were
// never drawn. Routing them all through chartAxisTicks removes that trap.

// chartAxisTickDensity describes how much room an axis has for its labels, so
// the tick interval can be chosen to fill it the way PowerPoint does rather than
// from a fixed tick-count band.
//
// A zero value means "unknown", and the interval falls back to niceStep.
type chartAxisTickDensity struct {
	// AxisLengthPt is the drawn length of the axis in points.
	AxisLengthPt float64
	// LabelExtentPt is the space one label needs along the axis, including the
	// clearance PowerPoint keeps between neighbours.
	LabelExtentPt float64
}

// maxTicks is how many labels fit along the axis without colliding.
func (d chartAxisTickDensity) maxTicks() int {
	if d.AxisLengthPt <= 0 || d.LabelExtentPt <= 0 {
		return 0
	}
	return int(math.Floor(d.AxisLengthPt / d.LabelExtentPt))
}

// chartAxisTicks returns the tick values an axis draws across [minV, maxV].
//
// The walk accumulates in floating point, so each step is re-rounded to shake
// off the drift that would otherwise make a tick land at 4.999999 and miss the
// loop bound.
func chartAxisTicks(minV, maxV float64, density chartAxisTickDensity) []float64 {
	rangeV := maxV - minV
	if rangeV <= 0 {
		rangeV = 1
	}
	step := chartAxisStep(rangeV, density)
	ticks := make([]float64, 0, 12)
	for tick := minV; tick <= maxV+step*1e-9; tick = math.Round((tick+step)*1e9) / 1e9 {
		ticks = append(ticks, tick)
	}
	return ticks
}

// chartAxisStep picks the tick interval: PowerPoint's own choice for the range,
// coarsened when the axis is too short to carry that many labels.
func chartAxisStep(rangeV float64, density chartAxisTickDensity) float64 {
	budget := maxChartAxisIntervals
	if limit := density.maxTicks(); limit >= minChartAxisTicks && limit-1 < budget {
		budget = limit - 1
	}
	return niceStepWithin(rangeV, budget)
}

// minChartAxisTicks guards against a plot so small that the density calculation
// would ask for one or two ticks and produce a meaningless axis.
const minChartAxisTicks = 3

// niceStepMultipliers are the tick intervals PowerPoint chooses between, scaled
// by the range's magnitude. Ascending, so the first that fits is the finest.
//
// One, two and five only: sixteen charts exported through PowerPoint, with data
// maxima from 1 to 1234, produced intervals of 0.2, 0.5, 1, 2, 5, 10, 20, 50,
// 100 and 200 and never a 2.5 or a 25.
//
//nolint:gochecknoglobals // Immutable lookup table shared by both step choosers.
var niceStepMultipliers = []float64{1, 2, 5}

// chartAxisTickLabels renders the ticks for [minV,maxV] as their drawn strings.
func chartAxisTickLabels(minV, maxV float64, valueFormat string, density chartAxisTickDensity) []string {
	ticks := chartAxisTicks(minV, maxV, density)
	labels := make([]string, 0, len(ticks))
	for _, tick := range ticks {
		labels = append(labels, formatTickValue(tick, valueFormat))
	}
	return labels
}
