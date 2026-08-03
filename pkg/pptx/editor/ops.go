package editor

import (
	"errors"
	"fmt"

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
	state := editorslide.AddSlideState{
		Slides:       e.slides,
		NextSlideNum: e.nextSlideNum,
		NextRelIDNum: e.nextRelIDNum,
		NextSlideID:  e.nextSlideID,
		SlideCount:   e.metadata.SlideCount,
		SlideWidth:   e.metadata.SlideSize.Width,
		SlideHeight:  e.metadata.SlideSize.Height,
	}
	next, index, err := editorslide.AddSlideToState(
		state,
		slide,
		func(content elements.SlideContent, part string, slideNumber int, width, height int64) (string, string, error) {
			return renderEditorSlideParts(e, content, part, slideNumber, "", width, height)
		},
		e.parts.Set,
	)
	if err != nil {
		return 0, err
	}
	e.slides = next.Slides
	e.nextSlideNum = next.NextSlideNum
	e.nextRelIDNum = next.NextRelIDNum
	e.nextSlideID = next.NextSlideID
	e.metadata.SlideCount = next.SlideCount
	return index, nil
}

// UpdateSlide replaces one slide content at index.
func (e *PresentationEditor) UpdateSlide(index int, slide elements.SlideContent) error {
	if e == nil {
		return errors.New("editor cannot be nil")
	}
	if index >= 0 && index < len(e.slides) {
		slide = editorslide.PreserveExistingSlideTransition(e.parts.Get, e.slides[index].Part, slide)
	}
	if err := editorslide.ValidateEditorSlideContent(slide); err != nil {
		return err
	}
	next, err := editorslide.UpdateSlideInState(
		editorslide.UpdateSlideState{
			Slides:      e.slides,
			SlideWidth:  e.metadata.SlideSize.Width,
			SlideHeight: e.metadata.SlideSize.Height,
		},
		index,
		slide,
		e.parts.Get,
		e.parts.Set,
		parseRelationshipsXML,
		func(content elements.SlideContent, part string, number int, notesTarget string, width, height int64) (string, string, error) {
			return renderEditorSlideParts(e, content, part, number, notesTarget, width, height)
		},
	)
	if err != nil {
		return err
	}
	e.slides = next.Slides
	return nil
}

// RemoveSlide removes one slide by index.
func (e *PresentationEditor) RemoveSlide(index int) error {
	if e == nil {
		return errors.New("editor cannot be nil")
	}
	next, err := editorslide.RemoveSlideAt(editorslide.RemoveSlideState{
		Slides:         e.slides,
		NotesInventory: e.notesInventory,
		SlideCount:     e.metadata.SlideCount,
	}, index, e.parts.Delete)
	if err != nil {
		return err
	}
	e.slides = next.Slides
	e.notesInventory = next.NotesInventory
	e.metadata.SlideCount = next.SlideCount
	return nil
}

