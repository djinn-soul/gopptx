package editor

import (
	"encoding/json"
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
	editormodcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/common"
)

func handleListPlaceholders(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}

	slide, err := e.GetSlide(slideIndex)
	if err != nil {
		return nil, err
	}

	placeholders, err := slide.Placeholders()
	if err != nil {
		return nil, err
	}

	result := make([]map[string]any, len(placeholders))
	for i, ph := range placeholders {
		result[i] = map[string]any{
			"index": ph.Index,
			"type":  ph.Type,
			"name":  ph.Name,
		}
	}

	return map[string]any{"placeholders": result}, nil
}

func handleSetPlaceholderContent(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}

	phIndex, err := requirePlaceholderIndex(p, v)
	if err != nil {
		return nil, err
	}

	// We support "text", "image_path", "text_style" inside the payload
	text := v.OptionalString(p, "text")
	imagePath := v.OptionalString(p, "image_path")
	tableSpec, hasTableSpec, err := editorcommand.ParsePlaceholderTableSpec(p)
	if err != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	chartFrame, hasChart, err := buildPlaceholderChartFrame(e, slideIndex, p)
	if err != nil {
		return nil, err
	}

	hasText := text != ""
	hasImagePath := imagePath != ""

	if err := validatePlaceholderContentKinds(hasText, hasImagePath, hasTableSpec, hasChart); err != nil {
		return nil, err
	}

	// For targeting, python-pptx typically relies on idx alone.
	// But we can allow optional ph_type if provided.
	phType := v.OptionalString(p, "ph_type")
	styleOpts := editormodcommon.ParsePlaceholderTextStyle(p)
	var imageRef *pptxxml.ImageRef
	if hasImagePath {
		imageRef, err = buildPlaceholderImageRef(e, slideIndex, imagePath, p)
		if err != nil {
			return nil, err
		}
	}

	// Surgical update: parse slide XML, find the shape, replace it
	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return nil, fmt.Errorf("read slide part %s: not found", partPath)
	}

	shapesList, err := parseSlideShapes(content)
	if err != nil {
		return nil, fmt.Errorf("parse shapes: %w", err)
	}

	shapeIndex, resolvedType, err := resolvePlaceholderTargetShape(shapesList, phIndex, phType)
	if err != nil {
		return nil, err
	}

	// Prepare the override spec for internal renderer
	phSpec := editormodcommon.BuildPlaceholderOverrideSpec(phIndex, resolvedType, text, imageRef, styleOpts)
	phSpec.Table = tableSpec
	phSpec.Chart = chartFrame

	newShapeXML := pptxxml.PlaceholderShape(phSpec, shapesList[shapeIndex].ID)

	newContent := replaceShapeNodes(content, shapesList, func(i int, _ *parsedShape) ([]byte, bool) {
		if i == shapeIndex {
			return []byte(newShapeXML), true
		}
		return nil, false
	})

	e.parts.Set(partPath, newContent)
	return map[string]bool{"updated": true}, nil
}
