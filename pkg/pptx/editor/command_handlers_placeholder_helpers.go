package editor

import (
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const placeholderBoundsLen = 4

func requirePlaceholderIndex(payload map[string]any, v *PayloadValidator) (int, error) {
	if val, ok := v.OptionalInt(payload, "index"); ok {
		return val, nil
	}
	if val, ok := v.OptionalInt(payload, "ph_index"); ok {
		return val, nil
	}
	return 0, NewBridgeError(ErrCodeMissingField, "missing index or ph_index")
}

func parsePlaceholderTextStyle(payload map[string]any) *shapes.PlaceholderOverrideOptions {
	styleMap, ok := payload["text_style"].(map[string]any)
	if !ok {
		return nil
	}

	styleOpts := &shapes.PlaceholderOverrideOptions{
		TextStyle: &shapes.PlaceholderTextStyle{},
	}
	if sizePt, ok := parseFloat(styleMap["size_pt"]); ok {
		s := int(sizePt)
		styleOpts.TextStyle.SizePt = &s
	}
	if bold, ok := styleMap["bold"].(bool); ok {
		styleOpts.TextStyle.Bold = &bold
	}
	if italic, ok := styleMap["italic"].(bool); ok {
		styleOpts.TextStyle.Italic = &italic
	}
	if color, ok := styleMap["color"].(string); ok {
		styleOpts.TextStyle.Color = &color
	}
	if font, ok := styleMap["font"].(string); ok {
		styleOpts.TextStyle.Font = &font
	}
	return styleOpts
}

func parsePlaceholderImageBounds(payload map[string]any) (float64, float64, float64, float64, error) {
	boundsRaw, ok := payload["bounds"].([]any)
	if !ok {
		return 0, 0, 0, 0, nil
	}
	if len(boundsRaw) != placeholderBoundsLen {
		return 0, 0, 0, 0, NewBridgeError(
			ErrCodeInvalidPayload,
			"bounds must be an array of 4 numbers [x, y, cx, cy]",
		)
	}
	vals := make([]float64, placeholderBoundsLen)
	for i, b := range boundsRaw {
		v, ok := parseFloat(b)
		if !ok {
			return 0, 0, 0, 0, NewBridgeError(
				ErrCodeInvalidPayload,
				fmt.Sprintf("bounds[%d] must be a number", i),
			)
		}
		vals[i] = v
	}
	return vals[0], vals[1], vals[2], vals[3], nil
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
	return &pptxxml.ImageRef{
		RelID: relID,
		Name:  imagePath,
		X:     int64(styling.Points(x)),
		Y:     int64(styling.Points(y)),
		CX:    int64(styling.Points(cx)),
		CY:    int64(styling.Points(cy)),
	}, nil
}

func findPlaceholderShapeIndex(
	shapesList []parsedShape,
	phIndex int,
	phType string,
) (int, int) {
	shapeIndex := -1
	matches := 0
	for i, s := range shapesList {
		if s.PhIndex != phIndex {
			continue
		}
		if phType != "" {
			targetType := pptxxml.NormalizePlaceholderType(phType)
			actualType := pptxxml.NormalizePlaceholderType(s.PhType)
			if targetType != actualType {
				continue
			}
		}
		shapeIndex = i
		matches++
	}
	return shapeIndex, matches
}

func resolvePlaceholderType(phType string, shape parsedShape) string {
	if phType != "" {
		return phType
	}
	return shape.PhType
}

func buildPlaceholderOverrideSpec(
	phIndex int,
	resolvedType string,
	text string,
	imageRef *pptxxml.ImageRef,
	styleOpts *shapes.PlaceholderOverrideOptions,
) pptxxml.PlaceholderOverrideSpec {
	phSpec := pptxxml.PlaceholderOverrideSpec{
		Index: phIndex,
		Type:  resolvedType,
		Text:  text,
		Image: imageRef,
	}
	if styleOpts != nil && styleOpts.TextStyle != nil {
		ts := styleOpts.TextStyle
		phSpec.TextStyle = &pptxxml.PlaceholderTextStyleSpec{
			SizePt: ts.SizePt,
			Color:  ts.Color,
			Bold:   ts.Bold,
			Italic: ts.Italic,
			Font:   ts.Font,
		}
	}
	return phSpec
}
