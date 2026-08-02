package editor

import (
	"encoding/json"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// handleAddChartUserShapes attaches shapes to a chart's drawing part.
//
// Payload:
//
//	{
//	  "slide_index": N,
//	  "chart_selector": {"index": 0},
//	  "shapes": [
//	    {"text": "Q4 spike", "from_x": 0.55, "from_y": 0.05,
//	     "to_x": 0.95, "to_y": 0.18, "font_size_pt": 12, "bold": true}
//	  ]
//	}
//
// Response: {"drawing_part": "ppt/drawings/drawing1.xml", "shape_count": N}.
func handleAddChartUserShapes(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}

	params := struct {
		ChartSelector common.ChartSelector `json:"chart_selector"`
		Shapes        []struct {
			Text       string  `json:"text"`
			FromX      float64 `json:"from_x"`
			FromY      float64 `json:"from_y"`
			ToX        float64 `json:"to_x"`
			ToY        float64 `json:"to_y"`
			FontSizePt float64 `json:"font_size_pt"`
			Bold       bool    `json:"bold"`
			Name       string  `json:"name"`
		} `json:"shapes"`
	}{}
	if err := json.Unmarshal(payload, &params); err != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}

	userShapes := make([]ChartUserShape, 0, len(params.Shapes))
	for _, shape := range params.Shapes {
		userShapes = append(userShapes, ChartUserShape{
			Text:       shape.Text,
			FromX:      shape.FromX,
			FromY:      shape.FromY,
			ToX:        shape.ToX,
			ToY:        shape.ToY,
			FontSizePt: shape.FontSizePt,
			Bold:       shape.Bold,
			Name:       shape.Name,
		})
	}

	drawingPart, addErr := e.AddChartUserShapes(slideIndex, params.ChartSelector, userShapes)
	if addErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, addErr.Error())
	}
	return map[string]any{
		"drawing_part": drawingPart,
		"shape_count":  len(userShapes),
	}, nil
}
