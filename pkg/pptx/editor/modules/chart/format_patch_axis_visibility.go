package chart

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const (
	axisTagCategory   = "catAx"
	axisTagDate       = "dateAx"
	axisTagValue      = "valAx"
	axisMajorTickTag  = "<c:majorTickMark"
	axisMinorTickTag  = "<c:minorTickMark"
	axisCrossAxTag    = "<c:crossAx"
	axisCrossesTag    = "<c:crosses"
	axisCrossesAtTag  = "<c:crossesAt"
	axisShapePropsTag = "<c:spPr"
	axisTextPropsTag  = "<c:txPr"

	// tickLabelRotationScale is the OOXML unit for text rotation: 60000ths of
	// a degree, the same unit a shape's <a:xfrm rot> uses.
	tickLabelRotationScale = 60000.0
	maxTickLabelRotation   = 360.0
)

var (
	reAxisDelete       = regexp.MustCompile(`<c:delete\b[^>]*/>`)
	reAxisCrossBetween = regexp.MustCompile(`<c:crossBetween\b[^>]*/>`)
	reAxisTxPrBodyRot  = regexp.MustCompile(`(?s)<c:txPr>\s*<a:bodyPr\b[^>]*/>`)
	reAxisTxPrBlock    = regexp.MustCompile(`(?s)<c:txPr>.*?</c:txPr>`)

	//nolint:gochecknoglobals // static lookup, behaves as a const set
	crossBetweenValues = map[string]bool{"between": true, "midCat": true}
)

// ValidateAxisVisibility checks the axis-visibility, rotation and tick-crossing
// fields of a format request.
func ValidateAxisVisibility(req common.ChartFormatUpdate) error {
	for _, rotation := range []*float64{
		req.CategoryAxisTickLabelRotation,
		req.ValueAxisTickLabelRotation,
	} {
		if rotation == nil {
			continue
		}
		if !isFinite(*rotation) || *rotation < -maxTickLabelRotation || *rotation > maxTickLabelRotation {
			return errors.New("axis tick label rotation must be between -360 and 360 degrees")
		}
	}
	if between := req.ValueAxisCrossBetween; between != nil &&
		!crossBetweenValues[strings.TrimSpace(*between)] {
		return errors.New("value_axis_cross_between must be one of between,midCat")
	}
	return nil
}

// PatchAxisVisibility writes axis visibility, tick-label rotation and where the
// category axis crosses relative to its tick marks.
func PatchAxisVisibility(xml string, req common.ChartFormatUpdate) string {
	for _, axisTag := range []string{axisTagCategory, axisTagDate} {
		xml = patchAxisVisible(xml, axisTag, req.CategoryAxisVisible)
		xml = patchAxisTickLabelRotation(xml, axisTag, req.CategoryAxisTickLabelRotation)
	}
	xml = patchAxisVisible(xml, axisTagValue, req.ValueAxisVisible)
	xml = patchAxisTickLabelRotation(xml, axisTagValue, req.ValueAxisTickLabelRotation)
	return patchAxisCrossBetween(xml, req.ValueAxisCrossBetween)
}

// patchAxisVisible writes <c:delete>, which is how an axis is hidden: the axis
// element stays, so the series keep referring to it, and val="1" tells
// PowerPoint not to draw it.
func patchAxisVisible(xml string, axisTag string, visible *bool) string {
	if visible == nil {
		return xml
	}
	node := `<c:delete val="1"/>`
	if *visible {
		node = `<c:delete val="0"/>`
	}
	return patchEachAxisBlock(xml, axisTag, func(block string) string {
		if reAxisDelete.MatchString(block) {
			return reAxisDelete.ReplaceAllLiteralString(block, node)
		}
		// CT_CatAx and CT_ValAx order axId, scaling, delete, axPos, ...
		return insertBeforeFirstOrClose(block, node, axisTag, []string{
			"<c:axPos", "<c:majorGridlines", "<c:minorGridlines", "<c:title",
			"<c:numFmt", axisMajorTickTag, axisMinorTickTag, axisTickLabelTagPrefix,
			axisShapePropsTag, axisTextPropsTag, axisCrossAxTag,
		})
	})
}

// patchAxisTickLabelRotation writes the rotation of an axis's tick labels
// through <c:txPr><a:bodyPr rot="">, which is where PowerPoint records the
// "Custom angle" of an axis label.
func patchAxisTickLabelRotation(xml string, axisTag string, degrees *float64) string {
	if degrees == nil {
		return xml
	}
	rotation := strconv.Itoa(int(*degrees * tickLabelRotationScale))
	return patchEachAxisBlock(xml, axisTag, func(block string) string {
		if reAxisTxPrBlock.MatchString(block) {
			return rewriteAxisTextRotation(block, rotation)
		}
		node := `<c:txPr><a:bodyPr rot="` + rotation + `" spcFirstLastPara="1" ` +
			`vertOverflow="ellipsis" vert="horz" wrap="square" anchor="ctr" anchorCtr="1"/>` +
			`<a:lstStyle/><a:p><a:pPr><a:defRPr/></a:pPr><a:endParaRPr lang="en-US"/></a:p></c:txPr>`
		// CT_CatAx/CT_ValAx order ..., spPr, txPr, crossAx, ...
		return insertBeforeFirstOrClose(block, node, axisTag, []string{
			axisCrossAxTag, axisCrossesTag, axisCrossesAtTag,
		})
	})
}

// rewriteAxisTextRotation sets rot on an axis's existing <a:bodyPr>, keeping
// every other body property the deck already carries.
func rewriteAxisTextRotation(block string, rotation string) string {
	return reAxisTxPrBlock.ReplaceAllStringFunc(block, func(txPr string) string {
		if !reAxisTxPrBodyRot.MatchString(txPr) {
			return txPr
		}
		return reAxisTxPrBodyRot.ReplaceAllStringFunc(txPr, func(bodyPr string) string {
			return setXMLAttribute(bodyPr, "rot", rotation)
		})
	})
}

// patchAxisCrossBetween writes <c:crossBetween>, the "between tick marks" /
// "on tick marks" choice. It is a child of CT_ValAx even though it describes
// where the category axis sits.
func patchAxisCrossBetween(xml string, between *string) string {
	if between == nil {
		return xml
	}
	node := `<c:crossBetween val="` + strings.TrimSpace(*between) + `"/>`
	return patchEachAxisBlock(xml, axisTagValue, func(block string) string {
		if reAxisCrossBetween.MatchString(block) {
			return reAxisCrossBetween.ReplaceAllLiteralString(block, node)
		}
		// CT_ValAx orders ..., crosses|crossesAt, crossBetween, majorUnit,
		// minorUnit, dispUnits, extLst.
		return insertBeforeFirstOrClose(block, node, axisTagValue, []string{
			"<c:majorUnit", "<c:minorUnit", "<c:dispUnits", chartExtensionListTagPrefix,
		})
	})
}

// setXMLAttribute replaces an attribute on a start tag, or adds it when the tag
// does not carry it yet.
func setXMLAttribute(tag, name, value string) string {
	pattern := regexp.MustCompile(`\b` + regexp.QuoteMeta(name) + `="[^"]*"`)
	if pattern.MatchString(tag) {
		return pattern.ReplaceAllLiteralString(tag, name+`="`+value+`"`)
	}
	insertAt := strings.LastIndex(tag, "/>")
	if insertAt < 0 {
		return tag
	}
	return tag[:insertAt] + ` ` + name + `="` + value + `"` + tag[insertAt:]
}
