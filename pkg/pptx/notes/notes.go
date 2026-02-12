package notes

import (
	"archive/zip"
	"fmt"
	"strings"
	"sync"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// RenderedNotesPart represents the XML components of a slide's speaker notes.
type RenderedNotesPart struct {
	SlideNumber int
	SlideXML    string
	RelsXML     string
}

// BuildRenderedNotesParts constructs the XML parts for speaker notes for all slides that have them.
func BuildRenderedNotesParts(slides []elements.SlideContent) []RenderedNotesPart {
	parts := make([]RenderedNotesPart, len(slides))
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
			parts[index] = RenderedNotesPart{
				SlideNumber: slideNumber,
				SlideXML:    pptxxml.NotesSlide(notes),
				RelsXML:     pptxxml.NotesSlideRelationships(slideNumber),
			}
		}(i, slide.Notes)
	}

	wg.Wait()

	filtered := make([]RenderedNotesPart, 0, len(slides))
	for _, part := range parts {
		if part.SlideNumber == 0 {
			continue
		}
		filtered = append(filtered, part)
	}
	return filtered
}

// NotesSlideNumbers extracts the slide numbers that have associated notes.
func NotesSlideNumbers(parts []RenderedNotesPart) []int {
	if len(parts) == 0 {
		return nil
	}
	out := make([]int, 0, len(parts))
	for _, part := range parts {
		out = append(out, part.SlideNumber)
	}
	return out
}

// NotesTargetBySlide returns a map of slide number to its notes slide XML path.
func NotesTargetBySlide(parts []RenderedNotesPart) map[int]string {
	if len(parts) == 0 {
		return nil
	}
	targets := make(map[int]string, len(parts))
	for _, part := range parts {
		targets[part.SlideNumber] = fmt.Sprintf("../notesSlides/notesSlide%d.xml", part.SlideNumber)
	}
	return targets
}

// WriteNotesFiles writes all notes-related XML files to the presentation package.
func WriteNotesFiles(zw *zip.Writer, parts []RenderedNotesPart) error {
	if len(parts) > 0 {
		if err := common.WriteFile(zw, "ppt/theme/theme2.xml", pptxxml.Theme(nil)); err != nil {
			return err
		}
	}

	for _, part := range parts {
		path := fmt.Sprintf("ppt/notesSlides/notesSlide%d.xml", part.SlideNumber)
		if err := common.WriteFile(zw, path, part.SlideXML); err != nil {
			return err
		}

		relsPath := fmt.Sprintf("ppt/notesSlides/_rels/notesSlide%d.xml.rels", part.SlideNumber)
		if err := common.WriteFile(zw, relsPath, part.RelsXML); err != nil {
			return err
		}
	}
	return nil
}
