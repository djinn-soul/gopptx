package editor

import (
	"encoding/json"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
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

	if hasImagePath && imagePath != "" {
		img := shapes.Image{Path: imagePath}

		// Optional bounds for crop/positioning within placeholder
		if boundsRaw, ok := p["bounds"].([]any); ok && len(boundsRaw) == 4 {
			if x, ok := parseFloat(boundsRaw[0]); ok {
				pt := styling.Points(x)
				img.X = pt
			}
			if y, ok := parseFloat(boundsRaw[1]); ok {
				pt := styling.Points(y)
				img.Y = pt
			}
			if cx, ok := parseFloat(boundsRaw[2]); ok {
				pt := styling.Points(cx)
				img.CX = pt
			}
			if cy, ok := parseFloat(boundsRaw[3]); ok {
				pt := styling.Points(cy)
				img.CY = pt
			}
		}
		override.Image = &img
	}

	slideTitle := e.slides[slideIndex].Title
	if slideTitle == "" {
		slideTitle = " "
	}
	spec := elements.SlideContent{
		Title:                slideTitle,
		PlaceholderOverrides: []shapes.PlaceholderContent{override},
	}

	if err := e.UpdateSlide(slideIndex, spec); err != nil {
		return nil, err
	}

	return map[string]bool{"updated": true}, nil
}

func parseFloat(v any) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case float32:
		return float64(val), true
	}
	return 0, false
}
