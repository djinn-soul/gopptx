package shape

import (
	"fmt"
	"regexp"
	"strings"
)

func UpdateShapeTransforms(xmlData []byte, x, y, w, h int) []byte {
	offRe := regexp.MustCompile(`(?s)<a:off\b[^>]*/>`)
	extRe := regexp.MustCompile(`(?s)<a:ext\b[^>]*/>`)

	res := offRe.ReplaceAllFunc(xmlData, func(_ []byte) []byte {
		return fmt.Appendf(nil, `<a:off x="%d" y="%d"/>`, x, y)
	})
	res = extRe.ReplaceAllFunc(res, func(_ []byte) []byte {
		return fmt.Appendf(nil, `<a:ext cx="%d" cy="%d"/>`, w, h)
	})
	return res
}

func ReplaceStyleInSpPr(
	xmlData []byte,
	styleXML string,
	applyFill bool,
	applyLine bool,
	applyEffects bool,
) []byte {
	spPrRe := regexp.MustCompile(`(?s)<p:spPr\b([^>]*)>(.*?)</p:spPr>`)
	match := spPrRe.FindSubmatchIndex(xmlData)
	if match == nil {
		return xmlData
	}
	inner := string(xmlData[match[4]:match[5]])
	inner = stripSelectiveStyleBlocks(inner, applyFill, applyLine, applyEffects)
	if styleXML != "" {
		inner = insertStyleAfterPresetGeometry(inner, styleXML)
	}
	inner = normalizeStyleAfterPresetGeometry(inner)
	replacement := fmt.Sprintf(`<p:spPr%s>%s</p:spPr>`, string(xmlData[match[2]:match[3]]), inner)
	updated := string(xmlData[:match[0]]) + replacement + string(xmlData[match[1]:])
	return []byte(updated)
}

func insertStyleAfterPresetGeometry(inner string, styleXML string) string {
	startIdx := strings.Index(inner, "<a:prstGeom")
	if startIdx < 0 {
		return styleXML + inner
	}

	search := inner[startIdx:]
	if closeRel := strings.Index(search, "</a:prstGeom>"); closeRel >= 0 {
		insertAt := startIdx + closeRel + len("</a:prstGeom>")
		return inner[:insertAt] + styleXML + inner[insertAt:]
	}
	if selfCloseRel := strings.Index(search, "/>"); selfCloseRel >= 0 {
		insertAt := startIdx + selfCloseRel + len("/>")
		return inner[:insertAt] + styleXML + inner[insertAt:]
	}
	return styleXML + inner
}

func stripSelectiveStyleBlocks(
	inner string,
	applyFill bool,
	applyLine bool,
	applyEffects bool,
) string {
	linePattern := regexp.MustCompile(`(?s)<a:ln\b[^>]*>.*?</a:ln>|<a:ln\b[^>]*/>`)
	restoreLineBlocks := map[string]string{}

	if applyFill && !applyLine {
		placeholderIndex := 0
		inner = linePattern.ReplaceAllStringFunc(inner, func(lineXML string) string {
			placeholder := fmt.Sprintf("__LINE_BLOCK_%d__", placeholderIndex)
			restoreLineBlocks[placeholder] = lineXML
			placeholderIndex++
			return placeholder
		})
	}
	if applyFill {
		inner = removeFillBlocks(inner)
	}
	if applyLine {
		inner = linePattern.ReplaceAllString(inner, "")
	}
	if applyEffects {
		effectPattern := regexp.MustCompile(`(?s)<a:effectLst\b[^>]*>.*?</a:effectLst>|<a:effectLst\b[^>]*/>`)
		inner = effectPattern.ReplaceAllString(inner, "")
	}
	for placeholder, lineXML := range restoreLineBlocks {
		inner = strings.ReplaceAll(inner, placeholder, lineXML)
	}
	return inner
}

func removeFillBlocks(inner string) string {
	solidPattern := regexp.MustCompile(`(?s)<a:solidFill\b[^>]*>.*?</a:solidFill>|<a:solidFill\b[^>]*/>`)
	noFillPattern := regexp.MustCompile(`(?s)<a:noFill\b[^>]*/>|<a:noFill\b[^>]*>.*?</a:noFill>`)
	gradPattern := regexp.MustCompile(`(?s)<a:gradFill\b[^>]*>.*?</a:gradFill>|<a:gradFill\b[^>]*/>`)
	patternPattern := regexp.MustCompile(`(?s)<a:pattFill\b[^>]*>.*?</a:pattFill>|<a:pattFill\b[^>]*/>`)
	blipPattern := regexp.MustCompile(`(?s)<a:blipFill\b[^>]*>.*?</a:blipFill>|<a:blipFill\b[^>]*/>`)
	groupPattern := regexp.MustCompile(`(?s)<a:grpFill\b[^>]*>.*?</a:grpFill>|<a:grpFill\b[^>]*/>`)

	inner = solidPattern.ReplaceAllString(inner, "")
	inner = noFillPattern.ReplaceAllString(inner, "")
	inner = gradPattern.ReplaceAllString(inner, "")
	inner = patternPattern.ReplaceAllString(inner, "")
	inner = blipPattern.ReplaceAllString(inner, "")
	inner = groupPattern.ReplaceAllString(inner, "")
	return inner
}

func normalizeStyleAfterPresetGeometry(inner string) string {
	geomPattern := regexp.MustCompile(`(?s)<a:prstGeom\b[^>]*(?:/>|>.*?</a:prstGeom>)`)
	geomLoc := geomPattern.FindStringIndex(inner)
	if geomLoc == nil {
		return inner
	}

	stylePattern := regexp.MustCompile(
		`(?s)` +
			`<a:solidFill\b[^>]*>.*?</a:solidFill>|<a:solidFill\b[^>]*/>|` +
			`<a:noFill\b[^>]*>.*?</a:noFill>|<a:noFill\b[^>]*/>|` +
			`<a:gradFill\b[^>]*>.*?</a:gradFill>|<a:gradFill\b[^>]*/>|` +
			`<a:pattFill\b[^>]*>.*?</a:pattFill>|<a:pattFill\b[^>]*/>|` +
			`<a:blipFill\b[^>]*>.*?</a:blipFill>|<a:blipFill\b[^>]*/>|` +
			`<a:grpFill\b[^>]*>.*?</a:grpFill>|<a:grpFill\b[^>]*/>|` +
			`<a:ln\b[^>]*>.*?</a:ln>|<a:ln\b[^>]*/>|` +
			`<a:effectLst\b[^>]*>.*?</a:effectLst>|<a:effectLst\b[^>]*/>`,
	)

	before := inner[:geomLoc[0]]
	styles := stylePattern.FindAllString(before, -1)
	if len(styles) == 0 {
		return inner
	}
	beforeWithoutStyles := stylePattern.ReplaceAllString(before, "")
	geom := inner[geomLoc[0]:geomLoc[1]]
	after := inner[geomLoc[1]:]
	return beforeWithoutStyles + geom + strings.Join(styles, "") + after
}
