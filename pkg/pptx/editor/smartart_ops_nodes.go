package editor

import (
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

// SetSmartArtNodes replaces the node tree of an existing SmartArt diagram.
func (e *PresentationEditor) SetSmartArtNodes(slideIndex, shapeID int, nodes []smartart.Node) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range", slideIndex)
	}
	slideRef := e.slides[slideIndex]

	refs, err := e.resolveSmartArtParts(slideRef, shapeID)
	if err != nil {
		return err
	}

	// Read existing layout URI from layout part.
	layoutXML, ok := e.parts.Get(refs.LayoutPath)
	if !ok {
		return fmt.Errorf("SmartArt layout part %q not found", refs.LayoutPath)
	}
	layoutURI := smartart.ExtractLayoutURI(string(layoutXML))
	if layoutURI == "" {
		return fmt.Errorf("could not determine SmartArt layout URI from %q", refs.LayoutPath)
	}

	// Build a SmartArt with the new nodes and validate.
	sa := smartart.NewSmartArt(smartart.CustomLayout(layoutURI))
	sa.Nodes = nodes
	if err := sa.Validate(slideIndex); err != nil {
		return err
	}

	spec := sa.ToSpec()
	e.parts.Set(refs.DataPath, []byte(pptxxml.SmartArtDataXML(spec)))
	e.parts.Set(refs.DrawingPath, []byte(pptxxml.SmartArtDrawingXML(spec)))
	return nil
}

func (e *PresentationEditor) readSmartArtUniqueID(partPath string) string {
	if partPath == "" {
		return ""
	}
	partXML, ok := e.parts.Get(partPath)
	if !ok {
		return ""
	}
	return smartart.ExtractUniqueID(string(partXML))
}
