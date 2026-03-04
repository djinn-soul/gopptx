package slide

import (
	"bytes"
	"regexp"
	"strings"
)

var (
	aTTitlePattern          = regexp.MustCompile(`(?s)<a:t(?:\s+[^>]*)?>(.*?)</a:t>`)
	titlePlaceholderPattern = regexp.MustCompile(
		`(?i)<p:ph\b[^>]*\btype\s*=\s*(?:"(?:title|ctrTitle)"|'(?:title|ctrTitle)')`,
	)
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
	indexes := aTTitlePattern.FindAllIndex(content, -1)
	if len(indexes) == 0 {
		return content, false
	}
	last := indexes[len(indexes)-1]
	start, end := last[0], last[1]
	match := content[start:end]
	replacement := replaceFn(match)

	out := make([]byte, 0, len(content)-len(match)+len(replacement))
	out = append(out, content[:start]...)
	out = append(out, replacement...)
	out = append(out, content[end:]...)
	return out, true
}

func replaceAllTextRuns(content []byte, replaceFn func(match []byte) []byte) ([]byte, bool) {
	modified := false
	replacedFirst := false
	res := aTTitlePattern.ReplaceAllFunc(content, func(match []byte) []byte {
		modified = true
		if !replacedFirst {
			replacedFirst = true
			return replaceFn(match)
		}
		return clearTextRun(match)
	})
	return res, modified
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
