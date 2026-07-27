package chart

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

var (
	reAxisTitle   = regexp.MustCompile(`(?s)<c:title>.*?</c:title>`)
	reAxisTitleTx = regexp.MustCompile(`(?s)<a:t>.*?</a:t>`)
	reAxisScaling = regexp.MustCompile(`(?s)<c:scaling(?:\s*/>|>.*?</c:scaling>)`)
	reAxisMax     = regexp.MustCompile(`<c:max val="[^"]*"/>`)
	reAxisMin     = regexp.MustCompile(`<c:min val="[^"]*"/>`)
	reAxisMajor   = regexp.MustCompile(`<c:majorUnit val="[^"]*"/>`)
	reAxisMinor   = regexp.MustCompile(`<c:minorUnit val="[^"]*"/>`)
	reAxisNumFmt  = regexp.MustCompile(`<c:numFmt [^>]*/>`)
)

// attributePatterns holds the attribute readers used by the axis and data-label
// patch loops. Hoisted out of attributeValue so the loops do not recompile a
// regexp on every call; the set of attribute names is fixed and small.
//
//nolint:gochecknoglobals // Package-level compiled patterns, as elsewhere here.
var attributePatterns = map[string]*regexp.Regexp{
	"formatCode":   regexp.MustCompile(`formatCode="([^"]*)"`),
	"sourceLinked": regexp.MustCompile(`sourceLinked="([^"]*)"`),
}

func validateAxisDetails(req common.ChartFormatUpdate) error {
	for _, pair := range []struct {
		name string
		min  *float64
		max  *float64
	}{
		{"category axis", req.CategoryAxisMinimumScale, req.CategoryAxisMaximumScale},
		{"value axis", req.ValueAxisMinimumScale, req.ValueAxisMaximumScale},
	} {
		if err := validateAxisScale(pair.name, pair.min, pair.max); err != nil {
			return err
		}
	}
	for _, unit := range []*float64{
		req.CategoryAxisMajorUnit, req.CategoryAxisMinorUnit,
		req.ValueAxisMajorUnit, req.ValueAxisMinorUnit,
	} {
		if unit != nil && (!isFinite(*unit) || *unit <= 0) {
			return errors.New("axis units must be finite positive numbers")
		}
	}
	for _, format := range []*string{req.CategoryAxisNumberFormat, req.ValueAxisNumberFormat} {
		if format != nil && strings.TrimSpace(*format) == "" {
			return errors.New("axis number format must not be empty")
		}
	}
	return nil
}

func validateAxisScale(name string, minimum *float64, maximum *float64) error {
	if (minimum != nil && !isFinite(*minimum)) || (maximum != nil && !isFinite(*maximum)) {
		return fmt.Errorf("%s scale values must be finite", name)
	}
	if minimum != nil && maximum != nil && *minimum >= *maximum {
		return fmt.Errorf("%s minimum scale must be less than maximum scale", name)
	}
	return nil
}

func isFinite(value float64) bool {
	return !math.IsInf(value, 0) && !math.IsNaN(value)
}

func patchAxisDetails(xml string, req common.ChartFormatUpdate) string {
	xml = patchAxisDetailSet(xml, "catAx", req.CategoryAxisTitle, req.CategoryAxisMinimumScale,
		req.CategoryAxisMaximumScale, req.CategoryAxisMajorUnit, req.CategoryAxisMinorUnit,
		req.CategoryAxisNumberFormat, req.CategoryAxisFormatLinked)
	xml = patchAxisDetailSet(xml, "dateAx", req.CategoryAxisTitle, req.CategoryAxisMinimumScale,
		req.CategoryAxisMaximumScale, req.CategoryAxisMajorUnit, req.CategoryAxisMinorUnit,
		req.CategoryAxisNumberFormat, req.CategoryAxisFormatLinked)
	return patchAxisDetailSet(xml, "valAx", req.ValueAxisTitle, req.ValueAxisMinimumScale,
		req.ValueAxisMaximumScale, req.ValueAxisMajorUnit, req.ValueAxisMinorUnit,
		req.ValueAxisNumberFormat, req.ValueAxisFormatLinked)
}

func patchAxisDetailSet(
	xml string,
	axisTag string,
	title *string,
	minimum *float64,
	maximum *float64,
	majorUnit *float64,
	minorUnit *float64,
	numberFormat *string,
	formatLinked *bool,
) string {
	if title == nil && minimum == nil && maximum == nil && majorUnit == nil && minorUnit == nil &&
		numberFormat == nil && formatLinked == nil {
		return xml
	}
	return patchEachAxisBlock(xml, axisTag, func(block string) string {
		block = patchAxisScaling(block, minimum, maximum)
		block = patchAxisTitle(block, title)
		block = patchAxisNumberFormat(block, numberFormat, formatLinked)
		block = patchAxisUnit(block, reAxisMajor, "majorUnit", majorUnit)
		return patchAxisUnit(block, reAxisMinor, "minorUnit", minorUnit)
	})
}

