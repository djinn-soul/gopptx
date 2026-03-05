package editor

import (
	"encoding/json"
	"errors"

	editormodcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
)

func handleBatchExecute(e *PresentationEditor, payload json.RawMessage) (any, error) {
	result, err := editormodcommand.HandleBatchExecute(
		payload,
		func(op string) (func(json.RawMessage) (any, error), bool) {
			handler, ok := commandHandlerFor(op)
			if !ok {
				return nil, false
			}
			return func(itemPayload json.RawMessage) (any, error) {
				return handler(e, itemPayload)
			}, true
		},
		func(err error) (editormodcommand.BridgeErrorView, bool) {
			var bridgeErr *BridgeError
			if errors.As(err, &bridgeErr) {
				return editormodcommand.BridgeErrorView{
					Code:    bridgeErr.Code,
					Message: bridgeErr.Message,
					Details: bridgeErr.Details,
				}, true
			}
			return editormodcommand.BridgeErrorView{}, false
		},
		editormodcommand.BatchOptions{
			BatchOp:       OpBatchExecute,
			UnknownOpCode: ErrCodeUnknownOp,
			OpFailedCode:  ErrCodeOpFailed,
		},
	)
	if err != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	return result, nil
}

func commandHandlerFor(op string) (commandHandler, bool) {
	for _, lookup := range []func(string) (commandHandler, bool){
		commandHandlerForSlides,
		commandHandlerForLayoutMetadata,
		commandHandlerForContent,
		commandHandlerForCommentsShapes,
		commandHandlerForNotesTables,
	} {
		if h, ok := lookup(op); ok {
			return h, true
		}
	}
	return nil, false
}
