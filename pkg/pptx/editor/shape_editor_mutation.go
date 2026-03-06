package editor

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
)

func (e *PresentationEditor) AddShape(slideIndex int, shapeType string, x, y, w, h float64) (int, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return 0, errors.New("slide index out of range")
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return 0, fmt.Errorf("read slide part %s: not found", partPath)
	}

	// Parse existing shapes to find max ID and last shape position
	// OPTIMIZATION: We only need the offsets, not the full properties.
	shapes, err := scanShapesWithOffsets(content, true) // true = skip properties parsing
	if err != nil {
		return 0, fmt.Errorf("parse shapes: %w", err)
	}

	maxID := editorshape.MaxObjectID(content, cNvPrIDPattern, cNvPrSubmatchSize)
	lastShapeEnd := int64(-1)
	for _, s := range shapes {
		if s.End > lastShapeEnd {
			lastShapeEnd = s.End
		}
	}
	newID := maxID + 1

	newShape := parsedShape{
		ID:   newID,
		Name: fmt.Sprintf("%s %d", shapeType, newID),
		Type: shapeType,
		Text: "",
		X:    int(x),
		Y:    int(y),
		W:    int(w),
		H:    int(h),
	}

	shapeXML, err := e.renderShapeXML(partPath, &newShape)
	if err != nil {
		return 0, err
	}

	// Insertion point: After last shape if exists, else before </p:spTree>
	var buf bytes.Buffer
	if lastShapeEnd != -1 {
		buf.Write(content[:lastShapeEnd])
		buf.Write(shapeXML)
		buf.Write(content[lastShapeEnd:])
	} else {
		endTree := []byte("</p:spTree>")
		idx := bytes.LastIndex(content, endTree)
		if idx == -1 {
			return 0, errors.New("invalid slide xml: missing spTree end")
		}
		buf.Write(content[:idx])
		buf.Write(shapeXML)
		buf.Write(content[idx:])
	}

	e.parts.Set(partPath, buf.Bytes())
	return newID, nil
}

var cNvPrIDPattern = regexp.MustCompile(`\bcNvPr\b[^>]*\bid="(\d+)"`)

const cNvPrSubmatchSize = 2
const placeholderTypeTitle = "title"

func (e *PresentationEditor) UpdateShape(slideIndex, shapeID int, updates common.ShapeUpdate) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return errors.New("slide index out of range")
	}

	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return fmt.Errorf("read slide part %s: not found", partPath)
	}

	shapes, err := parseSlideShapes(content) // parses basic properties and byte ranges
	if err != nil {
		return fmt.Errorf("parse shapes: %w", err)
	}

	updater := shapeUpdater{
		editor:    e,
		partPath:  partPath,
		shapeID:   shapeID,
		updates:   updates,
		origSlide: content,
	}
	newXML := replaceShapeNodes(content, shapes, updater.apply)
	if updater.err != nil {
		return updater.err
	}
	if !updater.found {
		return fmt.Errorf("shape id %d not found on slide %d", shapeID, slideIndex)
	}

	e.parts.Set(partPath, newXML)
	return nil
}

type shapeUpdater struct {
	editor    *PresentationEditor
	partPath  string
	shapeID   int
	updates   common.ShapeUpdate
	origSlide []byte
	found     bool
	err       error
}

func (u *shapeUpdater) apply(_ int, s *parsedShape) ([]byte, bool) {
	if u.err != nil || !u.matchesTarget(s) {
		return nil, false
	}
	u.found = true
	if hasPictureUpdateFields(u.updates) && s.Type != shapeTypePicture {
		u.err = fmt.Errorf("shape id %d is not a picture shape", u.shapeID)
		return nil, false
	}

	updatedXML := u.origSlide[s.Start:s.End]
	replaced := false

	updatedXML, replaced = u.applyTransforms(updatedXML, s, replaced)
	updatedXML, replaced, u.err = u.applyPicture(updatedXML, s, replaced)
	if u.err != nil {
		return nil, false
	}
	updatedXML, replaced, u.err = u.applyText(updatedXML, s, replaced)
	if u.err != nil {
		return nil, false
	}
	updatedXML, replaced, u.err = u.applyStyle(updatedXML, s, replaced)
	if u.err != nil {
		return nil, false
	}
	updatedXML, replaced, u.err = u.applyActions(updatedXML, s, replaced)
	if u.err != nil {
		return nil, false
	}
	if replaced {
		return updatedXML, true
	}
	return nil, false
}

func (u *shapeUpdater) matchesTarget(s *parsedShape) bool {
	return s.ID == u.shapeID || (s.PhType == placeholderTypeTitle && u.shapeID == 0)
}

