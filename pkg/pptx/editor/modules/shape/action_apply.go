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
	// a:hlinkHover is the correct element inside p:cNvPr. a:hlinkMouseOver is
	// matched too, so a shape written by an older build gets its dead element
	// replaced instead of ending up with both.
	hlinkHoverPattern := regexp.MustCompile(
		`(?s)<a:hlink(?:Hover|MouseOver)\b[^>]*/>|<a:hlink(?:Hover|MouseOver)\b[^>]*>.*?</a:hlink(?:Hover|MouseOver)>`,
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
