package editor

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

// extractSmartArtTexts extracts node text strings in order from a SmartArt data XML.
func extractSmartArtTexts(dataXML string) []string {
	re := regexp.MustCompile(`<a:t>([^<]*)</a:t>`)
	matches := re.FindAllStringSubmatch(dataXML, -1)
	texts := make([]string, 0, len(matches))
	for _, m := range matches {
		text := strings.TrimSpace(m[1])
		// Skip placeholder empty nodes.
		if text != "" {
			texts = append(texts, text)
		}
	}
	return texts
}

// extractSmartArtLayoutURI extracts the layout URI from a SmartArt layout XML part.
func extractSmartArtLayoutURI(layoutXML string) string {
	re := regexp.MustCompile(`uniqueId\s*=\s*["']([^"']+)["']`)
	if m := re.FindStringSubmatch(layoutXML); m != nil {
		return m[1]
	}
	// Fall back to xmlns or other URI patterns.
	re2 := regexp.MustCompile(`dgm:layoutDef[^>]*uniqueId\s*=\s*["']([^"']+)["']`)
	if m := re2.FindStringSubmatch(layoutXML); m != nil {
		return m[1]
	}
	return ""
}

// handleDeleteSmartArt removes a SmartArt diagram by shape ID.
//
// Payload: {"slide_index": N, "shape_id": N}.
// Response: {"deleted": true}.
func handleDeleteSmartArt(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}
	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	shapeID, ok := v.RequireInt(p, "shape_id")
	if !ok {
		return nil, v.Error()
	}
	if delErr := e.DeleteSmartArt(slideIndex, shapeID); delErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, delErr.Error())
	}
	return map[string]bool{"deleted": true}, nil
}

// handleChangeSmartArtLayout changes the layout of an existing SmartArt.
//
// Payload: {"slide_index": N, "shape_id": N, "layout": "layout_uri"}.
// Response: {"updated": true}.
func handleChangeSmartArtLayout(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}
	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	shapeID, ok := v.RequireInt(p, "shape_id")
	if !ok {
		return nil, v.Error()
	}
	layoutStr, ok := v.RequireString(p, "layout")
	if !ok {
		return nil, v.Error()
	}
	layout := smartart.CustomLayout(layoutStr)
	if updateErr := e.ChangeSmartArtLayout(slideIndex, shapeID, layout); updateErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, updateErr.Error())
	}
	return map[string]bool{"updated": true}, nil
}

// handleSetSmartArtStyle sets the quick style and/or color style of a SmartArt.
//
// Payload: {"slide_index": N, "shape_id": N, "quick_style": "...", "color_style": "..."}.
// Response: {"updated": true}.
func handleSetSmartArtStyle(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}
	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	shapeID, ok := v.RequireInt(p, "shape_id")
	if !ok {
		return nil, v.Error()
	}
	quickStyle := v.OptionalString(p, "quick_style")
	colorStyle := v.OptionalString(p, "color_style")
	if v.HasErrors() {
		return nil, v.Error()
	}
	if updateErr := e.SetSmartArtStyle(slideIndex, shapeID, quickStyle, colorStyle); updateErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, updateErr.Error())
	}
	return map[string]bool{"updated": true}, nil
}

// handleSetSmartArtNodes replaces the node tree of an existing SmartArt.
//
// Payload: {"slide_index": N, "shape_id": N, "items": ["text1", ...]}.
// Response: {"updated": true}.
func handleSetSmartArtNodes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}
	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	shapeID, ok := v.RequireInt(p, "shape_id")
	if !ok {
		return nil, v.Error()
	}
	items, _ := v.OptionalStringSlice(p, "items")
	if v.HasErrors() {
		return nil, v.Error()
	}
	nodes := make([]smartart.Node, len(items))
	for i, text := range items {
		nodes[i] = smartart.NewNode(text)
	}
	if updateErr := e.SetSmartArtNodes(slideIndex, shapeID, nodes); updateErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, updateErr.Error())
	}
	return map[string]bool{"updated": true}, nil
}
