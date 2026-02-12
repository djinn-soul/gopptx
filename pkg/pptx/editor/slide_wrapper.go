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
}

var (
	phPattern     = regexp.MustCompile(`<p:ph\b([^>]*)/?>`)
	phIdxPattern  = regexp.MustCompile(`idx="(\d+)"`)
	phTypePattern = regexp.MustCompile(`type="([^"]*)"`)
)

// Placeholders parses the slide XML and returns all placeholder elements found.
func (s *Slide) Placeholders() ([]Placeholder, error) {
	content, ok := s.editor.parts[s.PartName]
	if !ok {
		return nil, fmt.Errorf("slide part %q not found", s.PartName)
	}

	matches := phPattern.FindAllSubmatch(content, -1)
	result := make([]Placeholder, 0, len(matches))

	for _, match := range matches {
		ph := Placeholder{}
		attrs := string(match[1])

		if m := phIdxPattern.FindStringSubmatch(attrs); m != nil {
			val, _ := strconv.Atoi(m[1])
			ph.Index = val
		}
		if m := phTypePattern.FindStringSubmatch(attrs); m != nil {
			ph.Type = m[1]
		}

		result = append(result, ph)
	}
	return result, nil
}
