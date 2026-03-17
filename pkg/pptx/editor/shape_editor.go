package editor

import (
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

const (
	textRunFontSizeScale = 100
	minTextFrameColumns  = 1

	actionAttrCapHint = 3
)

// renderShapeXML reconstructs the XML for a shape based on its parsed properties.
func (e *PresentationEditor) renderShapeXML(partPath string, s *parsedShape) ([]byte, error) {
	if s.Type == shapeTypePicture {
		return nil, nil
	}

	txBody, err := renderTextBodyXML(e, partPath, s)
	if err != nil {
		return nil, err
	}

	clickXML, err := e.buildClickActionXML(partPath, s.ClickAction)
	if err != nil {
		return nil, err
	}
	hoverXML, err := e.buildHoverActionXML(partPath, s.HoverAction)
	if err != nil {
		return nil, err
	}
	styleXML, err := editorshape.RenderShapeStyleXML(
		s.Fill,
		s.Line,
		s.Shadow,
		s.Glow,
		s.Blur,
		s.SoftEdge,
		s.Reflection,
	)
	if err != nil {
		return nil, err
	}

	return editorshape.BuildPresetShapeXML(
		s.ID,
		s.Name,
		s.Type,
		clickXML,
		hoverXML,
		s.X,
		s.Y,
		s.W,
		s.H,
		styleXML,
		string(txBody),
	), nil
}
