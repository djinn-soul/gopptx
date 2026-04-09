package slide

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
)

var slideShowAttrPattern = regexp.MustCompile(`\s+show\s*=\s*("[^"]*"|'[^']*')`)

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
	if bytes.HasSuffix(tag, []byte("/>")) {
		return append(tag[:len(tag)-2], []byte(` show="0"/>`)...)
	}
	if bytes.HasSuffix(tag, []byte(">")) {
		return append(tag[:len(tag)-1], []byte(` show="0">`)...)
	}
	return tag
}
