package editor

import (
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	editormodcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/common"
)

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
