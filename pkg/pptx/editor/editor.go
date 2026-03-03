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
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/validation/logical"
	"github.com/djinn-soul/gopptx/pkg/pptx/validation/structural"
)

const (
	shapeTypePicture = "pic"
	minGroupShapes   = 2
	groupShapeTag    = "grpSp"
)

// Section describes a PowerPoint section entry.
type Section struct {
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

	metadata        common.Metadata
	nonSlideRels    []common.EditorRelationship
	presentationXML string
	embeddedFontLst string

	// Media inventory for deduplication (SHA1 -> PartPath)
	mediaInventory map[string]string
	nextMediaNum   int
	mediaMu        sync.Mutex
	imagePathCache map[string]imagePathCacheEntry
	imagePathMu    sync.RWMutex

	// Section management
	sections []Section

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

	// cleanupOnClose is an optional function called during Close().
	cleanupOnClose func()
}

// Metadata returns a pointer to the presentation-level metadata.
func (e *PresentationEditor) Metadata() *common.Metadata {
	return &e.metadata
}

// Close releases any resources held by the editor (e.g. the underlying file handle).
func (e *PresentationEditor) Close() error {
	if e == nil || e.parts == nil {
		return nil
	}
	err := e.parts.Close()
	if e.cleanupOnClose != nil {
		e.cleanupOnClose()
		e.cleanupOnClose = nil
	}
	return err
}

// SetCleanupOnClose registers a function to be called when the editor is closed.
func (e *PresentationEditor) SetCleanupOnClose(fn func()) {
	if e != nil {
		e.cleanupOnClose = fn
	}
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

// Validate performs a structural validation check on the underlying parts.
func (e *PresentationEditor) Validate() []structural.Issue {
	if e == nil || e.parts == nil {
		return nil
	}
	v := structural.NewValidator(e.parts)
	v.AddChecker(&logical.Checker{})
	return v.Validate()
}

// Repair attempts to automatically fix structural issues in the presentation.
func (e *PresentationEditor) Repair() (structural.RepairResult, error) {
	if e == nil || e.parts == nil {
		return structural.RepairResult{}, errors.New("nil editor or parts")
	}
	issues := e.Validate()
	if len(issues) == 0 {
		return structural.RepairResult{}, nil
	}
	r := structural.NewRepairer(e.parts)
	return r.Repair(issues), nil
}

// nextSlideRelID returns the next available relationship ID for a specific slide part.
func (e *PresentationEditor) nextSlideRelID(partPath string) (string, error) {
	relsPath := common.SlideRelsPartName(partPath)
	var rels []common.EditorRelationship
	if data, ok := e.parts.Get(relsPath); ok {
		parsed, err := parseRelationshipsXML(data)
		if err != nil {
			return "", fmt.Errorf("parse %s: %w", relsPath, err)
		}
		rels = parsed
	}

	nextNum := common.NextRelationshipNumber(rels)
	return fmt.Sprintf("rId%d", nextNum), nil
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
			// Note: Ignoring Get errors is intentional - slide titles are optional metadata
			// and we'll just skip titles for slides that can't be read
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
	if len(content) == 0 {
		return ""
	}
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
	newContent := replaceShapeNodes(content, shapes, func(i int, _ *parsedShape) ([]byte, bool) {
		if i == shapeIndex {
			return []byte{}, true
		}
		return nil, false
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

	// Find the group shape
	var groupShapeIndex int
	groupShapeIndex = -1
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
	children, err := extractGroupChildShapeNodes(groupXML)
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

	firstChildID := firstShapeIDInXML(children)
	if firstChildID == 0 {
		return shapeID, nil
	}
	return firstChildID, nil
}

func extractGroupChildShapeNodes(groupXML []byte) ([][]byte, error) {
	decoder := xml.NewDecoder(bytes.NewReader(groupXML))
	rootDepth := 0
	rootSeen := false
	children := make([][]byte, 0)

	for {
		startOffset := decoder.InputOffset()
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		switch t := token.(type) {
		case xml.StartElement:
			nextDepth, sawRoot, childNode, startErr := processGroupStartElement(
				decoder,
				groupXML,
				startOffset,
				rootDepth,
				rootSeen,
				t,
			)
			if startErr != nil {
				return nil, startErr
			}
			rootDepth = nextDepth
			rootSeen = sawRoot
			if len(childNode) > 0 {
				children = append(children, childNode)
			}
		case xml.EndElement:
			var done bool
			rootDepth, done = processGroupEndElement(rootDepth, rootSeen)
			if done {
				return children, nil
			}
		}
	}

	return children, nil
}

func processGroupStartElement(
	decoder *xml.Decoder,
	groupXML []byte,
	startOffset int64,
	rootDepth int,
	rootSeen bool,
	start xml.StartElement,
) (int, bool, []byte, error) {
	if !rootSeen {
		if start.Name.Local == groupShapeTag {
			return 1, true, nil, nil
		}
		return rootDepth, false, nil, nil
	}
	if rootDepth == 1 && isShapeElementLocal(start.Name.Local) {
		_, endOffset, err := extractShapeNode(groupXML, startOffset, decoder, start.Name.Local, true)
		if err != nil {
			return rootDepth, rootSeen, nil, err
		}
		child := bytes.TrimSpace(groupXML[startOffset:endOffset])
		return rootDepth, rootSeen, child, nil
	}
	return rootDepth + 1, rootSeen, nil, nil
}

func processGroupEndElement(rootDepth int, rootSeen bool) (int, bool) {
	if !rootSeen {
		return rootDepth, false
	}
	rootDepth--
	return rootDepth, rootDepth == 0
}

func isShapeElementLocal(name string) bool {
	switch name {
	case "sp", shapeTypePicture, "graphicFrame", groupShapeTag, "cxnSp":
		return true
	default:
		return false
	}
}

func firstShapeIDInXML(shapesXML [][]byte) int {
	for _, shapeXML := range shapesXML {
		match := cNvPrIDPattern.FindSubmatch(shapeXML)
		if len(match) < cNvPrSubmatchSize {
			continue
		}
		id, err := strconv.Atoi(string(match[1]))
		if err == nil {
			return id
		}
	}
	return 0
}
