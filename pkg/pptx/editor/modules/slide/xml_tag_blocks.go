package slide

import "strings"

func replaceFirstXMLTagBlock(source, tagName, replacement string) (string, bool) {
	start, end, found := findXMLTagBlock(source, tagName, 0)
	if !found {
		return source, false
	}
	return source[:start] + replacement + source[end:], true
}

func replaceAllXMLTagBlocks(source, tagName, replacement string) (string, bool) {
	return replaceAllXMLTagBlocksFunc(source, tagName, func(string) string { return replacement })
}

func replaceAllXMLTagBlocksFunc(source, tagName string, fn func(string) string) (string, bool) {
	searchPos := 0
	writePos := 0
	foundAny := false

	var b strings.Builder
	b.Grow(len(source))

	for {
		start, end, found := findXMLTagBlock(source, tagName, searchPos)
		if !found {
			break
		}
		foundAny = true
		b.WriteString(source[writePos:start])
		b.WriteString(fn(source[start:end]))
		searchPos = end
		writePos = end
	}
	if !foundAny {
		return source, false
	}
	b.WriteString(source[writePos:])
	return b.String(), true
}

func findXMLTagBlock(source, tagName string, offset int) (int, int, bool) {
	if offset < 0 {
		offset = 0
	}
	openPrefix := "<" + tagName
	closeTag := "</" + tagName + ">"
	searchFrom := offset

	for {
		rel := strings.Index(source[searchFrom:], openPrefix)
		if rel < 0 {
			return 0, 0, false
		}
		start := searchFrom + rel
		nameEnd := start + len(openPrefix)
		if nameEnd >= len(source) {
			return 0, 0, false
		}
		next := source[nameEnd]
		if !isTagBoundaryChar(next) {
			searchFrom = nameEnd
			continue
		}

		gtRel := strings.IndexByte(source[nameEnd:], '>')
		if gtRel < 0 {
			return 0, 0, false
		}
		openTagEnd := nameEnd + gtRel
		if isSelfClosingTag(source[start : openTagEnd+1]) {
			return start, openTagEnd + 1, true
		}

		closeRel := strings.Index(source[openTagEnd+1:], closeTag)
		if closeRel < 0 {
			searchFrom = openTagEnd + 1
			continue
		}
		end := openTagEnd + 1 + closeRel + len(closeTag)
		return start, end, true
	}
}

func isTagBoundaryChar(ch byte) bool {
	switch ch {
	case '>', '/', ' ', '\t', '\n', '\r':
		return true
	default:
		return false
	}
}

func isSelfClosingTag(tagOpen string) bool {
	gt := len(tagOpen) - 1
	if gt < 0 || tagOpen[gt] != '>' {
		return false
	}
	idx := gt - 1
	for idx >= 0 {
		switch tagOpen[idx] {
		case ' ', '\t', '\n', '\r':
			idx--
		case '/':
			return true
		default:
			return false
		}
	}
	return false
}
