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

// chartAxisStep picks the tick interval. With a known density it takes the
// finest "nice" interval whose tick count still fits; without one it falls back
// to the fixed-band heuristic in niceStep.
func chartAxisStep(rangeV float64, density chartAxisTickDensity) float64 {
	limit := density.maxTicks()
	if limit < minChartAxisTicks {
		return niceStep(rangeV)
	}
	magnitude := math.Pow(10, math.Floor(math.Log10(rangeV)))
	// Ascending, so the first interval that fits is the finest one that fits.
	for _, m := range niceStepMultipliers {
		step := m * magnitude
		if step <= 0 {
			continue
		}
		if int(math.Ceil(rangeV/step))+1 <= limit {
			return step
		}
	}
	return niceStep(rangeV)
}

// minChartAxisTicks guards against a plot so small that the density calculation
// would ask for one or two ticks and produce a meaningless axis.
const minChartAxisTicks = 3

// niceStepMultipliers are the tick intervals PowerPoint chooses between, scaled
// by the range's magnitude.
//
//nolint:gochecknoglobals // Immutable lookup table shared by both step choosers.
var niceStepMultipliers = []float64{0.1, 0.2, 0.25, 0.5, 1, 2, 2.5, 5, 10, 20, 25, 50, 100}

// chartAxisTickLabels renders the ticks for [minV,maxV] as their drawn strings.
func chartAxisTickLabels(minV, maxV float64, valueFormat string, density chartAxisTickDensity) []string {
	ticks := chartAxisTicks(minV, maxV, density)
	labels := make([]string, 0, len(ticks))
	for _, tick := range ticks {
		labels = append(labels, formatTickValue(tick, valueFormat))
	}
	return labels
}
