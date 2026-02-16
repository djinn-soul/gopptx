package editor

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Handle is an opaque identifier for a PresentationEditor instance.
type Handle uintptr

type registeredEditor struct {
	editor    *PresentationEditor
	lastError string
}

const initialEditorHandle = 1000

// EditorRegistry stores active editor handles.
type EditorRegistry struct {
	mu sync.RWMutex

	editors    map[Handle]*registeredEditor
	nextHandle uintptr
}

// NewEditorRegistry creates an empty registry.
func NewEditorRegistry() *EditorRegistry {
	return &EditorRegistry{
		editors:    make(map[Handle]*registeredEditor),
		nextHandle: initialEditorHandle, // Start above common low integer values.
	}
}

// RegisterEditor adds an editor to the registry and returns a handle.
func (r *EditorRegistry) RegisterEditor(e *PresentationEditor) Handle {
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
func (r *EditorRegistry) GetEditor(h Handle) (*PresentationEditor, bool) {
	if r == nil {
		return nil, false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	reg, ok := r.editors[h]
	if !ok {
		return nil, false
	}
	return reg.editor, true
}

// UnregisterEditor removes an editor from the registry and closes it.
func (r *EditorRegistry) UnregisterEditor(h Handle) {
	if r == nil {
		return
	}
	r.mu.Lock()
	reg, ok := r.editors[h]
	if ok {
		delete(r.editors, h)
	}
	r.mu.Unlock()

	if ok && reg.editor != nil {
		_ = reg.editor.Close()
	}
}

// SetHandleError records the last error encountered for a specific handle.
func (r *EditorRegistry) SetHandleError(h Handle, err error) {
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
func (r *EditorRegistry) GetHandleError(h Handle) string {
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
func RegisterEditor(registry *EditorRegistry, e *PresentationEditor) Handle {
	if registry == nil {
		return 0
	}
	return registry.RegisterEditor(e)
}

// GetEditor retrieves an editor from a registry by handle.
func GetEditor(registry *EditorRegistry, h Handle) (*PresentationEditor, bool) {
	if registry == nil {
		return nil, false
	}
	return registry.GetEditor(h)
}

// UnregisterEditor removes an editor from a registry and closes it.
func UnregisterEditor(registry *EditorRegistry, h Handle) {
	if registry == nil {
		return
	}
	registry.UnregisterEditor(h)
}

// SetHandleError records the last error for a handle in a registry.
func SetHandleError(registry *EditorRegistry, h Handle, err error) {
	if registry == nil {
		return
	}
	registry.SetHandleError(h, err)
}

// GetHandleError returns the last error for a handle in a registry.
func GetHandleError(registry *EditorRegistry, h Handle) string {
	if registry == nil {
		return fmt.Sprintf("invalid handle: %d", h)
	}
	return registry.GetHandleError(h)
}