func patchEachAxisBlock(xml string, axisTag string, patch func(string) string) string {
	startTag, endTag := "<c:"+axisTag+">", "</c:"+axisTag+">"
	start := strings.Index(xml, startTag)
	for start >= 0 {
		endRel := strings.Index(xml[start:], endTag)
		if endRel < 0 {
			return xml
		}
		end := start + endRel + len(endTag)
		block := patch(xml[start:end])
		xml = xml[:start] + block + xml[end:]
		next := strings.Index(xml[start+len(block):], startTag)
		if next < 0 {
			return xml
		}
		start += len(block) + next
	}
	return xml
}

func patchAxisScaling(block string, minimum *float64, maximum *float64) string {
	if minimum == nil && maximum == nil {
		return block
	}
	scaling := reAxisScaling.FindString(block)
	if scaling == "" {
		scaling = "<c:scaling/>"
		insertAt := axisScalingInsertIndex(block)
		if insertAt < 0 {
			return block
		}
		block = block[:insertAt] + scaling + block[insertAt:]
	}
	if maximum != nil {
		scaling = replaceOrInsertAxisNode(scaling, reAxisMax, "<c:max val=\""+formatAxisNumber(*maximum)+"\"/>")
	}
	if minimum != nil {
		scaling = replaceOrInsertAxisNode(scaling, reAxisMin, "<c:min val=\""+formatAxisNumber(*minimum)+"\"/>")
	}
	return reAxisScaling.ReplaceAllLiteralString(block, scaling)
}

func patchAxisTitle(block string, title *string) string {
	if title == nil {
		return block
	}
	node := `<c:title><c:tx><c:rich><a:bodyPr/><a:lstStyle/><a:p><a:r><a:rPr lang="en-US"/><a:t>` +
		xmlEscape(*title) + `</a:t></a:r></a:p></c:rich></c:tx><c:layout/><c:overlay val="0"/></c:title>`
	if current := reAxisTitle.FindString(block); current != "" {
		if reAxisTitleTx.MatchString(current) {
			node = reAxisTitleTx.ReplaceAllLiteralString(current, `<a:t>`+xmlEscape(*title)+`</a:t>`)
		}
		return strings.Replace(block, current, node, 1)
	}
	return insertAxisTitle(block, node)
}

func patchAxisNumberFormat(block string, format *string, linked *bool) string {
	if format == nil && linked == nil {
		return block
	}
	formatCode, sourceLinked := defaultChartNumberFormat, true
	if current := reAxisNumFmt.FindString(block); current != "" {
		formatCode = attributeValue(current, "formatCode", formatCode)
		sourceLinked = attributeValue(current, "sourceLinked", "1") == "1"
	}
	if format != nil {
		formatCode = *format
	}
	if linked != nil {
		sourceLinked = *linked
	}
	node := `<c:numFmt formatCode="` + xmlEscape(formatCode) + `" sourceLinked="` + boolToOneZero(sourceLinked) + `"/>`
	if reAxisNumFmt.MatchString(block) {
		return reAxisNumFmt.ReplaceAllLiteralString(block, node)
	}
	return insertAxisNumberFormat(block, node)
}

func patchAxisUnit(block string, re *regexp.Regexp, tag string, value *float64) string {
	if value == nil {
		return block
	}
	node := `<c:` + tag + ` val="` + formatAxisNumber(*value) + `"/>`
	if re.MatchString(block) {
		return re.ReplaceAllLiteralString(block, node)
	}
	return insertAxisUnit(block, node)
}

func replaceOrInsertAxisNode(value string, re *regexp.Regexp, node string) string {
	if re.MatchString(value) {
		return re.ReplaceAllLiteralString(value, node)
	}
	if prefix, found := strings.CutSuffix(value, "/>"); found {
		return prefix + ">" + node + "</c:scaling>"
	}
	return strings.Replace(value, "</c:scaling>", node+"</c:scaling>", 1)
}

func formatAxisNumber(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func attributeValue(xml string, name string, fallback string) string {
	re, ok := attributePatterns[name]
	if !ok {
		re = regexp.MustCompile(regexp.QuoteMeta(name) + `="([^"]*)"`)
	}
	match := re.FindStringSubmatch(xml)
	if len(match) != 2 {
		return fallback
	}
	return match[1]
}
