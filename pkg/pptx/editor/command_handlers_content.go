package editor

import (
	"encoding/json"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func handleListSlides(e *PresentationEditor, _ json.RawMessage) (any, error) {
	return map[string]any{"slides": e.Slides()}, nil
}

func handleFindAndReplace(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		Find    string `json:"find"`
		Replace string `json:"replace"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	count, err := e.FindAndReplaceInShapes(p.Find, p.Replace)
	if err != nil {
		return nil, err
	}
	return map[string]int{"replacements": count}, nil
}

func handleSearchShapes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var query common.ShapeSearchQuery
	if err := json.Unmarshal(payload, &query); err != nil {
		return nil, err
	}
	results, err := e.SearchShapes(query)
	if err != nil {
		return nil, err
	}
	return map[string]any{"results": results}, nil
}

func handleGetAuthors(e *PresentationEditor, _ json.RawMessage) (any, error) {
	authors, err := e.GetAuthors()
	if err != nil {
		return nil, err
	}
	return map[string]any{"authors": authors}, nil
}

func handleAddAuthor(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		Name     string `json:"name"`
		Initials string `json:"initials"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	author, err := e.AddAuthor(p.Name, p.Initials)
	if err != nil {
		return nil, err
	}
	return map[string]int64{"author_id": author.ID}, nil
}

func handleGetComments(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		SlideIndex int `json:"slide_index"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	comments, err := e.GetComments(p.SlideIndex)
	if err != nil {
		return nil, err
	}
	return map[string]any{"comments": comments}, nil
}

func handleAddComment(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		SlideIndex int    `json:"slide_index"`
		AuthorID   int64  `json:"author_id"`
		Text       string `json:"text"`
		X          int64  `json:"x"`
		Y          int64  `json:"y"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	if err := e.AddComment(p.SlideIndex, p.AuthorID, p.Text, p.X, p.Y); err != nil {
		return nil, err
	}
	return map[string]bool{"added": true}, nil
}

func handleRemoveComment(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		SlideIndex  int   `json:"slide_index"`
		AuthorID    int64 `json:"author_id"`
		AuthorIndex int   `json:"author_index"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	if err := e.RemoveComment(p.SlideIndex, p.AuthorID, p.AuthorIndex); err != nil {
		return nil, err
	}
	return map[string]bool{"removed": true}, nil
}

func handleSetModifyPassword(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		Password string `json:"password"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	e.Metadata().Protection.ModifyPassword = p.Password
	return map[string]bool{"updated": true}, nil
}

func handleSetMarkAsFinal(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		Final bool `json:"final"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	e.Metadata().Protection.MarkAsFinal = p.Final
	return map[string]bool{"updated": true}, nil
}
