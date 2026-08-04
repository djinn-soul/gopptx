package editor

import (
	"encoding/json"
)

// handleAddEquation inserts a native OMML equation as a text box.
//
// Payload:
//
//	{
//	  "slide_index": N,
//	  "latex": "\\frac{a}{b}",
//	  "x": 0, "y": 0, "w": 0, "h": 0,   // EMUs
//	  "font_size_pt": 24                 // optional
//	}
//
// Response: {"shape_id": N}.
func handleAddEquation(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	latex, _ := v.RequireString(p, "latex")
	if v.HasErrors() {
		return nil, v.Error()
	}

	params := struct {
		X          int     `json:"x"`
		Y          int     `json:"y"`
		W          int     `json:"w"`
		H          int     `json:"h"`
		FontSizePt float64 `json:"font_size_pt"`
	}{}
	if err := json.Unmarshal(payload, &params); err != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}

	shapeID, addErr := e.AddEquation(slideIndex, EquationRequest{
		LaTeX:      latex,
		X:          params.X,
		Y:          params.Y,
		W:          params.W,
		H:          params.H,
		FontSizePt: params.FontSizePt,
	})
	if addErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, addErr.Error())
	}
	return respShapeID(shapeID), nil
}
