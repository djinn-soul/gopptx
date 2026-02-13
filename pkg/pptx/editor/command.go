package editor

import (
	"encoding/json"
	"fmt"

	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// RequestEnvelope is the standard wrapper for all incoming commands.
type RequestEnvelope struct {
	APIVersion int             `json:"api_version"`
	Op         string          `json:"op"`
	Payload    json.RawMessage `json:"payload"`
	RequestID  string          `json:"request_id"`
}

// ResponseEnvelope is the standard wrapper for all outgoing responses.
type ResponseEnvelope struct {
	OK        bool         `json:"ok"`
	Result    interface{}  `json:"result,omitempty"`
	Error     *ErrorDetail `json:"error,omitempty"`
	RequestID string       `json:"request_id,omitempty"`
}

type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type commandHandler func(*PresentationEditor, json.RawMessage) (interface{}, error)

var commandHandlers = map[string]commandHandler{
	OpSlideCount:              handleSlideCount,
	OpAddSlide:                handleAddSlide,
	OpRemoveSlide:             handleRemoveSlide,
	OpMoveSlide:               handleMoveSlide,
	OpDuplicateSlide:          handleDuplicateSlide,
	OpGetMetadata:             handleGetMetadata,
	OpUpdateChartData:         handleUpdateChartData,
	OpListSlideCharts:         handleListSlideCharts,
	OpListSlideLayouts:        handleListSlideLayouts,
	OpRebindSlideLayout:       handleRebindSlideLayout,
	OpCloneLayoutMasterFamily: handleCloneLayoutMasterFamily,
}

// ExecuteCommand dispatches a JSON command to the appropriate editor method.
func ExecuteCommand(e *PresentationEditor, jsonInput string) string {
	var req RequestEnvelope
	if err := json.Unmarshal([]byte(jsonInput), &req); err != nil {
		return errorResponse("INVALID_JSON", err.Error(), "")
	}

	if req.APIVersion != 1 {
		return errorResponse("UNSUPPORTED_VERSION", fmt.Sprintf("API version %d not supported", req.APIVersion), req.RequestID)
	}

	handler, ok := commandHandlers[req.Op]
	if !ok {
		return errorResponse("UNKNOWN_OP", fmt.Sprintf("Operation %q not recognized", req.Op), req.RequestID)
	}

	result, err := handler(e, req.Payload)
	if err != nil {
		return errorResponse("OP_FAILED", err.Error(), req.RequestID)
	}

	resp := ResponseEnvelope{
		OK:        true,
		Result:    result,
		RequestID: req.RequestID,
	}
	out, marshalErr := json.Marshal(resp)
	if marshalErr != nil {
		return errorResponse("MARSHAL_ERROR", marshalErr.Error(), req.RequestID)
	}
	return string(out)
}

func handleSlideCount(e *PresentationEditor, _ json.RawMessage) (interface{}, error) {
	return map[string]int{"count": e.SlideCount()}, nil
}

func handleAddSlide(e *PresentationEditor, payload json.RawMessage) (interface{}, error) {
	var p struct {
		Title string `json:"title"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	index, err := e.AddSlide(elements.NewSlide(p.Title))
	if err != nil {
		return nil, err
	}
	return map[string]int{"index": index}, nil
}

func handleRemoveSlide(e *PresentationEditor, payload json.RawMessage) (interface{}, error) {
	var p struct {
		Index int `json:"index"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	return nil, e.RemoveSlide(p.Index)
}

func handleMoveSlide(e *PresentationEditor, payload json.RawMessage) (interface{}, error) {
	var p struct {
		From int `json:"from"`
		To   int `json:"to"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	return nil, e.MoveSlide(p.From, p.To)
}

func handleDuplicateSlide(e *PresentationEditor, payload json.RawMessage) (interface{}, error) {
	var p struct {
		Index    int `json:"index"`
		InsertAt int `json:"insert_at"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	newIdx, err := e.DuplicateSlide(p.Index, p.InsertAt)
	if err != nil {
		return nil, err
	}
	return map[string]int{"new_index": newIdx}, nil
}

func handleGetMetadata(e *PresentationEditor, _ json.RawMessage) (interface{}, error) {
	m := e.Metadata()
	return map[string]interface{}{
		"title":       m.Title,
		"slide_count": m.SlideCount,
		"size": map[string]int64{
			"width":  m.SlideSize.Width,
			"height": m.SlideSize.Height,
		},
	}, nil
}

func handleUpdateChartData(e *PresentationEditor, payload json.RawMessage) (interface{}, error) {
	var p struct {
		SlideIndex    int                    `json:"slide_index"`
		ChartSelector common.ChartSelector   `json:"chart_selector"`
		Data          common.ChartDataUpdate `json:"data"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	if err := e.UpdateChartData(p.SlideIndex, p.ChartSelector, p.Data); err != nil {
		return nil, err
	}
	return map[string]bool{"updated": true}, nil
}

func handleListSlideCharts(e *PresentationEditor, payload json.RawMessage) (interface{}, error) {
	var p struct {
		SlideIndex int `json:"slide_index"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	refs, err := e.ListSlideCharts(p.SlideIndex)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"charts": refs}, nil
}

func handleListSlideLayouts(e *PresentationEditor, _ json.RawMessage) (interface{}, error) {
	layouts, err := e.ListSlideLayouts()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"layouts": layouts}, nil
}

func handleRebindSlideLayout(e *PresentationEditor, payload json.RawMessage) (interface{}, error) {
	var p struct {
		SlideIndex int    `json:"slide_index"`
		LayoutPart string `json:"layout_part"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	if err := e.RebindSlideLayout(p.SlideIndex, p.LayoutPart); err != nil {
		return nil, err
	}
	return map[string]bool{"rebound": true}, nil
}

func handleCloneLayoutMasterFamily(e *PresentationEditor, payload json.RawMessage) (interface{}, error) {
	var p struct {
		LayoutPart string `json:"layout_part"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	result, err := e.CloneLayoutMasterFamily(p.LayoutPart)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"master_part": result.MasterPart,
		"theme_part":  result.ThemePart,
		"layout_map":  result.LayoutMap,
	}, nil
}

func errorResponse(code, message, reqID string) string {
	resp := ResponseEnvelope{
		OK:        false,
		RequestID: reqID,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
	out, err := json.Marshal(resp)
	if err != nil {
		// Fallback for catastrophic failure
		return `{"ok": false, "error": {"code": "INTERNAL_ERROR", "message": "Failed to marshal error response"}}`
	}
	return string(out)
}
