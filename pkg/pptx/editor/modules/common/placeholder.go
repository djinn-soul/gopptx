package common

import (
	"errors"
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const PlaceholderBoundsLen = 4

type PlaceholderShapeRef struct {
	Index int
	Type  string
}

func ParsePlaceholderTextStyle(payload map[string]any) *shapes.PlaceholderOverrideOptions {
	styleMap, ok := payload["text_style"].(map[string]any)
	if !ok {
		return nil
	}

	styleOpts := &shapes.PlaceholderOverrideOptions{
		TextStyle: &shapes.PlaceholderTextStyle{},
	}
	if sizePt, ok := ParseFloat64(styleMap["size_pt"]); ok {
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

func ParsePlaceholderImageBounds(payload map[string]any) (float64, float64, float64, float64, error) {
	boundsRaw, ok := payload["bounds"].([]any)
	if !ok {
		return 0, 0, 0, 0, nil
	}
	if len(boundsRaw) != PlaceholderBoundsLen {
		return 0, 0, 0, 0, errors.New("bounds must be an array of 4 numbers [x, y, cx, cy]")
	}
	vals := make([]float64, PlaceholderBoundsLen)
	for i, b := range boundsRaw {
		v, ok := ParseFloat64(b)
		if !ok {
			return 0, 0, 0, 0, fmt.Errorf("bounds[%d] must be a number", i)
		}
		vals[i] = v
	}
	return vals[0], vals[1], vals[2], vals[3], nil
}

func BuildPlaceholderImageRef(relID, imagePath string, x, y, cx, cy float64) *pptxxml.ImageRef {
	return &pptxxml.ImageRef{
		RelID: relID,
		Name:  imagePath,
		X:     int64(styling.Points(x)),
		Y:     int64(styling.Points(y)),
		CX:    int64(styling.Points(cx)),
		CY:    int64(styling.Points(cy)),
	}
}

func FindPlaceholderShapeIndex(shapesList []PlaceholderShapeRef, phIndex int, phType string) (int, int) {
	shapeIndex := -1
	matches := 0
	for i, s := range shapesList {
		if s.Index != phIndex {
			continue
		}
		if phType != "" {
			targetType := pptxxml.NormalizePlaceholderType(phType)
			actualType := pptxxml.NormalizePlaceholderType(s.Type)
			if targetType != actualType {
				continue
			}
		}
		shapeIndex = i
		matches++
	}
	return shapeIndex, matches
}

func ResolvePlaceholderType(phType, detectedType string) string {
	if phType != "" {
		return phType
	}
	return detectedType
}

func BuildPlaceholderOverrideSpec(
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
