package editor

import (
	"encoding/json"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	slidescatalog "github.com/djinn-soul/gopptx/pkg/pptx/editor/handlers/slidescatalog"
)

// handleUpdateChartCachedValues refreshes the numbers a chart displays without
// touching its workbook link.
//
// Payload: {"slide_index": N, "chart_selector": {...}, "data": {...}}.
// Response: {"updated": true}.
func handleUpdateChartCachedValues(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}

	var params struct {
		ChartSelector common.ChartSelector   `json:"chart_selector"`
		Data          common.ChartDataUpdate `json:"data"`
	}
	if err := json.Unmarshal(payload, &params); err != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}

	if err := e.UpdateChartCachedValues(slideIndex, params.ChartSelector, params.Data); err != nil {
		return nil, err
	}
	return slidescatalog.BuildUpdatedResponse(), nil
}

// handleGetChartDataSource reports whether a chart's data is embedded, linked
// to an external workbook, or absent.
//
// Payload: {"slide_index": N, "chart_selector": {...}}.
// Response: {"source": {...}}.
func handleGetChartDataSource(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}

	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}

	var params struct {
		ChartSelector common.ChartSelector `json:"chart_selector"`
	}
	if err := json.Unmarshal(payload, &params); err != nil {
		return nil, NewBridgeError(ErrCodeInvalidPayload, err.Error())
	}

	source, err := e.GetChartDataSource(slideIndex, params.ChartSelector)
	if err != nil {
		return nil, err
	}
	return map[string]any{"source": source}, nil
}
