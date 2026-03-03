package editor

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
)

// AddGroupShape adds a group shape to the target slide.
// When shapeIDs is non-empty, those shapes are nested under the new group.
func (e *PresentationEditor) AddGroupShape(slideIndex int, shapeIDs []int) (int, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return 0, errors.New("slide index out of range")
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return 0, fmt.Errorf("read slide part %s: not found", partPath)
	}

	shapeNodes, err := parseSlideShapes(content)
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

	if len(shapeIDs) == 0 {
		groupXML := []byte(buildGroupShapeXML(newID, 0, 0, 0, 0, nil))
		var out bytes.Buffer
		if lastShapeEnd != -1 {
			out.Write(content[:lastShapeEnd])
			out.Write(groupXML)
			out.Write(content[lastShapeEnd:])
		} else {
			endTree := []byte("</p:spTree>")
			idx := bytes.LastIndex(content, endTree)
			if idx == -1 {
				return 0, errors.New("invalid slide xml: missing spTree end")
			}
			out.Write(content[:idx])
			out.Write(groupXML)
			out.Write(content[idx:])
		}

		e.parts.Set(partPath, out.Bytes())
		return newID, nil
	}

	selected := make(map[int]parsedShape, len(shapeIDs))
	for _, id := range shapeIDs {
		found := false
		for _, shape := range shapeNodes {
			if shape.ID == id {
				selected[id] = shape
				found = true
				break
			}
		}
		if !found {
			return 0, fmt.Errorf("shape id %d not found on slide %d", id, slideIndex)
		}
	}

	order := slices.Clone(shapeIDs)
	slices.SortFunc(order, func(a, b int) int {
		shapeA := selected[a]
		shapeB := selected[b]
		switch {
		case shapeA.Start < shapeB.Start:
			return -1
		case shapeA.Start > shapeB.Start:
			return 1
		default:
			return 0
		}
	})

	childXML := make([]byte, 0)
	minX := selected[order[0]].X
	minY := selected[order[0]].Y
	maxX := selected[order[0]].X + selected[order[0]].W
	maxY := selected[order[0]].Y + selected[order[0]].H
	for _, id := range order {
		shape := selected[id]
		childXML = append(childXML, content[shape.Start:shape.End]...)
		if shape.X < minX {
			minX = shape.X
		}
		if shape.Y < minY {
			minY = shape.Y
		}
		if shape.X+shape.W > maxX {
			maxX = shape.X + shape.W
		}
		if shape.Y+shape.H > maxY {
			maxY = shape.Y + shape.H
		}
	}

	groupXML := []byte(buildGroupShapeXML(
		newID,
		minX,
		minY,
		maxX-minX,
		maxY-minY,
		childXML,
	))

	first := order[0]
	newContent := replaceShapeNodes(content, shapeNodes, func(_ int, shape *parsedShape) ([]byte, bool) {
		if _, ok := selected[shape.ID]; !ok {
			return nil, false
		}
		if shape.ID == first {
			return groupXML, true
		}
		return []byte{}, true
	})
	e.parts.Set(partPath, newContent)
	return newID, nil
}

func buildGroupShapeXML(id, x, y, w, h int, children []byte) string {
	return fmt.Sprintf(
		`<p:grpSp xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">`+
			`<p:nvGrpSpPr><p:cNvPr id="%d" name="Group %d"/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr>`+
			`<p:grpSpPr><a:xfrm><a:off x="%d" y="%d"/><a:ext cx="%d" cy="%d"/><a:chOff x="%d" y="%d"/><a:chExt cx="%d" cy="%d"/></a:xfrm></p:grpSpPr>`+
			`%s`+
			`</p:grpSp>`,
		id,
		id,
		x,
		y,
		w,
		h,
		x,
		y,
		w,
		h,
		string(children),
	)
}
