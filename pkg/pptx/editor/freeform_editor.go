package editor

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type freeformPoint struct {
	X int
	Y int
}

// AddFreeformShape inserts a custom-geometry shape based on provided points.
func (e *PresentationEditor) AddFreeformShape(
	slideIndex int,
	points []freeformPoint,
	closePath bool,
) (int, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return 0, errors.New("slide index out of range")
	}
	if len(points) < 2 {
		return 0, errors.New("freeform requires at least two points")
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return 0, fmt.Errorf("read slide part %s: not found", partPath)
	}

	shapeNodes, err := scanShapesWithOffsets(content, true)
	if err != nil {
		return 0, fmt.Errorf("parse shapes: %w", err)
	}

	maxID := maxObjectID(content)
	lastShapeEnd := int64(-1)
	for _, shape := range shapeNodes {
		if shape.End > lastShapeEnd {
			lastShapeEnd = shape.End
		}
	}
	newID := maxID + 1

	minX, minY := points[0].X, points[0].Y
	maxX, maxY := points[0].X, points[0].Y
	for _, pt := range points[1:] {
		if pt.X < minX {
			minX = pt.X
		}
		if pt.Y < minY {
			minY = pt.Y
		}
		if pt.X > maxX {
			maxX = pt.X
		}
		if pt.Y > maxY {
			maxY = pt.Y
		}
	}
	width := maxX - minX
	height := maxY - minY
	if width <= 0 {
		width = 1
	}
	if height <= 0 {
		height = 1
	}

	localPts := make([]freeformPoint, len(points))
	for i, pt := range points {
		localPts[i] = freeformPoint{X: pt.X - minX, Y: pt.Y - minY}
	}

	shapeXML := renderFreeformShapeXML(newID, minX, minY, width, height, localPts, closePath)

	var out bytes.Buffer
	if lastShapeEnd != -1 {
		out.Write(content[:lastShapeEnd])
		out.WriteString(shapeXML)
		out.Write(content[lastShapeEnd:])
	} else {
		endTree := []byte("</p:spTree>")
		idx := bytes.LastIndex(content, endTree)
		if idx == -1 {
			return 0, errors.New("invalid slide xml: missing spTree end")
		}
		out.Write(content[:idx])
		out.WriteString(shapeXML)
		out.Write(content[idx:])
	}

	e.parts.Set(partPath, out.Bytes())
	return newID, nil
}

func renderFreeformShapeXML(
	id, x, y, w, h int,
	points []freeformPoint,
	closePath bool,
) string {
	var path strings.Builder
	first := points[0]
	path.WriteString(fmt.Sprintf(`<a:moveTo><a:pt x="%d" y="%d"/></a:moveTo>`, first.X, first.Y))
	for _, pt := range points[1:] {
		path.WriteString(fmt.Sprintf(`<a:lnTo><a:pt x="%d" y="%d"/></a:lnTo>`, pt.X, pt.Y))
	}
	if closePath {
		path.WriteString(`<a:close/>`)
	}

	return fmt.Sprintf(
		`<p:sp xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">`+
			`<p:nvSpPr><p:cNvPr id="%d" name="Freeform %d"/><p:cNvSpPr/><p:nvPr/></p:nvSpPr>`+
			`<p:spPr>`+
			`<a:xfrm><a:off x="%d" y="%d"/><a:ext cx="%d" cy="%d"/></a:xfrm>`+
			`<a:custGeom><a:avLst/><a:gdLst/><a:ahLst/><a:cxnLst/><a:rect l="l" t="t" r="r" b="b"/>`+
			`<a:pathLst><a:path w="%d" h="%d">%s</a:path></a:pathLst>`+
			`</a:custGeom>`+
			`</p:spPr>`+
			`</p:sp>`,
		id,
		id,
		x,
		y,
		w,
		h,
		w,
		h,
		path.String(),
	)
}
