package editor

import "github.com/djinn-soul/gopptx/pkg/pptx/elements"

// toString safely converts an interface{} to string.
func toString(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// slideDataForJSON converts a SlideContent to a JSON-serializable map.
func slideDataForJSON(slide elements.SlideContent) map[string]any {
	data := map[string]any{
		placeholderTypeTitle: slide.Title,
		"layout":             slide.Layout,
		"bullets":            slide.Bullets,
		"notes":              slide.Notes,
	}
	if slide.Table != nil {
		data["table"] = map[string]any{
			"rows": slide.Table.Rows,
			"x":    int64(slide.Table.X),
			"y":    int64(slide.Table.Y),
			"cx":   int64(slide.Table.CX),
			"cy":   int64(slide.Table.CY),
		}
	}
	return data
}
