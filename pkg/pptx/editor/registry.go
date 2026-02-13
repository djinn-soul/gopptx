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

var (
	registryMu sync.RWMutex
	registry   = make(map[Handle]*registeredEditor)
	// Start at 1000 to avoid common low integer values
	nextHandle uintptr = 1000
)

// RegisterEditor adds an editor to the global registry and returns a handle.
func RegisterEditor(e *PresentationEditor) Handle {
	if e == nil {
		return 0
	}
	registryMu.Lock()
	defer registryMu.Unlock()

	h := Handle(atomic.AddUintptr(&nextHandle, 1))
	registry[h] = &registeredEditor{editor: e}
	return h
}

// GetEditor retrieves an editor from the registry by its handle.
func GetEditor(h Handle) (*PresentationEditor, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()

	reg, ok := registry[h]
	if !ok {
		return nil, false
	}
	return reg.editor, true
}

// UnregisterEditor removes an editor from the registry and closes it.
func UnregisterEditor(h Handle) {
	registryMu.Lock()
	reg, ok := registry[h]
	if ok {
		delete(registry, h)
	}
	registryMu.Unlock()

	if ok && reg.editor != nil {
		_ = reg.editor.Close()
	}
}

// SetHandleError records the last error encountered for a specific handle.
func SetHandleError(h Handle, err error) {
	if err == nil {
		return
	}
	registryMu.Lock()
	defer registryMu.Unlock()

	if reg, ok := registry[h]; ok {
		reg.lastError = err.Error()
	}
}

// GetHandleError returns the last error string for a specific handle.
func GetHandleError(h Handle) string {
	registryMu.RLock()
	defer registryMu.RUnlock()

	if reg, ok := registry[h]; ok {
		return reg.lastError
	}
	return fmt.Sprintf("invalid handle: %d", h)
}
