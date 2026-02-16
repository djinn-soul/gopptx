package editor

import (
	"encoding/json"
	"fmt"

	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
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
	Result    any          `json:"result,omitempty"`
	Error     *ErrorDetail `json:"error,omitempty"`
	RequestID string       `json:"request_id,omitempty"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type commandHandler func(*PresentationEditor, json.RawMessage) (any, error)

func commandHandlerFor(op string) (commandHandler, bool) {
	if handler, ok := commandHandlerForSlidesAndMeta(op); ok {
		return handler, true
	}
	return commandHandlerForShapesAndNotes(op)
}

func commandHandlerForSlidesAndMeta(op string) (commandHandler, bool) {
	switch op {
	case OpSlideCount:
		return handleSlideCount, true
	case OpAddSlide:
		return handleAddSlide, true
	case OpRemoveSlide:
		return handleRemoveSlide, true
	case OpMoveSlide:
		return handleMoveSlide, true
	case OpDuplicateSlide:
		return handleDuplicateSlide, true
	case OpGetMetadata:
		return handleGetMetadata, true
	case OpUpdateChartData:
		return handleUpdateChartData, true
	case OpListSlideCharts:
		return handleListSlideCharts, true
	case OpListSlideLayouts:
		return handleListSlideLayouts, true
	case OpRebindSlideLayout:
		return handleRebindSlideLayout, true
	case OpCloneLayoutMasterFamily:
		return handleCloneLayoutMasterFamily, true
	case OpAddSection:
		return handleAddSection, true
	case OpRemoveSection:
		return handleRemoveSection, true
	case OpRenameSection:
		return handleRenameSection, true
	case OpGetCoreProperties:
		return handleGetCoreProperties, true
	case OpSetCoreProperties:
		return handleSetCoreProperties, true
	case OpApplyTheme:
		return handleApplyTheme, true
	case OpSetSlideSize:
		return handleSetSlideSize, true
	default:
		return nil, false
	}
}

func commandHandlerForShapesAndNotes(op string) (commandHandler, bool) {
	switch op {
	case OpListShapes:
		return handleListShapes, true
	case OpAddShape:
		return handleAddShape, true
	case OpRemoveShape:
		return handleRemoveShape, true
	case OpUpdateShape:
		return handleUpdateShape, true
	case OpGetNotes:
		return handleGetNotes, true
	case OpSetNotes:
		return handleSetNotes, true
	default:
		return nil, false
	}
}

// ExecuteCommand dispatches a JSON command to the appropriate editor method.
func ExecuteCommand(e *PresentationEditor, jsonInput string) string {
	var req RequestEnvelope
	if err := json.Unmarshal([]byte(jsonInput), &req); err != nil {
		return errorResponse("INVALID_JSON", err.Error(), "")
	}

	if req.APIVersion != 1 {
		return errorResponse(
			"UNSUPPORTED_VERSION",
			fmt.Sprintf("API version %d not supported", req.APIVersion),
			req.RequestID,
		)
	}

	handler, ok := commandHandlerFor(req.Op)
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

func handleSlideCount(e *PresentationEditor, _ json.RawMessage) (any, error) {
	return map[string]int{"count": e.SlideCount()}, nil
}

func handleAddSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
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

func handleRemoveSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		Index int `json:"index"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	return nil, e.RemoveSlide(p.Index)
}

func handleMoveSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		From int `json:"from"`
		To   int `json:"to"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	return nil, e.MoveSlide(p.From, p.To)
}

func handleDuplicateSlide(e *PresentationEditor, payload json.RawMessage) (any, error) {
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

func handleUpdateChartData(e *PresentationEditor, payload json.RawMessage) (any, error) {
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

func handleListSlideCharts(e *PresentationEditor, payload json.RawMessage) (any, error) {
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
	return map[string]any{"charts": refs}, nil
}

func handleListSlideLayouts(e *PresentationEditor, _ json.RawMessage) (any, error) {
	layouts, err := e.ListSlideLayouts()
	if err != nil {
		return nil, err
	}
	return map[string]any{"layouts": layouts}, nil
}

func handleRebindSlideLayout(e *PresentationEditor, payload json.RawMessage) (any, error) {
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

func handleCloneLayoutMasterFamily(e *PresentationEditor, payload json.RawMessage) (any, error) {
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
	return map[string]any{
		"master_part": result.MasterPart,
		"theme_part":  result.ThemePart,
		"layout_map":  result.LayoutMap,
	}, nil
}

func handleAddSection(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		Name         string `json:"name"`
		SlideIndices []int  `json:"slide_indices"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	if err := e.AddSection(p.Name, p.SlideIndices); err != nil {
		return nil, err
	}
	return map[string]bool{"added": true}, nil
}

func handleRemoveSection(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	if err := e.RemoveSection(p.Name); err != nil {
		return nil, err
	}
	return map[string]bool{"removed": true}, nil
}

func handleRenameSection(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		OldName string `json:"old_name"`
		NewName string `json:"new_name"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	if err := e.RenameSection(p.OldName, p.NewName); err != nil {
		return nil, err
	}
	return map[string]bool{"renamed": true}, nil
}

func handleGetCoreProperties(e *PresentationEditor, _ json.RawMessage) (any, error) {
	return e.GetCoreProperties(), nil
}

func handleSetCoreProperties(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p common.CoreProperties
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	e.SetCoreProperties(p)
	return map[string]bool{"updated": true}, nil
}

func handleApplyTheme(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p struct {
		ThemeName string `json:"theme_name"`
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}

	var theme styling.Theme
	switch p.ThemeName {
	case "Corporate":
		theme = styling.ThemeCorporate
	case "Modern":
		theme = styling.ThemeModern
	case "Vibrant":
		theme = styling.ThemeVibrant
	case "Dark":
		theme = styling.ThemeDark
	case "Nature":
		theme = styling.ThemeNature
	case "Tech":
		theme = styling.ThemeTech
	case "Carbon":
		theme = styling.ThemeCarbon
	default:
		return nil, fmt.Errorf("unknown theme name %q", p.ThemeName)
	}

	if err := e.ApplyTheme(theme); err != nil {
		return nil, err
	}
	return map[string]bool{"applied": true}, nil
}

func handleSetSlideSize(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var p common.SlideSize
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	if err := e.SetSlideSize(p); err != nil {
		return nil, err
	}
	return map[string]bool{"updated": true}, nil
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
		Text       string             `json:"text"`       // Optional
		Properties *common.ShapeProps `json:"properties"` // Optional
	}
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}

	// Assuming PresentationEditor has AddShape(slideIndex, type, x, y, w, h)
	id, err := e.AddShape(p.SlideIndex, p.Type, p.X, p.Y, p.W, p.H)
	if err != nil {
		return nil, err
	}

	// Apply optional updates
	if p.Text != "" {
		updates := common.ShapeUpdate{
			Text: &p.Text,
		}
		if updateErr := e.UpdateShape(p.SlideIndex, id, updates); updateErr != nil {
			return nil, updateErr
		}
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
	// Assuming GetNotes returns just the text string for Phase 1
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
