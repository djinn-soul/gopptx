package editor

import (
	"fmt"
	"regexp"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

var (
	shapeCNvPrPattern = regexp.MustCompile(`(?s)(<p:cNvPr\b[^>]*?)(\s*/?>)`)
	shapeDescrPattern = regexp.MustCompile(`\s+descr="[^"]*"`)
	shapeTitlePattern = regexp.MustCompile(`\s+title="[^"]*"`)
)

const shapeMetadataMatchIndexCount = 6

func (u *shapeUpdater) applyAltTextAndTitle(
	xmlData []byte,
	_ *parsedShape,
	replaced bool,
) ([]byte, bool) {
	updateAltText := u.updates.AltText != nil || u.updates.Description != nil
	updateTitle := u.updates.Title != nil
	if !updateAltText && !updateTitle {
		return xmlData, replaced
	}

	altText := ""
	if u.updates.AltText != nil {
		altText = *u.updates.AltText
	} else if u.updates.Description != nil {
		altText = *u.updates.Description
	}
	title := ""
	if updateTitle {
		title = *u.updates.Title
	}

	content := string(xmlData)
	match := shapeCNvPrPattern.FindStringSubmatchIndex(content)
	if len(match) < shapeMetadataMatchIndexCount {
		return xmlData, replaced
	}
	attributes := content[match[2]:match[3]]
	if updateAltText {
		attributes = replaceOptionalShapeMetadataAttribute(
			attributes,
			shapeDescrPattern,
			"descr",
			altText,
		)
	}
	if updateTitle {
		attributes = replaceOptionalShapeMetadataAttribute(
			attributes,
			shapeTitlePattern,
			"title",
			title,
		)
	}
	updated := content[:match[2]] + attributes + content[match[3]:]
	return []byte(updated), true
}

func replaceOptionalShapeMetadataAttribute(
	attributes string,
	pattern *regexp.Regexp,
	name string,
	value string,
) string {
	attributes = pattern.ReplaceAllString(attributes, "")
	if value == "" {
		return attributes
	}
	return attributes + fmt.Sprintf(` %s="%s"`, name, common.XMLEscape(value))
}
