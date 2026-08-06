package editor

import (
	"encoding/json"

	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

// handleGetSmartArt reads back one SmartArt diagram.
//
// Payload: {"slide_index": N, "shape_id": N}.
// Response: {"shape_id": N, "layout": "...", "quick_style": "...",
// "color_style": "...", "nodes": [{"text": "...", "children": [...]}]}.
func handleGetSmartArt(e *PresentationEditor, payload json.RawMessage) (any, error) {
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

	info, readErr := e.GetSmartArt(slideIndex, shapeID)
	if readErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, readErr.Error())
	}
	return smartArtInfoResponse(info), nil
}

// handleListSmartArt reads back every SmartArt diagram on a slide.
//
// Payload: {"slide_index": N}.
// Response: {"diagrams": [ ... same shape as get_smartart ... ]}.
func handleListSmartArt(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}
	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}

	diagrams, readErr := e.ListSmartArt(slideIndex)
	if readErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, readErr.Error())
	}
	out := make([]any, 0, len(diagrams))
	for _, info := range diagrams {
		out = append(out, smartArtInfoResponse(info))
	}
	return map[string]any{"diagrams": out}, nil
}

func smartArtInfoResponse(info SmartArtInfo) map[string]any {
	return map[string]any{
		"shape_id":    info.ShapeID,
		"layout":      info.LayoutURI,
		"quick_style": info.QuickStyleID,
		"color_style": info.ColorStyleID,
		"nodes":       smartArtNodesResponse(info.Nodes),
	}
}

func smartArtNodesResponse(nodes []smartart.Node) []any {
	out := make([]any, 0, len(nodes))
	for _, node := range nodes {
		entry := map[string]any{"text": node.Text}
		if len(node.Children) > 0 {
			entry["children"] = smartArtNodesResponse(node.Children)
		}
		out = append(out, entry)
	}
	return out
}
