package editor

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

const prstEllipse = "ellipse"

// renderShapeXML reconstructs the XML for a shape based on its parsed properties.
func (e *PresentationEditor) renderShapeXML(partPath string, s *parsedShape) ([]byte, error) {
	// Helper for XML escaping
	escape := func(s string) string {
		var buf bytes.Buffer
		if err := xml.EscapeText(&buf, []byte(s)); err != nil {
			return s
		}
		return buf.String()
	}

	if s.Type == shapeTypePicture {
		return nil, nil
	}

	// Basic preset geometry mapping (Phase 1 supports common types)
	prst := "rect"
	switch strings.ToLower(s.Type) {
	case prstEllipse, "oval":
		prst = prstEllipse
	case "triangle":
		prst = "triangle"
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

	return fmt.Appendf(
		nil,
		`<p:sp xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">`+
			`<p:nvSpPr><p:cNvPr id="%d" name="%s">%s%s</p:cNvPr><p:cNvSpPr/><p:nvPr/></p:nvSpPr>`+
			`<p:spPr>`+
			`<a:xfrm><a:off x="%d" y="%d"/><a:ext cx="%d" cy="%d"/></a:xfrm>`+
			`%s`+
			`<a:prstGeom prst="%s"><a:avLst/></a:prstGeom>`+
			`</p:spPr>`+
			`%s`+
			`</p:sp>`,
		s.ID,
		escape(s.Name),
		clickXML,
		hoverXML,
		s.X,
		s.Y,
		s.W,
		s.H,
		styleXML,
		prst,
		string(txBody),
	), nil
}
