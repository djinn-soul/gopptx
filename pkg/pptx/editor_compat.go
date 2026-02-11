package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// Re-export types for backward compatibility.
type (
	PresentationEditor   = editor.PresentationEditor
	PresentationMetadata = common.PresentationMetadata
	SlideMetadata        = common.SlideMetadata
	SlideSize            = common.SlideSize
)

// Re-export constants.
var (
	SlideSize4x3  = common.SlideSize4x3
	SlideSize16x9 = common.SlideSize16x9
)

// OpenPresentationEditor re-exports the function from the editor package.
func OpenPresentationEditor(path string) (*PresentationEditor, error) {
	return editor.OpenPresentationEditor(path)
}