// DuplicateSlide clones a slide at srcIndex and inserts it at destIndex.
// All shared assets (layouts, images) are reused in the clone.
func (e *PresentationEditor) DuplicateSlide(srcIndex, destIndex int) (int, error) {
	if e == nil {
		return 0, errors.New("editor cannot be nil")
	}
	e.recalculateNextRelIDNum()
	cloneState := editorslide.CloneRelationshipState{
		NextChartNum:    e.nextChartNum,
		NextExcelNum:    e.nextExcelNum,
		NextNotesNum:    e.nextNotesNum,
		ChartEmbeddings: e.chartEmbeddings,
		NotesInventory:  e.notesInventory,
	}
	next, index, err := editorslide.DuplicateSlideInState(
		editorslide.DuplicateState{
			Slides:       e.slides,
			NextSlideNum: e.nextSlideNum,
			NextRelIDNum: e.nextRelIDNum,
			NextSlideID:  e.nextSlideID,
			SlideCount:   e.metadata.SlideCount,
		},
		srcIndex,
		destIndex,
		e.parts.Get,
		e.parts.Set,
		editorslide.AppendCopySuffixToXML,
		func(srcRelsBytes []byte, newPart string) ([]byte, error) {
			updatedRelsBytes, updatedState, err := editorslide.DeepCloneSlideRelationships(
				srcRelsBytes,
				newPart,
				cloneState,
				e.parts.Get,
				e.parts.Set,
				parseRelationshipsXML,
				renderRelationshipsXML,
				rewriteChartExternalData,
			)
			if err != nil {
				return nil, err
			}
			cloneState = updatedState
			return updatedRelsBytes, nil
		},
	)
	if err != nil {
		return 0, err
	}
	e.nextChartNum = cloneState.NextChartNum
	e.nextExcelNum = cloneState.NextExcelNum
	e.nextNotesNum = cloneState.NextNotesNum
	e.slides = next.Slides
	e.nextSlideNum = next.NextSlideNum
	e.nextRelIDNum = next.NextRelIDNum
	e.nextSlideID = next.NextSlideID
	e.metadata.SlideCount = next.SlideCount
	return index, nil
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
	next, err := editorslide.MoveSlideRefs(e.slides, from, to)
	if err != nil {
		return err
	}
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

// CopySlidesFromFile appends the selected slides of another PPTX package.
// Passing no indices copies every slide, which is what MergeFromFile does.
func (e *PresentationEditor) CopySlidesFromFile(filePath string, indices []int) error {
	other, err := OpenPresentationEditor(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = other.Close() }()
	return e.CopySlidesFromEditor(other, indices)
}

// CopySlidesFromEditor appends the selected slides of another editor instance,
// in the order the indices are given, so one slide can be lifted out of a deck
// without merging the whole thing. Passing no indices copies every slide.
func (e *PresentationEditor) CopySlidesFromEditor(other *PresentationEditor, indices []int) error {
	if err := editorslide.ValidateMergeEditorsNil(e == nil, other == nil); err != nil {
		return err
	}
	selected := other.slides
	if len(indices) > 0 {
		selected = make([]common.EditorSlideRef, 0, len(indices))
		for _, index := range indices {
			if index < 0 || index >= len(other.slides) {
				return fmt.Errorf("source slide index %d out of range [0,%d)", index, len(other.slides))
			}
			selected = append(selected, other.slides[index])
		}
	}
	return e.mergeSlideRefs(other, selected)
}

// MergeFromEditor appends slides from another editor instance.
func (e *PresentationEditor) MergeFromEditor(other *PresentationEditor) error {
	if err := editorslide.ValidateMergeEditorsNil(e == nil, other == nil); err != nil {
		return err
	}
	return e.mergeSlideRefs(other, other.slides)
}

// mergeSlideRefs copies the given source slides, with their assets, into this
// presentation.
func (e *PresentationEditor) mergeSlideRefs(
	other *PresentationEditor,
	sourceSlides []common.EditorSlideRef,
) error {
	// Done first, so every relationship id the import consumes is allocated
	// before the merge reserves its own.
	layoutImports, err := e.importSlideLayouts(other, sourceSlides)
	if err != nil {
		return err
	}
	next, err := editorslide.MergeSlidesFromSource(
		editorslide.MergeState{
			Slides:       e.slides,
			NextSlideNum: e.nextSlideNum,
			NextRelIDNum: e.nextRelIDNum,
			NextSlideID:  e.nextSlideID,
			SlideCount:   e.metadata.SlideCount,
		},
		sourceSlides,
		other.parts.Get,
		func(srcPart string, sourceRelsBytes []byte, newPart string) ([]byte, error) {
			return e.deepCloneSlideAssets(other, srcPart, sourceRelsBytes, newPart, layoutImports)
		},
		e.parts.Set,
	)
	if err != nil {
		return err
	}
	e.slides = next.Slides
	e.nextSlideNum = next.NextSlideNum
	e.nextRelIDNum = next.NextRelIDNum
	e.nextSlideID = next.NextSlideID
	e.metadata.SlideCount = next.SlideCount
	return nil
}

func (e *PresentationEditor) recalculateNextRelIDNum() {
	e.nextRelIDNum = editorslide.NextRelationshipIDNum(e.slides, e.nonSlideRels)
}

// SetSlideHidden marks or unmarks a slide as hidden (show="0" on p:sld root).
func (e *PresentationEditor) SetSlideHidden(index int, hidden bool) error {
	if e == nil {
		return errors.New("editor cannot be nil")
	}
	if index < 0 || index >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range", index)
	}
	e.slides[index].Hidden = hidden
	return nil
}

func (e *PresentationEditor) SetSlideTitle(index int, title string) error {
	if e == nil {
		return errors.New("editor cannot be nil")
	}
	next, err := editorslide.SetSlideTitleInState(e.slides, index, title, e.parts.Get, e.parts.Set)
	if err != nil {
		return err
	}
	e.slides = next
	return nil
}
