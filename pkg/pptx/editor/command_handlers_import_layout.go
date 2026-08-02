package editor

import (
	"encoding/json"
	"errors"
	"fmt"
)

// handleImportLayoutFrom copies a layout and its master family out of another
// presentation, so a slide in this deck can be bound to a layout defined
// elsewhere. clone_layout_master_family only works within one deck.
//
// Payload:
//
//	{
//	  "source_handle": N,        // one of source_handle / path is required
//	  "path": "other.pptx",
//	  "layout_part": "ppt/slideLayouts/slideLayout12.xml"
//	}
//
// Response: {"master_part": "...", "theme_part": "...", "layout_map": {...},
// "layout_part": "..."} where layout_part is the imported copy of the requested
// layout.
func handleImportLayoutFrom(e *PresentationEditor, payload json.RawMessage) (any, error) {
	if _, err := ParseRawPayload(payload); err != nil {
		return nil, err
	}

	var params struct {
		SourceHandle *int   `json:"source_handle"`
		Path         string `json:"path"`
		LayoutPart   string `json:"layout_part"`
	}
	if err := json.Unmarshal(payload, &params); err != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}
	if params.LayoutPart == "" {
		return nil, NewBridgeError(ErrCodeInvalidPayload, "import_layout_from requires layout_part")
	}
	if (params.SourceHandle == nil) == (params.Path == "") {
		return nil, NewBridgeError(
			ErrCodeInvalidPayload,
			"import_layout_from requires exactly one of source_handle or path",
		)
	}

	if params.Path != "" {
		src, err := OpenPresentationEditor(params.Path)
		if err != nil {
			return nil, NewBridgeError(ErrCodeOpFailed, err.Error())
		}
		defer func() { _ = src.Close() }()
		return importLayoutResponse(e, src, params.LayoutPart)
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
	// Same non-blocking source lock as merge_from_editor, so two threads
	// importing from each other cannot deadlock.
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
	return importLayoutResponse(e, src, params.LayoutPart)
}

func importLayoutResponse(e, src *PresentationEditor, layoutPart string) (any, error) {
	result, err := e.ImportLayoutMasterFamily(src, layoutPart)
	if err != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, err.Error())
	}
	imported, ok := result.LayoutMap[layoutPart]
	if !ok {
		return nil, NewBridgeError(
			ErrCodeOpFailed,
			fmt.Sprintf("layout %s was not part of the imported family", layoutPart),
		)
	}
	return map[string]any{
		"master_part": result.MasterPart,
		"theme_part":  result.ThemePart,
		"layout_map":  result.LayoutMap,
		"layout_part": imported,
	}, nil
}
