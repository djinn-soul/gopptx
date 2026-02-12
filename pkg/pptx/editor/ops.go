package editor

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

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

	slideXML, slideRelsXML, err := renderEditorSlideParts(e, slide, slideNumber, "", e.metadata.SlideSize.Width, e.metadata.SlideSize.Height)
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
	slideXML, relsXML, err := renderEditorSlideParts(e, slide, number, notesTarget, e.metadata.SlideSize.Width, e.metadata.SlideSize.Height)
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

// DuplicateSlide clones a slide at srcIndex and inserts it at destIndex.
// All shared assets (layouts, images) are reused in the clone.
func (e *PresentationEditor) DuplicateSlide(srcIndex, destIndex int) (int, error) {
	if e == nil {
		return 0, fmt.Errorf("editor cannot be nil")
	}
	if srcIndex < 0 || srcIndex >= len(e.slides) {
		return 0, fmt.Errorf("source slide index %d out of range [0,%d)", srcIndex, len(e.slides))
	}
	if destIndex < 0 || destIndex > len(e.slides) {
		return 0, fmt.Errorf("destination slide index %d out of range [0,%d]", destIndex, len(e.slides))
	}

	srcRef := e.slides[srcIndex]
	srcPart := srcRef.Part
	srcRelsPart := common.SlideRelsPartName(srcPart)

	slideBytes, ok := e.parts[srcPart]
	if !ok {
		return 0, fmt.Errorf("source slide part %q missing", srcPart)
	}
	relsBytes, ok := e.parts[srcRelsPart]
	if !ok {
		return 0, fmt.Errorf("source slide rels part %q missing", srcRelsPart)
	}

	// Allocate new identifiers
	// Recalculate max rel ID across ALL parts to avoid collisions if they were added out-of-sequence
	e.recalculateNextRelIDNum()

	slideNumber := e.nextSlideNum
	relID := fmt.Sprintf("rId%d", e.nextRelIDNum)
	newPart := fmt.Sprintf("ppt/slides/slide%d.xml", slideNumber)
	newRelsPart := common.SlideRelsPartName(newPart)

	// Clone bytes
	newSlideBytes := make([]byte, len(slideBytes))
	copy(newSlideBytes, slideBytes)
	newRelsBytes := make([]byte, len(relsBytes))
	copy(newRelsBytes, relsBytes)

	// Try to visually indicate copy in the XML
	newSlideBytes = appendCopySuffixToXML(newSlideBytes)
	e.parts[newPart] = newSlideBytes
	e.parts[newRelsPart] = newRelsBytes

	newRef := common.EditorSlideRef{
		SlideID: e.nextSlideID,
		RelID:   relID,
		Target:  fmt.Sprintf("slides/slide%d.xml", slideNumber),
		Part:    newPart,
		Title:   srcRef.Title + " (Copy)",
	}

	// Insert into slides slice
	e.slides = append(e.slides, common.EditorSlideRef{})
	copy(e.slides[destIndex+1:], e.slides[destIndex:])
	e.slides[destIndex] = newRef

	e.nextSlideID++
	e.nextRelIDNum++
	e.nextSlideNum++
	e.metadata.SlideCount = len(e.slides)

	return destIndex, nil
}

// DuplicateSlideAfter clones a slide at index and appends it immediately after.
func (e *PresentationEditor) DuplicateSlideAfter(index int) (int, error) {
	return e.DuplicateSlide(index, index+1)
}

// MoveSlide reorders a slide from one index to another.
func (e *PresentationEditor) MoveSlide(from, to int) error {
	if e == nil {
		return fmt.Errorf("editor cannot be nil")
	}
	if from < 0 || from >= len(e.slides) {
		return fmt.Errorf("from index %d out of range [0,%d)", from, len(e.slides))
	}
	if to < 0 || to >= len(e.slides) {
		return fmt.Errorf("to index %d out of range [0,%d)", to, len(e.slides))
	}
	if from == to {
		return nil
	}

	slide := e.slides[from]
	// Remove from slice
	e.slides = append(e.slides[:from], e.slides[from+1:]...)

	// Insert back at new position
	next := make([]common.EditorSlideRef, 0, len(e.slides)+1)
	next = append(next, e.slides[:to]...)
	next = append(next, slide)
	next = append(next, e.slides[to:]...)
	e.slides = next

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
	// Images are now supported!
	if slide.ChartKindCount() > 0 {
		return fmt.Errorf("editor add/update does not support chart authoring yet")
	}
	if err := slide.Validate(1); err != nil {
		return err
	}
	return nil
}

