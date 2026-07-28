package chart

import (
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// parseSeriesDataPointEffects returns, per point index, the effect children of
// each existing c:dPt shape properties. They are outside this API's model, so a
// rebuilt point has to carry them over or a recoloured point loses its shadow.
func parseSeriesDataPointEffects(ser string) map[int]string {
	effects := map[int]string{}
	for _, block := range reDataPointBlock.FindAllString(ser, -1) {
		index := trendlineIntValue(block, reDataPointIdx)
		if index == nil {
			continue
		}
		// Effects nested in <a:ln> belong to the outline, not the shape.
		shape := block
		if before, _, found := strings.Cut(shape, "<a:ln"); found {
			shape = before
		}
		if found := reDataPointEffects.FindAllString(shape, -1); len(found) > 0 {
			effects[*index] = strings.Join(found, "")
		}
	}
	return effects
}

// parseSeriesDataPoints reads the c:dPt run of one series back into the update
// payload shape.
func parseSeriesDataPoints(ser string) []common.ChartDataPoint {
	blocks := reDataPointBlock.FindAllString(ser, -1)
	out := make([]common.ChartDataPoint, 0, len(blocks))
	for _, block := range blocks {
		point := common.ChartDataPoint{}
		if index := trendlineIntValue(block, reDataPointIdx); index != nil {
			point.PointIndex = *index
		}
		point.InvertIfNegative = trendlineBoolValue(block, reDataPointInv)
		point.Bubble3D = trendlineBoolValue(block, reDataPointBub3D)
		point.Explosion = trendlineIntValue(block, reDataPointExpl)
		// The marker carries its own fill and line, so it is read and removed
		// before the point's own shape properties are parsed.
		marker := reDataPointMarker.FindString(block)
		parseDataPointMarker(marker, &point)
		shape := strings.Replace(block, marker, "", 1)
		point.LineWidthEMU = trendlineIntValue(shape, reDataPointWidth)
		point.LineColor = trendlineStringValue(shape, reDataPointLine)
		point.FillColor = parseDataPointFill(shape)
		out = append(out, point)
	}
	return out
}

// parseDataPointMarker reads a c:marker back into the point payload.
func parseDataPointMarker(marker string, point *common.ChartDataPoint) {
	if marker == "" {
		return
	}
	point.MarkerSymbol = trendlineStringValue(marker, reDataPointMarkerSym)
	point.MarkerSize = trendlineIntValue(marker, reDataPointMarkerSize)
	point.MarkerLineColor = trendlineStringValue(marker, reDataPointLine)
	point.MarkerFillColor = parseDataPointFill(marker)
}

// parseDataPointFill reads the shape fill, ignoring a solid fill that belongs
// to the line rather than the point.
func parseDataPointFill(block string) *string {
	fill := block
	if before, _, found := strings.Cut(block, "<a:ln"); found {
		fill = before
	}
	return trendlineStringValue(fill, reDataPointFill)
}

// parseDataPointState reads every c:dPt back, tagged with its series.
func parseDataPointState(xml string) []common.ChartDataPoint {
	out := make([]common.ChartDataPoint, 0)
	for seriesIndex, ser := range reSerBlocks.FindAllString(xml, -1) {
		for _, point := range parseSeriesDataPoints(ser) {
			point.SeriesIndex = seriesIndex
			out = append(out, point)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
