package editor

import (
	"fmt"

	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
)

// SlideShowMasterShapes reports whether a slide displays shapes inherited from
// its slide master.
func (e *PresentationEditor) SlideShowMasterShapes(slideIndex int) (bool, error) {
	slideXML, _, err := e.slideXMLAt(slideIndex)
	if err != nil {
		return false, err
	}
	visible, err := editorslide.ParseSlideShowMasterShapes(slideXML)
	if err != nil {
		return false, fmt.Errorf("parse master-shape visibility: %w", err)
	}
	return visible, nil
}

// SetSlideShowMasterShapes toggles shapes inherited from a slide's master.
func (e *PresentationEditor) SetSlideShowMasterShapes(
	slideIndex int,
	visible bool,
) error {
	slideXML, partName, err := e.slideXMLAt(slideIndex)
	if err != nil {
		return err
	}
	rewritten, err := editorslide.RewriteSlideShowMasterShapes(slideXML, visible)
	if err != nil {
		return fmt.Errorf("rewrite master-shape visibility: %w", err)
	}
	e.parts.Set(partName, rewritten)
	return nil
}

func (e *PresentationEditor) slideXMLAt(slideIndex int) ([]byte, string, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, "", fmt.Errorf("slide index %d out of range", slideIndex)
	}
	partName := e.slides[slideIndex].Part
	slideXML, ok := e.parts.Get(partName)
	if !ok {
		return nil, "", fmt.Errorf("slide part %q not found", partName)
	}
	return slideXML, partName, nil
}
