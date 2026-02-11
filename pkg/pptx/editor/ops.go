package editor

import (
	"fmt"

	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// AddSlide appends a new slide and returns its 0-based index.
func (e *PresentationEditor) AddSlide(slide elements.SlideContent) (int, error) {
	if e == nil {
		return 0, fmt.Errorf("editor cannot be nil")
	}
	if err := validateEditorSlideContent(slide); err != nil {
		return 0, err
	}

	slideNumber := e.nextSlideNum
	relID := fmt.Sprintf("rId%d", e.nextRelIDNum)
	part := fmt.Sprintf("ppt/slides/slide%d.xml", slideNumber)
	relsPart := common.SlideRelsPartName(part)

	slideXML, slideRelsXML, err := renderEditorSlideParts(slide, slideNumber, "", e.metadata.SlideSize.Width, e.metadata.SlideSize.Height)
	if err != nil {
		return 0, err
	}

	e.parts[part] = []byte(slideXML)
	e.parts[relsPart] = []byte(slideRelsXML)

	e.slides = append(e.slides, common.EditorSlideRef{
		SlideID: e.nextSlideID,
		RelID:   relID,
		Target:  fmt.Sprintf("slides/slide%d.xml", slideNumber),
		Part:    part,
		Title:   slide.Title,
	})

	e.nextSlideID++
	e.nextRelIDNum++
	e.nextSlideNum++
	e.metadata.SlideCount = len(e.slides)
	return len(e.slides) - 1, nil
}

// UpdateSlide replaces one slide content at index.
func (e *PresentationEditor) UpdateSlide(index int, slide elements.SlideContent) error {
	if e == nil {
		return fmt.Errorf("editor cannot be nil")
	}
	if index < 0 || index >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range [0,%d)", index, len(e.slides))
	}
	if err := validateEditorSlideContent(slide); err != nil {
		return err
	}

	ref := e.slides[index]
	existingRels, err := e.slideRelationships(ref.Part)
	if err != nil {
		return err
	}
	notesTarget, err := scanSupportedSlideRels(existingRels)
	if err != nil {
		return fmt.Errorf("slide %d cannot be updated safely: %w", index, err)
	}

	number, ok := parseSlidePartNumber(ref.Part)
	if !ok {
		return fmt.Errorf("unsupported slide part path %q", ref.Part)
	}
	slideXML, relsXML, err := renderEditorSlideParts(slide, number, notesTarget, e.metadata.SlideSize.Width, e.metadata.SlideSize.Height)
	if err != nil {
		return err
	}

	e.parts[ref.Part] = []byte(slideXML)
	e.parts[common.SlideRelsPartName(ref.Part)] = []byte(relsXML)
	e.slides[index].Title = slide.Title
	return nil
}

// RemoveSlide removes one slide by index.
func (e *PresentationEditor) RemoveSlide(index int) error {
	if e == nil {
		return fmt.Errorf("editor cannot be nil")
	}
	if index < 0 || index >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range [0,%d)", index, len(e.slides))
	}

	ref := e.slides[index]
	delete(e.parts, ref.Part)
	delete(e.parts, common.SlideRelsPartName(ref.Part))

	next := make([]common.EditorSlideRef, 0, len(e.slides)-1)
	next = append(next, e.slides[:index]...)
	next = append(next, e.slides[index+1:]...)
	e.slides = next
	e.metadata.SlideCount = len(e.slides)
	return nil
}

// MergeFromFile appends slides from another PPTX package.
func (e *PresentationEditor) MergeFromFile(filePath string) error {
	other, err := OpenPresentationEditor(filePath)
	if err != nil {
		return err
	}
	return e.MergeFromEditor(other)
}

// MergeFromEditor appends slides from another editor instance.
func (e *PresentationEditor) MergeFromEditor(other *PresentationEditor) error {
	if e == nil || other == nil {
		return fmt.Errorf("editors cannot be nil")
	}

	for idx, slide := range other.slides {
		rels, err := other.slideRelationships(slide.Part)
		if err != nil {
			return err
		}
		if _, err := scanSupportedSlideRels(rels); err != nil {
			return fmt.Errorf("source slide %d is not merge-supported: %w", idx, err)
		}

		slideNumber := e.nextSlideNum
		relID := fmt.Sprintf("rId%d", e.nextRelIDNum)
		part := fmt.Sprintf("ppt/slides/slide%d.xml", slideNumber)
		relsPart := common.SlideRelsPartName(part)

		sourceSlideBytes := other.parts[slide.Part]
		sourceRelsBytes := other.parts[common.SlideRelsPartName(slide.Part)]
		if len(sourceSlideBytes) == 0 || len(sourceRelsBytes) == 0 {
			return fmt.Errorf("source slide %d parts are missing", idx)
		}

		copiedSlide := make([]byte, len(sourceSlideBytes))
		copy(copiedSlide, sourceSlideBytes)
		// TODO: Implement deep-copy for slide relationships to remap rIds and media targets.
		// For now, only simple slides without complex internal rels are safely supported.
		copiedRels := make([]byte, len(sourceRelsBytes))
		copy(copiedRels, sourceRelsBytes)
		e.parts[part] = copiedSlide
		e.parts[relsPart] = copiedRels
		e.slides = append(e.slides, common.EditorSlideRef{
			SlideID: e.nextSlideID,
			RelID:   relID,
			Target:  fmt.Sprintf("slides/slide%d.xml", slideNumber),
			Part:    part,
			Title:   slide.Title,
		})

		e.nextSlideID++
		e.nextRelIDNum++
		e.nextSlideNum++
	}
	e.metadata.SlideCount = len(e.slides)
	return nil
}

func (e *PresentationEditor) slideRelationships(slidePart string) ([]common.EditorRelationship, error) {
	relsPart := common.SlideRelsPartName(slidePart)
	data, ok := e.parts[relsPart]
	if !ok {
		return nil, fmt.Errorf("missing slide relationships part %q", relsPart)
	}
	rels, err := parseRelationshipsXML(data)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", relsPart, err)
	}
	return rels, nil
}

func scanSupportedSlideRels(rels []common.EditorRelationship) (notesTarget string, err error) {
	for _, rel := range rels {
		switch rel.Type {
		case common.RelTypeSlideLayout:
		case common.RelTypeNotesSlide:
			notesTarget = rel.Target
		case common.RelTypeHyperlink:
		default:
			return "", fmt.Errorf("unsupported relationship type %q", rel.Type)
		}
	}
	return notesTarget, nil
}

func validateEditorSlideContent(slide elements.SlideContent) error {
	if slide.Notes != "" {
		return fmt.Errorf("editor add/update does not support notes authoring yet")
	}
	if len(slide.Images) > 0 {
		return fmt.Errorf("editor add/update does not support embedded image authoring yet")
	}
	if slide.ChartKindCount() > 0 {
		return fmt.Errorf("editor add/update does not support chart authoring yet")
	}
	if err := slide.Validate(1); err != nil {
		return err
	}
	return nil
}
