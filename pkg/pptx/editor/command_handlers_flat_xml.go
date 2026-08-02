package editor

import (
	"encoding/json"
	"strings"
)

// handleSaveFlatXML writes the presentation as a single PowerPoint XML
// Presentation file (upstream issue #1059).
//
// Payload: {"path": "deck.xml"}.
//
// Response: {"path": "deck.xml"}.
func handleSaveFlatXML(e *PresentationEditor, payload json.RawMessage) (any, error) {
	params := struct {
		Path string `json:"path"`
	}{}
	if err := json.Unmarshal(payload, &params); err != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	if strings.TrimSpace(params.Path) == "" {
		return nil, newPayloadValidationBridgeError(ErrCodeMissingField, "missing required field: path")
	}
	if err := e.SaveFlatXML(params.Path); err != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, err.Error())
	}
	return map[string]any{"path": params.Path}, nil
}
