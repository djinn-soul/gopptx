package shape

import (
	"fmt"
	"regexp"
	"strings"
)

func ReplaceOpenCloseCNvPrActions(
	xmlStr string,
	cNvPrOpenClose *regexp.Regexp,
	hlinkClickPattern *regexp.Regexp,
	hlinkHoverPattern *regexp.Regexp,
	hasClickAction bool,
	hasHoverAction bool,
	clickXML string,
	hoverXML string,
) (string, bool) {
	match := cNvPrOpenClose.FindStringSubmatchIndex(xmlStr)
	if match == nil {
		return "", false
	}
	inner := xmlStr[match[4]:match[5]]
	cleanInner := inner
	if hasClickAction {
		cleanInner = RemoveMatchedTags(cleanInner, hlinkClickPattern)
	}
	if hasHoverAction {
		cleanInner = RemoveMatchedTags(cleanInner, hlinkHoverPattern)
	}
	replacement := cleanInner + clickXML + hoverXML
	updated := xmlStr[:match[4]] + replacement + xmlStr[match[5]:]
	return updated, true
}

func ReplaceSelfClosingCNvPrActions(
	xmlStr string,
	cNvPrSelfClosing *regexp.Regexp,
	clickXML string,
	hoverXML string,
) (string, bool) {
	match := cNvPrSelfClosing.FindStringSubmatchIndex(xmlStr)
	if match == nil {
		return "", false
	}
	if clickXML == "" && hoverXML == "" {
		return xmlStr, true
	}
	attrs := xmlStr[match[2]:match[3]]
	replacement := fmt.Sprintf(`<p:cNvPr%s>%s%s</p:cNvPr>`, attrs, clickXML, hoverXML)
	return xmlStr[:match[0]] + replacement + xmlStr[match[1]:], true
}

func RemoveMatchedTags(input string, pattern *regexp.Regexp) string {
	matches := pattern.FindAllStringIndex(input, -1)
	if len(matches) == 0 {
		return input
	}
	var builder strings.Builder
	builder.Grow(len(input))
	last := 0
	for _, match := range matches {
		builder.WriteString(input[last:match[0]])
		last = match[1]
	}
	builder.WriteString(input[last:])
	return builder.String()
}
