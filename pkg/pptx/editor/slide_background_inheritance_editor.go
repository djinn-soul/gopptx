package editor

import (
	"fmt"

	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
)

// SetSlideFollowMasterBackground toggles whether a slide inherits its master
// background. Disabling inheritance adds a no-fill background when needed.
func (e *PresentationEditor) SetSlideFollowMasterBackground(
	slideIndex int,
	follow bool,
) error {
	slideXML, partName, err := e.slideXMLAt(slideIndex)
	if err != nil {
		return err
	}
	rewritten, err := editorslide.RewriteFollowMasterBackground(slideXML, follow)
	if err != nil {
		return fmt.Errorf("rewrite slide background inheritance: %w", err)
	}
	e.parts.Set(partName, rewritten)
	return nil
}
