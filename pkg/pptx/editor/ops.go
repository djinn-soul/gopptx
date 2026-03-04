package editor

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

const notesMasterThemeIndex = 2

// AddSlide appends a new slide and returns its 0-based index.
func (e *PresentationEditor) AddSlide(slide elements.SlideContent) (int, error) {
	if e == nil {
		return 0, errors.New("editor cannot be nil")
	}
	if err := editorslide.ValidateEditorSlideContent(slide); err != nil {
		return 0, err
	}

	slideNumber := e.nextSlideNum
	relID := fmt.Sprintf("rId%d", e.nextRelIDNum)
	part := fmt.Sprintf("ppt/slides/slide%d.xml", slideNumber)
	relsPart := common.SlideRelsPartName(part)

	slideXML, slideRelsXML, err := renderEditorSlideParts(
		e,
		slide,
		part,
		slideNumber,
		"",
		e.metadata.SlideSize.Width,
		e.metadata.SlideSize.Height,
	)
	if err != nil {
		return 0, err
	}

	e.parts.Set(part, []byte(slideXML))
	e.parts.Set(relsPart, []byte(slideRelsXML))

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
		return errors.New("editor cannot be nil")
	}
	if index < 0 || index >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range [0,%d)", index, len(e.slides))
	}
	ref := e.slides[index]
	slide = preserveExistingSlideTransition(e.parts, ref.Part, slide)
	if err := editorslide.ValidateEditorSlideContent(slide); err != nil {
		return err
	}

	existingRels, err := editorslide.SlideRelationships(ref.Part, e.parts.Get, parseRelationshipsXML)
	if err != nil {
		return err
	}
	notesTarget, err := editorslide.ScanSupportedSlideRels(existingRels)
	if err != nil {
		return fmt.Errorf("slide %d cannot be updated safely: %w", index, err)
	}
	if editorslide.SlideHasImageContent(slide) && !editorslide.HasSlideLayoutRelationship(existingRels) {
		return fmt.Errorf("slide %d cannot add images without a slideLayout relationship", index)
	}

	number, ok := common.ParseSlidePartNumber(ref.Part)
	if !ok {
		return fmt.Errorf("unsupported slide part path %q", ref.Part)
	}
	slideXML, relsXML, err := renderEditorSlideParts(
		e,
		slide,
		ref.Part,
		number,
		notesTarget,
		e.metadata.SlideSize.Width,
		e.metadata.SlideSize.Height,
	)
	if err != nil {
		return err
	}

	e.parts.Set(ref.Part, []byte(slideXML))
	e.parts.Set(common.SlideRelsPartName(ref.Part), []byte(relsXML))
	e.slides[index].Title = slide.Title
	return nil
}

// RemoveSlide removes one slide by index.
func (e *PresentationEditor) RemoveSlide(index int) error {
	if e == nil {
		return errors.New("editor cannot be nil")
	}
	if index < 0 || index >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range [0,%d)", index, len(e.slides))
	}

	ref := e.slides[index]
	if notesPart, ok := e.notesInventory[ref.Part]; ok {
		e.parts.Delete(notesPart)
		e.parts.Delete(common.SlideRelsPartName(notesPart))
		delete(e.notesInventory, ref.Part)
	}
	e.parts.Delete(ref.Part)
	e.parts.Delete(common.SlideRelsPartName(ref.Part))

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
		return 0, errors.New("editor cannot be nil")
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

	slideBytes, ok := e.parts.Get(srcPart)
	if !ok {
		return 0, fmt.Errorf("source slide part %q missing", srcPart)
	}
	relsBytes, ok := e.parts.Get(srcRelsPart)
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
	newSlideBytes = editorslide.AppendCopySuffixToXML(newSlideBytes)
	e.parts.Set(newPart, newSlideBytes)

	updatedRelsBytes, err := e.deepCloneSlideParts(srcRef.Part, relsBytes, newPart)
	if err != nil {
		return 0, fmt.Errorf("clone slide parts: %w", err)
	}
	e.parts.Set(newRelsPart, updatedRelsBytes)

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
		return errors.New("editor cannot be nil")
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
	defer func() { _ = other.Close() }()
	return e.MergeFromEditor(other)
}

