package editor

import (
	"encoding/json"

	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

// handleAddSmartArtNode inserts one node into an existing diagram.
//
// Payload: {"slide_index": N, "shape_id": N, "text": "...", "parent_path": [0],
// "index": 1, "color": "C00000", "image": "path.png"}. Omit parent_path for a
// top-level node and index to append.
// Response: {"updated": true}.
func handleAddSmartArtNode(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, v, slideIndex, shapeID, ok := smartArtNodeRequest(e, payload)
	if !ok {
		return nil, v.Error()
	}
	text := v.OptionalString(p, "text")
	node := smartart.NewNode(text)
	if color := v.OptionalString(p, "color"); color != "" {
		node = node.WithColor(color)
	}
	if image := v.OptionalString(p, "image"); image != "" {
		node = node.WithImage(image)
	}
	index, hasIndex := v.OptionalInt(p, "index")
	if !hasIndex {
		index = -1
	}
	if err := e.AddSmartArtNode(
		slideIndex, shapeID, smartArtNodePath(p, "parent_path"), index, node,
	); err != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, err.Error())
	}
	return respUpdated, nil
}

// handleRemoveSmartArtNode deletes a node and its children.
//
// Payload: {"slide_index": N, "shape_id": N, "path": [0, 1]}.
// Response: {"updated": true}.
func handleRemoveSmartArtNode(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, v, slideIndex, shapeID, ok := smartArtNodeRequest(e, payload)
	if !ok {
		return nil, v.Error()
	}
	if err := e.RemoveSmartArtNode(slideIndex, shapeID, smartArtNodePath(p, "path")); err != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, err.Error())
	}
	return respUpdated, nil
}

// handleUpdateSmartArtNode changes one node's text, colour or picture.
//
// Payload: {"slide_index": N, "shape_id": N, "path": [0], "text": "...",
// "color": "C00000", "image": "path.png"}. Omitted fields are left alone.
// Response: {"updated": true}.
func handleUpdateSmartArtNode(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, v, slideIndex, shapeID, ok := smartArtNodeRequest(e, payload)
	if !ok {
		return nil, v.Error()
	}
	text := v.OptionalString(p, "text")
	color := v.OptionalString(p, "color")
	image := v.OptionalString(p, "image")
	change := SmartArtNodeChange{Text: text, Color: color, ImagePath: image}
	if err := e.UpdateSmartArtNode(slideIndex, shapeID, smartArtNodePath(p, "path"), change); err != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, err.Error())
	}
	return respUpdated, nil
}

// smartArtNodeRequest reads the slide and shape every node op needs.
func smartArtNodeRequest(
	e *PresentationEditor,
	payload json.RawMessage,
) (map[string]any, *PayloadValidator, int, int, bool) {
	v := NewPayloadValidator()
	p, err := ParseRawPayload(payload)
	if err != nil {
		// Reading a missing field records the same malformed-payload error.
		v.RequireInt(nil, "slide_index")
		return nil, v, 0, 0, false
	}
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v, 0, 0, false
	}
	shapeID, ok := v.RequireInt(p, "shape_id")
	if !ok {
		return nil, v, 0, 0, false
	}
	return p, v, slideIndex, shapeID, true
}

// smartArtNodePath reads a node path: the index of each step down the tree.
func smartArtNodePath(payload map[string]any, key string) []int {
	raw, ok := payload[key].([]any)
	if !ok {
		return nil
	}
	out := make([]int, 0, len(raw))
	for _, item := range raw {
		switch value := item.(type) {
		case float64:
			out = append(out, int(value))
		case int:
			out = append(out, value)
		}
	}
	return out
}
