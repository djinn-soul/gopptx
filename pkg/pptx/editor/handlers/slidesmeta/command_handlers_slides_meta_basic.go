package slidesmeta

import (
	"encoding/json"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

func handleSetSlideSize(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.SlideSizeRequest, bool) {
			return editorcommand.ParseSlideSizeRequest(p, v.RequireInt64)
		},
		v.Error,
		func(request editorcommand.SlideSizeRequest) (any, error) {
			if err := e.SetSlideSize(common.SlideSize{Width: request.Width, Height: request.Height}); err != nil {
				return nil, err
			}
			return map[string]bool{"updated": true}, nil
		},
	)
}

func handleSetSlideTitle(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.SlideTitleRequest, bool) {
			return editorcommand.ParseSlideTitleRequest(
				p,
				func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
				v.RequireString,
			)
		},
		v.Error,
		func(request editorcommand.SlideTitleRequest) (any, error) {
			if err := e.SetSlideTitle(request.SlideIndex, request.Title); err != nil {
				return nil, err
			}
			return map[string]bool{"updated": true}, nil
		},
	)
}

func handleMergeFromFile(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.MergeFromFileRequest, bool) {
			return editorcommand.ParseMergeFromFileRequest(p, v.RequireString)
		},
		v.Error,
		func(request editorcommand.MergeFromFileRequest) (any, error) {
			if err := e.MergeFromFile(request.Path); err != nil {
				return nil, err
			}
			return map[string]bool{"merged": true}, nil
		},
	)
}

func handleUpdateSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.UpdateSlideRequest, bool) {
			return editorcommand.ParseUpdateSlideRequest(
				p,
				func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
				v.OptionalString,
				v.OptionalStringSlice,
			)
		},
		v.Error,
		func(request editorcommand.UpdateSlideRequest) (any, error) {
			title := request.Title
			if title == "" {
				title = e.slides[request.SlideIndex].Title
			}

			slide := elements.NewSlide(title)
			if request.Layout != "" {
				slide = slide.WithLayout(request.Layout)
			}
			for _, bullet := range request.Bullets {
				slide = slide.AddBullet(bullet)
			}
			if err := e.UpdateSlide(request.SlideIndex, slide); err != nil {
				return nil, err
			}
			return map[string]bool{"updated": true}, nil
		},
	)
}
