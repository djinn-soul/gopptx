package editor

import (
	"encoding/json"

	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
)

// handleAddOnlineVideo links a slide to a hosted video instead of embedding it.
//
// Payload: {"slide_index", "url", "left"/"top"/"width"/"height",
// optional "poster_data" (base64), "poster_format", "alt_text"}.
func handleAddOnlineVideo(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, placement, v, err := parseMediaInsertPayload(e, payload)
	if err != nil {
		return nil, err
	}

	videoURL, ok := v.RequireString(p, "url")
	if !ok {
		return nil, v.Error()
	}
	posterData, decodeErr := editorcommand.DecodeOptionalBase64Field(
		v.OptionalString(p, "poster_data"),
		maxMediaBase64,
		"poster",
	)
	if decodeErr != nil {
		return nil, decodeErr
	}
	posterFormat := v.OptionalString(p, "poster_format")
	altText := v.OptionalString(p, "alt_text")
	if vErr := v.Error(); vErr != nil {
		return nil, vErr
	}

	shapeID, err := e.AddOnlineVideo(
		placement.SlideIndex,
		videoURL,
		posterData,
		posterFormat,
		altText,
		placement.X,
		placement.Y,
		placement.W,
		placement.H,
	)
	if err != nil {
		return nil, err
	}
	return map[string]any{keyShapeID: shapeID}, nil
}
