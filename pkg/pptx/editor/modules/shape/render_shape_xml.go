package shape

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

const (
	prstEllipse  = "ellipse"
	prstOval     = "oval"
	prstTriangle = "triangle"
)

func escapeXMLText(value string) string {
	var buf bytes.Buffer
	if err := xml.EscapeText(&buf, []byte(value)); err != nil {
		return value
	}
	return buf.String()
}

// presetGeometry resolves a caller's shape type to an OOXML preset name.
//
// It used to recognize only ellipse, oval and triangle and silently emit
// prst="rect" for everything else, so every star, arrow, callout and flowchart
// shape gopptx exposes drew as a plain rectangle. Unknown names still fall back
// to "rect", because an invalid preset makes PowerPoint refuse the file.
func presetGeometry(shapeType string) string {
	switch strings.ToLower(strings.TrimSpace(shapeType)) {
	case prstEllipse, prstOval:
		return prstEllipse
	case prstTriangle:
		return prstTriangle
	}
	if canonical := shapes.NormalizeShapeType(shapeType); shapes.IsShapeType(canonical) {
		return canonical
	}
	return "rect"
}

// BuildPresetShapeXML constructs XML for standard preset-geometry shapes.
func BuildPresetShapeXML(
	id int,
	name string,
	shapeType string,
	clickXML string,
	hoverXML string,
	x int,
	y int,
	w int,
	h int,
	styleXML string,
	textBodyXML string,
) []byte {
	return fmt.Appendf(
		nil,
		`<p:sp xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">`+
			`<p:nvSpPr><p:cNvPr id="%d" name="%s">%s%s</p:cNvPr><p:cNvSpPr/><p:nvPr/></p:nvSpPr>`+
			`<p:spPr>`+
			`<a:xfrm><a:off x="%d" y="%d"/><a:ext cx="%d" cy="%d"/></a:xfrm>`+
			`<a:prstGeom prst="%s"><a:avLst/></a:prstGeom>`+
			`%s`+
			`</p:spPr>`+
			`%s`+
			`</p:sp>`,
		id,
		escapeXMLText(name),
		clickXML,
		hoverXML,
		x,
		y,
		w,
		h,
		presetGeometry(shapeType),
		styleXML,
		textBodyXML,
	)
}
