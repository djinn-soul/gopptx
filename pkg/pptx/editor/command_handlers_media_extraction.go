package editor

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

// handleListSlideMedia lists every media relationship on a slide: images,
// sounds and movies.
//
// Payload: {"slide_index": N}.
// Response: {"media": [...]}.
func handleListSlideMedia(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}

	media, err := e.ListSlideMedia(slideIndex)
	if err != nil {
		return nil, err
	}
	return map[string]any{"media": media}, nil
}

// handleExtractMedia returns the bytes of one media part, base64 encoded, so a
// caller can write an embedded movie or sound out as its own file.
//
// Payload: {"part_path": "ppt/media/media1.mp4"}.
// Response: {"data": "<base64>", "size_bytes": N}.
func handleExtractMedia(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	partPath, _ := v.RequireString(p, "part_path")
	if v.HasErrors() {
		return nil, v.Error()
	}

	data, err := e.ExtractMedia(partPath)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("media part is empty")
	}
	return map[string]any{
		"data":       base64.StdEncoding.EncodeToString(data),
		"size_bytes": len(data),
	}, nil
}
