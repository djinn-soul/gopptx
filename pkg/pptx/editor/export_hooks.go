package editor

import (
	"encoding/json"
	"errors"
)

// exportHandlers holds optional handlers registered by external packages to
// avoid the import cycle: editor → export → pptx → editor.
// Populated via RegisterExportHandler (called from bindings/c/bridge.go init).
//
//nolint:gochecknoglobals // Global registry is required for cross-package init-time export hook registration.
var exportHandlers = map[string]commandHandler{
	// Placeholder stubs ensure the contract tests (TestSupportedOpsMatchCommandHandlers)
	// always find a handler.  editorexport.Register() replaces these with real
	// implementations at startup.
	OpExportPDF: func(_ *PresentationEditor, _ json.RawMessage) (any, error) {
		return nil, errors.New("export not initialized: call editorexport.Register() at startup")
	},
	OpExportHTML: func(_ *PresentationEditor, _ json.RawMessage) (any, error) {
		return nil, errors.New("export not initialized: call editorexport.Register() at startup")
	},
}

// RegisterExportHandler registers a handler function for the given operation
// name, replacing any existing placeholder. This must be called before any
// bridge requests are processed (e.g. in an init()).
func RegisterExportHandler(op string, h func(*PresentationEditor, json.RawMessage) (any, error)) {
	exportHandlers[op] = h
}

// RegisteredExportOps returns the list of operation names registered via
// RegisterExportHandler.
func RegisteredExportOps() []string {
	ops := make([]string, 0, len(exportHandlers))
	for op := range exportHandlers {
		ops = append(ops, op)
	}
	return ops
}
