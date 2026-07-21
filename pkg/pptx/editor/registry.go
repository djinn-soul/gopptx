package editor

import (
	"encoding/xml"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/djinn-soul/gopptx/pkg/pptx/comments"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// Handle is an opaque identifier for a PresentationEditor instance.
type Handle uintptr

type registeredEditor struct {
	// mu serializes operations on editor. It is acquired only at the FFI
	// boundary (see LockEditor) — never inside PresentationEditor methods,
	// because command handlers such as export call back into e.Save and
	// would deadlock on a re-entrant acquire.
	mu        sync.Mutex
	editor    *PresentationEditor
	lastError string
}

const initialEditorHandle = 1000

// Registry stores active editor handles.
type Registry struct {
	mu sync.RWMutex

	editors    map[Handle]*registeredEditor
	nextHandle uintptr
}

// NewEditorRegistry creates an empty registry.
func NewEditorRegistry() *Registry {
	return &Registry{
		editors:    make(map[Handle]*registeredEditor),
		nextHandle: initialEditorHandle, // Start above common low integer values.
	}
}

func (r *Registry) RegisterEditor(e *PresentationEditor) Handle {
	if r == nil || e == nil {
		return 0
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	h := Handle(atomic.AddUintptr(&r.nextHandle, 1))
	r.editors[h] = &registeredEditor{editor: e}
	return h
}

// GetEditor retrieves an editor from the registry by its handle.
func (r *Registry) GetEditor(h Handle) (*PresentationEditor, bool) {
	if r == nil {
		return nil, false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	reg, ok := r.editors[h]
	if !ok || reg.editor == nil {
		return nil, false
	}
	return reg.editor, true
}

// UnregisterEditor removes an editor from the registry and closes it.
// It blocks until any in-flight operation on the handle completes.
func (r *Registry) UnregisterEditor(h Handle) {
	if r == nil {
		return
	}
	r.mu.Lock()
	reg, ok := r.editors[h]
	if ok {
		delete(r.editors, h)
	}
	r.mu.Unlock()

	if !ok {
		return
	}
	// Wait for any in-flight operation before closing, and clear the pointer so
	// a caller already blocked in LockEditor sees the handle as gone.
	reg.mu.Lock()
	defer reg.mu.Unlock()
	if reg.editor != nil {
		_ = reg.editor.Close()
		reg.editor = nil
	}
}

// SetHandleError records the last error encountered for a specific handle.
func (r *Registry) SetHandleError(h Handle, err error) {
	if r == nil || err == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	if reg, ok := r.editors[h]; ok {
		reg.lastError = err.Error()
	}
}

// GetHandleError returns the last error string for a specific handle.
func (r *Registry) GetHandleError(h Handle) string {
	if r == nil {
		return fmt.Sprintf("invalid handle: %d", h)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	if reg, ok := r.editors[h]; ok {
		return reg.lastError
	}
	return fmt.Sprintf("invalid handle: %d", h)
}

// RegisterEditor adds an editor to a registry and returns its handle.
func RegisterEditor(registry *Registry, e *PresentationEditor) Handle {
	if registry == nil {
		return 0
	}
	return registry.RegisterEditor(e)
}

// GetEditor retrieves an editor from a registry by handle.
func GetEditor(registry *Registry, h Handle) (*PresentationEditor, bool) {
	if registry == nil {
		return nil, false
	}
	return registry.GetEditor(h)
}

// UnregisterEditor removes an editor from a registry and closes it.
func UnregisterEditor(registry *Registry, h Handle) {
	if registry == nil {
		return
	}
	registry.UnregisterEditor(h)
}

// SetHandleError records the last error for a handle in a registry.
func SetHandleError(registry *Registry, h Handle, err error) {
	if registry == nil {
		return
	}
	registry.SetHandleError(h, err)
}

// GetHandleError returns the last error for a handle in a registry.
func GetHandleError(registry *Registry, h Handle) string {
	if registry == nil {
		return fmt.Sprintf("invalid handle: %d", h)
	}
	return registry.GetHandleError(h)
}

// Metadata re-exports common.Metadata.
type Metadata = common.Metadata

// PresentationMetadata re-exports common.Metadata (for backward compatibility).
type PresentationMetadata = common.Metadata

// SlideMetadata re-exports common.SlideMetadata.
type SlideMetadata = common.SlideMetadata

// SlideSize re-exports common.SlideSize.
type SlideSize = common.SlideSize

// GetSlideSizeName returns the string representation of a common.SlideSize.
func GetSlideSizeName(size common.SlideSize) string {
	switch size {
	case common.SlideSize4x3():
		return "screen4x3"
	case common.SlideSize16x9():
		return "screen16x9"
	default:
		return ""
	}
}

func SlideSize4x3() SlideSize { return common.SlideSize4x3() }
func SlideSize16x9() SlideSize {
	return common.SlideSize16x9()
}

const commentAuthorsPartName = "ppt/commentAuthors.xml"
const authorColorCycle = 7

type cmAuthorLst struct {
	Authors []cmAuthor `xml:"cmAuthor"`
}

type cmAuthor struct {
	ID       int64  `xml:"id,attr"`
	Name     string `xml:"name,attr"`
	Initials string `xml:"initials,attr"`
	LastIdx  int    `xml:"lastIdx,attr"`
	ClrIdx   int    `xml:"clrIdx,attr"`
}

// ensureAuthorsLoadedLocked reads the commentAuthors.xml part if cache is empty.
func (e *PresentationEditor) ensureAuthorsLoadedLocked() error {
	if e.authorCache != nil {
		return nil
	}
	parsedCache := make(map[int64]comments.Author)
	nextAuthorID := int64(1)
	content, ok := e.parts.Get(commentAuthorsPartName)
	if !ok {
		e.authorCache = parsedCache
		e.nextAuthorID = 1
		return nil
	}

	var lst cmAuthorLst
	if err := xml.Unmarshal(content, &lst); err != nil {
		return fmt.Errorf("parse commentAuthors.xml: %w", err)
	}

	for _, ca := range lst.Authors {
		parsedCache[ca.ID] = comments.Author{
			ID:         ca.ID,
			Name:       ca.Name,
			Initials:   ca.Initials,
			LastIndex:  ca.LastIdx,
			ColorIndex: ca.ClrIdx,
		}
		if ca.ID >= nextAuthorID {
			nextAuthorID = ca.ID + 1
		}
	}
	e.authorCache = parsedCache
	e.nextAuthorID = nextAuthorID
	return nil
}

// GetAuthors returns a snapshot of all registered authors.
func (e *PresentationEditor) GetAuthors() ([]comments.Author, error) {
	e.authorCacheMu.Lock()
	defer e.authorCacheMu.Unlock()

	if err := e.ensureAuthorsLoadedLocked(); err != nil {
		return nil, err
	}

	out := make([]comments.Author, 0, len(e.authorCache))
	for _, a := range e.authorCache {
		out = append(out, a)
	}
	return out, nil
}

// AddAuthor registers a new author or returns an existing one if name+initials match.
func (e *PresentationEditor) AddAuthor(name, initials string) (comments.Author, error) {
	e.authorCacheMu.Lock()
	defer e.authorCacheMu.Unlock()

	if err := e.ensureAuthorsLoadedLocked(); err != nil {
		return comments.Author{}, err
	}

	for _, a := range e.authorCache {
		if a.Name == name && a.Initials == initials {
			return a, nil
		}
	}

	id := e.nextAuthorID
	e.nextAuthorID++
	colorIdx := int(id) % authorColorCycle
	newAuthor := comments.Author{
		ID:         id,
		Name:       name,
		Initials:   initials,
		LastIndex:  0,
		ColorIndex: colorIdx,
	}
	e.authorCache[id] = newAuthor
	return newAuthor, nil
}

// updateAuthorLastIndex updates the tracking index for an author.
func (e *PresentationEditor) updateAuthorLastIndex(authorID int64, newIndex int) {
	e.authorCacheMu.Lock()
	defer e.authorCacheMu.Unlock()

	if a, ok := e.authorCache[authorID]; ok {
		if newIndex > a.LastIndex {
			a.LastIndex = newIndex
			e.authorCache[authorID] = a
		}
	}
}
