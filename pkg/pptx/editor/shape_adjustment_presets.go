package editor

// presetAdjustmentDefaults maps a preset geometry to its adjustment guides, in
// the order the preset defines them, with the default each guide carries.
//
// A partial <a:avLst> makes PowerPoint refuse to open the file: writing
// round2SameRect with only adj1 produces "the file is corrupted and
// unreadable", while writing adj1 and adj2 together opens fine. So an
// adjustment write always emits the preset's complete set, filling unnamed
// guides with these defaults.
//
// Presets absent from this table are rejected rather than guessed at, because a
// wrong guide name or a missing sibling is the same corruption.
//
//nolint:gochecknoglobals // static OOXML preset data, behaves as a const map
var presetAdjustmentDefaults = map[string][]presetGuide{
	// Rectangles with rounded or snipped corners.
	"roundRect":      {{guideAdj, 16667}},
	"round1Rect":     {{guideAdj, 16667}},
	"round2SameRect": {{guideAdj1, 16667}, {guideAdj2, 0}},
	"round2DiagRect": {{guideAdj1, 16667}, {guideAdj2, 0}},
	"snip1Rect":      {{guideAdj, 16667}},
	"snip2SameRect":  {{guideAdj1, 16667}, {guideAdj2, 0}},
	"snip2DiagRect":  {{guideAdj1, 16667}, {guideAdj2, 0}},
	"snipRoundRect":  {{guideAdj1, 16667}, {guideAdj2, 16667}},

	// Block arrows.
	"leftArrow":      {{guideAdj1, 50000}, {guideAdj2, 50000}},
	"rightArrow":     {{guideAdj1, 50000}, {guideAdj2, 50000}},
	"upArrow":        {{guideAdj1, 50000}, {guideAdj2, 50000}},
	"downArrow":      {{guideAdj1, 50000}, {guideAdj2, 50000}},
	"leftRightArrow": {{guideAdj1, 50000}, {guideAdj2, 50000}},
	"upDownArrow":    {{guideAdj1, 50000}, {guideAdj2, 50000}},

	// Single-guide shapes.
	"triangle":      {{guideAdj, 50000}},
	"parallelogram": {{guideAdj, 25000}},
	"trapezoid":     {{guideAdj, 25000}},
	"homePlate":     {{guideAdj, 50000}},
	"chevron":       {{guideAdj, 50000}},
	"can":           {{guideAdj, 25000}},
	"cube":          {{guideAdj, 25000}},
	"donut":         {{guideAdj, 25000}},
	"plus":          {{guideAdj, 25000}},
	"teardrop":      {{guideAdj, 100000}},
	"bevel":         {{guideAdj, 12500}},

	// Elbow connectors: the handle that decides where the bend sits, which is
	// the case upstream #1017 is about.
	"bentConnector3":   {{guideAdj1, 50000}},
	"curvedConnector3": {{guideAdj1, 50000}},
}

// Guide names shared by many presets.
const (
	guideAdj  = "adj"
	guideAdj1 = "adj1"
	guideAdj2 = "adj2"
)

// presetGuide is one adjustment guide of a preset geometry: its name and the
// value OOXML gives it when the shape is first drawn.
type presetGuide struct {
	Name    string
	Default int
}

// adjustableGuidesFor reports the guides a preset defines, and whether the
// preset is one this package knows how to adjust safely.
func adjustableGuidesFor(preset string) ([]presetGuide, bool) {
	guides, ok := presetAdjustmentDefaults[preset]
	return guides, ok
}
