package editor

import (
	"encoding/json"
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func handleListPlaceholders(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := v.RequireInt(p, "slide_index")
	if !ok {
		return nil, v.Error()
	}
	if !v.IndexBounds(slideIndex, 0, e.SlideCount(), "slide_index") {
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
	slideIndex, ok := v.RequireInt(p, "slide_index")
	if !ok {
		return nil, v.Error()
	}
	if !v.IndexBounds(slideIndex, 0, e.SlideCount(), "slide_index") {
		return nil, v.Error()
	}

	phIndex, ok := v.RequireInt(p, "ph_index")
	if !ok {
		return nil, v.Error()
	}

	// We support "text", "image_path", "text_style" inside the payload
	text := v.OptionalString(p, "text")
	imagePath := v.OptionalString(p, "image_path")

	hasText := text != ""
	hasImagePath := imagePath != ""

	if !hasText && !hasImagePath {
		return nil, NewBridgeError(ErrCodeInvalidPayload, "Must provide 'text' or 'image_path'")
	}

	override := shapes.PlaceholderContent{
		Index: phIndex,
	}

	// For targeting, python-pptx typically relies on idx alone.
	// But we can allow optional ph_type if provided.
	phType := v.OptionalString(p, "ph_type")
	if phType != "" {
		override.Type = phType
	}

	if hasText {
		override.Text = text
	}

	var styleOpts *shapes.PlaceholderOverrideOptions
	if styleMap, ok := p["text_style"].(map[string]any); ok {
		styleOpts = &shapes.PlaceholderOverrideOptions{
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
	}
	override.Override = styleOpts

	if hasImagePath {
		img := shapes.Image{Path: imagePath}

		// Optional bounds for crop/positioning within placeholder
		if boundsRaw, ok := p["bounds"].([]any); ok {
			if len(boundsRaw) != 4 {
				return nil, NewBridgeError(ErrCodeInvalidPayload, "bounds must be an array of 4 numbers [x, y, cx, cy]")
			}
			vals := make([]float64, 4)
			for i, b := range boundsRaw {
				v, ok := parseFloat(b)
				if !ok {
					return nil, NewBridgeError(ErrCodeInvalidPayload, fmt.Sprintf("bounds[%d] must be a number", i))
				}
				vals[i] = v
			}
			img.X, img.Y, img.CX, img.CY = styling.Points(vals[0]), styling.Points(vals[1]), styling.Points(vals[2]), styling.Points(vals[3])
		}

		// Since we're inserting an image, we need to register it and get a relationship ID
		relID, err := e.getOrCreateImageRelID(slideIndex, img.Path)
		if err != nil {
			return nil, err
		}
		img.RelID = relID
		override.Image = &img
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

	shapeIndex := -1
	matches := 0
	for i, s := range shapesList {
		if s.PhIndex == phIndex {
			// If phType was provided, we must match it to disambiguate collisions.
			// Both are normalized to handle defaults (like empty -> "obj").
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
	}

	if matches > 1 && phType == "" {
		return nil, NewBridgeError(ErrCodeInvalidValue, fmt.Sprintf("multiple placeholders with index %d found; provide 'ph_type' to disambiguate", phIndex))
	}

	if shapeIndex == -1 {
		msg := fmt.Sprintf("placeholder with index %d not found on slide", phIndex)
		if phType != "" {
			msg = fmt.Sprintf("placeholder with index %d and type %q not found on slide", phIndex, phType)
		}
		return nil, NewBridgeError(ErrCodeInvalidValue, msg)
	}

	// Prepare the override spec for internal renderer
	resolvedType := phType
	if resolvedType == "" {
		resolvedType = shapesList[shapeIndex].PhType
	}
	phSpec := pptxxml.PlaceholderOverrideSpec{
		Index: phIndex,
		Type:  resolvedType,
		Text:  override.Text,
	}
	if override.Image != nil {
		phSpec.Image = &pptxxml.ImageRef{
			RelID: override.Image.RelID,
			Name:  override.Image.Path,
			X:     int64(override.Image.X),
			Y:     int64(override.Image.Y),
			CX:    int64(override.Image.CX),
			CY:    int64(override.Image.CY),
		}
	}
	if override.Override != nil && override.Override.TextStyle != nil {
		ts := override.Override.TextStyle
		phSpec.TextStyle = &pptxxml.PlaceholderTextStyleSpec{
			SizePt: ts.SizePt,
			Color:  ts.Color,
			Bold:   ts.Bold,
			Italic: ts.Italic,
			Font:   ts.Font,
		}
	}

	newShapeXML := pptxxml.PlaceholderShape(phSpec, shapesList[shapeIndex].ID)

	newContent := replaceShapeNodes(content, shapesList, func(i int, p *parsedShape) ([]byte, bool) {
		if i == shapeIndex {
			return []byte(newShapeXML), true
		}
		return nil, false
	})

	e.parts.Set(partPath, newContent)
	return map[string]bool{"updated": true}, nil
}
