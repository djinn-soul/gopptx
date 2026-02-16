package editor

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/djinn-soul/gopptx/pkg/pptx/comments"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// EditorSection describes a PowerPoint section entry.
type EditorSection struct {
	Name     string
	GUID     string
	SlideIDs []int64
}

// PresentationEditor provides read/modify/save operations for existing PPTX files.
type PresentationEditor struct {
	parts *PartStore

	slides       []common.EditorSlideRef
	nextSlideID  int64
	nextRelIDNum int
	nextSlideNum int

	metadata        common.PresentationMetadata
	nonSlideRels    []common.EditorRelationship
	presentationXML string

	// Media inventory for deduplication (SHA1 -> PartPath)
	mediaInventory map[string]string
	nextMediaNum   int
	mediaMu        sync.Mutex
	imagePathCache map[string]imagePathCacheEntry
	imagePathMu    sync.RWMutex

	// Section management
	sections []EditorSection

	// Chart inventory (ChartPath -> EmbeddingPath)
	chartEmbeddings map[string]string
	nextChartNum    int
	nextExcelNum    int

	// Notes inventory (SlidePath -> NotesSlidePath)
	notesInventory map[string]string
	nextNotesNum   int

	// Comment authors
	authorCache   map[int64]comments.Author
	nextAuthorID  int64
	authorCacheMu sync.RWMutex
}

// Metadata returns presentation-level metadata parsed from the package.
func (e *PresentationEditor) Metadata() common.PresentationMetadata {
	return e.metadata
}

// Close releases any resources held by the editor (e.g. the underlying file handle).
func (e *PresentationEditor) Close() error {
	if e == nil || e.parts == nil {
		return nil
	}
	return e.parts.Close()
}

// SlideCount returns the number of slides currently tracked by the editor.
func (e *PresentationEditor) SlideCount() int {
	if e == nil {
		return 0
	}
	return len(e.slides)
}

// Slides returns ordered slide metadata snapshots (0-based indexes).
func (e *PresentationEditor) Slides() []common.SlideMetadata {
	if e == nil || len(e.slides) == 0 {
		return nil
	}
	out := make([]common.SlideMetadata, 0, len(e.slides))
	for idx, slide := range e.slides {
		out = append(out, common.SlideMetadata{
			Index:          idx,
			SlideID:        slide.SlideID,
			RelationshipID: slide.RelID,
			PartName:       slide.Part,
			Title:          slide.Title,
		})
	}
	return out
}

func nextSlideID(slides []common.EditorSlideRef) int64 {
	var maxID int64 = 255
	for _, slide := range slides {
		if slide.SlideID > maxID {
			maxID = slide.SlideID
		}
	}
	return maxID + 1
}

func nextRelationshipNumber(rels []common.EditorRelationship) int {
	maxNum := 0
	for _, rel := range rels {
		num, ok := parseRelationshipNumber(rel.ID)
		if ok && num > maxNum {
			maxNum = num
		}
	}
	return maxNum + 1
}

func nextSlidePartNumber(slides []common.EditorSlideRef) int {
	maxNum := 0
	for _, slide := range slides {
		num, ok := parseSlidePartNumber(slide.Part)
		if ok && num > maxNum {
			maxNum = num
		}
	}
	return maxNum + 1
}

func parseRelationshipNumber(id string) (int, bool) {
	trimmed := strings.TrimSpace(id)
	if !strings.HasPrefix(trimmed, "rId") {
		return 0, false
	}
	num, err := strconv.Atoi(strings.TrimPrefix(trimmed, "rId"))
	if err != nil || num <= 0 {
		return 0, false
	}
	return num, true
}

func parseSlidePartNumber(partPath string) (int, bool) {
	base := path.Base(strings.TrimSpace(partPath))
	if !strings.HasPrefix(base, "slide") || !strings.HasSuffix(base, ".xml") {
		return 0, false
	}
	num, err := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(base, "slide"), ".xml"))
	if err != nil || num <= 0 {
		return 0, false
	}
	return num, true
}

