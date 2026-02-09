package pptx

import (
	"fmt"
	"strings"
	"sync"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
)

type renderedNotesPart struct {
	slideNumber int
	slideXML    string
	relsXML     string
}

func buildRenderedNotesParts(slides []SlideContent) []renderedNotesPart {
	parts := make([]renderedNotesPart, len(slides))
	var wg sync.WaitGroup

	for i := range slides {
		slide := slides[i]
		if strings.TrimSpace(slide.Notes) == "" {
			continue
		}

		wg.Add(1)
		go func(index int, notes string) {
			defer wg.Done()
			slideNumber := index + 1
			parts[index] = renderedNotesPart{
				slideNumber: slideNumber,
				slideXML:    pptxxml.NotesSlide(notes),
				relsXML:     pptxxml.NotesSlideRelationships(slideNumber),
			}
		}(i, slide.Notes)
	}

	wg.Wait()

	filtered := make([]renderedNotesPart, 0, len(slides))
	for _, part := range parts {
		if part.slideNumber == 0 {
			continue
		}
		filtered = append(filtered, part)
	}
	return filtered
}

func notesSlideNumbers(parts []renderedNotesPart) []int {
	if len(parts) == 0 {
		return nil
	}
	out := make([]int, 0, len(parts))
	for _, part := range parts {
		out = append(out, part.slideNumber)
	}
	return out
}

func notesTargetBySlide(parts []renderedNotesPart) map[int]string {
	if len(parts) == 0 {
		return nil
	}
	targets := make(map[int]string, len(parts))
	for _, part := range parts {
		targets[part.slideNumber] = fmt.Sprintf("../notesSlides/notesSlide%d.xml", part.slideNumber)
	}
	return targets
}
