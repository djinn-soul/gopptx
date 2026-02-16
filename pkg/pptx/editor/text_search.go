package editor

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"regexp"
	"strings"
)

var textRunPattern = regexp.MustCompile(`(?s)(<a:t(?:\s+[^>]*)?>)(.*?)(</a:t>)`)

// FindAndReplaceInShapes performs a global text replacement across slide text runs.
// It returns the number of replacements made.
func (e *PresentationEditor) FindAndReplaceInShapes(findText, replaceText string) (int, error) {
	if e == nil {
		return 0, errors.New("editor cannot be nil")
	}
	if strings.TrimSpace(findText) == "" {
		return 0, errors.New("find text cannot be empty")
	}

	total := 0
	for i := range e.slides {
		partPath := e.slides[i].Part
		content, ok := e.parts.Get(partPath)
		if !ok {
			return 0, fmt.Errorf("read slide part %s: not found", partPath)
		}
		updated, count := replaceTextRuns(content, findText, replaceText)
		if count > 0 {
			total += count
			e.parts.Set(partPath, updated)
		}
	}
	return total, nil
}

func replaceTextRuns(content []byte, findText, replaceText string) ([]byte, int) {
	total := 0
	replaced := textRunPattern.ReplaceAllFunc(content, func(match []byte) []byte {
		sub := textRunPattern.FindSubmatch(match)
		if len(sub) < 4 {
			return match
		}
		openTag := string(sub[1])
		raw := string(sub[2])
		closeTag := string(sub[3])
		unescaped := html.UnescapeString(raw)
		count := strings.Count(unescaped, findText)
		if count == 0 {
			return match
		}
		total += count
		updated := strings.ReplaceAll(unescaped, findText, replaceText)
		return []byte(openTag + escapeXMLText(updated) + closeTag)
	})
	return replaced, total
}

func escapeXMLText(value string) string {
	var b bytes.Buffer
	_ = xml.EscapeText(&b, []byte(value))
	return b.String()
}
