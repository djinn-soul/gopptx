package main

/*
#include <stdlib.h>

typedef uintptr_t DeckHandle;
*/
import "C"

import (
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"unsafe"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

//nolint:gochecknoglobals // global bridge state
var (
	globalErrorMu sync.RWMutex
	globalError   string
	deckRegistry  = editor.NewEditorRegistry()
)

//nolint:gochecknoglobals // theme presets
var (
	ThemeCorporate = styling.ThemeCorporate
	ThemeModern    = styling.ThemeModern
	ThemeVibrant   = styling.ThemeVibrant
	ThemeDark      = styling.ThemeDark
	ThemeNature    = styling.ThemeNature
	ThemeTech      = styling.ThemeTech
	ThemeCarbon    = styling.ThemeCarbon
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

// main is required for cgo build but not used for library
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

	// We need to write it to a temp file to open it with the editor
	// (Current editor requires a file path)
	tmpFile, err := os.CreateTemp("", "gopptx-new-*.pptx")
	if err != nil {
		setGlobalError(err)
		return 0
	}
	tmpPath := tmpFile.Name()
	_ = tmpFile.Close()

	if err := os.WriteFile(tmpPath, data, 0o600); err != nil {
		setGlobalError(err)
		return 0
	}

	e, err := editor.OpenPresentationEditor(tmpPath)
	if err != nil {
		_ = os.Remove(tmpPath)
		setGlobalError(err)
		return 0
	}
	e.SetCleanupOnClose(func() {
		_ = os.Remove(tmpPath)
	})

	h := editor.RegisterEditor(deckRegistry, e)
	return C.DeckHandle(h)
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
