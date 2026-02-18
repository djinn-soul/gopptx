package editor

import (
	"encoding/json"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func handleListShapes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		SlideIndex int `json:"slide_index"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	shapes, err := e.GetShapes(p.SlideIndex)
	if err != nil {
		return nil, err
	}
	return map[string]any{"shapes": shapes}, nil
}

func handleAddShape(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		SlideIndex int                `json:"slide_index"`
		Type       string             `json:"type"`
		X          float64            `json:"x"`
		Y          float64            `json:"y"`
		W          float64            `json:"w"`
		H          float64            `json:"h"`
		Text       string             `json:"text"`
		Properties *common.ShapeProps `json:"properties"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}

	id, err := e.AddShape(p.SlideIndex, p.Type, p.X, p.Y, p.W, p.H)
	if err != nil {
		return nil, err
	}

	if p.Text != "" {
		updates := common.ShapeUpdate{Text: &p.Text}
		if updateErr := e.UpdateShape(p.SlideIndex, id, updates); updateErr != nil {
			return nil, updateErr
		}
	}
	return map[string]int{"shape_id": id}, nil
}

func handleAddImage(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		SlideIndex int     `json:"slide_index"`
		Path       string  `json:"path"`
		X          float64 `json:"x"`
		Y          float64 `json:"y"`
		W          float64 `json:"w"`
		H          float64 `json:"h"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	id, err := e.AddImage(p.SlideIndex, p.Path, p.X, p.Y, p.W, p.H)
	if err != nil {
		return nil, err
	}
	return map[string]int{"shape_id": id}, nil
}

func handleRemoveShape(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		SlideIndex int `json:"slide_index"`
		ShapeID    int `json:"shape_id"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	if err := e.RemoveShape(p.SlideIndex, p.ShapeID); err != nil {
		return nil, err
	}
	return map[string]bool{"removed": true}, nil
}

func handleUpdateShape(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		SlideIndex int                `json:"slide_index"`
		ShapeID    int                `json:"shape_id"`
		Updates    common.ShapeUpdate `json:"updates"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	if err := e.UpdateShape(p.SlideIndex, p.ShapeID, p.Updates); err != nil {
		return nil, err
	}
	return map[string]bool{"updated": true}, nil
}

func handleGetNotes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		SlideIndex int `json:"slide_index"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	notes, err := e.GetNotes(p.SlideIndex)
	if err != nil {
		return nil, err
	}
	return map[string]string{"text": notes}, nil
}

func handleSetNotes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		SlideIndex int    `json:"slide_index"`
		Text       string `json:"text"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	if err := e.SetNotes(p.SlideIndex, p.Text); err != nil {
		return nil, err
	}
	return map[string]bool{"updated": true}, nil
}
