package editor

import (
	"errors"
	"fmt"

	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

// insertBatchShapes allocates a contiguous block of count ids on a slide, has
// build render the batch XML starting from that block, and splices the result
// into the slide tree. Shared by the textbox and connector batch inserts, which
// differ only in what they render.
func (e *PresentationEditor) insertBatchShapes(
	slideIndex int,
	count int,
	build func(partPath string, startID int) ([]byte, []int, error),
) ([]int, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, errors.New("slide index out of range")
	}
	if count == 0 {
		return []int{}, nil
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return nil, fmt.Errorf("read slide part %s: not found", partPath)
	}

	maxID := e.maxObjectID(partPath, content)
	e.reserveObjectIDs(partPath, maxID+count)
	shapeXML, shapeIDs, err := build(partPath, maxID)
	if err != nil {
		return nil, err
	}

	newXML, err := insertShapeXML(content, shapeXML)
	if err != nil {
		return nil, err
	}
	e.parts.Set(partPath, newXML)
	return shapeIDs, nil
}

// maxObjectID returns the id to allocate above for partPath.
//
// It is the highest id present in the part, raised to the highest id this
// editor has ever handed out for that part. PowerPoint does not reuse
// <p:cNvPr id> values within a slide, and neither should we: the whole editor
// API is id-addressed, so recycling an id after a removal means a caller
// holding a stale id silently operates on an unrelated shape.
func (e *PresentationEditor) maxObjectID(partPath string, content []byte) int {
	live := editorshape.MaxObjectID(content, cNvPrIDPattern, cNvPrSubmatchSize)

	e.objectIDWatermarkMu.Lock()
	defer e.objectIDWatermarkMu.Unlock()
	if e.objectIDWatermark == nil {
		e.objectIDWatermark = make(map[string]int)
	}
	if mark := e.objectIDWatermark[partPath]; mark > live {
		live = mark
	}
	e.objectIDWatermark[partPath] = live
	return live
}

// rememberObjectIDs records the ids currently present in a part. Removal paths
// call it before mutating so that deleting the highest-numbered shape does not
// hand its id back to the next allocation.
func (e *PresentationEditor) rememberObjectIDs(partPath string, content []byte) {
	_ = e.maxObjectID(partPath, content)
}

// reserveObjectIDs records that ids up to and including highest have been used
// for partPath, so a later allocation starts above them.
func (e *PresentationEditor) reserveObjectIDs(partPath string, highest int) {
	e.objectIDWatermarkMu.Lock()
	defer e.objectIDWatermarkMu.Unlock()
	if e.objectIDWatermark == nil {
		e.objectIDWatermark = make(map[string]int)
	}
	if highest > e.objectIDWatermark[partPath] {
		e.objectIDWatermark[partPath] = highest
	}
}
