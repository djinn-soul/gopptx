package editor

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// AddSlide appends a new slide and returns its 0-based index.
func (e *PresentationEditor) AddSlide(slide elements.SlideContent) (int, error) {
	if e == nil {
		return 0, errors.New("editor cannot be nil")
	}
	if err := validateEditorSlideContent(slide); err != nil {
		return 0, err
	}

	slideNumber := e.nextSlideNum
	relID := fmt.Sprintf("rId%d", e.nextRelIDNum)
	part := fmt.Sprintf("ppt/slides/slide%d.xml", slideNumber)
	relsPart := common.SlideRelsPartName(part)

	slideXML, slideRelsXML, err := renderEditorSlideParts(
		e,
		slide,
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
	if slideHasImageContent(slide) && !hasSlideLayoutRelationship(existingRels) {
		return fmt.Errorf("slide %d cannot add images without a slideLayout relationship", index)
	}

	number, ok := parseSlidePartNumber(ref.Part)
	if !ok {
		return fmt.Errorf("unsupported slide part path %q", ref.Part)
	}
	slideXML, relsXML, err := renderEditorSlideParts(
		e,
		slide,
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
	newSlideBytes = appendCopySuffixToXML(newSlideBytes)
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

func (e *PresentationEditor) slideRelationships(slidePart string) ([]common.EditorRelationship, error) {
	relsPart := common.SlideRelsPartName(slidePart)
	data, ok := e.parts.Get(relsPart)
	if !ok {
		return nil, fmt.Errorf("missing slide relationships part %q", relsPart)
	}
	rels, err := parseRelationshipsXML(data)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", relsPart, err)
	}
	return rels, nil
}

func scanSupportedSlideRels(rels []common.EditorRelationship) (string, error) {
	notesTarget := ""
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

func hasSlideLayoutRelationship(rels []common.EditorRelationship) bool {
	for _, rel := range rels {
		if rel.Type == common.RelTypeSlideLayout {
			return true
		}
	}
	return false
}

func slideHasImageContent(slide elements.SlideContent) bool {
	if len(slide.Images) > 0 {
		return true
	}
	if slide.Background != nil && slide.Background.Type == elements.SlideBackgroundPicture &&
		slide.Background.PictureFill != nil {
		return true
	}
	for _, override := range slide.PlaceholderOverrides {
		if override.Image != nil {
			return true
		}
	}
	return false
}

func validateEditorSlideContent(slide elements.SlideContent) error {
	// Notes and Images/Charts are now supported!
	if err := slide.Validate(1); err != nil {
		return err
	}
	return nil
}

var aTTitlePattern = regexp.MustCompile(`(?s)<a:t>(.*?)</a:t>`)

func appendCopySuffixToXML(content []byte) []byte {
	res, _ := replaceTitleLikeText(content, replaceLastTextRun, func(match []byte) []byte {
		// Insert " (Copy)" before the closing tag.
		return []byte(strings.Replace(string(match), "</a:t>", " (Copy)</a:t>", 1))
	})
	return res
}

func replaceTitleLikeText(
	content []byte,
	runSelector func([]byte, func([]byte) []byte) ([]byte, bool),
	replaceFn func(match []byte) []byte,
) ([]byte, bool) {
	// Prefer replacing text inside the title placeholder shape when present.
	updated, ok := replaceTitlePlaceholderText(content, runSelector, replaceFn)
	if ok {
		return updated, true
	}
	// Fallback to first <a:t> to preserve behavior on unusual slide XML.
	return runSelector(content, replaceFn)
}

func replaceTitlePlaceholderText(
	content []byte,
	runSelector func([]byte, func([]byte) []byte) ([]byte, bool),
	replaceFn func(match []byte) []byte,
) ([]byte, bool) {
	const (
		shapeStart = "<p:sp"
		shapeEnd   = "</p:sp>"
	)

	searchFrom := 0
	for {
		startIdx := bytes.Index(content[searchFrom:], []byte(shapeStart))
		if startIdx < 0 {
			return content, false
		}
		start := searchFrom + startIdx

		endIdx := bytes.Index(content[start:], []byte(shapeEnd))
		if endIdx < 0 {
			return content, false
		}
		end := start + endIdx + len(shapeEnd)

		shape := content[start:end]
		if isTitlePlaceholderShape(shape) {
			replacedShape, replaced := runSelector(shape, replaceFn)
			if !replaced {
				return content, false
			}

			out := make([]byte, 0, len(content)-len(shape)+len(replacedShape))
			out = append(out, content[:start]...)
			out = append(out, replacedShape...)
			out = append(out, content[end:]...)
			return out, true
		}

		searchFrom = end
	}
}

func isTitlePlaceholderShape(shape []byte) bool {
	if !bytes.Contains(shape, []byte("<p:ph")) {
		return false
	}
	return bytes.Contains(shape, []byte(`type="title"`)) || bytes.Contains(shape, []byte(`type="ctrTitle"`))
}

func replaceFirstTextRun(content []byte, replaceFn func(match []byte) []byte) ([]byte, bool) {
	modified := false
	res := aTTitlePattern.ReplaceAllFunc(content, func(match []byte) []byte {
		if modified {
			return match
		}
		modified = true
		return replaceFn(match)
	})
	return res, modified
}

func replaceLastTextRun(content []byte, replaceFn func(match []byte) []byte) ([]byte, bool) {
	indexes := aTTitlePattern.FindAllIndex(content, -1)
	if len(indexes) == 0 {
		return content, false
	}
	last := indexes[len(indexes)-1]
	start, end := last[0], last[1]
	match := content[start:end]
	replacement := replaceFn(match)

	out := make([]byte, 0, len(content)-len(match)+len(replacement))
	out = append(out, content[:start]...)
	out = append(out, replacement...)
	out = append(out, content[end:]...)
	return out, true
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

	newContent, modified := replaceTitleLikeText(content, replaceFirstTextRun, func(_ []byte) []byte {
		return []byte("<a:t>" + common.XMLEscape(title) + "</a:t>")
	})
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

	e.mediaMu.Lock()
	defer e.mediaMu.Unlock()

	hash := sha256.Sum256(data)
	hexHash := hex.EncodeToString(hash[:])

	if path, ok := e.mediaInventory[hexHash]; ok {
		return path, nil
	}

	// New image
	partPath := fmt.Sprintf("ppt/media/image%d.%s", e.nextMediaNum, format)
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
	if name == "" {
		return errors.New("section name cannot be empty")
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

	e.sections = append(e.sections, Section{
		Name:     name,
		GUID:     guid,
		SlideIDs: ids,
	})
	return nil
}

// RemoveSection removes a section by name.
func (e *PresentationEditor) RemoveSection(name string) error {
	if e == nil {
		return errors.New("editor cannot be nil")
	}
	next := make([]Section, 0, len(e.sections))
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

// RenameSection renames a section.
func (e *PresentationEditor) RenameSection(oldName, newName string) error {
	if e == nil {
		return errors.New("editor cannot be nil")
	}
	if newName == "" {
		return errors.New("new section name cannot be empty")
	}
	found := false
	for i := range e.sections {
		if e.sections[i].Name == oldName {
			e.sections[i].Name = newName
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("section %q not found", oldName)
	}
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

	newChartData := cloneBytes(data)
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

	e.parts.Set(newExcel, cloneBytes(xdata))
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

	e.parts.Set(newNotesPart, cloneBytes(data))
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

func cloneBytes(b []byte) []byte {
	if b == nil {
		return nil
	}
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func (e *PresentationEditor) ensureNotesInfrastructure() {
	if e.parts.Has("ppt/notesMasters/notesMaster1.xml") {
		return
	}
	e.parts.Set("ppt/notesMasters/notesMaster1.xml", []byte(pptxxml.NotesMaster(nil)))
	e.parts.Set("ppt/notesMasters/_rels/notesMaster1.xml.rels", []byte(pptxxml.NotesMasterRelationships()))
	if !e.parts.Has("ppt/theme/theme2.xml") {
		e.parts.Set("ppt/theme/theme2.xml", []byte(pptxxml.Theme(nil)))
	}

	for _, rel := range e.nonSlideRels {
		if rel.Type == common.RelTypeNotesMaster {
			return
		}
	}
	e.recalculateNextRelIDNum()
	e.nonSlideRels = append(e.nonSlideRels, common.EditorRelationship{
		ID:     fmt.Sprintf("rId%d", e.nextRelIDNum),
		Type:   common.RelTypeNotesMaster,
		Target: "notesMasters/notesMaster1.xml",
	})
	e.nextRelIDNum++
}
