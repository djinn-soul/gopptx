package editor

import (
	"bytes"
	"errors"
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const layoutPartKind = "layout"

// AddLayoutShape appends an autoshape to a slide layout's shape tree and
// returns its new object ID.
//
// Upstream python-pptx has no way to do this (issue #1044): its layout objects
// expose shapes read-only, so anything drawn on every slide of a layout had to
// be repeated per slide.
func (e *PresentationEditor) AddLayoutShape(
	layoutPart string,
	shapeType string,
	x, y, w, h float64,
) (int, error) {
	return e.addShapeToPart(layoutPart, layoutPartKind, shapeType, x, y, w, h)
}

// AddLayoutTextBox appends a text box to a slide layout's shape tree.
func (e *PresentationEditor) AddLayoutTextBox(
	layoutPart string,
	text string,
	x, y, w, h float64,
) (int, error) {
	return e.addTextBoxToPart(layoutPart, layoutPartKind, text, x, y, w, h)
}

// AddMasterShape appends an autoshape to a slide master's shape tree.
func (e *PresentationEditor) AddMasterShape(
	masterPart string,
	shapeType string,
	x, y, w, h float64,
) (int, error) {
	return e.addShapeToPart(masterPart, "master", shapeType, x, y, w, h)
}

// AddMasterTextBox appends a text box to a slide master's shape tree.
func (e *PresentationEditor) AddMasterTextBox(
	masterPart string,
	text string,
	x, y, w, h float64,
) (int, error) {
	return e.addTextBoxToPart(masterPart, "master", text, x, y, w, h)
}

func (e *PresentationEditor) addShapeToPart(
	partPath string,
	kind string,
	shapeType string,
	x, y, w, h float64,
) (int, error) {
	if err := ValidateLayoutShapeExtents(w, h); err != nil {
		return 0, err
	}
	content, newID, err := e.prepareShapeTreePart(partPath, kind)
	if err != nil {
		return 0, err
	}

	newShape := parsedShape{
		ID:   newID,
		Name: fmt.Sprintf("%s %d", shapeType, newID),
		Type: shapeType,
		X:    int(x),
		Y:    int(y),
		W:    int(w),
		H:    int(h),
	}
	shapeXML, err := e.renderShapeXML(partPath, &newShape)
	if err != nil {
		return 0, err
	}
	return e.appendShapeToPart(partPath, content, shapeXML, newID)
}

func (e *PresentationEditor) addTextBoxToPart(
	partPath string,
	kind string,
	text string,
	x, y, w, h float64,
) (int, error) {
	if err := ValidateLayoutShapeExtents(w, h); err != nil {
		return 0, err
	}
	content, newID, err := e.prepareShapeTreePart(partPath, kind)
	if err != nil {
		return 0, err
	}

	newShape := parsedShape{
		ID:   newID,
		Name: fmt.Sprintf("TextBox %d", newID),
		Type: shapeTypeRect,
		Text: text,
		X:    int(x),
		Y:    int(y),
		W:    int(w),
		H:    int(h),
	}
	shapeXML, err := e.renderShapeXML(partPath, &newShape)
	if err != nil {
		return 0, err
	}
	shapeXML = bytes.Replace(
		shapeXML,
		[]byte("<p:cNvSpPr/>"),
		[]byte(`<p:cNvSpPr txBox="1"/>`),
		1,
	)
	return e.appendShapeToPart(partPath, content, shapeXML, newID)
}

// prepareShapeTreePart validates the target part and reserves the next object ID.
func (e *PresentationEditor) prepareShapeTreePart(
	partPath string,
	kind string,
) ([]byte, int, error) {
	if e == nil {
		return nil, 0, errors.New("editor cannot be nil")
	}
	if !e.isShapeTreePart(partPath, kind) {
		return nil, 0, fmt.Errorf("%q is not a slide %s part", partPath, kind)
	}
	content, ok := e.parts.Get(partPath)
	if !ok {
		return nil, 0, fmt.Errorf("read %s part %s: not found", kind, partPath)
	}

	newID := e.maxObjectID(partPath, content) + 1
	e.reserveObjectIDs(partPath, newID)
	return content, newID, nil
}

func (e *PresentationEditor) appendShapeToPart(
	partPath string,
	content []byte,
	shapeXML []byte,
	newID int,
) (int, error) {
	newXML, err := insertShapeXML(content, shapeXML)
	if err != nil {
		return 0, err
	}
	e.parts.Set(partPath, newXML)
	return newID, nil
}

// isShapeTreePart guards against editing an unrelated part by path typo.
func (e *PresentationEditor) isShapeTreePart(partPath string, kind string) bool {
	var known []common.SlideLayoutInfo
	switch kind {
	case layoutPartKind:
		layouts, err := e.ListSlideLayouts()
		if err != nil {
			return false
		}
		known = layouts
	case "master":
		masters, err := e.ListSlideMasters()
		if err != nil {
			return false
		}
		for _, master := range masters {
			if master.Part == partPath {
				return true
			}
		}
		return false
	default:
		return false
	}

	for _, layout := range known {
		if layout.Part == partPath {
			return true
		}
	}
	return false
}

// ValidateLayoutShapeExtents mirrors the slide-level extent rules so a layout
// shape cannot be created with a zero or negative size.
func ValidateLayoutShapeExtents(w, h float64) error {
	return validateShapeExtents(w, h)
}
