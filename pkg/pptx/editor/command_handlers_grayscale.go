package editor

import (
	"encoding/json"

	editorgrayscale "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/grayscale"
)

func handleConvertToGrayscale(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var opts editorgrayscale.Options
	if len(payload) > 0 {
		if err := json.Unmarshal(payload, &opts); err != nil {
			return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
		}
	}
	if err := e.ConvertToGrayscale(opts); err != nil {
		return nil, err
	}
	return respUpdated, nil
}
