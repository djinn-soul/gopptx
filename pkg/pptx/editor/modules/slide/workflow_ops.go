package slide

import (
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

type AddSlideState struct {
	Slides       []common.EditorSlideRef
	NextSlideNum int
	NextRelIDNum int
	NextSlideID  int64
	SlideCount   int
	SlideWidth   int64
	SlideHeight  int64
}

type RenderNewSlideFn func(slide elements.SlideContent, part string, slideNumber int, width, height int64) (string, string, error)
type RenderExistingSlideFn func(
	slide elements.SlideContent,
	part string,
	number int,
	notesTarget string,
	width, height int64,
) (string, string, error)

func AddSlideToState(
	state AddSlideState,
	slide elements.SlideContent,
	renderSlide RenderNewSlideFn,
	setPart SetPartFn,
) (AddSlideState, int, error) {
	slideNumber := state.NextSlideNum
	relID := fmt.Sprintf("rId%d", state.NextRelIDNum)
	part := fmt.Sprintf("ppt/slides/slide%d.xml", slideNumber)
	relsPart := common.SlideRelsPartName(part)

	slideXML, slideRelsXML, err := renderSlide(slide, part, slideNumber, state.SlideWidth, state.SlideHeight)
	if err != nil {
		return state, 0, err
	}

	setPart(part, []byte(slideXML))
	setPart(relsPart, []byte(slideRelsXML))

	state.Slides = append(state.Slides, common.EditorSlideRef{
		SlideID: state.NextSlideID,
		RelID:   relID,
		Target:  fmt.Sprintf("slides/slide%d.xml", slideNumber),
		Part:    part,
		Title:   slide.Title,
	})
	state.NextSlideID++
	state.NextRelIDNum++
	state.NextSlideNum++
	state.SlideCount = len(state.Slides)
	return state, len(state.Slides) - 1, nil
}

type RemoveSlideState struct {
	Slides         []common.EditorSlideRef
	NotesInventory map[string]string
	SlideCount     int
}

type DeletePartFn func(string)

func RemoveSlideAt(state RemoveSlideState, index int, deletePart DeletePartFn) (RemoveSlideState, error) {
	if index < 0 || index >= len(state.Slides) {
		return state, fmt.Errorf("slide index %d out of range [0,%d)", index, len(state.Slides))
	}

	ref := state.Slides[index]
	if notesPart, ok := state.NotesInventory[ref.Part]; ok {
		deletePart(notesPart)
		deletePart(common.SlideRelsPartName(notesPart))
		delete(state.NotesInventory, ref.Part)
	}
	deletePart(ref.Part)
	deletePart(common.SlideRelsPartName(ref.Part))

	next := make([]common.EditorSlideRef, 0, len(state.Slides)-1)
	next = append(next, state.Slides[:index]...)
	next = append(next, state.Slides[index+1:]...)
	state.Slides = next
	state.SlideCount = len(state.Slides)
	return state, nil
}

func MoveSlideRefs(slides []common.EditorSlideRef, from, to int) ([]common.EditorSlideRef, error) {
	if from < 0 || from >= len(slides) {
		return nil, fmt.Errorf("from index %d out of range [0,%d)", from, len(slides))
	}
	if to < 0 || to >= len(slides) {
		return nil, fmt.Errorf("to index %d out of range [0,%d)", to, len(slides))
	}
	if from == to {
		return slides, nil
	}

	slide := slides[from]
	slides = append(slides[:from], slides[from+1:]...)

	next := make([]common.EditorSlideRef, 0, len(slides)+1)
	next = append(next, slides[:to]...)
	next = append(next, slide)
	next = append(next, slides[to:]...)
	return next, nil
}

type CloneAssetsFn func(srcPart string, sourceRelsBytes []byte, newPart string) ([]byte, error)

type DuplicateState struct {
	Slides       []common.EditorSlideRef
	NextSlideNum int
	NextRelIDNum int
	NextSlideID  int64
	SlideCount   int
}

type CloneSlideRelationshipsFn func(srcRelsBytes []byte, newPart string) ([]byte, error)

func DuplicateSlideInState(
	state DuplicateState,
	srcIndex int,
	destIndex int,
	getPart GetPartFn,
	setPart SetPartFn,
	appendCopySuffix func([]byte) []byte,
	cloneRelationships CloneSlideRelationshipsFn,
) (DuplicateState, int, error) {
	if srcIndex < 0 || srcIndex >= len(state.Slides) {
		return state, 0, fmt.Errorf("source slide index %d out of range [0,%d)", srcIndex, len(state.Slides))
	}
	if destIndex < 0 || destIndex > len(state.Slides) {
		return state, 0, fmt.Errorf("destination slide index %d out of range [0,%d]", destIndex, len(state.Slides))
	}

	srcRef := state.Slides[srcIndex]
	srcPart := srcRef.Part
	srcRelsPart := common.SlideRelsPartName(srcPart)

	slideBytes, ok := getPart(srcPart)
	if !ok {
		return state, 0, fmt.Errorf("source slide part %q missing", srcPart)
	}
	relsBytes, ok := getPart(srcRelsPart)
	if !ok {
		return state, 0, fmt.Errorf("source slide rels part %q missing", srcRelsPart)
	}

	slideNumber := state.NextSlideNum
	relID := fmt.Sprintf("rId%d", state.NextRelIDNum)
	newPart := fmt.Sprintf("ppt/slides/slide%d.xml", slideNumber)
	newRelsPart := common.SlideRelsPartName(newPart)

	newSlideBytes := appendCopySuffix(CloneBytes(slideBytes))
	setPart(newPart, newSlideBytes)

	updatedRelsBytes, err := cloneRelationships(relsBytes, newPart)
	if err != nil {
		return state, 0, fmt.Errorf("clone slide parts: %w", err)
	}
	setPart(newRelsPart, updatedRelsBytes)

	newRef := common.EditorSlideRef{
		SlideID: state.NextSlideID,
		RelID:   relID,
		Target:  fmt.Sprintf("slides/slide%d.xml", slideNumber),
		Part:    newPart,
		Title:   srcRef.Title + " (Copy)",
	}

	state.Slides = append(state.Slides, common.EditorSlideRef{})
	copy(state.Slides[destIndex+1:], state.Slides[destIndex:])
	state.Slides[destIndex] = newRef

	state.NextSlideID++
	state.NextRelIDNum++
	state.NextSlideNum++
	state.SlideCount = len(state.Slides)
	return state, destIndex, nil
}

type MergeState struct {
	Slides       []common.EditorSlideRef
	NextSlideNum int
	NextRelIDNum int
	NextSlideID  int64
	SlideCount   int
}

func MergeSlidesFromSource(
	state MergeState,
	sourceSlides []common.EditorSlideRef,
	getSourcePart GetPartFn,
	cloneAssets CloneAssetsFn,
	setPart SetPartFn,
) (MergeState, error) {
	for idx, slide := range sourceSlides {
		slideNumber := state.NextSlideNum
		relID := fmt.Sprintf("rId%d", state.NextRelIDNum)
		part := fmt.Sprintf("ppt/slides/slide%d.xml", slideNumber)
		relsPart := common.SlideRelsPartName(part)

		sourceSlideBytes, _ := getSourcePart(slide.Part)
		sourceRelsBytes, _ := getSourcePart(common.SlideRelsPartName(slide.Part))
		if len(sourceSlideBytes) == 0 || len(sourceRelsBytes) == 0 {
			return state, fmt.Errorf("source slide %d parts are missing", idx)
		}

		copiedSlide := CloneBytes(sourceSlideBytes)
		copiedRels, err := cloneAssets(slide.Part, sourceRelsBytes, part)
		if err != nil {
			return state, fmt.Errorf("failed to clone slide assets: %w", err)
		}

		setPart(part, copiedSlide)
		setPart(relsPart, copiedRels)
		state.Slides = append(state.Slides, common.EditorSlideRef{
			SlideID: state.NextSlideID,
			RelID:   relID,
			Target:  fmt.Sprintf("slides/slide%d.xml", slideNumber),
			Part:    part,
			Title:   slide.Title,
		})

		state.NextSlideID++
		state.NextRelIDNum++
		state.NextSlideNum++
	}
	state.SlideCount = len(state.Slides)
	return state, nil
}

type UpdateSlideState struct {
	Slides      []common.EditorSlideRef
	SlideWidth  int64
	SlideHeight int64
}

func UpdateSlideInState(
	state UpdateSlideState,
	index int,
	slide elements.SlideContent,
	getPart GetPartFn,
	setPart SetPartFn,
	parseRelationships ParseRelationshipsFn,
	renderSlide RenderExistingSlideFn,
) (UpdateSlideState, error) {
	if index < 0 || index >= len(state.Slides) {
		return state, fmt.Errorf("slide index %d out of range [0,%d)", index, len(state.Slides))
	}

	ref := state.Slides[index]
	existingRels, err := Relationships(ref.Part, getPart, parseRelationships)
	if err != nil {
		return state, err
	}
	notesTarget, err := ScanSupportedSlideRels(existingRels)
	if err != nil {
		return state, fmt.Errorf("slide %d cannot be updated safely: %w", index, err)
	}
	if HasImageContent(slide) && !HasSlideLayoutRelationship(existingRels) {
		return state, fmt.Errorf("slide %d cannot add images without a slideLayout relationship", index)
	}

	number, ok := common.ParseSlidePartNumber(ref.Part)
	if !ok {
		return state, fmt.Errorf("unsupported slide part path %q", ref.Part)
	}
	slideXML, relsXML, err := renderSlide(slide, ref.Part, number, notesTarget, state.SlideWidth, state.SlideHeight)
	if err != nil {
		return state, err
	}

	setPart(ref.Part, []byte(slideXML))
	setPart(common.SlideRelsPartName(ref.Part), []byte(relsXML))
	state.Slides[index].Title = slide.Title
	return state, nil
}
