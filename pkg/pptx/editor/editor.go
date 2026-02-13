package editor

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"

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
		idx := idx
		wg.Add(1)
		go func() {
			defer wg.Done()
			data, _ := e.parts.Get(e.slides[idx].Part)
			title := extractFirstAText(data)
			ch <- result{index: idx, title: title}
		}()
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
			if err == io.EOF {
				return ""
			}
			return ""
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "t" {
			continue
		}
		var value string
		if err := decoder.DecodeElement(&value, &start); err != nil {
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
		return nil, fmt.Errorf("slide index out of range")
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
			Text: p.Text,
			X:    p.X,
			Y:    p.Y,
			W:    p.W,
			H:    p.H,
		}
	}
	return shapes, nil
}

// UpdateShape modifies the properties of a specific shape on a slide.
// Only fields that are non-zero in the 'updates' struct will be applied (if feasible),
// but for simple replacement logic we might replace the whole shape XML.
// Current implementation replaces the XML structure with a re-rendered version
// based on the parsed state + updates.
func (e *PresentationEditor) UpdateShape(slideIndex, shapeIndex int, updates common.Shape) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return fmt.Errorf("slide index out of range")
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
		return fmt.Errorf("shape index out of range")
	}

	// Apply updates to the target shape
	target := &shapes[shapeIndex]
	if updates.Text != "" {
		target.Text = updates.Text
	}
	if updates.X != 0 {
		target.X = updates.X
	}
	if updates.Y != 0 {
		target.Y = updates.Y
	}
	if updates.W != 0 {
		target.W = updates.W
	}
	if updates.H != 0 {
		target.H = updates.H
	}

	// Re-render only the modified shape
	newContent := replaceShapeNodes(content, shapes, func(i int, p *parsedShape) ([]byte, bool) {
		if i != shapeIndex {
			return nil, false
		}
		// Render the updated shape to XML
		// We need a helper to render a 'parsedShape' back to XML
		// NOT reusing elements.Shape logic because we might be dealing with p:pic or existing styles.
		// Limitation: This implementation recreates the shape structure, potentially losing
		// extra attributes (rotations, effects, custom geometry) if not captured in parsedShape.
		// For the smoke test "Find & Replace" text/pos, this might be acceptable for now,
		// but ideally we should update IN PLACE.
		//
		// Let's implement a 'renderParsedShape' in shape_editor.go
		return renderShapeXML(p), true
	})

	e.parts.Set(partPath, newContent)
	return nil
}

// RemoveShape removes a shape from the slide.
func (e *PresentationEditor) RemoveShape(slideIndex, shapeIndex int) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return fmt.Errorf("slide index out of range")
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
		return fmt.Errorf("shape index out of range")
	}

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
