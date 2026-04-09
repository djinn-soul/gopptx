package slide

import (
	"bytes"
	"regexp"
	"strings"
)

//nolint:gochecknoglobals // Precompiled patterns/tokens are immutable and hot-path.
var (
	titlePlaceholderPattern = regexp.MustCompile(
		`(?i)<p:ph\b[^>]*\btype\s*=\s*(?:"(?:title|ctrTitle)"|'(?:title|ctrTitle)')`,
	)
	aTextOpenPrefix = []byte("<a:t")
	aTextCloseTag   = []byte("</a:t>")
)

func AppendCopySuffixToXML(content []byte) []byte {
	res, _ := replaceTitleLikeText(content, replaceLastTextRun, func(match []byte) []byte {
		return []byte(strings.Replace(string(match), "</a:t>", " (Copy)</a:t>", 1))
	})
	return res
}

func ReplaceAllTitleTextRuns(content []byte, newTitle string) ([]byte, bool) {
	return replaceTitleLikeText(content, replaceAllTextRuns, func(_ []byte) []byte {
		return []byte("<a:t>" + newTitle + "</a:t>")
	})
}

func replaceTitleLikeText(
	content []byte,
	runSelector func([]byte, func([]byte) []byte) ([]byte, bool),
	replaceFn func(match []byte) []byte,
) ([]byte, bool) {
	updated, ok := replaceTitlePlaceholderText(content, runSelector, replaceFn)
	if ok {
		return updated, true
	}
	return runSelector(content, replaceFn)
}

func replaceTitlePlaceholderText(
	content []byte,
	runSelector func([]byte, func([]byte) []byte) ([]byte, bool),
	replaceFn func(match []byte) []byte,
) ([]byte, bool) {
	const (
		shapeStart = "<p:sp"
		shapeEnd   = "</p:sp>"
	)
	searchFrom := 0
	for {
		startIdx := bytes.Index(content[searchFrom:], []byte(shapeStart))
		if startIdx < 0 {
			return content, false
		}
		start := searchFrom + startIdx
		endIdx := bytes.Index(content[start:], []byte(shapeEnd))
		if endIdx < 0 {
			return content, false
		}
		end := start + endIdx + len(shapeEnd)
		shape := content[start:end]
		if isTitlePlaceholderShape(shape) {
			replacedShape, replaced := runSelector(shape, replaceFn)
			if !replaced {
				return content, false
			}
			out := make([]byte, 0, len(content)-len(shape)+len(replacedShape))
			out = append(out, content[:start]...)
			out = append(out, replacedShape...)
			out = append(out, content[end:]...)
			return out, true
		}
		searchFrom = end
	}
}

func isTitlePlaceholderShape(shape []byte) bool {
	if !bytes.Contains(shape, []byte("<p:ph")) {
		return false
	}
	return titlePlaceholderPattern.Match(shape)
}

func replaceLastTextRun(content []byte, replaceFn func(match []byte) []byte) ([]byte, bool) {
	lastStart, lastEnd := -1, -1
	searchFrom := 0
	for {
		start, end, next, ok := findNextATextRun(content, searchFrom)
		if !ok {
			break
		}
		lastStart, lastEnd = start, end
		searchFrom = next
	}
	if lastStart < 0 {
		return content, false
	}
	match := content[lastStart:lastEnd]
	replacement := replaceFn(match)

	out := make([]byte, 0, len(content)-len(match)+len(replacement))
	out = append(out, content[:lastStart]...)
	out = append(out, replacement...)
	out = append(out, content[lastEnd:]...)
	return out, true
}

func replaceAllTextRuns(content []byte, replaceFn func(match []byte) []byte) ([]byte, bool) {
	start, end, next, ok := findNextATextRun(content, 0)
	if !ok {
		return content, false
	}

	var out bytes.Buffer
	out.Grow(len(content))
	out.Write(content[:start])
	out.Write(replaceFn(content[start:end]))
	pos := end

	searchFrom := next
	for {
		runStart, runEnd, nextPos, found := findNextATextRun(content, searchFrom)
		if !found {
			break
		}
		out.Write(content[pos:runStart])
		out.Write(clearTextRun(content[runStart:runEnd]))
		pos = runEnd
		searchFrom = nextPos
	}
	out.Write(content[pos:])
	return out.Bytes(), true
}

func findNextATextRun(content []byte, searchFrom int) (int, int, int, bool) {
	openIdx := bytes.Index(content[searchFrom:], aTextOpenPrefix)
	if openIdx < 0 {
		return 0, 0, searchFrom, false
	}
	start := searchFrom + openIdx
	openEndRel := bytes.IndexByte(content[start:], '>')
	if openEndRel < 0 {
		return 0, 0, searchFrom, false
	}
	opened := start + openEndRel + 1
	closeIdxRel := bytes.Index(content[opened:], aTextCloseTag)
	if closeIdxRel < 0 {
		return 0, 0, searchFrom, false
	}
	end := opened + closeIdxRel + len(aTextCloseTag)
	return start, end, end, true
}

func clearTextRun(match []byte) []byte {
	endOpenTag := bytes.IndexByte(match, '>')
	if endOpenTag < 0 {
		return []byte("<a:t></a:t>")
	}
	openTag := match[:endOpenTag+1]
	out := make([]byte, 0, len(openTag)+len("</a:t>"))
	out = append(out, openTag...)
	out = append(out, []byte("</a:t>")...)
	return out
}