var aTTitlePattern = regexp.MustCompile(`(?s)<a:t>(.*?)</a:t>`)

func appendCopySuffixToXML(content []byte) []byte {
	modified := false
	res := aTTitlePattern.ReplaceAllFunc(content, func(match []byte) []byte {
		if modified {
			return match
		}
		modified = true
		s := string(match)
		// Insert " (Copy)" before the closing tag
		return []byte(strings.Replace(s, "</a:t>", " (Copy)</a:t>", 1))
	})
	return res
}

func (e *PresentationEditor) recalculateNextRelIDNum() {
	maxNum := 0
	// Scan slide references
	for _, slide := range e.slides {
		if num, ok := parseRelationshipNumber(slide.RelID); ok && num > maxNum {
			maxNum = num
		}
	}
	// Scan non-slide relationships (e.g. presentation.xml.rels)
	for _, rel := range e.nonSlideRels {
		if num, ok := parseRelationshipNumber(rel.ID); ok && num > maxNum {
			maxNum = num
		}
	}
	e.nextRelIDNum = maxNum + 1
}

// SetSlideTitle updates the text of the first title-like element in the slide XML.
func (e *PresentationEditor) SetSlideTitle(index int, title string) error {
	if e == nil {
		return fmt.Errorf("editor cannot be nil")
	}
	if index < 0 || index >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range", index)
	}

	ref := e.slides[index]
	content, ok := e.parts[ref.Part]
	if !ok {
		return fmt.Errorf("slide part %q missing", ref.Part)
	}

	modified := false
	newContent := aTTitlePattern.ReplaceAllFunc(content, func(match []byte) []byte {
		if modified {
			return match
		}
		modified = true
		return []byte("<a:t>" + common.XMLEscape(title) + "</a:t>")
	})

	e.parts[ref.Part] = newContent
	e.slides[index].Title = title
	return nil
}

// RegisterImage adds an image to the presentation or reuses an existing one based on SHA-1 hash.
func (e *PresentationEditor) RegisterImage(data []byte, format string) (string, error) {
	if e == nil {
		return "", fmt.Errorf("editor cannot be nil")
	}
	if len(data) == 0 {
		return "", fmt.Errorf("image data cannot be empty")
	}

	hash := sha1.Sum(data)
	hexHash := hex.EncodeToString(hash[:])

	if path, ok := e.mediaInventory[hexHash]; ok {
		return path, nil
	}

	// New image
	partPath := fmt.Sprintf("ppt/media/image%d.%s", e.nextMediaNum, format)
	e.nextMediaNum++

	e.parts[partPath] = data
	e.mediaInventory[hexHash] = partPath

	return partPath, nil
}

// AddSection creates a new grouped section for slides.
func (e *PresentationEditor) AddSection(name string, slideIndices []int) error {
	if e == nil {
		return fmt.Errorf("editor cannot be nil")
	}
	if name == "" {
		return fmt.Errorf("section name cannot be empty")
	}

	ids := make([]int64, 0, len(slideIndices))
	for _, idx := range slideIndices {
		if idx < 0 || idx >= len(e.slides) {
			return fmt.Errorf("slide index %d out of range", idx)
		}
		ids = append(ids, e.slides[idx].SlideID)
	}

	guid, err := common.NewGUID()
	if err != nil {
		return fmt.Errorf("generate section GUID: %w", err)
	}

	e.sections = append(e.sections, EditorSection{
		Name:     name,
		GUID:     guid,
		SlideIDs: ids,
	})
	return nil
}

// RemoveSection removes a section by name.
func (e *PresentationEditor) RemoveSection(name string) error {
	if e == nil {
		return fmt.Errorf("editor cannot be nil")
	}
	next := make([]EditorSection, 0, len(e.sections))
	found := false
	for _, s := range e.sections {
		if s.Name == name {
			found = true
			continue
		}
		next = append(next, s)
	}
	if !found {
		return fmt.Errorf("section %q not found", name)
	}
	e.sections = next
	return nil
}

// Sections returns the current section list.
func (e *PresentationEditor) Sections() []EditorSection {
	if e == nil {
		return nil
	}
	return e.sections
}