func (e *PresentationEditor) populateSlideTitlesConcurrently() {
	if e == nil || len(e.slides) == 0 {
		return
	}

	type result struct {
		index int
		title string
	}
	ch := make(chan result, len(e.slides))
	var wg sync.WaitGroup

	for idx := range e.slides {
		wg.Go(func() {
			data, _ := e.parts.Get(e.slides[idx].Part)
			title := extractFirstAText(data)
			ch <- result{index: idx, title: title}
		})
	}
	wg.Wait()
	close(ch)

	results := make([]result, 0, len(e.slides))
	for item := range ch {
		results = append(results, item)
	}
	sort.Slice(results, func(i, j int) bool { return results[i].index < results[j].index })
	for _, item := range results {
		e.slides[item.index].Title = item.title
	}
}

func extractFirstAText(content []byte) string {
	decoder := xml.NewDecoder(bytes.NewReader(content))
	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return ""
			}
			return ""
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "t" {
			continue
		}
		var value string
		if decodeErr := decoder.DecodeElement(&value, &start); decodeErr != nil {
			return ""
		}
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
}

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

	parsed, err := parseSlideShapes(content)
	if err != nil {
		return nil, fmt.Errorf("parse shapes: %w", err)
	}

	shapes := make([]common.Shape, len(parsed))
	for i, p := range parsed {
		shapes[i] = common.Shape{
			ID:   p.ID,
			Name: p.Name,
			Type: p.Type,
			Text: p.Text,
			X:    p.X,
			Y:    p.Y,
			W:    p.W,
			H:    p.H,
		}
	}
	return shapes, nil
}

// UpdateShapeByIndex modifies the properties of a specific shape on a slide by its index in the parsed shape list.
func (e *PresentationEditor) UpdateShapeByIndex(slideIndex, shapeIndex int, updates common.ShapeUpdate) error {
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

	return e.applyShapeUpdate(partPath, content, shapes, shapeIndex, updates)
}

// UpdateShape modifies the properties of a specific shape on a slide by its ID.
func (e *PresentationEditor) UpdateShape(slideIndex, shapeID int, updates common.ShapeUpdate) error {
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

	return e.applyShapeUpdate(partPath, content, shapes, shapeIndex, updates)
}

func (e *PresentationEditor) applyShapeUpdate(
	partPath string,
	content []byte,
	shapes []parsedShape,
	shapeIndex int,
	updates common.ShapeUpdate,
) error {
	target := &shapes[shapeIndex]
	if updates.Text != nil {
		target.Text = *updates.Text
	}
	if updates.X != nil {
		target.X = *updates.X
	}
	if updates.Y != nil {
		target.Y = *updates.Y
	}
	if updates.W != nil {
		target.W = *updates.W
	}
	if updates.H != nil {
		target.H = *updates.H
	}

	// Re-render only the modified shape
	newContent := replaceShapeNodes(content, shapes, func(i int, p *parsedShape) ([]byte, bool) {
		if i != shapeIndex {
			return nil, false
		}
		xml := renderShapeXML(p)
		if xml == nil {
			// If render returns nil (e.g. for p:pic), we must NOT replace, otherwise we delete it.
			// But for ID-based updates, we'll return an error in the caller context if we can.
			return nil, false
		}
		return xml, true
	})

	// Check if the shape was actually updated
	if bytes.Equal(content, newContent) && updates != (common.ShapeUpdate{}) {
		// This is a bit of a heuristic, but if renderShapeXML returns nil for the target,
		// newContent will equal content (replace=false).
		updatedShape := &shapes[shapeIndex]
		if updatedShape.Type == "pic" {
			return errors.New("updating shape of type 'pic' is not supported")
		}
	}

	e.parts.Set(partPath, newContent)
	return nil
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

func (e *PresentationEditor) applyShapeRemoval(
	partPath string,
	content []byte,
	shapes []parsedShape,
	shapeIndex int,
) error {
	// Replace with empty byte slice
	newContent := replaceShapeNodes(content, shapes, func(i int, p *parsedShape) ([]byte, bool) {
		if i == shapeIndex {
			return []byte{}, true
		}
		return nil, false
	})

	e.parts.Set(partPath, newContent)
	return nil
}
