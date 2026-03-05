package editor

import (
	"encoding/json"

	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
)

const (
	maxMediaBase64     = 50 * 1024 * 1024
	maxEmbeddingBase64 = 20 * 1024 * 1024
)

func handleAddVideo(e *PresentationEditor, payload json.RawMessage) (any, error) {
	return handleMediaInsertCommand(
		e,
		payload,
		editorcommand.NewVideoInsertSpec(
			maxMediaBase64,
			editorcommand.AdaptVideoBinaryInsert(e.AddVideo),
			editorcommand.AdaptVideoPathInsert(e.AddVideoFromFile),
		),
	)
}

func handleAddAudio(e *PresentationEditor, payload json.RawMessage) (any, error) {
	return handleMediaInsertCommand(
		e,
		payload,
		editorcommand.NewAudioInsertSpec(
			maxMediaBase64,
			editorcommand.AdaptAudioBinaryInsert(e.AddAudio),
			editorcommand.AdaptAudioPathInsert(e.AddAudioFromFile),
		),
	)
}

func handleAddOLEObject(e *PresentationEditor, payload json.RawMessage) (any, error) {
	return handleMediaInsertCommand(
		e,
		payload,
		editorcommand.NewOLEInsertSpec(
			maxEmbeddingBase64,
			editorcommand.AdaptOLEBinaryInsert(e.AddOLEObject),
			editorcommand.AdaptOLEPathInsert(e.AddOLEObjectFromFile),
		),
	)
}

type mediaInsertSpec = editorcommand.MediaInsertSpec

func handleMediaInsertCommand(
	e *PresentationEditor,
	payload json.RawMessage,
	spec mediaInsertSpec,
) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleMediaInsertCommand(
		payload,
		e.SlideCount(),
		parseRawPayloadBytes,
		v.RequireInt,
		v.RequireFloat64,
		v.IndexBounds,
		v.OptionalString,
		v.Error,
		func(shapeID int) any { return map[string]int{"shape_id": shapeID} },
		spec,
	)
}
