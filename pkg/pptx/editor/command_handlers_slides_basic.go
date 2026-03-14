package editor

import (
	"encoding/json"

	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func handleSlideCount(e *PresentationEditor, _ json.RawMessage) (any, error) {
	return map[string]int{"count": e.SlideCount()}, nil
}

func handleAddSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.AddSlideRequest, bool) {
			return editorcommand.ParseAddSlideRequest(p, v.OptionalString, v.OptionalStringSlice), true
		},
		v.Error,
		func(request editorcommand.AddSlideRequest) (any, error) {
			slide := elements.NewSlide(request.Title)
			if request.Layout != "" {
				slide = slide.WithLayout(request.Layout)
			}
			for _, bullet := range request.Bullets {
				slide = slide.AddBullet(bullet)
			}
			index, err := e.AddSlide(slide)
			if err != nil {
				return nil, err
			}
			return map[string]int{"index": index}, nil
		},
	)
}

func handleRemoveSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.SlideIndexRequest, bool) {
			return editorcommand.ParseSlideIndexOnlyRequest(p, v.RequireInt)
		},
		v.Error,
		func(request editorcommand.SlideIndexRequest) (any, error) {
			if !v.IndexBounds(request.Index, 0, e.SlideCount(), "index") {
				return nil, v.Error()
			}
			return nil, e.RemoveSlide(request.Index)
		},
	)
}

func handleMoveSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.MoveSlideRequest, bool) {
			return editorcommand.ParseMoveSlideRequest(p, v.RequireInt)
		},
		v.Error,
		func(request editorcommand.MoveSlideRequest) (any, error) {
			slideCount := e.SlideCount()
			if !v.IndexBounds(request.From, 0, slideCount, "from") {
				return nil, v.Error()
			}
			if !v.IndexBounds(request.To, 0, slideCount, "to") {
				return nil, v.Error()
			}
			return nil, e.MoveSlide(request.From, request.To)
		},
	)
}

func handleDuplicateSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.DuplicateSlideRequest, bool) {
			return editorcommand.ParseDuplicateSlideRequest(p, v.RequireInt, v.OptionalInt)
		},
		v.Error,
		func(request editorcommand.DuplicateSlideRequest) (any, error) {
			if !v.IndexBounds(request.Index, 0, e.SlideCount(), "index") {
				return nil, v.Error()
			}
			newIdx, err := e.DuplicateSlide(request.Index, request.InsertAt)
			if err != nil {
				return nil, err
			}
			return map[string]int{"new_index": newIdx}, nil
		},
	)
}

func handleGetMetadata(e *PresentationEditor, _ json.RawMessage) (any, error) {
	m := e.Metadata()
	return map[string]any{
		"title":       m.Title,
		"slide_count": m.SlideCount,
		"size": map[string]int64{
			"width":  m.SlideSize.Width,
			"height": m.SlideSize.Height,
		},
	}, nil
}
