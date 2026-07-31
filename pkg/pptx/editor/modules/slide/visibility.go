package slide

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
)

var slideShowAttrPattern = regexp.MustCompile(`\s+show\s*=\s*("[^"]*"|'[^']*')`)
var slideShowMasterShapesAttrPattern = regexp.MustCompile(
	`\s+showMasterSp\s*=\s*("[^"]*"|'[^']*')`,
)

func ParseSlideHidden(content []byte) (bool, error) {
	tag, _, _, err := locateSlideStartTag(content)
	if err != nil {
		return false, err
	}
	return slideTagShowIsHidden(tag), nil
}

func RewriteSlideHidden(content []byte, hidden bool) ([]byte, error) {
	tag, start, end, err := locateSlideStartTag(content)
	if err != nil {
		return nil, err
	}
	tag = stripShowAttr(tag)
	if hidden {
		tag = injectHiddenShowAttr(tag)
	}

	out := make([]byte, 0, len(content)-((end+1)-start)+len(tag))
	out = append(out, content[:start]...)
	out = append(out, tag...)
	out = append(out, content[end+1:]...)
	return out, nil
}

// ParseSlideShowMasterShapes reports whether shapes inherited from the slide
// master are visible. OOXML defaults showMasterSp to true when it is absent.
func ParseSlideShowMasterShapes(content []byte) (bool, error) {
	tag, _, _, err := locateSlideStartTag(content)
	if err != nil {
		return false, err
	}
	tagText := string(tag)
	return !strings.Contains(tagText, `showMasterSp="0"`) &&
		!strings.Contains(tagText, `showMasterSp="false"`) &&
		!strings.Contains(tagText, "showMasterSp='0'") &&
		!strings.Contains(tagText, "showMasterSp='false'"), nil
}

// RewriteSlideShowMasterShapes toggles shapes inherited from the slide master.
// The true state is represented by omitting the default-true attribute.
func RewriteSlideShowMasterShapes(content []byte, visible bool) ([]byte, error) {
	tag, start, end, err := locateSlideStartTag(content)
	if err != nil {
		return nil, err
	}
	tag = []byte(slideShowMasterShapesAttrPattern.ReplaceAllString(string(tag), ""))
	if !visible {
		tag = injectRootAttribute(tag, ` showMasterSp="0"`)
	}

	out := make([]byte, 0, len(content)-((end+1)-start)+len(tag))
	out = append(out, content[:start]...)
	out = append(out, tag...)
	out = append(out, content[end+1:]...)
	return out, nil
}

func locateSlideStartTag(content []byte) ([]byte, int, int, error) {
	searchStart := 0
	for {
		start := bytes.Index(content[searchStart:], []byte("<p:sld"))
		if start < 0 {
			return nil, 0, 0, errors.New("slide XML does not contain <p:sld> root")
		}
		start += searchStart
		nameEnd := start + len("<p:sld")
		if nameEnd >= len(content) {
			return nil, 0, 0, errors.New("slide XML has malformed <p:sld> root")
		}

		next := content[nameEnd]
		if next != ' ' && next != '\n' && next != '\r' && next != '\t' && next != '>' {
			searchStart = nameEnd
			continue
		}

		relEnd := bytes.IndexByte(content[start:], '>')
		if relEnd < 0 {
			return nil, 0, 0, errors.New("slide XML has unterminated <p:sld> root tag")
		}
		end := start + relEnd
		return content[start : end+1], start, end, nil
	}
}

func slideTagShowIsHidden(tag []byte) bool {
	tagText := string(tag)
	return strings.Contains(tagText, `show="0"`) ||
		strings.Contains(tagText, `show="false"`) ||
		strings.Contains(tagText, "show='0'") ||
		strings.Contains(tagText, "show='false'")
}

func stripShowAttr(tag []byte) []byte {
	return []byte(slideShowAttrPattern.ReplaceAllString(string(tag), ""))
}

func injectHiddenShowAttr(tag []byte) []byte {
	return injectRootAttribute(tag, ` show="0"`)
}

func injectRootAttribute(tag []byte, attribute string) []byte {
	if bytes.HasSuffix(tag, []byte("/>")) {
		return append(tag[:len(tag)-2], []byte(attribute+"/>")...)
	}
	if bytes.HasSuffix(tag, []byte(">")) {
		return append(tag[:len(tag)-1], []byte(attribute+">")...)
	}
	return tag
}
