package main

/*
#include <stdlib.h>
#include <stdint.h>

typedef uintptr_t DeckHandle;
*/
import "C"

import (
	"fmt"
	"runtime/debug"
	"sync"
	"unsafe"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
	"github.com/djinn-soul/gopptx/pkg/pptx/editorexport"
	"github.com/djinn-soul/gopptx/pkg/pptx/editorurlfetch"
)

//nolint:gochecknoglobals // global bridge state
var (
	globalErrorMu sync.RWMutex
	globalError   string
	deckRegistry  = editor.NewEditorRegistry()
)

// setGlobalError safely sets the global error message.
func setGlobalError(err error) {
	globalErrorMu.Lock()
	defer globalErrorMu.Unlock()
	if err != nil {
		globalError = err.Error()
	} else {
		globalError = ""
	}
}

//export deck_global_error
func deck_global_error() *C.char {
	globalErrorMu.RLock()
	defer globalErrorMu.RUnlock()
	if globalError == "" {
		return nil
	}
	return C.CString(globalError)
}

// main is required for cgo build but not used for library.
func init() { //nolint:gochecknoinits // required for cgo shared library registration
	editorexport.Register()
	editorurlfetch.Register()
	editor.RegisterEditorLookupFn(func(h int64) (*editor.PresentationEditor, bool) {
		return editor.GetEditor(deckRegistry, editor.Handle(h))
	})
}

func main() {}

// recoverPanic prevents Go panics from crashing the C host.
func recoverPanic(h editor.Handle) {
	if r := recover(); r != nil {
		err := fmt.Errorf("go panic: %v\n%s", r, debug.Stack())
		if h != 0 {
			editor.SetHandleError(deckRegistry, h, err)
		} else {
			setGlobalError(err)
		}
	}
}

//export deck_open
func deck_open(path *C.char) C.DeckHandle {
	defer recoverPanic(0)
	setGlobalError(nil) // clear previous error

	goPath := C.GoString(path)
	e, err := editor.OpenPresentationEditor(goPath)
	if err != nil {
		setGlobalError(err)
		return 0
	}

	h := editor.RegisterEditor(deckRegistry, e)
	return C.DeckHandle(h)
}

//export deck_new
func deck_new(title *C.char) C.DeckHandle {
	defer recoverPanic(0)
	setGlobalError(nil)

	goTitle := C.GoString(title)
	// Create a minimal 1-slide PPTX in memory
	data, err := pptx.Create(goTitle, 1)
	if err != nil {
		setGlobalError(err)
		return 0
	}

	// Open directly from bytes — no temp file write/read round-trip.
	// This also works on read-only filesystems where os.CreateTemp would fail.
	e, err := editor.OpenPresentationEditorFromBytes(data)
	if err != nil {
		setGlobalError(err)
		return 0
	}

	h := editor.RegisterEditor(deckRegistry, e)
	if h == 0 {
		return 0
	}
	return C.DeckHandle(h)
}

//export deck_open_bytes
func deck_open_bytes(data *C.char, length C.int) C.DeckHandle {
	defer recoverPanic(0)
	setGlobalError(nil)

	goBytes := C.GoBytes(unsafe.Pointer(data), length)
	e, err := editor.OpenPresentationEditorFromBytes(goBytes)
	if err != nil {
		setGlobalError(err)
		return 0
	}

	h := editor.RegisterEditor(deckRegistry, e)
	return C.DeckHandle(h)
}

//export deck_save_bytes
func deck_save_bytes(h C.DeckHandle, outLen *C.int) *C.char {
	handle := editor.Handle(h)
	defer recoverPanic(handle)

	e, ok := editor.GetEditor(deckRegistry, handle)
	if !ok {
		return nil
	}

	data, err := e.SaveToBytes()
	if err != nil {
		editor.SetHandleError(deckRegistry, handle, err)
		return nil
	}

	*outLen = C.int(len(data))
	buf := C.CBytes(data)
	return (*C.char)(buf)
}

//export deck_execute_json
func deck_execute_json(h C.DeckHandle, jsonInput *C.char) *C.char {
	handle := editor.Handle(h)
	defer recoverPanic(handle)

	e, ok := editor.GetEditor(deckRegistry, handle)
	if !ok {
		// Return a JSON error even if handle is invalid
		return C.CString(`{"ok": false, "error": {"code": "INVALID_HANDLE", "message": "Handle not found"}}`)
	}

	goInput := C.GoString(jsonInput)
	response := editor.ExecuteCommand(e, goInput)
	return C.CString(response)
}

//export deck_save
func deck_save(h C.DeckHandle, path *C.char) C.int {
	handle := editor.Handle(h)
	defer recoverPanic(handle)

	e, ok := editor.GetEditor(deckRegistry, handle)
	if !ok {
		return -1
	}

	goPath := C.GoString(path)
	if err := e.Save(goPath); err != nil {
		editor.SetHandleError(deckRegistry, handle, err)
		return 1
	}

	return 0
}

//export deck_last_error
func deck_last_error(h C.DeckHandle) *C.char {
	handle := editor.Handle(h)
	defer recoverPanic(handle)

	errMsg := editor.GetHandleError(deckRegistry, handle)
	return C.CString(errMsg)
}

//export deck_free_string
func deck_free_string(s *C.char) {
	if s != nil {
		C.free(unsafe.Pointer(s))
	}
}

//export deck_close
func deck_close(h C.DeckHandle) {
	handle := editor.Handle(h)
	defer recoverPanic(handle)

	editor.UnregisterEditor(deckRegistry, handle)
}