func (u *shapeUpdater) applyTransforms(
	xmlData []byte,
	s *parsedShape,
	replaced bool,
) ([]byte, bool) {
	if u.updates.X == nil && u.updates.Y == nil && u.updates.W == nil && u.updates.H == nil {
		return xmlData, replaced
	}
	if u.updates.X != nil {
		s.X = *u.updates.X
	}
	if u.updates.Y != nil {
		s.Y = *u.updates.Y
	}
	if u.updates.W != nil {
		s.W = *u.updates.W
	}
	if u.updates.H != nil {
		s.H = *u.updates.H
	}
	return editorshape.UpdateShapeTransforms(xmlData, s.X, s.Y, s.W, s.H), true
}

func (u *shapeUpdater) applyPicture(
	xmlData []byte,
	s *parsedShape,
	replaced bool,
) ([]byte, bool, error) {
	if !hasPictureUpdateFields(u.updates) {
		return xmlData, replaced, nil
	}
	if s.Type != shapeTypePicture {
		return xmlData, replaced, fmt.Errorf("shape id %d is not a picture shape", s.ID)
	}
	updatedXML, err := applyPictureShapeUpdates(xmlData, u.updates)
	if err != nil {
		return xmlData, replaced, err
	}
	return updatedXML, true, nil
}

func (u *shapeUpdater) applyText(xmlData []byte, s *parsedShape, replaced bool) ([]byte, bool, error) {
	if u.updates.Text == nil && u.updates.Runs == nil && u.updates.TextFrame == nil && u.updates.Paragraph == nil {
		return xmlData, replaced, nil
	}
	if u.updates.Text != nil {
		s.Text = *u.updates.Text
		s.Runs = nil
	}
	if u.updates.Runs != nil {
		s.Runs = *u.updates.Runs
	}
	if u.updates.TextFrame != nil {
		s.TextFrame = u.updates.TextFrame
	}
	if u.updates.Paragraph != nil {
		s.Paragraph = u.updates.Paragraph
	}
	updatedXML, err := replaceShapeTextBody(u.editor, u.partPath, xmlData, s)
	return updatedXML, true, err
}

func (u *shapeUpdater) applyStyle(xmlData []byte, s *parsedShape, replaced bool) ([]byte, bool, error) {
	applyFill := u.updates.Fill != nil
	applyLine := u.updates.Line != nil
	applyEffects := u.updates.Shadow != nil ||
		u.updates.Glow != nil ||
		u.updates.Blur != nil ||
		u.updates.SoftEdge != nil ||
		u.updates.Reflection != nil
	if u.updates.Fill == nil &&
		u.updates.Line == nil &&
		u.updates.Shadow == nil &&
		u.updates.Glow == nil &&
		u.updates.Blur == nil &&
		u.updates.SoftEdge == nil &&
		u.updates.Reflection == nil {
		return xmlData, replaced, nil
	}
	if u.updates.Fill != nil {
		s.Fill = u.updates.Fill
	}
	if u.updates.Line != nil {
		s.Line = u.updates.Line
	}
	if u.updates.Shadow != nil {
		s.Shadow = u.updates.Shadow
	}
	if u.updates.Glow != nil {
		s.Glow = u.updates.Glow
	}
	if u.updates.Blur != nil {
		s.Blur = u.updates.Blur
	}
	if u.updates.SoftEdge != nil {
		s.SoftEdge = u.updates.SoftEdge
	}
	if u.updates.Reflection != nil {
		s.Reflection = u.updates.Reflection
	}
	updatedXML, err := replaceShapeStyleSelective(
		xmlData,
		s.Fill,
		s.Line,
		s.Shadow,
		s.Glow,
		s.Blur,
		s.SoftEdge,
		s.Reflection,
		applyFill,
		applyLine,
		applyEffects,
	)
	return updatedXML, true, err
}

func (u *shapeUpdater) applyActions(xmlData []byte, s *parsedShape, replaced bool) ([]byte, bool, error) {
	if u.updates.ClickAction == nil && u.updates.HoverAction == nil {
		return xmlData, replaced, nil
	}
	if u.updates.ClickAction != nil {
		s.ClickAction = u.updates.ClickAction
	}
	if u.updates.HoverAction != nil {
		s.HoverAction = u.updates.HoverAction
	}
	updatedXML, err := replaceShapeActions(
		u.editor,
		u.partPath,
		xmlData,
		u.updates.ClickAction,
		u.updates.HoverAction,
	)
	return updatedXML, true, err
}
