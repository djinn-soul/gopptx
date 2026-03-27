package editor

import (
	"bytes"
	"errors"
	"fmt"
	"unsafe"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

// shapeCacheEntry caches the parsed shapes for one slide part.
// dataPtr is the backing-array address of the content slice at the time of caching.
// It acts as a staleness token: if PartStore returns a different slice (after a Set),
// the pointer changes and the cache automatically misses — no explicit invalidation needed.
// uintptr is intentional (not unsafe.Pointer) so it does NOT prevent GC of old data.
type shapeCacheEntry struct {
	dataPtr uintptr
	shapes  []common.Shape
}

const (
	shapeTypePicture = "pic"
	minGroupShapes   = 2
	groupShapeTag    = "grpSp"
)

// GetShapes returns a list of shapes found on the specified slide (0-based index).
func (e *PresentationEditor) GetShapes(slideIndex int) ([]common.Shape, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, errors.New("slide index out of range")
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return nil, fmt.Errorf("read slide part %s: not found", partPath)
	}

	// Fast path: return a shallow copy from the shape cache if the slide data
	// hasn't changed since it was last parsed.  The backing-array pointer is a
	// cheap staleness token — PartStore always stores a new slice on Set(), so a
	// different pointer means the content changed and we must re-parse.
	if len(content) > 0 {
		dataPtr := uintptr(unsafe.Pointer(&content[0]))
		if entry, hit := e.shapeCache[partPath]; hit && entry.dataPtr == dataPtr {
			out := make([]common.Shape, len(entry.shapes))
			copy(out, entry.shapes)
			return out, nil
		}
	}

	parsed, err := parseSlideShapes(content)
	if err != nil {
		return nil, fmt.Errorf("parse shapes: %w", err)
	}

	shapes := make([]common.Shape, len(parsed))
	for i, p := range parsed {
		var placeholderIndex *int
		if p.PhType != "" {
			idx := p.PhIndex
			placeholderIndex = &idx
		}
		shapes[i] = common.Shape{
			ID:               p.ID,
			Name:             p.Name,
			Type:             p.Type,
			Text:             p.Text,
			X:                p.X,
			Y:                p.Y,
			W:                p.W,
			H:                p.H,
			PlaceholderIndex: placeholderIndex,
			PlaceholderType:  p.PhType,
			Fill:             p.Fill,
			Line:             p.Line,
			Shadow:           p.Shadow,
			Adjustments:      p.Adjustments,
		}
	}

	// Populate cache for next call on the same unmodified slide.
	if len(content) > 0 {
		if e.shapeCache == nil {
			e.shapeCache = make(map[string]shapeCacheEntry)
		}
		e.shapeCache[partPath] = shapeCacheEntry{
			dataPtr: uintptr(unsafe.Pointer(&content[0])),
			shapes:  shapes,
		}
	}

	out := make([]common.Shape, len(shapes))
	copy(out, shapes)
	return out, nil
}

// RemoveShapeByIndex removes a shape from the slide by its index.
func (e *PresentationEditor) RemoveShapeByIndex(slideIndex, shapeIndex int) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return errors.New("slide index out of range")
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return fmt.Errorf("read slide part %s: not found", partPath)
	}

	shapes, err := parseSlideShapes(content)
	if err != nil {
		return fmt.Errorf("parse shapes: %w", err)
	}

	if shapeIndex < 0 || shapeIndex >= len(shapes) {
		return errors.New("shape index out of range")
	}

	return e.applyShapeRemoval(partPath, content, shapes, shapeIndex)
}

// RemoveShape removes a shape from the slide by its ID.
func (e *PresentationEditor) RemoveShape(slideIndex, shapeID int) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return errors.New("slide index out of range")
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return fmt.Errorf("read slide part %s: not found", partPath)
	}

	shapes, err := parseSlideShapes(content)
	if err != nil {
		return fmt.Errorf("parse shapes: %w", err)
	}

	shapeIndex := -1
	for i, s := range shapes {
		if s.ID == shapeID {
			shapeIndex = i
			break
		}
	}

	if shapeIndex == -1 {
		return fmt.Errorf("shape with ID %d not found", shapeID)
	}

	return e.applyShapeRemoval(partPath, content, shapes, shapeIndex)
}

// ClearShapes removes all shapes from the given slide in a single XML pass.
func (e *PresentationEditor) ClearShapes(slideIndex int) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return errors.New("slide index out of range")
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return fmt.Errorf("read slide part %s: not found", partPath)
	}

	shapes, err := parseSlideShapes(content)
	if err != nil {
		return fmt.Errorf("parse shapes: %w", err)
	}

	if len(shapes) == 0 {
		return nil
	}

	newContent := replaceShapeNodes(content, shapes, func(_ int, _ *parsedShape) ([]byte, bool) {
		return []byte{}, true
	})
	e.parts.Set(partPath, newContent)
	return nil
}

// GroupShapes groups the specified shapes on a slide into a group shape.
// Returns the ID of the created group shape.
func (e *PresentationEditor) GroupShapes(slideIndex int, shapeIDs []int) (int, error) {
	if len(shapeIDs) < minGroupShapes {
		return 0, errors.New("at least 2 shapes are required to form a group")
	}
	return e.AddGroupShape(slideIndex, shapeIDs)
}

// UngroupShapes ungroups a group shape, returning the ID of the first member shape.
func (e *PresentationEditor) UngroupShapes(slideIndex, shapeID int) (int, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return 0, errors.New("slide index out of range")
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return 0, errors.New("read slide part: not found")
	}

	shapes, err := parseSlideShapes(content)
	if err != nil {
		return 0, fmt.Errorf("parse shapes: %w", err)
	}

	groupShapeIndex := -1
	for i, s := range shapes {
		if s.ID == shapeID && s.IsGroup {
			groupShapeIndex = i
			break
		}
	}

	if groupShapeIndex == -1 {
		return 0, fmt.Errorf("group shape with ID %d not found", shapeID)
	}

	groupShape := shapes[groupShapeIndex]
	groupXML := content[groupShape.Start:groupShape.End]
	children, err := editorshape.ExtractGroupChildShapeNodes(groupXML, groupShapeTag, shapeTypePicture)
	if err != nil {
		return 0, fmt.Errorf("extract grouped shapes: %w", err)
	}
	if len(children) == 0 {
		return 0, fmt.Errorf("group shape with ID %d has no child shapes", shapeID)
	}

	replacement := bytes.Join(children, nil)
	newContent := replaceShapeNodes(content, shapes, func(i int, _ *parsedShape) ([]byte, bool) {
		if i == groupShapeIndex {
			return replacement, true
		}
		return nil, false
	})
	e.parts.Set(partPath, newContent)

	firstChildID := editorshape.FirstShapeIDInXML(children, cNvPrIDPattern, cNvPrSubmatchSize)
	if firstChildID == 0 {
		return shapeID, nil
	}
	return firstChildID, nil
}
