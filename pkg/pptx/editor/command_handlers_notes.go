package editor

import (
	"encoding/json"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorcommand "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/command"
)

func handleGetNotes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleSlideIndexRequest(
		payload,
		parseRawPayloadBytes,
		func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
		v.Error,
		func(slideIndex int) (any, error) {
			notes, err := e.GetNotes(slideIndex)
			if err != nil {
				return nil, err
			}
			hasNotesSlide, err := e.HasNotesSlide(slideIndex)
			if err != nil {
				return nil, err
			}
			rawPlaceholders, err := e.ListNotesPlaceholders(slideIndex)
			if err != nil {
				return nil, err
			}
			rawShapes, err := e.ListNotesShapes(slideIndex)
			if err != nil {
				return nil, err
			}
			placeholders := make([]common.PlaceholderInfo, 0, len(rawPlaceholders))
			for _, ph := range rawPlaceholders {
				placeholders = append(placeholders, common.PlaceholderInfo{
					Type:  ph.Type,
					Index: ph.Index,
					Name:  ph.Name,
				})
			}
			return editorcommand.BuildNotesResultDetailed(
				notes,
				hasNotesSlide,
				placeholders,
				rawShapes,
			), nil
		},
	)
}

func handleNotesSlideExists(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleSlideIndexRequest(
		payload,
		parseRawPayloadBytes,
		func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
		v.Error,
		func(slideIndex int) (any, error) {
			hasNotesSlide, err := e.HasNotesSlide(slideIndex)
			if err != nil {
				return nil, err
			}
			return map[string]bool{"notes_slide_exists": hasNotesSlide}, nil
		},
	)
}

func handleSetNotes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, text, ok := editorcommand.ParseSetNotesRequest(
		p,
		func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
		v.RequireString,
	)
	if !ok {
		return nil, v.Error()
	}
	if err := e.SetNotes(slideIndex, text); err != nil {
		return nil, err
	}
	return map[string]bool{"updated": true}, nil
}

func handleSetNotesShapeText(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.SetNotesShapeTextRequest, bool) {
			return editorcommand.ParseSetNotesShapeTextRequest(
				p,
				func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
				v.RequireInt,
				v.RequireString,
			)
		},
		v.Error,
		func(request editorcommand.SetNotesShapeTextRequest) (any, error) {
			if err := e.SetNotesShapeText(request.SlideIndex, request.ShapeID, request.Text); err != nil {
				return nil, err
			}
			return map[string]bool{"updated": true}, nil
		},
	)
}

func handleSetNotesShapeProps(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleSlideShapeRequestWithPayload(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (int, bool) { return requireSlideIndex(e, p, v) },
		v.RequireInt,
		v.Error,
		func(request editorcommand.SlideShapeRequest, p map[string]any) (any, error) {
			var updates common.ShapeUpdate
			if err := editorcommand.DecodeOptionalPayloadValue(p, "updates", &updates); err != nil {
				return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
			}
			if err := e.SetNotesShapeProperties(request.SlideIndex, request.ShapeID, updates); err != nil {
				return nil, err
			}
			return map[string]bool{"updated": true}, nil
		},
	)
}

func handleGetAuthors(e *PresentationEditor, _ json.RawMessage) (any, error) {
	authors, err := e.GetAuthors()
	if err != nil {
		return nil, err
	}
	return map[string]any{"authors": authors}, nil
}

func handleAddAuthor(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.AuthorAddRequest, bool) {
			return editorcommand.ParseAuthorAddRequest(p, v.RequireString)
		},
		v.Error,
		func(request editorcommand.AuthorAddRequest) (any, error) {
			author, err := e.AddAuthor(request.Name, request.Initials)
			if err != nil {
				return nil, err
			}
			return map[string]int64{"author_id": author.ID}, nil
		},
	)
}

func handleGetComments(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleSlideIndexRequest(
		payload,
		parseRawPayloadBytes,
		func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
		v.Error,
		func(slideIndex int) (any, error) {
			comments, err := e.GetComments(slideIndex)
			if err != nil {
				return nil, err
			}
			return map[string]any{"comments": comments}, nil
		},
	)
}

func handleAddComment(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.CommentAddRequest, bool) {
			return editorcommand.ParseCommentAddRequest(
				p,
				func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
				v.RequireInt64,
				v.RequireString,
			)
		},
		v.Error,
		func(request editorcommand.CommentAddRequest) (any, error) {
			if err := e.AddComment(
				request.SlideIndex,
				request.AuthorID,
				request.Text,
				request.X,
				request.Y,
			); err != nil {
				return nil, err
			}
			return map[string]bool{"added": true}, nil
		},
	)
}

func handleRemoveComment(e *PresentationEditor, payload json.RawMessage) (any, error) {
	v := NewPayloadValidator()
	return editorcommand.HandleParsedRequest(
		payload,
		parseRawPayloadBytes,
		func(p map[string]any) (editorcommand.CommentRemoveRequest, bool) {
			return editorcommand.ParseCommentRemoveRequest(
				p,
				func(payload map[string]any) (int, bool) { return requireSlideIndex(e, payload, v) },
				v.RequireInt64,
				v.RequireInt,
			)
		},
		v.Error,
		func(request editorcommand.CommentRemoveRequest) (any, error) {
			if err := e.RemoveComment(request.SlideIndex, request.AuthorID, request.AuthorIndex); err != nil {
				return nil, err
			}
			return map[string]bool{"removed": true}, nil
		},
	)
}