// MergeFromEditor appends slides from another editor instance.
func (e *PresentationEditor) MergeFromEditor(other *PresentationEditor) error {
	if e == nil || other == nil {
		return errors.New("editors cannot be nil")
	}

	for idx, slide := range other.slides {
		// Check removed to allow deep copier to handle supported types
		// if _, err := scanSupportedSlideRels(rels); err != nil {
		// 	return fmt.Errorf("source slide %d is not merge-supported: %w", idx, err)
		// }

		slideNumber := e.nextSlideNum
		relID := fmt.Sprintf("rId%d", e.nextRelIDNum)
		part := fmt.Sprintf("ppt/slides/slide%d.xml", slideNumber)
		relsPart := common.SlideRelsPartName(part)

		sourceSlideBytes, _ := other.parts.Get(slide.Part)
		sourceRelsBytes, _ := other.parts.Get(common.SlideRelsPartName(slide.Part))
		if len(sourceSlideBytes) == 0 || len(sourceRelsBytes) == 0 {
			return fmt.Errorf("source slide %d parts are missing", idx)
		}

		copiedSlide := make([]byte, len(sourceSlideBytes))
		copy(copiedSlide, sourceSlideBytes)

		// Use deep clone helper to handle assets (images, charts) and remapping
		copiedRels, err := e.deepCloneSlideAssets(other, slide.Part, sourceRelsBytes, part)
		if err != nil {
			return fmt.Errorf("failed to clone slide assets: %w", err)
		}

		e.parts.Set(part, copiedSlide)
		e.parts.Set(relsPart, copiedRels)
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

func (e *PresentationEditor) recalculateNextRelIDNum() {
	e.nextRelIDNum = editorslide.NextRelationshipIDNum(e.slides, e.nonSlideRels)
}

// SetSlideTitle replaces title placeholder text runs with the provided title.
func (e *PresentationEditor) SetSlideTitle(index int, title string) error {
	if e == nil {
		return errors.New("editor cannot be nil")
	}
	if index < 0 || index >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range", index)
	}

	ref := e.slides[index]
	content, ok := e.parts.Get(ref.Part)
	if !ok {
		return fmt.Errorf("slide part %q missing", ref.Part)
	}

	newContent, modified := editorslide.ReplaceAllTitleTextRuns(content, common.XMLEscape(title))
	if !modified {
		return fmt.Errorf("slide %d has no title text run to update", index)
	}

	e.parts.Set(ref.Part, newContent)
	e.slides[index].Title = title
	return nil
}

// RegisterImage adds an image to the presentation or reuses an existing one based on SHA-1 hash.
func (e *PresentationEditor) RegisterImage(data []byte, format string) (string, error) {
	if e == nil {
		return "", errors.New("editor cannot be nil")
	}
	if len(data) == 0 {
		return "", errors.New("image data cannot be empty")
	}

	// Calculate hash outside the lock to avoid blocking other image registrations.
	hash := sha256.Sum256(data)
	hexHash := hex.EncodeToString(hash[:])

	e.mediaMu.Lock()
	defer e.mediaMu.Unlock()

	if part, ok := e.mediaInventory[hexHash]; ok {
		return part, nil
	}

	// New image - use strconv for minor speedup over fmt.Sprintf
	partPath := "ppt/media/image" + strconv.Itoa(e.nextMediaNum) + "." + format
	e.nextMediaNum++

	e.parts.Set(partPath, data)
	e.mediaInventory[hexHash] = partPath

	return partPath, nil
}

// AddSection creates a new grouped section for slides.
func (e *PresentationEditor) AddSection(name string, slideIndices []int) error {
	if e == nil {
		return errors.New("editor cannot be nil")
	}
	ids, err := editorslide.BuildSectionSlideIDs(e.slides, slideIndices)
	if err != nil {
		return err
	}
	next, err := editorslide.AddSectionData(e.sections, name, ids, common.NewGUID)
	if err != nil {
		return err
	}
	e.sections = next
	return nil
}

// RemoveSection removes a section by name.
func (e *PresentationEditor) RemoveSection(name string) error {
	if e == nil {
		return errors.New("editor cannot be nil")
	}
	next, err := editorslide.RemoveSectionData(e.sections, name)
	if err != nil {
		return err
	}
	e.sections = next
	return nil
}

// RenameSection renames a section.
func (e *PresentationEditor) RenameSection(oldName, newName string) error {
	if e == nil {
		return errors.New("editor cannot be nil")
	}
	next, err := editorslide.RenameSectionData(e.sections, oldName, newName)
	if err != nil {
		return err
	}
	e.sections = next
	return nil
}

// Sections returns the current section list.
func (e *PresentationEditor) Sections() []Section {
	if e == nil {
		return nil
	}
	return e.sections
}

func (e *PresentationEditor) deepCloneSlideParts(
	_ string,
	srcSlideRelsBytes []byte,
	newSlidePart string,
) ([]byte, error) {
	rels, err := parseRelationshipsXML(srcSlideRelsBytes)
	if err != nil {
		return nil, err
	}

	changed := false
	for i, rel := range rels {
		newTarget, handled := e.cloneSlideRelationshipPart(rel, newSlidePart)
		if !handled {
			continue
		}
		rels[i].Target = newTarget
		changed = true
	}

	if changed {
		rendered := renderRelationshipsXML(rels)
		return []byte(rendered), nil
	}
	return srcSlideRelsBytes, nil
}

func (e *PresentationEditor) cloneSlideRelationshipPart(
	rel common.EditorRelationship,
	newSlidePart string,
) (string, bool) {
	switch rel.Type {
	case common.RelTypeChart:
		return e.cloneChartPart(rel)
	case common.RelTypeNotesSlide:
		return e.cloneNotesSlidePart(rel, newSlidePart)
	default:
		return "", false
	}
}

func (e *PresentationEditor) cloneChartPart(rel common.EditorRelationship) (string, bool) {
	srcChartPart := common.CanonicalPartPath(path.Join("ppt/slides", rel.Target))
	newChartPart := fmt.Sprintf("ppt/charts/chart%d.xml", e.nextChartNum)
	e.nextChartNum++

	data, chartOK := e.parts.Get(srcChartPart)
	if !chartOK {
		return "../charts/" + path.Base(newChartPart), true
	}

	newChartData := editorslide.CloneBytes(data)
	newChartData = e.cloneChartDependencies(srcChartPart, newChartPart, newChartData)
	e.parts.Set(newChartPart, newChartData)
	return "../charts/" + path.Base(newChartPart), true
}

func (e *PresentationEditor) cloneChartDependencies(srcChartPart, newChartPart string, newChartData []byte) []byte {
	srcChartRelsPath := common.SlideRelsPartName(srcChartPart)
	relsData, relsOK := e.parts.Get(srcChartRelsPath)
	if !relsOK {
		return newChartData
	}

	chartRels, _ := parseRelationshipsXML(relsData)
	for i, cr := range chartRels {
		if cr.Type != common.RelTypePackage {
			continue
		}

		updatedChartData, newExcel, copied := e.cloneChartEmbedding(srcChartPart, newChartData, cr)
		if !copied {
			continue
		}
		newChartData = updatedChartData
		chartRels[i].Target = "../embeddings/" + path.Base(newExcel)
		e.chartEmbeddings[newChartPart] = newExcel
	}

	newChartRelsPath := common.SlideRelsPartName(newChartPart)
	rendered := renderRelationshipsXML(chartRels)
	e.parts.Set(newChartRelsPath, []byte(rendered))
	return newChartData
}

func (e *PresentationEditor) cloneChartEmbedding(
	srcChartPart string,
	newChartData []byte,
	chartRel common.EditorRelationship,
) ([]byte, string, bool) {
	srcExcel := common.CanonicalPartPath(path.Join(path.Dir(srcChartPart), chartRel.Target))
	newExcel := fmt.Sprintf("ppt/embeddings/Microsoft_Excel_Worksheet%d.xlsx", e.nextExcelNum)
	e.nextExcelNum++

	xdata, excelOK := e.parts.Get(srcExcel)
	if !excelOK {
		return newChartData, "", false
	}

	e.parts.Set(newExcel, editorslide.CloneBytes(xdata))
	newChartData = rewriteChartExternalData(newChartData, chartRel.ID)
	return newChartData, newExcel, true
}

func (e *PresentationEditor) cloneNotesSlidePart(
	rel common.EditorRelationship,
	newSlidePart string,
) (string, bool) {
	srcNotesPart := common.CanonicalPartPath(path.Join("ppt/slides", rel.Target))
	newNotesPart := fmt.Sprintf("ppt/notesSlides/notesSlide%d.xml", e.nextNotesNum)
	e.nextNotesNum++

	data, notesOK := e.parts.Get(srcNotesPart)
	if !notesOK {
		return "../notesSlides/" + path.Base(newNotesPart), true
	}

	e.parts.Set(newNotesPart, editorslide.CloneBytes(data))
	e.cloneNotesRelationships(srcNotesPart, newNotesPart, newSlidePart)
	e.notesInventory[newSlidePart] = newNotesPart
	return "../notesSlides/" + path.Base(newNotesPart), true
}

func (e *PresentationEditor) cloneNotesRelationships(srcNotesPart, newNotesPart, newSlidePart string) {
	srcNotesRelsPath := common.SlideRelsPartName(srcNotesPart)
	relsData, relsOK := e.parts.Get(srcNotesRelsPath)
	if !relsOK {
		return
	}

	notesRels, _ := parseRelationshipsXML(relsData)
	for i, nr := range notesRels {
		if nr.Type == common.RelTypeSlide {
			notesRels[i].Target = "../slides/" + path.Base(newSlidePart)
		}
	}

	newNotesRelsPath := common.SlideRelsPartName(newNotesPart)
	rendered := renderRelationshipsXML(notesRels)
	e.parts.Set(newNotesRelsPath, []byte(rendered))
}

// UpdateNotesMaster configures the global notes master for the presentation.
//
//nolint:gocognit // Notes-master update coordinates validation, media registration, and rel wiring in one flow.
func (e *PresentationEditor) UpdateNotesMaster(master *elements.NotesMaster) error {
	if e == nil {
		return errors.New("editor cannot be nil")
	}
	if master != nil {
		if err := master.Validate(); err != nil {
			return err
		}
	}
	e.ensureNotesInfrastructure()
	e.ensureNotesMasterThemePart()

	backgroundRID, mediaNames, err := editorslide.ResolveNotesMasterBackgroundMedia(
		master,
		os.ReadFile,
		e.RegisterImage,
	)
	if err != nil {
		return err
	}

	spec := elements.MapNotesMasterToSpec(master, backgroundRID)
	e.parts.Set("ppt/notesMasters/notesMaster1.xml", []byte(pptxxml.NotesMaster(spec)))
	e.parts.Set(
		"ppt/notesMasters/_rels/notesMaster1.xml.rels",
		[]byte(pptxxml.NotesMasterRelationships(notesMasterThemeIndex, mediaNames)),
	)

	return nil
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

func (e *PresentationEditor) ensureNotesMasterThemePart() {
	editorslide.EnsureNotesMasterThemePart(e.parts.Has, e.parts.Get, e.parts.Set)
}
