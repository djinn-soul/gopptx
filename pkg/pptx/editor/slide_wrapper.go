package editor

import (
	"fmt"
	"regexp"
	"strconv"
)

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

// Placeholder describes a discovered placeholder in an existing slide.
type Placeholder struct {
	Index int
	Type  string
	Name  string
}

var (
	phPattern      = regexp.MustCompile(`(?i)<p:ph\b([^>]*)/?>`)
	phIdxPattern   = regexp.MustCompile(`(?i)\bidx\s*=\s*(?:"(\d+)"|'(\d+)')`)
	phTypePattern  = regexp.MustCompile(`(?i)\btype\s*=\s*(?:"([^"]*)"|'([^']*)')`)
	phNamePattern  = regexp.MustCompile(`(?i)<p:cNvPr\b[^>]*\bname\s*=\s*(?:"([^"]*)"|'([^']*)')`)
	shapeSPPattern = regexp.MustCompile(`(?s)<p:sp\b.*?</p:sp>`)
)

// Placeholders parses the slide XML and returns all placeholder elements found.
func (s *Slide) Placeholders() ([]Placeholder, error) {
	content, ok := s.editor.parts.Get(s.PartName)
	if !ok {
		return nil, fmt.Errorf("slide part %q not found", s.PartName)
	}

	return parsePlaceholdersFromSlideXML(content), nil
}

func parsePlaceholdersFromSlideXML(content []byte) []Placeholder {
	shapeMatches := shapeSPPattern.FindAll(content, -1)
	result := make([]Placeholder, 0, len(shapeMatches))
	for _, shape := range shapeMatches {
		match := phPattern.FindSubmatch(shape)
		if match == nil {
			continue
		}
		ph := Placeholder{}
		attrs := string(match[1])

		if m := phIdxPattern.FindStringSubmatch(attrs); m != nil {
			value := m[1]
			if value == "" {
				value = m[2]
			}
			val, _ := strconv.Atoi(value)
			ph.Index = val
		}
		if m := phTypePattern.FindStringSubmatch(attrs); m != nil {
			value := m[1]
			if value == "" {
				value = m[2]
			}
			ph.Type = value
		}
		if m := phNamePattern.FindSubmatch(shape); m != nil {
			value := string(m[1])
			if value == "" {
				value = string(m[2])
			}
			ph.Name = value
		}

		result = append(result, ph)
	}
	return result
}
