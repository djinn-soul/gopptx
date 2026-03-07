package editor

import (
	"errors"
	"fmt"
	"sync"

	"github.com/djinn-soul/gopptx/pkg/pptx/comments"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
	"github.com/djinn-soul/gopptx/pkg/pptx/validation/logical"
	"github.com/djinn-soul/gopptx/pkg/pptx/validation/structural"
)

// Section describes a PowerPoint section entry.
type Section = editorslide.SectionData

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

func (e *PresentationEditor) populateSlideTitlesConcurrently() {
	if e == nil || len(e.slides) == 0 {
		return
	}
	titles := editorslide.PopulateSlideTitlesConcurrently(e.slides, e.parts.Get)
	for idx, title := range titles {
		e.slides[idx].Title = title
	}
}
