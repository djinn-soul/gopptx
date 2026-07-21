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

// reportError writes err to the caller-supplied out-parameter when one is
// provided, falling back to the process-wide slot otherwise.
//
// The global slot is inherently racy — two threads opening decks concurrently
// will read each other's message — so it exists only to keep the pre-_ex
// exports working. New callers must use the *_ex variants and pass errOut.
// The returned string is owned by the caller and freed with deck_free_string.
func reportError(errOut **C.char, err error) {
	if errOut == nil {
		setGlobalError(err)
		return
	}
	if err != nil {
		*errOut = C.CString(err.Error())
	} else {
		*errOut = nil
	}
}

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
	editor.RegisterEditorTryLockFn(func(h int64) (func(), bool) {
		_, release, ok := editor.TryLockEditor(deckRegistry, editor.Handle(h))
		return release, ok
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

// recoverPanicTo prevents Go panics from crashing the C host, reporting through
// the caller's out-parameter rather than the shared global slot.
func recoverPanicTo(errOut **C.char) {
	if r := recover(); r != nil {
		reportError(errOut, fmt.Errorf("go panic: %v\n%s", r, debug.Stack()))
	}
}

// openDeck opens a presentation from a path, reporting failure via errOut.
func openDeck(path *C.char, errOut **C.char) C.DeckHandle {
	e, err := editor.OpenPresentationEditor(C.GoString(path))
	if err != nil {
		reportError(errOut, err)
		return 0
	}
	return C.DeckHandle(editor.RegisterEditor(deckRegistry, e))
}

// newDeck creates a minimal one-slide presentation, reporting failure via errOut.
func newDeck(title *C.char, errOut **C.char) C.DeckHandle {
	data, err := pptx.Create(C.GoString(title), 1)
	if err != nil {
		reportError(errOut, err)
		return 0
	}
	// Open directly from bytes — no temp file write/read round-trip.
	// This also works on read-only filesystems where os.CreateTemp would fail.
	e, err := editor.OpenPresentationEditorFromBytes(data)
	if err != nil {
		reportError(errOut, err)
		return 0
	}
	return C.DeckHandle(editor.RegisterEditor(deckRegistry, e))
}

// openDeckBytes opens a presentation from an in-memory buffer.
func openDeckBytes(data *C.char, length C.int, errOut **C.char) C.DeckHandle {
	e, err := editor.OpenPresentationEditorFromBytes(C.GoBytes(unsafe.Pointer(data), length))
	if err != nil {
		reportError(errOut, err)
		return 0
	}
	return C.DeckHandle(editor.RegisterEditor(deckRegistry, e))
}

// deck_open_ex opens a presentation and writes any error to *errOut, which the
// caller frees with deck_free_string. Prefer this over deck_open: it does not
// use the process-wide error slot, so concurrent opens cannot cross messages.
//
//export deck_open_ex
func deck_open_ex(path *C.char, errOut **C.char) C.DeckHandle {
	defer recoverPanicTo(errOut)
	reportError(errOut, nil)

	return openDeck(path, errOut)
}

// deck_new_ex creates a presentation, reporting errors via *errOut.
// See deck_open_ex.
//
//export deck_new_ex
func deck_new_ex(title *C.char, errOut **C.char) C.DeckHandle {
	defer recoverPanicTo(errOut)
	reportError(errOut, nil)

	return newDeck(title, errOut)
}

// deck_open_bytes_ex opens a presentation from memory, reporting errors via
// *errOut. See deck_open_ex.
//
//export deck_open_bytes_ex
func deck_open_bytes_ex(data *C.char, length C.int, errOut **C.char) C.DeckHandle {
	defer recoverPanicTo(errOut)
	reportError(errOut, nil)

	return openDeckBytes(data, length, errOut)
}

// deck_open reports errors through the process-wide slot read by
// deck_global_error, which races when decks are opened from several threads.
//
// Deprecated: use deck_open_ex.
//
//export deck_open
func deck_open(path *C.char) C.DeckHandle {
	defer recoverPanic(0)
	setGlobalError(nil) // clear previous error

	return openDeck(path, nil)
}

// deck_new reports errors through the racy process-wide slot.
//
// Deprecated: use deck_new_ex.
//
//export deck_new
func deck_new(title *C.char) C.DeckHandle {
	defer recoverPanic(0)
	setGlobalError(nil)

	return newDeck(title, nil)
}

// deck_open_bytes reports errors through the racy process-wide slot.
//
// Deprecated: use deck_open_bytes_ex.
//
//export deck_open_bytes
func deck_open_bytes(data *C.char, length C.int) C.DeckHandle {
	defer recoverPanic(0)
	setGlobalError(nil)

	return openDeckBytes(data, length, nil)
}

//export deck_save_bytes
func deck_save_bytes(h C.DeckHandle, outLen *C.int) *C.char {
	handle := editor.Handle(h)
	defer recoverPanic(handle)

	e, unlock, ok := editor.LockEditor(deckRegistry, handle)
	if !ok {
		if outLen != nil {
			*outLen = 0
		}
		return nil
	}
	defer unlock()

	data, err := e.SaveToBytes()
	if err != nil {
		editor.SetHandleError(deckRegistry, handle, err)
		if outLen != nil {
			*outLen = 0
		}
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

	e, unlock, ok := editor.LockEditor(deckRegistry, handle)
	if !ok {
		// Return a JSON error even if handle is invalid
		return C.CString(`{"ok": false, "error": {"code": "INVALID_HANDLE", "message": "Handle not found"}}`)
	}
	defer unlock()

	goInput := C.GoString(jsonInput)
	response := editor.ExecuteCommand(e, goInput)
	return C.CString(response)
}

//export deck_save
func deck_save(h C.DeckHandle, path *C.char) C.int {
	handle := editor.Handle(h)
	defer recoverPanic(handle)

	e, unlock, ok := editor.LockEditor(deckRegistry, handle)
	if !ok {
		return -1
	}
	defer unlock()

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
