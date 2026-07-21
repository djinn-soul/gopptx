package editor

import (
	"encoding/json"
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
	editormodcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
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
			keyIndex: ph.Index,
			keyType:  ph.Type,
			keyName:  ph.Name,
		}
	}

	return map[string]any{keyPlaceholder: result}, nil
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

	phSpec, phIndex, phType, forceRect, err := buildPlaceholderOverrideSpecForPayload(e, slideIndex, p, v)
	if err != nil {
		return nil, err
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

	phSpec.Type = resolvedType
	phSpec.GeometryXML = extractPlaceholderGeometryXML(content[shapesList[shapeIndex].Start:shapesList[shapeIndex].End])
	if forceRect != nil {
		phSpec.ForceRectGeometry = forceRect
	}

	newShapeXML := pptxxml.PlaceholderShape(phSpec, shapesList[shapeIndex].ID)

	newContent := replaceShapeNodes(content, shapesList, func(i int, _ *parsedShape) ([]byte, bool) {
		if i == shapeIndex {
			return []byte(newShapeXML), true
		}
		return nil, false
	})

	e.parts.Set(partPath, newContent)
	return respUpdated, nil
}

func buildPlaceholderOverrideSpecForPayload(
	e *PresentationEditor,
	slideIndex int,
	payload map[string]any,
	v *PayloadValidator,
) (pptxxml.PlaceholderOverrideSpec, int, string, *bool, error) {
	phIndex, err := requirePlaceholderIndex(payload, v)
	if err != nil {
		return pptxxml.PlaceholderOverrideSpec{}, 0, "", nil, err
	}
	phType := v.OptionalString(payload, "ph_type")
	text := v.OptionalString(payload, "text")
	imagePath := v.OptionalString(payload, "image_path")
	styleOpts := editormodcommon.ParsePlaceholderTextStyle(payload)
	tableSpec, hasTableSpec, err := editorcommand.ParsePlaceholderTableSpec(payload)
	if err != nil {
		return pptxxml.PlaceholderOverrideSpec{}, 0, "", nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	chartFrame, hasChart, err := buildPlaceholderChartFrame(e, slideIndex, payload)
	if err != nil {
		return pptxxml.PlaceholderOverrideSpec{}, 0, "", nil, err
	}
	hasImagePath := imagePath != ""
	if err := validatePlaceholderContentKinds(text != "", hasImagePath, hasTableSpec, hasChart); err != nil {
		return pptxxml.PlaceholderOverrideSpec{}, 0, "", nil, err
	}

	var imageRef *pptxxml.ImageRef
	if hasImagePath {
		imageRef, err = buildPlaceholderImageRef(e, slideIndex, imagePath, payload)
		if err != nil {
			return pptxxml.PlaceholderOverrideSpec{}, 0, "", nil, err
		}
	}
	phSpec := editormodcommon.BuildPlaceholderOverrideSpec(phIndex, phType, text, imageRef, styleOpts)
	phSpec.Table = tableSpec
	phSpec.Chart = chartFrame
	if err := applyPlaceholderBounds(payload, &phSpec); err != nil {
		return pptxxml.PlaceholderOverrideSpec{}, 0, "", nil, err
	}

	forceRect, hasForceRect := v.OptionalBool(payload, "force_rect_geometry")
	if hasForceRect {
		return phSpec, phIndex, phType, &forceRect, nil
	}
	return phSpec, phIndex, phType, nil, nil
}

func applyPlaceholderBounds(payload map[string]any, phSpec *pptxxml.PlaceholderOverrideSpec) error {
	if _, hasBounds := payload["bounds"]; !hasBounds {
		return nil
	}
	xPt, yPt, cxPt, cyPt, err := editormodcommon.ParsePlaceholderImageBounds(payload)
	if err != nil {
		return NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	xEMU := int64(styling.Points(xPt))
	yEMU := int64(styling.Points(yPt))
	cxEMU := int64(styling.Points(cxPt))
	cyEMU := int64(styling.Points(cyPt))
	phSpec.X = &xEMU
	phSpec.Y = &yEMU
	phSpec.CX = &cxEMU
	phSpec.CY = &cyEMU
	return nil
}

func validatePlaceholderContentKinds(hasText, hasImagePath, hasTableContent, hasChart bool) error {
	kinds := 0
	if hasText {
		kinds++
	}
	if hasImagePath {
		kinds++
	}
	if hasTableContent {
		kinds++
	}
	if hasChart {
		kinds++
	}
	if kinds == 0 {
		return NewBridgeError(
			ErrCodeInvalidPayload,
			"Must provide exactly one of 'text', 'image_path', 'table', or 'chart'",
		)
	}
	if kinds > 1 {
		return NewBridgeError(
			ErrCodeInvalidPayload,
			"Only one placeholder content kind is allowed: choose exactly one of 'text', 'image_path', 'table', or 'chart'",
		)
	}
	return nil
}

func requirePlaceholderIndex(payload map[string]any, v *PayloadValidator) (int, error) {
	if val, ok := v.OptionalInt(payload, "index"); ok {
		return val, nil
	}
	if val, ok := v.OptionalInt(payload, "ph_index"); ok {
		return val, nil
	}
	return 0, NewBridgeError(ErrCodeMissingField, "missing index or ph_index")
}

func parsePlaceholderImageBounds(payload map[string]any) (float64, float64, float64, float64, error) {
	x, y, cx, cy, err := editormodcommon.ParsePlaceholderImageBounds(payload)
	if err != nil {
		return 0, 0, 0, 0, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	return x, y, cx, cy, nil
}

func buildPlaceholderImageRef(
	e *PresentationEditor,
	slideIndex int,
	imagePath string,
	payload map[string]any,
) (*pptxxml.ImageRef, error) {
	x, y, cx, cy, err := parsePlaceholderImageBounds(payload)
	if err != nil {
		return nil, err
	}

	relID, err := e.getOrCreateImageRelID(slideIndex, imagePath)
	if err != nil {
		return nil, err
	}
	return editormodcommon.BuildPlaceholderImageRef(relID, imagePath, x, y, cx, cy), nil
}

func resolvePlaceholderTargetShape(
	shapesList []parsedShape,
	phIndex int,
	phType string,
) (int, string, error) {
	shapeRefs := make([]editormodcommon.PlaceholderShapeRef, len(shapesList))
	for i, shape := range shapesList {
		shapeRefs[i] = editormodcommon.PlaceholderShapeRef{
			Index: shape.PhIndex,
			Type:  shape.PhType,
		}
	}
	shapeIndex, matches := editormodcommon.FindPlaceholderShapeIndex(shapeRefs, phIndex, phType)
	if matches > 1 && phType == "" {
		return -1, "", NewBridgeError(
			ErrCodeInvalidValue,
			fmt.Sprintf("multiple placeholders with index %d found; provide 'ph_type' to disambiguate", phIndex),
		)
	}
	if shapeIndex == -1 {
		msg := fmt.Sprintf("placeholder with index %d not found on slide", phIndex)
		if phType != "" {
			msg = fmt.Sprintf("placeholder with index %d and type %q not found on slide", phIndex, phType)
		}
		return -1, "", NewBridgeError(ErrCodeInvalidValue, msg)
	}
	resolvedType := editormodcommon.ResolvePlaceholderType(phType, shapesList[shapeIndex].PhType)
	return shapeIndex, resolvedType, nil
}
