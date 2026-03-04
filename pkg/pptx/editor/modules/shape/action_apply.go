package shape

import (
	"errors"
	"regexp"
)

func ApplyCNvPrActions(
	xmlData []byte,
	hasClickAction bool,
	hasHoverAction bool,
	clickXML string,
	hoverXML string,
) ([]byte, error) {
	xmlStr := string(xmlData)
	hlinkClickPattern := regexp.MustCompile(`(?s)<a:hlinkClick\b[^>]*/>|<a:hlinkClick\b[^>]*>.*?</a:hlinkClick>`)
	hlinkHoverPattern := regexp.MustCompile(
		`(?s)<a:hlinkMouseOver\b[^>]*/>|<a:hlinkMouseOver\b[^>]*>.*?</a:hlinkMouseOver>`,
	)
	cNvPrOpenClose := regexp.MustCompile(`(?s)<p:cNvPr\b([^>]*)>(.*?)</p:cNvPr>`)

	if updated, ok := ReplaceOpenCloseCNvPrActions(
		xmlStr,
		cNvPrOpenClose,
		hlinkClickPattern,
		hlinkHoverPattern,
		hasClickAction,
		hasHoverAction,
		clickXML,
		hoverXML,
	); ok {
		return []byte(updated), nil
	}

	cNvPrSelfClosing := regexp.MustCompile(`<p:cNvPr\b([^>]*)/>`)
	if updated, ok := ReplaceSelfClosingCNvPrActions(xmlStr, cNvPrSelfClosing, clickXML, hoverXML); ok {
		return []byte(updated), nil
	}

	if hasClickAction || hasHoverAction {
		return nil, errors.New("shape has no cNvPr node for action update")
	}
	return xmlData, nil
}
