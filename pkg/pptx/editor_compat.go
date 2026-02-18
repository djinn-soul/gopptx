package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// PresentationEditor and related types are re-exported for backward compatibility.
type (
	// PresentationEditor is an alias for editor.PresentationEditor.
	PresentationEditor = editor.PresentationEditor
	// SlideMetadata is an alias for common.SlideMetadata.
	SlideMetadata = common.SlideMetadata
)

// OpenPresentationEditor re-exports the function from the editor package.
func OpenPresentationEditor(path string) (*PresentationEditor, error) {
	return editor.OpenPresentationEditor(path)
}

// OpenEditor is an alias for OpenPresentationEditor.
func OpenEditor(path string) (*PresentationEditor, error) {
	return OpenPresentationEditor(path)
}

// Save re-exports the Save method.
func Save(e *PresentationEditor, filePath string) error {
	return e.Save(filePath)
}
