package editor

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

// GetNotes returns the speaker notes for a specific slide.
func (e *PresentationEditor) GetNotes(slideIndex int) (string, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return "", errors.New("slide index out of range")
	}

	ref := e.slides[slideIndex]
	rels, err := e.slideRelationships(ref.Part)
	if err != nil {
		return "", err
	}

	var notesPart string
	for _, rel := range rels {
		if rel.Type == common.RelTypeNotesSlide {
			notesPart = path.Join("ppt/slides", rel.Target)
			break
		}
	}

	if notesPart == "" {
		return "", nil // No notes
	}

	data, ok := e.parts.Get(notesPart)
	if !ok {
		return "", nil // Missing part?
	}

	return extractAllText(data), nil
}

// SetNotes updates or creates the speaker notes for a specific slide.
func (e *PresentationEditor) SetNotes(slideIndex int, textContent string) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return errors.New("slide index out of range")
	}

	ref := e.slides[slideIndex]

	// Ensure infrastructure (master, theme)
	e.ensureNotesInfrastructure()

	// Check if notes slide already exists
	rels, err := e.slideRelationships(ref.Part)
	if err != nil {
		return err
	}

	var notesPart string

	for _, rel := range rels {
		if rel.Type == common.RelTypeNotesSlide {
			notesPart = path.Join("ppt/slides", rel.Target)
			break
		}
	}

	paras := []text.TextParagraph{{
		Runs: []text.TextRun{{Text: textContent}},
	}}
	xmlContent := pptxxml.NotesSlide(paras)

	if notesPart != "" {
		// Update existing
		e.parts.Set(notesPart, []byte(xmlContent))
		return nil
	}

	// Create new notes slide
	e.recalculateNextRelIDNum()
	if e.nextNotesNum < 1 {
		e.nextNotesNum = 1
	}
	notesNum := e.nextNotesNum
	e.nextNotesNum++

	newNotesPart := fmt.Sprintf("ppt/notesSlides/notesSlide%d.xml", notesNum)
	newNotesRelID := fmt.Sprintf("rId%d", e.nextRelIDNum)
	e.nextRelIDNum++

	e.parts.Set(newNotesPart, []byte(xmlContent))

	// Add relationship to Slide -> Notes
	slideRelsPart := common.SlideRelsPartName(ref.Part)
	slideRelsData, _ := e.parts.Get(slideRelsPart)

	slideRels, err := parseRelationshipsXML(slideRelsData)
	if err != nil {
		return err
	}

	slideRels = append(slideRels, common.EditorRelationship{
		ID:     newNotesRelID,
		Type:   common.RelTypeNotesSlide,
		Target: fmt.Sprintf("../notesSlides/notesSlide%d.xml", notesNum),
	})

	renderedSlideRels := renderRelationshipsXML(slideRels)
	e.parts.Set(slideRelsPart, []byte(renderedSlideRels))

	// Create Notes -> Slide relationship
	slideNum, ok := common.ParseSlidePartNumber(ref.Part)
	if !ok {
		return fmt.Errorf("could not parse slide number from %q", ref.Part)
	}
	notesRelsContent := pptxxml.NotesSlideRelationships(slideNum)

	e.parts.Set(common.SlideRelsPartName(newNotesPart), []byte(notesRelsContent))

	// Update inventory
	if e.notesInventory == nil {
		e.notesInventory = make(map[string]string)
	}
	e.notesInventory[ref.Part] = newNotesPart

	return nil
}

func extractAllText(content []byte) string {
	decoder := xml.NewDecoder(bytes.NewReader(content))
	var sb strings.Builder
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "t" {
			continue
		}
		var value string
		if decodeErr := decoder.DecodeElement(&value, &start); decodeErr == nil {
			sb.WriteString(value)
			sb.WriteString("\n")
		}
	}
	return strings.TrimSpace(sb.String())
}
