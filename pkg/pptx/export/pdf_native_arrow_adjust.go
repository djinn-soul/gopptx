package export

import (
	"math"
	"strconv"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// Arrow presets are sized by their a:avLst guides, not by fixed fractions of the
// bounding box. ECMA-376 defines rightArrow and its siblings against the shorter
// side of the box: the head is ss*adj2/100000 long, where ss is min(w, h), and
// the shaft is h*adj1/100000 tall. Scaling the head by the *width* instead —
// which this renderer did — turned a 158pt-wide, 11pt-tall connector arrow into
// a 79pt spearhead where PowerPoint draws 5.4pt.

const (
	// OOXML guide values are hundred-thousandths.
	ooxmlGuideScale = 100000.0

	// Preset defaults from ECMA-376 for the arrow family: half the shorter side
	// for the head, half the height for the shaft.
	defaultArrowShaftAdj = 50000.0
	defaultArrowHeadAdj  = 50000.0

	arrowAdjustmentShaft = "adj1"
	arrowAdjustmentHead  = "adj2"
)

// arrowGeometry is one arrow's resolved proportions in points.
type arrowGeometry struct {
	// head is the length of the arrowhead along the arrow's axis.
	head float64
	// shaft is the thickness of the arrow's body across that axis.
	shaft float64
}

// arrowGeometryFor resolves a shape's arrow guides against its box. w and h are
// the box in points; the caller passes them in drawing order, so an up or down
// arrow swaps them.
func arrowGeometryFor(adjustments []shapes.ShapeAdjustment, w, h float64) arrowGeometry {
	shaftAdj := ooxmlGuideValue(adjustments, arrowAdjustmentShaft, defaultArrowShaftAdj)
	headAdj := ooxmlGuideValue(adjustments, arrowAdjustmentHead, defaultArrowHeadAdj)

	shorter := math.Min(w, h)
	head := clampFloat(shorter*headAdj/ooxmlGuideScale, 0, w)
	shaft := clampFloat(h*shaftAdj/ooxmlGuideScale, 0, h)
	return arrowGeometry{head: head, shaft: shaft}
}

// ooxmlGuideValue reads one named a:gd, whose formula is "val <n>" for the
// adjustment guides a preset exposes. A missing or unparsable guide falls back
// to the preset's own default.
func ooxmlGuideValue(adjustments []shapes.ShapeAdjustment, name string, fallback float64) float64 {
	for _, adj := range adjustments {
		if !strings.EqualFold(strings.TrimSpace(adj.Name), name) {
			continue
		}
		formula := strings.TrimSpace(adj.Formula)
		if !strings.HasPrefix(formula, "val ") {
			continue
		}
		value, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimPrefix(formula, "val ")), 64)
		if err != nil {
			continue
		}
		return value
	}
	return fallback
}
