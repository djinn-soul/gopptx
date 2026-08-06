package editor

import (
	"fmt"

	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

// SmartArtInfo describes an existing SmartArt diagram: what it is, how it is
// styled, and the node tree it draws.
type SmartArtInfo struct {
	ShapeID      int             `json:"shape_id"`
	LayoutURI    string          `json:"layout"`
	QuickStyleID string          `json:"quick_style"`
	ColorStyleID string          `json:"color_style"`
	Nodes        []smartart.Node `json:"nodes"`
}

// GetSmartArt reads back the SmartArt diagram at shapeID.
//
// The node tree comes from the data model, which is what PowerPoint draws from,
// so what this returns is what the slide shows rather than what was last
// written.
func (e *PresentationEditor) GetSmartArt(slideIndex, shapeID int) (SmartArtInfo, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return SmartArtInfo{}, fmt.Errorf("slide index %d out of range", slideIndex)
	}
	slideRef := e.slides[slideIndex]

	refs, err := e.resolveSmartArtParts(slideRef, shapeID)
	if err != nil {
		return SmartArtInfo{}, err
	}

	dataXML, ok := e.parts.Get(refs.DataPath)
	if !ok {
		return SmartArtInfo{}, fmt.Errorf("SmartArt data part %q not found", refs.DataPath)
	}
	nodes, err := smartart.ParseDataModelNodes(dataXML)
	if err != nil {
		return SmartArtInfo{}, fmt.Errorf("parse SmartArt nodes: %w", err)
	}

	info := SmartArtInfo{
		ShapeID:      shapeID,
		QuickStyleID: e.readSmartArtUniqueID(refs.StylePath),
		ColorStyleID: e.readSmartArtUniqueID(refs.ColorPath),
		Nodes:        nodes,
	}
	if layoutXML, ok := e.parts.Get(refs.LayoutPath); ok {
		info.LayoutURI = smartart.ExtractLayoutURI(string(layoutXML))
	}
	return info, nil
}

// ListSmartArt reads back every SmartArt diagram on a slide.
func (e *PresentationEditor) ListSmartArt(slideIndex int) ([]SmartArtInfo, error) {
	shapes, err := e.GetShapes(slideIndex)
	if err != nil {
		return nil, err
	}
	out := make([]SmartArtInfo, 0, len(shapes))
	for _, shape := range shapes {
		info, err := e.GetSmartArt(slideIndex, shape.ID)
		if err != nil {
			// Shapes that are not SmartArt frames simply do not resolve.
			continue
		}
		out = append(out, info)
	}
	return out, nil
}
