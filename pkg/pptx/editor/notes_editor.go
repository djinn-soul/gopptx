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
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

// GetNotes returns the speaker notes for a specific slide.
func (e *PresentationEditor) GetNotes(slideIndex int) (string, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return "", errors.New("slide index out of range")
	}

	notesPart, err := e.resolveSlideNotesPart(slideIndex)
	if err != nil {
		return "", err
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

// ListNotesPlaceholders enumerates placeholder metadata from a slide's notes page.
func (e *PresentationEditor) ListNotesPlaceholders(slideIndex int) ([]Placeholder, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, errors.New("slide index out of range")
	}

	notesPart, err := e.resolveSlideNotesPart(slideIndex)
	if err != nil {
		return nil, err
	}
	if notesPart == "" {
		return []Placeholder{}, nil
	}

	data, ok := e.parts.Get(notesPart)
	if !ok {
		return []Placeholder{}, nil
	}
	return parsePlaceholdersFromSlideXML(data), nil
}

// ListNotesShapes enumerates shape metadata from a slide's notes page.
func (e *PresentationEditor) ListNotesShapes(slideIndex int) ([]common.NotesShapeInfo, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, errors.New("slide index out of range")
	}

	notesPart, err := e.resolveSlideNotesPart(slideIndex)
	if err != nil {
		return nil, err
	}
	if notesPart == "" {
		return []common.NotesShapeInfo{}, nil
	}

	data, ok := e.parts.Get(notesPart)
	if !ok {
		return []common.NotesShapeInfo{}, nil
	}
	parsed, err := scanShapesWithOffsets(data, false)
	if err != nil {
		return nil, err
	}
	shapes := make([]common.NotesShapeInfo, 0, len(parsed))
	for _, shape := range parsed {
		shapes = append(shapes, common.NotesShapeInfo{
			ID:               shape.ID,
			Name:             shape.Name,
			Type:             shape.Type,
			Text:             shape.Text,
			X:                float64(shape.X),
			Y:                float64(shape.Y),
			CX:               float64(shape.W),
			CY:               float64(shape.H),
			PlaceholderIndex: shape.PhIndex,
			PlaceholderType:  shape.PhType,
			HasTextFrame:     shape.TextFrame != nil,
		})
	}
	return shapes, nil
}

// HasNotesSlide reports whether a slide currently has a notes-slide relationship.
func (e *PresentationEditor) HasNotesSlide(slideIndex int) (bool, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return false, errors.New("slide index out of range")
	}

	ref := e.slides[slideIndex]
	rels, err := e.slideRelationships(ref.Part)
	if err != nil {
		return false, err
	}

	for _, rel := range rels {
		if rel.Type == common.RelTypeNotesSlide {
			return true, nil
		}
	}
	return false, nil
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

	paras := []text.Paragraph{{
		Runs: []text.Run{{Text: textContent}},
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

// UpdateNotesMaster configures the global notes master for the presentation.
func (e *PresentationEditor) UpdateNotesMaster(master *elements.NotesMaster) error {
	if e == nil {
		return errors.New("editor cannot be nil")
	}
	if master != nil {
		if err := master.Validate(); err != nil {
			return err
		}
	}
	editorslide.EnsureNotesMasterThemePart(e.parts.Has, e.parts.Get, e.parts.Set)
	if !e.parts.Has("ppt/notesMasters/notesMaster1.xml") {
		e.recalculateNextRelIDNum()
		e.nonSlideRels, e.nextRelIDNum = editorslide.EnsureNotesInfrastructure(
			e.parts.Has,
			e.parts.Set,
			e.nonSlideRels,
			e.nextRelIDNum,
			notesMasterThemeIndex,
		)
	}
	return editorslide.UpdateNotesMasterParts(
		master,
		e.parts.Set,
		notesMasterThemeIndex,
		e.RegisterImage,
	)
}

func (e *PresentationEditor) ensureNotesInfrastructure() {
	editorslide.EnsureNotesMasterThemePart(e.parts.Has, e.parts.Get, e.parts.Set)
	if e.parts.Has("ppt/notesMasters/notesMaster1.xml") {
		return
	}
	e.recalculateNextRelIDNum()
	e.nonSlideRels, e.nextRelIDNum = editorslide.EnsureNotesInfrastructure(
		e.parts.Has,
		e.parts.Set,
		e.nonSlideRels,
		e.nextRelIDNum,
		notesMasterThemeIndex,
	)
}

func (e *PresentationEditor) resolveSlideNotesPart(slideIndex int) (string, error) {
	ref := e.slides[slideIndex]
	rels, err := e.slideRelationships(ref.Part)
	if err != nil {
		return "", err
	}

	for _, rel := range rels {
		if rel.Type == common.RelTypeNotesSlide {
			return path.Join("ppt/slides", rel.Target), nil
		}
	}
	return "", nil
}
