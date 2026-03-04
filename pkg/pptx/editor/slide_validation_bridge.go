package editor

import (
	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func (e *PresentationEditor) slideRelationships(slidePart string) ([]common.EditorRelationship, error) {
	return editorslide.SlideRelationships(slidePart, e.parts.Get, parseRelationshipsXML)
}

func validateEditorSlideContent(slide elements.SlideContent) error {
	return editorslide.ValidateEditorSlideContent(slide)
}
