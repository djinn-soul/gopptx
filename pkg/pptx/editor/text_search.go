package editor

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

var textRunPattern = regexp.MustCompile(`(?s)(<a:t(?:\s+[^>]*)?>)(.*?)(</a:t>)`)

// FindAndReplaceInShapes replaces text across slide shapes only. Speaker notes,
// layouts and masters are untouched; use FindAndReplaceInScope to include them.
// It returns the number of replacements made.
func (e *PresentationEditor) FindAndReplaceInShapes(findText, replaceText string) (int, error) {
	return e.FindAndReplaceInScope(findText, replaceText, common.TextScopeSlides)
}

// FindAndReplaceInScope replaces text across the parts named by scope. An empty
// scope means common.TextScopeSlides, which is what FindAndReplaceInShapes has
// always done.
func (e *PresentationEditor) FindAndReplaceInScope(
	findText string,
	replaceText string,
	scope common.TextReplaceScope,
) (int, error) {
	if e == nil {
		return 0, errors.New("editor cannot be nil")
	}
	if strings.TrimSpace(findText) == "" {
		return 0, errors.New("find text cannot be empty")
	}
	parts, err := e.textReplacementParts(scope)
	if err != nil {
		return 0, err
	}

	total := 0
	for _, partPath := range parts {
		content, ok := e.parts.Get(partPath)
		if !ok {
			return 0, fmt.Errorf("read part %s: not found", partPath)
		}
		updated, count := replaceTextRuns(content, findText, replaceText)
		if count > 0 {
			total += count
			e.parts.Set(partPath, updated)
		}
	}
	return total, nil
}

// textReplacementParts lists the parts a scope covers, slides always first and
// in slide order so replacement counts stay stable across calls.
func (e *PresentationEditor) textReplacementParts(
	scope common.TextReplaceScope,
) ([]string, error) {
	parts := make([]string, 0, len(e.slides))
	for i := range e.slides {
		parts = append(parts, e.slides[i].Part)
	}

	switch scope {
	case common.TextScopeSlides, "":
		return parts, nil
	case common.TextScopeSlidesAndNotes:
		return append(parts, e.notesParts()...), nil
	case common.TextScopeAll:
		parts = append(parts, e.notesParts()...)
		for _, prefix := range []string{
			"ppt/slideLayouts/", "ppt/slideMasters/",
			"ppt/notesMasters/", "ppt/handoutMasters/",
		} {
			parts = append(parts, sortedXMLParts(e.parts.KeysWithPrefix(prefix))...)
		}
		return parts, nil
	default:
		return nil, fmt.Errorf(
			"unknown scope %q: want %q, %q or %q",
			scope, common.TextScopeSlides,
			common.TextScopeSlidesAndNotes, common.TextScopeAll,
		)
	}
}

func (e *PresentationEditor) notesParts() []string {
	return sortedXMLParts(e.parts.KeysWithPrefix("ppt/notesSlides/"))
}

// sortedXMLParts keeps only the XML parts of a directory, dropping the _rels
// entries a prefix scan also returns.
func sortedXMLParts(keys []string) []string {
	out := make([]string, 0, len(keys))
	for _, key := range keys {
		if strings.Contains(key, "/_rels/") || !strings.HasSuffix(key, ".xml") {
			continue
		}
		out = append(out, key)
	}
	sort.Strings(out)
	return out
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

// prefilterIsReliable reports whether scanning the raw part bytes can rule a
// slide out. The part stores text XML-escaped and split across <a:t> elements,
// so a needle carrying an escapable character or a space may be present in the
// decoded text while absent from every contiguous byte range. In that case the
// prefilter must not skip the slide.
func prefilterIsReliable(needle string) bool {
	return !strings.ContainsAny(needle, "&<>\"'") && !strings.ContainsAny(needle, " \t\r\n")
}

func contentMatchesTextNeedle(content []byte, query common.ShapeSearchQuery, needles shapeSearchNeedles) bool {
	if needles.textNeedle == nil || !prefilterIsReliable(string(needles.textNeedle)) {
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
