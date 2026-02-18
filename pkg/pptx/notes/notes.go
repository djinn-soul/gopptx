package notes

import (
	"fmt"
	"strings"
	"sync"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
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
		hasNotes := strings.TrimSpace(slide.Notes) != "" || len(slide.NotesBody) > 0
		if !hasNotes {
			continue
		}

		wg.Add(1)
		go func(index int, s elements.SlideContent) {
			defer wg.Done()
			slideNumber := index + 1

			var body []elements.Paragraph
			if len(s.NotesBody) > 0 {
				body = s.NotesBody
			} else {
				p := elements.NewParagraph()
				p.Runs = append(p.Runs, elements.NewRun(s.Notes))
				body = []elements.Paragraph{p}
			}

			parts[index] = RenderedNotesPart{
				SlideNumber: slideNumber,
				SlideXML:    pptxxml.NotesSlide(body),
				RelsXML:     pptxxml.NotesSlideRelationships(slideNumber),
			}
		}(i, slide)
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

// SlideNumbers extracts the slide numbers that have associated notes.
func SlideNumbers(parts []RenderedNotesPart) []int {
	if len(parts) == 0 {
		return nil
	}
	out := make([]int, 0, len(parts))
	for _, part := range parts {
		out = append(out, part.SlideNumber)
	}
	return out
}

// TargetBySlide returns a map of slide number to its notes slide XML path.
func TargetBySlide(parts []RenderedNotesPart) map[int]string {
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
func WriteNotesFiles(pw *pptxxml.PackageWriter, parts []RenderedNotesPart) error {
	// Note: We use theme1.xml from the main presentation, so we don't need to write theme2.xml anymore.

	for _, part := range parts {
		path := fmt.Sprintf("ppt/notesSlides/notesSlide%d.xml", part.SlideNumber)
		pw.AddPart(path, part.SlideXML)

		relsPath := fmt.Sprintf("ppt/notesSlides/_rels/notesSlide%d.xml.rels", part.SlideNumber)
		pw.AddPart(relsPath, part.RelsXML)
	}
	return nil
}
