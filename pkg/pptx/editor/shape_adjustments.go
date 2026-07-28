package editor

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// adjustmentScale is the OOXML unit for a preset-geometry adjustment: 100000
// is 1.0, matching how PowerPoint records a dragged adjustment handle.
const adjustmentScale = 100000.0

var (
	// prstGeomAnyPattern also matches a self-closing <a:prstGeom/>, which a
	// shape with no adjustments yet may carry.
	prstGeomAnyPattern = regexp.MustCompile(`(?s)<a:prstGeom\b[^>]*>.*?</a:prstGeom>|<a:prstGeom\b[^>]*/>`)
	prstGeomAttrsPat   = regexp.MustCompile(`(?s)^<a:prstGeom\b([^>]*?)/?>`)
	prstAttrPattern    = regexp.MustCompile(`\bprst="([^"]*)"`)
	avLstPattern       = regexp.MustCompile(`(?s)<a:avLst\b[^>]*>.*?</a:avLst>|<a:avLst\b[^>]*/>`)
)

// SetShapeAdjustments writes the adjustment values of a preset geometry — the
// yellow handles in PowerPoint's UI. An elbow connector routes through the
// point its adj1 names, so setting it is the difference between a connector
// that doubles back and one that goes straight (upstream #1017).
//
// Values are fractions in the same units the reader reports: 0.5 is the
// halfway point, which OOXML stores as "val 50000". Adjustments not named in
// the call are left alone, so one handle can be moved without disturbing the
// rest.
func (e *PresentationEditor) SetShapeAdjustments(
	slideIndex, shapeID int,
	adjustments []common.ShapeAdjustmentValue,
) error {
	if e == nil || e.parts == nil {
		return errors.New("editor cannot be nil")
	}
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return errors.New("slide index out of range")
	}
	if len(adjustments) == 0 {
		return errors.New("no adjustments given")
	}
	for _, adjustment := range adjustments {
		if strings.TrimSpace(adjustment.Name) == "" {
			return errors.New("every adjustment needs a name, for example adj1")
		}
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return fmt.Errorf("read slide part %s: not found", partPath)
	}
	shapes, err := parseSlideShapes(content)
	if err != nil {
		return fmt.Errorf("parse shapes: %w", err)
	}

	setter := &adjustmentSetter{shapeID: shapeID, adjustments: adjustments, source: content}
	updated := replaceShapeNodes(content, shapes, setter.apply)
	if setter.err != nil {
		return setter.err
	}
	if !setter.found {
		return fmt.Errorf("shape id %d not found in part %s", shapeID, partPath)
	}
	e.parts.Set(partPath, updated)
	return nil
}

type adjustmentSetter struct {
	shapeID     int
	adjustments []common.ShapeAdjustmentValue
	source      []byte
	found       bool
	err         error
}

func (a *adjustmentSetter) apply(_ int, s *parsedShape) ([]byte, bool) {
	if a.err != nil || s.ID != a.shapeID {
		return nil, false
	}
	a.found = true

	shapeXML := string(a.source[s.Start:s.End])
	geometry := prstGeomAnyPattern.FindString(shapeXML)
	if geometry == "" {
		a.err = fmt.Errorf(
			"shape id %d has no preset geometry; adjustments apply to preset shapes and connectors",
			a.shapeID,
		)
		return nil, false
	}

	rebuilt, err := setGeometryAdjustments(geometry, a.adjustments)
	if err != nil {
		a.err = err
		return nil, false
	}
	return []byte(strings.Replace(shapeXML, geometry, rebuilt, 1)), true
}

// setGeometryAdjustments writes the preset's complete guide set: the caller's
// values where given, the preset's defaults elsewhere, and any guide already on
// the shape in between.
//
// The completeness matters. PowerPoint refuses to open a file whose
// round2SameRect carries only adj1 — "the file is corrupted and unreadable" —
// but opens the same shape when adj1 and adj2 are both present.
func setGeometryAdjustments(geometry string, adjustments []common.ShapeAdjustmentValue) (string, error) {
	attrs := prstGeomAttrsPat.FindStringSubmatch(geometry)
	if len(attrs) < 2 {
		return "", errors.New("malformed preset geometry")
	}
	presetMatch := prstAttrPattern.FindStringSubmatch(attrs[1])
	if len(presetMatch) < 2 {
		return "", errors.New("preset geometry has no prst attribute")
	}
	preset := presetMatch[1]

	guides, known := adjustableGuidesFor(preset)
	if !known {
		return "", fmt.Errorf(
			"preset %q has no adjustment guides known to this package; "+
				"writing a partial or wrong guide set makes PowerPoint refuse the file",
			preset,
		)
	}

	requested := make(map[string]string, len(adjustments))
	for _, adjustment := range adjustments {
		requested[adjustment.Name] = adjustmentFormula(adjustment)
	}
	for name := range requested {
		if !guideDefined(guides, name) {
			return "", fmt.Errorf("preset %q has no adjustment named %q", preset, name)
		}
	}
	existing := parseGuideList(avLstPattern.FindString(geometry))

	var builder strings.Builder
	builder.WriteString("<a:prstGeom")
	builder.WriteString(attrs[1])
	builder.WriteString("><a:avLst>")
	for _, guide := range guides {
		fmt.Fprintf(&builder, `<a:gd name="%s" fmla="%s"/>`,
			guide.Name, resolveGuideFormula(guide, requested, existing))
	}
	builder.WriteString("</a:avLst></a:prstGeom>")
	return builder.String(), nil
}

// resolveGuideFormula picks what one guide is written as: the caller's value,
// else whatever the shape already carried, else the preset's default.
func resolveGuideFormula(
	guide presetGuide,
	requested, existing map[string]string,
) string {
	if formula, ok := requested[guide.Name]; ok {
		return formula
	}
	if formula, ok := existing[guide.Name]; ok {
		return formula
	}
	return "val " + strconv.Itoa(guide.Default)
}

func guideDefined(guides []presetGuide, name string) bool {
	for _, guide := range guides {
		if guide.Name == name {
			return true
		}
	}
	return false
}

func adjustmentFormula(adjustment common.ShapeAdjustmentValue) string {
	if adjustment.Formula != "" {
		return adjustment.Formula
	}
	return "val " + strconv.Itoa(int(adjustment.Value*adjustmentScale))
}

var guidePattern = regexp.MustCompile(`<a:gd\b[^>]*\bname="([^"]*)"[^>]*\bfmla="([^"]*)"[^>]*/?>`)

func parseGuideList(avLst string) map[string]string {
	guides := make(map[string]string)
	for _, match := range guidePattern.FindAllStringSubmatch(avLst, -1) {
		guides[match[1]] = match[2]
	}
	return guides
}
