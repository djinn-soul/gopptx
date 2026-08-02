package editor

import (
	"encoding/json"
)

// handleFitShapeText shrinks a shape's text to the largest size that fits.
//
// Payload:
//
//	{
//	  "slide_index": N,
//	  "shape_id": N,
//	  "font_path": "C:/Windows/Fonts/calibri.ttf", // optional; without it only
//	                                               // the autofit flags are set
//	  "max_size_pt": 18,   // optional
//	  "min_size_pt": 8,    // optional
//	  "word_wrap": true    // optional, defaults to true
//	}
//
// Response: {"font_size_pt": F, "fits": B, "line_count": N, "measured": B}.
func handleFitShapeText(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	shapeID, _ := v.RequireInt(p, "shape_id")
	if v.HasErrors() {
		return nil, v.Error()
	}

	params := struct {
		FontPath  string  `json:"font_path"`
		MaxSizePt float64 `json:"max_size_pt"`
		MinSizePt float64 `json:"min_size_pt"`
		WordWrap  *bool   `json:"word_wrap"`
	}{}
	if err := json.Unmarshal(payload, &params); err != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}

	wordWrap := true
	if params.WordWrap != nil {
		wordWrap = *params.WordWrap
	}

	result, fitErr := e.FitShapeText(slideIndex, shapeID, FitTextRequest{
		FontPath:  params.FontPath,
		MaxSizePt: params.MaxSizePt,
		MinSizePt: params.MinSizePt,
		WordWrap:  wordWrap,
	})
	if fitErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, fitErr.Error())
	}
	return map[string]any{
		"font_size_pt": result.FontSizePt,
		"fits":         result.Fits,
		"line_count":   result.LineCount,
		"measured":     result.Measured,
	}, nil
}

// handleFitShapeToText grows a shape until its text fits inside it.
//
// Payload:
//
//	{
//	  "slide_index": N,
//	  "shape_id": N,
//	  "font_path": "C:/Windows/Fonts/calibri.ttf", // required, to measure with
//	  "font_size_pt": 18,      // optional; defaults to the runs' largest size
//	  "min_height_emu": 0,     // optional
//	  "max_height_emu": 0,     // optional
//	  "shrink": false          // optional; allow the shape to get shorter
//	}
//
// Response: {"height_emu": N, "line_count": N, "resized": B}.
func handleFitShapeToText(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	shapeID, _ := v.RequireInt(p, "shape_id")
	if v.HasErrors() {
		return nil, v.Error()
	}

	params := struct {
		FontPath     string  `json:"font_path"`
		FontSizePt   float64 `json:"font_size_pt"`
		MinHeightEmu int     `json:"min_height_emu"`
		MaxHeightEmu int     `json:"max_height_emu"`
		Shrink       bool    `json:"shrink"`
	}{}
	if err := json.Unmarshal(payload, &params); err != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}

	result, fitErr := e.FitShapeToText(slideIndex, shapeID, FitShapeToTextRequest{
		FontPath:     params.FontPath,
		FontSizePt:   params.FontSizePt,
		MinHeightEmu: params.MinHeightEmu,
		MaxHeightEmu: params.MaxHeightEmu,
		Shrink:       params.Shrink,
	})
	if fitErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, fitErr.Error())
	}
	return map[string]any{
		"height_emu": result.HeightEmu,
		"line_count": result.LineCount,
		"resized":    result.Resized,
	}, nil
}
