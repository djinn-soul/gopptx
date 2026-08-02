package editor

import (
	"encoding/json"
	"errors"
	"fmt"
)

// handleCopySlidesFrom appends selected slides from another presentation,
// either an open editor handle or a file on disk. Unlike merge_from_editor,
// which takes the whole deck, this lifts out the slides the caller names.
//
// Payload:
//
//	{
//	  "source_handle": N,          // one of source_handle / path is required
//	  "path": "other.pptx",
//	  "slide_indices": [2]         // optional; omitted means every slide
//	}
//
// Response: {"slide_count": N, "first_index": I}.
func handleCopySlidesFrom(e *PresentationEditor, payload json.RawMessage) (any, error) {
	if _, err := ParseRawPayload(payload); err != nil {
		return nil, err
	}

	var params struct {
		SourceHandle *int   `json:"source_handle"`
		Path         string `json:"path"`
		SlideIndices []int  `json:"slide_indices"`
	}
	if err := json.Unmarshal(payload, &params); err != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	if (params.SourceHandle == nil) == (params.Path == "") {
		return nil, NewBridgeError(
			ErrCodeInvalidPayload,
			"copy_slides_from requires exactly one of source_handle or path",
		)
	}
	firstIndex := len(e.slides)
	if params.Path != "" {
		if copyErr := e.CopySlidesFromFile(params.Path, params.SlideIndices); copyErr != nil {
			return nil, NewBridgeError(ErrCodeOpFailed, copyErr.Error())
		}
		return copySlidesResponse(firstIndex, len(e.slides)-firstIndex), nil
	}

	if editorLookupFn == nil {
		return nil, errors.New("editor lookup not initialized: call editor.RegisterEditorLookupFn at startup")
	}
	src, found := editorLookupFn(int64(*params.SourceHandle))
	if !found {
		return nil, NewBridgeError(
			ErrCodeOpFailed,
			fmt.Sprintf("source handle %d not found", *params.SourceHandle),
		)
	}
	// Same locking rule as merge_from_editor: take the source lock without ever
	// blocking on it, so two threads copying from each other cannot deadlock.
	// Copying from ourselves needs no second acquire.
	if src != e && editorTryLockFn != nil {
		release, locked := editorTryLockFn(int64(*params.SourceHandle))
		if !locked {
			return nil, NewBridgeError(
				ErrCodeOpFailed,
				fmt.Sprintf("source handle %d is busy in another operation", *params.SourceHandle),
			)
		}
		defer release()
	}
	if copyErr := e.CopySlidesFromEditor(src, params.SlideIndices); copyErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, copyErr.Error())
	}
	return copySlidesResponse(firstIndex, len(e.slides)-firstIndex), nil
}

func copySlidesResponse(firstIndex, count int) map[string]any {
	return map[string]any{keyFirstIndex: firstIndex, keySlideCount: count}
}
