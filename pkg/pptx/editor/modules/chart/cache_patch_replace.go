package chart

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func getFieldPattern(tag string) *regexp.Regexp {
	return regexp.MustCompile(`(?s)<c:` + tag + `>.*?</c:` + tag + `>`)
}

func replaceFieldContent(
	seriesXML string,
	fieldTag string,
	formula string,
	strVals []string,
	numVals []float64,
	multiLevelVals [][]string,
) (string, error) {
	field := getFieldPattern(fieldTag).FindString(seriesXML)
	if field == "" {
		return "", fmt.Errorf("missing %s node", fieldTag)
	}

	updatedField := applyFieldFormula(field, formula)
	var err error
	switch {
	case len(multiLevelVals) > 0:
		updatedField, err = applyMultiLevelValues(fieldTag, updatedField, multiLevelVals)
	case len(strVals) > 0:
		updatedField, err = applyStringValues(fieldTag, updatedField, strVals)
	case len(numVals) > 0:
		updatedField, err = applyNumericValues(fieldTag, updatedField, numVals)
	default:
		err = fmt.Errorf("no values provided for %s", fieldTag)
	}
	if err != nil {
		return "", err
	}

	return strings.Replace(seriesXML, field, updatedField, 1), nil
}

func applyFieldFormula(field string, formula string) string {
	formulaNode := xmlFormulaPattern.FindString(field)
	if formulaNode == "" {
		return field
	}
	return strings.Replace(field, formulaNode, "<c:f>"+common.XMLEscape(formula)+"</c:f>", 1)
}

func applyStringValues(fieldTag string, field string, strVals []string) (string, error) {
	switch {
	case strCachePattern.MatchString(field):
		return strCachePattern.ReplaceAllString(field, buildStringData("strCache", strVals)), nil
	case strLitPattern.MatchString(field):
		return strLitPattern.ReplaceAllString(field, buildStringData("strLit", strVals)), nil
	case numCachePattern.MatchString(field):
		numeric, err := convertStringsToFloats(strVals, "numeric cache")
		if err != nil {
			return "", err
		}
		existing := numCachePattern.FindString(field)
		formatCode := extractFormatCode(existing)
		return numCachePattern.ReplaceAllString(field, buildNumberData("numCache", formatCode, numeric)), nil
	case numLitPattern.MatchString(field):
		numeric, err := convertStringsToFloats(strVals, "numeric literal")
		if err != nil {
			return "", err
		}
		existing := numLitPattern.FindString(field)
		formatCode := extractFormatCode(existing)
		return numLitPattern.ReplaceAllString(field, buildNumberData("numLit", formatCode, numeric)), nil
	default:
		return "", fmt.Errorf("missing data node for %s", fieldTag)
	}
}

func applyMultiLevelValues(fieldTag string, field string, vals [][]string) (string, error) {
	switch {
	case multiLvlCache.MatchString(field):
		return multiLvlCache.ReplaceAllString(field, buildMultiLevelData("multiLvlStrCache", vals)), nil
	case multiLvlLit.MatchString(field):
		return multiLvlLit.ReplaceAllString(field, buildMultiLevelData("multiLvlStrLit", vals)), nil
	default:
		return "", fmt.Errorf("missing multi-level category node for %s", fieldTag)
	}
}

func applyNumericValues(fieldTag string, field string, numVals []float64) (string, error) {
	switch {
	case numCachePattern.MatchString(field):
		existing := numCachePattern.FindString(field)
		formatCode := extractFormatCode(existing)
		return numCachePattern.ReplaceAllString(field, buildNumberData("numCache", formatCode, numVals)), nil
	case numLitPattern.MatchString(field):
		existing := numLitPattern.FindString(field)
		formatCode := extractFormatCode(existing)
		return numLitPattern.ReplaceAllString(field, buildNumberData("numLit", formatCode, numVals)), nil
	default:
		return "", fmt.Errorf("missing numeric data node for %s", fieldTag)
	}
}

func convertStringsToFloats(strVals []string, dest string) ([]float64, error) {
	values := make([]float64, 0, len(strVals))
	for _, s := range strVals {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, fmt.Errorf("category %q cannot be represented in %s", s, dest)
		}
		values = append(values, f)
	}
	return values, nil
}

func buildStringData(tag string, vals []string) string {
	var b strings.Builder
	b.WriteString("<c:")
	b.WriteString(tag)
	b.WriteString("><c:ptCount val=\"")
	b.WriteString(strconv.Itoa(len(vals)))
	b.WriteString("\"/>")
	for i, v := range vals {
		b.WriteString("<c:pt idx=\"")
		b.WriteString(strconv.Itoa(i))
		b.WriteString("\"><c:v>")
		b.WriteString(common.XMLEscape(v))
		b.WriteString("</c:v></c:pt>")
	}
	b.WriteString("</c:")
	b.WriteString(tag)
	b.WriteString(">")
	return b.String()
}

func buildNumberData(tag string, formatCode string, vals []float64) string {
	var b strings.Builder
	b.WriteString("<c:")
	b.WriteString(tag)
	b.WriteString("><c:formatCode>")
	b.WriteString(common.XMLEscape(formatCode))
	b.WriteString("</c:formatCode><c:ptCount val=\"")
	b.WriteString(strconv.Itoa(len(vals)))
	b.WriteString("\"/>")
	for i, v := range vals {
		b.WriteString("<c:pt idx=\"")
		b.WriteString(strconv.Itoa(i))
		b.WriteString("\"><c:v>")
		b.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
		b.WriteString("</c:v></c:pt>")
	}
	b.WriteString("</c:")
	b.WriteString(tag)
	b.WriteString(">")
	return b.String()
}

func buildMultiLevelData(tag string, levels [][]string) string {
	leafCount := len(levels[0])
	var b strings.Builder
	b.WriteString("<c:")
	b.WriteString(tag)
	b.WriteString("><c:ptCount val=\"")
	b.WriteString(strconv.Itoa(leafCount))
	b.WriteString("\"/>")
	for _, lvl := range levels {
		b.WriteString("<c:lvl>")
		for i, v := range lvl {
			b.WriteString("<c:pt idx=\"")
			b.WriteString(strconv.Itoa(i))
			b.WriteString("\"><c:v>")
			b.WriteString(common.XMLEscape(v))
			b.WriteString("</c:v></c:pt>")
		}
		b.WriteString("</c:lvl>")
	}
	b.WriteString("</c:")
	b.WriteString(tag)
	b.WriteString(">")
	return b.String()
}

func extractFormatCode(node string) string {
	matches := formatCodePattern.FindStringSubmatch(node)
	if len(matches) > 1 && strings.TrimSpace(matches[1]) != "" {
		return matches[1]
	}
	return "General"
}

func sheetRange(col string, n int) string {
	return fmt.Sprintf("Sheet1!$%s$2:$%s$%d", col, col, n+1)
}

func sheetRangeAcrossColumns(startCol int, endCol int, n int) string {
	start := ColumnName(startCol)
	end := ColumnName(endCol)
	return fmt.Sprintf("Sheet1!$%s$2:$%s$%d", start, end, n+1)
}
