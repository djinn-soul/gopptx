package chart

import (
	"regexp"
	"strings"
)

var (
	reDataLabelTextProperties = regexp.MustCompile(`(?s)<c:txPr>.*?</c:txPr>`)
	reDataLabelBodyProperties = regexp.MustCompile(`<a:bodyPr\b[^>]*/?>`)
	reBodyPropertiesWrap      = regexp.MustCompile(`\s+wrap="[^"]*"`)
	reDataLabelTextAnchor     = regexp.MustCompile(
		`<c:(?:dLblPos|showLegendKey|showVal|showCatName|showSerName|showPercent|showBubbleSize)\b`,
	)
)

const dataLabelWrapSquare = "square"

func patchDataLabelWordWrap(block string, wordWrap *bool) string {
	if wordWrap == nil {
		return block
	}
	wrap := xmlValueNone
	if *wordWrap {
		wrap = dataLabelWrapSquare
	}
	if textProperties := reDataLabelTextProperties.FindString(block); textProperties != "" {
		bodyProperties := reDataLabelBodyProperties.FindString(textProperties)
		if bodyProperties == "" {
			return block
		}
		updatedBodyProperties := setDataLabelBodyWrap(bodyProperties, wrap)
		updatedTextProperties := strings.Replace(
			textProperties,
			bodyProperties,
			updatedBodyProperties,
			1,
		)
		return strings.Replace(block, textProperties, updatedTextProperties, 1)
	}

	node := `<c:txPr><a:bodyPr wrap="` + wrap +
		`"/><a:lstStyle/><a:p><a:endParaRPr lang="en-US"/></a:p></c:txPr>`
	if anchor := reDataLabelTextAnchor.FindStringIndex(block); anchor != nil {
		return block[:anchor[0]] + node + block[anchor[0]:]
	}
	return strings.Replace(block, "</c:dLbls>", node+"</c:dLbls>", 1)
}

func setDataLabelBodyWrap(bodyProperties string, wrap string) string {
	attribute := ` wrap="` + wrap + `"`
	if reBodyPropertiesWrap.MatchString(bodyProperties) {
		return reBodyPropertiesWrap.ReplaceAllLiteralString(bodyProperties, attribute)
	}
	if prefix, found := strings.CutSuffix(bodyProperties, "/>"); found {
		return prefix + attribute + "/>"
	}
	return strings.TrimSuffix(bodyProperties, ">") + attribute + ">"
}
