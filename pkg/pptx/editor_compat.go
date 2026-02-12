package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// Re-export types for backward compatibility.
type (
	PresentationEditor = editor.PresentationEditor
	SlideMetadata      = common.SlideMetadata
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
