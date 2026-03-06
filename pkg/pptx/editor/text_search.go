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
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
)

var textRunPattern = regexp.MustCompile(`(?s)(<a:t(?:\s+[^>]*)?>)(.*?)(</a:t>)`)

const textRunPatternSubmatchSize = 4

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
		if len(sub) < textRunPatternSubmatchSize {
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

// SearchShapes scans all slides and returns shapes matching the query.
func (e *PresentationEditor) SearchShapes(query common.ShapeSearchQuery) ([]common.ShapeSearchResult, error) {
	if e == nil {
		return nil, errors.New("editor cannot be nil")
	}

	if !query.CaseSensitive {
		query.NameContains = strings.ToLower(query.NameContains)
		query.TypeEquals = strings.ToLower(query.TypeEquals)
		query.TextContains = strings.ToLower(query.TextContains)
	}

	results := make([]common.ShapeSearchResult, 0)
	for slideIndex := range e.slides {
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

// Slide is a high-level wrapper around an editable slide.
type Slide struct {
	ID       int64
	PartName string
	editor   *PresentationEditor
}

// GetSlide returns a Slide object for the given index (0-based).
func (e *PresentationEditor) GetSlide(index int) (*Slide, error) {
	if index < 0 || index >= len(e.slides) {
		return nil, fmt.Errorf("slide index %d out of range", index)
	}
	ref := e.slides[index]
	return &Slide{
		ID:       ref.SlideID,
		PartName: ref.Part,
		editor:   e,
	}, nil
}

func (e *PresentationEditor) slideRelationships(slidePart string) ([]common.EditorRelationship, error) {
	return editorslide.Relationships(slidePart, e.parts.Get, parseRelationshipsXML)
}

// Placeholder describes a discovered placeholder in an existing slide.
type Placeholder struct {
	Index int
	Type  string
	Name  string
}

// Placeholders parses the slide XML and returns all placeholder elements found.
func (s *Slide) Placeholders() ([]Placeholder, error) {
	content, ok := s.editor.parts.Get(s.PartName)
	if !ok {
		return nil, fmt.Errorf("slide part %q not found", s.PartName)
	}
	return parsePlaceholdersFromSlideXML(content), nil
}

func parsePlaceholdersFromSlideXML(content []byte) []Placeholder {
	parsed, _ := scanShapesWithOffsets(content, false)
	var result []Placeholder
	for _, s := range parsed {
		if s.PhIndex != -1 {
			result = append(result, Placeholder{
				Index: s.PhIndex,
				Type:  s.PhType,
				Name:  s.Name,
			})
		}
	}
	return result
}

// MoveShapeToFront moves the shape with the given ID to the front of the drawing order.
func (e *PresentationEditor) MoveShapeToFront(slideIndex, shapeID int) error {
	shapes, _, _, _, _, err := e.loadSlideShapesForReorder(slideIndex, shapeID)
	if err != nil {
		return err
	}
	return e.MoveShapeToIndex(slideIndex, shapeID, len(shapes)-1)
}

// MoveShapeToBack moves the shape with the given ID to the back of the drawing order.
func (e *PresentationEditor) MoveShapeToBack(slideIndex, shapeID int) error {
	return e.MoveShapeToIndex(slideIndex, shapeID, 0)
}

// MoveShapeToIndex moves the shape with the given ID to a specific drawing-order index.
// Index 0 is back-most, and len(shapes)-1 is front-most.
func (e *PresentationEditor) MoveShapeToIndex(slideIndex, shapeID, targetIndex int) error {
	shapes, shapeIndex, content, partPath, partFound, err := e.loadSlideShapesForReorder(slideIndex, shapeID)
	if err != nil {
		return err
	}
	if !partFound {
		return fmt.Errorf("read slide part %s: not found", partPath)
	}
	shapeCount := len(shapes)
	if targetIndex < 0 || targetIndex >= shapeCount {
		return fmt.Errorf("shape index %d out of range [0,%d)", targetIndex, shapeCount)
	}
	if shapeCount <= 1 || shapeIndex == targetIndex {
		return nil
	}

	return e.reorderShapeByIndex(content, partPath, shapes, shapeIndex, targetIndex)
}

func (e *PresentationEditor) loadSlideShapesForReorder(
	slideIndex, shapeID int,
) ([]parsedShape, int, []byte, string, bool, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, 0, nil, "", false, errors.New("slide index out of range")
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return nil, 0, nil, partPath, false, nil
	}

	shapes, err := scanShapesWithOffsets(content, false)
	if err != nil {
		return nil, 0, nil, partPath, true, fmt.Errorf("parse shapes: %w", err)
	}

	shapeIndex := -1
	for i, s := range shapes {
		if s.ID == shapeID {
			shapeIndex = i
			break
		}
	}
	if shapeIndex == -1 {
		return nil, 0, nil, partPath, true, fmt.Errorf("shape with ID %d not found", shapeID)
	}
	return shapes, shapeIndex, content, partPath, true, nil
}

func (e *PresentationEditor) reorderShapeByIndex(
	content []byte,
	partPath string,
	shapes []parsedShape,
	shapeIndex, targetIndex int,
) error {
	targetShape := shapes[shapeIndex]
	shapeXML := append([]byte(nil), content[targetShape.Start:targetShape.End]...)

	prefix := append([]byte(nil), content[:shapes[0].Start]...)
	suffix := append([]byte(nil), content[shapes[len(shapes)-1].End:]...)

	gaps := make([][]byte, 0, len(shapes)-1)
	for i := range len(shapes) - 1 {
		gap := append([]byte(nil), content[shapes[i].End:shapes[i+1].Start]...)
		gaps = append(gaps, gap)
	}

	remaining := make([]int, 0, len(shapes)-1)
	for i := range shapes {
		if i == shapeIndex {
			continue
		}
		remaining = append(remaining, i)
	}
	remaining = append(remaining[:targetIndex], append([]int{shapeIndex}, remaining[targetIndex:]...)...)

	var buf bytes.Buffer
	buf.Write(prefix)
	for pos, idx := range remaining {
		if idx == shapeIndex {
			buf.Write(shapeXML)
		} else {
			s := shapes[idx]
			buf.Write(content[s.Start:s.End])
		}
		if pos < len(gaps) {
			buf.Write(gaps[pos])
		}
	}
	buf.Write(suffix)

	e.parts.Set(partPath, buf.Bytes())
	return nil
}
