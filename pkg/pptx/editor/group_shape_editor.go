package editor

import (
	"bytes"
	"errors"
	"fmt"
	"slices"

	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

const minFreeformPoints = 2

type freeformPoint struct {
	X int
	Y int
}

// AddGroupShape adds a group shape to the target slide.
// When shapeIDs is non-empty, those shapes are nested under the new group.
func (e *PresentationEditor) AddGroupShape(slideIndex int, shapeIDs []int) (int, error) {
	partPath, content, shapeNodes, err := e.groupShapeContext(slideIndex)
	if err != nil {
		return 0, err
	}

	newID, lastShapeEnd := nextGroupShapeID(content, shapeNodes)

	if len(shapeIDs) == 0 {
		return e.insertEmptyGroupShape(partPath, content, newID, lastShapeEnd)
	}

	selected, err := selectShapesForGrouping(shapeNodes, shapeIDs, slideIndex)
	if err != nil {
		return 0, err
	}
	order := orderShapesByPosition(shapeIDs, selected)
	groupXML := buildGroupedShapeXML(content, newID, selected, order)

	newContent := replaceSelectedShapesWithGroup(content, shapeNodes, selected, order[0], groupXML)
	e.parts.Set(partPath, newContent)
	return newID, nil
}

func (e *PresentationEditor) groupShapeContext(slideIndex int) (string, []byte, []parsedShape, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return "", nil, nil, errors.New("slide index out of range")
	}
	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return "", nil, nil, fmt.Errorf("read slide part %s: not found", partPath)
	}
	shapeNodes, err := parseSlideShapes(content)
	if err != nil {
		return "", nil, nil, fmt.Errorf("parse shapes: %w", err)
	}
	return partPath, content, shapeNodes, nil
}

func nextGroupShapeID(content []byte, shapeNodes []parsedShape) (int, int64) {
	maxID := editorshape.MaxObjectID(content, cNvPrIDPattern, cNvPrSubmatchSize)
	lastShapeEnd := int64(-1)
	for _, shape := range shapeNodes {
		if shape.End > lastShapeEnd {
			lastShapeEnd = shape.End
		}
	}
	return maxID + 1, lastShapeEnd
}

func (e *PresentationEditor) insertEmptyGroupShape(
	partPath string,
	content []byte,
	newID int,
	lastShapeEnd int64,
) (int, error) {
	groupXML := []byte(buildGroupShapeXML(newID, 0, 0, 0, 0, nil))
	out, err := insertGroupXML(content, groupXML, lastShapeEnd)
	if err != nil {
		return 0, err
	}
	e.parts.Set(partPath, out)
	return newID, nil
}

func insertGroupXML(content, groupXML []byte, lastShapeEnd int64) ([]byte, error) {
	var out bytes.Buffer
	if lastShapeEnd != -1 {
		out.Write(content[:lastShapeEnd])
		out.Write(groupXML)
		out.Write(content[lastShapeEnd:])
		return out.Bytes(), nil
	}

	endTree := []byte("</p:spTree>")
	idx := bytes.LastIndex(content, endTree)
	if idx == -1 {
		return nil, errors.New("invalid slide xml: missing spTree end")
	}
	out.Write(content[:idx])
	out.Write(groupXML)
	out.Write(content[idx:])
	return out.Bytes(), nil
}

func selectShapesForGrouping(
	shapeNodes []parsedShape,
	shapeIDs []int,
	slideIndex int,
) (map[int]parsedShape, error) {
	selected := make(map[int]parsedShape, len(shapeIDs))
	for _, id := range shapeIDs {
		shape, found := findShapeByID(shapeNodes, id)
		if !found {
			return nil, fmt.Errorf("shape id %d not found on slide %d", id, slideIndex)
		}
		selected[id] = shape
	}
	return selected, nil
}

func findShapeByID(shapes []parsedShape, shapeID int) (parsedShape, bool) {
	for _, shape := range shapes {
		if shape.ID == shapeID {
			return shape, true
		}
	}
	return parsedShape{}, false
}

func orderShapesByPosition(shapeIDs []int, selected map[int]parsedShape) []int {
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
	return order
}

func buildGroupedShapeXML(content []byte, newID int, selected map[int]parsedShape, order []int) []byte {
	childXML, minX, minY, width, height := collectGroupGeometry(content, selected, order)
	return []byte(buildGroupShapeXML(newID, minX, minY, width, height, childXML))
}

func collectGroupGeometry(
	content []byte,
	selected map[int]parsedShape,
	order []int,
) ([]byte, int, int, int, int) {
	childXML := make([]byte, 0)
	minX := selected[order[0]].X
	minY := selected[order[0]].Y
	maxX := selected[order[0]].X + selected[order[0]].W
	maxY := selected[order[0]].Y + selected[order[0]].H

	for _, id := range order {
		shape := selected[id]
		childXML = append(childXML, content[shape.Start:shape.End]...)
		minX = min(minX, shape.X)
		minY = min(minY, shape.Y)
		maxX = max(maxX, shape.X+shape.W)
		maxY = max(maxY, shape.Y+shape.H)
	}
	return childXML, minX, minY, maxX - minX, maxY - minY
}

func replaceSelectedShapesWithGroup(
	content []byte,
	shapeNodes []parsedShape,
	selected map[int]parsedShape,
	firstShapeID int,
	groupXML []byte,
) []byte {
	return replaceShapeNodes(content, shapeNodes, func(_ int, shape *parsedShape) ([]byte, bool) {
		if _, ok := selected[shape.ID]; !ok {
			return nil, false
		}
		if shape.ID == firstShapeID {
			return groupXML, true
		}
		return []byte{}, true
	})
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

// AddFreeformShape inserts a custom-geometry shape based on provided points.
