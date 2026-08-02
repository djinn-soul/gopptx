package editor

import (
	"errors"
	"fmt"

	"github.com/djinn-soul/gopptx/pkg/pptx/mathml"
)

// EquationRequest places a native OMML equation on a slide.
type EquationRequest struct {
	// LaTeX is the source fragment, in the subset pkg/pptx/mathml supports.
	LaTeX string
	// X, Y, W and H are the shape bounds in EMUs.
	X, Y, W, H int
	// FontSizePt sets the equation's text size; zero leaves it to the theme.
	FontSizePt float64
}

// defaultEquationFontSizeHundredths is the sz attribute used when the caller
// names no size (24pt, PowerPoint's own default for a new equation box).
const defaultEquationFontSizeHundredths = 2400

// AddEquation inserts a text box holding a real PowerPoint equation, so the
// formula stays editable and scales with the slide instead of being pasted in
// as a picture (upstream issue #126).
func (e *PresentationEditor) AddEquation(slideIndex int, request EquationRequest) (int, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return 0, errors.New("slide index out of range")
	}
	if request.W <= 0 || request.H <= 0 {
		return 0, errors.New("equation shape needs a positive width and height")
	}
	paragraph, err := mathml.ParagraphXML(request.LaTeX)
	if err != nil {
		return 0, fmt.Errorf("render equation: %w", err)
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
	lastShapeEnd := int64(-1)
	for _, shape := range shapeNodes {
		if shape.End > lastShapeEnd {
			lastShapeEnd = shape.End
		}
	}

	newID := e.maxObjectID(partPath, content) + 1
	e.reserveObjectIDs(partPath, newID)

	sizeHundredths := defaultEquationFontSizeHundredths
	if request.FontSizePt > 0 {
		sizeHundredths = int(request.FontSizePt * 100) //nolint:mnd // sz is in hundredths of a point
	}

	shapeXML := fmt.Sprintf(
		`<p:sp xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" `+
			`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">`+
			`<p:nvSpPr><p:cNvPr id="%d" name="Equation %d"/>`+
			`<p:cNvSpPr txBox="1"/><p:nvPr/></p:nvSpPr>`+
			`<p:spPr><a:xfrm><a:off x="%d" y="%d"/><a:ext cx="%d" cy="%d"/></a:xfrm>`+
			`<a:prstGeom prst="rect"><a:avLst/></a:prstGeom><a:noFill/></p:spPr>`+
			`<p:txBody><a:bodyPr wrap="square" rtlCol="0"><a:spAutoFit/></a:bodyPr><a:lstStyle/>`+
			`%s`+
			`<a:p><a:endParaRPr lang="en-US" sz="%d" dirty="0"/></a:p>`+
			`</p:txBody></p:sp>`,
		newID, newID, request.X, request.Y, request.W, request.H, paragraph, sizeHundredths,
	)

	updatedContent, err := insertShapeIntoSlideTree(content, lastShapeEnd, shapeXML)
	if err != nil {
		return 0, err
	}
	e.parts.Set(partPath, updatedContent)
	return newID, nil
}
