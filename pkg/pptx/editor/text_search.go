package editor

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"regexp"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
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
	// Single-pass: FindAllSubmatchIndex locates all <a:t> runs in one regex scan.
	// The previous approach called ReplaceAllFunc + FindSubmatch, scanning twice.
	indices := textRunPattern.FindAllSubmatchIndex(content, -1)
	if len(indices) == 0 {
		return content, 0
	}

	total := 0
	var buf bytes.Buffer
	buf.Grow(len(content))
	pos := 0

	for _, idx := range indices {
		// idx layout: [fullStart, fullStop, openStart, openStop, textStart, textStop, closeStart, closeStop]
		fullStart, fullEnd := idx[0], idx[1]
		textStart, textEnd := idx[4], idx[5]

		raw := content[textStart:textEnd]
		unescaped := html.UnescapeString(string(raw))
		count := strings.Count(unescaped, findText)

		buf.Write(content[pos:fullStart])
		pos = fullEnd

		if count == 0 {
			buf.Write(content[fullStart:fullEnd])
			continue
		}

		total += count
		updated := strings.ReplaceAll(unescaped, findText, replaceText)
		buf.Write(content[idx[2]:idx[3]])         // openTag
		_ = xml.EscapeText(&buf, []byte(updated)) // escaped text, direct to buf
		buf.Write(content[idx[6]:idx[7]])         // closeTag
	}

	buf.Write(content[pos:])
	return buf.Bytes(), total
}

// SearchShapes scans all slides and returns shapes matching the query.
func (e *PresentationEditor) SearchShapes(query common.ShapeSearchQuery) ([]common.ShapeSearchResult, error) {
	if e == nil {
		return nil, errors.New("editor cannot be nil")
	}

	query, needles := prepareShapeSearchQuery(query)

	results := make([]common.ShapeSearchResult, 0)
	for slideIndex := range e.slides {
		partPath := e.slides[slideIndex].Part
		content, ok := e.parts.Get(partPath)
		if !ok {
			return nil, fmt.Errorf("read slide part %s: not found", partPath)
		}
		if !contentMatchesTextNeedle(content, query, needles) {
			continue
		}

		shapes, err := e.GetShapes(slideIndex)
		if err != nil {
			return nil, err
		}
		for _, shape := range shapes {
			if !shapeMatchesQuery(shape, query) {
				continue
			}
			results = append(results, common.ShapeSearchResult{
				SlideIndex: slideIndex,
				Shape:      shape,
			})
		}
	}
	return results, nil
}

type shapeSearchNeedles struct {
	textNeedle     []byte
	textNeedleOrig []byte
}

func prepareShapeSearchQuery(query common.ShapeSearchQuery) (common.ShapeSearchQuery, shapeSearchNeedles) {
	needles := shapeSearchNeedles{}
	if query.TextContains != "" && !query.CaseSensitive {
		needles.textNeedleOrig = []byte(query.TextContains)
	}
	if !query.CaseSensitive {
		query.NameContains = strings.ToLower(query.NameContains)
		query.TypeEquals = strings.ToLower(query.TypeEquals)
		query.TextContains = strings.ToLower(query.TextContains)
	}
	if query.TextContains != "" {
		needles.textNeedle = []byte(query.TextContains)
	}
	return query, needles
}

func contentMatchesTextNeedle(content []byte, query common.ShapeSearchQuery, needles shapeSearchNeedles) bool {
	if needles.textNeedle == nil {
		return true
	}
	if query.CaseSensitive {
		return bytes.Contains(content, needles.textNeedle)
	}
	return bytes.Contains(content, needles.textNeedleOrig) || asciiContainsFold(content, needles.textNeedle)
}

// asciiContainsFold reports whether b contains s using zero-allocation ASCII
// case-insensitive comparison. s must already be lowercased. Non-ASCII bytes
// are compared as-is (safe: PPTX text content is UTF-8 but search needles
// are typically ASCII).
func asciiContainsFold(b, s []byte) bool {
	n := len(s)
	if n == 0 {
		return true
	}
	bLen := len(b)
	if bLen < n {
		return false
	}
	first := s[0]
	for i := 0; i <= bLen-n; i++ {
		bc := b[i]
		if bc >= 'A' && bc <= 'Z' {
			bc += 'a' - 'A'
		}
		if bc != first {
			continue
		}
		match := true
		for j := 1; j < n; j++ {
			c := b[i+j]
			if c >= 'A' && c <= 'Z' {
				c += 'a' - 'A'
			}
			if c != s[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func shapeMatchesQuery(shape common.Shape, query common.ShapeSearchQuery) bool {
	name := shape.Name
	typ := shape.Type
	text := shape.Text
	qName := query.NameContains
	qType := query.TypeEquals
	qText := query.TextContains

	if !query.CaseSensitive {
		name = strings.ToLower(name)
		typ = strings.ToLower(typ)
		text = strings.ToLower(text)
	}

	if qName != "" && !strings.Contains(name, qName) {
		return false
	}
	if qType != "" && typ != qType {
		return false
	}
	if qText != "" && !strings.Contains(text, qText) {
		return false
	}
	return true
}
